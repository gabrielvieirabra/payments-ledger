CREATE TABLE IF NOT EXISTS idempotency_keys (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(255) NOT NULL,
    method        VARCHAR(10)  NOT NULL,
    path          VARCHAR(512) NOT NULL,
    status_code   INTEGER      NOT NULL,
    response_body JSONB        NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    expires_at    TIMESTAMPTZ  NOT NULL DEFAULT now() + INTERVAL '24 hours'
);

CREATE UNIQUE INDEX idx_idempotency_keys_unique ON idempotency_keys (idempotency_key, method, path);
