#!/usr/bin/env python3

import subprocess
import sys

backout = input('Do you rea!lly want to delete the current worlds? (y/N)')
if backout != 'y':
    print('Deletion aborted.')
    sys.exit()

else:
    print("Deleting the current worlds")
    subprocess.call(['sudo', 'msm', 'server', 'delete', 'remote'])
    subprocess.call(['sudo', 'msm', 'server', 'delete', 'youth_415'])
    subprocess.call(['sudo', 'msm', 'server', 'delete', 'youth_530'])
    subprocess.call(['sudo', 'msm', 'server', 'delete', 'teen_645'])
    subprocess.call(['sudo', 'msm', 'server', 'delete', 'survival_saturday'])
    print("Current worlds have been deleted")
