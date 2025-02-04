package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer() func() {
    exporter, err := otlptracehttp.New(context.Background(),
        otlptracehttp.WithEndpoint("localhost:4317"),
        otlptracehttp.WithInsecure(),
    )
    if err != nil {
        log.Fatalf("Failed to create OTLP exporter: %v", err)
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewSchemaless(
            semconv.ServiceNameKey.String("grafana-friends-meetup"),
        )),
    )

    otel.SetTracerProvider(tp)

    return func() {
        if err := tp.Shutdown(context.Background()); err != nil {
            log.Fatalf("Failed to shutdown tracer provider: %v", err)
        }
    }
}

func main() {
    shutdown := initTracer()
    defer shutdown()

    tracer := otel.Tracer("grafana-friends-meetup")

    ctx, span := tracer.Start(context.Background(), "main-function")
    defer span.End()

    startMeetup(ctx, tracer)
}

func startMeetup(ctx context.Context, tracer trace.Tracer) {
    _, span := tracer.Start(ctx, "startMeetup")
    defer span.End()

    log.Println("Starting the Grafana and Friends meetup...")
    time.Sleep(1 * time.Second)

    discussObservability(ctx, tracer)
    debugKubernetes(ctx, tracer) // Call the function that throws an error
}

func discussObservability(ctx context.Context, tracer trace.Tracer) {
    _, span := tracer.Start(ctx, "discussObservability")
    defer span.End()

    log.Println("Discussing observability and OpenTelemetry...")
    time.Sleep(1 * time.Second)

    span.AddEvent("O11y_demo", trace.WithAttributes(
        attribute.String("topic", "Tracing rocks"),
        attribute.String("feedback", "Works great"),
    ))

    log.Println("Observability discussion complete.")
    fmt.Println("Tracing complete.")
}

func debugKubernetes(ctx context.Context, tracer trace.Tracer) {
    _, span := tracer.Start(ctx, "debugKubernetes")
    defer span.End()

    log.Println("Attempting to load Kubernetes config...")

    // Try to open a non-existent configuration file to simulate a real Go error
    _, err := os.Open("non_existent_config.yaml")
    if err != nil {

        // Add the error as a span event
        span.AddEvent("Config file load failure", trace.WithAttributes(
            attribute.String("error.reason", err.Error()),
            attribute.String("file.name", "non_existent_config.yaml"),
        ))

        // Record the error and set span attributes
        span.RecordError(err)
        span.SetAttributes(
            attribute.Int("http.status_code", 400),
            attribute.String("http.status_text", "Bad Request"),
            attribute.String("otel.status_code", "ERROR"),
            attribute.String("error", err.Error()),
        )
    }

    log.Println("Execution complete.")
    time.Sleep(1 * time.Second)
}
