CREATE TYPE outbox_status AS ENUM(
    'PENDING',
    'PUBLISHED',
    'FAILED'
);

CREATE TABLE IF NOT EXISTS outbox_events(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL REFERENCES notifications(id),
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    status outbox_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    retry_count INT NOT NULL DEFAULT 0,
    published_at TIMESTAMPTZ
);
