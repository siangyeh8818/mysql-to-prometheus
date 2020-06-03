package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"

	exporter "github.com/siangyeh8818/mysql-to-prometheus/internal/exporter"
	lib "github.com/siangyeh8818/mysql-to-prometheus/internal/lib"
)

var metrics map[string]*prometheus.Desc
var appConfig lib.BaseConfig

func init() {
	metrics = exporter.AddMetrics()
}

func main() {
	log.Println("Exporter is start ro running")

	//db.DB_Handler()
	/*

		go func() {
			crawler.CallSelium()
		}()

		server.Run_Exporter_Server()
	*/

	// init Exporter
	exp := exporter.Exporter{
		DBMetrics: metrics,
		Config:    appConfig,
		//Cache:     appCache,
		//K8sClient: k8sclient,
	}
	// start the server
	exporter.NewServer(exp).Start()

}
