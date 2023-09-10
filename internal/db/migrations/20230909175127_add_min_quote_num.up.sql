BEGIN;

ALTER TABLE guilds ADD COLUMN quotes_required_reactions int DEFAULT 1;

COMMIT;