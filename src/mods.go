package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type manifest struct {
	Minecraft struct {
		Version    string `json:"version"`
		ModLoaders []struct {
			ID      string `json:"id"`
			Primary bool   `json:"primary"`
		} `json:"modLoaders"`
	} `json:"minecraft"`
	ManifestType    string                   `json:"manifestType"`
	ManifestVersion int                      `json:"manifestVersion"`
	Name            string                   `json:"name"`
	Version         string                   `json:"version"`
	Author          string                   `json:"author"`
	Overrides       string                   `json:"overrides"`
	Files           []map[string]interface{} `json:"files"`
}

func setupWorkspace() {
	check(os.MkdirAll(filepath.Join(modpackTemp, "overrides", "mods"), 0o755))
	check(os.MkdirAll(filepath.Join(modpackTemp, "overrides", "config"), 0o755))
	check(os.MkdirAll(modDownloads, 0o755))
	check(os.MkdirAll(filepath.Join(tempWorkspace, "mods"), 0o755))

	for _, f := range []string{"pack.toml", "index.toml"} {
		copyFile(f, filepath.Join(tempWorkspace, f))
	}
	filepath.WalkDir("mods", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(path, ".pw.toml") {
			copyFile(path, filepath.Join(tempWorkspace, path))
		}
		return nil
	})
}

func processMods(packName, version, mcVer, loaderID string) (manifest, string, int) {
	manifestData := manifest{
		ManifestType:    "minecraftModpack",
		ManifestVersion: 1,
		Name:            packName,
		Version:         version,
		Author:          "filipjaruska",
		Overrides:       "overrides",
	}
	manifestData.Minecraft.Version = mcVer
	manifestData.Minecraft.ModLoaders = []struct {
		ID      string `json:"id"`
		Primary bool   `json:"primary"`
	}{{ID: loaderID, Primary: true}}

	modListText := fmt.Sprintf("# Mods in %s v%s\n\nThis release includes the following mods:\n\n", packName, version)
	modCount := 0

	entries, _ := os.ReadDir("mods")
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), ".pw.toml") {
			continue
		}
		path := filepath.Join("mods", entry.Name())
		contentBytes, _ := os.ReadFile(path)
		content := string(contentBytes)

		name := parseString(content, `name\s*=\s*"([^"]+)"`)
		if name == "" {
			name = entry.Name()
		}
		modListText += fmt.Sprintf("- %s\n", name)
		modCount++

		if strings.Contains(content, `mode = "metadata:curseforge"`) {
			pid := parseInt(content, `project-id\s*=\s*(\d+)`)
			fid := parseInt(content, `file-id\s*=\s*(\d+)`)
			if pid > 0 && fid > 0 {
				fmt.Printf("Found CurseForge mod: Project ID: %d, File ID: %d (%s)\n", pid, fid, entry.Name())
				manifestData.Files = append(manifestData.Files, map[string]interface{}{"projectID": pid, "fileID": fid, "required": true})
			}
		} else {
			url := parseString(content, `url\s*=\s*"([^"]+)"`)
			filename := parseString(content, `filename\s*=\s*"([^"]+)"`)
			if url != "" && filename != "" {
				fmt.Printf("Downloading %s...\n", filename)
				dest := filepath.Join(modDownloads, filename)
				if download(url, dest) == nil {
					copyFile(dest, filepath.Join(modpackTemp, "overrides", "mods", filename))
				} else {
					fmt.Printf("Failed to download %s\n", filename)
				}
			}
		}
	}
	modListText += fmt.Sprintf("\nTotal: %d mods\n", modCount)
	return manifestData, modListText, modCount
}

func copyOverrides() {
	for _, f := range []string{"pack.toml", "index.toml", "options.txt"} {
		copyFile(f, filepath.Join(modpackTemp, "overrides", f))
	}
	if _, err := os.Stat("config"); err == nil {
		filepath.WalkDir("config", func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				rel, _ := filepath.Rel("config", path)
				copyFile(path, filepath.Join(modpackTemp, "overrides", "config", rel))
			}
			return nil
		})
	}
}
