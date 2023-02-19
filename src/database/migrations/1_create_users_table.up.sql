CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name varchar(255),
    last_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

CREATE OR REPLACE FUNCTION update_updated_at_column()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

insert into users (first_name, last_name, email) values ('Pavel', 'Bureš', 'pavel.bures@gmail.com');
insert into users (first_name, last_name, email) values ('Jana', 'Burešová', 'janik100@gmail.com');
insert into users (first_name, last_name, email) values ('Adam', 'Bureš', 'adam.bures.prg@gmail.com');
insert into users (first_name, last_name, email) values ('Karel', 'Bureš', 'karel.bures.prg@gmail.com');

CREATE TYPE mealtype AS ENUM ('breakfast', 'lunch', 'dinner');

CREATE TABLE meals (
    meal_id SERIAL PRIMARY KEY,
    author_id integer REFERENCES users ON DELETE CASCADE,
    meal_name varchar(255),
    meal_type mealtype,
    meal_date date,
    kcalories integer,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    UNIQUE(meal_date, meal_type)
);

CREATE INDEX meals_date_idx ON meals (meal_date);
CREATE TRIGGER update_users_modtime BEFORE UPDATE ON meals FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

CREATE TABLE consumers_meals (
    consumer_id integer REFERENCES users ON DELETE CASCADE,
    meal_id integer REFERENCES meals ON DELETE CASCADE
);