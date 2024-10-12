package resolver

import (
	"context"
	"fmt"
	"strings"

	"github.com/mircearoata/pubgrub-go/pubgrub"
	"github.com/mircearoata/pubgrub-go/pubgrub/semver"
)

type DependencyResolverError struct {
	pubgrub.SolvingError
	provider    Provider
	gameVersion int
}

func (e DependencyResolverError) Error() string {
	rootPkg := e.Cause().Terms()[0].Dependency()

	writer := pubgrub.NewStandardErrorWriter(rootPkg).WithIncompatibilityStringer(
		MakeDependencyResolverErrorStringer(e.provider, e.gameVersion),
	)
	e.WriteTo(writer)

	return writer.String()
}

type DependencyResolverErrorStringer struct {
	pubgrub.StandardIncompatibilityStringer
	provider     Provider
	packageNames map[string]string
	gameVersion  int
}

func MakeDependencyResolverErrorStringer(provider Provider, gameVersion int) *DependencyResolverErrorStringer {
	s := &DependencyResolverErrorStringer{
		provider:     provider,
		gameVersion:  gameVersion,
		packageNames: map[string]string{},
	}
	s.StandardIncompatibilityStringer = pubgrub.NewStandardIncompatibilityStringer().WithTermStringer(s)
	return s
}

func (w *DependencyResolverErrorStringer) getPackageName(pkg string) string {
	if pkg == factoryGamePkg {
		return "Satisfactory"
	}

	if name, ok := w.packageNames[pkg]; ok {
		return name
	}

	result, err := w.provider.GetModName(context.TODO(), pkg)
	if err != nil {
		return pkg
	}

	w.packageNames[pkg] = result.Name

	return result.Name
}

func (w *DependencyResolverErrorStringer) Term(t pubgrub.Term, includeVersion bool) string {
	name := w.getPackageName(t.Dependency())
	fullName := fmt.Sprintf("%s (%s)", name, t.Dependency())
	if name == t.Dependency() {
		fullName = t.Dependency()
	}

	if includeVersion {
		if t.Constraint().IsAny() {
			return fmt.Sprintf("every version of %s", fullName)
		}

		switch t.Dependency() {
		case factoryGamePkg:
			// Remove ".0.0" from the versions mentioned, since only the major is ever used
			return fmt.Sprintf("%s \"%s\"", fullName, strings.ReplaceAll(t.Constraint().String(), ".0.0", ""))
		default:
			res, err := w.provider.ModVersionsWithDependencies(context.TODO(), t.Dependency())
			if err != nil {
				return fmt.Sprintf("%s \"%s\"", fullName, t.Constraint())
			}

			var matched []semver.Version
			for _, v := range res {
				ver, err := semver.NewVersion(v.Version)
				if err != nil {
					// Assume it is contained in the constraint
					matched = append(matched, semver.Version{})
					continue
				}

				if t.Constraint().Contains(ver) {
					matched = append(matched, ver)
				}
			}

			if len(matched) == 1 {
				return fmt.Sprintf("%s \"%s\"", fullName, matched[0])
			}

			return fmt.Sprintf("%s \"%s\"", fullName, t.Constraint())
		}
	}

	return fullName
}

func (w *DependencyResolverErrorStringer) IncompatibilityString(incompatibility *pubgrub.Incompatibility, rootPkg string) string {
	terms := incompatibility.Terms()

	if len(terms) == 1 && terms[0].Dependency() == factoryGamePkg {
		return fmt.Sprintf("Satisfactory CL%d is installed", w.gameVersion)
	}

	return w.StandardIncompatibilityStringer.IncompatibilityString(incompatibility, rootPkg)
}
