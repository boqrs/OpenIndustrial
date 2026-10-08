CREATE TABLE resource_credentials (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    secret_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    consumed_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_resource_credentials_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_resource_credentials_resource_id
    ON resource_credentials(resource_id);


CREATE TABLE resource_identities (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    hardware_id VARCHAR(255),
    serial_number VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_resource_identities_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_resource_identities_resource_id
    ON resource_identities(resource_id);

CREATE INDEX idx_resource_identities_hardware_id
    ON resource_identities(hardware_id);

CREATE INDEX idx_resource_identities_serial_number
    ON resource_identities(serial_number);


CREATE TABLE resource_certificates (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    certificate_id VARCHAR(512) NOT NULL,
    certificate_serial_number VARCHAR(255),
    fingerprint VARCHAR(255) NOT NULL UNIQUE,
    subject TEXT,
    issuer TEXT,
    status VARCHAR(50) NOT NULL,
    not_before TIMESTAMPTZ NOT NULL,
    not_after TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    activated_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,

    CONSTRAINT fk_resource_certificates_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_resource_certificates_resource_id
    ON resource_certificates(resource_id);

CREATE INDEX idx_resource_certificates_certificate_serial_number
    ON resource_certificates(certificate_serial_number);