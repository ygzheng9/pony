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
-- Name: game_options; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.game_options (
    id uuid NOT NULL,
    game_id uuid NOT NULL,
    seq integer NOT NULL,
    pairs character varying(2000) NOT NULL,
    weights character varying(4000) NOT NULL,
    ratios character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.game_options OWNER TO postgres;

--
-- Name: games; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.games (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    criterion character varying(255) NOT NULL,
    options character varying(255) NOT NULL,
    pairs character varying(2000) NOT NULL,
    weights character varying(4000) NOT NULL,
    ratios character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.games OWNER TO postgres;

--
-- Name: inv_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inv_items (
    id uuid NOT NULL,
    company character varying(255) NOT NULL,
    warehouse character varying(255) NOT NULL,
    year character varying(255) NOT NULL,
    month character varying(255) NOT NULL,
    mat_name character varying(255) NOT NULL,
    mat_code character varying(255) NOT NULL,
    mat_spec character varying(255) NOT NULL,
    mat_style character varying(255) NOT NULL,
    mat_type character varying(255) NOT NULL,
    tree_type character varying(255) NOT NULL,
    cust_code character varying(255) NOT NULL,
    color character varying(255) NOT NULL,
    mat_unit character varying(255) NOT NULL,
    mat_qty numeric NOT NULL,
    mat_amt character varying(255) NOT NULL,
    mat_grade character varying(255) NOT NULL,
    cate1 character varying(255) NOT NULL,
    cate2 character varying(255) NOT NULL,
    surface character varying(255) NOT NULL,
    source character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.inv_items OWNER TO postgres;

--
-- Name: matrices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.matrices (
    id uuid NOT NULL,
    company character varying(255) NOT NULL,
    version character varying(255) NOT NULL,
    period character varying(255) NOT NULL,
    matrix character varying(255) NOT NULL,
    code character varying(255) NOT NULL,
    value character varying(255) NOT NULL,
    submit_user character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.matrices OWNER TO postgres;

--
-- Name: mo_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mo_items (
    id uuid NOT NULL,
    company character varying(255) NOT NULL,
    line character varying(255) NOT NULL,
    mo_date timestamp without time zone NOT NULL,
    mo_num character varying(255) NOT NULL,
    item_num character varying(255) NOT NULL,
    mat_name character varying(1000) NOT NULL,
    item_qty numeric NOT NULL,
    item_unit character varying(255) NOT NULL,
    cate2 character varying(255) NOT NULL,
    cate1 character varying(255) NOT NULL,
    inbound_date timestamp without time zone NOT NULL,
    warehouse character varying(255) NOT NULL,
    work_num character varying(255) NOT NULL,
    mat_code character varying(255) NOT NULL,
    mo_type character varying(255) NOT NULL,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    source character varying(255) NOT NULL,
    shift character varying(255) NOT NULL,
    step character varying(255) NOT NULL,
    mat_spec character varying(255) NOT NULL,
    main_mat_qty numeric NOT NULL,
    input_mat_qty1 numeric NOT NULL,
    input_mat_qty2 numeric NOT NULL,
    claim_qty1 numeric NOT NULL,
    claim_qty2 numeric NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.mo_items OWNER TO postgres;

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
-- Name: po_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.po_items (
    id uuid NOT NULL,
    company character varying(255) NOT NULL,
    po_date timestamp without time zone NOT NULL,
    po_num character varying(255) NOT NULL,
    vendor_name character varying(500) NOT NULL,
    mat_name character varying(500) NOT NULL,
    item_qty numeric NOT NULL,
    item_unit character varying(255) NOT NULL,
    unit_price numeric NOT NULL,
    cate2 character varying(255) NOT NULL,
    cate1 character varying(255) NOT NULL,
    operator character varying(255) NOT NULL,
    inbound_qty numeric NOT NULL,
    outstanding_qty numeric NOT NULL,
    inbound_date timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    line_num character varying(255) DEFAULT ''::character varying NOT NULL,
    dn_num character varying(255) DEFAULT ''::character varying NOT NULL,
    dn_item character varying(255) DEFAULT ''::character varying NOT NULL,
    item_status character varying(255) DEFAULT ''::character varying NOT NULL,
    mat_code character varying(255) DEFAULT ''::character varying NOT NULL,
    mat_spec character varying(255) DEFAULT ''::character varying NOT NULL,
    inbound_unit character varying(255) DEFAULT ''::character varying NOT NULL,
    arrived_book_qty numeric DEFAULT '0'::numeric NOT NULL,
    booked_qty character varying(255) DEFAULT '0'::character varying NOT NULL,
    unbooked_qty character varying(255) DEFAULT '0'::character varying NOT NULL,
    planned_date timestamp without time zone DEFAULT '2006-01-02 00:00:00'::timestamp without time zone NOT NULL,
    delayed_days numeric DEFAULT '0'::numeric NOT NULL
);


ALTER TABLE public.po_items OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: so_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.so_items (
    id uuid NOT NULL,
    company character varying(255) NOT NULL,
    order_num character varying(255) NOT NULL,
    cust_num character varying(255) NOT NULL,
    category character varying(255) NOT NULL,
    serial character varying(255) NOT NULL,
    mat_name character varying(255) NOT NULL,
    mat_model character varying(255) NOT NULL,
    item_qty numeric NOT NULL,
    period character varying(255) NOT NULL,
    move_type character varying(255) NOT NULL,
    sales_type character varying(255) NOT NULL,
    warehouse character varying(255) NOT NULL,
    wh_date timestamp without time zone NOT NULL,
    doc_date timestamp without time zone NOT NULL,
    book_party character varying(255) NOT NULL,
    remark character varying(255) NOT NULL,
    source character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.so_items OWNER TO postgres;

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
-- Name: words; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.words (
    id uuid NOT NULL,
    doc_name character varying(255) NOT NULL,
    word character varying(255) NOT NULL,
    count integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.words OWNER TO postgres;

--
-- Name: game_options game_options_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.game_options
    ADD CONSTRAINT game_options_pkey PRIMARY KEY (id);


--
-- Name: games games_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_pkey PRIMARY KEY (id);


--
-- Name: inv_items inv_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inv_items
    ADD CONSTRAINT inv_items_pkey PRIMARY KEY (id);


--
-- Name: matrices matrices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.matrices
    ADD CONSTRAINT matrices_pkey PRIMARY KEY (id);


--
-- Name: mo_items mo_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mo_items
    ADD CONSTRAINT mo_items_pkey PRIMARY KEY (id);


--
-- Name: notices notices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notices
    ADD CONSTRAINT notices_pkey PRIMARY KEY (id);


--
-- Name: po_items po_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.po_items
    ADD CONSTRAINT po_items_pkey PRIMARY KEY (id);


--
-- Name: so_items so_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.so_items
    ADD CONSTRAINT so_items_pkey PRIMARY KEY (id);


--
-- Name: surveys surveys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.surveys
    ADD CONSTRAINT surveys_pkey PRIMARY KEY (id);


--
-- Name: words words_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.words
    ADD CONSTRAINT words_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

