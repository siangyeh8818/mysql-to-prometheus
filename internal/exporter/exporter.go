package exporter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/siangyeh8818/mysql-to-prometheus/internal/lib"
	db "github.com/siangyeh8818/mysql-to-prometheus/internal/mysql"
)

type Exporter struct {
	DBMetrics map[string]*prometheus.Desc
	Config    lib.BaseConfig
	//remaining_device prometheus.Gauge
	//gaugeVec prometheus.GaugeVec
}

type Server struct {
	Handler  http.Handler
	exporter Exporter
}

// Start Running a Server instance
func (s *Server) Start() {
	log.Fatal(http.ListenAndServe(":8088", s.Handler))
}

func NewServer(exporter Exporter) *Server {
	log.Println(`
 	 This is a prometheus exporter for stream
  	Access: http://127.0.0.1:8088
  	`)

	r := http.NewServeMux()

	metricsPath := "/metrics"
	//listenAddress := ":8081"
	//metricsPrefix := "stream"
	//exporters := NewExporter(metricsPrefix)
	/*
	   	registry := prometheus.NewRegistry()
	       registry.MustRegister(metrics)
	*/
	//exporter.Options{}
	prometheus.MustRegister(&exporter)

	// Launch http service

	r.Handle(metricsPath, promhttp.Handler())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		 <head><title>Dummy Exporter</title></head>
		 <body>
		 <h1>Stream Exporter</h1>
		 <p><a href='` + metricsPath + `'>Metrics</a></p>
		 </body>
		 </html>`))
	})

	return &Server{Handler: r, exporter: exporter}
}

/*
func NewExporter(metricsPrefix string) *Exporter {
	remaining_device := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: metricsPrefix,
		Name:      "remaining_device",
		Help:      "This is a gauge metric example"})
	/*
		gaugeVec := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: metricsPrefix,
			Name:      "gauge_vec_metric",
			Help:      "This is a siang gauga vece metric"},
			[]string{"myLabel"})

	return &Exporter{
		DBMetrics: remaining_device,
	}
}
*/

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("-------Collect-----------")
	data, err := e.gatherData()
	//log.Println(data)
	if err != nil {
		log.Fatalf("Error gathering Data from Mysql server: %v", err)
		//log.Errorf("Error gathering Data from Mysql server: %v", err)
		return
	}

	err = e.processMetrics(data, ch)
	if err != nil {
		log.Fatalf("Error Processing Metrics", err)
	}
	log.Println("All Metrics successfully collected.")
	//log.Info("All Metrics successfully collected.")

}

// 讓exporter的prometheus屬性呼叫Describe方法

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	//e.remaining_device.Describe(ch)

	for _, m := range e.DBMetrics {
		ch <- m
	}
}

func (e *Exporter) gatherData() (lib.Data, error) {

	//var data lib.Data
	//d := new(lib.Datum)

	//DB邏輯
	data := db.DB_Handler()

	//data = append(data, d)
	return data, nil

}

func GetCsvContent(filepath string) float64 {
	for !Exists(filepath) {
		log.Println("-----wait for hema-im-exporter to write output.csv'")
		time.Sleep(1 * time.Second)
	}

	fin, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	bytes, err := ioutil.ReadAll(fin)
	if err != nil {
		panic(err)
	}
	csccontent := string(bytes)
	fmt.Println(csccontent)
	csccontent = strings.Replace(csccontent, ",", "", -1)
	csccontent = strings.Replace(csccontent, "\n", "", -1)
	csccontent = strings.Replace(csccontent, "\r", "", -1)
	csccontent = strings.Trim(csccontent, "\"")
	downloadcount, err := strconv.ParseFloat(csccontent, 64)
	if err != nil {
		panic(err)
	}
	log.Println("------downloadcount------")
	log.Println(downloadcount)
	return downloadcount
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
