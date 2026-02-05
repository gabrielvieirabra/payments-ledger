CREATE TABLE IF NOT EXISTS accounts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner      VARCHAR(255) NOT NULL,
    balance    BIGINT       NOT NULL DEFAULT 0,
    currency   VARCHAR(3)   NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_accounts_owner ON accounts (owner);
CREATE INDEX idx_accounts_currency ON accounts (currency);
