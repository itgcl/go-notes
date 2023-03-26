FROM golang:1.19

WORKDIR /build


RUN echo abcbedf>>a.txt

RUN echo 12345>>a.txt

RUN touch /var/log/dev.log
CMD tail -f /var/log/dev.log
