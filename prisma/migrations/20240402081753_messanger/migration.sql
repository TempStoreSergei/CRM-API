-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "messenger";

-- CreateTable
CREATE TABLE "messenger"."chats" (
    "id" TEXT NOT NULL,
    "user_id" TEXT NOT NULL,
    "text" TEXT,
    "media_url" TEXT,
    "type" TEXT NOT NULL DEFAULT 'text',
    "reply_at" TEXT,
    "read_at" TIMESTAMP(3),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,
    "deleted_at" TIMESTAMP(3),

    CONSTRAINT "chats_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "chats_created_at_user_id_idx" ON "messenger"."chats"("created_at", "user_id");

-- AddForeignKey
ALTER TABLE "messenger"."chats" ADD CONSTRAINT "chats_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
