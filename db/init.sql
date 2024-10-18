CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    birthday DATE NOT NULL,
    launchpad_id VARCHAR(50) NOT NULL,
    destination VARCHAR(50) NOT NULL,
    launch_date DATE NOT NULL
);
