package server

import (
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"
	handlers "github.com/lsmltesting/MicroBlog/services/api/internal/handlers/http"
)

type Config struct {
	MainPort       string
	PprofPort      string
	WithPprof      bool
	MaxHeaderBytes int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
}

type HTTPServer struct {
	userHttpHandler *handlers.UserHTTPHandler
	postHttpHandler *handlers.PostHTTPHandler
	likeHttpHandler *handlers.LikeHTTPHandler
	config          Config
	httpServer      *http.Server
}

func NewHTTPServer(
	config Config,
	userHttpHandler *handlers.UserHTTPHandler,
	postHttpHandler *handlers.PostHTTPHandler,
	likeHttpHandler *handlers.LikeHTTPHandler,
) *HTTPServer {
	return &HTTPServer{
		userHttpHandler: userHttpHandler,
		postHttpHandler: postHttpHandler,
		likeHttpHandler: likeHttpHandler,
		config:          config,
	}
}

func (s *HTTPServer) Run() error {
	s.RunPprof()

	router := mux.NewRouter()

	s.httpServer = &http.Server{
		Addr:           s.config.MainPort,
		Handler:        router,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
	}

	// Register methods from userHttpHandler
	router.Path("/register").Methods("POST").HandlerFunc(s.userHttpHandler.HandlerCreate)

	// Register methods from postHttpHandler
	router.Path("/posts").Methods("POST").HandlerFunc(s.postHttpHandler.HandlerCreate)
	router.Path("/posts").Methods("GET").HandlerFunc(s.postHttpHandler.HandlerGetAll)

	// Register method from likeHttpHandler
	router.Path("/posts/{post_id}/like").Methods("POST").Queries("user_id", "{user_id}").HandlerFunc(s.likeHttpHandler.HandlerCreate)
	router.Path("/likes").Methods("GET").HandlerFunc(s.likeHttpHandler.HandlerGetAll)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HTTPServer) RunPprof() {
	if !s.config.WithPprof {
		return
	}

	go func() {
		pprofRouter := mux.NewRouter()

		pprofServer := &http.Server{
			Addr:           s.config.PprofPort,
			Handler:        pprofRouter,
			MaxHeaderBytes: s.config.MaxHeaderBytes,
			ReadTimeout:    s.config.ReadTimeout,
			WriteTimeout:   s.config.WriteTimeout,
			IdleTimeout:    s.config.IdleTimeout,
		}

		pprofRouter.HandleFunc("/debug/pprof/", pprof.Index)
		pprofRouter.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofRouter.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofRouter.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofRouter.HandleFunc("/debug/pprof/trace", pprof.Trace)
		pprofRouter.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		pprofRouter.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		pprofRouter.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		pprofRouter.Handle("/debug/pprof/block", pprof.Handler("block"))
		pprofRouter.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
		pprofRouter.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))

		if err := pprofServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("pprofServer is failed: %v", err)
		}
	}()
}
