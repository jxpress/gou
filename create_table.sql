CREATE TABLE if not exists `karma`
(
    `id`       BIGINT   NOT NULL,
    `date`     DATETIME NOT NULL,
    `giver`    TEXT DEFAULT '',
    `receiver` TEXT DEFAULT '',
    `count`    INT  DEFAULT '1',
    `channel`  TEXT DEFAULT ''
);
