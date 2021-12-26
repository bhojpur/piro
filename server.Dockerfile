# build with Bhojpur GoRPA build :server-docker
FROM alpine:latest

COPY plugins--all/bin/* /app/plugins/
ENV PATH=$PATH:/app/plugins

COPY server/piro /app/piro
RUN chmod +x /app/piro
ENTRYPOINT [ "/app/piro" ]
