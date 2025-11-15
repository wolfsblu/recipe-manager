-- Create "ingredient" table
CREATE TABLE `ingredient` (`id` integer NOT NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Create "nutrient" table
CREATE TABLE `nutrient` (`id` integer NOT NULL, `name` text NOT NULL, `unit` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "nutrient_name" to table: "nutrient"
CREATE UNIQUE INDEX `nutrient_name` ON `nutrient` (`name`);
-- Create "ingredient_nutrient" table
CREATE TABLE `ingredient_nutrient` (`ingredient_id` integer NOT NULL, `nutrient_id` integer NOT NULL, `amount` real NOT NULL, PRIMARY KEY (`ingredient_id`, `nutrient_id`), CONSTRAINT `0` FOREIGN KEY (`nutrient_id`) REFERENCES `nutrient` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredient` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "permission" table
CREATE TABLE `permission` (`id` integer NOT NULL, `slug` text NOT NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "permission_slug" to table: "permission"
CREATE UNIQUE INDEX `permission_slug` ON `permission` (`slug`);
-- Create "recipe" table
CREATE TABLE `recipe` (`id` integer NOT NULL, `name` text NOT NULL, `servings` integer NOT NULL, `minutes` integer NOT NULL, `description` text NOT NULL, `created_by` integer NOT NULL, `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP), PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_recipes_created_at_desc" to table: "recipe"
CREATE INDEX `idx_recipes_created_at_desc` ON `recipe` (`created_at` DESC, `id` DESC);
-- Create index "idx_recipes_created_at_asc" to table: "recipe"
CREATE INDEX `idx_recipes_created_at_asc` ON `recipe` (`created_at`, `id`);
-- Create index "idx_recipes_name_asc" to table: "recipe"
CREATE INDEX `idx_recipes_name_asc` ON `recipe` (`name`, `id`);
-- Create index "idx_recipes_name_desc" to table: "recipe"
CREATE INDEX `idx_recipes_name_desc` ON `recipe` (`name` DESC, `id` DESC);
-- Create index "idx_recipes_servings_asc" to table: "recipe"
CREATE INDEX `idx_recipes_servings_asc` ON `recipe` (`servings`, `id`);
-- Create index "idx_recipes_servings_desc" to table: "recipe"
CREATE INDEX `idx_recipes_servings_desc` ON `recipe` (`servings` DESC, `id` DESC);
-- Create "role" table
CREATE TABLE `role` (`id` integer NOT NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "role_name" to table: "role"
CREATE UNIQUE INDEX `role_name` ON `role` (`name`);
-- Create "role_permission" table
CREATE TABLE `role_permission` (`role_id` integer NOT NULL, `permission_id` integer NOT NULL, CONSTRAINT `0` FOREIGN KEY (`permission_id`) REFERENCES `permission` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "tag" table
CREATE TABLE `tag` (`id` integer NOT NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "tag_name" to table: "tag"
CREATE UNIQUE INDEX `tag_name` ON `tag` (`name`);
-- Create index "idx_tags_name" to table: "tag"
CREATE INDEX `idx_tags_name` ON `tag` (`name`, `id`);
-- Create "unit" table
CREATE TABLE `unit` (`id` integer NOT NULL, `symbol` text NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "idx_units_name" to table: "unit"
CREATE INDEX `idx_units_name` ON `unit` (`name`, `id`);
-- Create "user" table
CREATE TABLE `user` (`id` integer NOT NULL, `email` text NOT NULL, `password_hash` text NOT NULL, `is_confirmed` boolean NOT NULL DEFAULT 0, `role_id` integer NOT NULL, `locale` text NOT NULL DEFAULT 'en', `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP), PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "user_email" to table: "user"
CREATE UNIQUE INDEX `user_email` ON `user` (`email`);
-- Create "user_registration" table
CREATE TABLE `user_registration` (`user_id` integer NOT NULL, `token` text NOT NULL, `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP), PRIMARY KEY (`user_id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "user_registration_token" to table: "user_registration"
CREATE UNIQUE INDEX `user_registration_token` ON `user_registration` (`token`);
-- Create "password_reset" table
CREATE TABLE `password_reset` (`user_id` integer NOT NULL, `token` text NOT NULL, `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP), PRIMARY KEY (`user_id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "password_reset_token" to table: "password_reset"
CREATE UNIQUE INDEX `password_reset_token` ON `password_reset` (`token`);
-- Create "recipe_image" table
CREATE TABLE `recipe_image` (`id` integer NOT NULL, `recipe_id` integer NOT NULL, `path` text NOT NULL, `sort_order` integer NOT NULL, `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP), PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`recipe_id`) REFERENCES `recipe` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "recipe_ingredient" table
CREATE TABLE `recipe_ingredient` (`id` integer NOT NULL, `step_id` integer NOT NULL, `ingredient_id` integer NOT NULL, `unit_id` integer NOT NULL, `amount` real NOT NULL, `sort_order` integer NOT NULL DEFAULT 0, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`unit_id`) REFERENCES `unit` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT, CONSTRAINT `1` FOREIGN KEY (`ingredient_id`) REFERENCES `ingredient` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `2` FOREIGN KEY (`step_id`) REFERENCES `recipe_step` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "recipe_ingredient_step_id_ingredient_id" to table: "recipe_ingredient"
CREATE UNIQUE INDEX `recipe_ingredient_step_id_ingredient_id` ON `recipe_ingredient` (`step_id`, `ingredient_id`);
-- Create index "recipe_ingredient_step_id_sort_order" to table: "recipe_ingredient"
CREATE UNIQUE INDEX `recipe_ingredient_step_id_sort_order` ON `recipe_ingredient` (`step_id`, `sort_order`);
-- Create index "idx_recipe_ingredients_sort_order" to table: "recipe_ingredient"
CREATE INDEX `idx_recipe_ingredients_sort_order` ON `recipe_ingredient` (`sort_order`);
-- Create "recipe_step" table
CREATE TABLE `recipe_step` (`id` integer NOT NULL, `recipe_id` integer NOT NULL, `instructions` text NOT NULL, `sort_order` integer NOT NULL DEFAULT 0, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`recipe_id`) REFERENCES `recipe` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "recipe_step_recipe_id_sort_order" to table: "recipe_step"
CREATE UNIQUE INDEX `recipe_step_recipe_id_sort_order` ON `recipe_step` (`recipe_id`, `sort_order`);
-- Create index "idx_recipe_steps_sort_order" to table: "recipe_step"
CREATE INDEX `idx_recipe_steps_sort_order` ON `recipe_step` (`sort_order`);
-- Create "recipe_tag" table
CREATE TABLE `recipe_tag` (`recipe_id` integer NOT NULL, `tag_id` integer NOT NULL, PRIMARY KEY (`recipe_id`, `tag_id`), CONSTRAINT `0` FOREIGN KEY (`tag_id`) REFERENCES `tag` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`recipe_id`) REFERENCES `recipe` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "meal_plan" table
CREATE TABLE `meal_plan` (`id` integer NOT NULL, `date` text NOT NULL DEFAULT (CURRENT_DATE), `user_id` integer NOT NULL, `recipe_id` integer NOT NULL, `sort_order` integer NOT NULL DEFAULT 0, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`recipe_id`) REFERENCES `recipe` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "meal_plan_date_sort_order" to table: "meal_plan"
CREATE UNIQUE INDEX `meal_plan_date_sort_order` ON `meal_plan` (`date`, `sort_order`);
-- Create index "idx_meal_plan_sort_order" to table: "meal_plan"
CREATE INDEX `idx_meal_plan_sort_order` ON `meal_plan` (`sort_order`);
-- Create "shopping_list" table
CREATE TABLE `shopping_list` (`id` integer NOT NULL, `user_id` integer NOT NULL, `name` text NOT NULL DEFAULT 'Shopping List', PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_shopping_lists_user_id" to table: "shopping_list"
CREATE INDEX `idx_shopping_lists_user_id` ON `shopping_list` (`user_id`);
-- Create index "idx_shopping_lists_name" to table: "shopping_list"
CREATE INDEX `idx_shopping_lists_name` ON `shopping_list` (`name`, `id`);
-- Create "shopping_list_item" table
CREATE TABLE `shopping_list_item` (`id` integer NOT NULL, `shopping_list_id` integer NOT NULL, `ingredient` text NOT NULL, `quantity` text NULL, `unit` text NULL, `done` boolean NOT NULL DEFAULT 0, `sort_order` integer NOT NULL DEFAULT 0, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`shopping_list_id`) REFERENCES `shopping_list` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "shopping_list_item_shopping_list_id_sort_order" to table: "shopping_list_item"
CREATE UNIQUE INDEX `shopping_list_item_shopping_list_id_sort_order` ON `shopping_list_item` (`shopping_list_id`, `sort_order`);
-- Create index "idx_shopping_list_items_shopping_list_id" to table: "shopping_list_item"
CREATE INDEX `idx_shopping_list_items_shopping_list_id` ON `shopping_list_item` (`shopping_list_id`);
-- Create index "idx_shopping_list_items_sort_order" to table: "shopping_list_item"
CREATE INDEX `idx_shopping_list_items_sort_order` ON `shopping_list_item` (`sort_order`);
