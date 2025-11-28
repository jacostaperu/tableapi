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

SERVICE=tableapi.service

 if !  systemctl stop $SERVICE; then
   echo "$SERVICE is not running!!!"
 else
   echo stopping $SERVICE
 fi

  if !  systemctl disable $SERVICE; then
    echo "$SERVICE is not enable to start at boot!!!"
  else
    echo "disabling $SERVICE to start at boot."
  fi

   if !  rm /etc/systemd/system/$SERVICE; then
     echo "there is no /etc/systemd/system/$SERVICE we are good"
   else
     echo "cleaning files  rm /etc/systemd/system/$SERVICE  OK."
   fi

 systemctl daemon-reload


rm -rf  /opt/tableapi/bin
echo "/opt/tableapi/tables/ is your data and is not going to be deleted (at least not now not by me)."



#looking for current instances if any  and kill them
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

# Verify removal
if systemctl status $SERVICE > /dev/null 2>&1; then
  echo "Error: $SERVICE still exists."
  exit 1
fi

if systemctl list-unit-files | grep -q $SERVICE; then
  echo "Error: $SERVICE still listed in unit files."
  exit 1
fi

echo "$SERVICE successfully removed."

exit 0
