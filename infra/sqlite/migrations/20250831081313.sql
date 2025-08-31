-- Create "permissions" table
CREATE TABLE `permissions` (`id` integer NULL, `slug` text NOT NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "permissions_slug" to table: "permissions"
CREATE UNIQUE INDEX `permissions_slug` ON `permissions` (`slug`);
-- Create "roles" table
CREATE TABLE `roles` (`id` integer NULL, `name` text NOT NULL, PRIMARY KEY (`id`));
-- Create index "roles_name" to table: "roles"
CREATE UNIQUE INDEX `roles_name` ON `roles` (`name`);
-- Create "role_permissions" table
CREATE TABLE `role_permissions` (`role_id` integer NOT NULL, `permission_id` integer NOT NULL, CONSTRAINT `0` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "user_roles" table
CREATE TABLE `user_roles` (`user_id` integer NOT NULL, `role_id` integer NOT NULL, CONSTRAINT `0` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
