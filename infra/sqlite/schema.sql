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

CREATE TABLE tags
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE units
(
    id   INTEGER PRIMARY KEY,
    code TEXT,
    name TEXT NOT NULL
);

CREATE TABLE ingredients
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE users
(
    id            INTEGER PRIMARY KEY,
    email         TEXT      NOT NULL UNIQUE,
    password_hash TEXT      NOT NULL,
    is_confirmed  BOOLEAN   NOT NULL DEFAULT 0,
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
    step_id       INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
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

CREATE INDEX idx_meal_plan_sort_order ON meal_plan (sort_order);
CREATE INDEX idx_recipe_ingredients_sort_order ON recipe_ingredients (sort_order);
CREATE INDEX idx_recipe_steps_sort_order ON recipe_steps (sort_order);
