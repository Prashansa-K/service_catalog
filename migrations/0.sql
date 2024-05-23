--- Creating Database
CREATE DATABASE IF NOT EXISTS servicecatalog;

--- Creating Tables
CREATE TABLE IF NOT EXISTS services (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  version_count INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS versions (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  service_id int NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  UNIQUE (name, service_id),
  FOREIGN KEY (service_id) REFERENCES services(id)
);
