-- phpMyAdmin SQL Dump
-- version 4.7.9
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Mar 26, 2018 at 02:34 PM
-- Server version: 5.7.21
-- PHP Version: 7.2.2

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT = @@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS = @@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION = @@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `dev`
--

-- --------------------------------------------------------

--
-- Table structure for table `groups`
--

CREATE TABLE `groups` (
  `id`         INT(11)                         NOT NULL,
  `name`       VARCHAR(128) CHARACTER SET utf8 NOT NULL,
  `updated_at` INT(11)                         NOT NULL,
  `created_at` INT(11)                         NOT NULL
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

-- --------------------------------------------------------

--
-- Table structure for table `instances`
--

CREATE TABLE `instances` (
  `id`              INT(11)    NOT NULL,
  `type_id`         TINYINT(4) NOT NULL,
  `url`             VARCHAR(128) DEFAULT NULL,
  `updated_at`      INT(11)      DEFAULT NULL,
  `created_at`      INT(11)      DEFAULT NULL,
  `creator_user_id` INT(11)      DEFAULT NULL
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1
  COMMENT ='External services access';

--
-- Triggers `instances`
--
DELIMITER $$
CREATE TRIGGER `instances`
  BEFORE UPDATE
  ON `instances`
  FOR EACH ROW
  BEGIN
    SET NEW.updated_at = UNIX_TIMESTAMP();
  END
$$
DELIMITER ;

-- --------------------------------------------------------

--
-- Table structure for table `instances_access`
--

CREATE TABLE `instances_access` (
  `id`               INT(11) UNSIGNED NOT NULL,
  `instance_id`      INT(11)          NOT NULL,
  `user_id`          INT(11) UNSIGNED NOT NULL,
  `group_id`         INT(11)          NOT NULL,
  `instance_user_id` INT(11) UNSIGNED DEFAULT NULL,
  `token`            TEXT,
  `updated_at`       INT(11) UNSIGNED DEFAULT NULL,
  `created_at`       INT(11) UNSIGNED DEFAULT NULL
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

-- --------------------------------------------------------

--
-- Table structure for table `tasks`
--

CREATE TABLE `tasks` (
  `id`               INT(11) NOT NULL,
  `group_id`         INT(11) NOT NULL,
  `instance_id`      INT(11) UNSIGNED                DEFAULT NULL,
  `instance_task_id` VARCHAR(16)                     DEFAULT NULL,
  `title`            VARCHAR(128) CHARACTER SET utf8 DEFAULT NULL,
  `updated_at`       INT(11)                         DEFAULT NULL,
  `created_at`       INT(11)                         DEFAULT NULL
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id`         INT(10) UNSIGNED NOT NULL,
  `username`   VARCHAR(16)      NOT NULL,
  `password`   VARCHAR(128) DEFAULT NULL,
  `updated_at` INT(11)      DEFAULT NULL,
  `created_at` INT(11)      DEFAULT NULL
)
  ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `groups`
--
ALTER TABLE `groups`
  ADD UNIQUE KEY `id` (`id`);

--
-- Indexes for table `instances`
--
ALTER TABLE `instances`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `instances_access`
--
ALTER TABLE `instances_access`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `tasks`
--
ALTER TABLE `tasks`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `unique_task_id` (`instance_id`, `instance_task_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `instances`
--
ALTER TABLE `instances`
  MODIFY `id` INT(11) NOT NULL AUTO_INCREMENT,
  AUTO_INCREMENT = 1;

--
-- AUTO_INCREMENT for table `instances_access`
--
ALTER TABLE `instances_access`
  MODIFY `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  AUTO_INCREMENT = 1;

--
-- AUTO_INCREMENT for table `tasks`
--
ALTER TABLE `tasks`
  MODIFY `id` INT(11) NOT NULL AUTO_INCREMENT,
  AUTO_INCREMENT = 1;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  AUTO_INCREMENT = 1;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT = @OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS = @OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION = @OLD_COLLATION_CONNECTION */;