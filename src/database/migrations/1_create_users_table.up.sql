CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name varchar(255),
    last_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL
);

insert into users (first_name, last_name, email) values ('Pavel', 'Bureš', 'pavel.bures@gmail.com');
insert into users (first_name, last_name, email) values ('Jana', 'Burešová', 'janik100@gmail.com');

/*
insert into users (first_name, last_name, email) values ('C', 'c', 'ccc');
*/

CREATE TYPE mealtype AS ENUM ('breakfast', 'lunch', 'dinner');

CREATE TABLE meals (
    meal_id SERIAL PRIMARY KEY,
    author_id integer REFERENCES users ON DELETE CASCADE,
    meal_name varchar(255),
    meal_type mealtype,
    meal_date date
);

CREATE INDEX meals_date_idx ON meals (meal_date);

/*
insert into meals (author_id, meal_name, meal_type, meal_date) values (1, 'AGulas',      'lunch', '2023-01-08');
insert into meals (author_id, meal_name, meal_type, meal_date) values (1, 'AChilli',     'lunch', '2023-01-08');
insert into meals (author_id, meal_name, meal_type, meal_date) values (1, 'AKoleno', 'breakfast', '2023-01-09');
insert into meals (author_id, meal_name, meal_type, meal_date) values (2, 'BFazolovka', 'dinner', '2023-01-09');
insert into meals (author_id, meal_name, meal_type, meal_date) values (2, 'BCibulacka', 'dinner', '2023-01-11');
insert into meals (author_id, meal_name, meal_type, meal_date) values (3, 'CRohlik', 'breakfast', '2023-01-12');
insert into meals (author_id, meal_name, meal_type, meal_date) values (3, 'CChleba', 'breakfast', '2023-01-12');
*/

CREATE TABLE consumers_meals (
    consumer_id integer REFERENCES users ON DELETE CASCADE,
    meal_id integer REFERENCES meals ON DELETE CASCADE
);

/*
insert into consumers_meals values (1, 1);
insert into consumers_meals values (1, 2);
insert into consumers_meals values (1, 3);
insert into consumers_meals values (2, 1);
insert into consumers_meals values (2, 2);
insert into consumers_meals values (2, 4);
*/
/*
Get meals and it's consumers:
select users.first_name as consumer, meal_name, meal_type from meals 
    join consumers_meals ON meals.meal_id=consumers_meals.meal_id
    join users ON user_id=consumers_meals.consumer_id 
    ORDER BY meal_id;
*/

/*
Get meals and it's authors;
select meal_name, meal_type, users.first_name as author from meals 
    join users ON meals.author_id=users.user_id 
    ORDER BY meal_name;
*/
