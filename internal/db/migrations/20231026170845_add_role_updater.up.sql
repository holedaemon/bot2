BEGIN;

CREATE TABLE role_updater (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,

    guild_id text NOT NULL,

    do_updates boolean NOT NULL DEFAULT 'false',

    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),

    UNIQUE(guild_id)
);

COMMIT;