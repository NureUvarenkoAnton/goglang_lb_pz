FROM golang:1.22.2-alpine3.19 as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o ./mig ./cmd/migrations/main.go
RUN go build -o ./app ./cmd/app/main.go

FROM alpine
COPY --from=builder /app/app /home/app
COPY --from=builder /app/mig /home/mig
COPY --from=builder /app/migrations /home/migrations
CMD /home/mig up ; /home/app 

