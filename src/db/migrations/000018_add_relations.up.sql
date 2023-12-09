ALTER TABLE lists
ADD CONSTRAINT fk_list_category_id
FOREIGN KEY (categoryId)
REFERENCES categories (id);

ALTER TABLE lists
ADD CONSTRAINT fk_list_user_id
FOREIGN KEY (userId)
REFERENCES users (id);

ALTER TABLE categories
ADD CONSTRAINT fk_category_user_id
FOREIGN KEY (userId)
REFERENCES users (id);