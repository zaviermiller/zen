package vm

import (
	"fmt"
	"strconv"
	"strings"
)

// current Zen version -- ONLY EDIT HERE
var VERSION Version = Version{Major: 1, Minor: 1, Patch: 0}

type Version struct {
	Major int
	Minor int
	Patch int
}

func ParseVersion(version string) Version {
	versions := strings.Split(version, ".")
	v := Version{Patch: -1}

	for i, ver := range versions {
		// dont rlly think there will be an error here :)
		ver, _ := strconv.Atoi(ver)

		switch i {
		case 0:
			v.Major = ver
		case 1:
			v.Minor = ver
		case 2:
			v.Patch = ver
		}
	}

	return v
}

// GreaterThan returns true if the calling version is
// greater
func (v Version) GreaterThan(v1 Version) bool {
	if v.Major > v1.Major {
		return true
	} else if v.Minor > v1.Minor {
		return true
	} else if v.Patch > v1.Patch {
		return true
	}

	return false
}

func (v Version) String() string {
	if v.Patch != -1 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}
