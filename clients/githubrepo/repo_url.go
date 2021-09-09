// Copyright 2020 Security Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package githubrepo

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ossf/scorecard/v2/clients"
	sce "github.com/ossf/scorecard/v2/errors"
)

const (
	githubOrgRepo  = ".github"
	scorecardOwner = "ossf"
	scorecardRepo  = "scorecard"
)

var ()

type repoURL struct {
	host, owner, repo string
	metadata          []string
}

// Parses input string into repoURL struct.
// Accepts "owner/repo" or "github.com/owner/repo".
func (r *repoURL) parse(input string) error {
	var t string

	const two = 2
	const three = 3

	c := strings.Split(input, "/")

	switch l := len(c); {
	// This will takes care for repo/owner format.
	// By default it will use github.com
	case l == two:
		t = "github.com/" + c[0] + "/" + c[1]
	case l >= three:
		t = input
	}

	// Allow skipping scheme for ease-of-use, default to https.
	if !strings.Contains(t, "://") {
		t = "https://" + t
	}

	u, e := url.Parse(t)
	if e != nil {
		//nolint:wrapcheck
		return sce.Create(sce.ErrScorecardInternal, fmt.Sprintf("url.Parse: %v", e))
	}

	const splitLen = 2
	split := strings.SplitN(strings.Trim(u.Path, "/"), "/", splitLen)
	if len(split) != splitLen {
		//nolint:wrapcheck
		return sce.Create(sce.ErrorInvalidURL, fmt.Sprintf("%v. Exepted full repository url", input))
	}

	r.host, r.owner, r.repo = u.Host, split[0], split[1]
	return nil
}

// URL implements Repo.URL.
func (r *repoURL) URL() string {
	return fmt.Sprintf("%s/%s/%s", r.host, r.owner, r.repo)
}

// String implements Repo.String.
func (r *repoURL) String() string {
	return fmt.Sprintf("%s-%s-%s", r.host, r.owner, r.repo)
}

// Org implements Repo.Org.
func (r *repoURL) Org() clients.Repo {
	return &repoURL{
		host:  r.host,
		owner: r.owner,
		repo:  githubOrgRepo,
	}
}

// IsValid implements Repo.IsValid.
func (r *repoURL) IsValid() error {
	switch r.host {
	case "github.com":
	default:
		//nolint:wrapcheck
		return sce.Create(sce.ErrorUnsupportedHost, r.host)
	}

	if strings.TrimSpace(r.owner) == "" || strings.TrimSpace(r.repo) == "" {
		//nolint:wrapcheck
		return sce.Create(sce.ErrorInvalidURL,
			fmt.Sprintf("%v. Expected the full reposiroty url", r.URL()))
	}
	return nil
}

func (r *repoURL) AppendMetadata(metadata ...string) {
	for _, arg := range metadata {
		r.metadata = append(r.metadata, arg)
	}
}

// Metadata implements Repo.Metadata.
func (r *repoURL) Metadata() []string {
	return r.metadata
}

// IsScorecardRepo implements Repo.IsScorecardRepo.
func (r *repoURL) IsScorecardRepo() bool {
	return r.owner == scorecardOwner && r.repo == scorecardRepo
}

// MakeGithubRepo takes input of form "owner/repo" or "github.com/owner/repo"
// and returns an implementation of clients.Repo interface.
func MakeGithubRepo(input string) (clients.Repo, error) {
	var repo repoURL
	if err := repo.parse(input); err != nil {
		return nil, fmt.Errorf("error during parse: %w", err)
	}
	if err := repo.IsValid(); err != nil {
		return nil, fmt.Errorf("error in IsValid: %w", err)
	}
	return &repo, nil
}
