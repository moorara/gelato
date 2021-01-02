package git

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	branch := strings.TrimPrefix(head.Name().String(), "refs/heads/")

	return hash, branch, nil
}

// Tag resolves a tag by its name.
func (g *Git) Tag(name string) (Tag, error) {
	ref, err := g.repo.Tag(name)
	if err != nil {
		return Tag{}, err
	}

	var tag Tag

	t, err := g.repo.TagObject(ref.Hash())
	switch err {
	// Annotated tag
	case nil:
		c, err := g.repo.CommitObject(t.Target)
		if err != nil {
			return Tag{}, err
		}
		tag = toAnnotatedTag(t, c)

	// Lightweight tag
	case plumbing.ErrObjectNotFound:
		c, err := g.repo.CommitObject(ref.Hash())
		if err != nil {
			return Tag{}, err
		}
		tag = toLightweightTag(ref, c)

	default:
		return Tag{}, err
	}

	return tag, nil
}

// Tags returns the list of all tags.
func (g *Git) Tags() (Tags, error) {
	refs, err := g.repo.Tags()
	if err != nil {
		return nil, err
	}

	tags := []Tag{}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		t, err := g.repo.TagObject(ref.Hash())
		switch err {
		// Annotated tag
		case nil:
			c, err := g.repo.CommitObject(t.Target)
			if err != nil {
				return err
			}
			tags = append(tags, toAnnotatedTag(t, c))

		// Lightweight tag
		case plumbing.ErrObjectNotFound:
			c, err := g.repo.CommitObject(ref.Hash())
			if err != nil {
				return err
			}
			tags = append(tags, toLightweightTag(ref, c))

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

// CreateTag creates a new annotated tag with a message.
// If successful, it returns the hash of the newly created tag.
func (g *Git) CreateTag(commit, name, message string) (string, error) {
	opts := &git.CreateTagOptions{Message: message}
	hash := plumbing.NewHash(commit)
	ref, err := g.repo.CreateTag(name, hash, opts)
	if err != nil {
		return "", err
	}

	return ref.Hash().String(), nil
}

func (g *Git) parentCommits(commitsMap map[plumbing.Hash]*object.Commit, h plumbing.Hash) error {
	if _, ok := commitsMap[h]; ok {
		return nil
	}

	c, err := g.repo.CommitObject(h)
	if err != nil {
		return err
	}

	commitsMap[c.Hash] = c

	for _, h := range c.ParentHashes {
		if err := g.parentCommits(commitsMap, h); err != nil {
			return err
		}
	}

	return nil
}

// CommitsIn returns all commits reachable from a revision.
func (g *Git) CommitsIn(rev string) (Commits, error) {
	h, err := g.repo.ResolveRevision(plumbing.Revision(rev))
	if err != nil {
		return nil, err
	}

	commitsMap := make(map[plumbing.Hash]*object.Commit)
	err = g.parentCommits(commitsMap, *h)
	if err != nil {
		return nil, err
	}

	commits := make([]Commit, 0)
	for _, c := range commitsMap {
		commits = append(commits, toCommit(c))
	}

	// Sort commits
	sort.Slice(commits, func(i, j int) bool {
		// The order of the commits should be from the most recent to the least recent
		return commits[i].Committer.After(commits[j].Committer)
	})

	return commits, nil
}

// CreateCommit stages a list of files in the working tree and then creates a new commit with a give message.
// If successful, it returns the hash of the newly created commit.
func (g *Git) CreateCommit(message string, files ...string) (string, error) {
	worktree, err := g.repo.Worktree()
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if _, err := worktree.Add(file); err != nil {
			return "", err
		}
	}

	hash, err := worktree.Commit(message, &git.CommitOptions{})
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

// Pull is same as git pull. It brings the changes from a remote repository into the current branch.
func (g *Git) Pull(ctx context.Context) error {
	worktree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	opts := &git.PullOptions{}

	if err = worktree.PullContext(ctx, opts); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return nil
		}
		return err
	}

	return nil
}

// Push performs a push to a remote repository.
func (g *Git) Push(ctx context.Context, remoteName string) error {
	return g.repo.PushContext(ctx, &git.PushOptions{
		RemoteName: remoteName,
	})
}

// PushTag pushes a tag a remote repository.
func (g *Git) PushTag(ctx context.Context, remoteName, tagName string) error {
	return g.repo.PushContext(ctx, &git.PushOptions{
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/tags/" + tagName + ":refs/tags/" + tagName),
		},
	})
}
