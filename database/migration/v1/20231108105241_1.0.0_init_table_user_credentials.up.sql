CREATE TABLE `user_credentials` (
    `user_uid` varchar(100) NOT NULL,
    `username` varchar(100) NOT NULL,
    `password` char(64) NOT NULL,
    `is_active` boolean NOT NULL DEFAULT 1,
    `created_at` datetime NOT NULL DEFAULT now(),
    PRIMARY KEY (`user_uid`, `password`),
    UNIQUE KEY (`username`),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;