# MacShutdownWS
Webservice to initiate shutdown procedure for Mac OS

### Install the process:
sudo launchctl load -w /Library/LaunchDaemons/com.macshutdown.plist

### Stopping the process:
sudo launchctl stop /Library/LaunchDaemons/com.macshutdown.plist

### Start the process:
sudo launchctl start -w /Library/LaunchDaemons/com.macshutdown.plist
