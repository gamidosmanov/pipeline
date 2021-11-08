FROM golang AS compilation_stage
RUN mkdir -p /go/src/pipeline
WORKDIR /go/src/pipeline
ADD bugger.go .
ADD main.go .
ADD pipeline.go .
ADD reader.go .
ADD stages.go .
ADD go.mod .
RUN go install .

FROM alpine:latest
LABEL version="1.0"
LABEL maintainer="Gamid Osmanov<gamid.osmanov@mail.ru>"
WORKDIR /root/
COPY --from=compilation_stage /go/bin/pipeline .
ENTRYPOINT ./pipeline