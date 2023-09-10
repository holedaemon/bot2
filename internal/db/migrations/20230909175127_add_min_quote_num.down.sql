BEGIN;

ALTER TABLE guilds DROP COLUMN quotes_required_reactions;

COMMIT;