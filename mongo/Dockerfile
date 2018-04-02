FROM ubuntu:16.04
MAINTAINER TonyWang <plantpark.net@gmail.com>
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv EA312927
RUN echo "deb http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.2 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-3.2.list
RUN apt-get update && apt-get install -y mongodb-org
RUN mkdir -p /data/db

ADD addusers.sh /addusers.sh

EXPOSE 27017 28017
CMD ["mongod", "--auth"]
CMD ["/set_mongodb_password.sh"]
ENTRYPOINT usr/bin/mongod

