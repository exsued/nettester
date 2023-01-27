#!/bin/bash
#if (whoami != root)
 # then echo "Please run as root"
  #else (
  make
  '/home/lazutin/Рабочий стол/dev/piTester/compiled/amd_pitester' -sites '/home/lazutin/Рабочий стол/dev/piTester/compiled/sites.txt' -log '/home/lazutin/Рабочий стол/dev/piTester/compiled/logs/' -onAlarm '/home/lazutin/Рабочий стол/dev/piTester/compiled/alarm.sh' -sessionServer "127.0.0.1:1289" -iface "enp5s0"
  #)
#fi
