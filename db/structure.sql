--
-- PostgreSQL database dump
--

\restrict c0FbiE23Lpu1cH2iAV84OvDoh9uiE1ttwYi3dgjyIoghGqxEP6M9VlD1Fqq3Rcn

-- Dumped from database version 17.10 (Debian 17.10-1.pgdg13+1)
-- Dumped by pg_dump version 17.10 (Debian 17.10-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: action_text_rich_texts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.action_text_rich_texts (
    id bigint NOT NULL,
    body text,
    created_at timestamp(6) without time zone NOT NULL,
    name character varying NOT NULL,
    record_id bigint NOT NULL,
    record_type character varying NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: action_text_rich_texts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.action_text_rich_texts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: action_text_rich_texts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.action_text_rich_texts_id_seq OWNED BY public.action_text_rich_texts.id;


--
-- Name: active_storage_attachments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.active_storage_attachments (
    id bigint NOT NULL,
    blob_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    name character varying NOT NULL,
    record_id bigint NOT NULL,
    record_type character varying NOT NULL
);


--
-- Name: active_storage_attachments_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.active_storage_attachments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: active_storage_attachments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.active_storage_attachments_id_seq OWNED BY public.active_storage_attachments.id;


--
-- Name: active_storage_blobs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.active_storage_blobs (
    id bigint NOT NULL,
    byte_size bigint NOT NULL,
    checksum character varying,
    content_type character varying,
    created_at timestamp(6) without time zone NOT NULL,
    filename character varying NOT NULL,
    key character varying NOT NULL,
    metadata text,
    service_name character varying NOT NULL
);


--
-- Name: active_storage_blobs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.active_storage_blobs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: active_storage_blobs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.active_storage_blobs_id_seq OWNED BY public.active_storage_blobs.id;


--
-- Name: active_storage_variant_records; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.active_storage_variant_records (
    id bigint NOT NULL,
    blob_id bigint NOT NULL,
    variation_digest character varying NOT NULL
);


--
-- Name: active_storage_variant_records_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.active_storage_variant_records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: active_storage_variant_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.active_storage_variant_records_id_seq OWNED BY public.active_storage_variant_records.id;


--
-- Name: ahoy_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ahoy_events (
    id bigint NOT NULL,
    name character varying,
    properties jsonb,
    "time" timestamp(6) without time zone,
    user_id bigint,
    visit_id bigint
);


--
-- Name: ahoy_events_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ahoy_events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ahoy_events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ahoy_events_id_seq OWNED BY public.ahoy_events.id;


--
-- Name: ahoy_visits; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ahoy_visits (
    id bigint NOT NULL,
    app_version character varying,
    browser character varying,
    city character varying,
    country character varying,
    device_type character varying,
    ip character varying,
    landing_page text,
    latitude double precision,
    longitude double precision,
    os character varying,
    os_version character varying,
    platform character varying,
    referrer text,
    referring_domain character varying,
    region character varying,
    started_at timestamp(6) without time zone,
    user_agent text,
    user_id bigint,
    utm_campaign character varying,
    utm_content character varying,
    utm_medium character varying,
    utm_source character varying,
    utm_term character varying,
    visit_token character varying,
    visitor_token character varying,
    ym_client_id character varying
);


--
-- Name: ahoy_visits_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ahoy_visits_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ahoy_visits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ahoy_visits_id_seq OWNED BY public.ahoy_visits.id;


--
-- Name: ai_chats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ai_chats (
    id bigint NOT NULL,
    ai_model_id bigint,
    created_at timestamp(6) without time zone NOT NULL,
    language_lesson_member_id bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: ai_chats_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ai_chats_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ai_chats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ai_chats_id_seq OWNED BY public.ai_chats.id;


--
-- Name: ai_messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ai_messages (
    id bigint NOT NULL,
    ai_chat_id bigint NOT NULL,
    ai_model_id bigint,
    ai_tool_call_id bigint,
    cache_creation_tokens integer,
    cached_tokens integer,
    content text,
    content_raw json,
    created_at timestamp(6) without time zone NOT NULL,
    input_tokens integer,
    output_tokens integer,
    role character varying NOT NULL,
    thinking_signature text,
    thinking_text text,
    thinking_tokens integer,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint
);


--
-- Name: ai_messages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ai_messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ai_messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ai_messages_id_seq OWNED BY public.ai_messages.id;


--
-- Name: ai_models; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ai_models (
    id bigint NOT NULL,
    capabilities jsonb DEFAULT '[]'::jsonb,
    context_window integer,
    created_at timestamp(6) without time zone NOT NULL,
    family character varying,
    knowledge_cutoff date,
    max_output_tokens integer,
    metadata jsonb DEFAULT '{}'::jsonb,
    modalities jsonb DEFAULT '{}'::jsonb,
    model_created_at timestamp(6) without time zone,
    model_id character varying NOT NULL,
    name character varying NOT NULL,
    pricing jsonb DEFAULT '{}'::jsonb,
    provider character varying NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: ai_models_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ai_models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ai_models_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ai_models_id_seq OWNED BY public.ai_models.id;


--
-- Name: ai_tool_calls; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ai_tool_calls (
    id bigint NOT NULL,
    ai_message_id bigint NOT NULL,
    arguments jsonb DEFAULT '{}'::jsonb,
    created_at timestamp(6) without time zone NOT NULL,
    name character varying NOT NULL,
    thought_signature text,
    tool_call_id character varying NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: ai_tool_calls_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ai_tool_calls_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ai_tool_calls_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ai_tool_calls_id_seq OWNED BY public.ai_tool_calls.id;


--
-- Name: ar_internal_metadata; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ar_internal_metadata (
    key character varying NOT NULL,
    value character varying,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: banners; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.banners (
    id bigint NOT NULL,
    background character varying DEFAULT 'cta_gradient'::character varying NOT NULL,
    body text NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    finishes_at timestamp(6) without time zone,
    locale character varying NOT NULL,
    starts_at timestamp(6) without time zone,
    state character varying DEFAULT 'draft'::character varying NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    url character varying
);


--
-- Name: banners_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.banners_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: banners_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.banners_id_seq OWNED BY public.banners.id;


--
-- Name: blog_post_likes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blog_post_likes (
    id bigint NOT NULL,
    blog_post_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint
);


--
-- Name: blog_post_likes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blog_post_likes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blog_post_likes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blog_post_likes_id_seq OWNED BY public.blog_post_likes.id;


--
-- Name: blog_post_related_language_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blog_post_related_language_items (
    id bigint NOT NULL,
    blog_post_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    language_id bigint NOT NULL,
    "order" integer,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: blog_post_related_language_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blog_post_related_language_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blog_post_related_language_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blog_post_related_language_items_id_seq OWNED BY public.blog_post_related_language_items.id;


--
-- Name: blog_posts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blog_posts (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    creator_id bigint NOT NULL,
    description character varying,
    language_id bigint,
    locale character varying,
    name character varying,
    related_language_items_count integer DEFAULT 0 NOT NULL,
    slug character varying,
    state character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: blog_posts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blog_posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blog_posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blog_posts_id_seq OWNED BY public.blog_posts.id;


--
-- Name: book_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.book_requests (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    state character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: book_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.book_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: book_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.book_requests_id_seq OWNED BY public.book_requests.id;


--
-- Name: course_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.course_categories (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    name character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: course_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.course_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: course_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.course_categories_id_seq OWNED BY public.course_categories.id;


--
-- Name: event_store_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.event_store_events (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    data bytea NOT NULL,
    event_id character varying(36) NOT NULL,
    event_type character varying NOT NULL,
    metadata bytea,
    valid_at timestamp(6) without time zone
);


--
-- Name: event_store_events_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.event_store_events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: event_store_events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.event_store_events_id_seq OWNED BY public.event_store_events.id;


--
-- Name: event_store_events_in_streams; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.event_store_events_in_streams (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    event_id character varying(36) NOT NULL,
    "position" integer,
    stream character varying NOT NULL
);


--
-- Name: event_store_events_in_streams_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.event_store_events_in_streams_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: event_store_events_in_streams_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.event_store_events_in_streams_id_seq OWNED BY public.event_store_events_in_streams.id;


--
-- Name: flipper_features; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.flipper_features (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    key character varying NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: flipper_features_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.flipper_features_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flipper_features_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.flipper_features_id_seq OWNED BY public.flipper_features.id;


--
-- Name: flipper_gates; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.flipper_gates (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    feature_key character varying NOT NULL,
    key character varying NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    value text
);


--
-- Name: flipper_gates_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.flipper_gates_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flipper_gates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.flipper_gates_id_seq OWNED BY public.flipper_gates.id;


--
-- Name: language_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_categories (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    description character varying,
    header character varying,
    locale character varying,
    name character varying,
    name_en character varying,
    name_ru character varying,
    slug character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_categories_id_seq OWNED BY public.language_categories.id;


--
-- Name: language_category_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_category_items (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    language_category_id bigint NOT NULL,
    language_landing_page_id bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_category_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_category_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_category_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_category_items_id_seq OWNED BY public.language_category_items.id;


--
-- Name: language_category_qna_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_category_qna_items (
    id bigint NOT NULL,
    answer character varying,
    created_at timestamp(6) without time zone NOT NULL,
    language_category_id bigint NOT NULL,
    question character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_category_qna_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_category_qna_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_category_qna_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_category_qna_items_id_seq OWNED BY public.language_category_qna_items.id;


--
-- Name: language_landing_page_qna_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_landing_page_qna_items (
    id bigint NOT NULL,
    answer character varying,
    created_at timestamp(6) without time zone NOT NULL,
    language_landing_page_id bigint NOT NULL,
    question character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_landing_page_qna_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_landing_page_qna_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_landing_page_qna_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_landing_page_qna_items_id_seq OWNED BY public.language_landing_page_qna_items.id;


--
-- Name: language_landing_pages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_landing_pages (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    description character varying,
    footer boolean,
    footer_name character varying,
    header character varying,
    landing_page_to_redirect_id bigint,
    language_category_id bigint,
    language_id bigint NOT NULL,
    listed boolean,
    locale character varying,
    main boolean,
    meta_description character varying,
    meta_title character varying,
    name character varying,
    "order" character varying,
    outcomes_description character varying,
    outcomes_header character varying,
    slug character varying,
    state character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    used_in_description character varying,
    used_in_header character varying
);


--
-- Name: language_landing_pages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_landing_pages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_landing_pages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_landing_pages_id_seq OWNED BY public.language_landing_pages.id;


--
-- Name: language_lesson_members; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_lesson_members (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    language_id bigint NOT NULL,
    language_member_id bigint NOT NULL,
    lesson_id bigint NOT NULL,
    messages_count integer DEFAULT 0,
    state character varying,
    updated_at timestamp without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: language_lesson_members_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_lesson_members_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_lesson_members_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_lesson_members_id_seq OWNED BY public.language_lesson_members.id;


--
-- Name: language_lesson_reviews; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_lesson_reviews (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    language_id bigint NOT NULL,
    language_lesson_id bigint NOT NULL,
    language_lesson_version_id bigint NOT NULL,
    language_lesson_version_info_id bigint NOT NULL,
    locale character varying NOT NULL,
    summary text DEFAULT ''::text NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_lesson_reviews_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_lesson_reviews_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_lesson_reviews_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_lesson_reviews_id_seq OWNED BY public.language_lesson_reviews.id;


--
-- Name: language_lesson_version_infos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_lesson_version_infos (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    definitions character varying,
    description character varying,
    instructions character varying,
    language_id bigint NOT NULL,
    language_lesson_id bigint NOT NULL,
    language_version_id bigint NOT NULL,
    locale character varying,
    name character varying,
    theory character varying,
    tips character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    version_id bigint NOT NULL
);


--
-- Name: language_lesson_version_infos_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_lesson_version_infos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_lesson_version_infos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_lesson_version_infos_id_seq OWNED BY public.language_lesson_version_infos.id;


--
-- Name: language_lesson_versions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_lesson_versions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    language_id bigint NOT NULL,
    language_version_id bigint NOT NULL,
    lesson_id bigint NOT NULL,
    module_version_id bigint NOT NULL,
    natural_order integer,
    "order" integer,
    original_code character varying,
    path_to_code character varying,
    prepared_code character varying,
    test_code character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_lesson_versions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_lesson_versions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_lesson_versions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_lesson_versions_id_seq OWNED BY public.language_lesson_versions.id;


--
-- Name: language_lessons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_lessons (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    language_id bigint,
    module_id bigint,
    natural_order integer,
    "order" integer,
    original_code text,
    path_to_code character varying(255),
    prepared_code text,
    slug character varying(255),
    state character varying(255),
    test_code text,
    updated_at timestamp without time zone NOT NULL,
    upload_id bigint
);


--
-- Name: language_lessons_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_lessons_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_lessons_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_lessons_id_seq OWNED BY public.language_lessons.id;


--
-- Name: language_members; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_members (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    finished_lessons_count integer DEFAULT 0 NOT NULL,
    language_id bigint NOT NULL,
    state character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: language_members_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_members_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_members_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_members_id_seq OWNED BY public.language_members.id;


--
-- Name: language_module_descriptions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_module_descriptions (
    id bigint NOT NULL,
    description text,
    inserted_at timestamp without time zone NOT NULL,
    language_id bigint,
    locale character varying(255),
    module_id bigint,
    name character varying(255),
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: language_module_descriptions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_module_descriptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_module_descriptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_module_descriptions_id_seq OWNED BY public.language_module_descriptions.id;


--
-- Name: language_module_version_infos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_module_version_infos (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    description character varying,
    language_id bigint NOT NULL,
    language_version_id bigint NOT NULL,
    locale character varying,
    name character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    version_id bigint NOT NULL
);


--
-- Name: language_module_version_infos_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_module_version_infos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_module_version_infos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_module_version_infos_id_seq OWNED BY public.language_module_version_infos.id;


--
-- Name: language_module_versions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_module_versions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    language_id bigint NOT NULL,
    language_version_id bigint NOT NULL,
    module_id bigint NOT NULL,
    "order" integer,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_module_versions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_module_versions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_module_versions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_module_versions_id_seq OWNED BY public.language_module_versions.id;


--
-- Name: language_modules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_modules (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    language_id bigint,
    "order" integer,
    slug character varying(255),
    state character varying(255),
    updated_at timestamp without time zone NOT NULL,
    upload_id bigint
);


--
-- Name: language_modules_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_modules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_modules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_modules_id_seq OWNED BY public.language_modules.id;


--
-- Name: language_version_infos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_version_infos (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    description character varying,
    header character varying,
    keywords character varying,
    language_id bigint NOT NULL,
    language_version_id bigint NOT NULL,
    locale character varying,
    seo_description text,
    title character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_version_infos_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_version_infos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_version_infos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_version_infos_id_seq OWNED BY public.language_version_infos.id;


--
-- Name: language_versions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.language_versions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    docker_image character varying,
    exercise_filename character varying,
    exercise_test_filename character varying,
    extension character varying,
    language_id bigint NOT NULL,
    learn_as character varying,
    lessons_count integer DEFAULT 0 NOT NULL,
    name character varying,
    progress character varying,
    result character varying,
    state character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: language_versions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.language_versions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_versions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.language_versions_id_seq OWNED BY public.language_versions.id;


--
-- Name: languages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.languages (
    id bigint NOT NULL,
    category_id bigint,
    created_at timestamp(6) without time zone NOT NULL,
    current_version_id bigint,
    docker_image character varying(255),
    exercise_filename character varying(255),
    exercise_test_filename character varying(255),
    extension character varying(255),
    hexlet_program_landing_page character varying,
    learn_as character varying,
    lessons_count integer DEFAULT 0 NOT NULL,
    members_count integer DEFAULT 0 NOT NULL,
    name character varying(255),
    "order" integer,
    progress character varying,
    slug character varying(255),
    state character varying(255),
    updated_at timestamp without time zone NOT NULL,
    upload_id bigint
);


--
-- Name: languages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.languages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: languages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.languages_id_seq OWNED BY public.languages.id;


--
-- Name: leads; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.leads (
    id bigint NOT NULL,
    ahoy_visit_id bigint,
    courses_data text,
    created_at timestamp(6) without time zone NOT NULL,
    email character varying,
    phone character varying,
    state character varying,
    survey_answers_data text,
    telegram character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL,
    whatsapp character varying,
    ym_client_id character varying
);


--
-- Name: leads_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.leads_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: leads_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.leads_id_seq OWNED BY public.leads.id;


--
-- Name: reviews; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reviews (
    id bigint NOT NULL,
    body text,
    created_at timestamp(6) without time zone NOT NULL,
    first_name character varying,
    language_id bigint NOT NULL,
    last_name character varying,
    locale character varying,
    pinned boolean,
    state character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: reviews_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.reviews_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: reviews_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.reviews_id_seq OWNED BY public.reviews.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    ip_address character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    user_agent character varying,
    user_id bigint NOT NULL
);


--
-- Name: sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.sessions_id_seq OWNED BY public.sessions.id;


--
-- Name: solid_cable_messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_cable_messages (
    id bigint NOT NULL,
    channel bytea NOT NULL,
    channel_hash bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    payload bytea NOT NULL
);


--
-- Name: solid_cable_messages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_cable_messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_cable_messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_cable_messages_id_seq OWNED BY public.solid_cable_messages.id;


--
-- Name: solid_cache_entries; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_cache_entries (
    id bigint NOT NULL,
    byte_size integer NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    key bytea NOT NULL,
    key_hash bigint NOT NULL,
    value bytea NOT NULL
);


--
-- Name: solid_cache_entries_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_cache_entries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_cache_entries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_cache_entries_id_seq OWNED BY public.solid_cache_entries.id;


--
-- Name: solid_queue_blocked_executions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_blocked_executions (
    id bigint NOT NULL,
    concurrency_key character varying NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    expires_at timestamp(6) without time zone NOT NULL,
    job_id bigint NOT NULL,
    priority integer DEFAULT 0 NOT NULL,
    queue_name character varying NOT NULL
);


--
-- Name: solid_queue_blocked_executions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_blocked_executions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_blocked_executions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_blocked_executions_id_seq OWNED BY public.solid_queue_blocked_executions.id;


--
-- Name: solid_queue_claimed_executions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_claimed_executions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    job_id bigint NOT NULL,
    process_id bigint
);


--
-- Name: solid_queue_claimed_executions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_claimed_executions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_claimed_executions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_claimed_executions_id_seq OWNED BY public.solid_queue_claimed_executions.id;


--
-- Name: solid_queue_failed_executions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_failed_executions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    error text,
    job_id bigint NOT NULL
);


--
-- Name: solid_queue_failed_executions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_failed_executions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_failed_executions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_failed_executions_id_seq OWNED BY public.solid_queue_failed_executions.id;


--
-- Name: solid_queue_jobs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_jobs (
    id bigint NOT NULL,
    active_job_id character varying,
    arguments text,
    class_name character varying NOT NULL,
    concurrency_key character varying,
    created_at timestamp(6) without time zone NOT NULL,
    finished_at timestamp(6) without time zone,
    priority integer DEFAULT 0 NOT NULL,
    queue_name character varying NOT NULL,
    scheduled_at timestamp(6) without time zone,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: solid_queue_jobs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_jobs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_jobs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_jobs_id_seq OWNED BY public.solid_queue_jobs.id;


--
-- Name: solid_queue_pauses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_pauses (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    queue_name character varying NOT NULL
);


--
-- Name: solid_queue_pauses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_pauses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_pauses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_pauses_id_seq OWNED BY public.solid_queue_pauses.id;


--
-- Name: solid_queue_processes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_processes (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    hostname character varying,
    kind character varying NOT NULL,
    last_heartbeat_at timestamp(6) without time zone NOT NULL,
    metadata text,
    name character varying NOT NULL,
    pid integer NOT NULL,
    supervisor_id bigint
);


--
-- Name: solid_queue_processes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_processes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_processes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_processes_id_seq OWNED BY public.solid_queue_processes.id;


--
-- Name: solid_queue_ready_executions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_ready_executions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    job_id bigint NOT NULL,
    priority integer DEFAULT 0 NOT NULL,
    queue_name character varying NOT NULL
);


--
-- Name: solid_queue_ready_executions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_ready_executions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_ready_executions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_ready_executions_id_seq OWNED BY public.solid_queue_ready_executions.id;


--
-- Name: solid_queue_recurring_executions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_recurring_executions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    job_id bigint NOT NULL,
    run_at timestamp(6) without time zone NOT NULL,
    task_key character varying NOT NULL
);


--
-- Name: solid_queue_recurring_executions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_recurring_executions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_recurring_executions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_recurring_executions_id_seq OWNED BY public.solid_queue_recurring_executions.id;


--
-- Name: solid_queue_recurring_tasks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_recurring_tasks (
    id bigint NOT NULL,
    arguments text,
    class_name character varying,
    command character varying(2048),
    created_at timestamp(6) without time zone NOT NULL,
    description text,
    key character varying NOT NULL,
    priority integer DEFAULT 0,
    queue_name character varying,
    schedule character varying NOT NULL,
    static boolean DEFAULT true NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: solid_queue_recurring_tasks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_recurring_tasks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_recurring_tasks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_recurring_tasks_id_seq OWNED BY public.solid_queue_recurring_tasks.id;


--
-- Name: solid_queue_scheduled_executions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_scheduled_executions (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    job_id bigint NOT NULL,
    priority integer DEFAULT 0 NOT NULL,
    queue_name character varying NOT NULL,
    scheduled_at timestamp(6) without time zone NOT NULL
);


--
-- Name: solid_queue_scheduled_executions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_scheduled_executions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_scheduled_executions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_scheduled_executions_id_seq OWNED BY public.solid_queue_scheduled_executions.id;


--
-- Name: solid_queue_semaphores; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.solid_queue_semaphores (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    expires_at timestamp(6) without time zone NOT NULL,
    key character varying NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    value integer DEFAULT 1 NOT NULL
);


--
-- Name: solid_queue_semaphores_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.solid_queue_semaphores_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: solid_queue_semaphores_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.solid_queue_semaphores_id_seq OWNED BY public.solid_queue_semaphores.id;


--
-- Name: staff_member_role_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.staff_member_role_permissions (
    id bigint NOT NULL,
    can_create boolean DEFAULT false NOT NULL,
    can_destroy boolean DEFAULT false NOT NULL,
    can_index boolean DEFAULT false NOT NULL,
    can_update boolean DEFAULT false NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    resource character varying NOT NULL,
    role_id bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: staff_member_role_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.staff_member_role_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: staff_member_role_permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.staff_member_role_permissions_id_seq OWNED BY public.staff_member_role_permissions.id;


--
-- Name: staff_member_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.staff_member_roles (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    description text,
    name character varying NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: staff_member_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.staff_member_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: staff_member_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.staff_member_roles_id_seq OWNED BY public.staff_member_roles.id;


--
-- Name: staff_members; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.staff_members (
    id bigint NOT NULL,
    allowed_locales character varying[] DEFAULT '{ru}'::character varying[] NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    role_id bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: staff_members_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.staff_members_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: staff_members_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.staff_members_id_seq OWNED BY public.staff_members.id;


--
-- Name: survey_answers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.survey_answers (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    state character varying,
    survey_id bigint NOT NULL,
    survey_item_id bigint,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: survey_answers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.survey_answers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: survey_answers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.survey_answers_id_seq OWNED BY public.survey_answers.id;


--
-- Name: survey_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.survey_items (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    "order" integer NOT NULL,
    slug character varying,
    state character varying,
    survey_id bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    value character varying
);


--
-- Name: survey_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.survey_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: survey_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.survey_items_id_seq OWNED BY public.survey_items.id;


--
-- Name: survey_scenario_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.survey_scenario_items (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    "order" integer,
    scenario_id bigint NOT NULL,
    survey_id bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: survey_scenario_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.survey_scenario_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: survey_scenario_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.survey_scenario_items_id_seq OWNED BY public.survey_scenario_items.id;


--
-- Name: survey_scenario_members; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.survey_scenario_members (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    event_name character varying,
    scenario_id bigint NOT NULL,
    state character varying,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: survey_scenario_members_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.survey_scenario_members_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: survey_scenario_members_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.survey_scenario_members_id_seq OWNED BY public.survey_scenario_members.id;


--
-- Name: survey_scenario_triggers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.survey_scenario_triggers (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    event_name character varying,
    event_threshold_count integer,
    scenario_id bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: survey_scenario_triggers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.survey_scenario_triggers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: survey_scenario_triggers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.survey_scenario_triggers_id_seq OWNED BY public.survey_scenario_triggers.id;


--
-- Name: survey_scenarios; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.survey_scenarios (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    locale character varying,
    name character varying,
    state character varying,
    survey_item_id bigint,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: survey_scenarios_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.survey_scenarios_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: survey_scenarios_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.survey_scenarios_id_seq OWNED BY public.survey_scenarios.id;


--
-- Name: surveys; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.surveys (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    description character varying,
    locale character varying,
    parent_survey_id bigint,
    parent_survey_item_id bigint,
    question character varying,
    run_after_finishing_lessons_count integer DEFAULT 0,
    run_always boolean DEFAULT false,
    slug character varying,
    state character varying,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: surveys_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.surveys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: surveys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.surveys_id_seq OWNED BY public.surveys.id;


--
-- Name: taggings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.taggings (
    id bigint NOT NULL,
    context character varying(128),
    created_at timestamp without time zone,
    tag_id bigint,
    taggable_id bigint,
    taggable_type character varying,
    tagger_id bigint,
    tagger_type character varying,
    tenant character varying(128)
);


--
-- Name: taggings_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.taggings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: taggings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.taggings_id_seq OWNED BY public.taggings.id;


--
-- Name: tags; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tags (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    name character varying,
    taggings_count integer DEFAULT 0,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: tags_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tags_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tags_id_seq OWNED BY public.tags.id;


--
-- Name: uploads; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.uploads (
    id bigint NOT NULL,
    inserted_at timestamp without time zone NOT NULL,
    language_name character varying(255),
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: uploads_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.uploads_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: uploads_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.uploads_id_seq OWNED BY public.uploads.id;


--
-- Name: user_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_accounts (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    provider character varying(255) NOT NULL,
    uid character varying(255) NOT NULL,
    updated_at timestamp(0) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: user_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_accounts_id_seq OWNED BY public.user_accounts.id;


--
-- Name: user_credentials; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_credentials (
    id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    external_id character varying NOT NULL,
    nickname character varying,
    public_key character varying NOT NULL,
    sign_count bigint DEFAULT 0 NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: user_credentials_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_credentials_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_credentials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_credentials_id_seq OWNED BY public.user_credentials.id;


--
-- Name: user_survey_pivots; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_survey_pivots (
    id bigint NOT NULL,
    coding_experience_item_id bigint,
    created_at timestamp(6) without time zone NOT NULL,
    goal_item_id bigint,
    study_plan_item_id bigint,
    updated_at timestamp(6) without time zone NOT NULL,
    user_id bigint NOT NULL
);


--
-- Name: user_survey_pivots_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_survey_pivots_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_survey_pivots_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_survey_pivots_id_seq OWNED BY public.user_survey_pivots.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    admin boolean,
    assistant_messages_count integer DEFAULT 0,
    confirmation_token character varying(255),
    contact_method character varying,
    contact_value character varying,
    created_at timestamp(6) without time zone NOT NULL,
    email character varying(255),
    email_delivery_state character varying(255),
    facebook_uid character varying(255),
    first_name character varying(255),
    github_uid integer,
    help boolean,
    last_name character varying(255),
    locale character varying(255),
    nickname character varying(255),
    password_digest character varying(255),
    phone character varying,
    phone_verified_at timestamp(6) without time zone,
    state character varying(255),
    updated_at timestamp without time zone NOT NULL,
    webauthn_id character varying
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: action_text_rich_texts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.action_text_rich_texts ALTER COLUMN id SET DEFAULT nextval('public.action_text_rich_texts_id_seq'::regclass);


--
-- Name: active_storage_attachments id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_storage_attachments ALTER COLUMN id SET DEFAULT nextval('public.active_storage_attachments_id_seq'::regclass);


--
-- Name: active_storage_blobs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_storage_blobs ALTER COLUMN id SET DEFAULT nextval('public.active_storage_blobs_id_seq'::regclass);


--
-- Name: active_storage_variant_records id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_storage_variant_records ALTER COLUMN id SET DEFAULT nextval('public.active_storage_variant_records_id_seq'::regclass);


--
-- Name: ahoy_events id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ahoy_events ALTER COLUMN id SET DEFAULT nextval('public.ahoy_events_id_seq'::regclass);


--
-- Name: ahoy_visits id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ahoy_visits ALTER COLUMN id SET DEFAULT nextval('public.ahoy_visits_id_seq'::regclass);


--
-- Name: ai_chats id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_chats ALTER COLUMN id SET DEFAULT nextval('public.ai_chats_id_seq'::regclass);


--
-- Name: ai_messages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_messages ALTER COLUMN id SET DEFAULT nextval('public.ai_messages_id_seq'::regclass);


--
-- Name: ai_models id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_models ALTER COLUMN id SET DEFAULT nextval('public.ai_models_id_seq'::regclass);


--
-- Name: ai_tool_calls id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_tool_calls ALTER COLUMN id SET DEFAULT nextval('public.ai_tool_calls_id_seq'::regclass);


--
-- Name: banners id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.banners ALTER COLUMN id SET DEFAULT nextval('public.banners_id_seq'::regclass);


--
-- Name: blog_post_likes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_post_likes ALTER COLUMN id SET DEFAULT nextval('public.blog_post_likes_id_seq'::regclass);


--
-- Name: blog_post_related_language_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_post_related_language_items ALTER COLUMN id SET DEFAULT nextval('public.blog_post_related_language_items_id_seq'::regclass);


--
-- Name: blog_posts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_posts ALTER COLUMN id SET DEFAULT nextval('public.blog_posts_id_seq'::regclass);


--
-- Name: book_requests id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.book_requests ALTER COLUMN id SET DEFAULT nextval('public.book_requests_id_seq'::regclass);


--
-- Name: course_categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_categories ALTER COLUMN id SET DEFAULT nextval('public.course_categories_id_seq'::regclass);


--
-- Name: event_store_events id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_store_events ALTER COLUMN id SET DEFAULT nextval('public.event_store_events_id_seq'::regclass);


--
-- Name: event_store_events_in_streams id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_store_events_in_streams ALTER COLUMN id SET DEFAULT nextval('public.event_store_events_in_streams_id_seq'::regclass);


--
-- Name: flipper_features id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.flipper_features ALTER COLUMN id SET DEFAULT nextval('public.flipper_features_id_seq'::regclass);


--
-- Name: flipper_gates id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.flipper_gates ALTER COLUMN id SET DEFAULT nextval('public.flipper_gates_id_seq'::regclass);


--
-- Name: language_categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_categories ALTER COLUMN id SET DEFAULT nextval('public.language_categories_id_seq'::regclass);


--
-- Name: language_category_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_category_items ALTER COLUMN id SET DEFAULT nextval('public.language_category_items_id_seq'::regclass);


--
-- Name: language_category_qna_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_category_qna_items ALTER COLUMN id SET DEFAULT nextval('public.language_category_qna_items_id_seq'::regclass);


--
-- Name: language_landing_page_qna_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_landing_page_qna_items ALTER COLUMN id SET DEFAULT nextval('public.language_landing_page_qna_items_id_seq'::regclass);


--
-- Name: language_landing_pages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_landing_pages ALTER COLUMN id SET DEFAULT nextval('public.language_landing_pages_id_seq'::regclass);


--
-- Name: language_lesson_members id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_members ALTER COLUMN id SET DEFAULT nextval('public.language_lesson_members_id_seq'::regclass);


--
-- Name: language_lesson_reviews id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_reviews ALTER COLUMN id SET DEFAULT nextval('public.language_lesson_reviews_id_seq'::regclass);


--
-- Name: language_lesson_version_infos id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_version_infos ALTER COLUMN id SET DEFAULT nextval('public.language_lesson_version_infos_id_seq'::regclass);


--
-- Name: language_lesson_versions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_versions ALTER COLUMN id SET DEFAULT nextval('public.language_lesson_versions_id_seq'::regclass);


--
-- Name: language_lessons id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lessons ALTER COLUMN id SET DEFAULT nextval('public.language_lessons_id_seq'::regclass);


--
-- Name: language_members id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_members ALTER COLUMN id SET DEFAULT nextval('public.language_members_id_seq'::regclass);


--
-- Name: language_module_descriptions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_descriptions ALTER COLUMN id SET DEFAULT nextval('public.language_module_descriptions_id_seq'::regclass);


--
-- Name: language_module_version_infos id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_version_infos ALTER COLUMN id SET DEFAULT nextval('public.language_module_version_infos_id_seq'::regclass);


--
-- Name: language_module_versions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_versions ALTER COLUMN id SET DEFAULT nextval('public.language_module_versions_id_seq'::regclass);


--
-- Name: language_modules id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_modules ALTER COLUMN id SET DEFAULT nextval('public.language_modules_id_seq'::regclass);


--
-- Name: language_version_infos id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_version_infos ALTER COLUMN id SET DEFAULT nextval('public.language_version_infos_id_seq'::regclass);


--
-- Name: language_versions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_versions ALTER COLUMN id SET DEFAULT nextval('public.language_versions_id_seq'::regclass);


--
-- Name: languages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.languages ALTER COLUMN id SET DEFAULT nextval('public.languages_id_seq'::regclass);


--
-- Name: leads id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.leads ALTER COLUMN id SET DEFAULT nextval('public.leads_id_seq'::regclass);


--
-- Name: reviews id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reviews ALTER COLUMN id SET DEFAULT nextval('public.reviews_id_seq'::regclass);


--
-- Name: sessions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions ALTER COLUMN id SET DEFAULT nextval('public.sessions_id_seq'::regclass);


--
-- Name: solid_cable_messages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_cable_messages ALTER COLUMN id SET DEFAULT nextval('public.solid_cable_messages_id_seq'::regclass);


--
-- Name: solid_cache_entries id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_cache_entries ALTER COLUMN id SET DEFAULT nextval('public.solid_cache_entries_id_seq'::regclass);


--
-- Name: solid_queue_blocked_executions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_blocked_executions ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_blocked_executions_id_seq'::regclass);


--
-- Name: solid_queue_claimed_executions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_claimed_executions ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_claimed_executions_id_seq'::regclass);


--
-- Name: solid_queue_failed_executions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_failed_executions ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_failed_executions_id_seq'::regclass);


--
-- Name: solid_queue_jobs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_jobs ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_jobs_id_seq'::regclass);


--
-- Name: solid_queue_pauses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_pauses ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_pauses_id_seq'::regclass);


--
-- Name: solid_queue_processes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_processes ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_processes_id_seq'::regclass);


--
-- Name: solid_queue_ready_executions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_ready_executions ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_ready_executions_id_seq'::regclass);


--
-- Name: solid_queue_recurring_executions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_recurring_executions ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_recurring_executions_id_seq'::regclass);


--
-- Name: solid_queue_recurring_tasks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_recurring_tasks ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_recurring_tasks_id_seq'::regclass);


--
-- Name: solid_queue_scheduled_executions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_scheduled_executions ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_scheduled_executions_id_seq'::regclass);


--
-- Name: solid_queue_semaphores id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_semaphores ALTER COLUMN id SET DEFAULT nextval('public.solid_queue_semaphores_id_seq'::regclass);


--
-- Name: staff_member_role_permissions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_member_role_permissions ALTER COLUMN id SET DEFAULT nextval('public.staff_member_role_permissions_id_seq'::regclass);


--
-- Name: staff_member_roles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_member_roles ALTER COLUMN id SET DEFAULT nextval('public.staff_member_roles_id_seq'::regclass);


--
-- Name: staff_members id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_members ALTER COLUMN id SET DEFAULT nextval('public.staff_members_id_seq'::regclass);


--
-- Name: survey_answers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_answers ALTER COLUMN id SET DEFAULT nextval('public.survey_answers_id_seq'::regclass);


--
-- Name: survey_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_items ALTER COLUMN id SET DEFAULT nextval('public.survey_items_id_seq'::regclass);


--
-- Name: survey_scenario_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_items ALTER COLUMN id SET DEFAULT nextval('public.survey_scenario_items_id_seq'::regclass);


--
-- Name: survey_scenario_members id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_members ALTER COLUMN id SET DEFAULT nextval('public.survey_scenario_members_id_seq'::regclass);


--
-- Name: survey_scenario_triggers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_triggers ALTER COLUMN id SET DEFAULT nextval('public.survey_scenario_triggers_id_seq'::regclass);


--
-- Name: survey_scenarios id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenarios ALTER COLUMN id SET DEFAULT nextval('public.survey_scenarios_id_seq'::regclass);


--
-- Name: surveys id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.surveys ALTER COLUMN id SET DEFAULT nextval('public.surveys_id_seq'::regclass);


--
-- Name: taggings id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.taggings ALTER COLUMN id SET DEFAULT nextval('public.taggings_id_seq'::regclass);


--
-- Name: tags id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags ALTER COLUMN id SET DEFAULT nextval('public.tags_id_seq'::regclass);


--
-- Name: uploads id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uploads ALTER COLUMN id SET DEFAULT nextval('public.uploads_id_seq'::regclass);


--
-- Name: user_accounts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_accounts ALTER COLUMN id SET DEFAULT nextval('public.user_accounts_id_seq'::regclass);


--
-- Name: user_credentials id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_credentials ALTER COLUMN id SET DEFAULT nextval('public.user_credentials_id_seq'::regclass);


--
-- Name: user_survey_pivots id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_survey_pivots ALTER COLUMN id SET DEFAULT nextval('public.user_survey_pivots_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: action_text_rich_texts action_text_rich_texts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.action_text_rich_texts
    ADD CONSTRAINT action_text_rich_texts_pkey PRIMARY KEY (id);


--
-- Name: active_storage_attachments active_storage_attachments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_storage_attachments
    ADD CONSTRAINT active_storage_attachments_pkey PRIMARY KEY (id);


--
-- Name: active_storage_blobs active_storage_blobs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_storage_blobs
    ADD CONSTRAINT active_storage_blobs_pkey PRIMARY KEY (id);


--
-- Name: active_storage_variant_records active_storage_variant_records_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_storage_variant_records
    ADD CONSTRAINT active_storage_variant_records_pkey PRIMARY KEY (id);


--
-- Name: ahoy_events ahoy_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ahoy_events
    ADD CONSTRAINT ahoy_events_pkey PRIMARY KEY (id);


--
-- Name: ahoy_visits ahoy_visits_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ahoy_visits
    ADD CONSTRAINT ahoy_visits_pkey PRIMARY KEY (id);


--
-- Name: ai_chats ai_chats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_chats
    ADD CONSTRAINT ai_chats_pkey PRIMARY KEY (id);


--
-- Name: ai_messages ai_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_messages
    ADD CONSTRAINT ai_messages_pkey PRIMARY KEY (id);


--
-- Name: ai_models ai_models_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_models
    ADD CONSTRAINT ai_models_pkey PRIMARY KEY (id);


--
-- Name: ai_tool_calls ai_tool_calls_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_tool_calls
    ADD CONSTRAINT ai_tool_calls_pkey PRIMARY KEY (id);


--
-- Name: ar_internal_metadata ar_internal_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ar_internal_metadata
    ADD CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key);


--
-- Name: banners banners_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.banners
    ADD CONSTRAINT banners_pkey PRIMARY KEY (id);


--
-- Name: blog_post_likes blog_post_likes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_post_likes
    ADD CONSTRAINT blog_post_likes_pkey PRIMARY KEY (id);


--
-- Name: blog_post_related_language_items blog_post_related_language_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_post_related_language_items
    ADD CONSTRAINT blog_post_related_language_items_pkey PRIMARY KEY (id);


--
-- Name: blog_posts blog_posts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_posts
    ADD CONSTRAINT blog_posts_pkey PRIMARY KEY (id);


--
-- Name: book_requests book_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.book_requests
    ADD CONSTRAINT book_requests_pkey PRIMARY KEY (id);


--
-- Name: course_categories course_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_categories
    ADD CONSTRAINT course_categories_pkey PRIMARY KEY (id);


--
-- Name: event_store_events_in_streams event_store_events_in_streams_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_store_events_in_streams
    ADD CONSTRAINT event_store_events_in_streams_pkey PRIMARY KEY (id);


--
-- Name: event_store_events event_store_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_store_events
    ADD CONSTRAINT event_store_events_pkey PRIMARY KEY (id);


--
-- Name: flipper_features flipper_features_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.flipper_features
    ADD CONSTRAINT flipper_features_pkey PRIMARY KEY (id);


--
-- Name: flipper_gates flipper_gates_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.flipper_gates
    ADD CONSTRAINT flipper_gates_pkey PRIMARY KEY (id);


--
-- Name: language_categories language_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_categories
    ADD CONSTRAINT language_categories_pkey PRIMARY KEY (id);


--
-- Name: language_category_items language_category_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_category_items
    ADD CONSTRAINT language_category_items_pkey PRIMARY KEY (id);


--
-- Name: language_category_qna_items language_category_qna_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_category_qna_items
    ADD CONSTRAINT language_category_qna_items_pkey PRIMARY KEY (id);


--
-- Name: language_landing_page_qna_items language_landing_page_qna_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_landing_page_qna_items
    ADD CONSTRAINT language_landing_page_qna_items_pkey PRIMARY KEY (id);


--
-- Name: language_landing_pages language_landing_pages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_landing_pages
    ADD CONSTRAINT language_landing_pages_pkey PRIMARY KEY (id);


--
-- Name: language_lesson_members language_lesson_members_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_members
    ADD CONSTRAINT language_lesson_members_pkey PRIMARY KEY (id);


--
-- Name: language_lesson_reviews language_lesson_reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_reviews
    ADD CONSTRAINT language_lesson_reviews_pkey PRIMARY KEY (id);


--
-- Name: language_lesson_version_infos language_lesson_version_infos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_version_infos
    ADD CONSTRAINT language_lesson_version_infos_pkey PRIMARY KEY (id);


--
-- Name: language_lesson_versions language_lesson_versions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_versions
    ADD CONSTRAINT language_lesson_versions_pkey PRIMARY KEY (id);


--
-- Name: language_lessons language_lessons_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lessons
    ADD CONSTRAINT language_lessons_pkey PRIMARY KEY (id);


--
-- Name: language_members language_members_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_members
    ADD CONSTRAINT language_members_pkey PRIMARY KEY (id);


--
-- Name: language_module_descriptions language_module_descriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_descriptions
    ADD CONSTRAINT language_module_descriptions_pkey PRIMARY KEY (id);


--
-- Name: language_module_version_infos language_module_version_infos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_version_infos
    ADD CONSTRAINT language_module_version_infos_pkey PRIMARY KEY (id);


--
-- Name: language_module_versions language_module_versions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_versions
    ADD CONSTRAINT language_module_versions_pkey PRIMARY KEY (id);


--
-- Name: language_modules language_modules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_modules
    ADD CONSTRAINT language_modules_pkey PRIMARY KEY (id);


--
-- Name: language_version_infos language_version_infos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_version_infos
    ADD CONSTRAINT language_version_infos_pkey PRIMARY KEY (id);


--
-- Name: language_versions language_versions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_versions
    ADD CONSTRAINT language_versions_pkey PRIMARY KEY (id);


--
-- Name: languages languages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT languages_pkey PRIMARY KEY (id);


--
-- Name: leads leads_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.leads
    ADD CONSTRAINT leads_pkey PRIMARY KEY (id);


--
-- Name: reviews reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: solid_cable_messages solid_cable_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_cable_messages
    ADD CONSTRAINT solid_cable_messages_pkey PRIMARY KEY (id);


--
-- Name: solid_cache_entries solid_cache_entries_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_cache_entries
    ADD CONSTRAINT solid_cache_entries_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_blocked_executions solid_queue_blocked_executions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_blocked_executions
    ADD CONSTRAINT solid_queue_blocked_executions_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_claimed_executions solid_queue_claimed_executions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_claimed_executions
    ADD CONSTRAINT solid_queue_claimed_executions_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_failed_executions solid_queue_failed_executions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_failed_executions
    ADD CONSTRAINT solid_queue_failed_executions_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_jobs solid_queue_jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_jobs
    ADD CONSTRAINT solid_queue_jobs_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_pauses solid_queue_pauses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_pauses
    ADD CONSTRAINT solid_queue_pauses_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_processes solid_queue_processes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_processes
    ADD CONSTRAINT solid_queue_processes_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_ready_executions solid_queue_ready_executions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_ready_executions
    ADD CONSTRAINT solid_queue_ready_executions_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_recurring_executions solid_queue_recurring_executions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_recurring_executions
    ADD CONSTRAINT solid_queue_recurring_executions_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_recurring_tasks solid_queue_recurring_tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_recurring_tasks
    ADD CONSTRAINT solid_queue_recurring_tasks_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_scheduled_executions solid_queue_scheduled_executions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_scheduled_executions
    ADD CONSTRAINT solid_queue_scheduled_executions_pkey PRIMARY KEY (id);


--
-- Name: solid_queue_semaphores solid_queue_semaphores_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_semaphores
    ADD CONSTRAINT solid_queue_semaphores_pkey PRIMARY KEY (id);


--
-- Name: staff_member_role_permissions staff_member_role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_member_role_permissions
    ADD CONSTRAINT staff_member_role_permissions_pkey PRIMARY KEY (id);


--
-- Name: staff_member_roles staff_member_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_member_roles
    ADD CONSTRAINT staff_member_roles_pkey PRIMARY KEY (id);


--
-- Name: staff_members staff_members_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_members
    ADD CONSTRAINT staff_members_pkey PRIMARY KEY (id);


--
-- Name: survey_answers survey_answers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_answers
    ADD CONSTRAINT survey_answers_pkey PRIMARY KEY (id);


--
-- Name: survey_items survey_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_items
    ADD CONSTRAINT survey_items_pkey PRIMARY KEY (id);


--
-- Name: survey_scenario_items survey_scenario_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_items
    ADD CONSTRAINT survey_scenario_items_pkey PRIMARY KEY (id);


--
-- Name: survey_scenario_members survey_scenario_members_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_members
    ADD CONSTRAINT survey_scenario_members_pkey PRIMARY KEY (id);


--
-- Name: survey_scenario_triggers survey_scenario_triggers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_triggers
    ADD CONSTRAINT survey_scenario_triggers_pkey PRIMARY KEY (id);


--
-- Name: survey_scenarios survey_scenarios_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenarios
    ADD CONSTRAINT survey_scenarios_pkey PRIMARY KEY (id);


--
-- Name: surveys surveys_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.surveys
    ADD CONSTRAINT surveys_pkey PRIMARY KEY (id);


--
-- Name: taggings taggings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.taggings
    ADD CONSTRAINT taggings_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: uploads uploads_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uploads
    ADD CONSTRAINT uploads_pkey PRIMARY KEY (id);


--
-- Name: user_accounts user_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_accounts
    ADD CONSTRAINT user_accounts_pkey PRIMARY KEY (id);


--
-- Name: user_credentials user_credentials_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_credentials
    ADD CONSTRAINT user_credentials_pkey PRIMARY KEY (id);


--
-- Name: user_survey_pivots user_survey_pivots_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_survey_pivots
    ADD CONSTRAINT user_survey_pivots_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_on_language_landing_page_id_98023e1f90; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_on_language_landing_page_id_98023e1f90 ON public.language_landing_page_qna_items USING btree (language_landing_page_id);


--
-- Name: idx_on_language_lesson_version_info_id_e5ef52eeca; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_on_language_lesson_version_info_id_e5ef52eeca ON public.language_lesson_reviews USING btree (language_lesson_version_info_id);


--
-- Name: index_action_text_rich_texts_uniqueness; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_action_text_rich_texts_uniqueness ON public.action_text_rich_texts USING btree (record_type, record_id, name);


--
-- Name: index_active_storage_attachments_on_blob_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_active_storage_attachments_on_blob_id ON public.active_storage_attachments USING btree (blob_id);


--
-- Name: index_active_storage_attachments_uniqueness; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_active_storage_attachments_uniqueness ON public.active_storage_attachments USING btree (record_type, record_id, name, blob_id);


--
-- Name: index_active_storage_blobs_on_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_active_storage_blobs_on_key ON public.active_storage_blobs USING btree (key);


--
-- Name: index_active_storage_variant_records_uniqueness; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_active_storage_variant_records_uniqueness ON public.active_storage_variant_records USING btree (blob_id, variation_digest);


--
-- Name: index_ahoy_events_on_name_and_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ahoy_events_on_name_and_time ON public.ahoy_events USING btree (name, "time");


--
-- Name: index_ahoy_events_on_properties; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ahoy_events_on_properties ON public.ahoy_events USING gin (properties jsonb_path_ops);


--
-- Name: index_ahoy_events_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ahoy_events_on_user_id ON public.ahoy_events USING btree (user_id);


--
-- Name: index_ahoy_events_on_visit_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ahoy_events_on_visit_id ON public.ahoy_events USING btree (visit_id);


--
-- Name: index_ahoy_visits_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ahoy_visits_on_user_id ON public.ahoy_visits USING btree (user_id);


--
-- Name: index_ahoy_visits_on_visit_token; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_ahoy_visits_on_visit_token ON public.ahoy_visits USING btree (visit_token);


--
-- Name: index_ahoy_visits_on_visitor_token_and_started_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ahoy_visits_on_visitor_token_and_started_at ON public.ahoy_visits USING btree (visitor_token, started_at);


--
-- Name: index_ai_chats_on_ai_model_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_chats_on_ai_model_id ON public.ai_chats USING btree (ai_model_id);


--
-- Name: index_ai_chats_on_language_lesson_member_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_chats_on_language_lesson_member_id ON public.ai_chats USING btree (language_lesson_member_id);


--
-- Name: index_ai_chats_on_user_id_and_language_lesson_member_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_ai_chats_on_user_id_and_language_lesson_member_id ON public.ai_chats USING btree (user_id, language_lesson_member_id);


--
-- Name: index_ai_messages_on_ai_chat_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_messages_on_ai_chat_id ON public.ai_messages USING btree (ai_chat_id);


--
-- Name: index_ai_messages_on_ai_model_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_messages_on_ai_model_id ON public.ai_messages USING btree (ai_model_id);


--
-- Name: index_ai_messages_on_ai_tool_call_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_messages_on_ai_tool_call_id ON public.ai_messages USING btree (ai_tool_call_id);


--
-- Name: index_ai_messages_on_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_messages_on_role ON public.ai_messages USING btree (role);


--
-- Name: index_ai_messages_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_messages_on_user_id ON public.ai_messages USING btree (user_id);


--
-- Name: index_ai_models_on_capabilities; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_models_on_capabilities ON public.ai_models USING gin (capabilities);


--
-- Name: index_ai_models_on_family; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_models_on_family ON public.ai_models USING btree (family);


--
-- Name: index_ai_models_on_modalities; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_models_on_modalities ON public.ai_models USING gin (modalities);


--
-- Name: index_ai_models_on_provider_and_model_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_ai_models_on_provider_and_model_id ON public.ai_models USING btree (provider, model_id);


--
-- Name: index_ai_tool_calls_on_ai_message_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_tool_calls_on_ai_message_id ON public.ai_tool_calls USING btree (ai_message_id);


--
-- Name: index_ai_tool_calls_on_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_ai_tool_calls_on_name ON public.ai_tool_calls USING btree (name);


--
-- Name: index_ai_tool_calls_on_tool_call_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_ai_tool_calls_on_tool_call_id ON public.ai_tool_calls USING btree (tool_call_id);


--
-- Name: index_banners_on_locale_and_state; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_banners_on_locale_and_state ON public.banners USING btree (locale, state);


--
-- Name: index_blog_post_likes_on_blog_post_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_blog_post_likes_on_blog_post_id ON public.blog_post_likes USING btree (blog_post_id);


--
-- Name: index_blog_post_likes_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_blog_post_likes_on_user_id ON public.blog_post_likes USING btree (user_id);


--
-- Name: index_blog_post_related_language_items_on_blog_post_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_blog_post_related_language_items_on_blog_post_id ON public.blog_post_related_language_items USING btree (blog_post_id);


--
-- Name: index_blog_post_related_language_items_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_blog_post_related_language_items_on_language_id ON public.blog_post_related_language_items USING btree (language_id);


--
-- Name: index_blog_posts_on_creator_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_blog_posts_on_creator_id ON public.blog_posts USING btree (creator_id);


--
-- Name: index_blog_posts_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_blog_posts_on_language_id ON public.blog_posts USING btree (language_id);


--
-- Name: index_blog_posts_on_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_blog_posts_on_slug ON public.blog_posts USING btree (slug);


--
-- Name: index_book_requests_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_book_requests_on_user_id ON public.book_requests USING btree (user_id);


--
-- Name: index_event_store_events_in_streams_on_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_event_store_events_in_streams_on_created_at ON public.event_store_events_in_streams USING btree (created_at);


--
-- Name: index_event_store_events_in_streams_on_event_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_event_store_events_in_streams_on_event_id ON public.event_store_events_in_streams USING btree (event_id);


--
-- Name: index_event_store_events_in_streams_on_stream_and_event_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_event_store_events_in_streams_on_stream_and_event_id ON public.event_store_events_in_streams USING btree (stream, event_id);


--
-- Name: index_event_store_events_in_streams_on_stream_and_position; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_event_store_events_in_streams_on_stream_and_position ON public.event_store_events_in_streams USING btree (stream, "position");


--
-- Name: index_event_store_events_on_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_event_store_events_on_created_at ON public.event_store_events USING btree (created_at);


--
-- Name: index_event_store_events_on_event_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_event_store_events_on_event_id ON public.event_store_events USING btree (event_id);


--
-- Name: index_event_store_events_on_event_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_event_store_events_on_event_type ON public.event_store_events USING btree (event_type);


--
-- Name: index_event_store_events_on_valid_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_event_store_events_on_valid_at ON public.event_store_events USING btree (valid_at);


--
-- Name: index_flipper_features_on_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_flipper_features_on_key ON public.flipper_features USING btree (key);


--
-- Name: index_flipper_gates_on_feature_key_and_key_and_value; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_flipper_gates_on_feature_key_and_key_and_value ON public.flipper_gates USING btree (feature_key, key, value);


--
-- Name: index_language_category_items_on_language_category_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_category_items_on_language_category_id ON public.language_category_items USING btree (language_category_id);


--
-- Name: index_language_category_items_on_language_landing_page_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_category_items_on_language_landing_page_id ON public.language_category_items USING btree (language_landing_page_id);


--
-- Name: index_language_category_qna_items_on_language_category_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_category_qna_items_on_language_category_id ON public.language_category_qna_items USING btree (language_category_id);


--
-- Name: index_language_landing_pages_on_landing_page_to_redirect_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_landing_pages_on_landing_page_to_redirect_id ON public.language_landing_pages USING btree (landing_page_to_redirect_id);


--
-- Name: index_language_landing_pages_on_language_category_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_landing_pages_on_language_category_id ON public.language_landing_pages USING btree (language_category_id);


--
-- Name: index_language_landing_pages_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_landing_pages_on_language_id ON public.language_landing_pages USING btree (language_id);


--
-- Name: index_language_lesson_members_on_language_member_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_members_on_language_member_id ON public.language_lesson_members USING btree (language_member_id);


--
-- Name: index_language_lesson_reviews_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_reviews_on_language_id ON public.language_lesson_reviews USING btree (language_id);


--
-- Name: index_language_lesson_reviews_on_language_lesson_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_reviews_on_language_lesson_id ON public.language_lesson_reviews USING btree (language_lesson_id);


--
-- Name: index_language_lesson_reviews_on_language_lesson_version_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_reviews_on_language_lesson_version_id ON public.language_lesson_reviews USING btree (language_lesson_version_id);


--
-- Name: index_language_lesson_version_infos_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_version_infos_on_language_id ON public.language_lesson_version_infos USING btree (language_id);


--
-- Name: index_language_lesson_version_infos_on_language_lesson_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_version_infos_on_language_lesson_id ON public.language_lesson_version_infos USING btree (language_lesson_id);


--
-- Name: index_language_lesson_version_infos_on_language_version_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_version_infos_on_language_version_id ON public.language_lesson_version_infos USING btree (language_version_id);


--
-- Name: index_language_lesson_version_infos_on_version_id_and_locale; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_version_infos_on_version_id_and_locale ON public.language_lesson_version_infos USING btree (version_id, locale);


--
-- Name: index_language_lesson_versions_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_versions_on_language_id ON public.language_lesson_versions USING btree (language_id);


--
-- Name: index_language_lesson_versions_on_language_version_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_versions_on_language_version_id ON public.language_lesson_versions USING btree (language_version_id);


--
-- Name: index_language_lesson_versions_on_lesson_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_versions_on_lesson_id ON public.language_lesson_versions USING btree (lesson_id);


--
-- Name: index_language_lesson_versions_on_module_version_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_lesson_versions_on_module_version_id ON public.language_lesson_versions USING btree (module_version_id);


--
-- Name: index_language_lessons_on_language_id_and_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_language_lessons_on_language_id_and_slug ON public.language_lessons USING btree (language_id, slug);


--
-- Name: index_language_members_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_members_on_language_id ON public.language_members USING btree (language_id);


--
-- Name: index_language_members_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_members_on_user_id ON public.language_members USING btree (user_id);


--
-- Name: index_language_module_version_infos_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_module_version_infos_on_language_id ON public.language_module_version_infos USING btree (language_id);


--
-- Name: index_language_module_version_infos_on_language_version_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_module_version_infos_on_language_version_id ON public.language_module_version_infos USING btree (language_version_id);


--
-- Name: index_language_module_versions_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_module_versions_on_language_id ON public.language_module_versions USING btree (language_id);


--
-- Name: index_language_module_versions_on_language_version_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_module_versions_on_language_version_id ON public.language_module_versions USING btree (language_version_id);


--
-- Name: index_language_module_versions_on_module_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_module_versions_on_module_id ON public.language_module_versions USING btree (module_id);


--
-- Name: index_language_version_infos_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_version_infos_on_language_id ON public.language_version_infos USING btree (language_id);


--
-- Name: index_language_version_infos_on_language_version_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_version_infos_on_language_version_id ON public.language_version_infos USING btree (language_version_id);


--
-- Name: index_language_versions_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_language_versions_on_language_id ON public.language_versions USING btree (language_id);


--
-- Name: index_languages_on_category_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_languages_on_category_id ON public.languages USING btree (category_id);


--
-- Name: index_languages_on_current_version_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_languages_on_current_version_id ON public.languages USING btree (current_version_id);


--
-- Name: index_leads_on_ahoy_visit_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_leads_on_ahoy_visit_id ON public.leads USING btree (ahoy_visit_id);


--
-- Name: index_leads_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_leads_on_user_id ON public.leads USING btree (user_id);


--
-- Name: index_reviews_on_language_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_reviews_on_language_id ON public.reviews USING btree (language_id);


--
-- Name: index_reviews_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_reviews_on_user_id ON public.reviews USING btree (user_id);


--
-- Name: index_sessions_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_sessions_on_user_id ON public.sessions USING btree (user_id);


--
-- Name: index_solid_cable_messages_on_channel; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_cable_messages_on_channel ON public.solid_cable_messages USING btree (channel);


--
-- Name: index_solid_cable_messages_on_channel_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_cable_messages_on_channel_hash ON public.solid_cable_messages USING btree (channel_hash);


--
-- Name: index_solid_cable_messages_on_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_cable_messages_on_created_at ON public.solid_cable_messages USING btree (created_at);


--
-- Name: index_solid_cache_entries_on_byte_size; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_cache_entries_on_byte_size ON public.solid_cache_entries USING btree (byte_size);


--
-- Name: index_solid_cache_entries_on_key_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_cache_entries_on_key_hash ON public.solid_cache_entries USING btree (key_hash);


--
-- Name: index_solid_cache_entries_on_key_hash_and_byte_size; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_cache_entries_on_key_hash_and_byte_size ON public.solid_cache_entries USING btree (key_hash, byte_size);


--
-- Name: index_solid_queue_blocked_executions_for_maintenance; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_blocked_executions_for_maintenance ON public.solid_queue_blocked_executions USING btree (expires_at, concurrency_key);


--
-- Name: index_solid_queue_blocked_executions_for_release; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_blocked_executions_for_release ON public.solid_queue_blocked_executions USING btree (concurrency_key, priority, job_id);


--
-- Name: index_solid_queue_blocked_executions_on_job_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_blocked_executions_on_job_id ON public.solid_queue_blocked_executions USING btree (job_id);


--
-- Name: index_solid_queue_claimed_executions_on_job_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_claimed_executions_on_job_id ON public.solid_queue_claimed_executions USING btree (job_id);


--
-- Name: index_solid_queue_claimed_executions_on_process_id_and_job_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_claimed_executions_on_process_id_and_job_id ON public.solid_queue_claimed_executions USING btree (process_id, job_id);


--
-- Name: index_solid_queue_dispatch_all; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_dispatch_all ON public.solid_queue_scheduled_executions USING btree (scheduled_at, priority, job_id);


--
-- Name: index_solid_queue_failed_executions_on_job_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_failed_executions_on_job_id ON public.solid_queue_failed_executions USING btree (job_id);


--
-- Name: index_solid_queue_jobs_for_alerting; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_jobs_for_alerting ON public.solid_queue_jobs USING btree (scheduled_at, finished_at);


--
-- Name: index_solid_queue_jobs_for_filtering; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_jobs_for_filtering ON public.solid_queue_jobs USING btree (queue_name, finished_at);


--
-- Name: index_solid_queue_jobs_on_active_job_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_jobs_on_active_job_id ON public.solid_queue_jobs USING btree (active_job_id);


--
-- Name: index_solid_queue_jobs_on_class_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_jobs_on_class_name ON public.solid_queue_jobs USING btree (class_name);


--
-- Name: index_solid_queue_jobs_on_finished_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_jobs_on_finished_at ON public.solid_queue_jobs USING btree (finished_at);


--
-- Name: index_solid_queue_pauses_on_queue_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_pauses_on_queue_name ON public.solid_queue_pauses USING btree (queue_name);


--
-- Name: index_solid_queue_poll_all; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_poll_all ON public.solid_queue_ready_executions USING btree (priority, job_id);


--
-- Name: index_solid_queue_poll_by_queue; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_poll_by_queue ON public.solid_queue_ready_executions USING btree (queue_name, priority, job_id);


--
-- Name: index_solid_queue_processes_on_last_heartbeat_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_processes_on_last_heartbeat_at ON public.solid_queue_processes USING btree (last_heartbeat_at);


--
-- Name: index_solid_queue_processes_on_name_and_supervisor_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_processes_on_name_and_supervisor_id ON public.solid_queue_processes USING btree (name, supervisor_id);


--
-- Name: index_solid_queue_processes_on_supervisor_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_processes_on_supervisor_id ON public.solid_queue_processes USING btree (supervisor_id);


--
-- Name: index_solid_queue_ready_executions_on_job_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_ready_executions_on_job_id ON public.solid_queue_ready_executions USING btree (job_id);


--
-- Name: index_solid_queue_recurring_executions_on_job_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_recurring_executions_on_job_id ON public.solid_queue_recurring_executions USING btree (job_id);


--
-- Name: index_solid_queue_recurring_executions_on_task_key_and_run_at; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_recurring_executions_on_task_key_and_run_at ON public.solid_queue_recurring_executions USING btree (task_key, run_at);


--
-- Name: index_solid_queue_recurring_tasks_on_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_recurring_tasks_on_key ON public.solid_queue_recurring_tasks USING btree (key);


--
-- Name: index_solid_queue_recurring_tasks_on_static; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_recurring_tasks_on_static ON public.solid_queue_recurring_tasks USING btree (static);


--
-- Name: index_solid_queue_scheduled_executions_on_job_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_scheduled_executions_on_job_id ON public.solid_queue_scheduled_executions USING btree (job_id);


--
-- Name: index_solid_queue_semaphores_on_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_semaphores_on_expires_at ON public.solid_queue_semaphores USING btree (expires_at);


--
-- Name: index_solid_queue_semaphores_on_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_solid_queue_semaphores_on_key ON public.solid_queue_semaphores USING btree (key);


--
-- Name: index_solid_queue_semaphores_on_key_and_value; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_solid_queue_semaphores_on_key_and_value ON public.solid_queue_semaphores USING btree (key, value);


--
-- Name: index_staff_member_role_permissions_on_role_id_and_resource; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_staff_member_role_permissions_on_role_id_and_resource ON public.staff_member_role_permissions USING btree (role_id, resource);


--
-- Name: index_staff_members_on_role_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_staff_members_on_role_id ON public.staff_members USING btree (role_id);


--
-- Name: index_staff_members_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_staff_members_on_user_id ON public.staff_members USING btree (user_id);


--
-- Name: index_survey_answers_on_survey_id_and_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_survey_answers_on_survey_id_and_user_id ON public.survey_answers USING btree (survey_id, user_id);


--
-- Name: index_survey_answers_on_survey_item_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_survey_answers_on_survey_item_id ON public.survey_answers USING btree (survey_item_id);


--
-- Name: index_survey_answers_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_survey_answers_on_user_id ON public.survey_answers USING btree (user_id);


--
-- Name: index_survey_items_on_survey_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_survey_items_on_survey_id ON public.survey_items USING btree (survey_id);


--
-- Name: index_survey_scenario_items_on_scenario_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_survey_scenario_items_on_scenario_id ON public.survey_scenario_items USING btree (scenario_id);


--
-- Name: index_survey_scenario_items_on_survey_id_and_scenario_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_survey_scenario_items_on_survey_id_and_scenario_id ON public.survey_scenario_items USING btree (survey_id, scenario_id);


--
-- Name: index_survey_scenario_members_on_scenario_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_survey_scenario_members_on_scenario_id ON public.survey_scenario_members USING btree (scenario_id);


--
-- Name: index_survey_scenario_members_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_survey_scenario_members_on_user_id ON public.survey_scenario_members USING btree (user_id);


--
-- Name: index_survey_scenario_triggers_on_event_name_and_scenario_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_survey_scenario_triggers_on_event_name_and_scenario_id ON public.survey_scenario_triggers USING btree (event_name, scenario_id);


--
-- Name: index_survey_scenario_triggers_on_scenario_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_survey_scenario_triggers_on_scenario_id ON public.survey_scenario_triggers USING btree (scenario_id);


--
-- Name: index_survey_scenarios_on_survey_item_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_survey_scenarios_on_survey_item_id ON public.survey_scenarios USING btree (survey_item_id);


--
-- Name: index_surveys_on_parent_survey_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_surveys_on_parent_survey_id ON public.surveys USING btree (parent_survey_id);


--
-- Name: index_surveys_on_parent_survey_item_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_surveys_on_parent_survey_item_id ON public.surveys USING btree (parent_survey_item_id);


--
-- Name: index_surveys_on_slug_and_locale; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_surveys_on_slug_and_locale ON public.surveys USING btree (slug, locale);


--
-- Name: index_taggings_on_context; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_context ON public.taggings USING btree (context);


--
-- Name: index_taggings_on_tag_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_tag_id ON public.taggings USING btree (tag_id);


--
-- Name: index_taggings_on_taggable_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_taggable_id ON public.taggings USING btree (taggable_id);


--
-- Name: index_taggings_on_taggable_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_taggable_type ON public.taggings USING btree (taggable_type);


--
-- Name: index_taggings_on_taggable_type_and_taggable_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_taggable_type_and_taggable_id ON public.taggings USING btree (taggable_type, taggable_id);


--
-- Name: index_taggings_on_tagger_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_tagger_id ON public.taggings USING btree (tagger_id);


--
-- Name: index_taggings_on_tagger_id_and_tagger_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_tagger_id_and_tagger_type ON public.taggings USING btree (tagger_id, tagger_type);


--
-- Name: index_taggings_on_tagger_type_and_tagger_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_tagger_type_and_tagger_id ON public.taggings USING btree (tagger_type, tagger_id);


--
-- Name: index_taggings_on_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_taggings_on_tenant ON public.taggings USING btree (tenant);


--
-- Name: index_tags_on_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_tags_on_name ON public.tags USING btree (name);


--
-- Name: index_user_accounts_on_provider_and_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_user_accounts_on_provider_and_uid ON public.user_accounts USING btree (provider, uid);


--
-- Name: index_user_credentials_on_external_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_user_credentials_on_external_id ON public.user_credentials USING btree (external_id);


--
-- Name: index_user_credentials_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_user_credentials_on_user_id ON public.user_credentials USING btree (user_id);


--
-- Name: index_user_survey_pivots_on_coding_experience_item_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_user_survey_pivots_on_coding_experience_item_id ON public.user_survey_pivots USING btree (coding_experience_item_id);


--
-- Name: index_user_survey_pivots_on_goal_item_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_user_survey_pivots_on_goal_item_id ON public.user_survey_pivots USING btree (goal_item_id);


--
-- Name: index_user_survey_pivots_on_study_plan_item_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_user_survey_pivots_on_study_plan_item_id ON public.user_survey_pivots USING btree (study_plan_item_id);


--
-- Name: index_user_survey_pivots_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_user_survey_pivots_on_user_id ON public.user_survey_pivots USING btree (user_id);


--
-- Name: index_users_on_LOWER_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX "index_users_on_LOWER_email" ON public.users USING btree (lower((email)::text));


--
-- Name: index_users_on_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_users_on_email ON public.users USING btree (email);


--
-- Name: index_users_on_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_users_on_phone ON public.users USING btree (phone) WHERE (phone IS NOT NULL);


--
-- Name: index_users_on_webauthn_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_users_on_webauthn_id ON public.users USING btree (webauthn_id) WHERE (webauthn_id IS NOT NULL);


--
-- Name: language_module_descriptions_module_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX language_module_descriptions_module_id_index ON public.language_module_descriptions USING btree (module_id);


--
-- Name: language_module_lessons_module_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX language_module_lessons_module_id_index ON public.language_lessons USING btree (module_id);


--
-- Name: language_module_lessons_upload_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX language_module_lessons_upload_id_index ON public.language_lessons USING btree (upload_id);


--
-- Name: language_modules_language_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX language_modules_language_id_index ON public.language_modules USING btree (language_id);


--
-- Name: language_modules_upload_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX language_modules_upload_id_index ON public.language_modules USING btree (upload_id);


--
-- Name: languages_slug_index; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX languages_slug_index ON public.languages USING btree (slug);


--
-- Name: languages_upload_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX languages_upload_id_index ON public.languages USING btree (upload_id);


--
-- Name: taggings_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX taggings_idx ON public.taggings USING btree (tag_id, taggable_id, taggable_type, context, tagger_id, tagger_type);


--
-- Name: taggings_idy; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX taggings_idy ON public.taggings USING btree (taggable_id, taggable_type, tagger_id, context);


--
-- Name: taggings_taggable_context_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX taggings_taggable_context_idx ON public.taggings USING btree (taggable_id, taggable_type, context);


--
-- Name: user_finished_lessons_language_module_lesson_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX user_finished_lessons_language_module_lesson_id_index ON public.language_lesson_members USING btree (lesson_id);


--
-- Name: user_finished_lessons_user_id_language_module_lesson_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX user_finished_lessons_user_id_language_module_lesson_id_index ON public.language_lesson_members USING btree (user_id, lesson_id);


--
-- Name: language_landing_pages fk_rails_00381fb5f4; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_landing_pages
    ADD CONSTRAINT fk_rails_00381fb5f4 FOREIGN KEY (landing_page_to_redirect_id) REFERENCES public.language_landing_pages(id);


--
-- Name: language_module_version_infos fk_rails_03a490d994; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_version_infos
    ADD CONSTRAINT fk_rails_03a490d994 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: language_version_infos fk_rails_11a0eeec04; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_version_infos
    ADD CONSTRAINT fk_rails_11a0eeec04 FOREIGN KEY (language_version_id) REFERENCES public.language_versions(id);


--
-- Name: language_category_items fk_rails_1653d28321; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_category_items
    ADD CONSTRAINT fk_rails_1653d28321 FOREIGN KEY (language_category_id) REFERENCES public.language_categories(id);


--
-- Name: language_module_descriptions fk_rails_1cc3025f38; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_descriptions
    ADD CONSTRAINT fk_rails_1cc3025f38 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: leads fk_rails_1d08b36969; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.leads
    ADD CONSTRAINT fk_rails_1d08b36969 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: book_requests fk_rails_1e4327d03b; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.book_requests
    ADD CONSTRAINT fk_rails_1e4327d03b FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: ai_tool_calls fk_rails_1ebf54e503; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_tool_calls
    ADD CONSTRAINT fk_rails_1ebf54e503 FOREIGN KEY (ai_message_id) REFERENCES public.ai_messages(id);


--
-- Name: language_lesson_members fk_rails_1ff6af1a44; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_members
    ADD CONSTRAINT fk_rails_1ff6af1a44 FOREIGN KEY (lesson_id) REFERENCES public.language_lessons(id);


--
-- Name: user_survey_pivots fk_rails_23a507246e; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_survey_pivots
    ADD CONSTRAINT fk_rails_23a507246e FOREIGN KEY (goal_item_id) REFERENCES public.survey_items(id);


--
-- Name: language_lesson_versions fk_rails_28cb0a511a; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_versions
    ADD CONSTRAINT fk_rails_28cb0a511a FOREIGN KEY (lesson_id) REFERENCES public.language_lessons(id);


--
-- Name: staff_members fk_rails_2925f894a8; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_members
    ADD CONSTRAINT fk_rails_2925f894a8 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: language_module_versions fk_rails_2a94cf2f1a; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_versions
    ADD CONSTRAINT fk_rails_2a94cf2f1a FOREIGN KEY (module_id) REFERENCES public.language_modules(id);


--
-- Name: language_lesson_version_infos fk_rails_2b2de76835; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_version_infos
    ADD CONSTRAINT fk_rails_2b2de76835 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: blog_post_likes fk_rails_2c9e0d4a09; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_post_likes
    ADD CONSTRAINT fk_rails_2c9e0d4a09 FOREIGN KEY (blog_post_id) REFERENCES public.blog_posts(id);


--
-- Name: reviews fk_rails_2d5506f396; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT fk_rails_2d5506f396 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: survey_scenario_items fk_rails_2eb91bfb82; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_items
    ADD CONSTRAINT fk_rails_2eb91bfb82 FOREIGN KEY (survey_id) REFERENCES public.surveys(id);


--
-- Name: solid_queue_recurring_executions fk_rails_318a5533ed; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_recurring_executions
    ADD CONSTRAINT fk_rails_318a5533ed FOREIGN KEY (job_id) REFERENCES public.solid_queue_jobs(id) ON DELETE CASCADE;


--
-- Name: blog_post_likes fk_rails_326553c5d9; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_post_likes
    ADD CONSTRAINT fk_rails_326553c5d9 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: language_lessons fk_rails_36dee53ae2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lessons
    ADD CONSTRAINT fk_rails_36dee53ae2 FOREIGN KEY (module_id) REFERENCES public.language_modules(id);


--
-- Name: survey_answers fk_rails_3869acb601; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_answers
    ADD CONSTRAINT fk_rails_3869acb601 FOREIGN KEY (survey_id) REFERENCES public.surveys(id);


--
-- Name: language_modules fk_rails_39957b735c; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_modules
    ADD CONSTRAINT fk_rails_39957b735c FOREIGN KEY (upload_id) REFERENCES public.uploads(id);


--
-- Name: blog_posts fk_rails_399b6b9958; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_posts
    ADD CONSTRAINT fk_rails_399b6b9958 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: solid_queue_failed_executions fk_rails_39bbc7a631; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_failed_executions
    ADD CONSTRAINT fk_rails_39bbc7a631 FOREIGN KEY (job_id) REFERENCES public.solid_queue_jobs(id) ON DELETE CASCADE;


--
-- Name: language_category_qna_items fk_rails_3ba389c392; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_category_qna_items
    ADD CONSTRAINT fk_rails_3ba389c392 FOREIGN KEY (language_category_id) REFERENCES public.language_categories(id);


--
-- Name: language_lesson_version_infos fk_rails_3e2e50d44d; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_version_infos
    ADD CONSTRAINT fk_rails_3e2e50d44d FOREIGN KEY (language_version_id) REFERENCES public.language_versions(id);


--
-- Name: language_module_versions fk_rails_413afabb56; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_versions
    ADD CONSTRAINT fk_rails_413afabb56 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: blog_post_related_language_items fk_rails_423f6248fd; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_post_related_language_items
    ADD CONSTRAINT fk_rails_423f6248fd FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: language_members fk_rails_4343116778; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_members
    ADD CONSTRAINT fk_rails_4343116778 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: ai_messages fk_rails_438986158e; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_messages
    ADD CONSTRAINT fk_rails_438986158e FOREIGN KEY (ai_model_id) REFERENCES public.ai_models(id);


--
-- Name: language_lessons fk_rails_439e3990bb; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lessons
    ADD CONSTRAINT fk_rails_439e3990bb FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: blog_posts fk_rails_45731ba462; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_posts
    ADD CONSTRAINT fk_rails_45731ba462 FOREIGN KEY (creator_id) REFERENCES public.users(id);


--
-- Name: survey_scenario_members fk_rails_463c98f012; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_members
    ADD CONSTRAINT fk_rails_463c98f012 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: ai_chats fk_rails_4831695bfa; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_chats
    ADD CONSTRAINT fk_rails_4831695bfa FOREIGN KEY (language_lesson_member_id) REFERENCES public.language_lesson_members(id);


--
-- Name: language_module_versions fk_rails_4b0d9ee90f; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_versions
    ADD CONSTRAINT fk_rails_4b0d9ee90f FOREIGN KEY (language_version_id) REFERENCES public.language_versions(id);


--
-- Name: language_lesson_versions fk_rails_4cae0d7625; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_versions
    ADD CONSTRAINT fk_rails_4cae0d7625 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: solid_queue_blocked_executions fk_rails_4cd34e2228; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_blocked_executions
    ADD CONSTRAINT fk_rails_4cd34e2228 FOREIGN KEY (job_id) REFERENCES public.solid_queue_jobs(id) ON DELETE CASCADE;


--
-- Name: language_category_items fk_rails_50f4828843; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_category_items
    ADD CONSTRAINT fk_rails_50f4828843 FOREIGN KEY (language_landing_page_id) REFERENCES public.language_landing_pages(id);


--
-- Name: user_survey_pivots fk_rails_535f0749ed; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_survey_pivots
    ADD CONSTRAINT fk_rails_535f0749ed FOREIGN KEY (coding_experience_item_id) REFERENCES public.survey_items(id);


--
-- Name: language_lesson_reviews fk_rails_58fe2f4daf; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_reviews
    ADD CONSTRAINT fk_rails_58fe2f4daf FOREIGN KEY (language_lesson_version_info_id) REFERENCES public.language_lesson_version_infos(id);


--
-- Name: language_modules fk_rails_5c574dddfe; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_modules
    ADD CONSTRAINT fk_rails_5c574dddfe FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: user_survey_pivots fk_rails_5ca50bd2d9; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_survey_pivots
    ADD CONSTRAINT fk_rails_5ca50bd2d9 FOREIGN KEY (study_plan_item_id) REFERENCES public.survey_items(id);


--
-- Name: staff_member_role_permissions fk_rails_5d1aecb760; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_member_role_permissions
    ADD CONSTRAINT fk_rails_5d1aecb760 FOREIGN KEY (role_id) REFERENCES public.staff_member_roles(id);


--
-- Name: survey_scenario_members fk_rails_5db85e9ca1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_members
    ADD CONSTRAINT fk_rails_5db85e9ca1 FOREIGN KEY (scenario_id) REFERENCES public.survey_scenarios(id);


--
-- Name: survey_answers fk_rails_621f80522c; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_answers
    ADD CONSTRAINT fk_rails_621f80522c FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: ai_messages fk_rails_69368b3cd9; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_messages
    ADD CONSTRAINT fk_rails_69368b3cd9 FOREIGN KEY (ai_tool_call_id) REFERENCES public.ai_tool_calls(id);


--
-- Name: language_lesson_versions fk_rails_6b52e26355; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_versions
    ADD CONSTRAINT fk_rails_6b52e26355 FOREIGN KEY (language_version_id) REFERENCES public.language_versions(id);


--
-- Name: language_versions fk_rails_6cc776ff38; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_versions
    ADD CONSTRAINT fk_rails_6cc776ff38 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: language_version_infos fk_rails_6ffda837cf; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_version_infos
    ADD CONSTRAINT fk_rails_6ffda837cf FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: language_lesson_version_infos fk_rails_749db9acb4; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_version_infos
    ADD CONSTRAINT fk_rails_749db9acb4 FOREIGN KEY (language_lesson_id) REFERENCES public.language_lessons(id);


--
-- Name: reviews fk_rails_74a66bd6c5; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT fk_rails_74a66bd6c5 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: sessions fk_rails_758836b4f0; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT fk_rails_758836b4f0 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: ai_chats fk_rails_768e14b856; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_chats
    ADD CONSTRAINT fk_rails_768e14b856 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: language_landing_pages fk_rails_77f646d57e; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_landing_pages
    ADD CONSTRAINT fk_rails_77f646d57e FOREIGN KEY (language_category_id) REFERENCES public.language_categories(id);


--
-- Name: surveys fk_rails_7c2d7aa117; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.surveys
    ADD CONSTRAINT fk_rails_7c2d7aa117 FOREIGN KEY (parent_survey_item_id) REFERENCES public.survey_items(id);


--
-- Name: language_lesson_members fk_rails_7e60189e01; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_members
    ADD CONSTRAINT fk_rails_7e60189e01 FOREIGN KEY (language_member_id) REFERENCES public.language_members(id);


--
-- Name: staff_members fk_rails_7fe1e4b2a6; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff_members
    ADD CONSTRAINT fk_rails_7fe1e4b2a6 FOREIGN KEY (role_id) REFERENCES public.staff_member_roles(id);


--
-- Name: solid_queue_ready_executions fk_rails_81fcbd66af; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_ready_executions
    ADD CONSTRAINT fk_rails_81fcbd66af FOREIGN KEY (job_id) REFERENCES public.solid_queue_jobs(id) ON DELETE CASCADE;


--
-- Name: language_members fk_rails_86f99b4837; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_members
    ADD CONSTRAINT fk_rails_86f99b4837 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: language_lesson_version_infos fk_rails_8872f568fe; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_version_infos
    ADD CONSTRAINT fk_rails_8872f568fe FOREIGN KEY (version_id) REFERENCES public.language_lesson_versions(id);


--
-- Name: language_lesson_reviews fk_rails_8bf293e8a9; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_reviews
    ADD CONSTRAINT fk_rails_8bf293e8a9 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: survey_answers fk_rails_8cb45b53a9; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_answers
    ADD CONSTRAINT fk_rails_8cb45b53a9 FOREIGN KEY (survey_item_id) REFERENCES public.survey_items(id);


--
-- Name: language_lesson_reviews fk_rails_9041c890b0; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_reviews
    ADD CONSTRAINT fk_rails_9041c890b0 FOREIGN KEY (language_lesson_version_id) REFERENCES public.language_lesson_versions(id);


--
-- Name: survey_scenario_items fk_rails_924a4eb1ae; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_items
    ADD CONSTRAINT fk_rails_924a4eb1ae FOREIGN KEY (scenario_id) REFERENCES public.survey_scenarios(id);


--
-- Name: language_lesson_reviews fk_rails_93e176fee3; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_reviews
    ADD CONSTRAINT fk_rails_93e176fee3 FOREIGN KEY (language_lesson_id) REFERENCES public.language_lessons(id);


--
-- Name: ai_messages fk_rails_958af724a1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_messages
    ADD CONSTRAINT fk_rails_958af724a1 FOREIGN KEY (ai_chat_id) REFERENCES public.ai_chats(id);


--
-- Name: active_storage_variant_records fk_rails_993965df05; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_storage_variant_records
    ADD CONSTRAINT fk_rails_993965df05 FOREIGN KEY (blob_id) REFERENCES public.active_storage_blobs(id);


--
-- Name: user_credentials fk_rails_9b162a81f6; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_credentials
    ADD CONSTRAINT fk_rails_9b162a81f6 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: language_module_version_infos fk_rails_9c7b71797b; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_version_infos
    ADD CONSTRAINT fk_rails_9c7b71797b FOREIGN KEY (version_id) REFERENCES public.language_module_versions(id);


--
-- Name: language_lessons fk_rails_9ca0d6ceaa; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lessons
    ADD CONSTRAINT fk_rails_9ca0d6ceaa FOREIGN KEY (upload_id) REFERENCES public.uploads(id);


--
-- Name: solid_queue_claimed_executions fk_rails_9cfe4d4944; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_claimed_executions
    ADD CONSTRAINT fk_rails_9cfe4d4944 FOREIGN KEY (job_id) REFERENCES public.solid_queue_jobs(id) ON DELETE CASCADE;


--
-- Name: taggings fk_rails_9fcd2e236b; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.taggings
    ADD CONSTRAINT fk_rails_9fcd2e236b FOREIGN KEY (tag_id) REFERENCES public.tags(id);


--
-- Name: languages fk_rails_bcc060b35a; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT fk_rails_bcc060b35a FOREIGN KEY (upload_id) REFERENCES public.uploads(id);


--
-- Name: survey_items fk_rails_bdccd8c655; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_items
    ADD CONSTRAINT fk_rails_bdccd8c655 FOREIGN KEY (survey_id) REFERENCES public.surveys(id);


--
-- Name: language_landing_page_qna_items fk_rails_c3a22366db; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_landing_page_qna_items
    ADD CONSTRAINT fk_rails_c3a22366db FOREIGN KEY (language_landing_page_id) REFERENCES public.language_landing_pages(id);


--
-- Name: active_storage_attachments fk_rails_c3b3935057; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_storage_attachments
    ADD CONSTRAINT fk_rails_c3b3935057 FOREIGN KEY (blob_id) REFERENCES public.active_storage_blobs(id);


--
-- Name: solid_queue_scheduled_executions fk_rails_c4316f352d; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.solid_queue_scheduled_executions
    ADD CONSTRAINT fk_rails_c4316f352d FOREIGN KEY (job_id) REFERENCES public.solid_queue_jobs(id) ON DELETE CASCADE;


--
-- Name: survey_scenario_triggers fk_rails_c8040740a8; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenario_triggers
    ADD CONSTRAINT fk_rails_c8040740a8 FOREIGN KEY (scenario_id) REFERENCES public.survey_scenarios(id);


--
-- Name: surveys fk_rails_c84299bbb9; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.surveys
    ADD CONSTRAINT fk_rails_c84299bbb9 FOREIGN KEY (parent_survey_id) REFERENCES public.surveys(id);


--
-- Name: event_store_events_in_streams fk_rails_c8d52b5857; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_store_events_in_streams
    ADD CONSTRAINT fk_rails_c8d52b5857 FOREIGN KEY (event_id) REFERENCES public.event_store_events(event_id);


--
-- Name: languages fk_rails_cd49b170ef; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT fk_rails_cd49b170ef FOREIGN KEY (category_id) REFERENCES public.language_categories(id);


--
-- Name: user_accounts fk_rails_d64ac9bcc2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_accounts
    ADD CONSTRAINT fk_rails_d64ac9bcc2 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: blog_post_related_language_items fk_rails_d74a547da2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_post_related_language_items
    ADD CONSTRAINT fk_rails_d74a547da2 FOREIGN KEY (blog_post_id) REFERENCES public.blog_posts(id);


--
-- Name: leads fk_rails_daafd5abf0; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.leads
    ADD CONSTRAINT fk_rails_daafd5abf0 FOREIGN KEY (ahoy_visit_id) REFERENCES public.ahoy_visits(id);


--
-- Name: language_module_descriptions fk_rails_dac1bb6244; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_descriptions
    ADD CONSTRAINT fk_rails_dac1bb6244 FOREIGN KEY (module_id) REFERENCES public.language_modules(id);


--
-- Name: ai_chats fk_rails_e112a98abc; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_chats
    ADD CONSTRAINT fk_rails_e112a98abc FOREIGN KEY (ai_model_id) REFERENCES public.ai_models(id);


--
-- Name: language_module_version_infos fk_rails_e4269a6401; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_module_version_infos
    ADD CONSTRAINT fk_rails_e4269a6401 FOREIGN KEY (language_version_id) REFERENCES public.language_versions(id);


--
-- Name: ai_messages fk_rails_eaeb97aaaa; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_messages
    ADD CONSTRAINT fk_rails_eaeb97aaaa FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: language_lesson_members fk_rails_eb00b8f36c; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_members
    ADD CONSTRAINT fk_rails_eb00b8f36c FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_survey_pivots fk_rails_eb41950872; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_survey_pivots
    ADD CONSTRAINT fk_rails_eb41950872 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: survey_scenarios fk_rails_f2414b87c2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survey_scenarios
    ADD CONSTRAINT fk_rails_f2414b87c2 FOREIGN KEY (survey_item_id) REFERENCES public.survey_items(id);


--
-- Name: language_landing_pages fk_rails_f28dfb0a77; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_landing_pages
    ADD CONSTRAINT fk_rails_f28dfb0a77 FOREIGN KEY (language_id) REFERENCES public.languages(id);


--
-- Name: languages fk_rails_f53e1946e0; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT fk_rails_f53e1946e0 FOREIGN KEY (current_version_id) REFERENCES public.language_versions(id);


--
-- Name: language_lesson_versions fk_rails_fd3c3cf805; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.language_lesson_versions
    ADD CONSTRAINT fk_rails_fd3c3cf805 FOREIGN KEY (module_version_id) REFERENCES public.language_module_versions(id);


--
-- PostgreSQL database dump complete
--

\unrestrict c0FbiE23Lpu1cH2iAV84OvDoh9uiE1ttwYi3dgjyIoghGqxEP6M9VlD1Fqq3Rcn

