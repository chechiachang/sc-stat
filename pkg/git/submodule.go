package git

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
)

func GitInit(url, directory, submodule string) {
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err == git.ErrRepositoryAlreadyExists {
		r, err = git.PlainOpen(directory)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	w, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	sub, err := w.Submodule(submodule)
	if err != nil {
		log.Fatal(err)
	}

	sr, err := sub.Repository()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("submodule initialized: ", submodule, sr)

}
