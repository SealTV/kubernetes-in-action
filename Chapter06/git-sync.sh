#!/bin/bash
trap "exit" SIGINT

echo GIT_SYNC_REPO $GIT_SYNC_REPO
echo GIT_SYNC_DEST $GIT_SYNC_DEST
echo GIT_SYNC_BRANCH $GIT_SYNC_BRANCH
echo GIT_SYNC_WAIT $GIT_SYNC_WAIT
echo $(git --work-tree=$GIT_SYNC_DEST --git-dir=$GIT_SYNC_DEST/.git remote -v)
git --work-tree=$GIT_SYNC_DEST --git-dir=$GIT_SYNC_DEST/.git remote set-url origin $GIT_SYNC_REPO
echo $(git --work-tree=$GIT_SYNC_DEST --git-dir=$GIT_SYNC_DEST/.git remote -v)

while :
do
    echo $(git --work-tree=$GIT_SYNC_DEST --git-dir=$GIT_SYNC_DEST/.git status) 
    git --work-tree=$GIT_SYNC_DEST --git-dir=$GIT_SYNC_DEST/.git pull origin $GIT_SYNC_BRANCH
    sleep $GIT_SYNC_WAIT
done