# Muzz Backend Technical

## Highlights

- **Graceful Shutdown**: The application gracefully handles shutdowns to ensure no data is lost and ongoing processes are completed.
- **Three Tier Layered Approach**:
    - **HTTP Layer**: Handles incoming HTTP requests and sends responses.
    - **Service Layer**: Contains business logic.
    - **Repository Layer**: Manages data persistence with MySQL and Elasticsearch.
- **Used MySQL**: Stores users, swipes, and matches.
- **Used Elasticsearch**: Facilitates user discovery.

## Developer Experience

- **Make Commands**: Simplifies common tasks such as imports, formatting, linting, and migrations.
- **Easy Migration Creation**: Create database migrations effortlessly.

## How to Run

1. **Setup Environment**:
    - Ensure Docker is installed and running.

2. **Run the Project**:
   ```sh
   make up
   ```
   
- Create a user
   ```sh
   curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{}'
   ```
- Login
```sh
curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "queen-ethelyn-beatty@muzz.com",
    "password" : "Password1"
}'
```
- discover
```sh
curl --location 'http://localhost:8080/discover?lat=10.0&lon=10.0&min_age=1&gender=male' \
--header 'x-pf-app: my_crm' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InF1ZWVuLWV0aGVseW4tYmVhdHR5QG11enouY29tIiwiZXhwIjoxNzE4ODM3ODMwLCJ1c2VyX2lkIjo2Mn0.tETHFLadyDuaCEwMOQ-8SOoIWk27IUYXZG5dZHKpfX8'
```
- swipe
```sh
curl --location 'http://localhost:8080/swipe' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InF1ZWVuLWV0aGVseW4tYmVhdHR5QG11enouY29tIiwiZXhwIjoxNzE4ODM3ODMwLCJ1c2VyX2lkIjo2Mn0.tETHFLadyDuaCEwMOQ-8SOoIWk27IUYXZG5dZHKpfX8' \
--data '{
    "user_id": 62,
    "target_id": 61,
    "preference": "yes"
}'
```

