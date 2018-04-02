#!/bin/bash
bin/elasticsearch-plugin list | grep 'ukrainian' &> /dev/null
if [ $? == 0 ]; then
   bin/elasticsearch-plugin install analysis-ukrainian
fi