#!/bin/bash

set -u

DIR="$(cd "$(dirname "$0")" && pwd)"
cd $DIR

go install .

compare_output() {
	args=$1

	fingOut=$(fing $args | sort)
	findOut=$(find $args | sort)

	result=$(diff <(echo "${fingOut}") <(echo "${findOut}"))
	if [ $? -ne 0 ]; then
		echo "--------------------"
		echo "output mismatch($1)"
		echo "--------------------"
		echo "[fing output]"
		echo "${fingOut}"
		echo "--------------------"
		echo "[find output]"
		echo "${findOut}"
		echo "--------------------"
		echo "[diff]"
		echo "${result}"
		echo "--------------------"
		exit 1
	fi
}

compare_output "testdata/jpg_dir/1.jpg"
compare_output "testdata -name *.jpg"
compare_output "testdata -name jpg_dir -prune -print"
compare_output "testdata -name jpg_dir -prune -o -type f"
compare_output "testdata -name jpg_dir -prune -false -o -type f"
compare_output "testdata -name jpg_dir -o -name png_dir"
