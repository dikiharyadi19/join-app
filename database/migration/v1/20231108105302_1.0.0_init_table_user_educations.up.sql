CREATE TABLE `user_educations` (
    `user_uid` varchar(100) NOT NULL,
    `college` text NOT NULL,
    `degree` varchar(100) NOT NULL,
    `major` varchar(100) NOT NULL,
    `from` year NOT NULL,
    `to` year NOT NULL,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_at` datetime NOT NULL DEFAULT now(),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;