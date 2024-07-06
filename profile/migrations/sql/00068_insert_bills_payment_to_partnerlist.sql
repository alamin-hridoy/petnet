-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO partner_list ( stype , name , status, service_name ) VALUES 
('ECP', 'Ecpay', 'ENABLED', 'BILLSPAYMENT'), 
('BYC', 'BayadCenter', 'ENABLED', 'BILLSPAYMENT'), 
('MLP', 'Multipay', 'ENABLED', 'BILLSPAYMENT');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DELETE FROM partner_list
WHERE stype IN('ECP', 'BYC', 'MLP');