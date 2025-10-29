package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"

	"github.com/ShopOnGO/search-service/configs"
	"github.com/ShopOnGO/search-service/internal/elastic"
	"github.com/ShopOnGO/search-service/internal/search"
	"github.com/ShopOnGO/search-service/migrations"
	"github.com/ShopOnGO/search-service/pkg/db"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ShopOnGO/search-service/internal/graph"
)

var (
	httpSrv *http.Server
)

type App struct {
	conf          *configs.Config
	kafkaConsumer *kafkaService.KafkaService
	searchSvc     *search.SearchService
}

func InitServices() *App {
	migrations.CheckForMigrations()
	conf := configs.LoadConfig()

	consoleLvl := conf.LogLevel
	fileLvl := conf.FileLogLevel
	logger.InitLogger(consoleLvl, fileLvl)
	logger.EnableFileLogging("TailorNado_search-service")

	_ = db.NewDB(conf)

	elastic.Init(conf)

	searchSvc := search.NewSearchService()

	kafkaConsumer := kafkaService.NewConsumer(
		conf.Kafka.Brokers,
		conf.Kafka.Topic,
		conf.Kafka.GroupID,
		conf.Kafka.ClientID,
	)

	return &App{
		conf:          conf,
		kafkaConsumer: kafkaConsumer,
		searchSvc:     searchSvc,
	}
}

func RunHTTPServer(app *App) {
	router := gin.Default()

	searchHandlerDeps := search.SearchHandlerDeps{
		SearchSvc: app.searchSvc,
	}

	search.NewSearchHandler(router, searchHandlerDeps)

	// GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.POST("/search", gin.WrapH(srv))
	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL", "/search").ServeHTTP(c.Writer, c.Request)
	})

	httpSrv = &http.Server{
		Addr:    ":8085",
		Handler: router,
	}

	go func() {
		logger.Info("HTTP server listening on :8085")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Infof("HTTP server error: %v\n", err)
		}
	}()
}

// func RunGRPCServer(app *App, wg *sync.WaitGroup) *grpc.Server {
// 	defer wg.Done()
// 	listener, err := net.Listen("tcp", ":50052")
// 	if err != nil {
// 		logger.Infof("TCP listener error: %v\n", err)
// 		return nil
// 	}
// 	grpcServer := grpc.NewServer()
// 	pb.RegisterReviewServiceServer(grpcServer, review.NewGrpcReviewService(app.reviewSvc))
// 	pb.RegisterQuestionServiceServer(grpcServer, question.NewGrpcQuestionService(app.questionSvc))
// 	logger.Info("gRPC server listening on :50052")
// 	if err := grpcServer.Serve(listener); err != nil {
// 		logger.Infof("gRPC server error: %v\n", err)
// 	}
// 	return grpcServer
// }

func RunKafkaConsumer(ctx context.Context, app *App) {
	defer app.kafkaConsumer.Close()

	dispatcher := kafkaService.NewDispatcher()
	dispatcher.Register("product-created", func(msg kafka.Message) error {
		return search.HandleProductEvent(msg.Value, string(msg.Key), app.searchSvc)
	})

	logger.Info("Kafka consumer started")
	app.kafkaConsumer.Consume(ctx, dispatcher.Dispatch)
}

func WaitForShutdown(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	logger.Info("Shutdown signal received")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownCancel()
	if httpSrv != nil {
		logger.Info("Shutting down HTTP server...")
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			logger.Infof("HTTP shutdown error: %v\n", err)
		}
	}
}
