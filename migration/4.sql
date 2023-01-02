-- /sales/incentive/client --

select emp_id, name, usr_position, mobile_number, region, outlet, isAdmin from user_master_table where deleted IS NULL and emp_id='69952' and quarter=3 and year=2022
 -- 69952	EMAN OBIED AL WAHIBI	OUTLET MANAGER	99111366	SHARQIYA	MASEERAH OUTLET	false --

 select distinct emp_id from user_hierarchy_table where deleted IS NULL and outlet_manager='69952' and quarter=3 and year=2022;
 -- [23101853 36204 10473865 10724657 23477481 69485] --