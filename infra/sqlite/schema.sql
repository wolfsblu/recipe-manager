CREATE TABLE ingredient
(
    id   INTEGER NOT NULL PRIMARY KEY,
    name TEXT    NOT NULL
);

CREATE TABLE nutrient
(
    id   INTEGER NOT NULL PRIMARY KEY,
    name TEXT    NOT NULL UNIQUE,
    unit TEXT    NOT NULL
);

CREATE TABLE ingredient_nutrient
(
    ingredient_id INTEGER NOT NULL REFERENCES ingredient (id) ON DELETE CASCADE,
    nutrient_id   INTEGER NOT NULL REFERENCES nutrient (id) ON DELETE CASCADE,
    amount        REAL    NOT NULL,
    PRIMARY KEY (ingredient_id, nutrient_id)
);

CREATE TABLE permission
(
    id   INTEGER NOT NULL PRIMARY KEY,
    slug TEXT    NOT NULL UNIQUE,
    name TEXT    NOT NULL
);

CREATE TABLE recipe
(
    id          INTEGER   NOT NULL PRIMARY KEY,
    name        TEXT      NOT NULL,
    servings    INTEGER   NOT NULL,
    minutes     INTEGER   NOT NULL,
    description TEXT      NOT NULL,
    created_by  INTEGER   NOT NULL REFERENCES user (id) ON DELETE CASCADE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role
(
    id   INTEGER NOT NULL PRIMARY KEY,
    name TEXT    NOT NULL UNIQUE
);

CREATE TABLE role_permission
(
    role_id       INTEGER NOT NULL REFERENCES role (id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permission (id) ON DELETE CASCADE
);

CREATE TABLE tag
(
    id   INTEGER NOT NULL PRIMARY KEY,
    name TEXT    NOT NULL UNIQUE
);

CREATE TABLE unit
(
    id     INTEGER NOT NULL PRIMARY KEY,
    symbol TEXT,
    name   TEXT    NOT NULL
);

CREATE TABLE user
(
    id            INTEGER   NOT NULL PRIMARY KEY,
    email         TEXT      NOT NULL UNIQUE,
    password_hash TEXT      NOT NULL,
    is_confirmed  BOOLEAN   NOT NULL DEFAULT 0,
    role_id       INTEGER   NOT NULL REFERENCES role (id) ON DELETE CASCADE,
    locale        TEXT      NOT NULL DEFAULT 'en',
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_registration
(
    user_id    INTEGER   NOT NULL PRIMARY KEY REFERENCES user (id) ON DELETE CASCADE,
    token      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE password_reset
(
    user_id    INTEGER   NOT NULL PRIMARY KEY REFERENCES user (id) ON DELETE CASCADE,
    token      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipe_image
(
    id         INTEGER   NOT NULL PRIMARY KEY,
    recipe_id  INTEGER   NOT NULL REFERENCES recipe (id) ON DELETE CASCADE,
    path       TEXT      NOT NULL,
    sort_order INTEGER   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipe_ingredient
(
    id            INTEGER NOT NULL PRIMARY KEY,
    step_id       INTEGER NOT NULL REFERENCES recipe_step (id) ON DELETE CASCADE,
    ingredient_id INTEGER NOT NULL REFERENCES ingredient (id) ON DELETE CASCADE,
    unit_id       INTEGER NOT NULL REFERENCES unit (id) ON DELETE RESTRICT,
    amount        REAL    NOT NULL,
    sort_order    INTEGER NOT NULL DEFAULT 0,
    UNIQUE (step_id, ingredient_id),
    UNIQUE (step_id, sort_order)
);

CREATE TABLE recipe_step
(
    id           INTEGER NOT NULL PRIMARY KEY,
    recipe_id    INTEGER NOT NULL REFERENCES recipe (id) ON DELETE CASCADE,
    instructions TEXT    NOT NULL,
    sort_order   INTEGER NOT NULL DEFAULT 0,
    UNIQUE (recipe_id, sort_order)
);

CREATE TABLE recipe_tag
(
    recipe_id INTEGER NOT NULL REFERENCES recipe (id) ON DELETE CASCADE,
    tag_id    INTEGER NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, tag_id)
);

CREATE TABLE meal_plan
(
    id         INTEGER NOT NULL PRIMARY KEY,
    date       TEXT    NOT NULL DEFAULT CURRENT_DATE,
    user_id    INTEGER NOT NULL REFERENCES user (id) ON DELETE CASCADE,
    recipe_id  INTEGER NOT NULL REFERENCES recipe (id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    UNIQUE (date, sort_order)
);

CREATE TABLE shopping_list
(
    id      INTEGER NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES user (id) ON DELETE CASCADE,
    name    TEXT    NOT NULL DEFAULT 'Shopping List'
);

CREATE TABLE shopping_list_item
(
    id               INTEGER NOT NULL PRIMARY KEY,
    shopping_list_id INTEGER NOT NULL REFERENCES shopping_list (id) ON DELETE CASCADE,
    ingredient       TEXT    NOT NULL,
    quantity         TEXT,
    unit             TEXT,
    done             BOOLEAN NOT NULL DEFAULT 0,
    sort_order       INTEGER NOT NULL DEFAULT 0,
    UNIQUE (shopping_list_id, sort_order)
);

CREATE INDEX idx_meal_plan_sort_order ON meal_plan (sort_order);
CREATE INDEX idx_recipe_ingredients_sort_order ON recipe_ingredient (sort_order);
CREATE INDEX idx_recipe_steps_sort_order ON recipe_step (sort_order);
CREATE INDEX idx_shopping_lists_user_id ON shopping_list (user_id);
CREATE INDEX idx_shopping_list_items_shopping_list_id ON shopping_list_item (shopping_list_id);
CREATE INDEX idx_shopping_list_items_sort_order ON shopping_list_item (sort_order);
CREATE INDEX idx_recipes_created_at_desc ON recipe (created_at DESC, id DESC);
CREATE INDEX idx_recipes_created_at_asc ON recipe (created_at ASC, id ASC);
CREATE INDEX idx_recipes_name_asc ON recipe (name ASC, id ASC);
CREATE INDEX idx_recipes_name_desc ON recipe (name DESC, id DESC);
CREATE INDEX idx_recipes_servings_asc ON recipe (servings ASC, id ASC);
CREATE INDEX idx_recipes_servings_desc ON recipe (servings DESC, id DESC);
CREATE INDEX idx_units_name ON unit (name, id);
CREATE INDEX idx_tags_name ON tag (name, id);
CREATE INDEX idx_shopping_lists_name ON shopping_list (name, id);
