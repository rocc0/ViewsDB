#!/usr/bin/env bash
echo "Creating mongo users..."
mongo admin --host localhost -u admin -p password --eval "db.createUser({"user": "hasher","pwd": "password","roles": ["readWrite"]});"
echo "Mongo users created."