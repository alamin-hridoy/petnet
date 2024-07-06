-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO partner_list ( stype , name , status, service_name )
VALUES
    ('RLN', 'RuralNet', 'ENABLED', 'MICROINSURANCE');
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
