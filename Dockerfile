FROM golang:1.13 as base

ENV APP /go/src
WORKDIR $APP

ADD https://github.com/golang-migrate/migrate/releases/download/v4.10.0/migrate.linux-amd64.tar.gz /bin/migrate.tar.gz
RUN tar -xzf /bin/migrate.tar.gz -C /bin/

COPY go.mod $APP
COPY go.sum $APP

RUN go mod download

COPY . $APP

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app

FROM base as test
# RUN go get -u github.com/stretchr/testify

# FROM scratch as release
FROM alpine as release

COPY --from=base /go/src/start.sh /app/
COPY --from=base /go/bin/app /app/
COPY --from=base /bin/migrate.linux-amd64 /app/
COPY --from=base /go/src/db/migrations /app/db/migrations/

CMD [ "bin/sh", "/app/start.sh" ]