package whatsapp

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func buildContextInfo(evt *events.Message, quote bool) *waE2E.ContextInfo {
	if evt == nil || !quote {
		return nil
	}

	return &waE2E.ContextInfo{
		StanzaID:      proto.String(evt.Info.ID),
		Participant:   proto.String(evt.Info.Sender.ToNonAD().String()),
		QuotedMessage: evt.Message,
	}
}

func (s *Service) ReplyText(ctx context.Context, evt *events.Message, text string, quote bool) error {
	var msg *waE2E.Message

	if contextInfo := buildContextInfo(evt, quote); contextInfo != nil {
		msg = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        proto.String(text),
				ContextInfo: contextInfo,
			},
		}
	} else {
		msg = &waE2E.Message{
			Conversation: proto.String(text),
		}
	}
	_, err := s.Client.SendMessage(ctx, evt.Info.Chat, msg)
	return err
}

func (s *Service) ReplyLocation(ctx context.Context, evt *events.Message, lat, long float64, name, address string, quote bool) error {
	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(lat),
			DegreesLongitude: proto.Float64(long),
			Name:             proto.String(name),
			Address:          proto.String(address),
			ContextInfo:      buildContextInfo(evt, quote),
		},
	}

	_, err := s.Client.SendMessage(ctx, evt.Info.Chat, msg)
	return err
}

func (s *Service) ReplyImage(ctx context.Context, evt *events.Message, imageData []byte, caption string, quote bool) error {
	uploaded, err := s.Client.Upload(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       proto.String(caption),
			Mimetype:      proto.String("image/jpeg"),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uploaded.FileLength),
			ContextInfo:   buildContextInfo(evt, quote),
		},
	}

	_, err = s.Client.SendMessage(ctx, evt.Info.Chat, msg)
	return err
}

func (s *Service) ReplyAudio(ctx context.Context, evt *events.Message, audioData []byte, isVoiceNote bool, quote bool) error {
	uploaded, err := s.Client.Upload(ctx, audioData, whatsmeow.MediaAudio)
	if err != nil {
		return fmt.Errorf("failed to upload audio: %w", err)
	}

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			Mimetype:      proto.String("audio/ogg; codecs=opus"),
			PTT:           proto.Bool(isVoiceNote),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uploaded.FileLength),
			ContextInfo:   buildContextInfo(evt, quote),
		},
	}

	_, err = s.Client.SendMessage(ctx, evt.Info.Chat, msg)
	return err
}

func (s *Service) ReplyDocument(ctx context.Context, evt *events.Message, docData []byte, fileName, mimetype string, quote bool) error {
	uploaded, err := s.Client.Upload(ctx, docData, whatsmeow.MediaDocument)
	if err != nil {
		return fmt.Errorf("failed to upload document: %w", err)
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			Mimetype:      proto.String(mimetype),
			Title:         proto.String(fileName),
			FileName:      proto.String(fileName),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uploaded.FileLength),
			ContextInfo:   buildContextInfo(evt, quote),
		},
	}

	_, err = s.Client.SendMessage(ctx, evt.Info.Chat, msg)
	return err
}
