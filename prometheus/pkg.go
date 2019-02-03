package prometheus

/*
Guide: https://prometheus.io/docs/guides/go-application/
Ref: https://github.com/prometheus/haproxy_exporter


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
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics records data on go_script
type Metrics struct {
	sync.Mutex
	Loops    prometheus.Counter
	FileSize prometheus.Gauge
}

// Inc counter
func (m *Metrics) Inc() {
	m.Lock()
	defer m.Unlock()
	m.Loops.Inc()

}

// Size of file
func (m *Metrics) Size(size int64) {
	m.Lock()
	defer m.Unlock()
	m.FileSize.SetToCurrentTime()
	m.FileSize.Set(float64(size))
}

// Init setup
func (m *Metrics) Init() {
	m.Lock()
	defer m.Unlock()

	loops := promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_script_processed_ops_total",
		Help: "The total number of looped events",
	})
	m.Loops = loops

	fileSize := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_script_fileSize",
		Help: "The total number of looped events",
	})
	m.FileSize = fileSize
}

// App main struct for Prometheus
type App struct {
	Router  *mux.Router
	Metrics Metrics
}

// Initilize can be used for testing
func (a *App) Initilize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()

}

func (a *App) initializeRoutes() {
	a.Metrics = Metrics{}
	a.Metrics.Init()

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

// CustomRegistry registers a custom sample
func CustomRegistry() string {
	temps := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "pond_temperature_celsius",
			Help: "The temperature of the frog pond.",
			//Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"species"},
	)

	// Simulate some observations.
	for i := 0; i < 1000; i++ {
		temps.WithLabelValues("litoria-caerulea").Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
		temps.WithLabelValues("lithobates-catesbeianus").Observe(32 + math.Floor(100*math.Cos(float64(i)*0.11))/10)
	}

	// Create a Summary without any observations.
	temps.WithLabelValues("leiopelma-hochstetteri")

	// Just for demonstration, let's check the state of the summary vector
	// by registering it with a custom registry and then let it collect the
	// metrics.
	reg := prometheus.NewRegistry()
	reg.MustRegister(temps)

	metricFamilies, err := reg.Gather()
	if err != nil || len(metricFamilies) != 1 {
		panic("unexpected behavior of custom test registry")
	}
	//fmt.Println(proto.MarshalTextString(metricFamilies[0]))
	return proto.MarshalTextString(metricFamilies[0])
}
