INSERT INTO roles (id, name)
VALUES (1, 'Administrator'),
       (2, 'Moderator'),
       (3, 'User');

INSERT INTO permissions (id, slug, name)
VALUES
    -- Units
    (1, 'can_create_unit', 'Create Unit'),
    (2, 'can_delete_unit', 'Delete Unit'),
    (3, 'can_update_unit', 'Update Unit'),
    (4, 'can_view_unit', 'View Unit'),

    -- Ingredients
    (5, 'can_create_ingredient', 'Create Ingredient'),
    (6, 'can_delete_ingredient', 'Delete Ingredient'),
    (7, 'can_update_ingredient', 'Update Ingredient'),
    (8, 'can_view_ingredient', 'View Ingredient'),

    -- Meal Plans
    (9, 'can_create_meal_plan', 'Create Meal Plan'),
    (10, 'can_delete_meal_plan', 'Delete Meal Plan'),
    (11, 'can_update_meal_plan', 'Update Meal Plan'),
    (12, 'can_view_meal_plan', 'View Meal Plan'),

    -- Recipes
    (13, 'can_create_recipe', 'Create Recipe'),
    (14, 'can_delete_recipe', 'Delete Recipe'),
    (15, 'can_update_recipe', 'Update Recipe'),
    (16, 'can_view_recipe', 'View Recipe'),

    -- Profile
    (17, 'can_update_profile', 'Update Profile'),
    (18, 'can_view_profile', 'View Profile'),

    -- Shopping Lists
    (19, 'can_create_shopping_list', 'Create Shopping List'),
    (20, 'can_delete_shopping_list', 'Delete Shopping List'),
    (21, 'can_update_shopping_list', 'Update Shopping List'),
    (22, 'can_view_shopping_list', 'View Shopping List');

-- Give all permissions to all roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p;