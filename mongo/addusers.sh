#!/bin/bash
echo "Creating mongo users..."
mongo admin -u admin -p password --eval "db.createUser({"user": "hasher","pwd": "password","roles": ["readWrite"]});"
echo "Mongo users created."