# Book search

## "Book search" - The first part of the "Home Library" project

A small utility for searching books on hard drives.
Information about the found books is stored in the database.

<br>
TODO:<br>
1. beutify code

<br>

## Config file:
if there is no configuration file (config.yaml), an example of such a file will be created<br>

example of config.yaml:

        paths:
            - ./
            - d:/
        ext:
            .djvu: true
            .pdf: true
        pgsql:
            psql_user: user
            psql_passqord: password
            psql_host: localhost
            psql_port: "5432"
            psql_database: sampledb

