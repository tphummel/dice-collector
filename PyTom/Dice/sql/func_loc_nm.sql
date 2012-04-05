DELIMITER $$

DROP FUNCTION IF EXISTS loc_nm$$
CREATE FUNCTION loc_nm(in_locationid INT)
	RETURNS VARCHAR(50)
BEGIN
	DECLARE v_location VARCHAR(50);
	SElECT name
	INTO v_location
	FROM location
	WHERE id = in_locationid;
	RETURN v_location;
END$$

DELIMITER ;
