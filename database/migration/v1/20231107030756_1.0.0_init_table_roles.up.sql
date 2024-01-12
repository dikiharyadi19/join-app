CREATE TABLE `roles` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `uid` varchar(100) NOT NULL,
    `name` varchar(50) NOT NULL,
    `is_deleted` boolean NOT NULL DEFAULT 0,
    `created_by` varchar(100) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_by` varchar(100) NOT NULL,
    `updated_at` datetime NOT NULL DEFAULT now(),
    FULLTEXT KEY (`name`),
    PRIMARY KEY (`id`),
    UNIQUE KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;