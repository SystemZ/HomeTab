ALTER TABLE files
    DROP user_id;

ALTER TABLE file_users
    ADD user_id INT NULL AFTER id;

CREATE TABLE `file_users`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `file_id`    int(11)  DEFAULT NULL,
    `user_id`    int(11)  DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

CREATE TABLE users
(
    id                 int(10) UNSIGNED NOT NULL,
    username           varchar(24)      NOT NULL,
    email              varchar(254)     NOT NULL,
    hash               char(60)         NOT NULL COMMENT 'bcrypt.GenerateFromPassword',
    default_project_id int(10) UNSIGNED NOT NULL DEFAULT 0,
    created_at         datetime         NOT NULL,
    updated_at         datetime         NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;


ALTER TABLE users
    ADD PRIMARY KEY (id),
    ADD UNIQUE KEY username (username),
    ADD UNIQUE KEY email (email);


ALTER TABLE users
    MODIFY id int(10) UNSIGNED NOT NULL AUTO_INCREMENT;