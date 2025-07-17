-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "username" VARCHAR NOT NULL UNIQUE,
  "password_hash" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

-- Create providers table
CREATE TABLE "providers" (
  "id" UUID DEFAULT uuid_generate_v4(),
  "user_id" UUID NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "base_url" VARCHAR(255) NOT NULL,
  "type" VARCHAR(255) NOT NULL,
  "deleted_at" TIMESTAMPTZ,
  PRIMARY KEY (id, user_id),
  CONSTRAINT fk_providers_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX ON "providers" ("user_id");

-- Create connections table
CREATE TABLE "connections" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" UUID NOT NULL,
  "provider_id" VARCHAR NOT NULL,
  "encrypted_api_key" TEXT NOT NULL,
  "name" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "deleted_at" TIMESTAMPTZ,
  CONSTRAINT connections_user_id_fkey FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
CREATE INDEX ON "connections" ("user_id");

-- Create models table
CREATE TABLE "models" (
  "id" UUID DEFAULT uuid_generate_v4(),
  "user_id" UUID NOT NULL,
  "connection_id" UUID,
  "proxy_model_id" VARCHAR(255) NOT NULL,
  "provider_model_id" VARCHAR(255) NOT NULL DEFAULT '',
  "thinking" BOOLEAN NOT NULL DEFAULT FALSE,
  "tools_usage" BOOLEAN NOT NULL DEFAULT FALSE,
  "price_input" DECIMAL(10, 8) NOT NULL DEFAULT 0.0,
  "price_output" DECIMAL(10, 8) NOT NULL DEFAULT 0.0,
  "deleted_at" TIMESTAMPTZ,
  PRIMARY KEY (id, user_id),
  CONSTRAINT fk_models_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT models_connection_id_fkey FOREIGN KEY (connection_id) REFERENCES connections(id) ON DELETE CASCADE
);
CREATE INDEX ON "models" ("user_id");

-- Update existing rows with a default value for provider_model_id
UPDATE "models"
SET
  "provider_model_id" = "proxy_model_id";

-- Create api_keys table
CREATE TABLE "api_keys" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" UUID NOT NULL,
  "key_hash" VARCHAR NOT NULL UNIQUE,
  "name" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "last_used_at" TIMESTAMPTZ,
  CONSTRAINT api_keys_user_id_fkey FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
CREATE INDEX ON "api_keys" ("user_id");

-- Create conversation_logs table
CREATE TABLE "conversation_logs" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_id" UUID NOT NULL,
    "model_id" UUID NOT NULL,
    "connection_id" UUID,
    "request_payload" JSONB NOT NULL,
    "response_payload" JSONB NOT NULL,
    "prompt_tokens" BIGINT,
    "completion_tokens" BIGINT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT conversation_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_connection FOREIGN KEY (connection_id) REFERENCES connections(id) ON DELETE CASCADE
);