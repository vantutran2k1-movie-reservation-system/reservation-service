CREATE TABLE IF NOT EXISTS user_registration_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_value VARCHAR(255) NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    expires_at TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);

CREATE INDEX idx_user_registration_tokens_user_id ON user_registration_tokens(user_id);
CREATE INDEX idx_user_registration_tokens_token_value ON user_registration_tokens(token_value);
CREATE INDEX idx_user_registration_tokens_expires_at ON user_registration_tokens(expires_at);