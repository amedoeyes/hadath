# Hadath

Hadath is a RESTful API designed for managing event creation, user authentication, and booking functionalities.

## Getting Started

### Perquisites

- Go
- PostgreSQL
- `.env` file with the appropriate configurations. Refer to `.env.example` for the required variables.

### Running the server

To start the API server, use the following command:

```sh
go run ./cmd/api/main.go
```

## API Endpoints

### Authentication

Endpoints related to user authentication.

- `/auth`
  - `/signup` `POST`
    Creates a new a account.
    - **Request**:
    ```json
    {
      "name": "string",
      "email": "string",
      "password": "string"
    }
    ```
  - `/signin` `POST`
    Authenticates an existing user.
    - **Request**:
    ```json
    {
      "email": "string",
      "password": "string"
    }
    ```
    - **Response**:
    ```json
    {
      "id": "integer",
      "name": "string",
      "email": "string",
      "api_key": "string"
    }
    ```
  - `/signout` `POST`
    Invalidates the API key.

### Events

Endpoints related to the creation, retrieval, updating, and deletion of events.

- `/events`
  - `POST`
    Creates a new event.
    - **Request**:
    ```json
    {
      "name": "string",
      "description": "string",
      "address": "string",
      "start_time": "string (RFC 3339 format)",
      "end_time": "string (RFC 3339 format)"
    }
    ```
  - `GET`
    Retrieves all events.
    - **Response**:
    ```json
    [
      {
        "name": "string",
        "description": "string",
        "address": "string",
        "start_time": "string (RFC 3339 format)",
        "end_time": "string (RFC 3339 format)"
      }
    ]
    ```
  - `/{id}`
    - `GET`
      Retrieves the event with the given id.
      - **Response**:
    ```json
    {
      "name": "string",
      "description": "string",
      "address": "string",
      "start_time": "string (RFC 3339 format)",
      "end_time": "string (RFC 3339 format)"
    }
    ```
    - `PUT`
      Updates the event with the given id.
      - **Request**: Same as the creation request.
    - `DELETE`
      Deletes the event with given id.

### Bookings

Endpoints related to user event bookings

- `/bookings`
  - `/user` `GET`
    Retrieves all events booked by the currently signed-in user.
    - **Response**:
    ```json
    [
      {
        "name": "string",
        "description": "string",
        "address": "string",
        "start_time": "string (RFC 3339 format)",
        "end_time": "string (RFC 3339 format)"
      }
    ]
    ```
  - `/event/{id}`
    - `GET`
      Retrieves a list of users booked for the event with the given id.
      - **Response**:
      ```json
      {
        "id": "integer",
        "name": "string",
        "email": "string"
      }
      ```
    - `POST`
      Books the currently signed-in user for the event with the given id.
    - `DELETE`
      Cancels the booking of the currently signed-in user for the event with the given id.

## Testing the Server Endpoints

To test the server endpoints, you can use the bash scripts located in the `./script/api/` directory. Below are examples of how to perform common tasks like creating an account, signing in, and using authorized endpoints.

### Signing up

```sh
./script/api/auth/signup "user" "user@email.com" "password"
```

### Signing in

```sh
./script/api/auth/signin "user@email.com" "password"
```

### Setting Up Authorization

Once you have the API key from signing in, export it as an environment variable `AUTHORIZATION` to authorize future requests:

```sh
export AUTHORIZATION="Bearer your_api_key_here"
```

For example:

```sh
export AUTHORIZATION="Bearer ebb37aca066d77bf492d5a26c3c1b7edcb5f327b22591222456e737631017534"
```

### Using Authorized Endpoints

With the `AUTHORIZATION` token set, you can now access endpoints that require authentication. For example, to create an event:

```sh
./script/api/events/post "Name" "Description" "Address" "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" "$(date -u +"%Y-%m-%dT%H:%M:%SZ" -d "+2 hours")"
```

And to book an event:

```sh
./script/api/bookings/event/id/post "b2d7bd82-ed58-4ea3-ba8b-ee59fe7f4f9e"
```
