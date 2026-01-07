# Backend Builder
FROM registry.access.redhat.com/ubi10/go-toolset:1.25 AS backend

WORKDIR /app
COPY backend/ /app/

RUN go build -o backend -ldflags="-s -w" ./cmd/backend


# Frontend Builder
FROM registry.access.redhat.com/ubi10/nodejs-24:latest AS frontend

WORKDIR /app
COPY frontend/ /app/

RUN npm install && npm run build


# Final image
FROM registry.access.redhat.com/ubi10:latest

WORKDIR /app

COPY --from=backend /app/backend /app/backend
COPY --from=frontend /app/dist/demo/ /app/public/

VOLUME /data

ENTRYPOINT ["/app/backend"]
CMD ["--dir", "/data", "--http", "0.0.0.0:8080" ]
