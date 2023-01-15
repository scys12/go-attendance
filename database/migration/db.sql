CREATE DATABASE attendance;
USE attendance;
set GLOBAL time_zone
= 'Asia/Jakarta';

CREATE TABLE users
(
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  full_name VARCHAR
(255) NOT NULL,
  password VARCHAR
(255) NOT NULL,
  email VARCHAR
(255) UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE attendances
(
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  check_in DATETIME,
  check_out DATETIME,
  attendance_date DATE,
  user_id int,
  FOREIGN KEY
(user_id) REFERENCES users
(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE activities
(
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  description VARCHAR
(255),
  attendance_id int,
  FOREIGN KEY
(attendance_id) REFERENCES attendances
(id),
user_id int,
  FOREIGN KEY
(user_id) REFERENCES users
(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);