CREATE TABLE credentials(
    user_id int8,
    service text,
    login text,
    password text,

    PRIMARY KEY (user_id, service)
);