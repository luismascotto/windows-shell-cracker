package main

import "sort"

// shellCommands are Windows Explorer shell: URIs from the cheat sheet (Shell Command column only).
var shellCommands = []string{
	"shell:AccountPictures",
	"shell:AddNewProgramsFolder",
	"shell:Administrative Tools",
	"shell:AppData",
	"shell:AppDataProgramData",
	"shell:Application Shortcuts",
	"shell:AppMods",
	"shell:AppsFolder",
	"shell:AppUpdatesFolder",
	"shell:ChangeRemoveProgramsFolder",
	"shell:Common AppData",
	"shell:Common Start Menu",
	"shell:Common Startup",
	"shell:ConnectionsFolder",
	"shell:ControlPanelFolder",
	"shell:CredentialManager",
	"shell:Desktop",
	"shell:Downloads",
	"shell:Fonts",
	"shell:ImplicitAppShortcuts",
	"shell:InternetFolder",
	"shell:Local AppData",
	"shell:LocalAppDataLow",
	"shell:NetworkPlacesFolder",
	"shell:OneDrive",
	"shell:Personal",
	"shell:PrintersFolder",
	"shell:ProgramFiles",
	"shell:Programs",
	"shell:SendTo",
	"shell:Start Menu",
	"shell:Startup",
	"shell:System",
	"shell:UsersLibrariesFolder",
	"shell:Windows",
}

func sortedShellCommands() []string {
	out := make([]string, len(shellCommands))
	copy(out, shellCommands)
	sort.Strings(out)
	return out
}
