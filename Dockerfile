FROM golang:1.24.3

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod tidy 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/go_final_project

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/go_final_project ./go_final_project
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=12345

EXPOSE ${TODO_PORT}

CMD ["./go_final_project"]
