-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "profile";

-- CreateTable
CREATE TABLE "profile"."profiles" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "project_id" TEXT,
    "display_name" TEXT,
    "avatar" TEXT,
    "status_message" TEXT,
    "phone_number" TEXT,
    "last_seen" TIMESTAMP(3),
    "is_online" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "profiles_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "profile"."profile_projects" (
    "profile_id" TEXT NOT NULL,
    "project_id" TEXT NOT NULL,

    CONSTRAINT "profile_projects_pkey" PRIMARY KEY ("profile_id","project_id")
);

-- CreateTable
CREATE TABLE "profile"."projects" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "task_id" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "projects_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "profile"."tasks" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "completed" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "tasks_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "profiles_userId_key" ON "profile"."profiles"("userId");

-- AddForeignKey
ALTER TABLE "profile"."profiles" ADD CONSTRAINT "profiles_userId_fkey" FOREIGN KEY ("userId") REFERENCES "auth"."users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "profile"."profiles" ADD CONSTRAINT "profiles_project_id_fkey" FOREIGN KEY ("project_id") REFERENCES "profile"."projects"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "profile"."profile_projects" ADD CONSTRAINT "profile_projects_profile_id_fkey" FOREIGN KEY ("profile_id") REFERENCES "profile"."profiles"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "profile"."profile_projects" ADD CONSTRAINT "profile_projects_project_id_fkey" FOREIGN KEY ("project_id") REFERENCES "profile"."projects"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "profile"."projects" ADD CONSTRAINT "projects_task_id_fkey" FOREIGN KEY ("task_id") REFERENCES "profile"."tasks"("id") ON DELETE SET NULL ON UPDATE CASCADE;
