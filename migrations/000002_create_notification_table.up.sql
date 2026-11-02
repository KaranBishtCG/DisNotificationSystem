CREATE TYPE notification_status AS ENUM(
    'PENDING',
    'PROCESSING',
    'SENT',
    'FAILED'
);

CREATE TYPE notification_channel AS ENUM(
    'EMAIL',
    'PUSH',
    'SMS'
);

CREATE TABLE IF NOT EXISTS notifications(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    recipient TEXT NOT NULL,
    channel notification_channel NOT NULL,
    status notification_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);