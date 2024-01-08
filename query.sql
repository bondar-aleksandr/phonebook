-- name: GetAllPersons :many
SELECT id, first_name, last_name, notes FROM person;

-- name: GetPersonByFname :many
SELECT id, first_name, last_name, notes FROM person WHERE first_name LIKE CONCAT('%', ?, '%');

-- name: GetPersonByLname :many
SELECT id, first_name, last_name, notes FROM person WHERE last_name LIKE ?;

-- name: GetPersonByFullname :many
SELECT id, first_name, last_name, notes FROM person WHERE first_name LIKE ? AND last_name LIKE ?;

-- name: GetPersonPhones :many
SELECT phone_number, phone_type FROM phone WHERE person_id = ?;

-- name: AddPerson :execlastid
INSERT INTO person (first_name, last_name, notes) VALUES (?, ?, ?);

-- name: AddPhone :exec
INSERT INTO phone (phone_number, phone_type, person_id) VALUES (?, ?, ?);