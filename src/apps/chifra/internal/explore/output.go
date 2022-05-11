// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package explorePkg

// EXISTING_CODE
import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

// RunExplore handles the explore command for the command line. Returns error only as per cobra.
func RunExplore(cmd *cobra.Command, args []string) (err error) {
	opts := ExploreFinishParse(args)
	// EXISTING_CODE
	// EXISTING_CODE
	err, _ = opts.ExploreInternal()
	return
}

// ServeExplore handles the explore command for the API. Returns error and a bool if handled
func ServeExplore(w http.ResponseWriter, r *http.Request) (err error, handled bool) {
	opts := ExploreFinishParseApi(w, r)
	// EXISTING_CODE
	// EXISTING_CODE
	return opts.ExploreInternal()
}

// ExploreInternal handles the internal workings of the explore command.  Returns error and a bool if handled
func (opts *ExploreOptions) ExploreInternal() (err error, handled bool) {
	err = opts.ValidateExplore()
	if err != nil {
		return err, true
	}

	// EXISTING_CODE
	if opts.Globals.ApiMode {
		return validate.Usage("Cannot use explore route in API mode."), true
	}

	for _, url := range urls {
		ret := url.getUrl(opts)
		fmt.Printf("Opening %s\n", ret)
		if !opts.Globals.TestMode {
			utils.OpenBrowser(ret)
		}
	}
	// EXISTING_CODE

	return
}

// EXISTING_CODE
func (u *ExploreUrl) getUrl(opts *ExploreOptions) string {

	var chain = opts.Globals.Chain

	// TODO: Multi-chain customize remote explorer strings per chain

	if opts.Google {
		var query = "https://www.google.com/search?q=[{TERM}]"
		query = strings.Replace(query, "[{TERM}]", u.term, -1)
		var exclusions = []string{
			"etherscan", "etherchain", "bloxy", "bitquery", "ethplorer", "tokenview", "anyblocks", "explorer",
		}
		for _, ex := range exclusions {
			query += ("+-" + ex)
		}
		return query
	}

	if u.termType == ExploreFourByte {
		var query = "https://www.4byte.directory/signatures/?bytes4_signature=[{TERM}]"
		query = strings.Replace(query, "[{TERM}]", u.term, -1)
		return query
	}

	if u.termType == ExploreEnsName {
		var query = "https://app.ens.domains/name/[{TERM}]/details"
		query = strings.Replace(query, "[{TERM}]", u.term, -1)
		return query
	}

	url := config.GetRemoteExplorer(chain)
	query := ""
	switch u.termType {
	case ExploreNone:
		// do nothing
	case ExploreTx:
		query = "tx/" + u.term
	case ExploreBlock:
		query = "block/" + u.term
	case ExploreAddress:
		fallthrough
	default:
		query = "address/" + u.term
	}

	if opts.Local {
		url = config.GetLocalExplorer(chain)
		query = strings.Replace(query, "tx/", "explorer/transactions/", -1)
		query = strings.Replace(query, "block/", "explorer/blocks/", -1)
	}

	return url + query
}

// EXISTING_CODE
