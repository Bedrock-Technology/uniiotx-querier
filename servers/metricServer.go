package servers

import (
	"context"
	"errors"
	"fmt"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type MetricServer struct {
	Logger common.Logger

	Addr   string
	server *http.Server
}

func (s *MetricServer) Start() {
	// Create server
	http.Handle("/metrics", promhttp.Handler())
	s.server = &http.Server{
		Addr:    s.Addr,
		Handler: http.DefaultServeMux,
	}

	// Start server
	s.Logger.Info("MetricServer started", "endpoint", fmt.Sprintf("http://%s/metrics", s.Addr))
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Fatal("failed to start MetricServer", err)
		}
	}
	s.Logger.Info("MetricServer closed")
}

func (s *MetricServer) Close() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.Logger.Error("failed to close MetricServer", err)
	}
}
