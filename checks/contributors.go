// Copyright 2020 OpenSSF Scorecard Authors
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

package checks

import (
	"github.com/ossf/scorecard/v4/checker"
	"github.com/ossf/scorecard/v4/checks/evaluation"
	"github.com/ossf/scorecard/v4/checks/raw"
	sce "github.com/ossf/scorecard/v4/errors"
)

// CheckContributors is the registered name for Contributors.
const CheckContributors = "Contributors"

//nolint:gochecknoinits
func init() {
	if err := registerCheck(CheckContributors, Contributors, nil); err != nil {
		// this should never happen
		panic(err)
	}
}

// Contributors run Contributors check.
func Contributors(c *checker.CheckRequest) checker.CheckResult {
	rawData, err := raw.Contributors(c.RepoClient)
	if err != nil {
		e := sce.WithMessage(sce.ErrScorecardInternal, err.Error())
		return checker.CreateRuntimeErrorResult(CheckContributors, e)
	}

	// Return raw results.
	if c.RawResults != nil {
		c.RawResults.ContributorsResults = rawData
	}

	// Return the score evaluation.
	return evaluation.Contributors(CheckContributors, c.Dlogger, &rawData)
}
