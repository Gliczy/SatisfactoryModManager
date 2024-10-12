package resolver

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/mircearoata/pubgrub-go/pubgrub"
	"github.com/mircearoata/pubgrub-go/pubgrub/helpers"
	"github.com/mircearoata/pubgrub-go/pubgrub/semver"
	"github.com/puzpuzpuz/xsync/v3"
)

type ficsitAPISource struct {
	provider        Provider
	lockfile        *LockFile
	toInstall       map[string]semver.Constraint
	requiredTargets map[TargetName]bool
	modVersionInfo  *xsync.MapOf[string, []ModVersion]
	gameVersion     semver.Version
}

var clientTargets = map[TargetName]bool{
	TargetNameWindows: true,
}

var serverTargets = map[TargetName]bool{
	TargetNameWindowsServer: true,
	TargetNameLinuxServer:   true,
}

func (f *ficsitAPISource) GetPackageVersions(pkg string) ([]pubgrub.PackageVersion, error) {
	// If root package, return the base list of dependencies
	if pkg == rootPkg {
		return []pubgrub.PackageVersion{{Version: semver.Version{}, Dependencies: f.toInstall}}, nil
	}

	// Ignore game dependency
	if pkg == factoryGamePkg {
		return []pubgrub.PackageVersion{{Version: f.gameVersion}}, nil
	}

	response, err := f.provider.ModVersionsWithDependencies(context.TODO(), pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mod %s: %w", pkg, err)
	}

	f.modVersionInfo.Store(pkg, response)

	versions := make([]pubgrub.PackageVersion, 0)
	for _, modVersion := range response {
		v, err := semver.NewVersion(modVersion.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version %s: %w", modVersion.Version, err)
		}

		matches, err := f.matchesTargetRequirements(modVersion)
		if err != nil {
			return nil, err
		}
		if !matches {
			continue
		}

		dependencies := make(map[string]semver.Constraint)
		optionalDependencies := make(map[string]semver.Constraint)
		for _, dependency := range modVersion.Dependencies {
			c, err := semver.NewConstraint(dependency.Condition)
			if err != nil {
				return nil, fmt.Errorf("failed to parse constraint %s: %w", dependency.Condition, err)
			}

			if dependency.Optional {
				optionalDependencies[dependency.ModID] = c
			} else {
				dependencies[dependency.ModID] = c
			}
		}

		// If a version range string is empty, no version will satisfy it
		if modVersion.GameVersion != "" {
			factoryGameConstraint, err := semver.NewConstraint(modVersion.GameVersion)
			if err != nil {
				return nil, fmt.Errorf("failed to parse game version constraint %s: %w", modVersion.GameVersion, err)
			}

			dependencies[factoryGamePkg] = factoryGameConstraint
		}

		versions = append(versions, pubgrub.PackageVersion{
			Version:              v,
			Dependencies:         dependencies,
			OptionalDependencies: optionalDependencies,
		})
	}

	return versions, nil
}

func (f *ficsitAPISource) PickVersion(pkg string, versions []semver.Version) semver.Version {
	if f.lockfile != nil {
		if existing, ok := f.lockfile.Mods[pkg]; ok {
			v, err := semver.NewVersion(existing.Version)
			if err == nil {
				if slices.ContainsFunc(versions, func(version semver.Version) bool {
					return v.Compare(version) == 0
				}) {
					return v
				}
			}
		}
	}

	return helpers.StandardVersionPriority(versions)
}

func (f *ficsitAPISource) matchesTargetRequirements(modVersion ModVersion) (bool, error) {
	if len(f.requiredTargets) == 0 {
		return true, nil
	}

	requiredClientTargets := make(map[TargetName]bool)
	requiredServerTargets := make(map[TargetName]bool)

	for target := range f.requiredTargets {
		if clientTargets[target] {
			requiredClientTargets[target] = true
		} else if serverTargets[target] {
			requiredServerTargets[target] = true
		} else {
			return false, fmt.Errorf("unknown requested target %s", target)
		}
	}

	missingClientTargets := maps.Clone(requiredClientTargets)
	missingServerTargets := maps.Clone(requiredServerTargets)
	for _, target := range modVersion.Targets {
		delete(missingClientTargets, target.TargetName)
		delete(missingServerTargets, target.TargetName)
	}

	if modVersion.RequiredOnRemote {
		// All targets must be present
		return len(missingClientTargets) == 0 && len(missingServerTargets) == 0, nil
	}

	// Don't consider as having all targets when no targets of that type were requested
	hasAllClient := len(requiredClientTargets) > 0 && len(missingClientTargets) == 0
	hasAllServer := len(requiredServerTargets) > 0 && len(missingServerTargets) == 0
	return hasAllClient || hasAllServer, nil
}
