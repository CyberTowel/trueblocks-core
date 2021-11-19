package transactionsPkg

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
/*
 * This file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
)

type TransactionsOptions struct {
	Transactions []string
	Articulate   bool
	Trace        bool
	Uniq         bool
	Reconcile    string
	Cache        bool
	Globals      globals.GlobalOptionsType
	BadFlag      error
}

func (opts *TransactionsOptions) TestLog() {
	logger.TestLog(len(opts.Transactions) > 0, "Transactions: ", opts.Transactions)
	logger.TestLog(opts.Articulate, "Articulate: ", opts.Articulate)
	logger.TestLog(opts.Trace, "Trace: ", opts.Trace)
	logger.TestLog(opts.Uniq, "Uniq: ", opts.Uniq)
	logger.TestLog(len(opts.Reconcile) > 0, "Reconcile: ", opts.Reconcile)
	logger.TestLog(opts.Cache, "Cache: ", opts.Cache)
	opts.Globals.TestLog()
}

func (opts *TransactionsOptions) ToCmdLine() string {
	options := ""
	if opts.Articulate {
		options += " --articulate"
	}
	if opts.Trace {
		options += " --trace"
	}
	if opts.Uniq {
		options += " --uniq"
	}
	if len(opts.Reconcile) > 0 {
		options += " --reconcile " + opts.Reconcile
	}
	if opts.Cache {
		options += " --cache"
	}
	options += " " + strings.Join(opts.Transactions, " ")
	options += fmt.Sprintf("%s", "") // silence go compiler for auto gen
	return options
}

func FromRequest(w http.ResponseWriter, r *http.Request) *TransactionsOptions {
	opts := &TransactionsOptions{}
	for key, value := range r.URL.Query() {
		switch key {
		case "transactions":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Transactions = append(opts.Transactions, s...)
			}
		case "articulate":
			opts.Articulate = true
		case "trace":
			opts.Trace = true
		case "uniq":
			opts.Uniq = true
		case "reconcile":
			opts.Reconcile = value[0]
		case "cache":
			opts.Cache = true
		default:
			if !globals.IsGlobalOption(key) {
				opts.BadFlag = validate.Usage("Invalid key ({0}) in {1} route.", key, "transactions")
				return opts
			}
		}
	}
	opts.Globals = *globals.FromRequest(w, r)

	return opts
}

var Options TransactionsOptions

func TransactionsFinishParse(args []string) *TransactionsOptions {
	// EXISTING_CODE
	Options.Transactions = args
	// EXISTING_CODE
	return &Options
}