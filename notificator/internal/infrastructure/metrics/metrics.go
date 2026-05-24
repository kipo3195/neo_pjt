package metrics

import "github.com/prometheus/client_golang/prometheus"

func MetricsRegister() {
	prometheus.MustRegister(
		ActiveConnections,
		MessagesTotal,
		RoomSubscriptions,
		DispatchDuration,
	)
}

func WSConnected() {
	ActiveConnections.Inc()
}

func WSDisconnected() {
	ActiveConnections.Dec()
}

func WSMessageReceived(status string) {
	MessagesTotal.WithLabelValues("inbound", status).Inc()
}

func WSMessageSent(status string) {
	MessagesTotal.WithLabelValues("outbound", status).Inc()
}

func SetRoomSubscriptions(count int) {
	RoomSubscriptions.Set(float64(count))
}

func ObserveDispatchDuration(seconds float64) {
	DispatchDuration.Observe(seconds)
}

var (
	ActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "notificator",
			Subsystem: "ws",
			Name:      "active_connections",
			Help:      "Current number of active websocket connections.",
		},
	)

	MessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "notificator",
			Subsystem: "ws",
			Name:      "messages_total",
			Help:      "Total number of websocket messages.",
		},
		[]string{"direction", "status"},
	)

	RoomSubscriptions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "notificator",
			Subsystem: "ws",
			Name:      "room_subscriptions",
			Help:      "Current number of room subscriptions loaded in memory.",
		},
	)

	DispatchDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "notificator",
			Subsystem: "ws",
			Name:      "dispatch_duration_seconds",
			Help:      "Message dispatch duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
	)
)
