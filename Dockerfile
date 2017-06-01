FROM alpine:3.2

ENV GOBIN /tweeter/bin
ENV PATH  $GOBIN:$PATH

WORKDIR /tweeter

EXPOSE 8080

ENTRYPOINT ["tweeter"]

COPY ./_output/bin/tweeter /tweeter/bin/tweeter
COPY ./assets /tweeter/assets
COPY ./templates /tweeter/templates
