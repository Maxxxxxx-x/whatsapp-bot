package whatsapp

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/types"
)

func (s *Service) GetPhoneNumber(ctx context.Context, sender types.JID) (string, error) {
	id, err := s.Client.Store.LIDs.GetPNForLID(context.Background(), sender)
	if err != nil {
		return "", fmt.Errorf("failed to get phone number from JID: %w", err)
	}

	return id.User, nil
}

func (s *Service) GetLocalDisplayName(ctx context.Context, sender types.JID, pushName string) string {
	targetJID := sender.ToNonAD()

	if targetJID.Server == types.HiddenUserServer {
		if pnJID, err := s.Client.Store.LIDs.GetPNForLID(ctx, targetJID); err != nil && !pnJID.IsEmpty() {
			targetJID = pnJID
		}
	}

	contact, err := s.Client.Store.Contacts.GetContact(ctx, targetJID)
	if err == nil && contact.Found {
		if contact.FullName != "" {
			return contact.FullName
		}
		if contact.BusinessName != "" {
			return contact.BusinessName
		}
	}

	if pushName != "" {
		return pushName
	}

	return targetJID.User
}
