-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_ingredients" table
CREATE TABLE `new_ingredients` (`id` integer NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Copy rows from old table "ingredients" to new temporary table "new_ingredients"
INSERT INTO `new_ingredients` (`id`, `name`) SELECT `id`, `name` FROM `ingredients`;
-- Drop "ingredients" table after copying rows
DROP TABLE `ingredients`;
-- Rename temporary table "new_ingredients" to "ingredients"
ALTER TABLE `new_ingredients` RENAME TO `ingredients`;
-- Create "nutrients" table
CREATE TABLE `nutrients` (`id` integer NULL, `name` text NOT NULL, `unit` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "nutrients_name" to table: "nutrients"
CREATE UNIQUE INDEX `nutrients_name` ON `nutrients` (`name`);
-- Create "ingredient_nutrients" table
CREATE TABLE `ingredient_nutrients` (`ingredient_id` integer NOT NULL, `nutrient_id` integer NOT NULL, `amount` real NOT NULL, PRIMARY KEY (`ingredient_id`, `nutrient_id`), CONSTRAINT `0` FOREIGN KEY (`nutrient_id`) REFERENCES `nutrients` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredients` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
