-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE org_profile
DROP COLUMN bus_info_id_photo_submitted,
DROP COLUMN bus_info_picture_submitted,
DROP COLUMN bus_info_nbi_clearance_submitted, 
DROP COLUMN bus_info_court_clearance_submitted,
DROP COLUMN bus_info_incorporation_papers_submitted,
DROP COLUMN bus_info_mayors_permit_submitted,
DROP COLUMN fin_info_financial_statement_submitted, 
DROP COLUMN fin_info_bank_statement_submitted, 
DROP COLUMN drp_info_questionnaire_submitted, 
DROP COLUMN acc_info_agree_terms_conditions, 
DROP COLUMN acc_info_agree_online_supplier_form, 

ADD COLUMN bus_info_id_photo_submitted smallint DEFAULT 0,
ADD COLUMN bus_info_picture_submitted smallint DEFAULT 0,
ADD COLUMN bus_info_nbi_clearance_submitted smallint DEFAULT 0,
ADD COLUMN bus_info_court_clearance_submitted smallint DEFAULT 0,
ADD COLUMN bus_info_incorporation_papers_submitted smallint DEFAULT 0,
ADD COLUMN bus_info_mayors_permit_submitted smallint DEFAULT 0,
ADD COLUMN fin_info_financial_statement_submitted smallint DEFAULT 0,
ADD COLUMN fin_info_bank_statement_submitted smallint DEFAULT 0,
ADD COLUMN drp_info_questionnaire_submitted smallint DEFAULT 0,
ADD COLUMN acc_info_agree_terms_conditions smallint DEFAULT 0,
ADD COLUMN acc_info_agree_online_supplier_form smallint DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
