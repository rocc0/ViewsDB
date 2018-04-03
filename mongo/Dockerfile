FROM mongo:latest

MAINTAINER <vlad>

ADD mongod.conf /etc/mongod.conf
ADD create_users.js /tmp/
RUN echo "Creating mongo users..."
RUN mongod \
    && sleep 5 \
    && mongo hashes /tmp/create_ddbb.js