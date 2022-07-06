package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/mem"
)

var (
	memTotalDesc     = prometheus.NewDesc("memexporter_memory_total_bytes", "The amount of total memory in bytes.", []string{}, nil)
	memUsedDesc      = prometheus.NewDesc("memexporter_memory_used_bytes", "The amount of used memory in bytes.", []string{}, nil)
	memAvailableDesc = prometheus.NewDesc("memexporter_memory_available_bytes", "The amount of available memory in bytes.", []string{}, nil)
	memFreeDesc      = prometheus.NewDesc("memexporter_memory_free_bytes", "The amount of free memory in bytes.", []string{}, nil)
)

type memCollector struct{}

func (mc memCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(mc, ch)
}

func (mc memCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	ch <- prometheus.MustNewConstMetric(memTotalDesc, prometheus.GaugeValue, float64(stats.Total))
	ch <- prometheus.MustNewConstMetric(memUsedDesc, prometheus.GaugeValue, float64(stats.Used))
	ch <- prometheus.MustNewConstMetric(memAvailableDesc, prometheus.GaugeValue, float64(stats.Available))
	ch <- prometheus.MustNewConstMetric(memFreeDesc, prometheus.GaugeValue, float64(stats.Free))
}

func main() {
	listenAddr := flag.String("listen-address", ":8080", "The address to listen on for metrics requests.")
	flag.Parse()

	prometheus.MustRegister(&memCollector{})

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
