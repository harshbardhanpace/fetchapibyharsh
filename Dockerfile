# temp container 
FROM golang:alpine as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN rm -rf /app/logs
RUN mkdir /app/logs
RUN apk --no-cache add build-base jq
# RUN go test
RUN go build -o space .
# Final build with minimal FS
FROM golang:alpine as finalBuild
RUN apk add --no-cache tzdata
WORKDIR /app
RUN rm -rf /app/logs
RUN mkdir /app/logs
COPY --from=builder /app/space /app/space
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/resources /app/resources
COPY --from=builder /app/.env /app/.env
ENV TZ="Asia/Kolkata"
CMD ["/app/space"]
