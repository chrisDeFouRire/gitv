package lib

import (
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// DirtyFolder checks if the folder is clean or dirty
func DirtyFolder(repo *git.Repository) (bool, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return false, err
	}
	status, err := wt.Status()
	if err != nil {
		return false, err
	}
	if status.IsClean() {
		return false, nil
	} else {
		return true, nil
	}

}

// FindLatestSemverTag returns the latest semver tag found on current branch
// returns "","", nil if no tag can be found
func FindLatestSemverTag(repo *git.Repository) (string, plumbing.Hash, error) {
	tagList := make(map[plumbing.Hash]string)
	/* Get all tags indexed by hash */

	tags, err := repo.Tags()
	if err != nil {
		return "", plumbing.ZeroHash, err
	}

	for ref, err := tags.Next(); err == nil; ref, err = tags.Next() {
		tagName := ref.Name().Short()
		tagList[ref.Hash()] = tagName
	}
	tags.Close()

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return "", plumbing.ZeroHash, err
	}
	defer iter.Close()

	for ref, err := iter.Next(); err == nil; ref, err = iter.Next() {
		tag, found := tagList[ref.Hash]
		if found {
			_, err := semver.NewVersion(tag)
			if err == nil {
				return tag, ref.Hash, nil
			}
		}
	}
	return "", plumbing.ZeroHash, nil
}