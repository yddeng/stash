package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prometheus2 "initialthree/node/common/prometheus"
	"initialthree/pkg/util"
	"net/http"
)

/*
	1. gate 与 game 的链接数
	2. gate 与 game 链接上的 stream数量
	3. gate 上 client 的数量
	4. 请求量(不包含登陆请求)
	5. 登陆请求成功量、失败量; 处理时间
*/
var (
	enablePrometheus bool = false

	// 1
	gameSocketCountFunc func() float64

	// 2
	gameSocketStreamCount *prometheus.GaugeVec

	// 3
	channelCount prometheus.Gauge

	// 4
	userRequestTotal prometheus.Counter

	// 5
	userLoginSuccessTotal prometheus.Counter
	userLoginFailedTotal  prometheus.Counter
	userLoginRtt          prometheus.Histogram
)

func RegisterGameSocketCountFunc(f func() float64) {
	gameSocketCountFunc = f
}

func GameSocketStreamInc(values ...string) {
	if enablePrometheus {
		gameSocketStreamCount.WithLabelValues(values...).Inc()
	}
}
func GameSocketStreamDec(values ...string) {
	if enablePrometheus {
		gameSocketStreamCount.WithLabelValues(values...).Dec()
	}
}

func ChannelCountInc() {
	if enablePrometheus {
		channelCount.Inc()
	}
}

func ChannelCountDec() {
	if enablePrometheus {
		channelCount.Dec()
	}
}

func UserRequestInc() {
	if enablePrometheus {
		userRequestTotal.Inc()
	}
}

func UserLoginSuccessInc() {
	if enablePrometheus {
		userLoginSuccessTotal.Inc()
	}
}

func UserLoginFailedInc() {
	if enablePrometheus {
		userLoginFailedTotal.Inc()
	}
}

func initMonitor() {
	enablePrometheus = true
	if gameSocketCountFunc != nil {
		prometheus2.NewPrometheusGaugeFunc(
			"initialthree",
			"gate",
			"game_socket_count",
			"Number of tcp connection to game.",
			gameSocketCountFunc,
		)
	}

	gameSocketStreamCount = prometheus2.NewPrometheusGaugeVec(
		"initialthree",
		"gate",
		"game_socket_stream_count",
		"Number of the stream",
		[]string{"logic"},
	)

	channelCount = prometheus2.NewPrometheusGauge("initialthree", "gate", "channel_count", "Number of channel.")

	userRequestTotal = prometheus2.NewPrometheusCounter("initialthree", "gate", "user_request_total", "Number of user request")
	userLoginSuccessTotal = prometheus2.NewPrometheusCounter("initialthree", "gate", "user_login_success_total", "Number of user login for success.")
	userLoginFailedTotal = prometheus2.NewPrometheusCounter("initialthree", "gate", "user_login_failed_total", "Number of user login for failed.")
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
