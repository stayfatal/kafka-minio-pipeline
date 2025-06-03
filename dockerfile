FROM golang:1.24.3-alpine AS mod

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./pkg ./pkg

COPY ./domain ./domain

FROM mod AS rest_builder

WORKDIR /app

COPY ./rest ./rest

RUN go build -o app rest/cmd/app/main.go

FROM alpine AS rest

WORKDIR /app

COPY --from=rest_builder app .

ENTRYPOINT [ "./app" ]


FROM mod AS pipe_builder

WORKDIR /app

COPY ./pipe ./pipe

RUN go build -o app pipe/cmd/app/main.go

FROM alpine AS pipe

WORKDIR /app

COPY --from=pipe_builder app .

ENTRYPOINT [ "./app" ]