# syntax=docker/dockerfile:1

FROM golang:1.27-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/openlane-mcp ./cmd/openlane-mcp

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/openlane-mcp /openlane-mcp
USER nonroot:nonroot
ENTRYPOINT ["/openlane-mcp"]
CMD ["serve"]
