
CREATE TABLE `modules` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `uid` VARCHAR(100) NULL DEFAULT NULL,
    `name` VARCHAR(255) NULL DEFAULT NULL,
    `description` text NULL DEFAULT NULL,
    `file` text NULL DEFAULT NULL,
    `is_deleted` boolean NULL DEFAULT 0,
    `created_by` VARCHAR(100) NULL DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_by` VARCHAR(100) NULL DEFAULT NULL,
    `updated_at` datetime NOT NULL DEFAULT now() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY (`uid`)
) COLLATE='utf8_general_ci' ENGINE=InnoDB;