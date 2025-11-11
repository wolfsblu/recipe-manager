-- Drop index "idx_recipes_created_at" from table: "recipes"
DROP INDEX `idx_recipes_created_at`;
-- Create index "idx_recipes_created_at_desc" to table: "recipes"
CREATE INDEX `idx_recipes_created_at_desc` ON `recipes` (`created_at` DESC, `id` DESC);
-- Create index "idx_recipes_created_at_asc" to table: "recipes"
CREATE INDEX `idx_recipes_created_at_asc` ON `recipes` (`created_at`, `id`);
-- Create index "idx_recipes_name_asc" to table: "recipes"
CREATE INDEX `idx_recipes_name_asc` ON `recipes` (`name`, `id`);
-- Create index "idx_recipes_name_desc" to table: "recipes"
CREATE INDEX `idx_recipes_name_desc` ON `recipes` (`name` DESC, `id` DESC);
-- Create index "idx_recipes_servings_asc" to table: "recipes"
CREATE INDEX `idx_recipes_servings_asc` ON `recipes` (`servings`, `id`);
-- Create index "idx_recipes_servings_desc" to table: "recipes"
CREATE INDEX `idx_recipes_servings_desc` ON `recipes` (`servings` DESC, `id` DESC);
