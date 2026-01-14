# Backend Builder
FROM registry.access.redhat.com/ubi10/go-toolset:1.25 AS backend

WORKDIR /app
COPY backend/ /app/

RUN go build -o backend -ldflags="-s -w" ./cmd/backend
RUN go build -o demo -ldflags="-s -w" ./cmd/demo


# Frontend Builder
FROM registry.access.redhat.com/ubi10/nodejs-24:latest AS frontend

USER root
WORKDIR /app
COPY frontend/ /app/

RUN npm install && npm run build


# Final image
FROM registry.access.redhat.com/ubi10/ubi-minimal:latest

WORKDIR /app

COPY --from=backend /app/backend /app/backend
COPY --from=backend /app/demo /app/demo
COPY --from=frontend /app/dist/demo/browser /app/public/

VOLUME /data

ENTRYPOINT ["/app/backend"]
CMD ["--dir", "/data", "serve", "--http", "0.0.0.0:8080" ]
