FROM ubuntu:20.04
USER root
RUN apt-get update || true \
	&& apt-get -y install  lvm2 
COPY textfile_collector /
RUN chmod +x textfile_collector
