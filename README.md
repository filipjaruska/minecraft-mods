# Friends Modpack

Vanilla~ish with qol mods. The latest release can be downloaded from [here](https://github.com/filipjaruska/minecraft-mods/releases/latest).

## Self hosting

todo

## Development Instructions

This modpack uses [packwiz](https://packwiz.infra.link/tutorials/creating/getting-started/) for mod management.

> Modpack already contains Forgified Fabric API and Sinytra Connector, no need to worry about compatibility, its possible to install Fabric only mods without worrying about compatibility with Forge mods.

1. Install Go (1.19 or newer) from https://golang.org/dl/
2. Run in terminal: `go install github.com/packwiz/packwiz@latest`

3. Add new mods using packwiz:

```bash
packwiz curseforge install [mod] # Install from CurseForge

packwiz modrinth install [mod] # Install from Modrinth
```

For new GUI only mods run: `packwiz mod side "name of mod" client`

4. Build the game locally by running`go run ./src` or `go run ./src --server` for server only build.
