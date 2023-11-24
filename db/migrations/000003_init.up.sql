CREATE SEQUENCE public.transfer_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

ALTER TABLE transfer ALTER COLUMN id TYPE bigint;
ALTER TABLE transfer ALTER COLUMN id SET NOT NULL;
ALTER TABLE transfer ALTER COLUMN id SET DEFAULT nextval('public.transfer_id_seq');
ALTER SEQUENCE public.transfer_id_seq OWNED BY transfer.id;
