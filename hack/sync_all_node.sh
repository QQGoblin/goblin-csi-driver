#!/bin/bash

docker image save goblin-hostpathplugin:canary --output output/goblin-hostpathplugin-canary.tar

scp output/goblin-hostpathplugin-canary.tar node2:/root/goblin-hostpathplugin-canary.tar
scp output/goblin-hostpathplugin-canary.tar node3:/root/goblin-hostpathplugin-canary.tar

ssh root@node2 "docker image load --input /root/goblin-hostpathplugin-canary.tar"
ssh root@node3 "docker image load --input /root/goblin-hostpathplugin-canary.tar"