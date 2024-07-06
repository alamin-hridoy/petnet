-- +goose Up
-- SQL in this section is executed when the migration is applied.
UPDATE partner_list set service_name = 'CASHINCASHOUT' WHERE stype IN (
  'GCASH', 'DRAGONPAY', 'PAYMAYA', 'COINS', 'PERAHUB', 'DISKARTECH');
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
