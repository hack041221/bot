package types

type Summary struct {
	Ratio float64 `json:"ratio"`
	Desc  string  `json:"desc"`
	TS    int     `json:"ts"`
}

type Ner struct {
	LOC []Tag `json:"LOC"`
	PER []Tag `json:"PER"`
	ORG []Tag `json:"ORG"`
}

type Tag struct {
	Tag string  `json:"tag"`
	WT  float64 `json:"wt"`
}

type Result struct {
	Summary []Summary `json:"summary"`
	Ner     Ner       `json:"ner"`
}

type StateMessage struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	URL       string `json:"url"`
	Error     string `json:"error"`
	Result    Result `json:"result"`
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
