CREATE TABLE IF NOT EXISTS `forbid` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `unique_id` varchar(40) NOT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4