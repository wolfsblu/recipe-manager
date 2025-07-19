-- Create "recipe_images" table
CREATE TABLE `recipe_images` (`id` integer NULL, `recipe_id` integer NOT NULL, `path` text NOT NULL, `sort_order` integer NOT NULL, `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP), PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
