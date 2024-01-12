CREATE TABLE `user_work_experiences` (
    `user_uid` varchar(100) NOT NULL,
    `role` text NOT NULL,
    `company` text NOT NULL,
    `industry` varchar(255) NOT NULL,
    `start_date` date NOT NULL,
    `end_date` date NOT NULL,
    `description` text NOT NULL,
    `created_at` datetime NOT NULL DEFAULT now(),
    `updated_at` datetime NOT NULL DEFAULT now(),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;