# Friends Modpack

[Note] Keep master vanilla ish with qol mods.

## Download Latest Release

The latest release can be downloaded from [here](https://github.com/filipjaruska/minecraft-mods/releases/latest).

## Overview

This modpack uses [packwiz](https://packwiz.infra.link/tutorials/creating/getting-started/) for mod management.

## Development Instructions

1. Install Go (1.19 or newer) from https://golang.org/dl/
2. Run in terminal: `go install github.com/packwiz/packwiz@latest`

3. Add mods to `mods.toml` using packwiz:

```bash
# Install from CurseForge
packwiz curseforge install [mod]

# Install from Modrinth
packwiz modrinth install [mod]
```

4. Build the game locally by running`go run ./src`.
