#!/usr/bin/env python3

from mc_classes import MCConfig

import argparse
from datetime import datetime
from pathlib import Path
import glob
import json
import subprocess



class EnterDir:
    """Context manager for changing the current working directory"""
    def __init__(self, newPath):
        self.newPath = os.path.expanduser(newPath)

    def __enter__(self):
        self.savedPath = os.getcwd()
        os.chdir(self.newPath)

    def __exit__(self, etype, value, traceback):
        os.chdir(self.savedPath)


print("Importing configuration from config.json")
with open('config.json', 'r') as config_file:
    config = json.load(config_file)


def create_world(day, jargroup, version, config=config):
    msm_server_path = config['msm_server_path']
    now = datetime.now()
    time_stamp = now.strftime('%a %b %d %X EST %Y')

    # 1 - create the world
    subprocess.call(['sudo', 'msm', 'server', 'create', day])
    subprocess.call(['sudo', 'msm', day, 'jar', jargroup])

    # 2 - create the eula.txt file
    eula_file_path = Path(msm_server_path, day, 'eula.txt')
    eula_file = open(eula_file_path, 'w')
    eula_file.write("#By changing the setting below to TRUE you are indicating your agreement to our EULA (https://account.mojang.com/documents/minecraft_eula).\n")
    eula_file.write(f'#{time_stamp}\n')
    eula_file.write('eula=true')
    eula_file.close()
    
    # 3 - update server.properties template and copy to the server folder
    server_properties = MCConfig('server.properties.template')
    server_properties.update_config('msm-version', f'minecraft/{version}')
    server_properties.update_config('motd', f'Autism Up - {day}')
    server_properties.write_config(f'{msm_server_path}/{day}/server.properties')

    print(f'World named "{day}" created')

if __name__ == '__main__':

    parser = argparse.ArgumentParser(description='Create new Minecraft worlds for the semester')
    parser.add_argument('--jargroupname', help='The name of the jargroup to create with')
    parser.add_argument('--version', help='The version of Minecraft the jargroup uses')
    args = parser.parse_args()

    print("Loading configuration")

    print("Creating new worlds")
    # create_world('remote', args.jargroupname, args.version)
    # create_world('youth_415', args.jargroupname, args.version)
    # create_world('youth_530', args.jargroupname, args.version)
    # create_world('teen_645', args.jargroupname, args.version)
    # create_world('survival_saturday', args.jargroupname, args.version)
    create_world('experiment', args.jargroupname, args.version)



