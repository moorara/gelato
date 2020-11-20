package git

import (
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	idPattern       = `[A-Za-z][0-9A-Za-z-]+[0-9A-Za-z]`
	domainPattern   = fmt.Sprintf(`%s\.[A-Za-z]{2,63}`, idPattern)
	repoPathPattern = fmt.Sprintf(`(%s/){1,20}(%s)`, idPattern, idPattern)
	httpsPattern    = fmt.Sprintf(`^https://(%s)/(%s)(.git)?$`, domainPattern, repoPathPattern)
	sshPattern      = fmt.Sprintf(`^git@(%s):(%s)(.git)?$`, domainPattern, repoPathPattern)
	httpsRE         = regexp.MustCompile(httpsPattern)
	sshRE           = regexp.MustCompile(sshPattern)
)

// Git provides Git functionalities.
type Git struct {
	repo *git.Repository
}

// New creates a new instance of Git.
func New(path string) (*Git, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		return nil, err
	}

	return &Git{
		repo: repo,
	}, nil
}

func parseRemoteURL(url string) (string, string, error) {
	// Parse the origin remote URL into a domain part a path part
	if m := httpsRE.FindStringSubmatch(url); len(m) == 6 { // HTTPS Git Remote URL
		//  Example:
		//    https://github.com/moorara/changelog.git
		//    m = []string{"https://github.com/moorara/changelog.git", "github.com", "moorara/changelog", "moorara/", "changelog", ".git"}
		return m[1], m[2], nil
	} else if m := sshRE.FindStringSubmatch(url); len(m) == 6 { // SSH Git Remote URL
		//  Example:
		//    git@github.com:moorara/changelog.git
		//    m = []string{"git@github.com:moorara/changelog.git", "github.com", "moorara/changelog, "moorara/", "changelog", ".git"}
		return m[1], m[2], nil
	}

	return "", "", fmt.Errorf("invalid git remote url: %s", url)
}

// Remote returns the domain part and path part of a Git remote repository URL.
// It assumes the remote repository is named origin.
func (g *Git) Remote(name string) (string, string, error) {
	remote, err := g.repo.Remote(name)
	if err != nil {
		return "", "", err
	}

	// TODO: Should we handle all URLs and not just the first one?
	var remoteURL string
	if config := remote.Config(); len(config.URLs) > 0 {
		remoteURL = config.URLs[0]
	}

	return parseRemoteURL(remoteURL)
}

// IsClean determines whether or not the working directory is clean.
func (g *Git) IsClean() (bool, error) {
	worktree, err := g.repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := worktree.Status()
	if err != nil {
		return false, err
	}

	return status.IsClean(), nil
}

// HEAD returns the hash and name (branch) of the HEAD reference.
func (g *Git) HEAD() (string, string, error) {
	head, err := g.repo.Head()
	if err != nil {
		return "", "", err
	}

	hash := head.Hash().String()
	branch := head.Name().String()

	return hash, branch, nil
}

// Pull is same as git pull. It brings the changes from a remote repository into the current branch.
func (g *Git) Pull(ctx context.Context) error {
	worktree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	opts := &git.PullOptions{}

	err = worktree.PullContext(ctx, opts)
	if err != nil {
		return err
	}

	return nil
}

// Tag resolves a tag by its name.
func (g *Git) Tag(name string) (Tag, error) {
	ref, err := g.repo.Tag(name)
	if err != nil {
		return Tag{}, err
	}

	var tag Tag

	tagObj, err := g.repo.TagObject(ref.Hash())
	switch err {
	// Annotated tag
	case nil:
		commitObj, err := g.repo.CommitObject(tagObj.Target)
		if err != nil {
			return Tag{}, err
		}
		tag = toAnnotatedTag(tagObj, commitObj)

	// Lightweight tag
	case plumbing.ErrObjectNotFound:
		commitObj, err := g.repo.CommitObject(ref.Hash())
		if err != nil {
			return Tag{}, err
		}
		tag = toLightweightTag(ref, commitObj)

	default:
		return Tag{}, err
	}

	return tag, nil
}

// Tags returns the list of all tags.
func (g *Git) Tags() ([]Tag, error) {
	refs, err := g.repo.Tags()
	if err != nil {
		return nil, err
	}

	tags := []Tag{}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		tagObj, err := g.repo.TagObject(ref.Hash())
		switch err {
		// Annotated tag
		case nil:
			commitObj, err := g.repo.CommitObject(tagObj.Target)
			if err != nil {
				return err
			}
			tag := toAnnotatedTag(tagObj, commitObj)
			tags = append(tags, tag)

		// Lightweight tag
		case plumbing.ErrObjectNotFound:
			commitObj, err := g.repo.CommitObject(ref.Hash())
			if err != nil {
				return err
			}
			tag := toLightweightTag(ref, commitObj)
			tags = append(tags, tag)

		default:
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort tags
	sort.Slice(tags, func(i, j int) bool {
		// The order of the tags should be from the most recent to the least recent
		return tags[i].After(tags[j])
	})

	return tags, nil
}
