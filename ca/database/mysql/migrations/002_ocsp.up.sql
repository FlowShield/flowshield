CREATE TABLE IF NOT EXISTS `ocsp_responses` (
  `serial_number` varchar(128) NOT NULL,
  `authority_key_identifier` varchar(128) NOT NULL,
  `body` text NOT NULL,
  `expiry` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`serial_number`,`authority_key_identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;