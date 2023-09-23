FROM ubuntu:22.04
RUN mkdir /app
COPY ./image_filter_oneapi /app/image_filter_oneapi
COPY ./ui /ui
COPY ./libtbb.so.12 /usr/lib/libtbb.so.12

ENTRYPOINT [ "/app/image_filter_oneapi" ]