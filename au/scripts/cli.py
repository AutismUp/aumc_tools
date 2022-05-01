from au import au
import click
import os
import sys

@click.group()
def cli():
    pass


config_file_path = os.environ.get('AU_CONFIG_FILE')
if not config_file_path:
    sys.exit('ERROR: AU_CONFIG_FILE environmental variable not set')

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
@click.option('--jargroup', help='Jargroup to use for the server')
@click.option('--version', help='Version of Minecraft that will be used')
def create_new_world(name, jargroup, version):
    '''Create a new world using Autism Up default configurations.'''

    app.create_new_world(name,jargroup, version)


if __name__ == '__main__':
    cli()