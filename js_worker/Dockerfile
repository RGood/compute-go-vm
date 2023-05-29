FROM node:12.18.1
WORKDIR /srv
COPY js_worker/package.json /srv
COPY js_worker/yarn.lock /srv
COPY protos /protos
COPY js_worker/src/ /srv/src
COPY js_worker/server /srv
RUN yarn && yarn build
