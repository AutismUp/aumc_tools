import jarmanager as jm
import click

@click.group()
def cli():
    pass

@cli.command()
@click.option('--version', default='latest', help='version of Minecraft to get')
def build_new_jar(version):
    """Builds the latest version of Spigot Minecraft and copies it to the
       Git repo for publication"""

    jm.build_new_jar(version=version)


@cli.command()
@click.option('--filename', help='Name of the jarfile to publish with')
def publish_new_jar(filename):
    '''Push the noted jarfile of Minecraft to GitHub and create a new JarGroup in MSM'''

    jm.publish_new_jar(filename)


if __name__ == '__main__':
    cli()