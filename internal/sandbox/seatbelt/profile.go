package seatbelt

import (
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
			paths := make([]string, 0, len(p.Filesystem.ROPaths))

			for i, roPath := range p.Filesystem.ROPaths {
				abs, err2 := filepath.Abs(roPath)
				if err2 != nil {
					return "", nil, errors.Wrapf(err2, "absolute path for ROPath %s", roPath)
				}
				pName := fmt.Sprintf("RO_ROOT_%d", i)

				paths = append(paths, fmt.Sprintf(`(subpath (param %q))`, pName))
				args = append(args, param{k: pName, v: abs})
			}

			fileReads := fmt.Sprintf(
				"; Additional allow file-read's\n(allow file-read*\n%s\n)\n",
				strings.Join(paths, "\n"),
			)
			_, _ = policy.WriteString(fileReads)
		}

		if len(p.Filesystem.RWPaths) != 0 {
			paths := make([]string, 0, len(p.Filesystem.RWPaths))

			for i, rwPath := range p.Filesystem.RWPaths {
				abs, err2 := filepath.Abs(rwPath)
				if err2 != nil {
					return "", nil, errors.Wrapf(err2, "absolute path for RWPath %s", rwPath)
				}
				pName := fmt.Sprintf("RW_ROOT_%d", i)

				paths = append(paths, fmt.Sprintf(`(subpath (param %q))`, pName))
				args = append(args, param{k: pName, v: abs})
			}

			fileWrites := fmt.Sprintf(
				"; Additional allow file-write's\n(allow file-write*\n%s\n)\n",
				strings.Join(paths, "\n"),
			)
			_, _ = policy.WriteString(fileWrites)
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
