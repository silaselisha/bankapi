ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "user_name_currency_key";
ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "accounts_first_name_fkey ";

ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "accounts_last_name_fkey ";

DROP TABLE IF EXISTS "users" CASCADE;

