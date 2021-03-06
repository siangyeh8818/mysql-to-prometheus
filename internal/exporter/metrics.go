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
		[]string{"matchId", "streamId", "location", "stateCategory", "dbtable", "Status", "sportName", "prikey", "level"}, nil,
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
	log.Println("--------processMetrics--------")
	for _, x := range data {

		//增加label的地方
		id := strconv.FormatUint(x.MatchID, 10)
		str1 := strconv.Itoa(x.StreamAPIStatus)
		level := strconv.Itoa(x.Level)
		//legid := strconv.FormatUint(x.LeagueId, 10)
		//log.Printf("x.StreamAPIStatus : " + str1)
		//log.Printf("string(x.StreamAPIStatus) : %s", string(x.StreamAPIStatus))
		//log.Printf("x.SportName : %s", x.SportName)
		ch <- prometheus.MustNewConstMetric(e.DBMetrics["stream_network_Score"], prometheus.GaugeValue, float64(x.Score), id, x.SpecialID, x.Location, x.StateCategory, x.DBtable, str1, x.SportName, x.PriKey, level)
		//ch <- prometheus.MustNewConstMetric(e.DBMetrics["Workflows"], prometheus.GaugeValue, x.WorkflowSize)
		//ch <- prometheus.MustNewConstMetric(e.DBMetrics["WorkflowInstances"], prometheus.GaugeValue, x.RunningWorkflowSize, "running", "")
	}
	return nil
}
