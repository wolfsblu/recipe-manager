INSERT INTO roles (id, name)
VALUES (1, 'Administrator'),
       (2, 'Moderator'),
       (3, 'User');

INSERT INTO permissions (id, slug, name)
VALUES (1, 'can_create_recipe', 'Create Recipe'),
       (2, 'can_delete_recipe', 'Delete Recipe'),
       (3, 'can_edit_recipe', 'Edit Recipe'),
       (4, 'can_view_recipe', 'View Recipe');

INSERT INTO role_permissions (role_id, permission_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (1, 4),
       (2, 1),
       (2, 2),
       (2, 3),
       (2, 4),
       (3, 1),
       (3, 2),
       (3, 3),
       (3, 4);