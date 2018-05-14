SELECT * FROM go_gift_record WHERE rev_user=100066


SELECT SUM(VALUE) FROM go_gift_record WHERE rev_user=100066

SELECT * FROM go_user WHERE uid=100066

 
 SELECT SUM(num) FROM go_gift_assigned_detail WHERE gift_record_id IN (SELECT id FROM go_gift_record WHERE rev_user=100066) AND identity =1
 
 SELECT *  FROM go_gift_assigned_detail WHERE gift_record_id IN (SELECT id FROM go_gift_record WHERE rev_user=100066) AND identity =1
 
 SELECT    id, FROM_UNIXTIME(record_time,'%Y年%m月%d') ,`value` FROM go_gift_record  WHERE rev_user=100066
 
  SELECT   SUM(`value`)  FROM go_gift_record  WHERE rev_user=100066
  
    SELECT  id, `value`  FROM go_gift_record  WHERE rev_user=100066
    
   SELECT * FROM go_gift_assigned_detail WHERE gift_record_id IN (32075,32076,32077)
   
   