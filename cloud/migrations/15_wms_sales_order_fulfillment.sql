ALTER TABLE sales_order_items
    ADD COLUMN IF NOT EXISTS reserved_quantity BIGINT NOT NULL DEFAULT 0;

ALTER TABLE sales_order_items
    ADD COLUMN IF NOT EXISTS fulfilled_quantity BIGINT NOT NULL DEFAULT 0;

ALTER TABLE sales_order_items
    ADD CONSTRAINT chk_sales_order_items_fulfillment_quantity
    CHECK (
        ordered_quantity > 0
        AND reserved_quantity >= 0
        AND fulfilled_quantity >= 0
        AND reserved_quantity + fulfilled_quantity <= ordered_quantity
    );


ALTER TABLE shipments
    ADD CONSTRAINT fk_shipments_sales_order
    FOREIGN KEY (sales_order_id)
    REFERENCES sales_orders(id)
    ON DELETE RESTRICT;


ALTER TABLE shipment_items
    ADD CONSTRAINT fk_shipment_items_device
    FOREIGN KEY (device_id)
    REFERENCES devices(id)
    ON DELETE RESTRICT;

ALTER TABLE shipment_items
    ADD CONSTRAINT fk_shipment_items_sales_order_item
    FOREIGN KEY (sales_order_item_id)
    REFERENCES sales_order_items(id)
    ON DELETE RESTRICT;


ALTER TABLE shipment_tracking_events
    ADD CONSTRAINT fk_shipment_tracking_events_shipment
    FOREIGN KEY (shipment_id)
    REFERENCES shipments(id)
    ON DELETE CASCADE;