package constant

const (
	DOWNLOAD_HELP_TEMPLATE = `📦 Download WordNet Japan sqlite3 database file from the official site.

You have to download WordNet Japan sqlite3 database file to use jrp at first.
jrp will download archive file from the official site and decompress it to the database file.

You can set the directory of the database file to the environment variable "JRP_WNJPN_DB_FILE_DIR".
The default directory is "~/.local/share/jrp" ("$XDG_DATA_HOME/jrp").

Usage:
  jrp download [flags]
  jrp dl       [flags]
  jrp d        [flags]

Flags:
  -h, --help   🤝 help for download
`
	DOWNLOAD_USE   = "download"
	DOWNLOAD_SHORT = "📦 Download WordNet Japan sqlite3 database file from the official site."
	DOWNLOAD_LONG  = `📦 Download WordNet Japan sqlite3 database file from the official site.

You have to download WordNet Japan sqlite3 database file to use jrp at first.
jrp will download archive file from the official site and decompress it to the database file.

You can set the directory of the database file to the environment variable "JRP_WNJPN_DB_FILE_DIR".
The default directory is "$XDG_DATA_HOME/jrp".
`
	DOWNLOAD_MESSAGE_DOWNLOADING        = "  📦 Downloading WordNet Japan sqlite3 database file from the official site..."
	DOWNLOAD_MESSAGE_SUCCEEDED          = "✅ Downloaded successfully! Now, you are ready to use jrp!"
	DOWNLOAD_MESSAGE_FAILED             = "❌ Failed to download... Please try again later..."
	DOWNLOAD_MESSAGE_ALREADY_DOWNLOADED = "✅ You are already ready to use jrp!"
)

func GetDownloadAliases() []string {
	return []string{"dl", "d"}
}