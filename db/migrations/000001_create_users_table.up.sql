CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `userName` varchar(10) NOT NULL,
    `passwordHash` varchar(100) NOT NULL,
    `isAdmin` tinyint NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_userName` (`userName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
