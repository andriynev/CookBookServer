CREATE TABLE `users` (
    `id` INT(11) auto_increment,
    `username` VARCHAR(255) NOT NULL,
    `password` blob NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username_users` (`username`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
