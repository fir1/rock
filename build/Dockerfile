FROM golang:1.18-alpine as app-builder

WORKDIR /app

ARG PERSONAL_ACCESS_TOKEN
ENV PERSONAL_ACCESS_TOKEN=$PERSONAL_ACCESS_TOKEN

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 go build -o /server-app cmd/*.go



FROM scratch
COPY --from=app-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=app-builder /server-app /server-app


EXPOSE 8080

CMD ["/server-app"]