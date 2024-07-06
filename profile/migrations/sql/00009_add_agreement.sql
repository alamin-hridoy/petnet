-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile 
ADD COLUMN acc_info_agree_terms_conditions BOOLEAN DEFAULT false,
ADD COLUMN acc_info_agree_online_supplier_form BOOLEAN DEFAULT false;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE org_profile
DROP COLUMN acc_info_agree_terms_conditions,
DROP COLUMN acc_info_agree_online_supplier_form;
