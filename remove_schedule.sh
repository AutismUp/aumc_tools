#!/bin/bash
echo "Removed Minecraft start and stop schedule"
crontab -r
crontab /mcdata/msm/no-scheduled-crontab
