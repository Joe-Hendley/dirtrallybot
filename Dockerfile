FROM golang:1.27-alpine AS build

WORKDIR /src

# Warm the module cache before copying the source so dependency downloads are
# cached independently of code changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# The templ views are committed, so a plain build is enough. CGO is disabled so
# the binary is fully static and can run on the distroless base below.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" \
    -o /out/dirtrallybot ./cmd/dirtrallybot

# godotenv treats a missing .env as a fatal error, so ship an empty one and let
# every real value arrive through the process environment. /data holds the bbolt
# database and is the intended volume mount point, so hand it to the distroless
# runtime user (uid 65532).
RUN mkdir -p /out/app /out/data \
    && touch /out/app/.env \
    && chown -R 65532:65532 /out/data

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/dirtrallybot /usr/local/bin/dirtrallybot
COPY --from=build --chown=65532:65532 /out/app /app
COPY --from=build --chown=65532:65532 /out/data /data

WORKDIR /app

# Keep the database on the /data volume rather than next to the binary.
ENV DBPATH=/data/rallybot.db

# Challenge viewer. Override WEBADDR to change the port or set it to "off" to
# disable the viewer.
EXPOSE 8080

ENTRYPOINT ["dirtrallybot"]
