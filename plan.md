# DisNotiSystem — Progress Plan

Distributed Notification System using the Transactional Outbox pattern with Kafka.

## 1. Setup & Infrastructure
- [x] Initialize Go module
- [x] Set up PostgreSQL (connection via pgx)
- [x] Config loader (`PORT`, `DATABASE_URL`)
- [x] Migration tool (`cmd/migrate`)
- [ ] Set up Kafka (Docker Compose or managed service)
- [ ] Add Kafka client dependency (e.g. `segmentio/kafka-go`)
- [ ] docker-compose.yml for local Postgres + Kafka

## 2. Database Schema
- [x] `users` table
- [x] `notifications` table (basic version)
- [x] `outbox_events` table (id, notification_id, event_type, payload JSONB, status, created_at, updated_at, retry_count)
- [ ] Align `notifications` table with design (add `updated_at`, review status enum)
- [ ] Indexes for efficient polling of `outbox_events` (status, created_at)

## 3. Notification API
- [x] `GET /healthz`
- [ ] `POST /notifications` endpoint
- [ ] Single DB transaction: insert into `notifications` + `outbox_events`
- [ ] Return `202 Accepted` on success

## 4. Outbox Worker
- [ ] New `cmd/outbox-worker` service
- [ ] Poll `outbox_events` for pending rows
- [ ] Publish event to corresponding Kafka topic
- [ ] Mark event as published (with retry + error handling)

## 5. Kafka Topics & Consumers
- [ ] Create topics: `notifications.email`, `notifications.sms`, `notifications.push`
- [ ] Define consumer groups per channel
- [ ] Handle message processing + offset commits

## 6. Notification Workers
- [ ] `cmd/email-worker` (+ provider integration, e.g. SendGrid)
- [ ] `cmd/sms-worker` (+ provider integration, e.g. Twilio)
- [ ] `cmd/push-worker` (+ provider integration, e.g. FCM)
- [ ] Retry logic + DLQ for failed sends

## 7. Observability & Resilience
- [ ] Structured logging
- [ ] Metrics (Prometheus/Grafana or similar)
- [ ] Retry/backoff handling
- [ ] Idempotency checks

## 8. Test & Deploy
- [ ] Unit tests
- [ ] Integration tests
- [ ] Load testing (optional)
- [ ] Dockerize services
- [ ] Deploy (Docker/Kubernetes)

---
### Next immediate step
Add `outbox_events` migration, then implement `POST /notifications` with transactional insert.
