-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile
DROP COLUMN bus_info_id_photo_submitted,
DROP COLUMN bus_info_picture_submitted,
DROP COLUMN bus_info_nbi_clearance_submitted,
DROP COLUMN bus_info_court_clearance_submitted,
DROP COLUMN bus_info_incorporation_papers_submitted,
DROP COLUMN bus_info_mayors_permit_submitted,
DROP COLUMN bus_info_id_photo_date_checked,
DROP COLUMN bus_info_picture_date_checked,
DROP COLUMN bus_info_nbi_clearance_date_checked,
DROP COLUMN bus_info_court_clearance_date_checked,
DROP COLUMN bus_info_incorporation_papers_date_checked,
DROP COLUMN bus_info_mayors_permit_date_checked,
DROP COLUMN fin_info_financial_statement_submitted,
DROP COLUMN fin_info_bank_statement_submitted,
DROP COLUMN fin_info_financial_statement_date_checked,
DROP COLUMN fin_info_bank_statement_date_checked,
DROP COLUMN drp_info_questionnaire_submitted,
DROP COLUMN drp_info_questionnaire_date_checked;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
