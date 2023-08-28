List and download files from Tor sites that use ALL_FILES

# Installation
`go install github.com/juanmera/allfiles/cmd/allfiles@latest`

# Examples
## Listing files
`allfiles ls -f ./ALL_FILES -totals -exts pdf,docx`

## Downloading files
`allfiles dl -f ./ALL_FILES -proxy "socks5://127.0.0.1:9050" -u "https://abc.onion/123/data" -outputdir "./data" -types document,bitmap`

# Config
All named parameters, can be set in a JSON file named "allfiles.json" in the current directory.
See allfiles.json.example

# Types
The types and exclude-types parameters support the following values with their respective file extensions:
* binary:   ".exe", ".dll"
* package:  ".cab", ".msi", ".pkg", ".deb", ".rpm", ".dmg", ".gpd", ".inf"
* archive:  ".zip", ".rar"
* video:    ".mov", ".avi", ".mp4", ".mkv", ".mpg", ".wmv"
* audio:    ".mp3"
* bitmap:   ".cr2", ".bmp", ".gif", ".jpg", ".jpeg", ".png", ".tif", ".tiff", ".ico", ".mdi"
* cad:      ".dwg", ".lcf"
* font:     ".otf", ".ttf"
* document: ".pdf", ".xls", ".xlsx", ".doc", ".docx", ".ppt", ".pptx", ".odt", ".ods"
* text:     ".txt", ".log", ".csv", ".xml"
* source:   ".php", ".c"
* database: ".accdb", ".mdb", ".sql", ".sqlite"
* config:   ".ini", ".conf", ".yml"
* others:   ".msg", ".tmp", ".lnk"