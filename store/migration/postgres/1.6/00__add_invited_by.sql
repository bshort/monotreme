-- Add invited_by column to user table to track who invited each user
ALTER TABLE "user" ADD COLUMN invited_by INTEGER REFERENCES "user"(id);