# build stage
FROM node:lts-alpine as build-env
ENV NODE_ENV=production
RUN \ 
  apk update && \
  apk add git gcc g++

ADD . /opt/app
WORKDIR /opt/app

RUN \
  yarn install && \
  yarn build

FROM nginx:alpine
COPY --from=build-env /opt/app/public /var/www
COPY --from=build-env /opt/app/envsubt.sh /etc/nginx/templates/envsubt.sh
COPY --from=build-env /opt/app/default.conf /etc/nginx/templates/default.conf.template
WORKDIR /var/www

RUN apk add --no-cache bash

ENV PORT=80
ENV REST_SERVER=localhost:8080

CMD ["/bin/sh", "-c", "/etc/nginx/templates/envsubt.sh && nginx -g \"daemon off;\""]
