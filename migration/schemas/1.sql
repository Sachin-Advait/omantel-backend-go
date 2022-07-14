CREATE TABLE IF NOT EXISTS public.user_master_table
(
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
emp_id bigint NOT NULL,
name character varying(64) COLLATE pg_catalog."default" NOT NULL,
usr_position character varying(32) COLLATE pg_catalog."default" NOT NULL,
region character varying(32) COLLATE pg_catalog."default" NOT NULL,
outlet character varying(64) COLLATE pg_catalog."default" NOT NULL,
remark character varying(128) COLLATE pg_catalog."default",
mobile_number character varying(32) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
CONSTRAINT user_master_table_pkey PRIMARY KEY (emp_id)
);

CREATE TABLE IF NOT EXISTS public.sale_transaction_table
(
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
id SERIAL PRIMARY KEY,
region character varying(32) COLLATE pg_catalog."default" NOT NULL,
outlet character varying(64) COLLATE pg_catalog."default" NOT NULL,
emp_id bigint NOT NULL,
category_type character varying(32) COLLATE pg_catalog."default" NOT NULL,
prev_product character varying(64) COLLATE pg_catalog."default",
curr_product character varying(64) COLLATE pg_catalog."default" NOT NULL,
quantity smallint NOT NULL,
category character varying(32) COLLATE pg_catalog."default" NOT NULL,
product_group character varying(32) COLLATE pg_catalog."default" NOT NULL,
sales_date date NOT NULL DEFAULT CURRENT_DATE,
quarter integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS public.sales_target_table
(
id SERIAL,
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
region character varying(32) COLLATE pg_catalog."default" NOT NULL,
outlet character varying(64) COLLATE pg_catalog."default" NOT NULL,
months integer NOT NULL,  
baqati integer NOT NULL,
hbb integer NOT NULL,
hayyak integer NOT NULL,
upgrade integer NOT NULL
);

CREATE TABLE IF NOT EXISTS public.product_plans_table
(
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
id SERIAL,
plan_name character varying(64) COLLATE pg_catalog."default" NOT NULL,
plan_price numeric(15,3) NOT NULL,
sip_plan_rate numeric(15,3) NOT NULL
);

CREATE TABLE IF NOT EXISTS public.user_hierarchy_table
(
    emp_id INTEGER PRIMARY KEY,
	outlet_manager integer NOT NULL DEFAULT 0,
	manager        integer NOT NULL DEFAULT 0,
	senior_manager  integer NOT NULL DEFAULT 0,
	general_manager integer NOT NULL DEFAULT 0,
	vice_president integer NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS public.outlet_details_table
(
id INTEGER PRIMARY KEY,
outlet character varying(32) COLLATE pg_catalog."default" NOT NULL,
region character varying(64) COLLATE pg_catalog."default" NOT NULL,
latitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
longitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
outlet_manager_id integer NOT NULL DEFAULT 0,
non_managerial_emp integer NOT NULL DEFAULT 0,
total_emp integer NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS public.user_authentication_table
(
emp_id SERIAL,
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
expirytime timestamp with time zone NOT NULL,
otp integer NOT NULL,
retry integer NOT NULL
);
ALter table public.user_authentication_table ADD UNIQUE (emp_id);

CREATE TABLE IF NOT EXISTS public.staff_daily_report_table
(
emp_id SERIAL,
created timestamp with time zone NOT NULL,
reporting_manager integer NOT NULL DEFAULT 0,
isquarterly bool NOT NULL DEFAULT false,
months integer NOT NULL DEFAULT 0,
quarter integer NOT NULL DEFAULT 0,
hayyak_sales_count integer NOT NULL DEFAULT 0,
hayyak_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
hayyak_achieved_percentage NUMERIC (10, 3) NOT NULL DEFAULT 0,
hayyak_payout NUMERIC (10, 3) NOT NULL DEFAULT 0,
hayyak_extra_payout NUMERIC (10, 3) NOT NULL DEFAULT 0,
hbb_sales_count integer NOT NULL DEFAULT 0,
hbb_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
hbb_achieved_percentage NUMERIC (10, 3) NOT NULL DEFAULT 0,
hbb_payout NUMERIC (10, 3) NOT NULL DEFAULT 0,
hbb_extra_payout NUMERIC (10, 3) NOT NULL DEFAULT 0,
baqati_sales_count integer NOT NULL DEFAULT 0,
baqati_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
baqati_achieved_percentage NUMERIC (10, 3) NOT NULL DEFAULT 0,
baqati_payout NUMERIC (10, 3) NOT NULL DEFAULT 0,
baqati_extra_payout NUMERIC (10, 3) NOT NULL DEFAULT 0,
upgrade_sales_count integer NOT NULL DEFAULT 0,
upgrade_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
uprade_achieved_percentage NUMERIC (10, 3) NOT NULL DEFAULT 0,
upgrade_payout NUMERIC (10, 3) NOT NULL DEFAULT 0,
upgrade_extra_payout NUMERIC (10, 3) NOT NULL DEFAULT 0,
total_payout numeric(10,3) NOT NULL DEFAULT 0,
hayyak_sales_count_per_plan text,
hbb_sales_count_per_plan text,
baqati_sales_count_per_plan text,
upgrade_sales_count_per_plan text,
CONSTRAINT ADD UNIQUE (emp_id, months, quarter)
);