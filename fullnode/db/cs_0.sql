-- +migrate Up

-- ----------------------------
-- Table structure for zta_client
-- ----------------------------
CREATE TABLE IF NOT EXISTS `zta_client` (
    `id` integer PRIMARY KEY  AUTOINCREMENT,
    `user_uuid` text(40)  NOT NULL DEFAULT '',
    `name` text(100)  NOT NULL DEFAULT '',
    `server_id` integer(20) NOT NULL DEFAULT '0',
    `uuid` text(40)  NOT NULL DEFAULT '',
    `port` integer(10) NOT NULL DEFAULT '0',
    `expire` integer(11) NOT NULL DEFAULT '0',
    `relay` blob COMMENT 'relay',
    `server` blob COMMENT 'server',
    `target` blob COMMENT 'target',
    `ca_pem` text(4000)  NOT NULL DEFAULT '',
    `cert_pem` text(3000)  NOT NULL DEFAULT '',
    `key_pem` text(4000)  NOT NULL DEFAULT '',
    `created_at` integer(20) NOT NULL DEFAULT '0',
    `updated_at` integer(20) NOT NULL DEFAULT '0'
);

-- ----------------------------
-- Table structure for zta_oauth2
-- ----------------------------
CREATE TABLE IF NOT EXISTS `zta_oauth2` (
  `id` integer  PRIMARY KEY  AUTOINCREMENT,
  `company` text(100)  NOT NULL DEFAULT '',
  `client_id` text(255)  NOT NULL DEFAULT '',
  `client_secret` text(255)  NOT NULL DEFAULT '',
  `redirect_url` text(255)  NOT NULL DEFAULT '',
  `scopes` blob,
  `auth_url` text(255)  NOT NULL DEFAULT '' ,
  `token_url` text(255)  NOT NULL DEFAULT '',
  `created_at` integer(20) NOT NULL DEFAULT '0',
  `updated_at` integer(20) NOT NULL DEFAULT '0'
);

INSERT INTO `zta_oauth2` (`id`, `company`, `client_id`, `client_secret`, `redirect_url`, `scopes`, `auth_url`, `token_url`, `created_at`, `updated_at`) VALUES (1, 'github', 'client_id', 'client_secret', 'http://your_domain/api/v1/user/oauth2/callback/github', '["user"]', '', '', '1656380163', '1656380163');

-- ----------------------------
-- Table structure for zta_relay
-- ----------------------------
CREATE TABLE IF NOT EXISTS `zta_relay` (
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `user_uuid` text(40) NOT NULL DEFAULT '',
  `name` text(100) NOT NULL DEFAULT '',
  `uuid` text(40) NOT NULL DEFAULT '',
  `host` text(255) NOT NULL DEFAULT '',
  `port` integer(10) NOT NULL DEFAULT '0',
  `out_port` integer(10) NOT NULL DEFAULT '0',
  `ca_pem` text(4000) NOT NULL DEFAULT '',
  `cert_pem` text(3000) NOT NULL DEFAULT '',
  `key_pem` text(4000) NOT NULL DEFAULT '',
  `created_at` integer(20) NOT NULL DEFAULT '0',
  `updated_at` integer(20) NOT NULL DEFAULT '0'
);

-- ----------------------------
-- Table structure for zta_resource
-- ----------------------------
CREATE TABLE IF NOT EXISTS `zta_resource` (
  `id` integer  PRIMARY KEY  AUTOINCREMENT,
  `user_uuid` text(40)  NOT NULL DEFAULT '',
  `name` text(100)  NOT NULL DEFAULT '',
  `uuid` text(40)  NOT NULL DEFAULT '',
  `type` text(40)  NOT NULL DEFAULT '',
  `host` text(255)  NOT NULL DEFAULT '',
  `port` text(255)  NOT NULL DEFAULT '',
  `created_at` integer(20) NOT NULL DEFAULT '0',
  `updated_at` integer(20) NOT NULL DEFAULT '0'
);

-- ----------------------------
-- Table structure for zta_server
-- ----------------------------
CREATE TABLE IF NOT EXISTS `zta_server` (
  `id` integer  PRIMARY KEY  AUTOINCREMENT,
  `user_uuid` text(40)  NOT NULL DEFAULT '',
  `resource_id` text(100)  NOT NULL DEFAULT '',
  `name` text(100)  NOT NULL DEFAULT '',
  `uuid` text(40)  NOT NULL DEFAULT '',
  `host` text(255)  NOT NULL DEFAULT '',
  `port` integer(10) NOT NULL DEFAULT '0',
  `out_port` integer(10) NOT NULL DEFAULT '0',
  `ca_pem` text(4000)  NOT NULL DEFAULT '',
  `cert_pem` text(3000)  NOT NULL DEFAULT '',
  `key_pem` text(4000)  NOT NULL DEFAULT '',
  `created_at` integer(20) NOT NULL DEFAULT '0',
  `updated_at` integer(20) NOT NULL DEFAULT '0'
);

-- ----------------------------
-- Table structure for zta_user
-- ----------------------------
CREATE TABLE IF NOT EXISTS `zta_user` (
  `id` integer  PRIMARY KEY  AUTOINCREMENT,
  `email` text(100)  NOT NULL DEFAULT '',
  `avatar_url` text(100)  NOT NULL DEFAULT '',
  `uuid` text(40)  NOT NULL DEFAULT '',
  `created_at` integer(20) NOT NULL DEFAULT '0',
  `updated_at` integer(20) NOT NULL DEFAULT '0'
);
-- +migrate Down
