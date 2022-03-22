#!/usr/bin/env python3

import argparse
from pathlib import Path
import glob
import json
import os
import shutil
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



def build_new_jar(config):
    '''Cleans up previous jar build content and builds the desired version of Spigot
       (latest is default)
    '''

    print(f"Removing temporary directories from {Path(config['build_directory'])}")
    for directory in config['temp_folders']:
        try:
            print(Path(config['build_directory'], directory))
            shutil.rmtree(Path(config['build_directory'], directory))
        except OSError as e:
            print(f"Error: {Path(config['build_directory'], directory)}, {e.strerror}")

    print(f"Removing temporary files from {Path(config['build_directory'])}")
    for temp_file in config['temp_files']:
        try:
            print(Path(config['build_directory'], temp_file))
            Path(config['build_directory'], temp_file).unlink()
        except OSError as e:
            print("Error: %s : %s" % (temp_file, e.strerror))

    if config['delete_spigot_jars']:
        print(f"Removing old spigot jars from {Path(config['build_directory'])}")
        old_jar_files = Path(config['build_directory']).glob('spigot*.jar')
        for jar_file in old_jar_files:
            try:
                jar_file.unlink()
            except OSError as e:
                print("Error: %s : %s" % (jar_file, e.strerror))

    print('Running BuildTools')
    with EnterDir(Path(config['build_directory'])):
        subprocess.call(['java', '-jar', 'BuildTools.jar', '--rev', config['minecraft_version']])
        new_jar_files = Path(config['build_directory']).glob('spigot*.jar')
        for jar_file in new_jar_files:
            shutil.copy(jar_file, Path(config['jar_git_repo'], 'jars'))
    print('BuildTools complete')

def publish_new_jar(config):
    '''Push the noted jarfile of Minecraft to GitHub and create a new JarGroup in MSM'''
    
    filename = f"spigot-{config['minecraft_version']}.jar"

    with EnterDir(Path(config['jar_git_repo'])):

        if os.path.isfile(f'jars/{filename}'):
            print('Commiting new jar to Github')
            subprocess.call(['git', 'add', '-A'])
            subprocess.call(['git', 'commit', '-m', 'New Minecraft jar added via script'])
            subprocess.call(['git', 'push', 'origin', 'master'])
            print('New jar committed to Github')
        else:
            raise FileNotFoundError(f"The jar file is not in {config['jar_git_repo']}. Did you build it?")

    
    latest_jargroup_name = filename.replace('spigot-', '').replace('.jar', '').replace('.', '_').replace('-', '_')

    latest_jargroup_url = f'https://github.com/thehatchcloud/minecraft_jars/raw/master/jars/{filename}'
    print("Adding the new jargroup")
    subprocess.call(['sudo', 'msm', 'jargroup', 'create', latest_jargroup_name, latest_jargroup_url])
    subprocess.call(['sudo', 'msm', 'jargroup', 'getlatest', latest_jargroup_name])

    print("Adding the latest jar is complete!")
    print(f"New jargroup name is {latest_jargroup_name}")


def main(args):

    print(f"Importing configuration from {args.config}")
    with open(args.config, 'r') as config_file:
        config = json.load(config_file)['build_config']

    build_new_jar(config)
    publish_new_jar(config)


if __name__ == '__main__':

    parser = argparse.ArgumentParser(description='Build new Minecraft jar file and jargroup')
    parser.add_argument('--config', help='Configuration file to reference for build')
    args = parser.parse_args()

    main(args)
