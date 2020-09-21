package semver_test

import (
	"fmt"

	"github.com/moorara/gelato/pkg/semver"
)

func ExampleParse() {
	v, _ := semver.Parse("v0.1.0-rc.1+sha.abcdeff.20200820")
	fmt.Printf("Semantic Version: %s\n", v)
}

func ExampleParse_invalid() {
	str := "v1.0"
	if _, ok := semver.Parse(str); !ok {
		fmt.Printf("Invalid semantic version: %s\n", str)
	}
}

func ExampleSemVer_AddPrerelease() {
	v := semver.SemVer{Major: 0, Minor: 1, Patch: 0}
	v.AddPrerelease("beta")
	fmt.Printf("Semantic Version: %s\n", v)
}

func ExampleSemVer_AddMetadata() {
	v := semver.SemVer{Major: 0, Minor: 1, Patch: 0}
	v.AddMetadata("20200920")
	fmt.Printf("Semantic Version: %s\n", v)
}

func ExampleSemVer_Next() {
	v := semver.SemVer{Major: 0, Minor: 1, Patch: 0}
	fmt.Printf("Next: %s\n", v.Next())
}

func ExampleSemVer_ReleasePatch() {
	v := semver.SemVer{Major: 0, Minor: 1, Patch: 0}
	fmt.Printf("Patch Release: %s\n", v.ReleasePatch())
}

func ExampleSemVer_ReleaseMinor() {
	v := semver.SemVer{Major: 0, Minor: 1, Patch: 0}
	fmt.Printf("Minor Release: %s\n", v.ReleaseMinor())
}

func ExampleSemVer_ReleaseMajor() {
	v := semver.SemVer{Major: 0, Minor: 1, Patch: 0}
	fmt.Printf("Major Release: %s\n", v.ReleaseMajor())
}
