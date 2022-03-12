#!/bin/bash

id=''
user=''
verbose='false'

while getopts 'i:u:' flag; do
  case "${flag}" in
    i) id="${OPTARG}" ;;
    u) user="${OPTARG}" ;;
    *) exit 1 ;;
  esac
done
datetime=$(date +%s); path="/var/www/网倡本/${id}-${datetime}.png"
wget --save-cookies="/tmp/zc_${datetime}" --keep-session-cookies --post-data "name=admin&pass=zabbix&enter=Sign+in" -O /dev/null -q "http://$ZBXIP/zabbix/index.php?login=1"
wget --load-cookies="/tmp/zc_${datetime}"  -O ${path} -q "http://$ZBXIP/zabbix/chart3.php?items[0][itemid]=${id}&width=1280&from=now-30d&to=now"
chart_url="https://$ZBXIP/${id}-${datetime}.png"; rm -f /tmp/zc_${datetime}
payload="payload={\"attachments\": [{\"title\": \"<@${user}>\",\"color\": \"${color}\",\"image_url\": \"${chart_url}\"}]}"
curl -sS --data-urlencode "${payload}" "https://hooks.slack.com/services/$TENANT/$CHANNEL/$WEBHOOK"
