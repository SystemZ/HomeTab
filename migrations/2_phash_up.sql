ALTER TABLE files
  ADD COLUMN phash char(64);

CREATE TABLE `files_distance` (
  `a_b`     TEXT UNIQUE,
  `distance` INTEGER
);

CREATE INDEX a_b
  ON files_distance (a_b);