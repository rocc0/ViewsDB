FROM mongo:latest

MAINTAINER <vlad>

ADD mongod.conf /etc/mongod.conf
ADD create_users.js /tmp/
RUN echo "Creating mongo users..."
RUN mongod --fork --auth --logpath /var/log/mongodb.log \
    && sleep 5 \
    && mongo hashes /tmp/create_users.js