package github

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chechiachang/sc-stat/pkg/utils"
	"github.com/google/go-github/v63/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
)

var client *github.Client
var ctx = context.Background()

func getRef(sourceOwner, sourceRepo, commitBranch, baseBranch *string) (ref *github.Reference, err error) {
	if ref, _, err = client.Git.GetRef(ctx, *sourceOwner, *sourceRepo, "refs/heads/"+*commitBranch); err == nil {
		return ref, nil
	}

	// We consider that an error means the branch has not been found and needs to
	// be created.
	if *commitBranch == *baseBranch {
		return nil, errors.New("the commit branch does not exist but `-base-branch` is the same as `-commit-branch`")
	}

	if *baseBranch == "" {
		return nil, errors.New("the `-base-branch` should not be set to an empty string when the branch specified by `-commit-branch` does not exists")
	}

	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, *sourceOwner, *sourceRepo, "refs/heads/"+*baseBranch); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + *commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, *sourceOwner, *sourceRepo, newRef)
	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(ref *github.Reference, sourceFiles, sourceOwner, sourceRepo *string) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	// Load each file into the tree.
	for _, fileArg := range strings.Split(*sourceFiles, ",") {
		file, content, err := getFileContent(fileArg, sourceFiles)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &github.TreeEntry{
			Path:    github.String(strings.SplitN(file, "/", 2)[1]), // remove root from tree path for submodule
			Type:    github.String("blob"),
			Content: github.String(string(content)),
			Mode:    github.String("100644")})
	}

	tree, _, err = client.Git.CreateTree(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string, sourceFiles *string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = os.ReadFile(localFile)
	return targetName, b, err
}

func pushCommit(ref *github.Reference, tree *github.Tree, commitMessage, sourceOwner, sourceRepo, authorName, authorEmail, privateKey *string) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &github.Timestamp{Time: date}, Name: authorName, Email: authorEmail}
	commit := &github.Commit{Author: author, Message: commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	opts := github.CreateCommitOptions{}
	if *privateKey != "" {
		armoredBlock, e := os.ReadFile(*privateKey)
		if e != nil {
			return e
		}
		keyring, e := openpgp.ReadArmoredKeyRing(bytes.NewReader(armoredBlock))
		if e != nil {
			return e
		}
		if len(keyring) != 1 {
			return errors.New("expected exactly one key in the keyring")
		}
		key := keyring[0]
		opts.Signer = github.MessageSignerFunc(func(w io.Writer, r io.Reader) error {
			return openpgp.ArmoredDetachSign(w, key, r, nil)
		})
	}
	newCommit, _, err := client.Git.CreateCommit(ctx, *sourceOwner, *sourceRepo, commit, &opts)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, *sourceOwner, *sourceRepo, ref, false)
	return err
}

func CommitPush() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	// prepare commit message
	date := time.Now()
	commitMessage := fmt.Sprintf("chore: upload data on %s", date.Format("2006-1-2"))
	sourceOwner := "chechiachang"
	sourceRepo := "sc-stat-data"
	commitBranch := "main"
	baseBranch := "main"
	authorName := "sc-stat-automation"
	authorEmail := "chechiachang999@gmail.com"
	privateKey := ""

	// prepare commit files
	// TODO list files in data
	files := utils.Glob("data", func(s string) bool {
		return filepath.Ext(s) == ".csv" &&
			(strings.Contains(s, date.Format("2006-1-2")) ||
				strings.Contains(s, date.AddDate(0, 0, -1).Format("2006-1-2"))) // today and yesterday
	})
	sourceFiles := strings.Join(files, ",")

	if sourceOwner == "" || sourceRepo == "" || commitBranch == "" || sourceFiles == "" || authorName == "" || authorEmail == "" {
		log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	}
	client = github.NewClient(nil).WithAuthToken(token)

	ref, err := getRef(&sourceOwner, &sourceRepo, &commitBranch, &baseBranch)
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}
	if ref == nil {
		log.Fatalf("No error where returned but the reference is nil")
	}

	tree, err := getTree(ref, &sourceFiles, &sourceOwner, &sourceRepo)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	if err := pushCommit(ref, tree, &commitMessage, &sourceOwner, &sourceRepo, &authorName, &authorEmail, &privateKey); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	log.Infof("Commit created successfully: %s\n", sourceFiles)

}
