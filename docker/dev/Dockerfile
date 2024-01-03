FROM postgres:15


RUN apt-get update && apt-get  install -y postgresql-15-postgis-3  


CMD ["/usr/local/bin/docker-entrypoint.sh","postgres"]