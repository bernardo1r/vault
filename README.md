# Running the project locally
For running the project locally you can use [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/install/).
But first it is needed to generate custom certificates for the client and server.

## Generating self-signed certificates
Run the command below to generate the server certificate:
```bash
openssl req -newkey rsa:4096 -x509 -sha256 -nodes -out cert/backend/cert.pem -keyout cert/backend/key.pem
```

And to generate the client certificate:
```bash
openssl req -newkey rsa:4096 -x509 -sha256 -nodes -out cert/frontend/cert.pem -keyout cert/frontend/key.pem
```

## Running locally
Run Docker Compose:
```bash
docker-compose up
```
To open the application just open the browser on page `localhost:5000`
