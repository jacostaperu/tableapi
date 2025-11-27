#!/bin/bash
#
# create folders

sudo mkdir -p /opt/tableapi
sudo chmod 755 /opt/tableapi

sudo mv ./tableapi/* /opt/tableapi/
sudo chmod 755 /opt/tableapi/tableapi.sh
sudo cp ./server /opt/tableapi/


sudo mv tableapi.service /etc/systemd/system/tableapi.service
sudo systemctl daemon-reload
sudo systemctl enable tableapi
sudo systemctl start tableapi
