CREATE TABLE `jwk` (
    `id` varchar(255) NOT NULL,
    `key` text NOT NULL,
    `expired_at` datetime NOT NULL,
    `is_deleted` boolean NOT NULL DEFAULT 0,
    `created_at` datetime NOT NULL DEFAULT now(),
    PRIMARY KEY (`id`, `expired_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;