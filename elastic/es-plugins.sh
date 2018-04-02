#!/bin/bash
bin/plugin list | grep 'ukrainian' &> /dev/null
if [ $? == 0 ]; then
   bin/elasticsearch-plugin install analysis-ukrainian
fi