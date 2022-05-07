# Autism Up Minecraft Tool

## Overview

The Autism Up Minecraft Tool (aumc) is a wrapper script around the various Minecraft utilities needed to
support AU's specific Minecraft configuration. Though not specifically required to support
Minecraft, the tool makes creating, updating and support simpler.

## Installation and setup

To install:

```bash

pip install aumc

```

After installing, run the command line to generate a default configuration file, and a default Minecraft server.properties file:

```
au

WARNING: AU_CONFIG_FILE environmental variable not set
Do you want to create a new configuration file? (y/n)

```

When "y" is entered, two files will be created: config.json and server.properties

Place them in a desired location, update config.json file to reference the location of the server.properties file, and set the AU_CONFIG_FILE environmental variable to be equal to the path of the config.json file.

You're ready to start using aumc!


## Test and development setup

To test Autism Up Minecraft environment and develop the tool, a development environment should
be created on your local workstation. To setup your workstation to enable this, you'll need a 
few extra pieces of software:

1. [**VirtualBox**](https://www.virtualbox.org/wiki/Downloads) - This is Oracle's open source computer virtualization software. It will allow you to create a virtual minecraft server on your computer.
2. [**Vagrant**](https://www.vagrantup.com/downloads) - Vagrant is an automation tool that will automatically configure the VirtualBox computer the same way the live Minecraft server will be configured.

Download and install these tools before getting started.

After installing VirtualBox and Vagrant, open your command line console, change to a directory where you want to keep your development environment, and run the following commands:

```bash

git clone git@github.com:AutismUp/aumc_tools.git
cd aumc_tools
vagrant up

```

After Vagrant completes the build of the local test server, log on to the server by SSH'ing in using Vagrant:

```bash

vagrant ssh

```

Finally, in your Minecraft multiplayer serve list add an additional server:

Server Name: AU Test and Development
Serever Address: localhost

You're now ready to test and develop the Autism Up Minecraft experience!!
