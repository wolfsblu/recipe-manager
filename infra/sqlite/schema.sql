CREATE TABLE ingredients
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE nutrients
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    unit TEXT NOT NULL
);

CREATE TABLE ingredient_nutrients
(
    ingredient_id INTEGER NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    nutrient_id   INTEGER NOT NULL REFERENCES nutrients (id) ON DELETE CASCADE,
    amount        REAL    NOT NULL,
    PRIMARY KEY (ingredient_id, nutrient_id)
);

CREATE TABLE permissions
(
    id   INTEGER PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL
);

CREATE TABLE recipes
(
    id          INTEGER PRIMARY KEY,
    name        TEXT      NOT NULL,
    servings    INTEGER   NOT NULL,
    minutes     INTEGER   NOT NULL,
    description TEXT      NOT NULL,
    created_by  INTEGER   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE role_permissions
(
    role_id       INTEGER NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions (id) ON DELETE CASCADE
);

CREATE TABLE tags
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE units
(
    id     INTEGER PRIMARY KEY,
    symbol TEXT,
    name   TEXT NOT NULL
);

CREATE TABLE users
(
    id            INTEGER PRIMARY KEY,
    email         TEXT      NOT NULL UNIQUE,
    password_hash TEXT      NOT NULL,
    is_confirmed  BOOLEAN   NOT NULL DEFAULT 0,
    role_id       INTEGER   NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    locale        TEXT      NOT NULL DEFAULT 'en',
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_registrations
(
    user_id    INTEGER PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    token      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE password_resets
(
    user_id    INTEGER PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    token      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipe_images
(
    id         INTEGER PRIMARY KEY,
    recipe_id  INTEGER   NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    path       TEXT      NOT NULL,
    sort_order INTEGER   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipe_ingredients
(
    id            INTEGER PRIMARY KEY,
    step_id       INTEGER NOT NULL REFERENCES recipe_steps (id) ON DELETE CASCADE,
    ingredient_id INTEGER NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    unit_id       INTEGER NOT NULL REFERENCES units (id) ON DELETE RESTRICT,
    amount        REAL    NOT NULL,
    sort_order    INTEGER NOT NULL DEFAULT 0,
    UNIQUE (step_id, ingredient_id),
    UNIQUE (step_id, sort_order)
);

CREATE TABLE recipe_steps
(
    id           INTEGER PRIMARY KEY,
    recipe_id    INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    instructions TEXT    NOT NULL,
    sort_order   INTEGER NOT NULL DEFAULT 0,
    UNIQUE (recipe_id, sort_order)
);

CREATE TABLE recipe_tags
(
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    tag_id    INTEGER NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, tag_id)
);

CREATE TABLE meal_plan
(
    id         INTEGER PRIMARY KEY,
    date       TEXT    NOT NULL DEFAULT CURRENT_DATE,
    user_id    INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    recipe_id  INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    UNIQUE (date, sort_order)
);

CREATE TABLE recipe_votes
(
    id        INTEGER PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    user_id   INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    vote      INTEGER NOT NULL CHECK (vote IN (-1, 1)),
    UNIQUE (recipe_id, user_id)
);

CREATE TABLE shopping_lists
(
    id      INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name    TEXT    NOT NULL DEFAULT 'Shopping List'
);

CREATE TABLE shopping_list_items
(
    id               INTEGER PRIMARY KEY,
    shopping_list_id INTEGER NOT NULL REFERENCES shopping_lists (id) ON DELETE CASCADE,
    ingredient       TEXT    NOT NULL,
    quantity         TEXT,
    unit             TEXT,
    done             BOOLEAN NOT NULL DEFAULT 0,
    sort_order       INTEGER NOT NULL DEFAULT 0,
    UNIQUE (shopping_list_id, sort_order)
);

CREATE INDEX idx_meal_plan_sort_order ON meal_plan (sort_order);
CREATE INDEX idx_recipe_ingredients_sort_order ON recipe_ingredients (sort_order);
CREATE INDEX idx_recipe_steps_sort_order ON recipe_steps (sort_order);
CREATE INDEX idx_recipe_votes_recipe_id ON recipe_votes (recipe_id);
CREATE INDEX idx_recipe_votes_user_id ON recipe_votes (user_id);
CREATE INDEX idx_shopping_lists_user_id ON shopping_lists (user_id);
CREATE INDEX idx_shopping_list_items_shopping_list_id ON shopping_list_items (shopping_list_id);
CREATE INDEX idx_shopping_list_items_sort_order ON shopping_list_items (sort_order);
