# Barista (Backend)

A backend for the blogging app Barista.

## Usage

1. Install project-external dependencies

- [Go](https://go.dev/) (programming language)
- [Goose](https://github.com/pressly/goose/) (database migration tool)
- [SQLite](https://sqlite.org/) (database driver)

2. Clone the project locally and navigate to project folder

```sh
git clone https://github.com/LucDeCaf/barista-backend
cd barista-backend
```

3. Install project dependencies

```sh
go mod download
```

4. Create an empty database using Goose
```sh
cd migrations
goose sqlite3 ../db.db up
```

5. Run the project locally
```sh
cd ..
go run main.go
```

The app should now be running locally on port 8080.

```sh
<< curl localhost:8080/v1/blogs
>> []
```

## Specifications
- Go `net/http` package
- Goose SQL migration management tool
- SQLite database

## How to (roughly) use the API

This API is not expected to be exposed directly to the internet but is instead meant to run as an internal
process on a server which is _also_ running a separate frontend that will be exposed to the internet.

| Endpoint | Description | Response | Required headers |
| --- | --- | --- | --- |
| GET `/v1/blogs` | Get all blogs in JSON format | Blog array (JSON) ||
| POST `/v1/blogs` | Create a new blog | Blog (JSON) | `Authorization` |
| GET `/v1/blogs/{id}` | Get blog with the specified ID | Blog (JSON) ||
| DELETE `/v1/blogs/{id}` | Delete blog with the specified ID | Blog (JSON) | `Authorization` |
| POST `/v1/users` | Perform an action involving users | Varies | `Server-Action`, possibly `Authorization` |
| POST `/login` | Login the user as specified in request body | Authorization JWT ||
| POST `/register` | Register a new user as specified in request body | User (JSON) ||
