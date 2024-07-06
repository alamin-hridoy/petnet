-- +goose Up
-- SQL in this section is executed when the migration is applied.
UPDATE partner_list set service_name = 'REMITTANCE' WHERE stype IN (
  'CEB', 'IE', 'IR', 'TF', 'RIA', 'MB', 'RM', 'BPI', 'USSC', 'JPR', 'IC', 'UNT', 'AYA', 'WU', 'WISE', 'CEBI');
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
