package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/imega/commerceml2teleport/health"
	"github.com/imega/commerceml2teleport/shutdown"
	"github.com/improbable-eng/go-httpwares/logging/logrus"
	"github.com/improbable-eng/go-httpwares/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	logger *logrus.Entry
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
	})
	logger = logrus.WithField("channel", "commerce-ml2teleport")

	grpcSrv := grpc.NewServer()
	health.RegisterHealthServer(grpcSrv)

	l, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		logger.Errorf("failed to listen on the TCP network address 0.0.0.0:9000, %s", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/{uuid}", handler)

	m := http.NewServeMux()
	m.Handle("/", router)

	hm := http_logrus.Middleware(logger, http_logrus.WithRequestFieldExtractor(func(req *http.Request) map[string]interface{} {
		return map[string]interface{}{
			"http.request.x-req-id": "unset",
		}
	}))(m)
	s := &http.Server{
		Addr:         "0.0.0.0:80",
		Handler:      hm,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	shutdown.RegisterShutdownFunc(func() {
		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
		s.Shutdown(ctx)
	})

	go grpcSrv.Serve(l)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logrus.Errorf("failed to listen on the TCP network address %s and handle requests on incoming connections, %s", s.Addr, err)
		}
	}()

	logger.Info("server is started")
	shutdown.LoopUntilShutdown(15 * time.Second)
	logger.Info("server is stopped")
}

func handler(w http.ResponseWriter, req *http.Request) {
	var (
		ctx    = req.Context()
		logger = ctxlogrus.Extract(ctx)
		uuid   = mux.Vars(req)["uuid"]
	)

	if len(uuid) < 1 {
		logger.Errorf("url path doest not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Infof("url path is %s", uuid)
	//go parser.Parse(req.URL.Path)
}
