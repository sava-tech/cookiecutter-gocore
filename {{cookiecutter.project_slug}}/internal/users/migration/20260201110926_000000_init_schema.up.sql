-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2026-02-01T11:10:27.718Z

CREATE TABLE "users" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email_verified" bool NOT NULL DEFAULT false,
  "account_type" varchar NOT NULL DEFAULT 'user',
  "phone_number" varchar UNIQUE NOT NULL,
  "avatar" varchar NOT NULL,
  "age" integer NOT NULL DEFAULT 1,
  "gender" varchar NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "hashed_password" varchar(255) NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "updated_at" timestamptz NOT NULL DEFAULT 'now()',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);


-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- SESSIONS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    refresh_token TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    client_ip VARCHAR(45) NOT NULL,
    is_blocked BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sessions indexes
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_email ON sessions(email);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Composite index for common queries
CREATE INDEX idx_sessions_user_active ON sessions(user_id, is_blocked, expires_at);

-- ============================================
-- VERIFICATIONS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(255) NOT NULL,
    identifier_type VARCHAR(50) NOT NULL CHECK (identifier_type IN ('email', 'phone')),
    identifier VARCHAR(255) NOT NULL,
    verification_type VARCHAR(50) NOT NULL ,
    expired_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CHECK (verification_type IN ('registration', 'password_reset', 'email_change'))

-- Verifications indexes
CREATE INDEX idx_verifications_identifier ON verifications(identifier);
CREATE INDEX idx_verifications_code ON verifications(code);
CREATE INDEX idx_verifications_type ON verifications(verification_type);
CREATE INDEX idx_verifications_expired_at ON verifications(expired_at);
CREATE INDEX idx_verifications_used ON verifications(used);

-- Composite indexes for performance
CREATE INDEX idx_verifications_lookup ON verifications(identifier, code, verification_type, used, expired_at);
CREATE INDEX idx_verifications_cleanup ON verifications(created_at, used, expired_at);

-- ============================================
-- OPTIONAL: Foreign key constraints
-- ============================================

-- If you have a users table, uncomment this:
-- ALTER TABLE sessions 
--     ADD CONSTRAINT fk_sessions_user 
--     FOREIGN KEY (user_id) 
--     REFERENCES users(id) 
--     ON DELETE CASCADE;

-- ============================================
-- OPTIONAL: Clean up old records (scheduled job)
-- ============================================

-- Create a function to clean up expired sessions
CREATE OR REPLACE FUNCTION cleanup_expired_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM sessions WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Create a function to clean up old verifications
CREATE OR REPLACE FUNCTION cleanup_old_verifications(days integer DEFAULT 7)
RETURNS void AS $$
BEGIN
    DELETE FROM verifications 
    WHERE created_at < NOW() - (days || ' days')::interval
    AND used = true;
END;
$$ LANGUAGE plpgsql;


COMMENT ON COLUMN "users"."email" IS 'the email of the registered user';

COMMENT ON COLUMN "users"."first_name" IS 'the first name of the registered user';

COMMENT ON COLUMN "users"."last_name" IS 'the last name of the registered user';

COMMENT ON COLUMN "users"."full_name" IS 'the full name of the registered user';

COMMENT ON COLUMN "users"."email_verified" IS 'this indicates if the email is verified or not';

COMMENT ON COLUMN "users"."account_type" IS 'the account_type of the registered user';

COMMENT ON COLUMN "users"."phone_number" IS 'the user phone number';

COMMENT ON COLUMN "users"."avatar" IS 'user game avatar';

COMMENT ON COLUMN "users"."age" IS 'the user actual age';

COMMENT ON COLUMN "users"."gender" IS 'the user gender';

COMMENT ON COLUMN "users"."is_active" IS 'this indicates if the account is active or not';

COMMENT ON COLUMN "users"."hashed_password" IS 'hashed user password';

COMMENT ON COLUMN "users"."password_changed_at" IS 'password changed date';

COMMENT ON COLUMN "users"."updated_at" IS 'user account updated date';

COMMENT ON COLUMN "users"."created_at" IS 'user account created date';
