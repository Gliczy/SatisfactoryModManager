package resolver

import (
	"context"
	"errors"
)

var _ Provider = (*MockProvider)(nil)

type MockProvider struct{}

var commonTargets = []Target{
	{
		TargetName: "Windows",
		Hash:       "698df20278b3de3ec30405569a22050c6721cc682389312258c14948bd8f38ae",
	},
	{
		TargetName: "WindowsServer",
		Hash:       "7be01ed372e0cf3287a04f5cb32bb9dcf6f6e7a5b7603b7e43669ec4c6c1457f",
	},
	{
		TargetName: "LinuxServer",
		Hash:       "bdbd4cb1b472a5316621939ae2fe270fd0e3c0f0a75666a9cbe74ff1313c3663",
	},
}

func (m MockProvider) ModVersionsWithDependencies(_ context.Context, modID string) ([]ModVersion, error) {
	sml3 := Dependency{
		ModID:     "SML",
		Condition: "^3.6.0",
		Optional:  false,
	}

	switch modID {
	case "RefinedPower":
		return []ModVersion{
			{
				Version:          "3.2.13",
				RequiredOnRemote: true,
				Dependencies: []Dependency{
					{
						ModID:     "ModularUI",
						Condition: "^2.1.11",
						Optional:  false,
					},
					{
						ModID:     "RefinedRDLib",
						Condition: "^1.1.7",
						Optional:  false,
					},
					{
						ModID:     "SML",
						Condition: "^3.6.1",
						Optional:  false,
					},
				},
				Targets: commonTargets,
			},
			{
				Version:          "3.2.11",
				RequiredOnRemote: true,
				Dependencies: []Dependency{
					{
						ModID:     "ModularUI",
						Condition: "^2.1.10",
						Optional:  false,
					},
					{
						ModID:     "RefinedRDLib",
						Condition: "^1.1.6",
						Optional:  false,
					},
					sml3,
				},
				Targets: commonTargets,
			},
			{
				Version:          "3.2.10",
				RequiredOnRemote: true,
				Dependencies: []Dependency{
					{
						ModID:     "ModularUI",
						Condition: "^2.1.9",
						Optional:  false,
					},
					{
						ModID:     "RefinedRDLib",
						Condition: "^1.1.5",
						Optional:  false,
					},
					sml3,
				},
				Targets: commonTargets,
			},
		}, nil
	case "RefinedRDLib":
		return []ModVersion{
			{
				Version:          "1.1.7",
				RequiredOnRemote: true,
				Dependencies: []Dependency{
					{
						ModID:     "SML",
						Condition: "^3.6.1",
						Optional:  false,
					},
				},
				Targets: commonTargets,
			},
			{
				Version:          "1.1.6",
				RequiredOnRemote: true,
				Dependencies:     []Dependency{sml3},
				Targets:          commonTargets,
			},
			{
				Version:          "1.1.5",
				RequiredOnRemote: true,
				Dependencies:     []Dependency{sml3},
				Targets:          commonTargets,
			},
		}, nil
	case "ModularUI":
		return []ModVersion{
			{
				Version:          "2.1.12",
				RequiredOnRemote: true,
				Dependencies: []Dependency{
					{
						ModID:     "SML",
						Condition: "^3.6.1",
						Optional:  false,
					},
				},
				Targets: commonTargets,
			},
			{
				Version:          "2.1.11",
				RequiredOnRemote: true,
				Dependencies:     []Dependency{sml3},
				Targets:          commonTargets,
			},
			{
				Version:          "2.1.10",
				RequiredOnRemote: true,
				Dependencies:     []Dependency{sml3},
				Targets:          commonTargets,
			},
		}, nil
	case "ThisModDoesNotExist$$$":
		return []ModVersion{}, errors.New("mod not found")
	case "ComplexMod":
		return []ModVersion{
			{
				Version:          "3.0.0",
				RequiredOnRemote: true,
				Dependencies:     []Dependency{sml3},
				Targets: []Target{
					{
						TargetName: "LinuxServer",
						Hash:       "8739c76e681f900923b900c9df0ef75cf421d39cabb54650c4b9ad19b6a76d85",
					},
				},
			},
			{
				Version:          "2.0.0",
				RequiredOnRemote: true,
				Dependencies:     []Dependency{sml3},
				Targets:          commonTargets,
			},
			{
				Version:          "1.0.0",
				RequiredOnRemote: true,
				Dependencies:     []Dependency{sml3},
				Targets: []Target{
					{
						TargetName: "Windows",
						Hash:       "62f5c84eca8480b3ffe7d6c90f759e3b463f482530e27d854fd48624fdd3acc9",
					},
				},
			},
		}, nil
	case "SML":
		return []ModVersion{
			{
				Version:          "2.2.1",
				GameVersion:      ">=125236",
				RequiredOnRemote: true,
				Targets:          []Target{},
			},
			{
				Version:          "3.3.2",
				GameVersion:      ">=194714",
				RequiredOnRemote: true,
				Targets: []Target{
					{
						TargetName: TargetNameWindows,
						Hash:       "unknown",
					},
				},
			},
			{
				Version:          "3.6.0",
				GameVersion:      ">=264901",
				RequiredOnRemote: true,
				Targets: []Target{
					{
						TargetName: TargetNameWindows,
						Hash:       "unknown",
					},
					{
						TargetName: TargetNameWindowsServer,
						Hash:       "unknown",
					},
					{
						TargetName: TargetNameLinuxServer,
						Hash:       "unknown",
					},
				},
			},
			{
				Version:          "3.6.1",
				GameVersion:      ">=264901",
				RequiredOnRemote: true,
				Targets: []Target{
					{
						TargetName: TargetNameWindows,
						Hash:       "unknown",
					},
					{
						TargetName: TargetNameWindowsServer,
						Hash:       "unknown",
					},
					{
						TargetName: TargetNameLinuxServer,
						Hash:       "unknown",
					},
				},
			},
		}, nil
	case "ClientOnlyMod":
		return []ModVersion{
			{
				Version:          "1.0.0",
				RequiredOnRemote: false,
				Targets: []Target{
					{
						TargetName: "Windows",
						Hash:       "8739c76e681f900923b900c9df0ef75cf421d39cabb54650c4b9ad19b6a76d85",
					},
				},
			},
		}, nil
	case "ServerOnlyMod":
		return []ModVersion{
			{
				Version:          "2.0.0",
				RequiredOnRemote: false,
				Targets: []Target{
					{
						TargetName: "WindowsServer",
						Hash:       "8739c76e681f900923b900c9df0ef75cf421d39cabb54650c4b9ad19b6a76d85",
					},
				},
			},
			{
				Version:          "1.0.0",
				RequiredOnRemote: false,
				Targets: []Target{
					{
						TargetName: "WindowsServer",
						Hash:       "8739c76e681f900923b900c9df0ef75cf421d39cabb54650c4b9ad19b6a76d85",
					},
					{
						TargetName: "LinuxServer",
						Hash:       "8739c76e681f900923b900c9df0ef75cf421d39cabb54650c4b9ad19b6a76d85",
					},
				},
			},
		}, nil
	}

	panic("ModVersionsWithDependencies: " + modID)
}

func (m MockProvider) GetModName(_ context.Context, modReference string) (*ModName, error) {
	switch modReference {
	case "RefinedPower":
		return &ModName{
			ID:           "DGiLzB3ZErWu2V",
			ModReference: "RefinedPower",
			Name:         "Refined Power",
		}, nil
	case "RefinedRDLib":
		return &ModName{
			ID:           "B24emzbs6xVZQr",
			ModReference: "RefinedRDLib",
			Name:         "RefinedRDLib",
		}, nil
	case "ComplexMod":
		return &ModName{
			ID:           "asd32rfewqhy4",
			ModReference: "ComplexMod",
			Name:         "ComplexMod",
		}, nil
	case "ClientOnlyMod":
		return &ModName{
			ID:           "asd32rfewqhy4",
			ModReference: "ClientOnlyMod",
			Name:         "ClientOnlyMod",
		}, nil
	case "ServerOnlyMod":
		return &ModName{
			ID:           "asd32rfewqhy4",
			ModReference: "ServerOnlyMod",
			Name:         "ServerOnlyMod",
		}, nil
	case "SML":
		return &ModName{
			ID:           "SML",
			ModReference: "SML",
			Name:         "Satisfactory Mod Loader",
		}, nil
	}

	panic("GetModName: " + modReference)
}
