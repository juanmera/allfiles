package extensions

var byTypes = map[string][]string{
	"binary":   {".exe", ".dll"},                                                                   // Windows binaries
	"package":  {".cab", ".msi", ".pkg", ".deb", ".rpm", ".dmg", ".gpd", ".inf"},                   // App packages
	"archive":  {".zip", ".rar"},                                                                   // Compressed
	"video":    {".mov", ".avi", ".mp4", ".mkv", ".mpg", ".wmv"},                                   // Video
	"audio":    {".mp3"},                                                                           // Audio
	"bitmap":   {".cr2", ".bmp", ".gif", ".jpg", ".jpeg", ".png", ".tif", ".tiff", ".ico", ".mdi"}, // Images (.mdi is a MS scanned document)
	"cad":      {".dwg", ".lcf"},                                                                   // AutoCAD
	"font":     {".otf", ".ttf"},                                                                   // Fonts
	"document": {".pdf", ".xls", ".xlsx", ".doc", ".docx", ".ppt", ".pptx", ".odt", ".ods"},        // Documents
	"text":     {".txt", ".log", ".csv", ".xml"},                                                   // Logs
	"source":   {".php", ".c"},                                                                     // Source code
	"database": {".accdb", ".mdb", ".sql", ".sqlite"},                                              // Databases
	"config":   {".ini", ".conf", ".yml"},                                                          // Configuration
	"others":   {".msg", ".tmp", ".lnk"},                                                           // Others
}

func Types() []string {
	types := make([]string, len(byTypes))
	i := 0
	for k := range byTypes {
		types[i] = k
		i++
	}
	return types
}

func Get(types ...string) []string {
	extensions := make([]string, 0, 2)
	for _, v := range types {
		if byType, ok := byTypes[v]; ok {
			extensions = append(extensions, byType...)
		}
	}
	return extensions
}
