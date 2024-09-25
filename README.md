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

- `/auth/login` - Log in with basic auth

#### OAuth2

- `/auth/google_login` - Log in with Google
- `/auth/google_callback` - Google login callback
- `/auth/github_login` - Log in with GitHub
- `/auth/github_callback` - GitHub login callback

### User Management

- `/users`
  - Create users
  - Edit users
  - Retrieve all users (paginated)
  - Retrieve user details
  - Delete users

### Kubernetes Resources

#### Events

- `/events` - List events
- `/events/:id` - Get event details

#### Pods

- `/pods`
  - Create pods
  - Edit pods
  - Retrieve all pods (paginated)
  - Retrieve pod details
  - Delete pods

#### Deployments

- `/deployments`
  - Create deployments
  - Edit deployments
  - Retrieve all deployments (paginated)
  - Retrieve deployment details
  - Delete deployments

#### Namespaces

- `/namespaces`
  - Create namespaces
  - Edit namespaces
  - Retrieve all namespaces (paginated)
  - Retrieve namespace details
  - Delete namespaces

## Swagger Documentation

You can access the Swagger documentation at the following URL:

```sh
/swagger/index.html
```

## Installation

### Prerequisites

- Docker
- Go (version 1.22.3 or later)
- kubectl configured with access to your Kubernetes cluster

### Postgres Setup

To set up the PostgreSQL database for the Kubernetes API, follow these steps:

1. **Create Environment File**

   Copy the example environment file to create your .env file:

   ```sh
   cp .env_example .env
   ```

   Ensure that you update the .env file with your specific environment variables.

2. **Run PostgreSQL Setup Script**

   Execute the postgres_setup.sh script to initialize the PostgreSQL database:

   ```sh
   ./postgres_setup.sh
   ```

   This script will:
   - Set up the PostgreSQL container if it's not already running.
   - Create the necessary tables in the database.

### Building and Running the API

1. Clone the repository:

   ```sh
   git clone https://github.com/fleimkeipa/kubernetes-api.git
   cd kubernetes-api
   ```

2. Install dependencies:

   ```sh
   go mod download
   ```

3. Build the application:

   ```sh
   go build
   ```

4. Run the application:

   ```sh
   ./kubernetes-api
   ```

## Usage

Here are some example API calls using curl:

1. Login (Basic Auth):

   ```sh
   curl -X POST http://localhost:8080/auth/login -u username:password
   ```

## Contributing

We welcome contributions to the kubernetes-api project! Here's how you can contribute:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please make sure to update tests as appropriate and adhere to the project's coding standards.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Support

If you encounter any problems or have any questions, please open an issue in the GitHub repository.
