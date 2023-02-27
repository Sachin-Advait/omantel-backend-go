       
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

CREATE TABLE public.vop_table
(
    product character varying(64) NOT NULL,
    achieved_percent integer NOT NULL,
    vop_value NUMERIC (10, 3) NOT NULL,
    extra_pay NUMERIC (10, 3) NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,  
    PRIMARY KEY (product, achieved_percent, quarter, year)
);

CREATE TABLE public.max_payout_table
(
    position character varying(32) NOT NULL,
    product character varying(64) NOT NULL,
    max_pay NUMERIC (10, 3) NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,  
    PRIMARY KEY (position, product, quarter, year)
);

CREATE TABLE public.max_extra_payout_table
(
    position character varying(32) NOT NULL,
    product character varying(64) NOT NULL,
    max_pay NUMERIC (10, 3) NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    comment character varying(128) NOT NULL,
    quarter integer NOT NULL,
    year integer NOT NULL,  
    PRIMARY KEY (position, product, quarter, year)
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
    subcategory character varying(32) NOT NULL DEFAULT '',
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