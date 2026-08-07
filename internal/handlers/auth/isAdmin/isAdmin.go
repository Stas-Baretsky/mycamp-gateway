package isadmin

import (
	"context"
	"errors"
	resp "gateway-api/internal/lib/api/response"
	"gateway-api/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type isAdminRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type isAdminResponse struct {
	resp.Response
	IsAdmin bool `json:"is_admin"`
}

type UserAdminChecker interface {
	IsAdmin(
		ctx context.Context,
		UserId int64,
	) (bool, error)
}

var valid = validator.New()

func New(log *slog.Logger, adminChecker UserAdminChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.isAdmin.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req isAdminRequest

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info(
			"request body decoded",
			slog.Int("email", int(req.UserID)),
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

		isAdmin, err := adminChecker.IsAdmin(
			r.Context(),
			req.UserID,
		)

		if err != nil {
			log.Error("is admin check failed", sl.Err(err))
			render.Status(r, http.StatusBadGateway)
			render.JSON(w, r, resp.Error("login error"))

			return
		}

		log.Info(
			"loggin successful",
			slog.Int64("user_id", req.UserID),
			slog.Bool("is_admin", isAdmin),
		)

		render.Status(r, http.StatusFound)

		render.JSON(w, r, isAdminResponse{
			resp.OK(),
			isAdmin,
		})

	}
}
