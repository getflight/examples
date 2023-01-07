CREATE TABLE notes
(
    id         INT           NOT NULL,
    name       VARCHAR(256)  NOT NULL,
    note       VARCHAR(4000) NOT NULL,
    created_at TIMESTAMP     NOT NULL,
    updated_at TIMESTAMP     NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
);