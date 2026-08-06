package isadmin

import (
	"context"
	resp "gateway-api/internal/lib/api/response"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
)

type isAdminRequest struct {
	UserID int64 `json:"user_id"`
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
	}
}
