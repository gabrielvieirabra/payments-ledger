CREATE TABLE IF NOT EXISTS transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_account_id UUID   NOT NULL REFERENCES accounts (id),
    to_account_id   UUID   NOT NULL REFERENCES accounts (id),
    amount          BIGINT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_from_account_id ON transactions (from_account_id);
CREATE INDEX idx_transactions_to_account_id ON transactions (to_account_id);
CREATE INDEX idx_transactions_from_to ON transactions (from_account_id, to_account_id);
