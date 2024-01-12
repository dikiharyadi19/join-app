CREATE TABLE `user_languages` (
    `user_uid` varchar(100) NOT NULL,
    `language` varchar(100) NOT NULL,
    `is_deleted` boolean NOT NULL DEFAULT 0,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_at` datetime NOT NULL DEFAULT now(),
    PRIMARY KEY (`user_uid`, `language`),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;