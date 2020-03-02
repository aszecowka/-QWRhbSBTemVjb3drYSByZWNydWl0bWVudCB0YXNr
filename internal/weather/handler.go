package weather

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

//go:generate mockery -name=BulkGetter -output=automock -outpkg=automock -case=underscore
type BulkGetter interface {
	GetWeatherForCities(ctx context.Context, cities []string) (map[string]*OpenWeatherResponse, error)
}

func NewHandler(getter BulkGetter, logger logrus.FieldLogger) *handler {
	return &handler{
		getter: getter,
		logger: logger,
	}
}

type handler struct {
	getter BulkGetter
	logger logrus.FieldLogger
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cities := req.URL.Query()["city"]

	if len(cities) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(rw).Encode(BadRequestResponse{Message: "`city` query parameter is required"})
		if err != nil {
			h.logger.Warnf("Got error on encoding Bad Request Response: %v", err)
		}
		return
	}

	out, err := h.getter.GetWeatherForCities(req.Context(), cities)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.logger.Warnf("Got error on getting weather for cities: %v", err)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(rw).Encode(out); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.logger.Warnf("Got error on encoding response: %v", err)
		return
	}
}
