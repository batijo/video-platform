FROM node:alpine

# Also exposing VSCode debug ports
EXPOSE 8000 9929 9230

RUN apk add --no-cache nasm python make g++ automake autoconf libtool dpkg pkgconfig libpng libpng-dev \
    && apk add vips-dev fftw-dev build-base \
    --update-cache \
    --repository https://alpine.global.ssl.fastly.net/alpine/v3.11/community \
    --repository https://alpine.global.ssl.fastly.net/alpine/v3.11/main \
    && rm -fR /var/cache/apk/*

RUN npm install -g gatsby-cli

WORKDIR /app
COPY ./package.json .
RUN yarn install && yarn cache clean
COPY . .
CMD ["yarn", "develop", "-H", "0.0.0.0" ]