BEGIN;

CREATE TABLE mch.project
(
    id              uuid primary key DEFAULT uuid_generate_v4(),
    name            text not null,
    image           text,
    industry        text not null,
    owner_id        uuid references mch.user (id) ON DELETE CASCADE,
    description     text,
    readiness_stage text not null,
    tags            text[],
    launch_date     bigint
);

CREATE TABLE mch.user_project
(
    project_id uuid references mch.project (id) ON DELETE CASCADE,
    user_id    uuid references mch.user (id) ON DELETE CASCADE,
    viewed     boolean DEFAULT true,
    applied    boolean DEFAULT false,
    confirmed  boolean DEFAULT false,
    rejected   boolean DEFAULT false,
    PRIMARY KEY (project_id, user_id)
);

CREATE TABLE mch.user_project_viewed
(
    project_id uuid references mch.project (id) ON DELETE CASCADE,
    user_id    uuid references mch.user (id) ON DELETE CASCADE,
    confirmed  boolean DEFAULT false,
    PRIMARY KEY (project_id, user_id)
);

COMMIT;