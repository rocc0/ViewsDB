FROM mongo:latest

MAINTAINER <vlad>

ADD mongod.conf /etc/mongod.conf
ADD create_users.js /tmp/
RUN echo "Creating mongo users..."
RUN mongod --fork --auth --logpath /var/log/mongodb.log \
    && echo "Waiting..." \
    && sleep 5 \
    && echo "Start..." \
    && mongo hashes /tmp/create_users.js -u admin -p password --authenticationDatabase admin