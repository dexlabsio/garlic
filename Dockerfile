# garlic/Dockerfile
FROM golang:1.24 AS garlic-source

WORKDIR /garlic
COPY . .
