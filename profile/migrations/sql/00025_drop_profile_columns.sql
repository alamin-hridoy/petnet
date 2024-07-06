-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile
DROP COLUMN wu_coy,
DROP COLUMN wu_operator_id,
DROP COLUMN wu_terminal_id,
DROP COLUMN bus_info_id_photo_urls,
DROP COLUMN bus_info_picture_urls,
DROP COLUMN bus_info_nbi_clearance_urls,
DROP COLUMN bus_info_court_clearance_url,
DROP COLUMN bus_info_incorporation_paper_urls,
DROP COLUMN bus_info_mayors_permit_url,
DROP COLUMN fin_info_financial_statement_urls,
DROP COLUMN fin_info_bank_statement_urls,
DROP COLUMN drp_info_service,
DROP COLUMN drp_info_questionnaire_urls;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
