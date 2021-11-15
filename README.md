# Incentive Calculator

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




/**** server details ****/
ip: 93.115.22.132
user: root
pass: LGNt9TW5k7zdebu7!


/**** db details ****/
user: postgres
password: gis@123






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



