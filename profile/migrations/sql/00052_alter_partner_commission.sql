-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE partner_commission_config ADD CONSTRAINT partner_commission_config_remit_type_bound_type_partner_tra_uk UNIQUE (remit_type, bound_type, partner, transaction_type, tier_type);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE partner_commission_config DROP CONSTRAINT partner_commission_config_remit_type_bound_type_partner_tra_uk;