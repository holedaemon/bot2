BEGIN;

ALTER TABLE guilds ADD COLUMN quotes_required_reactions int;

COMMIT;