CREATE TABLE
    "accounts" (
        "id" bigserial PRIMARY KEY,
        "owner" varchar NOT NULL,
        "balance" integer NOT NULL,
        "currency" varchar NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE
    "entries" (
        "id" bigserial PRIMARY KEY,
        "account_id" bigserial NOT NULL,
        "amount" integer,
        "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE
    "transfers" (
        "id" bigserial PRIMARY KEY,
        "from_account_id" bigserial NOT NULL,
        "to_account_id" bigserial NOT NULL,
        "amount" integer NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now())
    );

ALTER TABLE "entries"
ADD
    FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
ADD
    FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
ADD
    FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
