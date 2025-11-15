CREATE TABLE fiscal_document_types (
    document_type_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    prefix VARCHAR(10) UNIQUE NOT NULL
);

INSERT INTO fiscal_document_types (name, prefix)
VALUES
    ('Factura', 'F-'),
    ('Boleta', 'B-'),
    ('Nota de Crédito', 'NC-'),
    ('Nota de Débito', 'ND-'),
    ('Recibo', 'R-');