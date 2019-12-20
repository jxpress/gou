CREATE TABLE if not exists `karma`
(
    `id`       INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    `date`     DATETIME DEFAULT CURRENT_TIMESTAMP,
    `giver`    TEXT     DEFAULT '',
    `receiver` TEXT     DEFAULT '',
    `count`    FLOAT    DEFAULT 1.0,
    `channel`  TEXT     DEFAULT ''
);

CREATE TABLE if not exists `user`
(
    `id` TEXT NOT NULL PRIMARY KEY,
    `name` TEXT,
    `display_name` TEXT,
    `team_id` TEXT,
    `image_url` TEXT
);
