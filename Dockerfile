# sudo docker run --rm -p 8081:8081 graham/plebis
FROM debian:stable
MAINTAINER Graham King <graham@gkgk.org>
RUN mkdir -p /opt/plebis && chown www-data /opt/plebis
USER www-data
WORKDIR /opt/plebis
COPY plebis /opt/plebis/
COPY index.html /opt/plebis/
EXPOSE 8081
CMD ["/opt/plebis/plebis", "-p", "8081", "-x", "/opt/plebis/index.html", "-s", "/opt/plebis/store.dat"]
