-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS revenue_sharing_report (
   id uuid NOT NULL DEFAULT uuid_generate_v4(),
   org_id uuid NOT NULL,
   dsa_code text NOT NULL DEFAULT '',
   year_month text NOT NULL DEFAULT '',
   remittance_count int	DEFAULT 0,
   cico_count int	DEFAULT 0,
   bills_payment_count int	DEFAULT 0,
   insurance_count int	DEFAULT 0,
   dsa_commission text NOT NULL DEFAULT '',
   dsa_commission_type text NOT NULL DEFAULT '',
   status smallint	DEFAULT 0,
   created timestamptz NOT NULL DEFAULT NOW(),
   updated timestamptz NOT NULL DEFAULT NOW(),
   CONSTRAINT revenue_sharing_report_org_id_year_month UNIQUE (org_id, year_month)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS revenue_sharing_report;
