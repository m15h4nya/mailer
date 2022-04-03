package server

import (
	"apitask/config"
	handlers "apitask/server/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

func NewServer(log *zap.SugaredLogger, db *gorm.DB, cfg *config.Config, ch chan struct{}) *http.Server {
	handler := handlers.NewHandler(log, db, ch)

	serveMux := mux.NewRouter()

	docHandler := middleware.Redoc(middleware.RedocOpts{SpecURL: "/swagger.yaml"}, nil)
	getHandler := serveMux.Methods(http.MethodGet).Subrouter()
	getHandler.Handle("/docs", docHandler)
	getHandler.Handle("/swagger.yaml", http.FileServer(http.Dir("/opt/apitask/docs/")))

	patchHandler := serveMux.Methods(http.MethodPatch).Subrouter()
	patchHandler.HandleFunc("/client/{id}", handler.UpdateClient)
	patchHandler.HandleFunc("/mailing/{id}", handler.UpdateMailing)

	getHandler.HandleFunc("/statistics", handler.GetStatistics)
	getHandler.HandleFunc("/mailing/statistics", handler.GetMailingStatistic)

	postHandler := serveMux.Methods(http.MethodPost).Subrouter()
	postHandler.HandleFunc("/client", handler.AddClient)
	postHandler.HandleFunc("/mailing", handler.AddMailing)
	postHandler.HandleFunc("/mailing/start", handler.StartMailing)

	deleteHandler := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteHandler.HandleFunc("/client/{id}", handler.DeleteClient)
	deleteHandler.HandleFunc("/mailing/{id}", handler.DeleteMailing)

	server := &http.Server{
		Handler: serveMux,
		Addr:    cfg.ListenAddr,
	}

	return server
}
