-- Dropping the usecase scheme's objects.
DROP TABLE IF EXISTS usecase.employees;

DROP TABLE IF EXISTS usecase.citizenships;

DROP TABLE IF EXISTS usecase.titles;

DROP TABLE IF EXISTS usecase.divisions;

DROP SCHEMA IF EXISTS usecase;

-- Dropping the app_realm scheme's objects.
DROP TABLE IF EXISTS app_realm.users;

DROP TABLE IF EXISTS app_realm.role_rights_mapping;

DROP TABLE IF EXISTS app_realm.rights;

DROP TABLE IF EXISTS app_realm.roles;

DROP SCHEMA IF EXISTS app_realm;
