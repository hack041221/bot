package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"gitlab.com/dreamteam-hack/hack041221/telegram-bot/pkg/types"
)

func main() {
	ch := jobQueue.Listen()
	go func() {
		for {
			if stop {
				break
			}

			msg, ok := <-ch
			if !ok {
				l.Debug().Msgf("skip %s message", msg)
				continue
			}

			v := &types.JobMessage{}
			if err := json.Unmarshal([]byte(*msg.Body), v); err != nil {
				l.Error().Err(err)
				continue
			}

			if err := process(v, uuid.New()); err != nil {
				l.Error().Err(err).Str("url", v.URL).Msg("process error")
				stateErrorReply(v, err)
			}
			if err := jobQueue.Remove(msg); err != nil {
				l.Error().Err(err).Str("url", v.URL).Msg("remove message from sqs error")
			}
		}
	}()

	done := make(chan struct{})
	go func() {
		l.Info().Msg("Listening signals...")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		stop = true
		close(done)
	}()
	<-done
	l.Info().Msg("Done")
}
