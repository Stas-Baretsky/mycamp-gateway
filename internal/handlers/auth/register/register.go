package register

import (
	"context"
	"errors"
	"gateway-api/internal/domain/auth"
	resp "gateway-api/internal/lib/api/response"
	"gateway-api/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

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

type UserRegistrator interface {
	Register(
		ctx context.Context,
		params auth.RegisterParams,
	) (int64, error)
}

var valid = validator.New()

func New(log *slog.Logger, userRegister UserRegistrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.register.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req RegisterRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info(
			"request body decoded",
			slog.String("email", req.Email),
		)

		var validateErr validator.ValidationErrors
		if err = valid.Struct(req); err != nil {
			if errors.As(err, &validateErr) {

				log.Error("invalid request", sl.Err(err))

				render.JSON(w, r, resp.ValidationError(validateErr))

				return
			}
			render.JSON(w, r, resp.Error("validation failed"))
			return
		}

		id, err := userRegister.Register(
			r.Context(),
			auth.RegisterParams{
				Email:    req.Email,
				Password: req.Password,
			},
		)

		//TODO: добавить возврат ошибки
		// switch status.Code(err) {
		// 	case codes.AlreadyExists:
		// 		// 409 Conflict
		// 	case codes.InvalidArgument:
		// 		// 400 Bad Request
		// 	default:
		// 		// 500 Internal Server Error

		if err != nil {
			log.Error("register failed", sl.Err(err))

			render.JSON(w, r, resp.Error("register error"))

			return
		}

		log.Info("user registrated", slog.Int64("user_id", id))

		render.Status(r, http.StatusCreated)

		render.JSON(w, r, RegisterResponse{
			resp.OK(),
			id,
		})
	}
}
