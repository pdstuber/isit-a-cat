FROM oven/bun as base
WORKDIR /app

FROM base AS install
ADD package.json .
ADD bun.lockb .
RUN bun install --frozen-lockfile

FROM base AS build
COPY --from=install /app/node_modules /app/node_modules
COPY . .
RUN NODE_OPTIONS=--openssl-legacy-provider bun run build

FROM pierrezemb/gostatic:latest
COPY --from=build /app/dist /srv/http
ENTRYPOINT ["/goStatic"]
CMD [ "-enable-health" ]