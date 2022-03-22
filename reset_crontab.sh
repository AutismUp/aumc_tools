#!/bin/bash
crontab -r
crontab /mcdata/msm/scheduled-crontab
echo "Enabled Minecraft start stop schedule"
