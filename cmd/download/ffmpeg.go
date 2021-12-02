package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// ffmpeg -i bup.webm -map 0:a -acodec pcm_s16le -ac 1 -ar 8000 out.wav
func ffmpegExtractAudio(src, dst string) error {
	out, err := ffmpeg(fmt.Sprintf("-i %s -map 0:a -acodec pcm_s16le -ac 1 -ar 8000 %s", src, dst))
	if err != nil {
		return err
	}
	l.Debug().Msgf("%s", out)
	return nil
}

// ffmpeg -i bup.webm -pix_fmt bgr8 -r 1/1 ./frames/%5d.jpg
func ffmpegExtractFrames(src, dstDir string) error {
	cmd := fmt.Sprintf("-i %s -pix_fmt bgr8 -r 1/1 %s/", src, dstDir) + "%5d.jpg"
	out, err := ffmpeg(cmd)
	if err != nil {
		l.
			Error().
			Err(err).
			Str("cmd", cmd).
			Msg("ffmpeg extract frame")
		return err
	}
	l.Debug().Msgf("%s", out)
	return nil
}

func ffmpeg(arg string) (out []byte, err error) {
	args := []string{"-nostdin"}
	args = append(args, strings.Split(arg, " ")...)
	l.Debug().Msgf("ffmpeg %s", strings.Join(args, " "))
	return exec.Command("ffmpeg", args...).Output()
}
