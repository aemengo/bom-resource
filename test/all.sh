#!/bin/bash

set -e

$(dirname $0)/image.sh
# $(dirname $0)/check.sh
# $(dirname $0)/get.sh


echo -e '\e[32mall tests passed!\e[0m'
