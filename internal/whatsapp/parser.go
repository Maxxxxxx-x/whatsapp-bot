package whatsapp

import "go.mau.fi/whatsmeow/proto/waE2E"

func ExtractText(msg *waE2E.Message) string {
	if msg == nil {
		return ""
	}

	if msg.GetConversation() != "" {
		return msg.GetConversation()
	}

	if msg.GetExtendedTextMessage() != nil {
		return msg.GetExtendedTextMessage().GetText()
	}

	return ""
}
