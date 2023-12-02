CREATE TABLE `categories` (
    `id` int(32) NOT NULL AUTO_INCREMENT,
    `name` varchar(12) NOT NULL,
    `description` varchar(500) NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;