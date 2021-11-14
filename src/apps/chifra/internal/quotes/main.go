package quotesPkg

/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/root"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) error {
	options := ""
	if Options.Freshen {
		options += " --freshen"
	}
	if len(Options.Period) > 0 {
		options += " --period " + Options.Period
	}
	if len(Options.Pair) > 0 {
		options += " --pair " + Options.Pair
	}
	if len(Options.Feed) > 0 {
		options += " --feed " + Options.Feed
	}
	arguments := ""
	for _, arg := range args {
		arguments += " " + arg
	}

	return root.PassItOn("getQuotes", &Options.Globals, options, arguments)
}
