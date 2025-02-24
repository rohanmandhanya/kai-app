# **Vulnerability Scanner API**

This project provides a RESTful API for scanning and querying vulnerabilities from GitHub repositories. The API supports two main operations:

1. **Scan** - Fetches JSON files from a specified GitHub repository, processes the vulnerabilities contained in these files, and stores them in an SQLite database.
2. **Query** - Allows querying the stored vulnerabilities by specific filters (e.g., severity).

---

## **Prerequisites**

- **Docker** and **Docker Compose** installed on your machine.
- Go 1.16 or higher (for development).
- SQLite (for database).
- GitHub API access (for fetching repository data).

---

## **Installation**

1. **Clone the repository:**

`git clone https://github.com/yourusername/vulnerability-scanner-api.git`
`cd vulnerability-scanner-api`

2. **Install dependencies:**

`go mod tidy`

## **Running the Application**

- **Docker**
To run the application using Docker Compose
`docker compose up`

- **Manually**
To run the application manually
`go run cmd/main.go`

## **Testing**
This application includes unit and integration tests to ensure that the logic and API endpoints function correctly.
`go test ./,, -v`
