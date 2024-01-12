CREATE TABLE `access` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `uid` varchar(100) NOT NULL,
    `role_uid` varchar(255) NOT NULL,
    `resource_uid` varchar(255) NOT NULL,
    `is_deleted` boolean NOT NULL DEFAULT 0,
    `created_by` varchar(100) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_by` varchar(100) NOT NULL,
    `updated_at` datetime NOT NULL DEFAULT now(),
    PRIMARY KEY (`id`),
    UNIQUE KEY (`role_uid`, `resource_uid`),
    FOREIGN KEY (`role_uid`) REFERENCES roles(`uid`),
    FOREIGN KEY (`resource_uid`) REFERENCES resources(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;