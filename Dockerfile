FROM docker.io/golang:1.26.2 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build

FROM gcr.io/distroless/static-debian12
COPY --from=build /go/src/app/token-swap-discourse-jwt /

EXPOSE 3400
CMD ["/token-swap-discourse-jwt"]
