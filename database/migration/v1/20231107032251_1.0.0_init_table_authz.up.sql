CREATE TABLE `authz` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `uid` varchar(100) NOT NULL,
    `user_uid` varchar(100) NOT NULL,
    `role_uid` varchar(100) NOT NULL,
    `is_deleted` boolean NOT NULL DEFAULT 0,
    `last_active` datetime NOT NULL,
    `created_by` varchar(100) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_by` varchar(100) NOT NULL,
    `updated_at` datetime NOT NULL DEFAULT now(),
    PRIMARY KEY (`id`, `user_uid`, `role_uid`),
    UNIQUE KEY (`uid`),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`),
    FOREIGN KEY (`role_uid`) REFERENCES roles(`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;