ALTER TABLE `certificates` CHANGE `common_name` `common_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '';

CREATE INDEX `common_name_idx` ON `certificates`(`common_name`) USING BTREE;

CREATE INDEX `revoked_at_idx` ON `certificates`(`revoked_at`) USING BTREE;

CREATE INDEX `expiry_idx` ON `certificates`(`expiry`) USING BTREE;

CREATE INDEX `not_before_idx` ON `certificates`(`not_before`) USING BTREE;

CREATE INDEX `ca_label_idx` ON `certificates`(`ca_label`) USING BTREE;

CREATE INDEX `status_idx` ON `certificates`(`status`) USING BTREE;

CREATE INDEX `unique_id_idx` ON `forbid`(`unique_id`) USING BTREE;

CREATE INDEX `deleted_at_idx` ON `forbid`(`deleted_at`) USING BTREE;

CREATE INDEX `name_idx` ON `self_keypair`(`name`) USING BTREE;