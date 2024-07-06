-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS mltf_risk_assesment_question (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    org_id uuid NOT NULL,
    qid text NOT NULL,
    qtype text NOT NULL,
    customers_total text DEFAULT NULL,
    hr_total text DEFAULT NULL,
    impact_score text DEFAULT NULL,
    updated timestamptz NOT NULL DEFAULT now(),
    created timestamptz NOT NULL DEFAULT now()
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS mltf_risk_assesment_question;
