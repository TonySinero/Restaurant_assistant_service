CREATE
EXTENSION cube;
CREATE
EXTENSION earthdistance;

CREATE TABLE order_statuses
(
    id          INT PRIMARY KEY,
    description VARCHAR
);

CREATE TABLE types
(
    id         UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    name       varchar   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE categories
(
    id          INT PRIMARY KEY,
    description VARCHAR
);

CREATE TABLE delivery_types
(
    id          INT PRIMARY KEY,
    description VARCHAR
);

CREATE TABLE payment_types
(
    id          INT PRIMARY KEY,
    description VARCHAR
);

CREATE TABLE restaurants
(
    id              UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    title           varchar   NOT NULL,
    rating          DECIMAL,
    medium_price    DECIMAL,
    user_id         INT       NOT NULL,
    address         varchar,
    image           varchar,
    time_work_start time without time zone,
    time_work_end   time without time zone,
    is_active       bool,
    number          numeric,
    email           varchar,
    description     varchar,
    created_at      TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    latitude        DECIMAL,
    longitude       DECIMAL
);

CREATE TABLE location
(
    restaurant_id UUID references restaurants (id) on delete cascade NOT NULL,
    latitude      varchar                                            NOT NULL,
    longitude     varchar                                            NOT NULL
);

CREATE TABLE dishes
(
    id            UUID PRIMARY KEY                                   NOT NULL DEFAULT gen_random_uuid(),
    restaurant_id UUID references restaurants (id) on delete cascade NOT NULL,
    type          UUID references types (id)                         NOT NULL,
    cost          DECIMAL                                            NOT NULL,
    name          varchar                                            NOT NULL,
    cooking_time  integer,
    image         varchar,
    weight        DECIMAL,
    description   varchar,
    status        varchar                                            NOT NULL DEFAULT 'available',
    created_at    TIMESTAMP                                          NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
        CONSTRAINT check_status CHECK ( status IN ('available', 'unavailable') )
);

CREATE TABLE orders
(
    id                  UUID PRIMARY KEY                                   NOT NULL,
    row_id              SERIAL                                             NOT NULL unique,
    restaurant_id       UUID references restaurants (id) on delete cascade NOT NULL,
    cost                DECIMAL,
    delivery_time       timestamp                                          NOT NULL,
    client_full_name    VARCHAR                                            NOT NULL,
    client_phone_number TEXT                                               NOT NULL,
    address             VARCHAR                                            NOT NULL,
    delivery_type       INT references delivery_types (id),
    payment_type        INT                                                NOT NULL DEFAULT 1,
    courier_service     INT,
    status              INT REFERENCES order_statuses (id)                 NOT NULL DEFAULT 1,
    view_status         BOOLEAN                                                     DEFAULT FALSE,
    created_at          TIMESTAMP                                          NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE feedbacks
(
    id            UUID PRIMARY KEY                                   NOT NULL DEFAULT gen_random_uuid(),
    restaurant_id UUID references restaurants (id) on delete cascade NOT NULL,
    order_id      UUID references orders (id) on delete cascade      NOT NULL,
    feedback      VARCHAR,
    rating        INTEGER
);

CREATE TABLE order_dishes
(
    id         UUID PRIMARY KEY                                       DEFAULT gen_random_uuid(),
    order_id   UUID REFERENCES orders (id) on delete cascade NOT NULL,
    dish_id    UUID REFERENCES dishes (id)                   NOT NULL,
    amount     integer                                       NOT NULL,
    rate       DECIMAL,
    created_at TIMESTAMP                                     NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE category_restaurants
(
    id            UUID PRIMARY KEY                                              DEFAULT gen_random_uuid(),
    category_id   integer REFERENCES categories (id) on delete cascade NOT NULL,
    restaurant_id UUID REFERENCES restaurants (id)                     NOT NULL,
    created_at    TIMESTAMP                                            NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

INSERT INTO payment_types(id, description)
VALUES (1, 'card'),
       (2, 'cash'),
       (3, 'card online');

INSERT INTO order_statuses(id, description)
VALUES (1, 'New'),
       (2, 'In progress'),
       (3, 'Ready for delivery'),
       (4, 'Completed'),
       (5, 'Canceled');

INSERT INTO delivery_types(id, description)
VALUES (1, 'Restaurant delivery'),
       (2, 'Service delivery');

INSERT INTO categories(id, description)
VALUES (1, 'Asian'),
       (2, 'Burgers'),
       (3, 'Breakfast'),
       (4, 'Pizza'),
       (5, 'Lunch'),
       (6, 'Fast-food'),
       (7, 'Dessert'),
       (8, 'Coffee'),
       (9, 'Steaks');

INSERT INTO restaurants (id, user_id, title, rating, medium_price, address, time_work_start,
                         time_work_end, is_active,
                         image, latitude, longitude, description)
VALUES ('01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 1, 'Meat & Fish', 8.5, 51.20, 'Kulman street, 14k1, Minsk',
        '08:00:00',
        '15:00:00', TRUE,
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b/37de00e0-ceab-4dd9-843b-ce68e01e3159.jpg',
        53.921435397746784, 27.581316198274365,
        'Chain steakhouse and fish restaurant. The concept of the establishment is a harmonious combination of fish and meat dishes. The menu focuses on domestic grain-fed beef and seafood from Kamchatka, Yakutia and Murmansk.'),
       ('02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 28, 'KFC', 7.1, 15.00, 'Pritytskogo street, 101a, Minsk', '09:00:00',
        '17:00:00', TRUE,
        'https://onlineshop.fra1.digitaloceanspaces.com/02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b/1ce759d9-2467-44e4-b9b2-591388e1cb80.png',
        53.906670487925986, 27.435194827109605,
        'Kentucky Fried Chicken, abbreviated as KFC, is an international chain of catering restaurants specializing in chicken dishes.'),
       ('03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 1, 'Burger KING', 6.5, 12.00, 'Pritytskogo street, 19a, Minsk',
        '10:00:00',
        '18:00:00', TRUE,
        'https://onlineshop.fra1.digitaloceanspaces.com/03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b/7ade7bc2-1dfa-4fe1-86ab-272b2c4fe221.png',
        53.909390439793654, 27.496377398273868,
        'Burger King Corporation is an American company that owns the Burger King global fast food chain specializing in hamburgers.'),
       ('04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b', 1, 'McDonald’s', 9, 11.20, 'Pritytskogo street 28, Minsk', '08:00:00',
        '23:00:00', TRUE,
        'https://onlineshop.fra1.digitaloceanspaces.com/04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b/02013725-971f-4a58-b375-ebc0c66b7765.png',
        53.909390439793654, 27.496377398273868,
        'McDonald’s is an American corporation operating in the field of catering, the world''s largest chain of fast food restaurants.'),
       ('04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b', 1, 'Dodo Pizza', 5, 25.00, 'st. Romanovskaya Sloboda 8, Minsk',
        '00:08:00',
        '00:00:00',
        TRUE,
        'https://onlineshop.fra1.digitaloceanspaces.com/1738cc31-6075-4588-8efe-9a1c1a0f0ae4/0ab0fed5-75af-4ed4-94c4-a61ed2e100d5.jpg',
        53.902777, 27.547575,
        'Huge selection on the menu. The perfect combination of flavors and the thinnest dough is the key to delicious pizza. Delicious pizza delivery with the best ingredients.'),
       ('04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b', 1, 'Domino’s Pizza', 5, 26.00, 'Ave. Pobediteley 1, Minsk', '00:08:00',
        '00:00:00',
        TRUE,
        'https://onlineshop.fra1.digitaloceanspaces.com/1738cc31-6075-4588-8efe-9a1c1a0f0ae4/acf12143-ba1a-4f36-b712-6f31e460c8e3.png',
        53.905143, 27.552264,
        'Dominos Pizza is an American catering company. Operates the world''s largest chain of pizzerias.'),
       ('05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b', 1, 'NEFT', 9.5, 60.00, 'Aranskaya street 8, Minsk', '00:00:00',
        '00:00:00',
        TRUE,
        'https://onlineshop.fra1.digitaloceanspaces.com/05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b/e7d9e569-5d51-4a7e-b39e-b79a0c89b6c3.jpg',
        53.885351410688884, 27.569866540601108,
        'The NEFT Lounge Restaurant is a unique establishment with a conceptual name, the main goal of which is to create a truly high-quality and atmospheric vacation. NEFT combines all the best from the restaurant and hookah industry: a high-quality menu thought out to the smallest detail, a high level of service and serving dishes, an individual approach to each guest, a perfectly matched hookah park and a huge selection of premium positions of tobacco-free nicotine-free blends, as well as a wide tea room map, which has more than fifty types of tea.'),
       ('06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b', 1, 'KOKOS', 8, 45.00, 'Bogdan Khmelnitsky street 4, Minsk',
        '10:00:00',
        '18:00:00', TRUE,
        'https://onlineshop.fra1.digitaloceanspaces.com/06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b/d42a1d19-bc13-41c9-9aed-5b9d6751aba2.jpg',
        53.922992038531554, 27.596746998274412, 'The atmosphere of real Asian cuisine is felt in the institution.');

INSERT INTO types (id, name)
VALUES ('11fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 'Steaks'),
       ('11fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 'Lunch'),
       ('11fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 'Salads'),
       ('11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 'Snacks'),

       ('22fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 'Burgers'),
       ('22fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 'Baskets'),
       ('22fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 'Chicken'),
       ('23fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 'Sauce'),

       ('33fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 'Beef burgers'),
       ('33fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 'Chicken burgers'),
       ('33fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 'Fish burgers'),

       ('33fb44e3-5f18-04eb-80a1-d8b4e8a33f1b', 'Pizzas'),
       ('33fb44e3-5f18-04eb-80a1-d8b4e8a23f1b', 'Desserts'),
       ('33fb44e3-5f18-04eb-80a1-d8b4e8a11f1b', 'Drinks'),

       ('77fb44e3-3f18-03eb-82a1-d8b4e8a22f1b', 'Paste'),
       ('77fb44e3-3f18-04eb-82a1-d8b4e8a22f1b', 'Meat'),

       ('99fb44e3-1f18-01eb-20a1-d8b4e8a22f1b', 'Curry'),
       ('99fb44e3-1f18-02eb-20a1-d8b4e8a22f1b', 'Soups'),
       ('99fb44e3-1f18-03eb-20a1-d8b4e8a22f1b', 'Wok');


INSERT INTO dishes (id, restaurant_id, type, cost, name, image, description, weight)
VALUES ('01aa44e3-5f18-45eb-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 35, 'Filet mignon steak',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/c45b77a2-4865-476c-8be8-da23d432d0f4.jpeg',
        'Filet mignon steak  signature spice rubbed, grass-finished steak, grilled medium rare', 300),
       ('02aa44e3-0000-45eb-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 25.90, 'Pork steaks',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/c9b2d180-99ab-4fe6-ba9f-3152eb458ab2.jpg',
        'Pork steaks -  grass-finished steak grilled medium rare, drizzled with our fresh chimichurri sauce', 200),
       ('03aa44e3-0000-45eb-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 6.80, 'Mini Caesar salad with chicken',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/e4ca6f37-23dd-49dc-9301-773a11b8d8bb.jpg',
        'Mini Caesar salad with chicken romaine lettuce, baby kale, never frozen chicken, tomatoes, garlic croutons, shaved asiago, caesar dressing',
        160),
       ('04aa44e3-0000-45eb-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 6.80, 'Mini Caesar salad with salmon',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/06610fa1-7d57-4028-ad9d-e7695c3a1590.jpg',
        'Mini Caesar salad with salmon romaine lettuce, organic baby kale, grilled salmon, tomatoes, garlic croutons, shaved asiago, caesar dressing ',
        150),
       ('05aa44e3-0000-45eb-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 16.90, 'Meat salad',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/75134d88-2e54-4ad4-bb03-1bfbf8e2a006.jpg',
        'Meet calad up lettuce blend, steak, pineapple-mango salsa, jicama, mint, coconut roasted cashews, mild jalapeno lime dressing, carrot',
        200),
       ('06aa44e3-0000-45eb-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 10.30, 'Greek salad',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/2aee1d8f-832c-4ec9-ae63-aa8ef907ef83.jpg',
        'Greek salad marinated yellow & red beets, cucumber,  arugula, artisan lettuce, roasted tomatoes, carrot, crisp jicama, raw walnuts, hemp seeds, goat cheese ',
        200),
       ('07aa44e3-0000-45eb-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 11, 'Vegetable snack',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/85d69d91-684e-4544-911b-6d4874c112b2.jpg',
        'Vegetable snack cucumbers, celery, mint, pickled onions, raw walnuts, grapes, blue cheese, apple, carrot, pepper ',
        200),
       ('08aa44e3-0000-45eb-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 19, 'Carpaccio',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/3baa1429-f9e9-4a76-83c9-612b921c0aa7.jpg',
        'Carpaccio fresh beef tenderloin, arugula, watercress, shaved Parmigiano-Reggiano cheese,  freshly ground pepper ',
        110),

       ('01bb44e3-0000-45eb-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '22fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 4, 'Chef Burger junior',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/b9bd80ff-3f12-4a10-9950-436457ba8587.png',
        'Chef Burger junior - 2 strips, iceberg lettuce, tomatoes, wheat bun with black and white sesame seeds', 140),
       ('02bb44e3-0000-45eb-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '22fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 4, 'Chef Burger junior spicy',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/3efc0757-04c1-4c8e-8be7-6b1691146a19.jpg',
        'Chef Burger junior spicy - Hot & spicy spicy chicken, lettuce, pickled cucumbers, onions, signature Burger sauce, sesame bun. ',
        150),
       ('03bb44e3-0000-45eb-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '22fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 14.10, 'Basket S',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/694af036-329b-4482-b45f-de9d646b646d.png',
        'Basket S - 12 spicy breaded chicken wings', 190),
       ('04bb44e3-0000-45eb-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '22fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 19.80, 'Basket M',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/bdcc6a1e-2062-4021-a095-3e376de15b57.png',
        'Basket M – 20 spicy breaded chicken wings', 260),
       ('05bb44e3-0000-45eb-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 4.69, 'Four legs',
        'https://onlineshop.fra1.digitaloceanspaces.com/05bb44e3-0000-45eb-80a1-d8b4e8a22f1b/89e7ca43-fdf4-4904-8605-020e45d47a1e.png',
        'Four legs - bread crumbs, freshly ground black pepper, 4  chicken legs, split ', 100),
       ('06bb44e3-0000-45eb-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 5, 'Three strips',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/927adc3d-9c78-4107-9dae-15e6342f7100.png',
        'Three strips - bread crumbs, freshly ground black pepper, 3  chicken fillet ', 110),
       ('07bb44e3-0000-45eb-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '23fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 0.9, 'Cheese sauce',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/e3ed28ac-4322-4cf3-8444-c58498455c50.jpg',
        '', 30),
       ('08bb44e3-0000-45eb-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '23fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 0.9, 'Barbecue sauce',
        'https://onlineshop.fra1.digitaloceanspaces.com/08bb44e3-0000-45eb-80a1-d8b4e8a22f1b/688f8fc6-7979-4a77-917c-86edddbfbfc5.png',
        '', 30),

       ('01ee44e3-0000-45eb-80a1-d8b4e8a22f1b', '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '33fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 6.30, 'Big King',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/d037b4e8-b3cb-46db-bc97-3ef7f459b17a.jpg',
        'Big King - beef with tomatoes, fresh chopped lettuce, mayonnaise, pickled cucumbers and chopped white onions on a bun with sesame topping',
        100),
       ('02ee44e3-0000-45eb-80a1-d8b4e8a22f1b', '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '33fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 7, 'Whopper',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/6738b828-aa21-42e7-8968-4485cc1f127c.jpg',
        'Whopper - Double steak with Salsa sauce, tomatoes, pickled cucumbers, Iceberg lettuce on cheese roll, Mozzarella cheese ',
        110),
       ('03ee44e3-0000-45eb-80a1-d8b4e8a22f1b', '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '33fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 5.90, 'Chicken King',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/2ddd8e8c-70ac-4604-a93e-2eff32522238.jpg',
        'Chicken King - breaded chicken cutlets, cucumbers, fresh onion and sesame-sprinkled bun ', 130),
       ('04ee44e3-0000-45eb-80a1-d8b4e8a22f1b', '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '33fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 9.50, 'Mozzarella Chicken',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/3e265f1b-8aa7-470c-b70c-e77f9f17e597.jpg',
        'Mozzarella Chicken - beef, tomatoes, fresh chopped lettuce, mayonnaise, cucumbers and chopped white onions, sesame seed bun and cheese',
        150),
       ('05ee44e3-0000-45eb-80a1-d8b4e8a22f1b', '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '33fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 7.90, 'Shrimp King',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/16c82df5-5876-42b3-8f02-32b2a618c56e.jpg',
        'Shrimp King - breaded king prawns, fresh Iceberg lettuce, cheese slices, pickled cucumbers, Caesar sauce on a potato bun',
        110),
       ('07ee44e3-0000-45eb-80a1-d8b4e8a22f1b', '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 12.50, 'King Bouquet',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/a3c6f71d-6a3a-4391-922e-70e9106ebadd.jpg',
        'King Bouquet - breaded king prawns, potato slices ', 350),
       ('08ee44e3-0000-45eb-80a1-d8b4e8a22f1b', '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 5.30, 'King Nuggets M',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/421648f7-034a-442e-8ca2-5e6068a05219.jpg',
        'King Nuggets M – garlic black pepper, boneless skinless chicken breasts,  breadcrumbs ', 120),

       ('01ff44e3-0000-45eb-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b',
        '33fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 10.10, 'Big Tasty',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/497f6c1b-fbd8-4877-a6af-99caf81c9e82.jpg',
        'Big Tasty - beef, tomatoes, fresh chopped lettuce, mayonnaise, cucumbers and chopped white onions, sesame seed bun and cheese',
        100),
       ('02ff44e3-0000-45eb-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b',
        '33fb44e3-5f18-01eb-80a1-d8b4e8a22f1b', 6.50, 'Big Mac',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/5451021c-a9d3-45ea-b5f7-db0133ee7558.jpg',
        'Big Mac - steak with Salsa sauce, tomatoes, pickled cucumbers, Iceberg lettuce on cheese roll, Mozzarella cheese ',
        110),
       ('03ff44e3-0000-45eb-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b',
        '33fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 2.70, 'Chicken burger',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/5940dd0d-0c26-4249-bdf3-492ba595dab8.png',
        'Chicken burger - chicken cutlets, fresh Iceberg lettuce, tomato,  cheese, onion and sesame-sprinkled bun ',
        130),
       ('04ff44e3-0000-45eb-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b',
        '33fb44e3-5f18-02eb-80a1-d8b4e8a22f1b', 6, 'McChicken',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/5c687172-3f60-4542-8484-5f33e33e6d08.jpg',
        'McChicken - breaded chicken cutlets, cucumbers, fresh onion, sesame seed bun and cheese ', 130),
       ('05ff44e3-0000-45eb-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b',
        '33fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 6.50, 'Filet o fish',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/b0d7d522-b4a7-411b-9df2-a9c45dc9654b.jpg',
        'Fillet o Fish - fish cakes, cucumbers, fresh onion and sesame seed bun ', 130),
       ('06ff44e3-0000-45eb-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b',
        '33fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 7.60, 'Double Fillet o Fish',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/3c099b3c-89f2-4122-acfc-9349e41d3c86.jpg',
        'Double Fillet o Fish - 2 fish cakes, cucumbers, fresh onion and sesame seed bun ', 150),
       ('07ff44e3-0000-45eb-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 4.10, 'French fries large',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/4f967c62-4b4e-4be7-b602-6a6b6466d8db.jpg',
        'French fries - large potato slices ', 150),
       ('08ff44e3-0000-45eb-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 5.80, 'Chicken McNuggets 6 pcs',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/5c770f70-908a-49a3-b368-724f116cf32d.jpg',
        'Chicken McNuggets – garlic black pepper, boneless skinless chicken breasts,  breadcrumbs ', 120),

       ('01aa44e3-0000-45aa-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a33f1b', 26, 'Pepperoni Pizza M',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/61321b0f-f5a4-42c6-b950-23dd73069f4c.jpeg',
        'Pepperoni Pizza M -  barbecue sauce, pepperoni, brisket (pork), mozzarella cheese, basil ', 310),
       ('02aa44e3-0000-45aa-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a33f1b', 28, 'Barbecue Pizza M',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/8d51fecd-f299-4d10-87a5-62be90e854c7.png',
        'Barbecue Pizza M - pizza sauce, brisket (pork), chicken fillet, fresh onion, barbecue sauce, mozzarella cheese, basil ',
        150),
       ('03aa44e3-0000-45aa-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a23f1b', 8, 'Syrniki',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/92b442a7-901f-4435-97f6-2e507ed2c49c.jpg',
        'Syrniki - cottage cheese, eggs, flour, sour cream ', 160),
       ('04aa44e3-0000-45aa-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a23f1b', 5.39, 'Cupcake',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/f8121c1a-295e-413b-850a-a0e7f82158eb.jpg',
        'Cupcake - dough, sour  cream, eggs, jam ', 80),
       ('05aa44e3-0000-45aa-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a11f1b', 3, 'Coca Cola 0.5',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/970c2158-beea-462e-94c5-66c03c65904e.jpg',
        '', 500),
       ('06aa44e3-0000-45aa-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a11f1b', 3, 'Pepsi 0.5',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/86d89786-b3fa-4daf-907b-07c60cb0abe1.png',
        '', 500),
       ('07aa44e3-0000-45aa-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b',
        '11fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 6.8, 'Mini Caesar salad with chicken',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/e4ca6f37-23dd-49dc-9301-773a11b8d8bb.jpg',
        'Mini Caesar salad with chicken - romaine lettuce, baby kale, never frozen chicken, tomatoes, garlic croutons, shaved asiago, caesar dressing ',
        160),
       ('08aa44e3-0000-45aa-80a1-d8b4e8a22f1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b',
        '11fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 6.8, 'Mini Caesar salad with salmon',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/06610fa1-7d57-4028-ad9d-e7695c3a1590.jpg',
        'Mini Caesar salad with salmon  - romaine lettuce, grilled salmon, tomatoes, garlic croutons, shaved asiago, caesar dressing ',
        150),





       ('01aa44e3-0000-45aa-80a2-d8b4e8a22a1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a33f1b', 26, 'Pepperoni Pizza M',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/61321b0f-f5a4-42c6-b950-23dd73069f4c.jpeg',
        'Pepperoni Pizza M -  barbecue sauce, pepperoni, brisket (pork), mozzarella cheese, basil ', 310),
       ('02aa44e3-0000-45aa-80a2-d8b4e8a22a1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a33f1b', 28, 'Barbecue Pizza M',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/8d51fecd-f299-4d10-87a5-62be90e854c7.png',
        'Barbecue Pizza M - pizza sauce, brisket (pork), chicken fillet, fresh onion, barbecue sauce, mozzarella cheese, basil ',
        150),
       ('03aa44e3-0000-45aa-80a2-d8b4e8a22a1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a23f1b', 8, 'Syrniki',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/92b442a7-901f-4435-97f6-2e507ed2c49c.jpg',
        'Syrniki - cottage cheese, eggs, flour, sour cream ', 160),
       ('04aa44e3-0000-45aa-80a2-d8b4e8a22a1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a23f1b', 5.39, 'Cupcake',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/f8121c1a-295e-413b-850a-a0e7f82158eb.jpg',
        'Cupcake - dough, sour  cream, eggs, jam ', 80),
       ('05aa44e3-0000-45aa-80a2-d8b4e8a22a1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a11f1b', 3, 'Coca Cola 0.5',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/970c2158-beea-462e-94c5-66c03c65904e.jpg',
        '', 300),
       ('06aa44e3-0000-45aa-80a2-d8b4e8a22a1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b',
        '33fb44e3-5f18-04eb-80a1-d8b4e8a11f1b', 3, 'Pepsi 0.5',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/86d89786-b3fa-4daf-907b-07c60cb0abe1.png',
        '', 300),
       ('07aa44e3-0000-45aa-80a2-d8b4e8a22a1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b',
        '11fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 6.8, 'Mini Caesar salad with chicken',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/e4ca6f37-23dd-49dc-9301-773a11b8d8bb.jpg',
        'Mini Caesar salad with chicken - romaine lettuce, baby kale, never frozen chicken, tomatoes, garlic croutons, shaved asiago, caesar dressing ',
        160),
       ('08aa44e3-0000-45aa-80a2-d8b4e8a22a1b', '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b',
        '11fb44e3-5f18-03eb-80a1-d8b4e8a22f1b', 6.8, 'Mini Caesar salad with salmon',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/06610fa1-7d57-4028-ad9d-e7695c3a1590.jpg',
        'Mini Caesar salad with salmon  - romaine lettuce, grilled salmon, tomatoes, garlic croutons, shaved asiago, caesar dressing ',
        150),


       ('01aa44e3-0000-45eb-80a1-d8b4e8a22f7b', '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b',
        '99fb44e3-1f18-02eb-20a1-d8b4e8a22f1b', 20, 'Pumpkin soup',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/81828995-ec4a-42ef-9aef-7dc3ca50ace6.jpg',
        'Pumpkin soup pumpkin, onion, garlic, milk ', 220),
       ('02aa44e3-0000-45eb-80a1-d8b4e8a22f7b', '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b',
        '99fb44e3-1f18-02eb-20a1-d8b4e8a22f1b', 19, 'Tomato soup',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/73979981-530a-4c1d-9021-f20d23e97a67.jpg',
        'Tomato soup Combine butter, onion, and tomatoes', 250),
       ('03aa44e3-0000-45eb-80a1-d8b4e8a22f7b', '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b',
        '77fb44e3-3f18-03eb-82a1-d8b4e8a22f1b', 23, 'Pasta Carbonara',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/cfded70b-3754-4d7a-941c-95a5792431a3.jpg',
        'Pasta Carbonara - spaghetti, cream, chicken, spices, pepper ', 220),
       ('04aa44e3-0000-45eb-80a1-d8b4e8a2271b', '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b',
        '77fb44e3-3f18-03eb-82a1-d8b4e8a22f1b', 24, 'Spaghetti with turkey',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/504b6bb1-52b8-473f-b6cc-464c3abfda25.jpg',
        'Spaghetti with turkey -   spaghetti, chicken, ketchup, basil, spices ', 250),
       ('05aa44e3-0000-45eb-80a1-d8b4e8a22f7b', '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b',
        '77fb44e3-3f18-04eb-82a1-d8b4e8a22f1b', 16.90, 'Filet mignon steak',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/c45b77a2-4865-476c-8be8-da23d432d0f4.jpeg',
        'Filet mignon steak  signature spice rubbed, grass-finished steak, grilled medium rare', 300),
       ('06aa44e3-0000-45eb-80a1-d8b4e8a22f7b', '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b',
        '77fb44e3-3f18-04eb-82a1-d8b4e8a22f1b', 10.30, 'Pork steaks',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/c9b2d180-99ab-4fe6-ba9f-3152eb458ab2.jpg',
        'Pork steaks -  grass-finished steak grilled medium rare, drizzled with our fresh chimichurri sauce ', 200),
       ('07aa44e3-0000-45eb-80a1-d8b4e8a22f7b', '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 16, 'Vegetable snack',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/85d69d91-684e-4544-911b-6d4874c112b2.jpg',
        'Vegetable snack - cucumbers, celery, mint, pickled onions, raw walnuts, grapes, blue cheese, apple, carrot, pepper ',
        200),
       ('08aa44e3-0000-45eb-80a1-d8b4e8a22f7b', '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 17, 'Carpaccio',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/3baa1429-f9e9-4a76-83c9-612b921c0aa7.jpg',
        'Carpaccio-  fresh beef tenderloin, arugula, watercress, shaved Parmigiano-Reggiano cheese,  freshly ground pepper ',
        110),

       ('01aa44e3-0000-45eb-80a1-d8b4e8a22f8b', '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b',
        '99fb44e3-1f18-01eb-20a1-d8b4e8a22f1b', 15, 'Curry with vegetables',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/c0569f0f-798f-40e0-b2de-9d5b60c5d02e.jpg',
        'Curry with vegetables - tomatoes, cucumbers, Iceberg lettuce on cheese roll, Mozzarella chees, curry. ketchup ',
        110),
       ('02aa44e3-0000-45eb-80a1-d8b4e8a22f8b', '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b',
        '99fb44e3-1f18-01eb-20a1-d8b4e8a22f1b', 17, 'Curry with chicken',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/56617f40-b872-4ae3-9c7e-01348697c9d8.jpg',
        'Curry with chicken -  curry, onion, tomatoes, chicken, ketchup, basil ', 150),
       ('03aa44e3-0000-45eb-80a1-d8b4e8a22f8b', '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b',
        '99fb44e3-1f18-02eb-20a1-d8b4e8a22f1b', 20, 'Pumpkin soup',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/81828995-ec4a-42ef-9aef-7dc3ca50ace6.jpg',
        'Pumpkin soup pumpkin, onion, garlic, milk ', 220),
       ('04aa44e3-0000-45eb-80a1-d8b4e8a22f8b', '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b',
        '99fb44e3-1f18-02eb-20a1-d8b4e8a22f1b', 19, 'Tomato soup',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/73979981-530a-4c1d-9021-f20d23e97a67.jpg',
        'Tomato soup Combine butter, onion, and tomatoes ', 250),
       ('05aa44e3-0000-45eb-80a1-d8b4e8a22f8b', '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b',
        '99fb44e3-1f18-03eb-20a1-d8b4e8a22f1b', 16.90, 'Udo noodles with bacon',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/27aebbb5-5b9c-4163-bac1-d8f21394cc80.jpg',
        'Udo noodles with bacon - bacon, tomato, eggs, udon noodles, sliced vegetables, a protein and a savoury sauce ',
        190),
       ('06aa44e3-0000-45eb-80a1-d8b4e8a22f8b', '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b',
        '99fb44e3-1f18-03eb-20a1-d8b4e8a22f1b', 17.90, 'Shrimp udo noodles',
        'https://onlineshop.fra1.digitaloceanspaces.com/01aa44e3-0000-45aa-80a2-d8b4e8a22a1b/aba3e2fb-911a-46ef-8d2f-b04019f684a7.jpg',
        'Shrimp udo noodles - udon noodles, sliced vegetables, a protein and a savoury sauce ', 170),
       ('07aa44e3-0000-45eb-80a1-d8b4e8a22f8b', '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 16, 'Vegetable snack',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/85d69d91-684e-4544-911b-6d4874c112b2.jpg',
        'Vegetable snack - cucumbers, celery, mint, pickled onions, raw walnuts, grapes, blue cheese, apple, carrot, pepper ',
        200),
       ('08aa44e3-0000-45eb-80a1-d8b4e8a22f8b', '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b',
        '11fb44e3-5f18-04eb-80a1-d8b4e8a22f1b', 17, 'Carpaccio',
        'https://onlineshop.fra1.digitaloceanspaces.com/01fb44e3-0000-45eb-80a1-d8b4e8a22f1b/3baa1429-f9e9-4a76-83c9-612b921c0aa7.jpg',
        'Carpaccio - fresh beef tenderloin, arugula, watercress, shaved Parmigiano-Reggiano cheese,  freshly ground pepper ',
        110);

INSERT INTO orders (id, restaurant_id, payment_type, cost, delivery_time, address, client_phone_number,
                    client_full_name)
VALUES ('01ff44e3-5f18-0000-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 32, '2016-06-22 19:10:25-07',
        'ул. Новоуфимская 11, Минск', '291111111', 'Иван Иванов'),
       ('02ff44e3-5f18-0001-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 15, '2020-06-22 19:10:25-07',
        'улица Притыцкого 29', '291111112', 'Петя Петров'),
       ('03ff44e3-5f18-0002-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 15.3,
        '2022-02-22 19:10:25-07',
        'улица Притыцкого 30', '291111113', 'Катя Костевич'),
       ('01ff44e3-5f18-0004-80a1-d8b4e8a22f1b', '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 32, '2016-06-22 19:10:25-07',
        'ул. Новоуфимская 6, Минск', '291111111', 'Саня Санчез'),
       ('02ff44e3-5f18-0005-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 15, '2020-06-22 19:10:25-07',
        'улица Притыцкого 5', '291111112', 'Ваня Вантуз'),
       ('03ff44e3-5f18-0006-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 15.3,
        '2022-02-22 19:10:25-07',
        'улица Притыцкого 4', '291111113', 'Дима Димыч'),
       ('02ff44e3-5f18-0007-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 15, '2020-06-22 19:10:25-07',
        'улица Притыцкого 3', '291111112', 'Гена Геныч'),
       ('02ff44e3-5f18-0008-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 15, '2020-06-22 19:10:25-07',
        'улица Притыцкого 2', '291111112', 'Афанасий Афанасий'),
       ('02ff44e3-5f18-0009-80a1-d8b4e8a22f1b', '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b', 2, 15, '2020-06-22 19:10:25-07',
        'улица Притыцкого 1', '291111112', 'Ира Смелая');

INSERT INTO order_dishes (order_id, dish_id, amount)
VALUES ('01ff44e3-5f18-0000-80a1-d8b4e8a22f1b', '01ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 2),
       ('01ff44e3-5f18-0000-80a1-d8b4e8a22f1b', '02ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 2),

       ('02ff44e3-5f18-0001-80a1-d8b4e8a22f1b', '03ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 3),

       ('03ff44e3-5f18-0002-80a1-d8b4e8a22f1b', '04ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 1),
       ('03ff44e3-5f18-0002-80a1-d8b4e8a22f1b', '05ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 1),
       ('03ff44e3-5f18-0002-80a1-d8b4e8a22f1b', '06ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 1),
       ('03ff44e3-5f18-0002-80a1-d8b4e8a22f1b', '07ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 1),
       ('01ff44e3-5f18-0004-80a1-d8b4e8a22f1b', '01ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 2),
       ('02ff44e3-5f18-0005-80a1-d8b4e8a22f1b', '02ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 2),

       ('03ff44e3-5f18-0006-80a1-d8b4e8a22f1b', '03ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 3),

       ('02ff44e3-5f18-0007-80a1-d8b4e8a22f1b', '04ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 1),
       ('02ff44e3-5f18-0008-80a1-d8b4e8a22f1b', '05ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 1),
       ('02ff44e3-5f18-0009-80a1-d8b4e8a22f1b', '06ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 1),
       ('02ff44e3-5f18-0009-80a1-d8b4e8a22f1b', '07ff44e3-0000-45eb-80a1-d8b4e8a22f1b', 1);

INSERT INTO category_restaurants(category_id, restaurant_id)
VALUES (1, '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b'),
       (2, '01fb44e3-5f18-45eb-80a1-d8b4e8a22f1b'),
       (1, '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b'),
       (2, '02fb44e3-5f18-45eb-80a1-d8b4e8a22f1b'),
       (3, '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b'),
       (5, '03fb44e3-5f18-45eb-80a1-d8b4e8a22f1b'),
       (2, '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b'),
       (3, '04fb44e3-5f18-41eb-80a1-d8b4e8a99f1b'),
       (1, '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b'),
       (5, '05fb44e3-3f18-43eb-82a1-d8b4e8a22f1b'),
       (4, '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b'),
       (4, '04fb44e3-5f18-41eb-80a1-d8b4e8a77f1b'),
       (4, '04fb44e3-5f18-41eb-80a1-d8b4e8a11f1b'),
       (5, '06fb44e3-1f18-45eb-20a1-d8b4e8a22f1b');