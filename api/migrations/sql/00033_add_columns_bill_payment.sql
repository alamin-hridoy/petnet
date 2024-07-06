-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE bill_payment
    ADD COLUMN IF NOT EXISTS client_reference_number text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS bill_partner_id text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS partner_charge text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS reference_number text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS validation_number text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS receipt_validation_number text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS tpa_id text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS type text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS txnid text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE bill_payment
DROP COLUMN client_reference_number,
DROP COLUMN bill_partner_id,
DROP COLUMN partner_charge,
DROP COLUMN reference_number,
DROP COLUMN validation_number,
DROP COLUMN receipt_validation_number,
DROP COLUMN tpa_id,
DROP COLUMN type,
DROP COLUMN txnid;
