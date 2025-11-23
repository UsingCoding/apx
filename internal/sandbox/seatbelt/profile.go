package seatbelt

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/sandbox"
)

// params for sandbox-exec parameters
type params []param

func (p params) flat() (res []string) {
	for _, v := range p {
		res = append(res, fmt.Sprintf("-D%s=%s", v.k, v.v))
	}
	return
}

type param struct {
	k, v string
}

var (
	//go:embed seatbelt_base_policy.sbpl
	seatbeltBasePolicy []byte
	//go:embed seatbelt_network_policy.sbpl
	seatbeltNetworkPolicy []byte
)

func makeProfile(p sandbox.Policy) (_ string, args params, err error) {
	var policy strings.Builder

	wd, err := os.Getwd()
	if err != nil {
		return "", nil, errors.Wrap(err, "get current working directory")
	}

	_, _ = policy.WriteString(string(seatbeltBasePolicy) + "\n")
	// params from base
	args = params{
		{k: "WRITABLE_ROOT_0", v: wd},
		{k: "WRITABLE_ROOT_0_RO_0", v: path.Join(wd, ".git")},
		// allow access to /tmp, /var via /private to avoid mismatch on macOS
		{k: "WRITABLE_ROOT_1", v: "/private/tmp"},
		{k: "WRITABLE_ROOT_2", v: "/private/var"},
		{k: "WRITABLE_ROOT_3", v: "/tmp"},
		{k: "WRITABLE_ROOT_4", v: "/var"},
	}

	// user's home directory parameter for SBPL rules that exclude it from broad read access
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil, errors.Wrap(err, "get home dir")
	}
	args = append(args, param{k: "USER_HOME_DIR", v: home})

	//nolint:gocritic
	{
		if p.Filesystem.FullDiskReadAccess {
			const diskReadAccess = "; allow read-only file operations\n(allow file-read*)"

			_, _ = policy.WriteString(diskReadAccess)
		}

		if len(p.Filesystem.ROPaths) != 0 {
			policy.WriteString("; Additional allow file-read's\n(allow file-read*\n")

			for i, roPath := range p.Filesystem.ROPaths {
				abs, err2 := filepath.Abs(roPath)
				if err2 != nil {
					return "", nil, errors.Wrapf(err2, "absolute path for ROPath %s", roPath)
				}
				pName := fmt.Sprintf("RO_ROOT_%d", i)

				policy.WriteString(fmt.Sprintf("	(subpath (param %q))\n", pName))
				args = append(args, param{k: pName, v: abs})
			}

			policy.WriteString(")\n")
		}

		if len(p.Filesystem.RWPaths) != 0 {
			policy.WriteString("; Additional allow file-write's\n(allow file-write*\n")

			for i, rwPath := range p.Filesystem.RWPaths {
				abs, err2 := filepath.Abs(rwPath)
				if err2 != nil {
					return "", nil, errors.Wrapf(err2, "absolute path for RWPath %s", rwPath)
				}
				pName := fmt.Sprintf("RW_ROOT_%d", i)

				policy.WriteString(fmt.Sprintf("	(subpath (param %q))\n", pName))
				args = append(args, param{k: pName, v: abs})
			}

			policy.WriteString(")\n")
		}

		if len(p.Filesystem.DenyPaths) != 0 {
			_, _ = policy.WriteString("(deny file-read* file-read-metadata\n")

			for i, denyPath := range p.Filesystem.DenyPaths {
				abs, err2 := filepath.Abs(denyPath)
				if err2 != nil {
					return "", nil, errors.Wrapf(err2, "absolute path for DenyPath %s", denyPath)
				}
				pName := fmt.Sprintf("DENY_PATH_%d", i)

				_, _ = policy.WriteString(fmt.Sprintf("	(subpath (param %q))\n", pName))
				args = append(args, param{k: pName, v: abs})
			}

			_, _ = policy.WriteString(")\n")
		}

		if p.Filesystem.NoCache {
			userCacheDir, err2 := os.UserCacheDir()
			if err2 != nil {
				return "", nil, errors.Wrapf(err2, "get user cache dir")
			}

			// allow cache read and write
			const cacheROAccess = "(allow file-read*\n  (subpath (param \"DARWIN_USER_CACHE_DIR\"))\n)"
			const cacheRWAccess = "(allow file-write*\n  (subpath (param \"DARWIN_USER_CACHE_DIR\"))\n)"

			_, _ = policy.WriteString(cacheROAccess)
			_, _ = policy.WriteString(cacheRWAccess)
			args = append(args, param{k: "DARWIN_USER_CACHE_DIR", v: userCacheDir})
		}
	}

	if !p.Network.Deny {
		_, _ = policy.WriteString(string(seatbeltNetworkPolicy) + "\n")
	}

	return policy.String(), args, nil
}
