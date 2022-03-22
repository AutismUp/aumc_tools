from datetime import datetime

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

