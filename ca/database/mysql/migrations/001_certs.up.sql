CREATE TABLE IF NOT EXISTS `certificates` (
  `serial_number` varchar(128) NOT NULL,
  `authority_key_identifier` varchar(128) NOT NULL,
  `ca_label` varchar(128) DEFAULT NULL,
  `status` varchar(128) NOT NULL,
  `reason` int(11) DEFAULT NULL,
  `expiry` timestamp NULL DEFAULT NULL,
  `revoked_at` timestamp NULL DEFAULT NULL,
  `pem` text NOT NULL,
  `issued_at` timestamp NULL DEFAULT NULL,
  `not_before` timestamp NULL DEFAULT NULL,
  `metadata` json DEFAULT NULL,
  `sans` json DEFAULT NULL,
  `common_name` text,
  PRIMARY KEY (`serial_number`,`authority_key_identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;