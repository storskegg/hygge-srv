package hygge_srv

import "github.com/prometheus/client_golang/prometheus"

var (
	gaugeOptsHumi = prometheus.GaugeOpts{
		Name: "hygge_humi",
		Help: "Humidity in % RH",
	}
	gaugeOptsTemp = prometheus.GaugeOpts{
		Name: "hygge_temp",
		Help: "Temperature in C",
	}
	gaugeOptsBatt = prometheus.GaugeOpts{
		Name: "hygge_batt",
		Help: "Battery level in V",
	}

	gaugeHumi = prometheus.NewGauge(gaugeOptsHumi)
	gaugeTemp = prometheus.NewGauge(gaugeOptsTemp)
	gaugeBatt = prometheus.NewGauge(gaugeOptsBatt)
)
