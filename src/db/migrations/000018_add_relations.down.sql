ALTER TABLE lists
DROP FOREIGN KEY fk_list_category_id;

ALTER TABLE lists
DROP FOREIGN KEY fk_list_user_id;

ALTER TABLE categories
DROP FOREIGN KEY fk_category_user_id;