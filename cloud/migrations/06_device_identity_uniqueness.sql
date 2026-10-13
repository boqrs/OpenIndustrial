-- Each resource can have at most one canonical identity.
CREATE UNIQUE INDEX IF NOT EXISTS
    uq_resource_identities_resource_id
ON resource_identities(resource_id);


CREATE UNIQUE INDEX IF NOT EXISTS uq_resource_identities_hardware_id
    ON resource_identities (BTRIM(hardware_id))
    WHERE hardware_id IS NOT NULL
      AND BTRIM(hardware_id) <> '';

CREATE UNIQUE INDEX IF NOT EXISTS uq_resource_identities_serial_number
    ON resource_identities (BTRIM(serial_number))
    WHERE serial_number IS NOT NULL
      AND BTRIM(serial_number) <> '';