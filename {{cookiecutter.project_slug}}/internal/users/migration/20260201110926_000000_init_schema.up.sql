-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2026-02-01T11:10:27.718Z

CREATE TABLE "users" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "avatar" varchar NOT NULL,
  "age" integer NOT NULL DEFAULT 1,
  "gender" varchar NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "updated_at" timestamptz NOT NULL DEFAULT 'now()',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

COMMENT ON COLUMN "users"."email" IS 'the email of the registered user';

COMMENT ON COLUMN "users"."username" IS 'the username of the registered user';

COMMENT ON COLUMN "users"."phone_number" IS 'the user phone number';

COMMENT ON COLUMN "users"."avatar" IS 'user game avatar';

COMMENT ON COLUMN "users"."age" IS 'the user actual age';

COMMENT ON COLUMN "users"."gender" IS 'the user gender';

COMMENT ON COLUMN "users"."is_active" IS 'this indicates if the account is active or not';

COMMENT ON COLUMN "users"."hashed_password" IS 'hashed user password';

COMMENT ON COLUMN "users"."password_changed_at" IS 'password changed date';

COMMENT ON COLUMN "users"."updated_at" IS 'user account updated date';

COMMENT ON COLUMN "users"."created_at" IS 'user account created date';
