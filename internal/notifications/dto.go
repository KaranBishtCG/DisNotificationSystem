package notifications

import (
	"strings"
	"time"

	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/validation"
)

type CreateNotificationRequest struct {
	UserID    string  `json:"user_id"`
	Channel   Channel `json:"channel"`
	Recipient string  `json:"recipient"`
	Message   string  `json:"message"`
}

type outboxPayload struct {
	NotificationID string `json:"notification_id"`
	UserID         string `json:"user_id"`
	Channel        string `json:"channel"`
	Recipient      string `json:"recipient"`
	Message        string `json:"message"`
}

type CreateNotificationResponse struct {
	NotificationId string    `json:"notification_id"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}

func (req CreateNotificationRequest) Validate() error {
	if strings.TrimSpace(string(req.Channel)) == "" {
		return &validation.ValidationError{
			Field: "Channel",
			Msg:   "Channel is required",
		}
	}

	if strings.TrimSpace(req.Recipient) == "" {
		return &validation.ValidationError{
			Field: "Recipient",
			Msg:   "Recipient is required",
		}
	}

	if strings.TrimSpace(req.Message) == "" {
		return &validation.ValidationError{
			Field: "Message",
			Msg:   "Message is required",
		}
	}

	return nil
}
