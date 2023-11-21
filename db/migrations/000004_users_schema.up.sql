CREATE TABLE
    "users" (
        "username" varchar PRIMARY KEY,
        "fullname" varchar NOT NULL,
        "email" varchar UNIQUE NOT NULL,
        "password" varchar NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now()),
        "updated_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z')
    );

ALTER TABLE IF EXISTS "accounts"
ADD
    FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE IF EXISTS "accounts"
ADD
    CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");