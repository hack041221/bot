package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/dreamteam-hack/hack041221/telegram-bot/pkg/config"
	"net/http"
)

var downloader *Downloader
var l zerolog.Logger
var jobQueue *SQSQueue
var stateQueue *SQSQueue
var frameQueue *SQSQueue
var audioQueue *SQSQueue
var c = new(cfg)
var stop = false
var uploader *s3manager.Uploader

func createSess(c *cfg) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(c.AwsRegion),
		Credentials:      credentials.NewStaticCredentials(c.AwsAccessKey, c.AwsSecretKey, ""),
	})
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	l = log.
		With().
		Str("app", "downloader").
		Logger()
	config.LoadConfig(c)
	downloader = NewDownloader(&http.Client{}, 4) // @todo set request timeouts
	sess, err := createSess(c)
	if err != nil {
		l.Fatal().Err(err).Msg("create aws session error")
	}

	uploader = s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024
	})

	s := sqs.New(sess)
	jobQueue = NewSQSQueue(s, c.JobURL)
	stateQueue = NewSQSQueue(s, c.StateURL)
	audioQueue = NewSQSQueue(s, c.AudioURL)
	frameQueue = NewSQSQueue(s, c.FrameURL)
}
