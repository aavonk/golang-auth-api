
CREATE TABLE IF NOT EXISTS "users" (
    "id" TEXT NOT NULL,
    "first_name" VARCHAR(50) NOT NULL,
    "last_name" VARCHAR(80) NOT NULL,
    "email" TEXT NOT NULL UNIQUE,
    "email_confirmed" BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY ("id")
);