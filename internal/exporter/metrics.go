package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/siangyeh8818/mysql-to-prometheus/internal/lib"
)

func AddMetrics() map[string]*prometheus.Desc {

	DBMetrics := make(map[string]*prometheus.Desc)

	DBMetrics["stream_network_Score"] = prometheus.NewDesc(
		prometheus.BuildFQName("solar", "stream", "network_Score"),
		"A metric with a constant '0' value labeled by matchId,roomId,locationId from DB table.",
		[]string{"matchId", "streamId", "location", "stateCategory", "dbtable"}, nil,
	)
	/*
		DBMetrics["Workflows"] = prometheus.NewDesc(
			prometheus.BuildFQName("pentium", "marvin", "workflows"),
			"Size of workflows",
			[]string{}, nil,
		)

		DBMetrics["WorkflowInstances"] = prometheus.NewDesc(
			prometheus.BuildFQName("pentium", "marvin", "workflow_instances"),
			"Details information of workflow instances",
			[]string{"status", "reason"}, nil,
		)


	*/
	log.Println("Metrics added.....")

	return DBMetrics
}

func (e *Exporter) processMetrics(data lib.Data, ch chan<- prometheus.Metric) error {
	// APIMetrics - range through the data slice
	for _, x := range data {

		//增加label的地方
		id := strconv.FormatUint(x.MatchID, 10)

		ch <- prometheus.MustNewConstMetric(e.DBMetrics["stream_network_Score"], prometheus.GaugeValue, float64(x.Score), id, x.SpecialID, x.Location, x.StateCategory, x.DBtable)
		//ch <- prometheus.MustNewConstMetric(e.DBMetrics["Workflows"], prometheus.GaugeValue, x.WorkflowSize)
		//ch <- prometheus.MustNewConstMetric(e.DBMetrics["WorkflowInstances"], prometheus.GaugeValue, x.RunningWorkflowSize, "running", "")
	}
	return nil
}
