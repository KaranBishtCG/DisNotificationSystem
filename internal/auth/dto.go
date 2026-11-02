package auth

import (
	"net/mail"
	"strings"
	"time"

	"github.com/Karan-Singh-Bisht/DisNotiSystem/internal/validation"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	ID         string    `json:"id"`
	Created_At time.Time `json:"created_at"`
}

func (req SignUpRequest) ValidateSignUp() error {
	if strings.TrimSpace(req.Username) == "" {
		return &validation.ValidationError{
			Field: "Username",
			Msg:   "Username is required!!",
		}
	}

	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return &validation.ValidationError{
			Field: "Email",
			Msg:   "Enter valid email!!",
		}
	}

	if len(req.Password) < 8 {
		return &validation.ValidationError{
			Field: "Password",
			Msg:   "Password must be atleast 8 characters!!",
		}
	}

	return nil
}
