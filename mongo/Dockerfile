FROM mongo:latest

MAINTAINER <vlad>

ADD create_users.js /tmp/
RUN echo "Creating mongo users..."
RUN mongod --fork --auth --logpath /var/log/mongodb.log \
    && sleep 5 \
    && mongo <YOUR DATABASE> /tmp/create_ddbb.js