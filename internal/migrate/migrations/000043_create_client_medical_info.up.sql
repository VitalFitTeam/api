-- Create client_medical_info table
CREATE TABLE IF NOT EXISTS client_medical_info (
    medical_info_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    medical_conditions TEXT,
    medical_risks TEXT,
    warnings TEXT,
    allergies TEXT,
    medications TEXT,
    emergency_contact TEXT,
    blood_type VARCHAR(10),
    last_updated_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_last_updated_by
        FOREIGN KEY (last_updated_by)
        REFERENCES users(user_id)
        ON DELETE RESTRICT
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_client_medical_info_user_id ON client_medical_info(user_id);

-- Create index on deleted_at for soft delete queries
CREATE INDEX IF NOT EXISTS idx_client_medical_info_deleted_at ON client_medical_info(deleted_at);

-- Create medical_info_audit_log table
CREATE TABLE IF NOT EXISTS medical_info_audit_log (
    audit_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_info_id UUID NOT NULL,
    user_id UUID NOT NULL,
    modified_by UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    fields_changed TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_medical_info_id
        FOREIGN KEY (medical_info_id)
        REFERENCES client_medical_info(medical_info_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_audit_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_modified_by
        FOREIGN KEY (modified_by)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);

-- Create indexes for audit log queries
CREATE INDEX IF NOT EXISTS idx_medical_audit_user_id ON medical_info_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_medical_audit_medical_info_id ON medical_info_audit_log(medical_info_id);
CREATE INDEX IF NOT EXISTS idx_medical_audit_created_at ON medical_info_audit_log(created_at DESC);

-- Add comment to tables
COMMENT ON TABLE client_medical_info IS 'Stores encrypted medical information for clients';
COMMENT ON TABLE medical_info_audit_log IS 'Audit trail for all medical information access and modifications';
