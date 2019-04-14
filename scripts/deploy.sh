#!/usr/bin/env sh

PROJECT_DIR=$(pwd)/[!.][!vendor][!src]*
BIN_FILE="/srv/go_blog/bin/server"
POSTS_DIR="/srv/go_blog/posts"


testForCredentials() {
  [ -z "$DO_USER" ] \
    && { echo "SSH user environment variable not set."; exit 1; }
  [ -z "$DO_BLOG_IP" ] \
    && { echo "SSH ip address environment variable not set."; exit 1; }
}

deploy() {
  testForCredentials;

  echo "Building...";
  make build-deploy;
  echo "Done.";
  echo "Uploading to server...";
  ssh $DO_USER@$DO_BLOG_IP "if [ -z $PROJECT_DIR ]; then mkdir /srv/go_blog; else echo 'Folder go_blog exists'; fi;";
  ssh $DO_USER@$DO_BLOG_IP "[ -z $BIN_FILE ] && rm -f $BIN_FILE;";
  scp -rp $PROJECT_DIR $DO_USER@$DO_BLOG_IP:/srv/go_blog;
  ssh $DO_USER@$DO_BLOG_IP "if [ -z $POSTS_DIR ]; then mkdir $POSTS_DIR; else echo 'Folder posts exists'; fi;";
  scp -rp ./posts/* $DO_USER@$DO_BLOG_IP:/srv/go_blog/posts;
  echo "Done.";
}

deploy;
