package statistics

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type messageStatistic struct {
	mtx          sync.RWMutex
	sent         int64
	sentCounter  prometheus.Counter
	failed       int64
	sumRTT       time.Duration
	rttHistogram prometheus.Histogram
}

func newMessageStatistic() *messageStatistic {
	return &messageStatistic{
		sentCounter:  newPrometheusCounter("message", "sent_total", "the total number of messages sent by robot"),
		rttHistogram: newPrometheusHistogram("message", "rtt_seconds", "the histogram of message rtt seconds", 0.1, 0.5, 1),
	}
}

func (r *messageStatistic) Received(rtt time.Duration) {
	r.mtx.Lock()
	r.sent++
	r.sentCounter.Inc()
	r.sumRTT += rtt
	r.rttHistogram.Observe(rtt.Seconds())
	r.mtx.Unlock()
}

func (r *messageStatistic) SentFailed() {
	r.mtx.Lock()
	r.failed++
	r.mtx.Unlock()
}

func (r *messageStatistic) Output(w io.Writer) error {
	r.mtx.RLock()
	sent := r.sent
	failed := r.failed
	delay := time.Duration(0)
	if sent > 0 {
		delay = r.sumRTT / time.Duration(r.sent) / 2
	}
	r.mtx.RUnlock()

	_, err := io.WriteString(w, fmt.Sprintf(`message:
	sent: %d
	failed: %d
	delay: %vms
`,
		sent,
		failed,
		delay.Milliseconds(),
	))

	return err
}

type loginStatistic struct {
	mtx            sync.RWMutex
	success        int64
	successCounter prometheus.Counter
	failed         int64
	failedCounter  prometheus.Counter
	maxTime        time.Duration
	maxTimeGauge   prometheus.Gauge
	sumTime        time.Duration
	timeHistogram  prometheus.Summary
}

func newLoginStatistic() *loginStatistic {
	return &loginStatistic{
		successCounter: newPrometheusCounter("login", "success_total", "the total number of successfull login"),
		failedCounter:  newPrometheusCounter("login", "failed_total", "the total number of login failures"),
		maxTimeGauge:   newPrometheusGauge("login", "max_use_time_seconds", "the maximum use time seconds of login"),
		timeHistogram:  newPrometheusHistogram("login", "use_time_seconds", "the histogram of login use time seconds", 1, 5, 10),
	}
}

func (l *loginStatistic) Success(time time.Duration) {
	l.mtx.Lock()
	l.success++
	l.successCounter.Inc()
	l.timeHistogram.Observe(time.Seconds())
	l.sumTime += time
	if time > l.maxTime {
		l.maxTime = time
		l.maxTimeGauge.Set(l.maxTime.Seconds())
	}
	l.mtx.Unlock()
}

func (l *loginStatistic) Fail() {
	l.mtx.Lock()
	l.failed++
	l.failedCounter.Inc()
	l.mtx.Unlock()
}

func (l *loginStatistic) Output(w io.Writer) error {
	l.mtx.RLock()
	success := l.success
	failed := l.failed
	maxTime := l.maxTime
	averageTime := time.Duration(0)
	if success > 0 {
		averageTime = l.sumTime / time.Duration(success)
	}
	l.mtx.RUnlock()

	_, err := io.WriteString(w, fmt.Sprintf(`login:
	success: %d
	failed: %d
	maxTime: %fs
	averageTime: %fs
`,
		success,
		failed,
		maxTime.Seconds(),
		averageTime.Seconds(),
	))

	return err
}

type statistics struct {
	message *messageStatistic
	login   *loginStatistic
}

// func (s *statistics) Message() *messageStatistic { return s.message }

// func (s *statistics) Login() *loginStatistic { return s.login }

func (s *statistics) Output(w io.Writer) error {
	var err error

	if _, err = io.WriteString(w, "<<<<<<<<<<<<<<<<<<<<<< statistics >>>>>>>>>>>>>>>>>>>>>\n"); err != nil {
		return err
	}

	if err = s.message.Output(w); err != nil {
		return err
	}

	return s.login.Output(w)
}

var gStatistics = &statistics{
	message: newMessageStatistic(),
	login:   newLoginStatistic(),
}

var gServer *http.Server

var errMetricsExposed = errors.New("metrics exposed")
var errMetricsConcealed = errors.New("metrics concealed")

// 暴露指标，通过 http://addr:port/metrics 拉取
func ExposeMetrics(port int, errLogger *log.Logger) error {
	if gServer != nil {
		return errMetricsExposed
	}

	gServer = &http.Server{
		Addr:     ":" + strconv.Itoa(port),
		ErrorLog: errLogger,
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	gServer.Handler = mux

	go func() {
		if err := gServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errLogger.Fatalln("expose metrics:", err)
		}
	}()

	return nil
}

// 隐藏指标
func ConcealMetrics() error {
	if gServer == nil {
		return errMetricsConcealed
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return gServer.Shutdown(ctx)
}

func Statistics() *statistics { return gStatistics }

func Message() *messageStatistic { return gStatistics.message }

func Login() *loginStatistic { return gStatistics.login }
