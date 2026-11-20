# syntax=docker/dockerfile:1
FROM node:24-alpine AS build
WORKDIR /src
RUN npm install -g pnpm@12.21.0
COPY pnpm-workspace.yaml pnpm-lock.yaml ./
COPY elzzes-calculator-web/package.json elzzes-calculator-web/
COPY elzzes-ui-lib/package.json elzzes-ui-lib/
RUN pnpm install --frozen-lockfile
COPY elzzes-calculator-web elzzes-calculator-web
COPY elzzes-ui-lib elzzes-ui-lib
RUN pnpm --filter elzzes-calculator-web build

FROM nginxinc/nginx-unprivileged:1.27-alpine
COPY elzzes-calculator-web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/elzzes-calculator-web/dist /usr/share/nginx/html
EXPOSE 8080
