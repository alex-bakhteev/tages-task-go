CREATE TABLE IF NOT EXISTS products
(
    id    SERIAL PRIMARY KEY,
    name  TEXT           NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    product_id  INT REFERENCES products (id),
    quantity    INT            NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL
);