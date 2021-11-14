ALTER TABLE users CHANGE isAdmin is_admin tinyint NOT NULL;
ALTER TABLE users CHANGE passwordHash password_hash varchar(100) NOT NULL;