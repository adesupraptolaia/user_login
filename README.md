# User Login

There are 2 apps here, auth-app and user-app.
Username and password are saved to auth-app in table `users` and User's profile saved to user-app in table `user_profiles`.
User-app and Auth-app save data to mysql.

## Port

User-app running on port 8000 and Auth-app running on 9000 (public) and 9001 (private).
Port 9001 (actually) will not expose to external, it can only access internally by user-app.

## Flow

- Login User Flow
  ![create_new_user_flow](https://github.com/adesupraptolaia/user_login/blob/main/assert/login_flow.png)

- Create New User Flow
  ![create_new_user_flow](https://github.com/adesupraptolaia/user_login/blob/main/assert/create_new_user_flow.png)

- Delete User Flow
  ![create_new_user_flow](https://github.com/adesupraptolaia/user_login/blob/main/assert/delete_user_flow.png)

## How to running

You can use `docker-compose` or running it manually.

- using docker-compose

  ```
  make run
  ```

- manually
  make sure you have mysql, and set database config in `./config/congfig.yml`
  OR you can using docker to start mysql

  ```
  docker run --name mysql-8.0 \
      -e MYSQL_ROOT_PASSWORD=root \
      -e MYSQL_DATABASE=user_login \
      -p 3306:3306 \
      -d mysql:8.0
  ```

  after you have mysql, run these command.

  ```
  go mod tidy
  go run main.go auth

  ### open new terminal ###
  go run main.go user
  ```

## Swagger

You can access the Swagger after running the app.
user-app `http://localhost:8000/swagger/index.html`
and auth-app `http://localhost:9000/swagger/index.html`

## Postman

for testing purpose, use postman collection by importing file [postman.json](https://github.com/adesupraptolaia/user_login/blob/main/assert/postman.json) to your postman app
