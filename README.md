# s3_ip_server

Assuming you have a list of ips in a format like this https://www.dshield.org/ipsascii.html?limit=100 in an S3 bucket, this server will nicely format them into HTML.

## Requirements

Environment varables:

* `AWS_REGION`
* `AWS_BUCKET`
* `AWS_OBJECT`


## User Data


```

#!/bin/bash

# Install Golang
export ARTIFACT=go1.11.5.linux-amd64.tar.gz
curl -o /tmp/$ARTIFACT "https://dl.google.com/go/$ARTIFACT"
tar -C /usr/local -xzf $ARTIFACT
GOPATH=/opt

# Install our server
go get github.com/onetwopunch

# Create a systemd service to get it to run in the background
cat << 'SVC' > /etc/systemd/system/dpg.service
[Unit]
Description=Distributed Password Guessing Scenario
After=network.target

[Service]
Type=simple
User=ec2-user
ExecStart=/opt/bin/s3_ip_server
Restart=on-failure
Environment=PATH=$PATH:/usr/local/go/bin
Environment=GOPATH=/opt

## EDIT THESE ##
Environment=AWS_REGION=[edit me]
Environment=AWS_BUCKET=[edit me]
Environment=AWS_OBJECT=[edit me]
################

[Install]
WantedBy=multi-user.target

SVC

systemctl daemon-reload
systemctl start dpg


```