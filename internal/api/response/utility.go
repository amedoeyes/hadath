package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amedoeyes/hadath/internal/api"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func WriteJSONError(w http.ResponseWriter, status int, message string, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := ErrorResponse{
		Message: message,
		Errors:  errors,
	}
	json.NewEncoder(w).Encode(response)
}

func WriteJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func HandleError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := err.Error()
	var errors map[string]string = nil

	switch err := err.(type) {
	case validator.ValidationErrors:
		status = http.StatusBadRequest
		message = "Validation failed"
		errors = make(map[string]string)
		for _, e := range err {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
			case "gt":
				errors[field] = fmt.Sprintf("%s must be greater than %s", field, e.Param())
			case "gte":
				errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
			case "gtfield":
				errors[field] = fmt.Sprintf("%s must be greater than %s", field, e.Param())
			case "lt":
				errors[field] = fmt.Sprintf("%s must be less than %s", field, e.Param())
			case "lte":
				errors[field] = fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
			case "ltfield":
				errors[field] = fmt.Sprintf("%s must be less than %s", field, e.Param())
			default:
				errors[field] = fmt.Sprintf("Invalid value for %s", field)
			}
		}
	case *pgconn.PgError:
		switch err.Code {
		case "23505":
			status = http.StatusConflict
			switch err.ConstraintName {
			case "users_email_key":
				message = "An account with this email already exists"
			case "bookings_user_id_event_id_key":
				message = "User already booked to this event"
			default:
				message = err.Error()
			}
		default:
			status = http.StatusInternalServerError
			message = err.Error()
		}
	default:
		switch err {
		case io.EOF:
			status = http.StatusBadRequest
			message = "Empty request"
		case bcrypt.ErrMismatchedHashAndPassword:
			status = http.StatusUnauthorized
			message = "Invalid credentials"
		case api.ErrUnauthorized:
			status = http.StatusUnauthorized
			message = "Unauthorized"
		}
	}

	WriteJSONError(w, status, message, errors)
}
