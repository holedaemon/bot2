BEGIN;

ALTER TABLE roles DROP COLUMN role_name;

DROP TABLE quotes;

DROP TABLE guilds;

DROP TABLE discord_tokens;

DROP TABLE sessions;

COMMIT;