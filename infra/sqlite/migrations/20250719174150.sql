-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_recipe_images" table
CREATE TABLE `new_recipe_images` (`id` integer NULL, `recipe_id` integer NOT NULL, `filename` text NOT NULL, `sort_order` integer NOT NULL, `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP), PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Copy rows from old table "recipe_images" to new temporary table "new_recipe_images"
INSERT INTO `new_recipe_images` (`id`, `recipe_id`, `sort_order`, `created_at`) SELECT `id`, `recipe_id`, `sort_order`, `created_at` FROM `recipe_images`;
-- Drop "recipe_images" table after copying rows
DROP TABLE `recipe_images`;
-- Rename temporary table "new_recipe_images" to "recipe_images"
ALTER TABLE `new_recipe_images` RENAME TO `recipe_images`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
