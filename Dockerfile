FROM ubuntu:20.04

# Prevent interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

# Set up working directory
WORKDIR /opt

# Update system packages
RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get install -y \
    screen \
    rsync \
    zip \
    jq \
    wget \
    curl \
    git \
    sudo \
    python3-pip \
    python-is-python3 \
    && rm -rf /var/lib/apt/lists/*

# Install OpenJDK 17 (more reliable than Oracle JDK)
RUN apt-get update && \
    apt-get install -y openjdk-17-jdk-headless && \
    rm -rf /var/lib/apt/lists/*

# Create minecraft user
RUN useradd minecraft --system --home /opt/msm && \
    mkdir -p /opt/msm && \
    chown -R minecraft:minecraft /opt/msm

# Install MSM (Minecraft Server Manager)
ENV UPDATE_URL="https://raw.githubusercontent.com/msmhq/msm/master"
RUN mkdir -p /tmp/msm-install && \
    cd /tmp/msm-install && \
    wget -q ${UPDATE_URL}/msm.conf -O msm.conf.orig && \
    wget -q ${UPDATE_URL}/cron/msm -O msm.cron.orig && \
    wget -q ${UPDATE_URL}/init/msm -O msm.init.orig && \
    sed 's#USERNAME="minecraft"#USERNAME="minecraft"#g' msm.conf.orig | \
        sed "s#/opt/msm#/opt/msm#g" | \
        sed "s#UPDATE_URL=.*\$#UPDATE_URL=\"$UPDATE_URL\"#" > msm.conf && \
    install -b -m0644 msm.conf /etc/msm.conf && \
    install -b msm.init.orig /etc/init.d/msm && \
    ln -s /etc/init.d/msm /usr/local/bin/msm && \
    chmod +x /etc/init.d/msm && \
    /etc/init.d/msm update --noinput && \
    /etc/init.d/msm jargroup create minecraft minecraft && \
    rm -rf /tmp/msm-install

# Install Spigot BuildTools
RUN mkdir /opt/build_tools && \
    curl -o /opt/build_tools/BuildTools.jar https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar && \
    chown -R minecraft:minecraft /opt/build_tools

# Create test user
RUN useradd -m -s /bin/bash testuser && \
    usermod -aG sudo testuser && \
    usermod -aG minecraft testuser && \
    echo "testuser ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Configure Git
RUN git config --global user.name "Test User" && \
    git config --global user.email "test@example.com"

# Expose Minecraft port
EXPOSE 25565

# Set up entrypoint
WORKDIR /workspace
CMD ["/bin/bash"]
