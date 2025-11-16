-- +migrate Up
CREATE TYPE refund_status AS ENUM ('Pending', 'Processed', 'Failed');

CREATE TABLE refunds (
    refund_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_id UUID NOT NULL REFERENCES payments(payment_id),
    invoice_id UUID REFERENCES invoices(invoice_id),
    reason VARCHAR(255) NOT NULL,
    amount_refunded DECIMAL(10, 2) NOT NULL,
    currency_refunded VARCHAR(3) NOT NULL,
    amount_base DECIMAL(10, 2) NOT NULL,
    currency_base VARCHAR(3) NOT NULL DEFAULT 'USD',
    exchange_rate DECIMAL(18, 8) NOT NULL,    
    refund_method_id UUID REFERENCES payment_methods(method_id),
    transaction_id VARCHAR(255),
    status refund_status NOT NULL DEFAULT 'Pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_invoice_id ON refunds(invoice_id);
CREATE INDEX idx_refunds_status ON refunds(status);