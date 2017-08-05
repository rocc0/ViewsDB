FROM jenkins:2.60.2
MAINTAINER Vlad Kot <vk@clc.com.ua>

USER root

RUN mkdir /var/log/jenkins
RUN mkdir /var/cache/jenkins


RUN chown -R jenkins:jenkins /var/log/jenkins
RUN chown -R jenkins:jenkins /var/cache/jenkins
USER jenkins
ENV JAVA_OPTS="-Xmx8192m"

EXPOSE 8080
