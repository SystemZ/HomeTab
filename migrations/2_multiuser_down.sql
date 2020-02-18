ALTER TABLE files ADD user_id INT DEFAULT NULL AFTER id;

ALTER TABLE file_tags DROP user_id;

DROP TABLE file_users;

DROP TABLE users;