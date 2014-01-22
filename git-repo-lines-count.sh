#!/bin/bash

if [[ -z "$1" ]];
then
    echo "Usage: $0 git_repo_path"
    exit 1
fi

cd $1
while read author
do
    for cid in `git log --author "$author" --oneline | awk '{ print $1 }' | xargs`;
    do
        git show $cid --stat | tail -n 1 | awk '{print $4,$6 }';
    done | awk -v author="$author" '{ inserted+=$1; removed+=$2 } END{ printf("%-50s %-5d %-5d\n", author,inserted,removed); }'
done < <(git log --pretty="%an" | sort -u)

