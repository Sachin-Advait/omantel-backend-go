       
    ---------GREEN PLUM-------------

CREATE TABLE public.logs_table
(
    emp_id character varying(32),
    url character varying(64), 
    type character varying(32), 
    msg character varying(64), 
    code integer NOT NULL,
    text text, 
    created timestamp with time zone NOT NULL,
    PRIMARY KEY (emp_id, created)
);


CREATE TABLE public.user_master_table
(
    emp_id character varying(32),
    isAdmin bool DEFAULT false,
    isActive bool DEFAULT false,
    name character varying(64) NOT NULL DEFAULT '',
    usr_position character varying(32) NOT NULL DEFAULT '',
    region character varying(32) NOT NULL DEFAULT '',
    outlet character varying(64) NOT NULL DEFAULT '',
    remark character varying(128) NOT NULL DEFAULT '',
    mobile_number character varying(32) NOT NULL DEFAULT '',
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,
    PRIMARY KEY (emp_id, quarter, year)
);


CREATE TABLE public.regions_table
(
    region_code character varying(64),
    region_name character varying(64) NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,
    PRIMARY KEY (region_code, quarter, year)
);


CREATE TABLE public.outlets_table
(
    outlet_code character varying(64),
    outlet_name character varying(64) NOT NULL,
    region_code character varying(64) NOT NULL,
    region_name character varying(64) NOT NULL,
    subregion_name character varying(64) NOT NULL,
    latitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
    longitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,
    PRIMARY KEY (outlet_code, quarter, year)
);

CREATE TABLE public.user_hierarchy_table
(
    emp_id character varying(32),
    outlet_manager character varying(32) NOT NULL DEFAULT '',
    team_lead character varying(32) NOT NULL DEFAULT '',
    manager        character varying(32) NOT NULL DEFAULT '',
    senior_manager  character varying(32) NOT NULL DEFAULT '',
    general_manager character varying(32) NOT NULL DEFAULT '',
    vice_president character varying(32) NOT NULL DEFAULT '',
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,
    PRIMARY KEY (emp_id, quarter, year)
);


CREATE TABLE public.products_table
(
    product_code character varying(32),
    product_name character varying(64) NOT NULL,
    weightage NUMERIC (10, 3) NOT NULL DEFAULT 0,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,
    PRIMARY KEY (product_code, quarter, year)
);


CREATE TABLE public.plans_table
(
    plan_code character varying(128), -- TODO: once planCode is given this will be 32 in length
    plan_name character varying(128) NOT NULL,
    product_code character varying(32) NOT NULL,
    plan_price NUMERIC (10, 3) NOT NULL DEFAULT 0,
    sip_plan_rate NUMERIC (10, 3) NOT NULL DEFAULT 0,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,
    PRIMARY KEY (plan_code, quarter, year)
);


CREATE TABLE public.targets_table
(
    target_code character varying(32),
    region character varying(32) NOT NULL,
    outlet character varying(32) NOT NULL,
    month integer NOT NULL, 
    baqati integer NOT NULL,
    hbb integer NOT NULL,
    hayyak integer NOT NULL,
    upgrade integer NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,  
    PRIMARY KEY (target_code, quarter, year)
);

CREATE TABLE public.threshold_table
(
    id character varying(32),
    position character varying(32) NOT NULL,
    product character varying(64) NOT NULL,
    achieved_percent integer NOT NULL,
    payout NUMERIC (10, 3) NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,  
    PRIMARY KEY (position, product, achieved_percent, quarter, year)
);

CREATE TABLE public.voc_criteria_table
(
    unsatisfactory NUMERIC (10, 3) NOT NULL,
    almost_achieved NUMERIC (10, 3) NOT NULL,
    achieved NUMERIC (10, 3) NOT NULL,
    well_achieved NUMERIC (10, 3) NOT NULL,
    outstanding NUMERIC (10, 3) NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    year integer NOT NULL,  
    PRIMARY KEY (year)
);

CREATE TABLE public.voc_max_pay_table
(
    position character varying(32) NOT NULL,
    max_pay NUMERIC (10, 3) NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    year integer NOT NULL,  
    PRIMARY KEY (position, year)
);

CREATE TABLE public.voc_user_table
(
    emp_id character varying(32),
    achieved_percent NUMERIC (10, 3) NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    year integer NOT NULL,  
    PRIMARY KEY (emp_id, year)
);


CREATE TABLE public.sale_transaction_table
(
    id integer PRIMARY KEY,
    id_sbsc integer NOT NULL,
    region character varying(32) NOT NULL DEFAULT '',
    outlet character varying(64) NOT NULL DEFAULT '',
    product_group character varying(32) NOT NULL DEFAULT '',
    transaction_type character varying(32) NOT NULL DEFAULT '', -- new or upgrade
    prev_plan character varying(128) NOT NULL DEFAULT '',
    curr_plan character varying(128) NOT NULL DEFAULT '',
    quantity smallint NOT NULL,
    category character varying(32) NOT NULL DEFAULT '',
    business_unit character varying(32) NOT NULL DEFAULT '',
    subcategory character varying(32) NOT NULL DEFAULT '',
    final_srvc_cat character varying(32) NOT NULL DEFAULT '',
    sales_date date NOT NULL,
    terminate_date date NOT NULL,
    quarter integer NOT NULL DEFAULT 1,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    remark character varying(128) NOT NULL DEFAULT '',
    comment character varying(128) NOT NULL
);

CREATE TABLE public.sale_terminated_transactions
(
    id SERIAL PRIMARY KEY,
    id_sbsc integer NOT NULL,
    product_group character varying(32) NOT NULL DEFAULT '',
    curr_plan character varying(128) NOT NULL DEFAULT '',
    subcategory character varying(32) NOT NULL DEFAULT '',
    sales_date date NOT NULL,
    terminate_date date NOT NULL,
    quarter integer NOT NULL DEFAULT 1,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    remark character varying(128) NOT NULL DEFAULT '',
    comment character varying(128) NOT NULL
);


CREATE TABLE public.user_authentication_table
(
    emp_id character varying(32) PRIMARY KEY,
    user_active BOOLEAN NOT NULL DEFAULT false,
    otp integer NOT NULL,
    otp_resend_count integer NOT NULL DEFAULT 0,
    otp_verify_failed_count integer NOT NULL DEFAULT 0,
    otp_used BOOLEAN NOT NULL DEFAULT false,
    otp_create_time timestamp with time zone NOT NULL,
    otp_update_time timestamp with time zone NOT NULL,
    login_attempt_count integer NOT NULL DEFAULT 0,
    last_login_time timestamp with time zone,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    comment character varying(32) NOT NULL
);


CREATE TABLE public.user_reports_table
(
    emp_id character varying(32),
    position character varying(32) NOT NULL DEFAULT '',
    reporting_manager character varying(32) NOT NULL DEFAULT '',
    reporting_manager_position character varying(32) NOT NULL DEFAULT '',
    region character varying(32) NOT NULL DEFAULT '',
    outlet character varying(64) NOT NULL DEFAULT'',
    latitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
    longitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
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
    isquarterly bool NOT NULL DEFAULT false,
    month integer NOT NULL DEFAULT 0,
    quarter integer NOT NULL DEFAULT 0,
    year integer NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    PRIMARY KEY (emp_id, month, quarter, year)
);

CREATE TABLE public.manager_reports_table
(
    mngr_id character varying(32),
    position character varying(32) NOT NULL DEFAULT '',
    reporting_manager character varying(32) NOT NULL DEFAULT '',
    reporting_manager_position character varying(32) NOT NULL DEFAULT '',
    region character varying(32) NOT NULL DEFAULT '',
    outlet character varying(64) NOT NULL DEFAULT'',
    latitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
    longitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
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
    isquarterly bool NOT NULL DEFAULT false,
    month integer NOT NULL DEFAULT 0,
    quarter integer NOT NULL DEFAULT 0,
    year integer NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    PRIMARY KEY (mngr_id, month, quarter, year)
);

CREATE TABLE public.region_reports_table
(
    region character varying(32),
    hbb_sales_count integer NOT NULL DEFAULT 0,
    hbb_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
    hayyak_sales_count integer NOT NULL DEFAULT 0,
    hayyak_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
    baqati_sales_count integer NOT NULL DEFAULT 0,
    baqati_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
    upgrade_sales_count integer NOT NULL DEFAULT 0,
    upgrade_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
    isquarterly bool NOT NULL DEFAULT false,
    month integer NOT NULL DEFAULT 0,
    quarter integer NOT NULL DEFAULT 0,
    year integer NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    PRIMARY KEY (region, month, quarter, year)
);

CREATE TABLE public.outlet_reports_table
(
    outlet character varying(32),
    region character varying(32) NOT NULL DEFAULT '',
    latitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
    longitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
    hbb_sales_count integer NOT NULL DEFAULT 0,
    hbb_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
    hayyak_sales_count integer NOT NULL DEFAULT 0,
    hayyak_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
    baqati_sales_count integer NOT NULL DEFAULT 0,
    baqati_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
    upgrade_sales_count integer NOT NULL DEFAULT 0,
    upgrade_target NUMERIC (10, 3) NOT NULL DEFAULT 0,
    isquarterly bool NOT NULL DEFAULT false,
    month integer NOT NULL DEFAULT 0,
    quarter integer NOT NULL DEFAULT 0,
    year integer NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    PRIMARY KEY (outlet, month, quarter, year)
);

CREATE TABLE public.user_sale_count_table
(
    id SERIAL PRIMARY KEY,
    emp_id character varying(32),
    quantity smallint NOT NULL,
    category character varying(32) NOT NULL DEFAULT '',
    business_unit character varying(32) NOT NULL DEFAULT '',
    subcategory character varying(32) NOT NULL DEFAULT '',
    curr_plan character varying(128) NOT NULL DEFAULT '',
    final_srvc_cat character varying(32) NOT NULL DEFAULT '',
    isquarterly bool NOT NULL DEFAULT false,
    month integer NOT NULL DEFAULT 0,
    quarter integer NOT NULL DEFAULT 0,
    year integer NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
);