--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET timezone = UTC;


CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Create the roles, if they don't already exist
--

CREATE OR REPLACE FUNCTION create_role_if_not_exist(u NAME, p TEXT)
RETURNS TEXT AS
$$
BEGIN
  IF NOT EXISTS (
    SELECT pg_roles.rolname AS r
    FROM   pg_catalog.pg_roles
    WHERE  pg_roles.rolname = u
  ) THEN
    EXECUTE format('CREATE ROLE %I WITH CREATEROLE LOGIN PASSWORD %L', u, p);
    RETURN format('CREATE ROLE %I WITH CREATEROLE', u);
  ELSE
    RETURN format('ROLE %I EXISTS', u);
  END IF;
END;
$$
LANGUAGE plpgsql;

-- We set these passwords, but we don't use them. Our setup
-- depends on physical access to the linked app container in
-- order to execute commands. This can therefore be safely used
-- on local test boxes.

SELECT create_role_if_not_exist('ddvote', 'f1ght4urr1ght');
