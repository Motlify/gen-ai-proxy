DROP TABLE IF EXISTS "conversation_logs";
DROP TABLE IF EXISTS "api_keys";
DROP TABLE IF EXISTS "models";
DROP TABLE IF EXISTS "connections";
DROP TABLE IF EXISTS "providers";
DROP TABLE IF EXISTS "users";
DROP EXTENSION IF EXISTS "uuid-ossp";

ALTER TABLE "models"
DROP COLUMN "provider_model_id";

ALTER TABLE "models"
RENAME COLUMN "proxy_model_id" TO "name";