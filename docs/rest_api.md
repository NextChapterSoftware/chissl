## User Management REST API

This REST API provides endpoints to manage users, including authentication, retrieval, creation, updating, and deletion of user records. It uses Basic Authentication for securing access to the endpoints.

### Known limitations and TODOs
- Use client side hashed passwords instead of plaintext
 
### Endpoints

#### Authentication Middleware

* **Basic Auth Middleware**: Validates the username and password from the `Authorization` header in each request.

#### User Endpoints

* **Get All Users**
    * **Endpoint**: `GET /users`
    * **Description**: Retrieves a list of all users.
    * **Response**: JSON array of users.

* **Get User by Username**
    * **Endpoint**: `GET /user/{username}`
    * **Description**: Retrieves details of a specific user by username.
    * **Response**: JSON object of the user details.

* **Add New User**
    * **Endpoint**: `POST /user`
    * **Description**: Adds a new user.
    * **Request Body**: JSON object with user details.
    * **Response**: Status 201 Created on success.

* **Update Existing User**
    * **Endpoint**: `PUT /user`
    * **Description**: Updates details of an existing user.
    * **Request Body**: JSON object with updated user details.
    * **Response**: Status 202 Accepted on success.

* **Delete User by Username**
    * **Endpoint**: `DELETE /user/{username}`
    * **Description**: Deletes a user by username.
    * **Response**: Status 202 Accepted on success.

* **Upload Authfile**
    * **Endpoint**: `POST /authfile`
    * **Description**: Uploads a new authfile to reset user configurations.
    * **Request Body**: JSON array of users.
    * **Response**: Status 202 Accepted on success.

### Error Handling

* **400 Bad Request**: Returned for invalid request payloads or incorrect URL formats.
* **401 Unauthorized**: Returned for unauthorized access due to missing or invalid credentials.
* **404 Not Found**: Returned when a requested user is not found.
* **409 Conflict**: Returned when attempting to add a user that already exists.
* **500 Internal Server Error**: Returned for server-side errors.

### Example Usage

```sh
# Get all users
curl -u username:password -X GET http://localhost:8080/users

# Get a specific user
curl -u username:password -X GET http://localhost:8080/user/johndoe

# Add a new user
curl -u username:password -X POST http://localhost:8080/user -d '{"username": "janedoe", "password": "password123", "is_admin": true}'

# Update an existing user
curl -u username:password -X PUT http://localhost:8080/user -d '{"username": "janedoe", "password": "newpassword123", "is_admin": true}'

# Delete a user
curl -u username:password -X DELETE http://localhost:8080/user/janedoe

# Upload an auth file
curl -u username:password -X POST http://localhost:8080/authfile -d '[{"username": "admin", "password": "adminpass", "is_admin": true}, {"username": "user1", "password": "user1pass", "is_admin": false}]'
