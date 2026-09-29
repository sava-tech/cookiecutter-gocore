-- Drop functions if created
DROP FUNCTION IF EXISTS cleanup_expired_sessions();
DROP FUNCTION IF EXISTS cleanup_old_verifications(integer);

-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_email;
DROP INDEX IF EXISTS idx_sessions_refresh_token;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_user_active;

DROP INDEX IF EXISTS idx_verifications_identifier;
DROP INDEX IF EXISTS idx_verifications_code;
DROP INDEX IF EXISTS idx_verifications_type;
DROP INDEX IF EXISTS idx_verifications_expired_at;
DROP INDEX IF EXISTS idx_verifications_used;
DROP INDEX IF EXISTS idx_verifications_lookup;
DROP INDEX IF EXISTS idx_verifications_cleanup;

-- Drop tables
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS verifications;

DROP TABLE IF EXISTS "users";