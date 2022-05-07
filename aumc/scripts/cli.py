from au import au
import click
import os
import sys

@click.group()
def cli():
    pass


config_file_path = os.environ.get('AU_CONFIG_FILE')
if not config_file_path:
    print('WARNING: AU_CONFIG_FILE environmental variable not set')
    create = input('Do you want to create a new configuration file? (y/n)')
    if create == 'y':
        from au.config_templates import default_config, default_server_properties

        with open('config.json', 'w') as config_file:
            config_file.write(default_config)
        
        with open('server.properties.template', 'w') as default_server_properties_file:
            default_server_properties_file.write(default_server_properties)
        
        print('A default configuration file (config.json) and a default server properties file (server.properties.template) have been created in the current working directory.')
        print('Update the files with desired configurations, place them in the desired location, and set the AU_CONFIG_FILE environmental variable to the path of the config file.')
        sys.exit()
    else:
        sys.exit('Set the AU_CONFIG_FILE environmental variable to the path of the config file.')
        


app = au.AuMc(config_file_path)

@cli.command()
def check_config():
    '''Prints the current configuration to the screen to check it.'''
    click.echo(app.config)


@cli.command()
def reload_config():
    '''Reload the configuration file'''
    
    app.__init__(config_file_path)


@cli.command()
def build_new_jar():
    """Builds the latest version of Spigot Minecraft and copies it to the
       Git repo for publication"""

    app.build_new_jar()


@cli.command()
@click.option('--filename', help='Name of the jarfile to publish with')
def publish_new_jar(filename):
    '''Push the noted jarfile of Minecraft to GitHub and create a new JarGroup in MSM'''

    click.echo('This publishes the new jar to GitHub')


@cli.command()
@click.option('--name', help='Name of the Minecraft server to create')
@click.option('--from_config', is_flag=True, help='All new worlds listed in the configuration file')
@click.option('--jargroup', help='Jargroup to use for the server')
@click.option('--version', help='Version of Minecraft that will be used')
def create_new_world(name, from_config, jargroup, version):
    '''Create a new world using Autism Up default configurations.'''

    if from_config:
        click.echo("Creating worlds defined in the configuration file.")
        for world in app.config['world_config']['world_names']:
            app.create_new_world(world, jargroup, version)
    else:
        click.echo(f"Creating individual world named {name}")
        app.create_new_world(name, jargroup, version)


@cli.command()
@click.option('--name', help='Name of the Minecraft server to delete')
@click.option('--from_config', is_flag=True, help='All new worlds listed in the configuration file')
def delete_world(name, from_config):
    '''Deletes a world using Autism Up configurations'''

    if from_config:
        click.echo("Deleteing all the worlds from the config file.")
        confirm = input("Are you sure? (y/n)")
        if confirm == 'y':
            for world in app.config['world_config']['world_names']:
                app.delete_world(world)
        else:
            sys.exit('Deletion aborted')
    else:
        click.echo(f"Deleteing world named {name}")
        confirm = input("Are you sure? (y/n)")
        if confirm == 'y':
            app.delete_world(name)
 


if __name__ == '__main__':
    cli()