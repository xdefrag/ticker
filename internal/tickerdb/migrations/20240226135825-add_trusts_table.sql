
-- +migrate Up
CREATE TABLE trusts (
    code text NOT NULL,
    issuer_account text NOT NULL,
    source text NOT NULL,
    priority bigint NOT NULL,
    updated_at timestamptz NOT NULL
);

ALTER TABLE ONLY public.trusts
    ADD CONSTRAINT trusts_code_issuer_account_unique UNIQUE (code, issuer_account);

-- +migrate Down
DROP TABLE trusts;
