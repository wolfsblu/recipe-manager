-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_recipe_ingredients" table
CREATE TABLE `new_recipe_ingredients` (
  `id` integer NULL,
  `step_id` integer NOT NULL,
  `ingredient_id` integer NOT NULL,
  `unit_id` integer NOT NULL,
  `amount` real NOT NULL,
  `sort_order` integer NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  CONSTRAINT `0` FOREIGN KEY (`unit_id`) REFERENCES `units` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT `1` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredients` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `2` FOREIGN KEY (`step_id`) REFERENCES `recipe_steps` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Copy rows from old table "recipe_ingredients" to new temporary table "new_recipe_ingredients"
INSERT INTO `new_recipe_ingredients` (`id`, `step_id`, `ingredient_id`, `unit_id`, `amount`, `sort_order`) SELECT `id`, `step_id`, `ingredient_id`, `unit_id`, `amount`, `sort_order` FROM `recipe_ingredients`;
-- Drop "recipe_ingredients" table after copying rows
DROP TABLE `recipe_ingredients`;
-- Rename temporary table "new_recipe_ingredients" to "recipe_ingredients"
ALTER TABLE `new_recipe_ingredients` RENAME TO `recipe_ingredients`;
-- Create index "recipe_ingredients_step_id_ingredient_id" to table: "recipe_ingredients"
CREATE UNIQUE INDEX `recipe_ingredients_step_id_ingredient_id` ON `recipe_ingredients` (`step_id`, `ingredient_id`);
-- Create index "recipe_ingredients_step_id_sort_order" to table: "recipe_ingredients"
CREATE UNIQUE INDEX `recipe_ingredients_step_id_sort_order` ON `recipe_ingredients` (`step_id`, `sort_order`);
-- Create index "idx_recipe_ingredients_sort_order" to table: "recipe_ingredients"
CREATE INDEX `idx_recipe_ingredients_sort_order` ON `recipe_ingredients` (`sort_order`);
-- Drop "user_roles" table
DROP TABLE `user_roles`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
