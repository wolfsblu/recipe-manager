-- Add column "locale" to table: "users"
ALTER TABLE `users` ADD COLUMN `locale` text NOT NULL DEFAULT 'en';
