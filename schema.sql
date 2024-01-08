CREATE TABLE IF NOT EXISTS person (
    id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(50) NOT NULL,
	last_name VARCHAR(50) NOT NULL,
    notes VARCHAR(200),
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    modified TIMESTAMP ON UPDATE CURRENT_TIMESTAMP()
);

CREATE TABLE IF NOT EXISTS phone (
    id INT PRIMARY KEY AUTO_INCREMENT,
    phone_number VARCHAR(20) UNIQUE NOT NULL, 
	phone_type SMALLINT NOT NULL DEFAULT 3,
    person_id INT NOT NULL,
    CONSTRAINT phone_owner FOREIGN KEY (person_id) REFERENCES person(id) ON DELETE CASCADE
);