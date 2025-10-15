-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Drop "recipe_votes" table
DROP TABLE `recipe_votes`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
