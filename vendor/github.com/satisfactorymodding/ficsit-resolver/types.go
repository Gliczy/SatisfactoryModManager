package resolver

import "context"

type Provider interface {
	ModVersionsWithDependencies(context context.Context, modID string) ([]ModVersion, error)
	GetModName(context context.Context, modReference string) (*ModName, error)
}

type TargetName string

const (
	TargetNameLinuxServer   TargetName = "LinuxServer"
	TargetNameWindows       TargetName = "Windows"
	TargetNameWindowsServer TargetName = "WindowsServer"
)

type ModName struct {
	ID           string `json:"id"`
	ModReference string `json:"mod_reference"`
	Name         string `json:"name"`
}

type VersionDependency struct {
	ModReference string
	Constraint   string
	Optional     bool
}

type ModVersion struct {
	Version          string       `json:"version"`
	GameVersion      string       `json:"game_version"`
	Dependencies     []Dependency `json:"dependencies"`
	Targets          []Target     `json:"targets"`
	RequiredOnRemote bool         `json:"required_on_remote"`
}

type Dependency struct {
	ModID     string `json:"mod_id"`
	Condition string `json:"condition"`
	Optional  bool   `json:"optional"`
}

type Target struct {
	TargetName TargetName `json:"target_name"`
	Link       string     `json:"link"`
	Hash       string     `json:"hash"`
	Size       int64      `json:"size"`
}
