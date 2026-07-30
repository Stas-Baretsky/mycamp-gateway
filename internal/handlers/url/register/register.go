package register

import (
	resp "gateway-api/internal/lib/api/response"
	"gateway-api/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// TODO :: // Добавить в сервис Auth возврат ошибки при попытке регистрации
type RegisterResponse struct {
	resp.Response
	UserID int64 `json:"userid,omitempty"`
}

type UserRegister interface {
	Register(email string, password string) int
}

func New(log *slog.Logger, userRegister UserRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.register.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqId(r.Context()))
		)
		var req RegisterRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}
		 
	}
}
