from datetime import datetime
import glob
import json
import os
from pathlib import Path
import shutil
import subprocess


class AuServerCreationException(Exception):
    pass

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


class MCConfig(object):

    def __init__(self, filepath):
        self.filepath = filepath

        self.config = {}
        with open(filepath, 'r') as file_obj:
            self.file_lines = file_obj.readlines()

            for line in self.file_lines:
                if line.startswith('#'):
                    pass
                else:
                    line = line.rstrip()
                    line_config = line.split('=')
                    self.config[line_config[0]] = line_config[-1]    

    def update_config(self, property, value):
        self.config[property] = value

    def write_config(self, filepath):
        with open(filepath, 'w') as file_obj:
            
            now = datetime.now()
            time_stamp = now.strftime('%a %b %d %X EST %Y')

            file_obj.write('#Minecraft server properties\n')
            file_obj.write(f'#{time_stamp}\n')
            for key, value in self.config.items():
                file_obj.write(f'{key}={value}\n')


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
        msm_path = self.config['msm_path']
       
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


    def create_new_world(self, name, jargroup, version):

        msm_server_path = f"{self.config['msm_path']}/servers"
        time_stamp = datetime.now().strftime('%a %b %d %X EST %Y')

        # 1 - create the world
        subprocess.call(['sudo', 'msm', 'server', 'create', name])
        subprocess.call(['sudo', 'msm', name, 'jar', jargroup])

        # 2 - create the eula.txt file
        eula_file_path = Path(msm_server_path, name, 'eula.txt')
        eula_file = open(eula_file_path, 'w')
        eula_file.write("#By changing the setting below to TRUE you are indicating your agreement to our EULA (https://account.mojang.com/documents/minecraft_eula).\n")
        eula_file.write(f'#{time_stamp}\n')
        eula_file.write('eula=true')
        eula_file.close()
        
        # 3 - update server.properties template and copy to the server folder
        server_properties = MCConfig(self.config['world_config']['server_properties_template'])
        server_properties.update_config('msm-version', f'minecraft/{version}')
        server_properties.update_config('motd', f'Autism Up - {name}')
        server_properties.write_config(f'{msm_server_path}/{name}/server.properties')



        subprocess.call(['sudo', 'msm', name, 'start'])

        for operator in self.config['op_usernames']:
            subprocess.call(['sudo', 'msm', name, 'op', 'add', operator])
        
        subprocess.call(['sudo', 'msm', name, 'stop', 'now'])
        subprocess.call(['sudo', 'msm', name, 'worlds', 'ram', 'world'])

        subprocess.run(['sudo', 'chown', '-R', 'minecraft', f'{msm_server_path}/{name}'])
        subprocess.run(['sudo', 'chgrp', '-R', 'minecraft', f'{msm_server_path}/{name}'])

        print(f'World named "{name}" created')


    def delete_world(self, name):

        msm_archive_path = f"{self.config['msm_path']}/archives"
        subprocess.call(['sudo', 'msm', name, 'backup'])
        
        list_of_files = glob.glob(f'{msm_archive_path}/backups/{name}*')
        latest_file = max(list_of_files, key=os.path.getctime)
        shutil.copy2(latest_file, Path.home())

        subprocess.call(['sudo', 'msm', 'server', 'delete', name])
        subprocess.call(['sudo', 'rm', '-rf', f'{msm_archive_path}/backups/{name}'])
        subprocess.call(['sudo', 'rm', '-rf', f'{msm_archive_path}/logs/{name}'])
        subprocess.call(['sudo', 'rm', '-rf', f'{msm_archive_path}/worlds/{name}'])


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
        subprocess.run(['chown', '-R', 'minecraft', f'{self.msm_path}/servers/{world_name}'])
        subprocess.run(['chgrp', '-R', 'minecraft', f'{self.msm_path}/servers/{world_name}'])
