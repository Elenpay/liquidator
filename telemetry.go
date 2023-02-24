package main

import (
	context "context"
	"errors"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	log "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type metrics struct {
	channelBalanceGauge prometheus.GaugeVec
	// L1 Fees of Swaps
	onchainFees *prometheus.CounterVec
	// L2 Fees of Swaps
	offchainFees *prometheus.CounterVec
	// Fees of the swap provider (i.e. Loop)
	providerFees *prometheus.CounterVec
}

// Inits the prometheusMetrics global metric struct
func initMetrics(reg prometheus.Registerer) {

	log.Debug("Registering prometheus metrics")

	m := &metrics{
		channelBalanceGauge: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "liquidator_channel_balance",
			Help: "The total number of processed events",
		},
			[]string{"chan_id", "local_node_pubkey", "remote_node_pubkey", "local_node_alias", "remote_node_alias", "active", "initiator"},
		),
		onchainFees: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "liquidator_onchain_fees",
			Help: "Onchain fees of swaps in sats",
		},
			[]string{"node_alias", "swap_type", "chan_id", "provider"},
		),
		offchainFees: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "liquidator_offchain_fees",
			Help: "Offchain fees of swaps in sats",
		},
			[]string{"node_alias", "swap_type", "chan_id", "provider"}),
		providerFees: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "liquidator_provider_fees",
			Help: "Fees of the swap provider (i.e. Loop) in sats",
		},
			[]string{"node_alias", "swap_type", "chan_id", "provider"},
		),
	}
	//Register custom metrics
	reg.MustRegister(m.channelBalanceGauge)
	reg.MustRegister(m.onchainFees)
	reg.MustRegister(m.offchainFees)
	reg.MustRegister(m.providerFees)

	//Golang collector
	reg.MustRegister(collectors.NewGoCollector())
	//Process collector
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	prometheusMetrics = m
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(semconv.ServiceNameKey.String(OTELServiceName)),
		resource.WithFromEnv(),
	)
	if err != nil {
		log.Fatalf("Failed to detect environment resource: %v", err)
	}

	return r
}

func spanExporter() (*otlptrace.Exporter, error) {
	var otlpEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint != "" {
		log.Infof("exporting to OTLP collector at %s", otlpEndpoint)
		traceClient := otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(otlpEndpoint),
		)
		return otlptrace.New(context.Background(), traceClient)
	}
	return nil, errors.New("OTEL_EXPORTER_OTLP_ENDPOINT must not be empty")
}

// Init opentelemetry tracer
func initTracer(ctx context.Context) (*trace.TracerProvider, error) {

	//TracerProvider
	res := newResource()
	tp := trace.NewTracerProvider(trace.WithResource(res))

	// span exporter
	exp, err := spanExporter()
	if err != nil {
		log.Fatal("failed to initialize Span exporter")
	}

	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithSampler(trace.AlwaysSample()),
			trace.WithResource(res),
			trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exp)),
		),
	)

	return tp, nil

}
