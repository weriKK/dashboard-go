FROM alpine

RUN apk add --no-cache ca-certificates
RUN apk add --update curl && rm -rf /var/cache/apk/*

WORKDIR /dashboard
USER nobody:nobody

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD curl -f http://localhost:8080/status || exit 1

EXPOSE 8080
ENTRYPOINT ["./api-service"]


COPY --chown=nobody:nobody api-service /dashboard