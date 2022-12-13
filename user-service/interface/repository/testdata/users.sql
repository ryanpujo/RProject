CREATE TABLE public.users (
  id integer NOT NULL PRIMARY KEY,
  first_name character varying(255),
  last_name character varying(255),
  username character varying(255),
  email character varying(255),
  password character varying(255),
  created_at timestamp without time zone,
  updated_at timestamp without time zone
);

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
  SEQUENCE NAME public.users_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1
)