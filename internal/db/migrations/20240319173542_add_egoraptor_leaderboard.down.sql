BEGIN;

DROP TABLE egoraptor_mentions;

ALTER TABLE egoraptor_settings RENAME TO egoraptor_mentions;

COMMIT;