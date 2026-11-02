package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/httpx"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/middleware"
	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/validation"
	"github.com/jackc/pgx/v5/pgconn"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
	}
}

func (ah AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx)
	ah.logger.With("request_id", requestId)

	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.logger.Error("Failed to decode sign up data", "err", err)
		httpx.Error(w, http.StatusBadRequest, "Bad Request", httpx.CodeValidationFailed)
		return
	}

	if err := req.ValidateSignUp(); err != nil {
		var verr *validation.ValidationError
		errors.As(err, &verr)
		ah.logger.Error("Sign up request validation failed!!", "err", err)
		httpx.ValidationError(w, http.StatusBadRequest, verr.Error(), httpx.CodeValidationFailed, verr.Field)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		ah.logger.Error("Failed to hash password in sign up", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	row := ah.db.QueryRowContext(ctx,
		`INSERT INTO users (username,email,password)
	VALUES ($1,$2,$3) RETURNING id,created_at`, req.Username, req.Email, hashedPassword)

	var out SignupResponse
	if err := row.Scan(&out.ID, &out.Created_At); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			httpx.Error(w, http.StatusConflict, "Email Already Taken", httpx.CodeConflict)
			return
		}
		ah.logger.Error("Failed to scan into signup DTO", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	ah.logger.Info("New User Registered", "user_id", out.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)
}
