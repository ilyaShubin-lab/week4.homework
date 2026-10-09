package health

import (
	"net/http"

	"go.uber.org/zap"

	"boilerplates/platform/pkg/logger"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		logger.Error(r.Context(), "failed to write health response", zap.Error(err))
	}
}
