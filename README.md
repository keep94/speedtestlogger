# speedtestlogger

The app logs your home internet speeds. It is designed for MacOS.

## Prerequisites

You must have the Speedtest cli by Ookla. Go [here](https://www.speedtest.net/apps/cli) for instructions on how to install.

## Architecture

This app actually contains 3 programs:

- **stlinit:** You run this once during setup to initialize the sqlite file that will hold the internet speeds.
- **stllog:** This program uses the CSV file from the Speedtest cli by Ookla to log the current internet speed in the sqlite file created by stlinit. Typically a cron job runs this program along with the Speedtest cli by Ookla once every hour.
- **stlview** This program is the web server that serves the pages showing historical internet speeds.

## Setup

1. Clone this repository to ~/go/src/github.com/keep94/speedtestlogger.
2. Compile the three programs from the previous section by running `go install ...` The 3 compiled executables can be found in ~/go/bin.
3. Create the ~/stl directory to store the internet speeds.
4. Run `~/go/bin/stlinit -db ~/stl/stl.db` to create the sqlite3 file.
5. Create the bash script to record the current internet speed and store in `~/go/bin/speedtest.sh`

```bash
#!/bin/zsh

if /opt/homebrew/bin/speedtest -f csv &> stl_out.csv; then
    /Users/youruserid/go/bin/stllog -db stl.db -csv stl_out.csv
else
    /Users/youruserid/go/bin/stllog -db stl.db
fi
```
This tells the Speedtest cli to log the internet speeds to a csv file called stl_out.csv. If successful it runs the stllog program with that csv file. If unsuccessful, it runs stllog without a csv file, which causes stllog to record 0 upload speed and 0 download speed. The current working directory of this script will be `~/stl`

6. Create the plist file to run speedtest.sh

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>Label</key>
        <string>yourlastname.speedtest</string>

        <key>UserName</key>
        <string>youruserid</string>

        <key>ProgramArguments</key>
        <array>
          <string>/Users/youruserid/go/bin/speedtest.sh</string>
        </array>

        <key>StartCalendarInterval</key>
        <array>
            <dict>
                <key>Minute</key>
                <integer>(choose an int between 0 and 59 inclusive)</integer>
            </dict>
        </array>

        <key>WorkingDirectory</key>
        <string>/Users/youruserid/stl</string>

        <key>StandardErrorPath</key>
        <string>/Users/youruserid/logs/speedtest/stderr.log</string>

        <key>StandardOutPath</key>
        <string>/Users/youruserid/logs/speedtest/stdout.log</string>

</dict>
</plist>
```
Store in `/Library/LaunchDaemons` in a file called `yourlastname.speedtest.plist`

7. Create the plist file for stlview

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>Label</key>
        <string>yourlastname.stlview</string>

        <key>KeepAlive</key>
        <true/>

        <key>UserName</key>
        <string>youruserid</string>

        <key>ProgramArguments</key>
        <array>
          <string>/Users/youruserid/go/bin/stlview</string>
          <string>-db</string>
          <string>/Users/youruserid/stl/stl.db</string>
          <string>-http</string>
          <string>:9792</string>
        </array>

        <key>WorkingDirectory</key>
        <string>/Users/youruserid/stl</string>

        <key>StandardErrorPath</key>
        <string>/Users/youruserid/logs/stlview/stderr.log</string>

        <key>StandardOutPath</key>
        <string>/Users/youruserid/logs/stlview/stdout.log</string>

</dict>
</plist>
```
You can see the download speeds by going to [http://localhost:9792](http://localhost:9792). You can change this by entering a different port number besides 9792 in the plist file.
