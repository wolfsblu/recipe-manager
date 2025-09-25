-- Create "shopping_lists" table
CREATE TABLE `shopping_lists` (`id` integer NULL, `user_id` integer NOT NULL, `name` text NOT NULL DEFAULT 'Shopping List', PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_shopping_lists_user_id" to table: "shopping_lists"
CREATE INDEX `idx_shopping_lists_user_id` ON `shopping_lists` (`user_id`);
-- Create "shopping_list_items" table
CREATE TABLE `shopping_list_items` (`id` integer NULL, `shopping_list_id` integer NOT NULL, `ingredient` text NOT NULL, `quantity` text NULL, `unit` text NULL, `done` boolean NOT NULL DEFAULT 0, `sort_order` integer NOT NULL DEFAULT 0, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`shopping_list_id`) REFERENCES `shopping_lists` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "shopping_list_items_shopping_list_id_sort_order" to table: "shopping_list_items"
CREATE UNIQUE INDEX `shopping_list_items_shopping_list_id_sort_order` ON `shopping_list_items` (`shopping_list_id`, `sort_order`);
-- Create index "idx_shopping_list_items_shopping_list_id" to table: "shopping_list_items"
CREATE INDEX `idx_shopping_list_items_shopping_list_id` ON `shopping_list_items` (`shopping_list_id`);
-- Create index "idx_shopping_list_items_sort_order" to table: "shopping_list_items"
CREATE INDEX `idx_shopping_list_items_sort_order` ON `shopping_list_items` (`sort_order`);
