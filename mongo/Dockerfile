FROM ubuntu:16.04
MAINTAINER TonyWang <plantpark.net@gmail.com>
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv EA312927
RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | tee /etc/apt/sources.list.d/mongodb.list
RUN apt-get update && apt-get install -y mongodb-10gen
RUN mkdir -p /data/db

ADD addusers.sh /addusers.sh

EXPOSE 27017 28017
CMD ["mongod", "--auth"]
CMD ["/set_mongodb_password.sh"]
ENTRYPOINT usr/bin/mongod

