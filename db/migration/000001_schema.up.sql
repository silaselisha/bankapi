CREATE TABLE
    "accounts"(
        id bigserial PRIMARY KEY,
        first_name varchar NOT NULL,
        last_name varchar NOT NULL,
        gender varchar NOT NULL,
        balance bigint NOT NULL,
        currency varchar NOT NULL,
        created_at timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE
    "entries"(
        id bigserial PRIMARY KEY,
        account_id bigserial,
        amount bigint NOT NULL,
        created_at timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE
    "transfers"(
        id bigserial PRIMARY KEY,
        from_account_id bigserial,
        to_account_id bigserial,
        amount bigint NOT NULL,
        created_at timestamptz NOT NULL DEFAULT (now())
    );

CREATE INDEX ON "accounts" ("first_name");

CREATE INDEX ON "accounts" ("last_name");

CREATE INDEX ON "accounts" ("first_name", "last_name");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

ALTER TABLE "entries"
ADD
    FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
ADD
    FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
ADD
    FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE accounts ALTER COLUMN gender TYPE VARCHAR(255);

ALTER TABLE accounts ALTER COLUMN currency TYPE VARCHAR(255);