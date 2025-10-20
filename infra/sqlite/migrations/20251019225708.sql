-- Create index "idx_recipes_created_at" to table: "recipes"
CREATE INDEX `idx_recipes_created_at` ON `recipes` (`created_at` DESC, `id` DESC);
-- Create index "idx_tags_name" to table: "tags"
CREATE INDEX `idx_tags_name` ON `tags` (`name`, `id`);
-- Create index "idx_units_name" to table: "units"
CREATE INDEX `idx_units_name` ON `units` (`name`, `id`);
-- Create index "idx_shopping_lists_name" to table: "shopping_lists"
CREATE INDEX `idx_shopping_lists_name` ON `shopping_lists` (`name`, `id`);
