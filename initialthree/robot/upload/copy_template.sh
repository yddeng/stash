#!/bin/sh

set -ex
cd $(dirname $0)

test -f inc.sh || cp inc.sh.template inc.sh
test -f Makefile || cp Makefile.template Makefile