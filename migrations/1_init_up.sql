CREATE TABLE `files` (
  `id`        INTEGER PRIMARY KEY AUTOINCREMENT,
  `last_path` VARCHAR(4096) NULL,
  `size`      INTEGER       NULL,
  `mime`      VARCHAR(32)   NULL,
  `md5`       VARCHAR(32)   NULL,
  `sha1`      VARCHAR(40)   NULL,
  `sha256`    VARCHAR(64)   NULL,
  `created`   DATE          NULL
);

CREATE TABLE `tags` (
  `id`   INTEGER PRIMARY KEY AUTOINCREMENT,
  `fid`  INTEGER,
  `name` VARCHAR(64)
);