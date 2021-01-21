CREATE TABLE employee (
    id uuid NOT NULL
);

CREATE TABLE project (
    id integer NOT NULL,
    name text,
    state smallint,
    progress numeric,
    owner uuid
);


CREATE TABLE project_employee (
    project_id bigint NOT NULL,
    employee_id uuid NOT NULL
);



CREATE SEQUENCE project_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE project_id_seq OWNED BY project.id;

ALTER TABLE ONLY project ALTER COLUMN id SET DEFAULT nextval('project_id_seq'::regclass);


ALTER TABLE ONLY employee
    ADD CONSTRAINT employee_pkey PRIMARY KEY (id);


ALTER TABLE ONLY project_employee
    ADD CONSTRAINT project_employee_pkey PRIMARY KEY (project_id, employee_id);

ALTER TABLE ONLY project
    ADD CONSTRAINT project_name_unique UNIQUE (name);

ALTER TABLE ONLY project
    ADD CONSTRAINT project_pkey PRIMARY KEY (id);


ALTER TABLE ONLY project_employee
    ADD CONSTRAINT employee_id FOREIGN KEY (employee_id) REFERENCES employee(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;

ALTER TABLE ONLY project
    ADD CONSTRAINT owner FOREIGN KEY (owner) REFERENCES employee(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


ALTER TABLE ONLY project_employee
    ADD CONSTRAINT project_id FOREIGN KEY (project_id) REFERENCES project(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;
