package main

type cfg struct {
	BotToken     string `env:"TELEGRAM_TOKEN"`
	AwsSecretKey string `env:"AWS_SECRET_KEY"`
	AwsAccessKey string `env:"AWS_ACCESS_KEY"`
	AwsRegion    string `env:"AWS_REGION"`
	AwsBucket    string `env:"AWS_BUCKET"`
	JobURL       string `env:"JOB_URL"`
	StateURL     string `env:"STATE_URL"`
	YoutubeURL   string `env:"YOUTUBE_URL"`
}