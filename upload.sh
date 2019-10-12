#!/bin/bash
echo $RCLONE_CONFIG | base64 -d > ~/.rclone.conf
cat ~/.rclone.conf
RAND_STR=$(pwgen -N 1 -s 32)
URL="https://s3.eu-central-1.amazonaws.com/tmp.lvlup.pro/$RAND_STR/tasktab.apk"
echo "Uploading to URL: $URL"
rclone -q copyto app/build/outputs/apk/debug/app-debug.apk s3:tmp.lvlup.pro/$RAND_STR/tasktab.apk && curl -s \
--form-string "token=$PUSHOVER_TOKEN" \
--form-string "user=$PUSHOVER_USER" \
--form-string "title=New TaskTab Android version" \
--form-string "message=Click here to update" \
--form-string "priority=-1" \
--form-string "url=$URL" \
--form-string "url_title=Download *.apk here" \
https://api.pushover.net/1/messages.json