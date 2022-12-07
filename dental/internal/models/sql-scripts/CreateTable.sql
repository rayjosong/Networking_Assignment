CREATE TABLE users (
  uid INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	username varchar(100) NOT NULL UNIQUE,
  password varchar(100) NOT NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  role varchar(30) NOT NULL,
  created_at TIMESTAMP NOT NULL
);
