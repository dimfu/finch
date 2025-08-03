FROM node:22-alpine AS build
WORKDIR /var/www/client

COPY client .
RUN corepack enable
RUN npm i
RUN npm run build

FROM node:22-alpine
WORKDIR /var/www/client

COPY --from=build /var/www/client/.output/ .
ENV PORT=8081
ENV HOST=0.0.0.0
EXPOSE 8081
CMD ["node", "server/index.mjs"]