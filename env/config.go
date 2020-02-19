package env

import log "github.com/sirupsen/logrus"

type Config struct {
	PrettyLogOutput bool
	LogLevel        string
	SourceFile      string
	OutputMode      string
}

func (c *Config) Dump() {
	log.Printf("\n ------- CONFIG DUMP ------ \n"+
		"\tPrettyLogOutput: %t\n"+
		"\tLogLevel: %s\n"+
		"\tSourceFile: %s\n"+
		"\tOutputMode: %s\n"+
		"-----------------------",
		c.PrettyLogOutput,
		c.LogLevel,
		c.SourceFile,
		c.OutputMode,
	)
}
