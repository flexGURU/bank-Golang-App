-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2025-02-28T06:19:38.308Z

CREATE TABLE "users" (
  "username" VARCHAR PRIMARY KEY,
  "hashed_password" VARCHAR NOT NULL,
  "full_name" VARCHAR NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "is_verified" BOOLEAN NOT NULL DEFAULT false,
  "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00+00:00',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE "verify_emails" (
  "id" SERIAL PRIMARY KEY,
  "username" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "secret_code" VARCHAR NOT NULL,
  "is_used" BOOLEAN NOT NULL DEFAULT false,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "expires_at" TIMESTAMPTZ NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "account" (
  "id" SERIAL PRIMARY KEY,
  "owner" VARCHAR NOT NULL,
  "balance" INT NOT NULL,
  "currency" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE "entries" (
  "id" SERIAL PRIMARY KEY,
  "account_id" INT NOT NULL,
  "amount" INT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transfer" (
  "id" SERIAL PRIMARY KEY,
  "from_account_id" INT NOT NULL,
  "to_account_id" INT NOT NULL,
  "amount" INT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE "sessions" (
  "id" UUID PRIMARY KEY,
  "username" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "user_agent" VARCHAR NOT NULL,
  "client_ip" VARCHAR NOT NULL,
  "is_blocked" BOOLEAN NOT NULL DEFAULT false,
  "expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE UNIQUE INDEX "owner_currency_key" ON "account" ("owner", "currency");

CREATE INDEX "idx_account_owner" ON "account" ("owner");

CREATE INDEX "idx_entries_account_id" ON "entries" ("account_id");

CREATE INDEX "idx_transfer_from_account" ON "transfer" ("from_account_id");

CREATE INDEX "idx_transfer_to_account" ON "transfer" ("to_account_id");

CREATE INDEX "idx_transfer_from_to_account" ON "transfer" ("from_account_id", "to_account_id");

ALTER TABLE "account" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("from_account_id") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("to_account_id") REFERENCES "account" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
