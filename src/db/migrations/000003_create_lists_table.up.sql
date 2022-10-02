CREATE TABLE `lists` (
    `id` int(32) NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL,
    `userId` int(32) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_lists_name` (`name`, `userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `listItems` (
    `id` int(32) NOT NULL AUTO_INCREMENT,
    `listId` int(32) NOT NULL,
    `title` varchar(50) NOT NULL,
    `description` varchar(200) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
