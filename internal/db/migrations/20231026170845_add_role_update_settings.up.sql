BEGIN;

CREATE TABLE role_update_settings (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,

    guild_id text NOT NULL,

    do_updates boolean NOT NULL DEFAULT 'true',

    steam_user_id text NOT NULL,
    steam_app_id int NOT NULL,

    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),

    UNIQUE(guild_id)
);

COMMIT;