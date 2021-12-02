package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

func NewSQSQueue(s *sqs.SQS, uri string) *SQSQueue {
	return &SQSQueue{
		s:        s,
		uri:      uri,
		messages: make(chan *sqs.Message),
	}
}

type SQSQueue struct {
	s        *sqs.SQS
	messages chan *sqs.Message
	uri      string
}

func (w *SQSQueue) Send(msg interface{}) (err error) {
	body, err := json.Marshal(&msg)
	if err != nil {
		l.Error().Err(err).Msg("json marshal failed")
		return err
	}

	if err := w.Publish(body); err != nil {
		l.Error().Err(err).Msg("sqs publish failed")
		return err
	}
	return nil
}

func (w *SQSQueue) Publish(msg []byte) (err error) {
	_, err = w.s.SendMessage(&sqs.SendMessageInput{
		MessageBody:            aws.String(string(msg)),
		QueueUrl:               aws.String(w.uri),
		MessageDeduplicationId: aws.String(uuid.New().String()),
		MessageGroupId:         aws.String(uuid.New().String()),
	})
	return
}

func (w *SQSQueue) Listen() <-chan *sqs.Message {
	go func() {
		input := &sqs.ReceiveMessageInput{QueueUrl: &w.uri}

		for {
			output, err := w.s.ReceiveMessage(input)
			if err != nil {
				l.Error().Err(err).Msg("sqs ReceiveMessage error")
				continue
			}

			for _, message := range output.Messages {
				l.Info().Msgf("Received message: %s", message)
				w.messages <- message
			}
		}
	}()

	return w.messages
}

func (w *SQSQueue) Remove(message *sqs.Message) error {
	_, err := w.s.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(w.uri),
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}
