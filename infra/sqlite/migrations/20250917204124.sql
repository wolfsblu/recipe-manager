-- Create "recipe_votes" table
CREATE TABLE `recipe_votes` (`id` integer NULL, `recipe_id` integer NOT NULL, `user_id` integer NOT NULL, `vote` integer NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CHECK (vote IN (-1, 1)));
-- Create index "recipe_votes_recipe_id_user_id" to table: "recipe_votes"
CREATE UNIQUE INDEX `recipe_votes_recipe_id_user_id` ON `recipe_votes` (`recipe_id`, `user_id`);
-- Create index "idx_recipe_votes_recipe_id" to table: "recipe_votes"
CREATE INDEX `idx_recipe_votes_recipe_id` ON `recipe_votes` (`recipe_id`);
-- Create index "idx_recipe_votes_user_id" to table: "recipe_votes"
CREATE INDEX `idx_recipe_votes_user_id` ON `recipe_votes` (`user_id`);
