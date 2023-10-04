<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#setup">Setup</a></li>
        <li><a href="#testing">Testing</a></li>
      </ul>
    </li>
    <li>
      <a href="#design">Design</a>
      <ul>
        <li><a href="#database">Database</a></li>
        <li><a href="#api">API</a></li>
      </ul>
    </li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

## About the Project
A REST API for an Authorization+Authentication service in Golang built for the Houseware Backend Octernship Assignment

### Built With

* [Gin](https://github.com/gin-gonic/gin) - Gin is the Golang Web framework used for building the API. It was chosen because it is the most popular web framework for Go and is easy to use, has high performance and a robust middleware system that allows developers to easily add functionality to their APIs.
* [PostgreSQL](https://www.postgresql.org) - I decided to use a relational database for the project because there is potential for joins and relationships involving users in the database. PostgreSQL was chosen as it is a feature-rich database that is open-source and is highly scalable.
* [GORM](https://github.com/go-gorm/gorm) - It is the ORM used to interface with Postgres.  GORM was chosen as it is a developer-friendly, feature-rich ORM built on the `database/sql` package, possessing many functionalities like auto migration, logging, contexts, prepared statements, associations, constraints etc. The GORM package takes the code-first approach and uses structs as the model for interacting with databases which makes it very developer friendly. Using GORM doesn’t trade off most of the functionalities you’ll get from writing raw SQL queries.

## Getting Started
### Prerequisites
*  Install [Go](https://go.dev/doc/install)
*  A PostgreSQL database has to be setup locally. This can be done using Docker.
	* Install [Docker](https://docs.docker.com/get-docker/) 
	* Run 
		```sh
		docker run --name postgres-container -e POSTGRES_PASSWORD=<PASSWORD> -e POSTGRES_USER=<USER> -e POSTGRES_DB=<DATABASE> -p <HOST_PORT>:5432 -d postgres
		```

### Setup
* Fill out the`.env` file
	
	| Variable             | Description                                                 | Default        |
	| -------------------- | ----------------------------------------------------------- | -------------- |
	| SERVER_PORT          | Port on which the API server should run                     |
	| CLIENT_DOMAIN        | Domain on which the refresh cookie has to be sent           | localhost      |
	| DB_HOST              | Host where the database is running                          | localhost      |
	| DB_PORT              | Port on which the database is running                       |                |
	| DB_PASSWORD          | Password for the database user                              |                |
	| DB_USER              | User for the database                                       |                |
	| DB_SSLMODE           | SSL mode for the database                                   | disable        |
	| DB_NAME              | Name of the databse                                         |                |
	| ACCESS_TOKEN_SECRET  | Secret to sign access token                                 |                |
	| REFRESH_TOKEN_SECRET | Secret to sign refresh token                                |                |
	
* In the root folder run
	```sh
	cd backend; go run cmd/main.go
	```
	
### Testing
Unit tests have been written for all the handlers. Use a different database to run the tests to prevent conflicts with existing data as well as to prevent deleting data after a test during clean up.

* Fill out the `.env.testing` file
* In the root folder run
	```sh
	cd backend; go test ./...
	```
## Design
### Database 

Users of different organizations are isolated using PostgreSQL schemas. PostgreSQL schemas let us hold multiple instances of the same set of tables inside a single database. They’re essentially **namespaces for tables**. The following is diagram of the database structure -

<img width="935" alt="image" src="https://user-images.githubusercontent.com/76054921/227722806-da1b2fd4-9dca-460f-ae1c-c7f7862ca942.png">


Multi-tenant solutions range from one database per tenant (shared nothing) to one row per tenant (shared everything).

"Shared nothing" = separate database per tenant, most expensive per client, highest data isolation, simple disaster recovery, theoretically harder maintenance, easily customizable, lowest number of rows per table.

"Shared everything" = shared table, least expensive per tenant, lowest data isolation, complicated disaster recovery, simpler structural maintenance, highest number of rows per table.

"Shared schema" = tenants share a database, each tenant has its own named schema, cost falls between "shared nothing" and "shared everything", better isolation than "shared everything", easier maintenance than "shared nothing", more active tenants per server, disaster recovery for a single tenant is easy or hard depending on DBMS.

In the assignment, each of the organizations have their own schemas with isolated User tables. This approach was chosen as it provided a balance between isolation and cost/maintenance.

### API
Postman spec for the API can be found [here](https://www.postman.com/rafikun/workspace/public/collection/18135506-6b30d9b6-7133-4fa4-a0fa-0fde1462a312?action=share&creator=18135506)

## Contact
Abdullah Rafi - abdullahrafi.1308@gmail.com
