-- create databases
CREATE DATABASE IF NOT EXISTS `delivery` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

GRANT Usage ON *.* TO 'delivery'@'%';
GRANT Alter ON `delivery`.* TO 'delivery'@'%';
GRANT Create ON `delivery`.* TO 'delivery'@'%';
GRANT Create view ON `delivery`.* TO 'delivery'@'%';
GRANT Delete ON `delivery`.* TO 'delivery'@'%';
GRANT Delete history ON `delivery`.* TO 'delivery'@'%';
GRANT Grant option ON `delivery`.* TO 'delivery'@'%';
GRANT Index ON `delivery`.* TO 'delivery'@'%';
GRANT Insert ON `delivery`.* TO 'delivery'@'%';
GRANT References ON `delivery`.* TO 'delivery'@'%';
GRANT Select ON `delivery`.* TO 'delivery'@'%';
GRANT Show view ON `delivery`.* TO 'delivery'@'%';
GRANT Trigger ON `delivery`.* TO 'delivery'@'%';
GRANT Update ON `delivery`.* TO 'delivery'@'%';
GRANT Execute ON `delivery`.* TO 'delivery'@'%';
FLUSH PRIVILEGES;

CREATE TABLE IF NOT EXISTS `delivery`.orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  distance INT UNSIGNED NOT NULL,
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  CONSTRAINT order_PK PRIMARY KEY (id)
)
ENGINE=InnoDB;