#!/bin/sh
supervisord > /tmp/log 2>&1 &
echo Running test suite for $PLUGIN
ginkgo -r
