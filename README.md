# Friends Modpack

A custom Minecraft modpack for Minecraft 1.21.1 on NeoForge.

## Overview

This modpack uses [packwiz](https://packwiz.infra.link/tutorials/creating/getting-started/) for modpack management.

Modpack config is in [pack.toml](/pack.toml).

## Development Instructions

1. Install Go (1.19 or newer) from https://golang.org/dl/
2. Run in terminal: `go install github.com/packwiz/packwiz@latest`

## Mod Management

```bash
# Install from CurseForge
packwiz curseforge install [mod]

# Install from Modrinth
packwiz modrinth install [mod]
```

## CurseForge Distribution

To export the modpack for CurseForge:

```bash
packwiz curseforge export
```

This will create a `.zip` file that can be imported directly into the CurseForge app.
