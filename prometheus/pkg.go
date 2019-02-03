package prometheus

/*
Guide: https://prometheus.io/docs/guides/go-application/


go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp

Instead of using the following:

		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9912", nil)

	You could use:

		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())

		s := &http.Server{Addr: ":9912", Handler: h}
		s.ListenAndServe()


*/

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

// App main struct for Prometheus
type App struct {
	Router *mux.Router
}

// Initilize can be used for testing
func (a *App) Initilize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	recordMetrics()
	a.Router.Handle("/metrics", promhttp.Handler())
	a.Router.HandleFunc("/", a.info).Methods("GET")
}

// Run is to be used on live server
func (a *App) Run(addr string, writeTimeout int, readTimeout int) {
	srv := &http.Server{
		Handler: a.Router,
		Addr:    fmt.Sprintf(":%s", addr),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func (a *App) info(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	s := fmt.Sprintf("{repo:\"%s\"}",
		"https://github.com/mchirico/go_script")
	_, err := w.Write([]byte(s))
	if err != nil {
		log.Printf("Can not write response: %v\n", err)
	}

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Printf("Can not write response: %v\n", response)
	}
}
