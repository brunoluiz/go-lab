## To-do

- Review the tests, it is pretty much AI slop from the agent
- Support to Postgres
- Observability (metrics, logging, tracing)
- Re-think error handling middlewares
- Support to REST via OpenAPI

## Observability

This service is instrumented with OpenTelemetry for tracing and metrics. Configuration is done via environment variables.

### Environment Variables

- `OTEL_SERVICE_NAME`: Service name (default: auto-detected)
- `OTEL_SERVICE_VERSION`: Service version
- `OTEL_TRACES_EXPORTER`: Traces exporter (default: `otlp`)
- `OTEL_METRICS_EXPORTER`: Metrics exporter (default: `console`)
- `OTEL_RESOURCE_ATTRIBUTES`: Additional resource attributes

### Tracing

Traces are exported via OTLP to the OpenTelemetry Collector, which forwards them to Jaeger. To view traces:

1. Start dependencies (including Jaeger and OpenTelemetry Collector):

   ```bash
   docker-compose up -d
   ```

2. Access Jaeger UI at <http://localhost:16686>

3. Run the service:

   ```bash
   export OTEL_SERVICE_NAME=todo-service
   export OTEL_TRACES_EXPORTER=otlp
   export DB_DSN="postgres://todo_user:todo_pass@localhost:5432/todo?sslmode=disable"
   export OTEL_EXPORTER_OTLP_INSECURE=true
   ./cmd/connectrpc/connectrpc
   ```

4. Make requests to see traces.

### Metrics

Metrics are exported via OTLP to the collector by default, which logs them to stdout. They are sent periodically (every 3 seconds).
