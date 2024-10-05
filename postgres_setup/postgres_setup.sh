#!/bin/bash

# Load environment variables from .env file if it exists
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

# Check if the container with the name 'postgres' is running
if [ ! "$(docker ps -q -f name=postgres)" ]; then
    # Check if the container exists but is in an exited (stopped) state
    if [ "$(docker ps -aq -f status=exited -f name=postgres)" ]; then
        # Start the container if it exists but is currently stopped
        docker start postgres
    else
        # Create and run a new container if it does not already exist
        docker run --name postgres -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD -e POSTGRES_DB=$POSTGRES_DB -d -p 5432:5432 postgres:latest
    fi

    # Wait for 10 seconds to ensure the container is fully up and running
    sleep 10
fi

# Create the 'users' table in the 'public' schema if it does not already exist
docker exec -it postgres psql -U postgres -d $POSTGRES_DB -c "
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users') THEN
        CREATE TABLE public.users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            email VARCHAR(100) UNIQUE,
            password VARCHAR(255) NOT NULL,
            role_id INTEGER,
            deleted_at TIMESTAMPTZ
        );
    END IF;
END
\$\$;
"

# Execute SQL command to insert the admin user into the users table
# Only insert if the user does not already exist
docker exec -it postgres psql -U postgres -d $POSTGRES_DB -c "
INSERT INTO public.users (username, email, password, role_id)
SELECT '$ADMIN_USERNAME', '$ADMIN_EMAIL', '$ADMIN_PASSWORD', '7'
WHERE NOT EXISTS (
    SELECT 1 FROM public.users WHERE username = '$ADMIN_USERNAME'
);
"
