CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS subscriptions 
(
	subscription_id SERIAL PRIMARY KEY,
	service_name VARCHAR(32) NOT NULL,
	price INTEGER NOT NULL,
	user_id UUID NOT NULL,
	start_date VARCHAR(10) NOT NULL CHECK (
        start_date ~ '^(0[1-9]|1[0-2])-[2-9][0-9]{3}$'
        AND TO_DATE(start_date, 'MM-YYYY') IS NOT NULL
 ),
	end_date VARCHAR(10) CHECK (
        start_date ~ '^(0[1-9]|1[0-2])-[2-9][0-9]{3}$'
        AND TO_DATE(start_date, 'MM-YYYY') IS NOT NULL
    )
);
