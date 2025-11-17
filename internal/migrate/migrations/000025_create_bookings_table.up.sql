CREATE TYPE IF NOT EXISTS booking_status_enum AS ENUM ('Confirmed', 'CancelledByUser', 'CancelledBySystem');

CREATE TABLE IF NOT EXISTS bookings (
    booking_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    class_id UUID NOT NULL,
    status booking_status_enum NOT NULL DEFAULT 'Confirmed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (class_id) REFERENCES classes(class_id),
    CONSTRAINT unique_user_class_booking UNIQUE (user_id, class_id)
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_class ON bookings (user_id, class_id);
