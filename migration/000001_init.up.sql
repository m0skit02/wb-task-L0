-- Таблица заказов
CREATE TABLE orders (
                        order_uid         VARCHAR PRIMARY KEY,
                        track_number      VARCHAR NOT NULL,
                        entry             VARCHAR NOT NULL,
                        locale            VARCHAR(5),
                        internal_signature TEXT,
                        customer_id       VARCHAR NOT NULL,
                        delivery_service  VARCHAR,
                        shard_key          VARCHAR,
                        sm_id             INTEGER,
                        date_created      TIMESTAMP WITH TIME ZONE NOT NULL,
                        oof_shard         VARCHAR
);

-- Таблица доставки
CREATE TABLE deliveries (
                          delivery_id   VARCHAR PRIMARY KEY,
                          order_uid     VARCHAR NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                          name          VARCHAR NOT NULL,
                          phone         VARCHAR,
                          zip           VARCHAR,
                          city          VARCHAR,
                          address       TEXT,
                          region        VARCHAR,
                          email         VARCHAR
);

-- Таблица оплаты
CREATE TABLE payments (
                         payment_id    VARCHAR PRIMARY KEY,
                         order_uid     VARCHAR NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                         transaction   VARCHAR NOT NULL,
                         request_id    VARCHAR,
                         currency      VARCHAR(10) NOT NULL,
                         provider      VARCHAR NOT NULL,
                         amount        NUMERIC(12,2) NOT NULL,
                         payment_dt    BIGINT,
                         bank          VARCHAR,
                         delivery_cost NUMERIC(12,2),
                         goods_total   NUMERIC(12,2),
                         custom_fee    NUMERIC(12,2)
);

-- Таблица товаров
CREATE TABLE items (
                       item_id       VARCHAR PRIMARY KEY,
                       order_uid     VARCHAR NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                       chrt_id       BIGINT NOT NULL,
                       track_number  VARCHAR,
                       price         NUMERIC(12,2) NOT NULL,
                       rid           VARCHAR,
                       name          VARCHAR NOT NULL,
                       sale          NUMERIC(5,2),
                       size          VARCHAR,
                       total_price   NUMERIC(12,2),
                       nm_id         BIGINT,
                       brand         VARCHAR,
                       status        INTEGER
);
