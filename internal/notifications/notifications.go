package notifications

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/httpx"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/middleware"
)

type Channel string

// ENUM
const (
	ChannelEmail Channel = "EMAIL"
	ChannelSMS   Channel = "SMS"
	ChannelPush  Channel = "PUSH"
)

func (c Channel) Valid() bool {
	switch c {
	case ChannelEmail, ChannelPush, ChannelSMS:
		return true
	default:
		return false
	}
}

type NotificationHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewNotificationHandler(db *sql.DB, logger *slog.Logger) *NotificationHandler {
	return &NotificationHandler{
		db:     db,
		logger: logger,
	}
}

func (nh NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx)

	var req CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nh.logger.Error("Decoding failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "Failed to decode", httpx.CodeMalformedJSON)
		return
	}

	tx, err := nh.db.BeginTx(ctx, nil)
	if err != nil {
		nh.logger.Error("Transaction Failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}
	defer tx.Rollback()

	var out CreateNotificationResponse

	err = tx.QueryRowContext(ctx,
		`INSERT INTO notifications (user_id, channel, recipient, message) VALUES ($1,$2,$3,$4) RETURNING id,message,created_at`,
		req.UserID, req.Channel, req.Recipient, req.Message,
	).Scan(&out.NotificationId, &out.Message, &out.CreatedAt)

	if err != nil {
		nh.logger.Error("Insertion into notifications failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	payload, err := json.Marshal(outboxPayload{
		NotificationID: out.NotificationId,
		UserID:         req.UserID,
		Channel:        string(req.Channel),
		Recipient:      req.Recipient,
		Message:        req.Message,
	})

	if err != nil {
		nh.logger.Error("Marshalling Failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO outbox_events (notification_id,event_type,payload) 
		VALUES($1,$2,$3)`,
		out.NotificationId, "NOTIFICATION_CREATED", payload,
	)

	if err != nil {
		nh.logger.Error("DB Execution Failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	if err := tx.Commit(); err != nil {
		nh.logger.Error("Transaction Commit Failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		nh.logger.Error("Writing response failed", "request_id", requestId, "err", err)
	}
}
