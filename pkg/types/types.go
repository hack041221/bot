package types

type StateMessage struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	URL       string `json:"url"`
	Error     string `json:"error"`
	VideoID   string `json:"video_id"`
}

type JobMessage struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	URL       string `json:"url"`
}

type FrameMessage struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	FramesURL string `json:"frames_url"`
	VideoID   string `json:"video_id"`
}

type AudioMessage struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	AudioURL  string `json:"audio_url"`
	VideoID   string `json:"video_id"`
}