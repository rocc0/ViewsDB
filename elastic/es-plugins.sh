#!/bin/bash
sh bin/plugin list | grep 'ukrainian' &> /dev/null
if [ $? == 0 ]; then
   sh bin/elasticsearch-plugin install analysis-ukrainian
fi