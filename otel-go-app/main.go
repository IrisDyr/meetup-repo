package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"  // Renamed to avoid conflict
    "go.opentelemetry.io/otel/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer() func() {
    // Set up the OTLP exporter with a specific Collector endpoint
    exporter, err := otlptracehttp.New(context.Background(),
        otlptracehttp.WithEndpoint("localhost:4317"),  // Replace this if needed
        otlptracehttp.WithInsecure(),  // Disable TLS since this is local
    )
    if err != nil {
        log.Fatalf("Failed to create OTLP exporter: %v", err)
    }

    // Create the trace provider with the OTLP exporter
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewSchemaless(
            semconv.ServiceNameKey.String("grafana-friends-meetup"),
        )),
    )

    otel.SetTracerProvider(tp)

    // Return a function to clean up the tracer
    return func() {
        if err := tp.Shutdown(context.Background()); err != nil {
            log.Fatalf("Failed to shutdown tracer provider: %v", err)
        }
    }
}

func main() {
    // Initialize the tracer and defer its shutdown
    shutdown := initTracer()
    defer shutdown()

    // Get a tracer instance
    tracer := otel.Tracer("grafana-friends-meetup")

    ctx, span := tracer.Start(context.Background(), "main-function")
    defer span.End()

    // Call functions with instrumentation
    startMeetup(ctx, tracer)
}

func startMeetup(ctx context.Context, tracer trace.Tracer) {
    _, span := tracer.Start(ctx, "startMeetup")
    defer span.End()

    log.Println("Starting the Grafana and Friends meetup...")
    time.Sleep(2 * time.Second) // Simulate work

    // Call the discussObservability function with custom span event
    discussObservability(ctx, tracer)
}

func discussObservability(ctx context.Context, tracer trace.Tracer) {
    _, span := tracer.Start(ctx, "discussObservability")
    defer span.End()

    log.Println("Discussing observability and OpenTelemetry...")
    time.Sleep(1 * time.Second) // Simulate discussion

    // Add a custom event to the span
    span.AddEvent("O11y_demo", trace.WithAttributes(
        attribute.String("topic", "Tracing rocks"),
        attribute.String("feedback", "Works great"),
    ))

    log.Println("Observability discussion complete.")
    fmt.Println("Tracing complete.")
}
