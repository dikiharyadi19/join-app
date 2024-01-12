CREATE TABLE `user_socials` (
    `user_uid` varchar(100) NOT NULL,
    `title` varchar(100) NOT NULL,
    `link` text NOT NULL,
    `is_deleted` boolean NOT NULL DEFAULT 0,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_at` datetime NOT NULL DEFAULT now(),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;