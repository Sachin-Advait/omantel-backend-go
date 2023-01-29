# Incentive Calculator


##File Operation

1. Read sales input file containing following parameter:
    Id 
    Name     
    Position 
    Outlet   
    Region   
    CategoryType    
    PreviousProduct 
    CurrentProduct  
    Quantity
    Category        
    ProductGroup    

2. Db insertion for sale transaction
    if client exists:
        insert sales transaction:
            EmployeeID
            Outlet   
            Region   
            CategoryType    
            PreviousProduct 
            CurrentProduct  
            Quantity
            Category        
            ProductGroup  

3. Target Table
        Region Outlet Baqati HBB Hayyak Upgrade



## Db Details






/**** server details ****/
ip: 93.115.22.132
user: root
pass: MJh6VL4g5sK4tMve!
      USXXajHTzAKyaZkt!
      Bitool2021#$%


/**** new server details ****/
ip: 93.114.128.202
user: root
pass: 8tNk89VrAM4ct6XE!     http://omantel.moheoman.live


vim /var/lib/pgsql/13/data/log/postgresql-




/**** db details ****/
user: postgres
password: Secure@#!2021v
          Oracle12#$%fget

/*
	VP
	GM
	SM
	Manager
	OutletManager
	Promoter/CSR
*/







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
    mobile_number character varying(32) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    CONSTRAINT user_master_table_pkey PRIMARY KEY (emp_id)
)





CREATE TABLE IF NOT EXISTS public.user_authentication_table
(
emp_id bigint PRIMARY KEY,
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
modified_reason character varying(32) NOT NULL
)



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
)



CREATE TABLE IF NOT EXISTS public.product_plans_table
(
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
id SERIAL,
plan_name character varying(64) COLLATE pg_catalog."default" NOT NULL,
plan_price numeric(15,3) NOT NULL,
sip_plan_rate numeric(15,3) NOT NULL
)



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
CONSTRAINT user_master_table_pkey PRIMARY KEY (emp_id)
)



CREATE TABLE IF NOT EXISTS public.user_authentication_table
(
emp_id bigint PRIMARY KEY,
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
modified_reason character varying(32) NOT NULL
)



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
)



CREATE TABLE IF NOT EXISTS public.user_hierarchy_table
(
    emp_id INTEGER PRIMARY KEY,
	outlet_manager integer NOT NULL DEFAULT 0,
	manager        integer NOT NULL DEFAULT 0,
	senior_manager  integer NOT NULL DEFAULT 0,
	general_manager integer NOT NULL DEFAULT 0,
	vice_president integer NOT NULL DEFAULT 0
)

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
)

CREATE TABLE IF NOT EXISTS public.regions_table
(
region_code character varying(32) PRIMARY KEY,
region_name character varying(64) NOT NULL,
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
comment character varying(128) NOT NULL
)

CREATE TABLE IF NOT EXISTS public.outlets_table
(
outlet_code character varying(32) PRIMARY KEY,
outlet_name character varying(64) NOT NULL,
region_code character varying(32) NOT NULL,
latitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
longitude NUMERIC (10, 3) NOT NULL DEFAULT 0,
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
comment character varying(128) NOT NULL,
CONSTRAINT fk_region FOREIGN KEY(region_code)
REFERENCES regions_table(region_code)
ON DELETE CASCADE
ON UPDATE CASCADE
)

CREATE TABLE IF NOT EXISTS public.products_table
(
product_code character varying(32) PRIMARY KEY,
product_name character varying(64) NOT NULL,
weightage NUMERIC (10, 3) NOT NULL DEFAULT 0,
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
comment character varying(128) NOT NULL
)

CREATE TABLE IF NOT EXISTS public.plans_table
(
plan_code character varying(32) PRIMARY KEY,
plan_name character varying(64) NOT NULL,
product_code character varying(32) NOT NULL,
plan_price NUMERIC (10, 3) NOT NULL DEFAULT 0,
sip_plan_rate NUMERIC (10, 3) NOT NULL DEFAULT 0,
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
comment character varying(128) NOT NULL,
CONSTRAINT fk_region FOREIGN KEY(product_code)
REFERENCES products_table(product_code)
ON DELETE CASCADE
ON UPDATE CASCADE
)

CREATE TABLE IF NOT EXISTS public.targets_table
(
target_code character varying(64) PRIMARY KEY,
region_code character varying(32) NOT NULL,
outlet_code character varying(32) NOT NULL,
month integer NOT NULL, 
year integer NOT NULL,  
baqati integer NOT NULL,
hbb integer NOT NULL,
hayyak integer NOT NULL,
upgrade integer NOT NULL,
created timestamp with time zone NOT NULL,
modified timestamp with time zone NOT NULL,
comment character varying(128) NOT NULL,
CONSTRAINT fk_targets FOREIGN KEY(outlet_code)
REFERENCES outlets_table(outlet_code)
ON DELETE CASCADE
ON UPDATE CASCADE
)

CREATE TABLE IF NOT EXISTS public.logs_table
(
    emp_id INTEGER NOT NULL,
    url character varying(64), 
    type character varying(32), 
    msg character varying(64), 
    text text, 
    created timestamp with time zone NOT NULL,
    PRIMARY KEY(emp_id, created)
)

select a.emp_id, b.outlet_manager, b.manager, b.senior_manager, b.general_manager, b.vice_president from 
user_master_table a inner join user_hierarchy_table b 
on a.emp_id = b.emp_id where a.usr_position= 'PROMOTER'

insert into user_hierarchy_table(emp_id, outlet_manager, manager, senior_manager, general_manager, vice_president)
values(201,0,0,301,0,0)

delete from user_hierarchy_table


userid:     root
login:      ssh root@93.114.128.202
password:   fQdxzNQWj8LUgqad!

Postgress
Host- localhost
User - gis
Databases- gisdb
Password- gis@#42746

And port 8090 is routed on http://omantel.moheoman.live/


https://stackoverflow.com/questions/38822764/how-to-send-a-https-request-with-a-certificate-golang


./incentive_calculator --port 8090 --env uat

scp ./incentive_calculator root@93.115.22.132:/root/go/src/gitlab.com/rendezvousPoint/incentive_calculator

ssh root@93.115.22.132
sseLbscK3xZ9AE2P!

NEW
ssh root@86.105.252.94
pZ9WzHCaf8JrnP