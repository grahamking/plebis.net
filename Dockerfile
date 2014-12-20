# sudo docker run --name=plebis -d localhost:5000/plebis
FROM debian:stable
MAINTAINER Graham King <graham@gkgk.org>
RUN mkdir -p /opt/plebis && chown www-data:www-data /opt/plebis
WORKDIR /opt/plebis
ADD plebis /opt/plebis/
ADD index.html /opt/plebis/
RUN chown www-data:www-data plebis index.html
EXPOSE 8081
USER www-data
RUN mkdir -p /opt/plebis/data && chown www-data /opt/plebis/data
VOLUME ["/opt/plebis/data/"]
CMD ["/opt/plebis/plebis", "-p", "8081", "-x", "/opt/plebis/index.html", "-s", "/opt/plebis/data/store.dat"]
