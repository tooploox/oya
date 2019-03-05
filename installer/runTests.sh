#!/usr/bin/env bash

TEST_IMG_NAME="oya-installer-test"
OUTPUT="$(mktemp -d)/install.log"

runTest="docker run -v $(pwd):/oya -t $TEST_IMG_NAME"
errors=0

red=$'\e[1;31m'
grn=$'\e[1;32m'
end=$'\e[0m'

echo "Test log file $OUTPUT"
printf "Building docker container $TEST_IMG_NAME: "
docker build -t $TEST_IMG_NAME . 2>&1 >> $OUTPUT
printf "%s\n" "${grn}DONE${end}"


printf "Oya installer fails if sops is missing: "
$runTest /oya/noSopsTest.sh 2>&1 >> $OUTPUT
if [ $? -eq 0 ]; then
    printf "%s\n" "${red}FAIL${end}"
    errors=1
else
    printf "%s\n" "${grn}PASS${end}"
fi

printf "Oya installer passes if all good: "
$runTest /oya/successTest.sh 2>&1 >> $OUTPUT
if [ $? -eq 0 ]; then
    printf "%s\n" "${grn}PASS${end}"
else
    printf "%s\n" "${red}FAIL${end}"
    errors=1
fi

exit $errors
