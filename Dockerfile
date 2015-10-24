FROM alpine:3.2

ENV GOBIN /oinker/bin
ENV PATH  $GOBIN:$PATH

WORKDIR /oinker

EXPOSE 8080

ENTRYPOINT ["oinker"]

COPY ./_output/bin/oinker /oinker/bin/oinker
COPY ./assets /oinker/assets
COPY ./templates /oinker/templates
