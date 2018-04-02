# https://github.com/elastic/elasticsearch-docker
#FROM elasticsearch:5.3.0 # deprecated images https://github.com/deviantony/docker-elk/issues/119
FROM docker.elastic.co/elasticsearch/elasticsearch:5.3.0
RUN chmod 755 /data/elastic/es-plugins.sh
RUN bin/elasticsearch-plugin install analysis-ukrainian
# Add your elasticsearch plugins setup here
# Example: RUN elasticsearch-plugin install analysis-icu