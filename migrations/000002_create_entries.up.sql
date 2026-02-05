CREATE TABLE IF NOT EXISTS entries (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID    NOT NULL REFERENCES accounts (id),
    amount     BIGINT  NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_entries_account_id ON entries (account_id);
