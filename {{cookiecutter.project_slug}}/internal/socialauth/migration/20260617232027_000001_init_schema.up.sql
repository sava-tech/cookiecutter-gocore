CREATE TABLE IF NOT EXISTS social_auth_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    avatar_url TEXT,
    access_token TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_id)
);

CREATE INDEX idx_social_auth_users_user_id ON social_auth_users(user_id);
CREATE INDEX idx_social_auth_users_provider ON social_auth_users(provider, provider_id);
CREATE INDEX idx_social_auth_users_email ON social_auth_users(email);