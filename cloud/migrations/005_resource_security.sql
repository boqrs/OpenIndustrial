-- +migrate Up
-- This migration establishes the "Resource Security Kernel".
-- It creates tables for managing credentials, identities, and certificates for any resource,
-- separating security concerns from the core resource model.
-- As per the finalized architecture, these tables do not contain a redundant TenantID,
-- as tenancy is owned by the parent Resource.

-- 1. Credential Enums
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'credential_type') THEN
        CREATE TYPE credential_type AS ENUM ('bootstrap');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'credential_status') THEN
        CREATE TYPE credential_status AS ENUM ('active', 'consumed', 'revoked');
    END IF;
END$$;

-- 2. Certificate Enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'certificate_status') THEN
        CREATE TYPE certificate_status AS ENUM ('pending', 'active', 'revoked', 'expired');
    END IF;
END$$;

-- 3. Resource Credentials Table
-- Stores bootstrap/provisioning credentials for resources.
CREATE TABLE IF NOT EXISTS resource_credentials (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    type credential_type NOT NULL,
    status credential_status NOT NULL,
    secret_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    consumed_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_resource_credentials_resource
        FOREIGN KEY(resource_id) 
        REFERENCES resources(uuid)
        ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_resource_credentials_resource_id ON resource_credentials(resource_id);

-- 4. Resource Identities Table
-- Stores intrinsic, hardware-based identifiers for resources.
CREATE TABLE IF NOT EXISTS resource_identities (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL UNIQUE,
    hardware_id VARCHAR(255),
    serial_number VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_resource_identities_resource
        FOREIGN KEY(resource_id) 
        REFERENCES resources(uuid)
        ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_resource_identities_hardware_id ON resource_identities(hardware_id);
CREATE INDEX IF NOT EXISTS idx_resource_identities_serial_number ON resource_identities(serial_number);

-- 5. Resource Certificates Table
-- Stores X.509 certificate information associated with resources for TLS/MQTT authentication.
CREATE TABLE IF NOT EXISTS resource_certificates (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    certificate_id  BIGINT NOT NULL,
    certificate_serial_number VARCHAR(255),
    fingerprint VARCHAR(255) NOT NULL UNIQUE,
    subject TEXT,
    issuer TEXT,
    status certificate_status NOT NULL,
    not_before TIMESTAMPTZ NOT NULL,
    not_after TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    activated_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,

    CONSTRAINT fk_resource_certificates_resource
        FOREIGN KEY(resource_id) 
        REFERENCES resources(uuid)
        ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_resource_certificates_resource_id ON resource_certificates(resource_id);
CREATE INDEX IF NOT EXISTS idx_resource_certificates_serial_number ON resource_certificates(certificate_serial_number);


-- +migrate Down
DROP TABLE IF EXISTS resource_certificates;
DROP TABLE IF EXISTS resource_identities;
DROP TABLE IF EXISTS resource_credentials;

DROP TYPE IF EXISTS certificate_status;
DROP TYPE IF EXISTS credential_status;
DROP TYPE IF EXISTS credential_type;