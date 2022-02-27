FROM golang:1.17.7-alpine3.15 as base

ENV APP /go/src
WORKDIR $APP

RUN apk add -u build-base

ADD https://github.com/golang-migrate/migrate/releases/download/v4.10.0/migrate.linux-amd64.tar.gz /bin/migrate.tar.gz
RUN tar -xzf /bin/migrate.tar.gz -C /bin/ && rm /bin/migrate.tar.gz

COPY go.mod $APP
COPY go.sum $APP

RUN go mod download

COPY . $APP

FROM base as build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app ./cmd/api/main.go

FROM alpine as release

COPY --from=build /go/src/start.sh /app/
COPY --from=build /go/bin/app /app/
COPY --from=build /bin/migrate.linux-amd64 /app/
COPY --from=build /go/src/db/migrations /app/db/migrations/

CMD [ "bin/sh", "/app/start.sh" ]