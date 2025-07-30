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
-- Drop "user_awards" table
DROP TABLE `user_awards`;
-- Drop "user_reputation" table
DROP TABLE `user_reputation`;
-- Drop "recipe_votes" table
DROP TABLE `recipe_votes`;
-- Drop "actions" table
DROP TABLE `actions`;
-- Drop "awards" table
DROP TABLE `awards`;
-- Create index "meal_plan_date_sort_order" to table: "meal_plan"
CREATE UNIQUE INDEX `meal_plan_date_sort_order` ON `meal_plan` (`date`, `sort_order`);
-- Create index "idx_meal_plan_sort_order" to table: "meal_plan"
CREATE INDEX `idx_meal_plan_sort_order` ON `meal_plan` (`sort_order`);
-- Create "recipe_ingredients" table
CREATE TABLE `recipe_ingredients` (`id` integer NULL, `step_id` integer NOT NULL, `ingredient_id` integer NOT NULL, `unit_id` integer NOT NULL, `amount` real NOT NULL, `sort_order` integer NOT NULL DEFAULT 0, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`unit_id`) REFERENCES `units` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT, CONSTRAINT `1` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredients` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `2` FOREIGN KEY (`step_id`) REFERENCES `recipes` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "recipe_ingredients_step_id_ingredient_id" to table: "recipe_ingredients"
CREATE UNIQUE INDEX `recipe_ingredients_step_id_ingredient_id` ON `recipe_ingredients` (`step_id`, `ingredient_id`);
-- Create index "recipe_ingredients_step_id_sort_order" to table: "recipe_ingredients"
CREATE UNIQUE INDEX `recipe_ingredients_step_id_sort_order` ON `recipe_ingredients` (`step_id`, `sort_order`);
-- Create index "idx_recipe_ingredients_sort_order" to table: "recipe_ingredients"
CREATE INDEX `idx_recipe_ingredients_sort_order` ON `recipe_ingredients` (`sort_order`);
-- Create "recipe_steps" table
CREATE TABLE `recipe_steps` (`id` integer NULL, `recipe_id` integer NOT NULL, `instructions` text NOT NULL, `sort_order` integer NOT NULL DEFAULT 0, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "recipe_steps_recipe_id_sort_order" to table: "recipe_steps"
CREATE UNIQUE INDEX `recipe_steps_recipe_id_sort_order` ON `recipe_steps` (`recipe_id`, `sort_order`);
-- Create index "idx_recipe_steps_sort_order" to table: "recipe_steps"
CREATE INDEX `idx_recipe_steps_sort_order` ON `recipe_steps` (`sort_order`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
