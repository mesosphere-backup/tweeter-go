FROM alpine:3.2

COPY ./_output/bin/oinker /oinker/bin/oinker
COPY ./assets /oinker/assets
COPY ./templates /oinker/templates

ENV GOBIN /oinker/bin
ENV PATH  $GOBIN:$PATH

WORKDIR /oinker

EXPOSE 8080

#ENTRYPOINT []
ENTRYPOINT ["oinker"]
