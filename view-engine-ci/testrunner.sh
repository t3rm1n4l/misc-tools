#!/bin/bash

cd ~/couchbase/testrunner

while true;
do
    TEST_REPORT="~/tests/$(date +%Y-%m-%d-%H:%M)"

    mkdir -p $TEST_REPORT

    make simple-test &> $TEST_REPORT/simple.txt
    make test-views-pre-merge &> $TEST_REPORT/premerge.txt
    make test-viewmerge &> $TEST_REPORT/viewmerge.txt
    sleep 14400
done


