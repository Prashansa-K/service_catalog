--- Creating Database
CREATE DATABASE IF NOT EXISTS `servicecatalog`;

--- Creating Tables
CREATE TABLE IF NOT EXISTS `servicecatalog`.`service` (
  `id` SERIAL PRIMARY KEY,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL
  `version_count` INT NOT NULL DEFAULT 0,
)

CREATE TABLE IF NOT EXISTS `servicecatalog`.`versions` (
  `id` VARCHAR(255) PRIMARY KEY,
  `name` VARCHAR(255) NOT NULL,
  `service_id` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  UNIQUE (`name`, `service_id`),
  FOREIGN KEY (`service_id`) REFERENCES `servicecatalog`.`service`(`id`)
);
