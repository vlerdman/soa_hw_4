# SOA_HW_4

Simple REST API service with JWT auth, RabbitMQ and Postgres.

## API + Examples

### POST /users/register

```bash
curl 0.0.0.0:8080/users/register -X POST -d \
'{
  "username":"Bob", 
  "password": "secret", 
  "sex": "male", 
  "email": "bob_best@gmail.com", 
  "avatar":"https://img.freepik.com/free-photo/red-white-cat-i-white-studio_155003-13189.jpg?w=1060&t=st=1685911204~exp=1685911804~hmac=41b2d86724e961bd540888757d3c1bf2c2e0f050230f82c3ca36d8aecc44cd94"
}'
```

```bash
{
    "id":"fe5cd01f-92de-4671-b720-18c93e8fd0c4",
    "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZmU1Y2QwMWYtOTJkZS00NjcxLWI3MjAtMThjOTNlOGZkMGM0IiwiZXhwIjoxNjg1OTk0NDY5fQ.VAtm8tr0IJDzq9gR17V6wXzNqWavcIxLdCZdWIAb87w"
}
```

Responce fields: id (uuid of registered user), token (auth token with 1 hour TTL).

### POST /users/auth

```bash
curl 0.0.0.0:8080/users/auth -X POST -d \
'{
    "username": "Bob",
    "password": "secret"
}'
```

```bash
{
    "id":"fe5cd01f-92de-4671-b720-18c93e8fd0c4",
    "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZmU1Y2QwMWYtOTJkZS00NjcxLWI3MjAtMThjOTNlOGZkMGM0IiwiZXhwIjoxNjg1OTk0NTY2fQ.lkwQb2AbuQhgyrQlB0Iwi2vf7EPt7yUfdfZoCPvXDOc"
}
```

### GET /users?usernames=

```bash
curl "0.0.0.0:8080/users?usernames=Alice,Bob" -H \
"Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZmU1Y2QwMWYtOTJkZS00NjcxLWI3MjAtMThjOTNlOGZkMGM0IiwiZXhwIjoxNjg1OTk0NTY2fQ.lkwQb2AbuQhgyrQlB0Iwi2vf7EPt7yUfdfZoCPvXDOc"
```

```bash
{
    "users":[
        {
            "id":"4c40a729-6cf8-47c3-89c2-62323ae05c44",
            "registration_time":"2023-06-05T18:47:38.764342Z",
            "password":"******",
            "username":"Alice",
            "avatar":"https://i.natgeofe.com/n/548467d8-c5f1-4551-9f58-6817a8d2c45e/NationalGeographic_2572187_square.jpg",
            "sex":"female",
            "email":"alice_best@gmail.com"
        },
        {
            "id":"fe5cd01f-92de-4671-b720-18c93e8fd0c4",
            "registration_time":"2023-06-05T18:47:49.639602Z",
            "password":"******",
            "username":"Bob",
            "avatar":"https://img.freepik.com/free-photo/red-white-cat-i-white-studio_155003-13189.jpg?w=1060\u0026t=st=1685911204~exp=1685911804~hmac=41b2d86724e961bd540888757d3c1bf2c2e0f050230f82c3ca36d8aecc44cd94",
            "sex":"male",
            "email":"bob_best@gmail.com"
        }
        ]}
```

### POST /users/edit

```bash
curl 0.0.0.0:8080/users/edit -X POST -H \
"Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZmU1Y2QwMWYtOTJkZS00NjcxLWI3MjAtMThjOTNlOGZkMGM0IiwiZXhwIjoxNjg1OTk0NTY2fQ.lkwQb2AbuQhgyrQlB0Iwi2vf7EPt7yUfdfZoCPvXDOc" -d \
'{
    "username":"Bobby", 
    "password": "secret", 
    "sex": "male", 
    "email": "bobby_best@gmail.com", 
    "avatar":"https://img.freepik.com/free-photo/red-white-cat-i-white-studio_155003-13189.jpg?w=1060&t=st=1685911204~exp=1685911804~hmac=41b2d86724e961bd540888757d3c1bf2c2e0f050230f82c3ca36d8aecc44cd94
}'
```

### POST /users/stats

Starts async task of creating user resume (user is determined by auth token).

```bash
curl 0.0.0.0:8080/users/stats -X POST -H \
"Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZmU1Y2QwMWYtOTJkZS00NjcxLWI3MjAtMThjOTNlOGZkMGM0IiwiZXhwIjoxNjg1OTk0NTY2fQ.lkwQb2AbuQhgyrQlB0Iwi2vf7EPt7yUfdfZoCPvXDOc"
```

```bash
{
    "id":"59ebccf1-a6ac-4634-8c95-7090e5df0f59"
}
```

Id in responce - id of task for getting result.

### GET /users/stats/{id}

```bash
curl 0.0.0.0:8080/users/stats/59ebccf1-a6ac-4634-8c95-7090e5df0f59 -H \
"Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZmU1Y2QwMWYtOTJkZS00NjcxLWI3MjAtMThjOTNlOGZkMGM0IiwiZXhwIjoxNjg1OTk0NTY2fQ.lkwQb2AbuQhgyrQlB0Iwi2vf7EPt7yUfdfZoCPvXDOc" --output result.pdf
```

Example of generated [pdf](examples/result.pdf)