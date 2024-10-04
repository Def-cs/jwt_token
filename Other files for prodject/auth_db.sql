-- Adminer 4.8.1 PostgreSQL 15.8 (Debian 15.8-1.pgdg120+1) dump

DROP TABLE IF EXISTS "tokens";
CREATE TABLE "public"."tokens" (
    "user_id" uuid NOT NULL,
    "token" character varying NOT NULL,
    CONSTRAINT "tokens_token" UNIQUE ("token"),
    CONSTRAINT "tokens_user_id" UNIQUE ("user_id")
) WITH (oids = false);


DROP TABLE IF EXISTS "users";
CREATE TABLE "public"."users" (
    "id" uuid DEFAULT uuid_generate_v4() NOT NULL,
    "login" character varying(256) NOT NULL,
    "password" character varying(256) NOT NULL,
    "email" character varying(256) NOT NULL,
    CONSTRAINT "users_login" UNIQUE ("login"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


ALTER TABLE ONLY "public"."tokens" ADD CONSTRAINT "tokens_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE NOT DEFERRABLE;

-- 2024-10-04 16:58:03.683545+00
