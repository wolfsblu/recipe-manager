-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_units" table
CREATE TABLE `new_units` (`id` integer NULL, `symbol` text NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Copy rows from old table "units" to new temporary table "new_units"
INSERT INTO `new_units` (`id`, `name`, `symbol`) SELECT `id`, `name`, `code` FROM `units`;
-- Drop "units" table after copying rows
DROP TABLE `units`;
-- Rename temporary table "new_units" to "units"
ALTER TABLE `new_units` RENAME TO `units`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
