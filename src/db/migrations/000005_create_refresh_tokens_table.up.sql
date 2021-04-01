CREATE TABLE `refresh_tokens` (
    `id` int(32) NOT NULL AUTO_INCREMENT,
    `userId` int(32) NOT NULL,
    `refreshToken` varchar(250) NOT NULL,
    `expirationDate` timestamp NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_refresh_token` (`refreshToken`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
