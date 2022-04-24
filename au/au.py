import glob
import json
import os
from pathlib import Path
import shutil
import subprocess


class AuRestoreException(Exception):
    pass


class EnterDir:
    """Context manager for changing the current working directory"""
    def __init__(self, newPath):
        self.newPath = os.path.expanduser(newPath)

    def __enter__(self):
        self.savedPath = os.getcwd()
        os.chdir(self.newPath)

    def __exit__(self, etype, value, traceback):
        os.chdir(self.savedPath)


class AuMc(object):

    def __init__(self, config_file_path):

        self.config_file_path = config_file_path
        with open(config_file_path, 'r') as config_file:
            self.config = json.load(config_file) 


    def build_new_jar(self):
        '''Cleans up previous jar build content and builds the desired version of Spigot
        (latest is default)
        '''

        build_config = self.config['build_config']

        print(f"Removing temporary directories from {Path(build_config['build_directory'])}")
        for directory in build_config['temp_folders']:
            try:
                print(Path(build_config['build_directory'], directory))
                shutil.rmtree(Path(build_config['build_directory'], directory))
            except OSError as e:
                print(f"Error: {Path(build_config['build_directory'], directory)}, {e.strerror}")

        print(f"Removing temporary files from {Path(build_config['build_directory'])}")
        for temp_file in build_config['temp_files']:
            try:
                print(Path(build_config['build_directory'], temp_file))
                Path(build_config['build_directory'], temp_file).unlink()
            except OSError as e:
                print("Error: %s : %s" % (temp_file, e.strerror))

        if build_config['delete_spigot_jars']:
            print(f"Removing old spigot jars from {Path(build_config['build_directory'])}")
            old_jar_files = Path(build_config['build_directory']).glob('spigot*.jar')
            for jar_file in old_jar_files:
                try:
                    jar_file.unlink()
                except OSError as e:
                    print("Error: %s : %s" % (jar_file, e.strerror))

        print('Running BuildTools')
        with EnterDir(Path(build_config['build_directory'])):
            subprocess.call(['java', '-jar', 'BuildTools.jar', '--rev', build_config['minecraft_version']])
            new_jar_files = Path(build_config['build_directory']).glob('spigot*.jar')
            for jar_file in new_jar_files:
                shutil.copy(jar_file, Path(build_config['jar_git_repo'], 'jars'))
        print('BuildTools complete')


    def publish_new_jar(self):
        '''Push the noted jarfile of Minecraft to GitHub and create a new JarGroup in MSM'''

        build_config = self.config['build_config'] 
       
        filename = f"spigot-{build_config['minecraft_version']}.jar"

        with EnterDir(Path(build_config['jar_git_repo'])):

            if os.path.isfile(f'jars/{filename}'):
                print('Commiting new jar to Github')
                subprocess.call(['git', 'add', '-A'])
                subprocess.call(['git', 'commit', '-m', 'New Minecraft jar added via script'])
                subprocess.call(['git', 'push', 'origin', 'master'])
                print('New jar committed to Github')
            else:
                raise FileNotFoundError(f"The jar file is not in {build_config['jar_git_repo']}. Did you build it?")

        
        latest_jargroup_name = filename.replace('spigot-', '').replace('.jar', '').replace('.', '_').replace('-', '_')

        latest_jargroup_url = f'https://github.com/thehatchcloud/minecraft_jars/raw/master/jars/{filename}'
        print("Adding the new jargroup")
        subprocess.call(['sudo', 'msm', 'jargroup', 'create', latest_jargroup_name, latest_jargroup_url])
        subprocess.call(['sudo', 'msm', 'jargroup', 'getlatest', latest_jargroup_name])

        print("Adding the latest jar is complete!")
        print(f"New jargroup name is {latest_jargroup_name}")


    def restore_world(self, world_name, restore_to_date):

        # Step 1 - Check that the world exists
        status = subprocess.run(['msm', world_name, 'status'], check=True)
        if status == f'There is no server with the name "{world_name}".':
            raise AuRestoreException('The Minecraft world does not exist')
        
        # Step 2 - Check that the identified backup exists
        for backup_file in os.listdir(self.backup_world_path):
            if backup_file.startswith(restore_to_date):
                backup_filename = f"{self.backup_world_path}/{backup_file}"
        if not backup_filename:
            raise AuRestoreException('Backup file for the provided date not available.')

        # Step 3 - Stop the world
        subprocess.run(['msm', world_name, 'stop', 'now'])

        # Step 4 - Backup the current world content
        subprocess.run(['msm', world_name, 'backup'])

        # Step 5 - Remove the current world content
        subprocess.run(['rm', '-rf', f'{self.msm_path}/servers/{world_name}/worldstorage/*'])
        subprocess.run(['rm', '-rf', f'{self.msm_path}/servers/{world_name}/*'])

        # Step 6 - Extract the desired backup to the world folder
        subprocess.run(['unzip' f'{backup_filename}', '-d' f'{self.msm_path}/servers/{world_name}/worldstorage'])

        # Step 7 - Restore directory ownership to the "minecraft" Linux user
        subprocess.run(['chown', '-R', 'minecraft', f'{self.msm_path}/servers/{world_name}/worldstorage/world'])
        subprocess.run(['chgrp', '-R', 'minecraft', f'{self.msm_path}/servers/{world_name}/worldstorage/world'])
