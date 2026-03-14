-- Creating the usecase schema's objects.
CREATE SCHEMA IF NOT EXISTS usecase;

CREATE TABLE IF NOT EXISTS usecase.divisions (
    id                SERIAL PRIMARY KEY,
    name              TEXT NOT NULL,
    type              TEXT NOT NULL CHECK (type in ('division', 'directorate', 'department', 'unit', 'group')),
    state_size        INT NOT NULL,
    superdivision_id  INT REFERENCES usecase.divisions(id) ON DELETE RESTRICT,
    UNIQUE(name, type)
);

CREATE TABLE IF NOT EXISTS usecase.titles (
    id    SERIAL PRIMARY KEY,
    name  TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS usecase.citizenships (
    id    SERIAL PRIMARY KEY,
    name  TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS usecase.employees (
    id              SERIAL PRIMARY KEY,
    tin_num         TEXT UNIQUE NOT NULL,
    snils_num       TEXT UNIQUE NOT NULL,
    passport_data   TEXT UNIQUE NOT NULL,
    phone_num       TEXT UNIQUE NOT NULL,
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    patronymic      TEXT,
    address         TEXT NOT NULL,
    title_id        INT REFERENCES usecase.titles(id) ON DELETE RESTRICT,
    hiring_date     DATE NOT NULL,
    unit_id         INT REFERENCES usecase.divisions(id) ON DELETE RESTRICT,
    education       TEXT NOT NULL,
    salary          NUMERIC(10, 2) NOT NULL,
    citizenship_id  INT NOT NULL REFERENCES usecase.citizenships(id) ON DELETE RESTRICT
);

-- Creating the app_realm's schema's objects.
CREATE SCHEMA IF NOT EXISTS app_realm;

CREATE TABLE IF NOT EXISTS app_realm.roles (
    id    SERIAL PRIMARY KEY,
    name  TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS app_realm.rights (
    id    SERIAL PRIMARY KEY,
    name  TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS app_realm.role_rights_mapping (
    role_id   INT REFERENCES app_realm.roles(id) ON DELETE RESTRICT,
    right_id  INT REFERENCES app_realm.rights(id) ON DELETE RESTRICT,
    PRIMARY KEY(role_id, right_id)
);

CREATE TABLE IF NOT EXISTS app_realm.users (
    id        SERIAL PRIMARY KEY,
    login     TEXT UNIQUE NOT NULL,
    hash_pwd  TEXT NOT NULL,
    salt      INT  NOT NULL,
    role      INT REFERENCES app_realm.roles(id) ON DELETE RESTRICT
);

-- Configuring the app_realm's objects.
INSERT INTO app_realm.roles (name)
VALUES ('analyst'),
       ('hr-manager'),
       ('admin');

INSERT INTO app_realm.rights (name)
VALUES ('read-employees'),
       ('write-employees'),
       ('read-divisions'),
       ('write-divisions'),
       ('write-citizenships'),
       ('write-titles'),
       ('read-users'),
       ('write-users');

INSERT INTO app_realm.role_rights_mapping (role_id, right_id)
SELECT app_realm.roles.id, app_realm.rights.id
FROM app_realm.roles CROSS JOIN app_realm.rights
WHERE 
    app_realm.roles.name='analyst' AND (
        app_realm.rights.name='read-employees'
        OR
        app_realm.rights.name='read-divisions'
    )
    OR
    app_realm.roles.name='hr-manager' AND (
        app_realm.rights.name='read-employees'
        OR
        app_realm.rights.name='write-employees'
        OR
        app_realm.rights.name='read-divisions'
    )
    OR
    app_realm.roles.name='admin' AND (
        app_realm.rights.name='read-employees'
        OR
        app_realm.rights.name='write-employees'
        OR
        app_realm.rights.name='read-divisions'
        OR
        app_realm.rights.name='write-divisions'
        OR
        app_realm.rights.name='write-citizenships'
        OR
        app_realm.rights.name='write-titles'
        OR
        app_realm.rights.name='read-users'
        OR
        app_realm.rights.name='write-users'
    );
