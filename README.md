This program/API will listen to designated endpoints&port and acts accordingly to each HTTP feeding. Each call made will out to slack nicely as ephemeral reply.
Combine with system-daemon such as systemd/sysvinit to allocate worker/PID for each endpoints call made.
Can be used literally/conceptually for similiar use case where automatic/macro action/operation are heavily required as well.

Refference used/related in this repo:
- [Slack App](https://api.slack.com/authentication/basics)
- [Expect tcl](https://man7.org/linux/man-pages/man1/expect.1.html)
- [JIRA API](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/)
- [GIO API](https://vdc-download.vmware.com/vmwb-repository/dcr-public/1b6cf07d-adb3-4dba-8c47-9c1c92b04857/241956dd-e128-4fcc-8131-bf66e1edd895/vcloud_sp_api_guide_30_0.pdf)
- [Nextcloud API](https://docs.nextcloud.com/server/latest/developer_manual/client_apis/index.html)
