CREATE TABLE `counters` (
    `id` int(32) NOT NULL AUTO_INCREMENT,
    `name` varchar(10) NOT NULL,
    `value` int(32) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_counters_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;