package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prometheus2 "initialthree/node/common/prometheus"
	"initialthree/pkg/util"
	"net/http"
)

var (
	enablePrometheus bool = false

	onlinePlayerFunc func() float64
	onlinePlayer     prometheus.GaugeFunc
)

func RegisterOnlinePlayerFunc(f func() float64) {
	onlinePlayerFunc = f
}

func initMonitor() {
	enablePrometheus = true
	if onlinePlayerFunc != nil {
		onlinePlayer = prometheus2.NewPrometheusGaugeFunc("initialthree", "game", "online_player", "", onlinePlayerFunc)
	}
}

func StartPrometheus(address string) {
	initMonitor()

	address = util.ParseAddress(address)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(address, nil)
		if err != nil {
			panic(err)
		}
	}()

}
