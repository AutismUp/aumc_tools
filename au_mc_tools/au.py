from difflib import restore
import subprocess
import os

class AuRestoreException(Exception):
    pass


class AuMc(object):

    def __init__(self):
        self.msm_path = "/mnt/volumen_nyc1_03/msm"
        self.backup_world_path = f"{self.msm_path}/archives/worlds/experiment"


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
