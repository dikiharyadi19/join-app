CREATE TABLE `user_expertises` (
    `user_uid` varchar(100) NOT NULL,
    `expertise_uid` varchar(100) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_at` datetime NOT NULL DEFAULT now(),
    PRIMARY KEY (`user_uid`, `expertise_uid`),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`),
    FOREIGN KEY (`expertise_uid`) REFERENCES expertises(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;