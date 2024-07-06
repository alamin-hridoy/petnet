-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE remit_cache RENAME COLUMN remittance TO remit;
ALTER TABLE remit_cache RENAME COLUMN remit_type TO partner_remit_type;
ALTER TABLE remit_cache DROP COLUMN dsa_order_id;
ALTER TABLE remit_cache ADD COLUMN remit_type text NOT NULL DEFAULT '';
ALTER TABLE remit_history ADD COLUMN remit_type text NOT NULL DEFAULT '';
UPDATE remit_history
SET remit_type = 'DISBURSE'
WHERE remittance->>'txn_type' = 'PO';
UPDATE remit_history
SET remit_type = 'CREATE'
WHERE remittance->>'txn_type' = 'SO';
UPDATE remit_cache
SET remit_type = 'DISBURSE'
WHERE partner_remit_type = 'PO';
UPDATE remit_cache
SET remit_type = 'CREATE'
WHERE partner_remit_type = 'SO';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
