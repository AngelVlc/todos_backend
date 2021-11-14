ALTER TABLE users CHANGE is_admin isAdmin tinyint NOT NULL;
ALTER TABLE users CHANGE password_hash passwordHash varchar(100) NOT NULL;
