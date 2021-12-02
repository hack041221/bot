package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
	"gitlab.com/dreamteam-hack/hack041221/telegram-bot/pkg/bot"
)

var c = new(cfg)

func init() {
	LoadConfig(c)
}

func createSess(c *cfg) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(c.AwsRegion),
		Credentials:      credentials.NewStaticCredentials(c.AwsAccessKey, c.AwsSecretKey, ""),
	})
}

func main() {
	sess, err := createSess(c)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create aws session")
		return
	}

	u := bot.NewUploader(c.AwsBucket, sess)
	q := bot.NewQueue(sqs.New(sess))
	b := bot.NewBot(c.BotToken, c.JobURL, c.StateURL, c.YoutubeURL, q, u)
	if err := b.Init(); err != nil {
		log.Fatal().Err(err).Msg("failed to init telegram bot")
	}
	b.Listen()
}
