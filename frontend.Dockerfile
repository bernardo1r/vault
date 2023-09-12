FROM alpine
COPY ./frontend /home/frontend
WORKDIR /home/frontend
RUN apk add python3 py3-pip
RUN pip install -r requirements.txt
COPY ./cert/frontend/cert.pem .
COPY ./cert/frontend/key.pem .
COPY ./cert/backend/cert.pem server.pem
