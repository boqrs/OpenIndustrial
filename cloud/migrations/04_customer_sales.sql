CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(32) NOT NULL,
    contact_name VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(64),
    address TEXT,
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customers_name
    ON customers(name);

CREATE INDEX idx_customers_type
    ON customers(type);

CREATE INDEX idx_customers_status
    ON customers(status);

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

CREATE TABLE production_plan_allocations (
    id BIGSERIAL PRIMARY KEY,
    production_plan_id BIGINT NOT NULL,
    sales_order_item_id BIGINT NOT NULL,
    allocated_quantity BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_production_plan_allocations_plan
        FOREIGN KEY (production_plan_id)
        REFERENCES production_plans(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_production_plan_allocations_sales_order_item
        FOREIGN KEY (sales_order_item_id)
        REFERENCES sales_order_items(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_production_plan_allocations_plan_id
    ON production_plan_allocations(production_plan_id);

CREATE INDEX idx_production_plan_allocations_sales_order_item_id
    ON production_plan_allocations(sales_order_item_id);