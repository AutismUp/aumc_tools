default_config = """{
    "msm_path": "",
    "build_config": {
        "build_directory": "",
        "temp_folders": [
            "BuildData",
            "Bukkit",
            "CraftBukkit",
            "Spigot",
            "apache-maven-3.6.0",
            "work"
        ],
        "temp_files": [
            "BuildTools.log.txt"
        ],
        "minecraft_version": "",
        "delete_spigot_jars": true,
        "jar_git_repo": ""
    },
    "world_config": {
        "server_properties_template": "",
        "world_names": [
            "world1",
            "world2"
        ]
    },
    "op_usernames": [
        "op1",
        "op2"
    ] 
}
"""

default_server_properties = """#Minecraft server properties
#Wed Dec 30 16:44:40 EST 2020
msm-version=minecraft/1.18.2
spawn-protection=16
max-tick-time=60000
query.port=25565
generator-settings=
sync-chunk-writes=true
force-gamemode=false
allow-nether=true
enforce-whitelist=true
gamemode=creative
broadcast-console-to-ops=true
enable-query=false
player-idle-timeout=0
text-filtering-config=
difficulty=easy
spawn-monsters=true
broadcast-rcon-to-ops=true
op-permission-level=4
pvp=true
entity-broadcast-range-percentage=100
snooper-enabled=true
level-type=default
hardcore=false
enable-status=true
enable-command-block=false
max-players=20
network-compression-threshold=256
resource-pack-sha1=
max-world-size=29999984
function-permission-level=2
rcon.port=25575
server-port=25565
debug=false
server-ip=
spawn-npcs=true
allow-flight=false
level-name=world
view-distance=10
resource-pack=
spawn-animals=true
white-list=true
rcon.password=
generate-structures=true
max-build-height=256
online-mode=true
level-seed=
use-native-transport=true
prevent-proxy-connections=false
enable-jmx-monitoring=false
enable-rcon=false
rate-limit=0
motd=AU
"""