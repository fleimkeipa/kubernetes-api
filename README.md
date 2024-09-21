# kubernetes-api

This project provides a Kubernetes API wrapper with authentication and various Kubernetes resource management endpoints.

## Features

- Authentication (Basic Auth, Google, and GitHub)
- User management
- Kubernetes resource management (Pods, Deployments, Namespaces)
- Event viewing
- Swagger documentation

## Technologies Used

- **Web Framework**: Echo v4
- **Logging**: ZapLogger
- **Kubernetes Client**: k8s.io/client-go
- **Database**: PostgreSQL
- **Configuration**: yaml and viper for environment mapping

## API Endpoints

### Authentication

#### Basic Auth

Log in with basic auth

### OAuth2

Log in with Google
Google login callback
Log in with GitHub
GitHub login callback

### User Management

- Create users
- Edit users
- Retrieve all users (paginated)
- Retrieve user details
- Delete users

### Kubernetes Resources

#### Events

List events
Get event details

#### Pods

- Create pods
- Edit pods
- Retrieve all pods (paginated)
- Retrieve pod details
- Delete pods

#### Deployments

- Create deployments
- Edit deployments
- Retrieve all deployments (paginated)
- Retrieve deployment details
- Delete deployments

#### Namespaces

- Create namespaces
- Edit namespaces
- Retrieve all namespaces (paginated)
- Retrieve namespace details
- Delete namespaces

## Swagger endpoint

You can access the Swagger documentation at the following URL:

```sh
/swagger/index.html
```

## Postgres Setup

To set up the PostgreSQL database for the Kubernetes API, follow these steps:

### Create Environment File

Copy the example environment file to create your .env file:

```sh
cp .env_example .env
```

Ensure that you update the .env file with your specific environment variables.

### Run PostgreSQL Setup Script

Execute the postgres_setup.sh script to initialize the PostgreSQL database:

```sh
./postgres_setup.sh
````

This script will:

Set up the PostgreSQL container if it's not already running.
Create the necessary tables in the database.
