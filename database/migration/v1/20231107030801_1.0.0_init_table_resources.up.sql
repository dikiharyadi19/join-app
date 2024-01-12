CREATE TABLE `resources` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `uid` varchar(100) NOT NULL,
    `name` varchar(255) NOT NULL,
    `parent_uid` varchar(100) DEFAULT NULL,
    `type` smallint(1) NOT NULL,
    `action` varchar(50) NOT NULL,
    `is_deleted` boolean NOT NULL DEFAULT 0,
    `created_by` varchar(100) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_by` varchar(100) NOT NULL,
    `updated_at` datetime NOT NULL DEFAULT now(),
    PRIMARY KEY (`id`),
    UNIQUE KEY (`uid`),
    FOREIGN KEY (`parent_uid`) REFERENCES resources(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;