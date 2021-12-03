package bot

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/dreamteam-hack/hack041221/telegram-bot/pkg/types"
	tb "gopkg.in/tucnak/telebot.v2"
)

func NewBot(token, jobUri, stateUri, youtubeUri string, q *queue, u *uploader) *bot {
	return &bot{
		token:      token,
		q:          q,
		u:          u,
		stateUri:   stateUri,
		jobUri:     jobUri,
		youtubeUri: youtubeUri,
	}
}

type bot struct {
	stateUri   string
	jobUri     string
	youtubeUri string
	token      string
	q          *queue
	u          *uploader
	tbot       *tb.Bot
}

func (b *bot) handleState(msg []byte) error {
	s := &types.StateMessage{}
	if err := json.Unmarshal(msg, s); err != nil {
		log.Error().Err(err).Msg("json.Unmarshal")
		return err
	}

	log.Info().Msgf("handleState: %v", s)
	replyRecipient := &tb.Message{
		ID:   s.MessageID,
		Chat: &tb.Chat{ID: s.ChatID},
	}

	log.Debug().Msgf("handleState = %s", s)
	var replyMsg string
	if len(s.Error) > 0 {
		replyMsg = fmt.Sprintf("Произошла ошибка при обработке видео: %s", s.Error)
		log.Info().Msgf(replyMsg)
	} else {
		msgBody := ""
		for _, s := range s.Result.Summary {
			msgBody += fmt.Sprintf("Ratio: %s\n%s\n\n", s.Ratio, s.Desc)
		}
		msgBody += "--------\n"
		msgBody += fmt.Sprintf("LOC: %s\n\n", strings.Join(s.Result.Ner.LOC, ", "))
		msgBody += fmt.Sprintf("PER: %s\n\n", strings.Join(s.Result.Ner.PER, ", "))
		msgBody += fmt.Sprintf("ORG: %s\n\n", strings.Join(s.Result.Ner.ORG, ", "))

		replyMsg = fmt.Sprintf(fmt.Sprintf("Видео готово:\n\n%s", msgBody))
	}
	if _, err := b.tbot.Reply(replyRecipient, replyMsg); err != nil {
		log.Error().Err(err).Msgf("failed to send reply %s", s.URL)
	}

	return nil
}

func (b *bot) Listen() {
	go b.q.Listen(b.stateUri, b.handleState)

	b.tbot.Handle(tb.OnText, func(m *tb.Message) {
		b.handleText(m)
	})
	b.tbot.Handle(tb.OnVideo, func(m *tb.Message) {
		b.handleUpload(m, &m.Video.File)
	})
	b.tbot.Handle(tb.OnVideoNote, func(m *tb.Message) {
		b.handleUpload(m, &m.VideoNote.File)
	})

	b.tbot.Handle("/start", func(m *tb.Message) {
		b.tbot.Reply(m, "Привет друг! Я принимаю только видео файлы и видео сообщения, все остальное будет проигнорировано.")
	})

	b.tbot.Start()
}

func (b *bot) Init() error {
	tbot, err := tb.NewBot(tb.Settings{
		Token:  b.token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize telegram bot")
		return err
	}
	b.tbot = tbot
	return nil
}

func (b *bot) handleText(m *tb.Message) {
	// 	for _, l := range hasYoutubeLink(m.Text) {
	// 		msg := fmt.Sprintf("Обнаружено youtube видео (%s), начата обработка", l)
	// 		if _, err := b.tbot.Reply(m, msg); err != nil {
	// 			log.Error().Err(err).Msg("Не удалось отправить уведомление")
	// 		}

	// 		payload := &youtubePayload{
	// 			ChatID:    m.Chat.ID,
	// 			MessageID: m.ID,
	// 			VideoURL:  l,
	// 		}
	// 		if err := b.q.Send(b.jobUri, payload); err != nil {
	// 			log.Error().Err(err).Msg("failed to upload video file")
	// 			return
	// 		}
	// 	}

	payload := &types.JobMessage{
		ChatID:    m.Chat.ID,
		MessageID: m.ID,
		URL:       m.Text,
	}
	if err := b.q.Send(b.jobUri, payload); err != nil {
		log.Error().Err(err).Msg("failed to upload video file")
		return
	}
}

func (b *bot) handleUpload(m *tb.Message, f *tb.File) {
	filePath := fmt.Sprintf("/%s", f.UniqueID)
	reader, err := b.tbot.GetFile(f)
	if err != nil {
		log.Error().Err(err).Msg("failed to get file from telegram server")
		return
	}

	log.Debug().Msgf("uploaded: %s", filePath)
	if err := b.u.Upload(filePath, reader); err != nil {
		log.Error().Err(err).Msg("failed to upload video file")
		return
	}

	if _, err := b.tbot.Reply(m, "Видео получено, начата обработка"); err != nil {
		log.Error().Err(err).Msg("Не удалось отправить уведомление")
	}

	payload := &types.JobMessage{
		ChatID:    m.Chat.ID,
		MessageID: m.ID,
		URL:       filePath,
	}
	if err := b.q.Send(b.jobUri, payload); err != nil {
		log.Error().Err(err).Msg("failed to upload video file")
		return
	}
}
