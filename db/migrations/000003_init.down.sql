ALTER TABLE transfer ALTER COLUMN id TYPE integer;
ALTER TABLE transfer ALTER COLUMN id DROP DEFAULT;

DROP SEQUENCE public.transfer_id_seq;
