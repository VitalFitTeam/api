-- Add new booking statuses for check-in and waitlist functionality
ALTER TYPE booking_status_enum ADD VALUE IF NOT EXISTS 'CheckedIn';
ALTER TYPE booking_status_enum ADD VALUE IF NOT EXISTS 'Waitlist';
