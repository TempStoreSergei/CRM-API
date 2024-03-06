-- CreateTable
CREATE TABLE "Almanac" (
    "id" TEXT NOT NULL,
    "user" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "body" TEXT NOT NULL,
    "startTime" TIMESTAMP(3) NOT NULL,
    "endTime" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Almanac_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "Almanac_user_idx" ON "Almanac"("user");
