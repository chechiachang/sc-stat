package git

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
)

func GitInit(url, directory string) {
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err == git.ErrRepositoryAlreadyExists {
		r, err = git.PlainOpen(directory)
		if err != nil {
			log.Panic(err)
		}
	} else if err != nil {
		log.Panic(err)
	}

	w, err := r.Worktree()
	if err != nil {
		log.Panic(err)
	}

	log.Info("Git initialized: ", r, w)

}
