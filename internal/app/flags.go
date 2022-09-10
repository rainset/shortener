package app

import "flag"

func (a *App) regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}

func (a *App) getStringFlag(name string) string {
	return flag.Lookup(name).Value.(flag.Getter).Get().(string)
}

func (a *App) InitFlags() {
	var flagA, flagB, flagF, flagD string
	a.regStringVar(&flagA, "a", a.Config.ServerAddress, "set env SERVER_ADDRESS")
	a.regStringVar(&flagB, "b", a.Config.ServerBaseURL, "set env BASE_URL")
	a.regStringVar(&flagF, "f", a.Config.ServerFileStoragePath, "set env FILE_STORAGE_PATH")
	a.regStringVar(&flagD, "d", a.Config.DatabaseDSN, "set env DATABASE_DSN")
	flag.Parse()
	a.Config.ServerAddress = a.getStringFlag("a")
	a.Config.ServerBaseURL = a.getStringFlag("b")
	a.Config.ServerFileStoragePath = a.getStringFlag("f")
	a.Config.DatabaseDSN = a.getStringFlag("d")
}
