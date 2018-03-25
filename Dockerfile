FROM alpine:3.7
ARG yversion
LABEL version=$yversion \
      description="YAGES gRPC server" \
      maintainer="michael.hausenblas@gmail.com"
COPY ./srv-yages /app/srv-yages
WORKDIR /app
RUN chown -R 1001:1 /app
USER 1001
RUN chmod +x srv-yages
EXPOSE 9000
CMD ["/app/srv-yages"]