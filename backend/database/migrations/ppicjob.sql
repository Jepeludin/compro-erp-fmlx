-- Table: public.ppic_links

CREATE SEQUENCE IF NOT EXISTS ppic_links_id_seq;

CREATE TABLE IF NOT EXISTS public.ppic_links
(
    id bigint NOT NULL DEFAULT nextval('ppic_links_id_seq'::regclass),
    source_schedule_id bigint NOT NULL,
    target_schedule_id bigint NOT NULL,
    link_type character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT '0'::character varying,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ppic_links_pkey PRIMARY KEY (id),
    CONSTRAINT ppic_links_source_schedule_id_fkey FOREIGN KEY (source_schedule_id)
        REFERENCES public.ppic_schedules (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT ppic_links_target_schedule_id_fkey FOREIGN KEY (target_schedule_id)
        REFERENCES public.ppic_schedules (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.ppic_links
    OWNER to postgres;

CREATE INDEX IF NOT EXISTS idx_ppic_links_source
    ON public.ppic_links USING btree
    (source_schedule_id ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_ppic_links_target
    ON public.ppic_links USING btree
    (target_schedule_id ASC NULLS LAST)
    TABLESPACE pg_default;
