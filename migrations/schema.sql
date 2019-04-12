--
-- PostgreSQL database dump
--

-- Dumped from database version 11.2 (Debian 11.2-1.pgdg90+1)
-- Dumped by pg_dump version 11.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: t_email; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.t_email (
    "ID" integer NOT NULL,
    "SEND_TO" character varying(50),
    "CRE_DTE" character varying(40),
    "USR_NME" character varying(255),
    "SUBJ" character varying(255),
    "CTENT" character varying(8000),
    "STS" character varying(10),
    "SENT_DTE" character varying(40)
);


ALTER TABLE public.t_email OWNER TO postgres;

--
-- Name: TABLE t_email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.t_email IS '邮件信息';


--
-- Name: T_EMAIL_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."T_EMAIL_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."T_EMAIL_id_seq" OWNER TO postgres;

--
-- Name: T_EMAIL_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."T_EMAIL_id_seq" OWNED BY public.t_email."ID";


--
-- Name: blogs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.blogs (
    id uuid NOT NULL,
    title character varying(255) NOT NULL,
    content character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.blogs OWNER TO postgres;

--
-- Name: matrices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.matrices (
    id uuid NOT NULL,
    company character varying(255) NOT NULL,
    period character varying(255) NOT NULL,
    matrix character varying(255) NOT NULL,
    code character varying(255) NOT NULL,
    value numeric NOT NULL,
    submit_user character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.matrices OWNER TO postgres;

--
-- Name: notices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notices (
    id uuid NOT NULL,
    send_to character varying(255) NOT NULL,
    user_name character varying(255) NOT NULL,
    subject character varying(255) NOT NULL,
    content character varying(8000) NOT NULL,
    status character varying(255) NOT NULL,
    sent_date character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.notices OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: surveys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.surveys (
    id uuid NOT NULL,
    survey_no character varying(255) NOT NULL,
    submit_user character varying(255) NOT NULL,
    question_no character varying(255) NOT NULL,
    answers character varying(255) NOT NULL,
    submit_date character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.surveys OWNER TO postgres;

--
-- Name: t_email ID; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.t_email ALTER COLUMN "ID" SET DEFAULT nextval('public."T_EMAIL_id_seq"'::regclass);


--
-- Name: t_email T_EMAIL_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.t_email
    ADD CONSTRAINT "T_EMAIL_pkey" PRIMARY KEY ("ID");


--
-- Name: blogs blogs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blogs
    ADD CONSTRAINT blogs_pkey PRIMARY KEY (id);


--
-- Name: matrices matrices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.matrices
    ADD CONSTRAINT matrices_pkey PRIMARY KEY (id);


--
-- Name: notices notices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notices
    ADD CONSTRAINT notices_pkey PRIMARY KEY (id);


--
-- Name: surveys surveys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.surveys
    ADD CONSTRAINT surveys_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

