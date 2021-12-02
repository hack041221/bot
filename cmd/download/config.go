package main

type cfg struct {
	AwsSecretKey string `env:"AWS_SECRET_KEY"`
	AwsAccessKey string `env:"AWS_ACCESS_KEY"`
	AwsRegion    string `env:"AWS_REGION"`
	AwsBucket    string `env:"AWS_BUCKET"`
	JobURL       string `env:"JOB_URL"`
	StateURL     string `env:"STATE_URL"`
	FrameURL     string `env:"FRAME_URL"`
	AudioURL     string `env:"AUDIO_URL"`
}
