-- Create "tags" table
CREATE TABLE `tags` (`id` integer NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "tags_name" to table: "tags"
CREATE UNIQUE INDEX `tags_name` ON `tags` (`name`);
-- Create "recipe_tags" table
CREATE TABLE `recipe_tags` (`recipe_id` integer NOT NULL, `tag_id` integer NOT NULL, PRIMARY KEY (`recipe_id`, `tag_id`), CONSTRAINT `0` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
