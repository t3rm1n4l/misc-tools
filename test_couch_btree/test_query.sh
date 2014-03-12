#!/bin/bash

BUILDDIR="/Users/sarath/development/couchbase/couchdb"
export ERL_FLAGS="$ERL_FLAGS -pa $BUILDDIR/test/etap/ -pa $BUILDDIR/src/mapreduce/ -pa $BUILDDIR/src/couch_index_merger/ -pa $BUILDDIR/src/couch_view_parser/ -pa $BUILDDIR/src/couchdb/ -pa $BUILDDIR/src/snappy"
erlc -I $BUILDDIR/src/couchdb test_query.erl
erl -s test_query start -s init stop

