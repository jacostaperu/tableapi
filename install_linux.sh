#!/bin/bash
#
# create folders
if [ -n "$SUDO_USER" ]; then
  echo "This command was run by user: $SUDO_USER"
else
  echo "Not run with sudo or SUDO_USER missing"
  echo "please try again with sudo "
  exit
fi

#looking for current instances
PID=$(ss -tlnp | grep :8087 | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2)

if [ -n "$PID" ]; then
  echo "Found server process listening on port 8087 with PID: $PID"
  read -p "Put a hit on $PID? [y/N]: " -n 1 -r
  echo    # move to new line
  if [[ $REPLY =~ ^[Yy]$ ]]; then

    kill "$PID"
    echo "Process $PID has been killed."

    PID=$(ss -tlnp | grep :8087 | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2)

    if [ -n "$PID" ]; then
      echo "Oh no, I couldn't eliminate : $PID"
      echo "You need to kill $PID yourself...  Good Bye."
      exit 0
    fi


  else
    echo "We're giving pid $PID a pass this time...Good Bye."

    exit 0
  fi
else
  echo "No server process found listening on port 8087."
fi


echo Deploying updated binaries


mkdir -p /opt/tableapi/bin
mkdir -p /opt/tableapi/tables
chmod 755 /opt/tableapi/bin
chmod 755 /opt/tableapi/tables

#copy the binaries
if ! install server /opt/tableapi/bin; then
  echo "Failed to install "server" to /opt/tableapi/bin" >&2
  echo "please fix the error and try again...Good bye"
  exit 1
else
    echo Binaries installed OK!!!!
fi

cp tables/PIN_Table.csv /opt/tableapi/tables

if [ -e  /opt/tableapi/tables/PIN_Table.csv ]; then
  echo "Warning: /opt/tableapi/tables/PIN_Table.csv already exists, then NOT overwriting it!" >&2
else
  cp ./server /opt/tableapi/bin
fi

#create the systemd service unit
cat <<EOF > /etc/systemd/system/tableapi.service
[Unit]
Description=table api
After=network.target

[Service]
Type=simple
User=$SUDO_USER
WorkingDirectory=/opt/tableapi
ExecStart=/opt/tableapi/bin/server -tablespath /opt/tableapi/tables
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

sudo chown -R $SUDO_USER /opt/tableapi

systemctl daemon-reload
systemctl enable tableapi
systemctl start tableapi
systemctl status tableapi


#testing the presence of the running instance
ss -lptun | grep 8087
