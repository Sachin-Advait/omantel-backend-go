-- /sales/incentive/client --

select emp_id, name, usr_position, mobile_number, region, outlet, isAdmin from user_master_table where deleted IS NULL 
and emp_id='69963' and quarter=3 and year=2022

select outlet_manager from user_hierarchy_table where deleted IS NULL and emp_id='69963' and quarter=3 and year=2022

SELECT outlet_code, outlet_name, region_code, region_name, latitude, longitude FROM outlets_table
where deleted IS NULL and region_name='DHOFAR' and outlet_name='AWQAD OUTLET' and quarter=3 and year=2022

select id, prev_plan, curr_plan, quantity, product_group, transaction_type, subcategory, sales_date, terminate_date 
from sale_transaction_table where emp_id='69963' and region='DHOFAR' and outlet='AWQAD OUTLET' 
and quarter=3 and EXTRACT('year' from sales_date)=2022

select count(*) from user_master_table where deleted IS NULL and isactive=true and region='DHOFAR' and outlet='AWQAD OUTLET' 
and (usr_position='PROMOTER' OR usr_position='CSR') and quarter=3 and year=2022

-- Outlet Manager --

select emp_id, name, usr_position, mobile_number, region, outlet, isAdmin from user_master_table where deleted IS NULL 
and emp_id='50795' and quarter=3 and year=2022

select manager from user_hierarchy_table where deleted IS NULL and outlet_manager='50795' and quarter=3 and year=2022

SELECT outlet_code, outlet_name, region_code, region_name, latitude, longitude FROM outlets_table
where deleted IS NULL and region_name='DHOFAR' and outlet_name='AWQAD OUTLET' and quarter=3 and year=2022

select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.emp_id AND a.quarter = b.quarter 
AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL 
AND b.outlet_manager='50795' AND a.quarter=3 AND a.year=2022

select id, prev_plan, curr_plan, quantity, product_group, transaction_type, subcategory, sales_date, terminate_date 
from sale_transaction_table where emp_id='69963' and region='DHOFAR' and outlet='AWQAD OUTLET' 
and quarter=3 and EXTRACT('year' from sales_date)=2022

select count(*) from user_master_table where deleted IS NULL and isactive=true and region='DHOFAR' and outlet='AWQAD OUTLET' 
and (usr_position='PROMOTER' OR usr_position='CSR') and quarter=3 and year=2022

-- Manager --

select emp_id, name, usr_position, mobile_number, region, outlet, isAdmin from user_master_table where deleted IS NULL 
and emp_id='33220' and quarter=3 and year=2022

select distinct senior_manager from user_hierarchy_table where deleted IS NULL and manager='33220' and quarter=3 and year=2022

select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.outlet_manager AND a.quarter = b.quarter 
AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL 
AND b.manager='33220' AND a.quarter=3 AND a.year=2022 and a.isactive=true 

-- Senior Manager --

select * from user_master_table where deleted IS NULL 
and emp_id='36177' and quarter=3 and year=2022

select distinct general_manager from user_hierarchy_table where deleted IS NULL and senior_manager='36177' and quarter=3 and year=2022

select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.outlet_manager AND a.quarter = b.quarter 
AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL 
AND b.senior_manager='36177' AND a.quarter=3 AND a.year=2022

-- General Manager --
select * from user_master_table where deleted IS NULL 
and emp_id='12050' and quarter=3 and year=2023
select vice_president from user_hierarchy_table where deleted IS NULL and general_manager='12050' and quarter=3 and year=2023 limit 1
SELECT outlet_code, outlet_name, region_code, region_name, latitude, longitude FROM outlets_table where deleted IS NULL and 
region_name='MUSCAT' and outlet_name='MUSCAT' and quarter=3 and year=2023
select distinct a.emp_id from user_master_table a INNER JOIN user_hierarchy_table b on a.emp_id = b.outlet_manager AND a.quarter = b.quarter 
AND a.year = b.year WHERE a.deleted IS NULL AND b.deleted IS NULL AND b.general_manager='12050' AND a.quarter=3 AND a.year=2023

-- Query to fetch all users who logged in today --
select emp_id from user_authentication_table WHERE DATE(last_login_time) >= CURRENT_DATE AND DATE(last_login_time) < CURRENT_DATE + INTERVAL '1 DAY'