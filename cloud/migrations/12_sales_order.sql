CREATE TABLE sales_orders (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    order_no VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    order_date TIMESTAMPTZ NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_sales_orders_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_sales_orders_customer_id
    ON sales_orders(customer_id);

CREATE INDEX idx_sales_orders_status
    ON sales_orders(status);

CREATE INDEX idx_sales_orders_order_date
    ON sales_orders(order_date);


CREATE TABLE sales_order_items (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    ordered_quantity BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_sales_order_items_sales_order
        FOREIGN KEY (sales_order_id)
        REFERENCES sales_orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_sales_order_items_product
        FOREIGN KEY (product_id)
        REFERENCES product_models(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_sales_order_items_sales_order_id
    ON sales_order_items(sales_order_id);

CREATE INDEX idx_sales_order_items_product_id
    ON sales_order_items(product_id);