package repo

import "github.com/bilus/oya/pkg/types"

// Open opens a library containing all versions of a single Oya pack.
func Open(importPath types.ImportPath) (*GithubRepo, error) {
	if importPath.Host() != types.HostGithub {
		return nil, ErrNotGithub{ImportPath: importPath}
	}
	repoUri, basePath, packPath, err := parseImportPath(importPath)
	if err != nil {
		return nil, err
	}
	return &GithubRepo{
		repoUri:    repoUri,
		packPath:   packPath,
		basePath:   basePath,
		importPath: importPath,
	}, nil
}
