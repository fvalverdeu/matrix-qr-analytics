# Cloud Observability

This project relies on Google Cloud native observability for Cloud Run services.
No external logging, metrics, tracing, or agent stack is required for this challenge.

## What Cloud Run Provides Automatically

Cloud Run already provides:

- Request logs in Cloud Logging for each HTTP request.
- Container stdout and stderr ingestion into Cloud Logging.
- Request metrics in Cloud Monitoring (including 5xx response counts/rates).

Because both services write to stdout/stderr, application logs are captured automatically.

## Service Separation in Logs

Use Cloud Run service name to distinguish logs:

- Public Go service: matrix-qr-api
- Private Node service: matrix-qr-statistics

In Cloud Logging, filter by:

- resource.type="cloud_run_revision"
- resource.labels.service_name="matrix-qr-api" or "matrix-qr-statistics"

## Useful gcloud Log Commands

Recent logs for Go service:

```bash
gcloud logging read \
  'resource.type="cloud_run_revision" AND resource.labels.service_name="matrix-qr-api"' \
  --project=matrix-qr-analytics-507313-s6 \
  --limit=100 \
  --format='value(timestamp,severity,textPayload,jsonPayload.message)'
```

Recent logs for Node service:

```bash
gcloud logging read \
  'resource.type="cloud_run_revision" AND resource.labels.service_name="matrix-qr-statistics"' \
  --project=matrix-qr-analytics-507313-s6 \
  --limit=100 \
  --format='value(timestamp,severity,textPayload,jsonPayload.message)'
```

Go request-level 5xx entries:

```bash
gcloud logging read \
  'resource.type="cloud_run_revision" AND resource.labels.service_name="matrix-qr-api" AND httpRequest.status>=500' \
  --project=matrix-qr-analytics-507313-s6 \
  --limit=100
```

## How Downstream Failures Appear

When matrix-qr-api returns STATISTICS_UNAVAILABLE (HTTP 502):

1. Check Go logs for downstream failure context from the statistics client.
2. Check Node request logs around the same timestamp:
   - If no corresponding Node request exists, likely upstream-to-downstream connectivity or auth invocation issue.
   - If Node request exists and returns 5xx, inspect Node application error logs for controller/service failure.

This preserves generic client-facing error envelopes while providing internal diagnostics in logs.

## Minimum Logging Policy

Do log:

- Layer/context for internal errors (handler/client/controller).
- Downstream status code and operation path.
- Error type for unexpected failures when useful.

Do not log:

- Raw matrix payloads (request or response).
- Authorization headers or ID tokens.
- Full request headers.
- Full downstream response bodies.
- Full downstream URLs.
- Full transport error strings when they may include unstable or sensitive detail.
- Secrets or credentials.

## Cloud Monitoring 5xx Alert (Minimum)

Target signal: abnormal HTTP 5xx responses on matrix-qr-api.

Recommended metric source:

- Metric type: `run.googleapis.com/request_count`
- Monitored resource: `cloud_run_revision`
- Minimum filters:
  - `resource.labels.service_name = "matrix-qr-api"`
  - `resource.labels.location = "southamerica-west1"`
  - `metric.labels.response_code_class = "5xx"`

Console path:

1. Cloud Monitoring -> Alerting -> Create policy.
2. Select metric for Cloud Run request count/rate.
3. Filter resource labels:
   - service_name = matrix-qr-api
   - location = southamerica-west1
4. Add response class/code filter for 5xx.
5. Condition example:
   - Trigger when 5xx request count is above 0 for 5 minutes,
     or above a small baseline rate suitable for this environment.
6. Configure notification channel and policy name.

Suggested initial policy name:

- matrix-qr-api-5xx-alert

After creation, verify policy by:

- Reviewing recent incidents.
- Correlating with Cloud Logging request and application logs.

## Operational Notes

- Public endpoint health and QR smoke tests in CD validate external behavior.
- Observability for diagnostics is primarily request logs + stderr/stdout application logs.
- This is intentionally minimal and dependency-free for challenge scope.
