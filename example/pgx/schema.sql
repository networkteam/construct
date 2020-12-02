CREATE TABLE projects
(
    id    uuid NOT NULL
        CONSTRAINT projects_pk PRIMARY KEY,
    title text NOT NULL
);

CREATE TABLE todos
(
    id           uuid NOT NULL
        CONSTRAINT todos_pk PRIMARY KEY,
    project_id   uuid NOT NULL
        CONSTRAINT todos_projects_id_fk REFERENCES projects (id),
    title        text NOT NULL,
    completed_at timestamp WITH TIME ZONE
);
