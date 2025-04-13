# Friends Modpack

A custom Minecraft modpack for Minecraft 1.21.1 on NeoForge.

## Overview

This modpack uses [packwiz](https://packwiz.infra.link/tutorials/creating/getting-started/) for modpack management.

Modpack config is in [pack.toml](/pack.toml).

## Instructions

1. Install Go (1.19 or newer) from https://golang.org/dl/
2. Run in terminal: `go install github.com/packwiz/packwiz@latest`
3. To get the file run: `packwiz curseforge export`

## Mod Management

```bash
# Install from CurseForge
packwiz curseforge install [mod]

# Install from Modrinth
packwiz modrinth install [mod]
```
