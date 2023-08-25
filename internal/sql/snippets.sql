-- Create a new UTF-8 `templatemaker` database.
CREATE DATABASE templatemaker CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Switch to using the `templatemaker` database.
USE templatemaker;

-- Create a `templates` table.
CREATE TABLE templates (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    subject TEXT NOT NULL,
    description TEXT NOT NULL,
    assessment TEXT NOT NULL,
    recommendation TEXT NOT NULL,
    query TEXT NULL,
    user_id CHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO templates (id, name, subject, description, assessment, recommendation, user_id) VALUES (
    '4af9f070-c1f84da2baccdde10aa0',
    'Name 1',
    'Subject 1',
    'Description 1',
    'Assessment 1',
    'Recommendation 1',
    '6433d953-7eb9d1bc674a2b55c3d0'
);