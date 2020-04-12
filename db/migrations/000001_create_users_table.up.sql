CREATE TABLE `users` (
    `id` int(32) NOT NULL AUTO_INCREMENT,
    `name` varchar(10) NOT NULL,
    `password_hash` varchar(100) NOT NULL,
    `is_admin` tinyint NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
