CREATE DATABASE simple_bank;
\connect simple_bank;

CREATE TABLE IF NOT EXISTS "account" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,
    "balance" bigint NOT NULL,
    "currency" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "transfer" (
    "id" integer PRIMARY KEY,
    "from_account_id" bigint NOT NULL,
    "to_account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX IF NOT EXISTS "owner_idx" ON "account" ("owner");

CREATE INDEX IF NOT EXISTS "account_idx" ON "entries" ("account_id");

CREATE INDEX IF NOT EXISTS "from_account_idx" ON "transfer" ("from_account_id");

CREATE INDEX IF NOT EXISTS "to_account_idx" ON "transfer" ("to_account_id");

CREATE INDEX IF NOT EXISTS "from_account_idx_to_account_idx" ON "transfer" (
    "from_account_id", "to_account_id"
);

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfer"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY (
    "account_id"
) REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY (
    "from_account_id"
) REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY (
    "to_account_id"
) REFERENCES "account" ("id");
