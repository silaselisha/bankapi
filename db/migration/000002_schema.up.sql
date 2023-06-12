CREATE TABLE
    "users"(
        "id" bigserial PRIMARY key,
        "first_name" varchar NOT NULL UNIQUE,
        "last_name" varchar NOT NULL UNIQUE,
        "gender" varchar NOT NULL,
        "email" varchar UNIQUE NOT NULL,
        "password" varchar NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now()),
        "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:0Z'
    );

ALTER TABLE "accounts"
ADD
    FOREIGN KEY ("first_name") REFERENCES "users" ("first_name");
ALTER TABLE "accounts"
ADD
    FOREIGN KEY ("last_name") REFERENCES "users" ("last_name");

ALTER TABLE "accounts"
    ADD CONSTRAINT "user_name_currency_key" UNIQUE ("first_name", "last_name", "currency");