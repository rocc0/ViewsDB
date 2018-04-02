#!/bin/bash
bash bin/plugin list | grep 'ukrainian' &> /dev/null
if [ $? == 0 ]; then
   bash bin/elasticsearch-plugin install analysis-ukrainian
fi