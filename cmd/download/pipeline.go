package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/dreamteam-hack/hack041221/telegram-bot/pkg/types"
)

func process(v *types.JobMessage, uid uuid.UUID) error {
	videoDir, err := ioutil.TempDir(os.TempDir(), "video")
	if err != nil {
		l.
			Error().
			Err(err).
			Msg("ioutil.TempFile")
		return err
	}
	defer os.RemoveAll(videoDir)

	f := fmt.Sprintf("%s/%s", videoDir, filepath.Base(v.URL))
	l.Debug().Str("path", f).Msg("temporary file")

	if err := downloader.download(v.URL, f); err != nil {
		return err
	}

	frameDir, err := ioutil.TempDir(os.TempDir(), "frame")
	if err != nil {
		l.Error().Err(err).Msg("ioutil.TempDir")
		return err
	}
	defer os.RemoveAll(frameDir)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if err := ffmpegExtractFrames(f, frameDir); err != nil {
			log.Error().Err(err).Msg("ffmpegExtractFrames")
			wg.Done()
			return
		}

		uploadDir := fmt.Sprintf("frame/%s", uid.String())
		var objects []s3manager.BatchUploadObject

		files, err := ioutil.ReadDir(frameDir)
		if err != nil {
			log.Error().Err(err).Msg("ioutil.ReadDir")
			wg.Done()
			return
		}

		for _, fi := range files {
			log.Info().Msgf("%s", fi.Name())
			f, err := os.Open(fi.Name())
			if err != nil {
				log.Error().Err(err).Msg("os.Open")
				continue
			}
			uploadKey := fmt.Sprintf("%s/%s", uploadDir, filepath.Base(fi.Name()))
			objects = append(objects, s3manager.BatchUploadObject{
				Object: &s3manager.UploadInput{
					Key:    aws.String(uploadKey),
					Bucket: aws.String(c.AwsBucket),
					Body:   f,
				},
			})
		}

		iter := &s3manager.UploadObjectsIterator{Objects: objects}
		if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
			log.Error().Err(err).Msg("uploader.UploadWithIterator")
			wg.Done()
			return
		}

		frameJob := &types.FrameMessage{
			ChatID:    v.ChatID,
			MessageID: v.MessageID,
			FramesURL: uploadDir, // @todo
			VideoID:   uid.String(),
		}
		if err := frameQueue.Send(frameJob); err != nil {
			log.Error().Err(err).Msg("frameQueue.Send error")
		}

		wg.Done()
	}()

	go func() {
		audioDst := fmt.Sprintf("%s.wav", f)
		if err := ffmpegExtractAudio(f, audioDst); err != nil {
			log.Error().Err(err).Msg("ffmpegExtractAudio")
			wg.Done()
			return
		}

		f, err := os.Open(audioDst)
		if err != nil {
			log.Error().Err(err).Msg("os.Open")
			wg.Done()
			return
		}
		uploadKey := fmt.Sprintf("audio/%s/%s", uid.String(), filepath.Base(audioDst))

		if _, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(c.AwsBucket),
			Key:    aws.String(uploadKey),
			Body:   f,
		}); err != nil {
			log.Error().Err(err).Msg("uploader.Upload")
			wg.Done()
			return
		}

		audioJob := &types.AudioMessage{
			ChatID:    v.ChatID,
			MessageID: v.MessageID,
			AudioURL:  uploadKey,
			VideoID:   uid.String(),
		}
		if err := audioQueue.Send(audioJob); err != nil {
			log.Error().Err(err).Msg("audioQueue.Send error")
			wg.Done()
			return
		}

		wg.Done()
	}()

	wg.Wait()

	return nil
}
