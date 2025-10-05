FROM docker.io/node:lts-alpine AS frontend

WORKDIR /webapp

COPY webapp/package.json webapp/package-lock.json ./
RUN --mount=type=cache,target=/webapp/.npm \
    npm set cache /webapp/.npm && \
    npm ci

COPY webapp .
RUN npm run build

FROM docker.io/golang:1.25-alpine AS backend

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
COPY --from=frontend /webapp/dist ./webapp/dist

RUN --mount=type=cache,target=/go/pkg/mod go generate && go build -o /app/recipe-manager .

FROM scratch

COPY --from=backend /app/recipe-manager /bin/recipe-manager
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Atlas requires a tmp directory
COPY --from=backend --chmod=1777 /tmp /tmp
# Needed to apply migrations
COPY --from=docker.io/arigaio/atlas /atlas /bin/atlas

ENV PATH="$PATH:/bin"

ENTRYPOINT ["recipe-manager"]