package job

import (
	"github.com/go-git/go-git/v5"
	config2 "github.com/go-git/go-git/v5/config"
	"os"
)

func GitCleanCloneOrPull(repositoyUrl, destinationDirectory string) error {
	//checkout MC-Remapper
	err := os.MkdirAll(destinationDirectory, 0777)
	if err != nil {
		return err
	}

	var repo *git.Repository
	repo, err = git.PlainOpen(destinationDirectory)
	if err != nil && err != git.ErrRepositoryNotExists {
		return err
	}

	if err == git.ErrRepositoryNotExists {
		repo, err = git.PlainClone(destinationDirectory, false, &git.CloneOptions{
			URL:        repositoyUrl,
			Progress:   os.Stdout,
			RemoteName: "origin",
		})
		if err != nil {
			return err
		}
	}

	if repo == nil {
		panic("repo is nil (╯°□°）╯︵ ┻━┻")
	}

	_, err = repo.Remote("origin")
	if err == git.ErrRemoteNotFound {
		_, err = repo.CreateRemote(&config2.RemoteConfig{
			Name: "origin",
			URLs: []string{repositoyUrl},
		})
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Clean(&git.CleanOptions{
		Dir: true,
	})
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}
