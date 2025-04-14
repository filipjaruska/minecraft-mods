# Friends Modpack

A custom Minecraft modpack for Minecraft 1.21.1 on NeoForge.

## Download Latest Release

The latest release can be downloaded from [here](https://github.com/filipjaruska/minecraft-mods/releases/latest).

## Overview

This modpack uses [packwiz](https://packwiz.infra.link/tutorials/creating/getting-started/) for modpack management.

## Development Instructions

1. Install Go (1.19 or newer) from https://golang.org/dl/
2. Run in terminal: `go install github.com/packwiz/packwiz@latest`

### Building the Modpack

Run `.\build-local.ps1` to build the game locally

## Mod Management

```bash
# Install from CurseForge
packwiz curseforge install [mod]

# Install from Modrinth
packwiz modrinth install [mod]
```
