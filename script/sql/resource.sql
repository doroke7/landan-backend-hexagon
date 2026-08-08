# ************************************************************
# Sequel Ace SQL dump
# 版本號： 20104
#
# https://sequel-ace.com/
# https://github.com/Sequel-Ace/Sequel-Ace
#
# 主機: 127.0.0.1 (MySQL 8.0.46)
# 資料庫: resource
# 產生時間: 2026-08-08 02:40:11 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# 傾印（Dump）資料表 tx-admin_users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tx-admin_users`;

CREATE TABLE `tx-admin_users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NOT NULL DEFAULT '2038-01-19 03:14:07',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `da` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `tx-admin_users` WRITE;
/*!40000 ALTER TABLE `tx-admin_users` DISABLE KEYS */;

INSERT INTO `tx-admin_users` (`id`, `name`, `password`, `created_at`, `updated_at`, `deleted_at`)
VALUES
	(1,'admin','55b764578f3e645474b770f25ed9eab0','2026-08-08 01:20:25','2026-08-08 01:20:25','2038-01-19 03:14:07');

/*!40000 ALTER TABLE `tx-admin_users` ENABLE KEYS */;
UNLOCK TABLES;


# 傾印（Dump）資料表 tx-games
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tx-games`;

CREATE TABLE `tx-games` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NOT NULL DEFAULT '2038-01-19 03:14:07',
  PRIMARY KEY (`id`),
  UNIQUE KEY `k` (`key`),
  KEY `n-da` (`name`,`deleted_at`),
  KEY `da` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 傾印（Dump）資料表 tx-table_record_logs
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tx-table_record_logs`;

CREATE TABLE `tx-table_record_logs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `game_id` int unsigned NOT NULL DEFAULT '0',
  `table_record_id` int unsigned NOT NULL DEFAULT '0',
  `state` tinyint NOT NULL DEFAULT '0',
  `text` json NOT NULL,
  `image` json NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NOT NULL DEFAULT '2038-01-19 03:14:07',
  PRIMARY KEY (`id`),
  KEY `tri-gi-s` (`table_record_id`,`game_id`,`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 傾印（Dump）資料表 tx-table_records
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tx-table_records`;

CREATE TABLE `tx-table_records` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '2026-0322-0077-0000-0000-0001' COMMENT '年-月日-桌號-局號-期號',
  `game_id` int unsigned NOT NULL DEFAULT '0',
  `table_id` int unsigned NOT NULL DEFAULT '0',
  `state` tinyint NOT NULL DEFAULT '0',
  `text` json NOT NULL,
  `image` json NOT NULL,
  `result` json NOT NULL,
  `started_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ended_at` timestamp NOT NULL DEFAULT '2038-01-19 03:14:07',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NOT NULL DEFAULT '2038-01-19 03:14:07',
  PRIMARY KEY (`id`),
  UNIQUE KEY `n` (`no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 傾印（Dump）資料表 tx-tables
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tx-tables`;

CREATE TABLE `tx-tables` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '2026-0322-0077-0000-0000-0001' COMMENT '年-月日-桌號-局號-期號',
  `game_id` int unsigned NOT NULL DEFAULT '0',
  `key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `state` tinyint unsigned NOT NULL DEFAULT '0',
  `description` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `result` json NOT NULL,
  `started_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ended_at` timestamp NOT NULL DEFAULT '2038-01-19 03:14:07',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NOT NULL DEFAULT '2038-01-19 03:14:07',
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`),
  KEY `da` (`deleted_at`),
  KEY `n-da` (`no`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
