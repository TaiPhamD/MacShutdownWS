# MacShutdownWS
Webservice to initiate shutdown procedure for Mac OS

## Build process

Install go lang dev tool and build the .go file. Then create a config.txt in the same path as the compiled binary.
The config.txt should consist of 2 lines where the first line is your endpoint basic auth password and the 2nd line is your webservice port

### Install the process:
sudo launchctl load -w /Library/LaunchDaemons/com.macshutdown.plist

### Stopping the process:
sudo launchctl stop /Library/LaunchDaemons/com.macshutdown.plist

### Start the process:
sudo launchctl start -w /Library/LaunchDaemons/com.macshutdown.plist
