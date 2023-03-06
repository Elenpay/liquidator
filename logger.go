package main

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func init() {

	//add a hook to the logger to add the fields in the context
	log.AddHook(&logrusContextHook{})

	//If log level is debug
	// Add this line for logging filename and line number!
	if log.StandardLogger().GetLevel() == log.DebugLevel {
		log.SetReportCaller(true)
	}

}

type logrusContextHook struct {
}

func (hook *logrusContextHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook *logrusContextHook) Fire(entry *log.Entry) error {

	//If there is a key with value "span" in the data, convert the value to a span
	if span, ok := entry.Data["span"]; ok {

		//Convert the value to a span
		span := span.(trace.Span)

		//Remove the key "span" from the data
		delete(entry.Data, "span")

		//Add the fields of trace id and span id to the log entry
		entry.Data["dd.trace_id"] = convertTraceID(span.SpanContext().TraceID().String())
		entry.Data["dd.span_id"] = convertTraceID(span.SpanContext().SpanID().String())
	}

	return nil
}

// Took from DD https://docs.datadoghq.com/tracing/other_telemetry/connect_logs_and_traces/opentelemetry?tab=go
func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
