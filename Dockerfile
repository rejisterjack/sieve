FROM golang:1.27-bookworm AS build
ENV GOTOOLCHAIN=auto
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/sieve ./cmd/sieve

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/sieve /sieve
COPY config.demo.yaml /config.demo.yaml
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/sieve", "serve", "-config", "/config.demo.yaml"]
