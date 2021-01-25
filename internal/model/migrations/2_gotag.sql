-- +migrate Up
-- MariaDB dump 10.17  Distrib 10.5.1-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: dev
-- ------------------------------------------------------
-- Server version	10.5.1-MariaDB-1:10.5.1+maria~bionic

/*!40101 SET @OLD_CHARACTER_SET_CLIENT = @@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS = @@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION = @@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE = @@TIME_ZONE */;
/*!40103 SET TIME_ZONE = '+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS = @@UNIQUE_CHECKS, UNIQUE_CHECKS = 0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;
/*!40101 SET @OLD_SQL_MODE = @@SQL_MODE, SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES = @@SQL_NOTES, SQL_NOTES = 0 */;

--
-- Table structure for table `file_tags`
--

DROP TABLE IF EXISTS `file_tags`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `file_tags`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `user_id`    int(11)  DEFAULT NULL,
    `file_id`    int(11)  DEFAULT NULL,
    `tag_id`     int(11)  DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `file_users`
--

DROP TABLE IF EXISTS `file_users`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
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
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `files`
--

DROP TABLE IF EXISTS `files`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `files`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `file_name`  varchar(255) CHARACTER SET utf8mb4  DEFAULT NULL,
    `file_path`  varchar(4096) CHARACTER SET utf8mb4 DEFAULT NULL,
    `size_b`     int(11)                             DEFAULT NULL,
    `mime_id`    int(11)                             DEFAULT NULL,
    `phash_a`    bigint(16)                          DEFAULT NULL,
    `phash_b`    bigint(16)                          DEFAULT NULL,
    `phash_c`    bigint(16)                          DEFAULT NULL,
    `phash_d`    bigint(16)                          DEFAULT NULL,
    `sha256`     char(64)                            DEFAULT NULL,
    `created_at` datetime                            DEFAULT NULL,
    `updated_at` datetime                            DEFAULT NULL,
    `deleted_at` datetime                            DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `logs`
--

DROP TABLE IF EXISTS `logs`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `logs`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `body`       text    NOT NULL,
    `created_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mimes`
--

DROP TABLE IF EXISTS `mimes`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mimes`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `mime`       varchar(255) DEFAULT NULL,
    `created_at` datetime     DEFAULT NULL,
    `updated_at` datetime     DEFAULT NULL,
    `deleted_at` datetime     DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tags`
--

DROP TABLE IF EXISTS `tags`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tags`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `tag`        varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
    `created_at` datetime                           DEFAULT NULL,
    `updated_at` datetime                           DEFAULT NULL,
    `deleted_at` datetime                           DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

/*!40103 SET TIME_ZONE = @OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE = @OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS = @OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT = @OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS = @OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION = @OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES = @OLD_SQL_NOTES */;

-- Dump completed on 2020-02-19 17:25:26


DROP
    FUNCTION IF EXISTS HAMMINGDISTANCE;

CREATE
    FUNCTION HAMMINGDISTANCE(A0 BIGINT, A1 BIGINT, A2 BIGINT, A3 BIGINT,
                             B0 BIGINT, B1 BIGINT, B2 BIGINT, B3 BIGINT)
    RETURNS INT DETERMINISTIC
    RETURN
            BIT_COUNT(A0 ^ B0) +
            BIT_COUNT(A1 ^ B1) +
            BIT_COUNT(A2 ^ B2) +
            BIT_COUNT(A3 ^ B3);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE people;
