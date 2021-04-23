ARG NGINX_VERSION=1.19.7
ARG NGINX_RTMP_VERSION=1.2.1
ARG NGINX_VOD_MODULE_VERSION=1.27

##############################
# Build the NGINX-build image.
FROM alpine:3.13 as build
ARG NGINX_VERSION
ARG NGINX_RTMP_VERSION
ARG NGINX_VOD_MODULE_VERSION

# Build dependencies.
RUN apk add --update \
  build-base \
  ca-certificates \
  curl \
  gcc \
  libc-dev \
  libgcc \
  linux-headers \
  make \
  musl-dev \
  openssl \
  openssl-dev \
  pcre \
  pcre-dev \
  pkgconf \
  pkgconfig \
  zlib-dev \
  ffmpeg \
  ffmpeg-dev

# Get nginx source.
RUN cd /tmp && \
  wget https://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz && \
  tar zxf nginx-${NGINX_VERSION}.tar.gz && \
  rm nginx-${NGINX_VERSION}.tar.gz

# Get nginx rtmp module.
RUN cd /tmp && \
  wget https://github.com/arut/nginx-rtmp-module/archive/v${NGINX_RTMP_VERSION}.tar.gz && \
  tar zxf v${NGINX_RTMP_VERSION}.tar.gz && rm v${NGINX_RTMP_VERSION}.tar.gz

# Get nginx vod module.
RUN cd /tmp && \
  wget https://github.com/kaltura/nginx-vod-module/archive/${NGINX_VOD_MODULE_VERSION}.tar.gz && \
  tar zxf ${NGINX_VOD_MODULE_VERSION}.tar.gz && rm ${NGINX_VOD_MODULE_VERSION}.tar.gz

# Compile nginx with nginx-rtmp module.
RUN cd /tmp/nginx-${NGINX_VERSION} && \
  ./configure \
  --prefix=/usr/local/nginx \
  --add-module=/tmp/nginx-rtmp-module-${NGINX_RTMP_VERSION} \
  --add-module=/tmp/nginx-vod-module-${NGINX_VOD_MODULE_VERSION} \
  --conf-path=/usr/local/nginx/conf/nginx.conf \
  --with-threads \
  --with-file-aio \
  --with-http_ssl_module \
  --with-debug \
  --with-cc-opt="-Wimplicit-fallthrough=0" && \
  cd /tmp/nginx-${NGINX_VERSION} && make && make install

# Cleanup.
RUN rm -rf /var/cache/* /tmp/*


##########################
# Build the release image.
FROM alpine:3.13

# Set default ports.
ENV HTTP_PORT 80
ENV HTTPS_PORT 443
ENV RTMP_PORT 1935

RUN apk add --update \
  ca-certificates \
  gettext \
  openssl \
  pcre \
  lame \
  libogg \
  curl \
  libass \
  libvpx \
  libvorbis \
  libwebp \
  libtheora \
  opus \
  rtmpdump \
  x264-dev \
  x265-dev \
  ffmpeg

# COPY --from=build /usr/local /usr/local
# COPY --from=build /usr/lib /usr/lib
COPY --from=build /usr/local/nginx /usr/local/nginx


# Add NGINX path, config and static files.
ENV PATH "${PATH}:/usr/local/nginx/sbin"
COPY ./nginx.conf /usr/local/nginx/conf/nginx.conf.template
RUN mkdir -p /opt/data && mkdir /www
# COPY ./static /www/static

EXPOSE 1935
EXPOSE 80

CMD envsubst "$(env | sed -e 's/=.*//' -e 's/^/\$/g')" < \
  /usr/local/nginx/conf/nginx.conf.template > /usr/local/nginx/conf/nginx.conf && \
  nginx
