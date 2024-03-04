-- Create a new database
CREATE DATABASE IF NOT EXISTS dbservice_test;

-- Create a new user and grant privileges
CREATE USER 'testuser'@'%' IDENTIFIED WITH mysql_native_password BY 'password';
GRANT ALL PRIVILEGES ON dbservice_test.* TO 'testuser'@'%';

FLUSH PRIVILEGES;
