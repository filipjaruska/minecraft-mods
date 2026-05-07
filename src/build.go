package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	manifestName  = "manifest.json"
	modsListName  = "mods_list.txt"
	modpackTemp   = "modpack_temp"
	tempWorkspace = "temp_workspace"
	modDownloads  = "mod_downloads"
)

func main() {
	isServer := flag.Bool("server", false, "Build a server pack instead of a client pack (filters out side=\"client\" mods)")
	flag.Parse()

	version := os.Getenv("VERSION")
	if version == "" {
		version = "local-" + time.Now().Format("20060102")
	}

	packTOML, err := os.ReadFile("pack.toml")
	if err != nil {
		fmt.Println("Error reading pack.toml")
		os.Exit(1)
	}

	content := string(packTOML)
	packName := parseString(content, `name\s*=\s*"([^"]+)"`)
	if packName == "" {
		packName = "Friends"
	}
	mcVer := parseString(content, `minecraft\s*=\s*"([^"]+)"`)
	loaderID := "neoforge-" + parseString(content, `neoforge\s*=\s*"([^"]+)"`)

	// Create zip name representation without spaces
	modpackZipName := strings.ReplaceAll(packName, " ", "")
	
	// Modify names if building for server
	if *isServer {
		packName += " Server"
		modpackZipName += "-Server"
	}

	fmt.Printf("Building %s (%s)\n", packName, version)

	checkPackwiz()

	// Initial cleanup
	cleanup()

	// Ensure cleanup runs when the script finishes (successfully or via panic/check)
	defer func() {
		cleanup()
		fmt.Println("Cleanup complete")
	}()

	// 1. Prepare target folders
	setupWorkspace()

	// 2. Extract configurations and bundle Modrinth jars
	manifestData, modListText, modCount := processMods(packName, version, mcVer, loaderID, *isServer)

	// Export MOD_COUNT to GitHub Actions if running in CI
	if envFile := os.Getenv("GITHUB_ENV"); envFile != "" {
		f, err := os.OpenFile(envFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			f.WriteString(fmt.Sprintf("MOD_COUNT=%d\n", modCount))
			f.Close()
		}
	}

	// 3. Write out the human-readable mod list
	check(os.WriteFile(modsListName, []byte(modListText), 0o644))
	copyFile(modsListName, filepath.Join(modpackTemp, modsListName))

	// 4. Inject configs and other local overrides
	copyOverrides(*isServer)

	// 5. Save the generated manifest.json
	manifestJSON, _ := json.MarshalIndent(manifestData, "", "  ")
	check(os.WriteFile(filepath.Join(modpackTemp, manifestName), append(manifestJSON, '\n'), 0o644))
	fmt.Println("Manifest.json created")

	// 6. Zip into a CurseForge pack
	zipName := fmt.Sprintf("%s-%s.zip", modpackZipName, version)
	check(createZip(modpackTemp, zipName))

	fmt.Printf("Done! Modpack created at: %s\n", zipName)
}

func checkPackwiz() {
	_, err := exec.LookPath("packwiz")
	if err != nil {
		fmt.Println("packwiz not found. Please install it first.\nRun: go install github.com/packwiz/packwiz@latest")
		os.Exit(1)
	}
}

func cleanup() {
	paths := []string{modpackTemp, tempWorkspace, modDownloads}
	// Only delete mods_list.txt if we aren't running in a GitHub Action
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		paths = append(paths, modsListName)
	}
	
	for _, p := range paths {
		os.RemoveAll(p)
	}
}
