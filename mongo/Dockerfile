# https://github.com/elastic/elasticsearch-docker
#FROM elasticsearch:5.3.0 # deprecated images https://github.com/deviantony/docker-elk/issues/119
FROM mongo:latest
RUN mongod --auth
RUN echo "Creating mongo users..."
RUN mongo admin --host localhost -u admin -p password --eval "db.createUser({"user": "hasher","pwd": "password","roles": ["readWrite"]});"
RUN echo "Mongo users created."
# Add your elasticsearch plugins setup here
# Example: RUN edocklasticsearch-plugin install analysis-icu