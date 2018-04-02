#!/bin/bash
sudo bin/plugin list | grep 'ukrainian' &> /dev/null
if [ $? == 0 ]; then
   sudo bin/elasticsearch-plugin install analysis-ukrainian
fi