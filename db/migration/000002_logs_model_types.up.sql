ALTER TABLE "conversation_logs" RENAME TO "logs";
ALTER TABLE "logs" ADD COLUMN "type" VARCHAR(255) NOT NULL DEFAULT 'llm';
UPDATE "logs" SET "type" = 'llm';

ALTER TABLE "models" ADD COLUMN "type" VARCHAR(255) NOT NULL DEFAULT 'llm';
UPDATE "models" SET "type" = 'llm';