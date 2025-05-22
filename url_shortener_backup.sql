--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Ubuntu 16.9-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.9 (Ubuntu 16.9-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: url_table; Type: TABLE; Schema: public; Owner: kerem
--

CREATE TABLE public.url_table (
    id integer NOT NULL,
    original_url text NOT NULL,
    short_url character varying(16) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    expires_at timestamp without time zone,
    usage_count integer DEFAULT 0
);


ALTER TABLE public.url_table OWNER TO kerem;

--
-- Name: url_table_id_seq; Type: SEQUENCE; Schema: public; Owner: kerem
--

CREATE SEQUENCE public.url_table_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.url_table_id_seq OWNER TO kerem;

--
-- Name: url_table_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: kerem
--

ALTER SEQUENCE public.url_table_id_seq OWNED BY public.url_table.id;


--
-- Name: url_table id; Type: DEFAULT; Schema: public; Owner: kerem
--

ALTER TABLE ONLY public.url_table ALTER COLUMN id SET DEFAULT nextval('public.url_table_id_seq'::regclass);


--
-- Data for Name: url_table; Type: TABLE DATA; Schema: public; Owner: kerem
--

COPY public.url_table (id, original_url, short_url, created_at, deleted_at, expires_at, usage_count) FROM stdin;
27	https://www.google.com	pmCYRstL	2025-05-21 14:56:32.969065	\N	2025-05-21 15:01:32.969065	0
28	youtube.com	YSisXqcI	2025-05-21 16:02:01.455702	\N	2025-05-21 16:07:01.455702	0
\.


--
-- Name: url_table_id_seq; Type: SEQUENCE SET; Schema: public; Owner: kerem
--

SELECT pg_catalog.setval('public.url_table_id_seq', 28, true);


--
-- Name: url_table url_table_pkey; Type: CONSTRAINT; Schema: public; Owner: kerem
--

ALTER TABLE ONLY public.url_table
    ADD CONSTRAINT url_table_pkey PRIMARY KEY (id);


--
-- Name: url_table url_table_short_url_key; Type: CONSTRAINT; Schema: public; Owner: kerem
--

ALTER TABLE ONLY public.url_table
    ADD CONSTRAINT url_table_short_url_key UNIQUE (short_url);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT ALL ON SCHEMA public TO kerem;


--
-- PostgreSQL database dump complete
--

