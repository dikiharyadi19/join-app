CREATE TABLE `user_feedbacks` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `uid` varchar(100) NOT NULL,
    `mentor_uid` varchar(100) NOT NULL,
    `rating` decimal NOT NULL ,
    `comment` text NOT NULL,
    `created_at` datetime NOT NULL now(),
    `created_by` varchar(100) NOT NULL,
    `updated_at` datetime NOT NULL now(),
    `updated_by` varchar(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`uid`),
    UNIQUE KEY (`mentor_uid`, `created_by`),
    FOREIGN KEY (`mentor_uid`) REFERENCES users(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;