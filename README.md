# kubernetes-api

## Swagger endpoint

```sh
http://localhost:8080/swagger/index.html
```

## Postgres Setup

Create .env with .env_example after run postgrest_setup.sh

Kubernetes API
Swagger Endpoint
You can access the Swagger documentation at the following URL:

```bash
http://localhost:8080/swagger/index.html
```

PostgreSQL Setup
To set up the PostgreSQL database for the Kubernetes API, follow these steps:

Create Environment File:

Copy the example environment file to create your .env file:

```sh
cp .env_example .env
```

Ensure that you update the .env file with your specific environment variables.

Run PostgreSQL Setup Script:

Execute the postgres_setup.sh script to initialize the PostgreSQL database:

```sh
./postgres_setup.sh
````

This script will:

Set up the PostgreSQL container if it's not already running.
Create the necessary tables in the database.
