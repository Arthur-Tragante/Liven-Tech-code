# Liven Tech code test

## Getting Started

Follow the steps below to set up and run the project.

### Clone the Repository

```bash
git clone https://github.com/Arthur-Tragante/liven-tech-code-test.git
```

### Create Environment File
Navigate to the backend folder 

```
cd backend
```
create an .env file with the following content (you can change the values if you want):

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=liventechuser
DB_PASSWORD=liventechpassword
DB_NAME=liventechdatabase
JWT_SECRET=jwtsecret
POSTGRES_USER=liventechuser
POSTGRES_PASSWORD=liventechpassword
POSTGRES_DB=liventechdatabase
```

### Install Docker Desktop

Make sure Docker Desktop is installed and running on your machine.

Start the backend services using Docker:

```
docker-compose up -d
```
Run the Go application:

```
go run main.go
```
To run the backend tests:
```
go test ./...
```
### Frontend Setup

Navigate to the frontend folder:
```
cd frontend
```
Install the required dependencies:
```
npm install
```
Start the frontend application:
```
npm start
```
To run the frontend tests:
```
npm test
```
