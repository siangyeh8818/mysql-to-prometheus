package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	remaining_device prometheus.Gauge
	//gaugeVec prometheus.GaugeVec
}

func Run_Exporter_Server() {
	log.Println(`
  This is a prometheus exporter for hema.im
  Access: http://127.0.0.1:8081
  `)

	metricsPath := "/metrics"
	listenAddress := ":8081"
	metricsPrefix := "hema_im"
	exporters := NewExporter(metricsPrefix)
	/*
	   	registry := prometheus.NewRegistry()
	       registry.MustRegister(metrics)
	*/
	prometheus.MustRegister(exporters)

	// Launch http service

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		 <head><title>Dummy Exporter</title></head>
		 <body>
		 <h1>Dummy Exporter</h1>
		 <p><a href='` + metricsPath + `'>Metrics</a></p>
		 </body>
		 </html>`))
	})
	log.Println(http.ListenAndServe(listenAddress, nil))
}

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
	*/
	return &Exporter{
		remaining_device: remaining_device,
		//gaugeVec: gaugeVec,
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	//e.gauge.Set(float64(100))

	e.remaining_device.Set(GetCsvContent("output.csv"))
	//e.gaugeVec.WithLabelValues("hello").Set(float64(0))
	e.remaining_device.Collect(ch)
	//e.gaugeVec.Collect(ch)
}

// 讓exporter的prometheus屬性呼叫Describe方法

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.remaining_device.Describe(ch)
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
