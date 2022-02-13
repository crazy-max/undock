package config

import "github.com/alecthomas/kong"

type Cli struct {
	Version kong.VersionFlag

	LogLevel   string `kong:"name=log-level,env=LOG_LEVEL,default=info,help='Set log level.'"`
	LogJSON    bool   `kong:"name=log-json,env=LOG_JSON,default=false,help='Enable JSON logging output.'"`
	LogCaller  bool   `kong:"name=log-caller,env=LOG_CALLER,default=false,help='Add file:line of the caller to log output.'"`
	LogNoColor bool   `kong:"name=log-nocolor,env=LOG_NOCOLOR,default=false,help='Disable colorized output.'"`

	CacheDir string `kong:"name=cachedir,type=path,env=UNDOCK_CACHE_DIR,help='Set cache path. (eg. ~/.local/share/undock/cache)'"`
	Platform string `kong:"name=platform,help='Enforce platform for source image. (eg. linux/amd64)'"`

	All      bool     `kong:"name=all,default=false,help='Extract all architectures if source is a manifest list.'"`
	Includes []string `kong:"name=include,help='Include a subset of files/dirs from the source image.'"`
	Insecure bool     `kong:"name=insecure,default=false,help='Allow contacting the registry or docker daemon over HTTP, or HTTPS with failed TLS verification.'"`
	RmDist   bool     `kong:"name=rm-dist,default=false,help='Removes dist folder.'"`
	Wrap     bool     `kong:"name=wrap,default=false,help='For a manifest list, merge output in dist folder.'"`

	Source string `kong:"arg,required,name=source,help='Source image. (eg. alpine:latest)'"`
	Dist   string `kong:"arg,required,name=dist,type=path,help='Dist folder. (eg. ./dist)'"`
}
