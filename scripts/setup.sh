#!/bin/bash

cd $WORKSPACE
echo "****************************************************"
echo "datanimbus.io.b2b.agent.watcher :: Fetching dependencies"
echo "****************************************************"
go get -u github.com/kardianos/service
go get -u github.com/ian-kent/go-log/log