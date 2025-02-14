CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL
);

CREATE TYPE pet_status AS ENUM ('available', 'pending', 'sold');

CREATE TABLE pets (
                      id SERIAL PRIMARY KEY,
                      name VARCHAR(255) NOT NULL,
                      category_id INT NOT NULL REFERENCES categories(id),
                      status pet_status NOT NULL,
                      deleted_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE TYPE order_status AS ENUM ('placed', 'approved', 'delivered');

CREATE TABLE orders(
    id SERIAL PRIMARY KEY,
    pet_id INT NOT NULL REFERENCES pets(id),
    quantity INT NOT NULL,
    ship_date TIMESTAMP WITHOUT TIME ZONE,
    status order_status NOT NULL,
    complete BOOLEAN NOT NULL
);

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email  VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    user_status INT NOT NULL
);