package main

func stateErrorReply(v *JobMessage, err error) {
	// Reply in telegram with error
	stateMsg := &StateMessage{
		ChatID:    v.ChatID,
		MessageID: v.MessageID,
		URL:       v.URL,
		Error:     err.Error(),
	}
	if err := stateQueue.Send(stateMsg); err != nil {
		l.Error().Err(err).Msg("change state")
	}
}