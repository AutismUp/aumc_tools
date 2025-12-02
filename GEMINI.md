# Project Overview

This project, the Autism Up Minecraft Tool (`aumc`), is a command-line utility designed to simplify the management of Minecraft servers for Autism Up. It is a hybrid project that leverages both Go and Python. The core server management functionalities are being developed in Go, while the user-facing CLI and configuration management are handled by Python.

The tool is designed to be installed as a Python package and provides a command-line interface for initializing the configuration, managing the Minecraft server, and interacting with the underlying Go utilities.

The project also includes a `Vagrantfile` to create a consistent and reproducible development environment using VirtualBox.

## Building and Running

### Python CLI

The primary way to use the tool is through the Python CLI.

**Installation:**

```bash
pip install git+https://github.com/AutismUp/aumc_tools.git
```

**Initialization:**

After installation, run the `au` command to initialize the configuration:

```bash
au
```

This will create a `config.json` and `server.properties` file.

**Running Commands:**

The Python CLI provides the following commands:

*   `au init`: Initializes the configuration files.
*   `au msm`: Interacts with the Minecraft Server Manager (MSM).
*   `au server`: Manages the Minecraft server.

### Go CLI

The Go part of the project provides the core functionalities.

**Building:**

To build the Go binary, run the following command from the root of the project:

```bash
go build -o aumc_go
```

**Running Commands:**

The Go CLI provides the following commands:

*   `./aumc_go check-config`: Displays the current configuration.
*   `./aumc_go reload-config`: Reloads the configuration file.
*   `./aumc_go create-new-world`: Creates a new Minecraft world.
*   `./aumc_go delete-world`: Deletes a Minecraft world.

**Note:** The world creation and deletion functionalities in the Go CLI are not yet fully implemented.

### Development Environment

To set up a development environment, you need to have VirtualBox and Vagrant installed.

1.  Clone the repository:

    ```bash
    git clone git@github.com:AutismUp/aumc_tools.git
    ```

2.  Change into the project directory:

    ```bash
    cd aumc_tools
    ```

3.  Start the Vagrant environment:

    ```bash
    vagrant up
    ```

4.  SSH into the Vagrant box:

    ```bash
    vagrant ssh
    ```

## Development Conventions

*   The project uses Go for the core logic and Python for the user-facing CLI.
*   The Go CLI is built using the [Cobra](https://github.com/spf13/cobra) library.
*   The Python CLI is built using the [Click](https://click.palletsprojects.com/en/8.1.x/) library.
*   Configuration is managed through a `config.json` file.
*   The project uses a `Vagrantfile` to define the development environment.
