-- name: GetAllPersons :many
SELECT * FROM person;

-- name: GetPersonByID :one
SELECT * FROM person WHERE id = ?;

-- name: GetPersonByFname :many
SELECT * FROM person WHERE first_name LIKE CONCAT('%', ?, '%');

-- name: GetPersonByLname :many
SELECT * FROM person WHERE last_name LIKE CONCAT('%', ?, '%');

-- name: GetPersonByFullname :many
SELECT * FROM person WHERE first_name LIKE CONCAT('%', ?, '%') AND last_name LIKE CONCAT('%', ?, '%');

-- name: GetPersonPhones :many
SELECT phone_number, phone_type FROM phone WHERE person_id = ?;

-- name: GetPersonIDByPhone :many
SELECT person_id FROM phone WHERE phone_number LIKE CONCAT('%', ?, '%');

-- name: AddPerson :execlastid
INSERT INTO person (first_name, last_name, notes) VALUES (?, ?, ?);

-- name: AddPhone :exec
INSERT INTO phone (phone_number, phone_type, person_id) VALUES (?, ?, ?);

-- name: DeletePersonByFname :execrows
DELETE FROM person WHERE first_name LIKE CONCAT('%', ?, '%');

-- name: DeletePhoneByNumber :execrows
DELETE FROM phone WHERE phone_number LIKE CONCAT('%', ?, '%');