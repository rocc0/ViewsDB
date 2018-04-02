#!/bin/bash
if bin/elasticsearch-plugin list | grep -q 'ukrainian'; then
   echo "exists"
else
   bin/elasticsearch-plugin install analysis-ukrainian
fi