CREATE TABLE IF NOT EXISTS `self_keypair` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(40) NOT NULL,
    `private_key` text,
    `certificate` text,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4