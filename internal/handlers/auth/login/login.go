package login

import (
	"context"
	"errors"
	"gateway-api/internal/domain/auth"
	resp "gateway-api/internal/lib/api/response"
	"gateway-api/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	AppID    int32  `json:"app_id"`
}

type LoginResponse struct {
	resp.Response
	Token string `json:"token"`
}

type UserLoginer interface {
	Login(
		ctx context.Context,
		params auth.LoginParams,
	) (string, error)
}

var valid = validator.New()

func New(log *slog.Logger, userLoginer UserLoginer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req LoginRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info(
			"request body decoded",
			slog.String("email", req.Email),
			slog.Int("app_id", int(req.AppID)),
		)

		var validateErr validator.ValidationErrors

		if err = valid.Struct(req); err != nil {
			if errors.As(err, &validateErr) {
				log.Error("invalid request", sl.Err(err))

				render.JSON(w, r, resp.ValidationError(validateErr))
			}
			render.JSON(w, r, resp.Error("validation failed"))
			return
		}
		////TODO:: извлекать app ID из заголовка или домена
		token, err := userLoginer.Login(
			r.Context(),
			auth.LoginParams{
				Email:    req.Email,
				Password: req.Password,
				AppID:    req.AppID,
			})

		////TODO:: добавить обработку статус кодов

		if err != nil {
			log.Error("login failed", sl.Err(err))

			render.JSON(w, r, resp.Error("login error"))

			return
		}

		log.Info("loggin successful", slog.String("user", req.Email))

		render.Status(r, http.StatusOK)

		render.JSON(w, r, LoginResponse{
			resp.OK(),
			token,
		})
	}
}
