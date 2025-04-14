# Local build script for Minecraft modpack
# This script mimics the GitHub Actions workflow but runs locally on Windows

# Configuration
$VERSION = "local-" + (Get-Date -Format "yyyyMMdd")
$MODPACK_NAME = "FriendsModpack"
$TEMP_WORKSPACE = "temp_workspace"

# Check if packwiz is installed
try {
    $packwizVersion = packwiz --version
    Write-Host "Using Packwiz: $packwizVersion"
} catch {
    Write-Host "Error: packwiz not found. Please install it first."
    Write-Host "Run: go install github.com/packwiz/packwiz@latest"
    exit 1
}

# Create temp directories
Write-Host "Creating temporary directories..."
$modpackTemp = "modpack_temp"

# Clean up any existing temporary directories
if (Test-Path -Path $modpackTemp) {
    Remove-Item -Path $modpackTemp -Recurse -Force
}
if (Test-Path -Path $TEMP_WORKSPACE) {
    Remove-Item -Path $TEMP_WORKSPACE -Recurse -Force
}
if (Test-Path -Path "mod_downloads") {
    Remove-Item -Path "mod_downloads" -Recurse -Force
}

# Create temporary directories
New-Item -ItemType Directory -Force -Path $modpackTemp
New-Item -ItemType Directory -Force -Path "$modpackTemp\overrides\mods"
New-Item -ItemType Directory -Force -Path "$modpackTemp\overrides\config"
New-Item -ItemType Directory -Force -Path "mod_downloads"

# Create a temporary workspace to avoid modifying original files
Write-Host "Creating temporary workspace..."
New-Item -ItemType Directory -Force -Path $TEMP_WORKSPACE
New-Item -ItemType Directory -Force -Path "$TEMP_WORKSPACE\mods"
Copy-Item -Path "pack.toml" -Destination $TEMP_WORKSPACE
Copy-Item -Path "index.toml" -Destination $TEMP_WORKSPACE
Copy-Item -Path "mods\*.pw.toml" -Destination "$TEMP_WORKSPACE\mods\"

# Generate mods list directly from original files
Write-Host "Generating mods list..."
$modsList = "# Mods in $MODPACK_NAME v$VERSION`r`n`r`nThis release includes the following mods:`r`n`r`n"
$modCount = 0

Get-ChildItem -Path "mods" -Filter "*.pw.toml" | ForEach-Object {
    $modFile = Get-Content $_.FullName -Raw
    
    # Extract mod name from the file
    if ($modFile -match 'name\s*=\s*"([^"]+)"') {
        $modName = $matches[1]
    } else {
        $modName = $_.BaseName
    }
    
    $modsList += "- $modName`r`n"
    $modCount++
}

$modsList += "`r`nTotal: $modCount mods"
Set-Content -Path "mods_list.txt" -Value $modsList
Copy-Item -Path "mods_list.txt" -Destination "$modpackTemp\"
Write-Host "Mods list generated with $modCount mods"

# Copy base files to overrides
Copy-Item -Path "pack.toml" -Destination "$modpackTemp\overrides\"
Copy-Item -Path "index.toml" -Destination "$modpackTemp\overrides\"
Copy-Item -Path "options.txt" -Destination "$modpackTemp\overrides\"

# Copy config files
if (Test-Path -Path "config") {
    Write-Host "Copying configuration files..."
    Copy-Item -Path "config\*" -Destination "$modpackTemp\overrides\config\" -Recurse
}

# Create manifest.json
Write-Host "Creating manifest.json..."
$manifestContent = @"
{
  "minecraft": {
    "version": "1.21.1",
    "modLoaders": [
      {
        "id": "neoforge-21.1.147",
        "primary": true
      }
    ]
  },
  "manifestType": "minecraftModpack",
  "manifestVersion": 1,
  "name": "Friends Modpack",
  "version": "$VERSION",
  "author": "filipjaruska",
  "overrides": "overrides",
  "files": [
"@

$firstMod = $true

# Process each mod file
Write-Host "Processing mod files..."
Get-ChildItem -Path "mods" -Filter "*.pw.toml" | ForEach-Object {
    $modFile = Get-Content $_.FullName -Raw
    
    if ($modFile -match 'mode = "metadata:curseforge"') {
        # Extract CurseForge project ID and file ID
        $projectId = if ($modFile -match 'project-id\s*=\s*(\d+)') { $matches[1] } else { $null }
        $fileId = if ($modFile -match 'file-id\s*=\s*(\d+)') { $matches[1] } else { $null }
        
        if ($projectId -and $fileId) {
            Write-Host "Found CurseForge mod: Project ID: $projectId, File ID: $fileId"
            
            # Add comma for all mods except the first one
            if (-not $firstMod) {
                $manifestContent += ","
            } else {
                $firstMod = $false
            }
            
            # Add the mod to manifest.json
            $manifestContent += @"

    {
      "projectID": $projectId,
      "fileID": $fileId,
      "required": true
    }
"@
        }
    } else {
        # Non-CurseForge mod (like Modrinth)
        Write-Host "Processing non-CurseForge mod: $($_.Name)"
        
        # Extract download URL and filename
        $url = if ($modFile -match 'url\s*=\s*"([^"]+)"') { $matches[1] } else { $null }
        $filename = if ($modFile -match 'filename\s*=\s*"([^"]+)"') { $matches[1] } else { $null }
        
        if ($url -and $filename) {
            Write-Host "Downloading $filename from $url"
            
            try {
                # Download the actual mod .jar file
                $downloadPath = "mod_downloads\$filename"
                Invoke-WebRequest -Uri $url -OutFile $downloadPath
                
                if (Test-Path -Path $downloadPath) {
                    Write-Host "Successfully downloaded $filename"
                    Copy-Item -Path $downloadPath -Destination "$modpackTemp\overrides\mods\"
                } else {
                    Write-Host "Failed to download $filename"
                }
            } catch {
                Write-Host "Error downloading $filename`: $_"
            }
        } else {
            Write-Host "Warning: Could not extract download information from $($_.Name)"
        }
    }
}

# Close the JSON structure
$manifestContent += @"

  ]
}
"@

Set-Content -Path "$modpackTemp\manifest.json" -Value $manifestContent
Write-Host "Manifest.json created"

# Create ZIP file
Write-Host "Creating modpack ZIP..."
$zipFilePath = "$MODPACK_NAME-$VERSION.zip"

# Remove existing file if it exists
if (Test-Path $zipFilePath) {
    Remove-Item -Path $zipFilePath -Force
}

# Create the ZIP file
Add-Type -AssemblyName System.IO.Compression.FileSystem
[System.IO.Compression.ZipFile]::CreateFromDirectory($modpackTemp, $zipFilePath)

Write-Host "Done! Modpack created at: $zipFilePath"
Write-Host "You can import this file into CurseForge launcher."

# Clean up
Write-Host "Cleaning up temporary files..."
Remove-Item -Path $modpackTemp -Recurse -Force
Remove-Item -Path "mod_downloads" -Recurse -Force
Remove-Item -Path $TEMP_WORKSPACE -Recurse -Force
Remove-Item -Path "mods_list.txt" -Recurse -Force