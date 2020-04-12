FROM golang:1.13 as base

ADD https://github.com/golang-migrate/migrate/releases/download/v4.10.0/migrate.linux-amd64.tar.gz /migrate.tar.gz
RUN tar -xzf /migrate.tar.gz -C /

ENV APP /go/src
WORKDIR $APP

COPY go.mod $APP
COPY go.sum $APP

RUN go mod download

COPY . $APP

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app

FROM base as test
# RUN go get -u github.com/stretchr/testify

FROM scratch as release
COPY --from=base /go/bin/app /
COPY --from=base /migrate.linux-amd64 /migrate
CMD [ "./app" ]