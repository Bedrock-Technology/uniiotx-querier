package servers

import (
	"context"
	"errors"
	"fmt"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/Bedrock-Technology/uniiotx-querier/interactors"
	"github.com/Bedrock-Technology/uniiotx-querier/servers/middlewares"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
	"net/http"
)

// DataServer data gateway service
type DataServer struct {
	Logger common.Logger

	If *interactors.InteractorFactory

	Addr   string
	server *http.Server
}

func (s *DataServer) Start() {
	// Create server
	h := s.newHandler()
	s.server = &http.Server{
		Addr:    s.Addr,
		Handler: h,
	}

	// Start server
	s.Logger.Info("DataServer started", "endpoint", fmt.Sprintf("http://%s/docs", s.Addr))

	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Fatal("failed to start DataServer", err)
		}
	}
	s.Logger.Info("DataServer closed")
}

func (s *DataServer) Close() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.Logger.Error("failed to close DataServer", err)
		return
	}
}

func (s *DataServer) newHandler() http.Handler {
	r := web.NewService(openapi3.NewReflector())

	// Init API documentation schema.
	r.OpenAPISchema().SetTitle("uniIOTX contract querier")
	r.OpenAPISchema().SetDescription("This app provides a REST API for querying uniIOTX project's contract status")
	r.OpenAPISchema().SetVersion("v0.0.1")

	// Setup middlewares
	loggerMiddleware := middlewares.LoggerMiddleware(s.Logger)
	rateLimitMiddleware := middlewares.RateLimiterMiddleware(5, 10)

	r.Use(loggerMiddleware, rateLimitMiddleware)

	// Route data services
	r.Post("/stakedDelegates", s.If.ListStakedDelegatesFn()())
	r.Post("/stakedBuckets", s.If.ListStakedBucketsFn()())

	r.Post("/redeemedBuckets", s.If.ListRedeemedBucketsFn()())

	r.Post("/managerRewards", s.If.GetManagerRewardsFn()())
	r.Post("/managerRewardsByYear", s.If.ListManagerRewardsByYearFn()())
	r.Post("/managerRewardsByMonth", s.If.ListManagerRewardsByMonthFn()())

	r.Post("/assetStatistics", s.If.GetAssetStatisticsFn()())
	r.Post("/assetStatisticsByYear", s.If.ListAssetStatisticsByYearFn()())
	r.Post("/assetStatisticsByMonth", s.If.ListAssetStatisticsByMonthFn()())

	// Swagger UI endpoint at /docs.
	r.Docs("/docs", swgui.New)

	return r
}
