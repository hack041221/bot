package main

import "gitlab.com/dreamteam-hack/hack041221/telegram-bot/pkg/types"

func stateErrorReply(v *types.JobMessage, err error) {
	// Reply in telegram with error
	stateMsg := &types.StateMessage{
		ChatID:    v.ChatID,
		MessageID: v.MessageID,
		URL:       v.URL,
		Error:     err.Error(),
	}
	if err := stateQueue.Send(stateMsg); err != nil {
		l.Error().Err(err).Msg("change state")
	}
}
