FROM node:12.18.1 as build
WORKDIR /build
COPY js_worker/package.json /build
COPY js_worker/yarn.lock /build
COPY js_worker/src/ /build/src
COPY js_worker/server /build
RUN yarn && yarn build

FROM node:alpine as production
WORKDIR /srv

RUN apk add --no-cache bash

COPY --from=build /build/node_modules /srv/node_modules
COPY --from=build /build/dist /srv/dist
COPY protos /protos
COPY js_worker/server /srv
