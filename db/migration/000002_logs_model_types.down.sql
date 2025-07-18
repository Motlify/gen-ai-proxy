ALTER TABLE "logs" RENAME TO "conversation_logs";
ALTER TABLE "conversation_logs" DROP COLUMN "type";

ALTER TABLE "models" DROP COLUMN "type";
