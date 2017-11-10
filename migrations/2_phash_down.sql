CREATE TABLE t1_backup AS
  SELECT
    id,
    last_path,
    size,
    mime,
    md5,
    sha1,
    sha256,
    created
  FROM files;
DROP TABLE files;
ALTER TABLE t1_backup
  RENAME TO files;
/* FIXME add index and autoincrement */
CREATE INDEX id
  ON files (id);

DROP TABLE files_distance;