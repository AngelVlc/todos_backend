ALTER TABLE `lists` ADD `itemsCount` int(32) NOT NULL DEFAULT 0;

UPDATE lists SET itemsCount = (SELECT COUNT(li.id) FROM listItems li WHERE li.listId = lists .id);
