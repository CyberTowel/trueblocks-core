package chunksPkg

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
 * The file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

import (
	"net/http"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/root"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type ChunksOptionsType struct {
	Blocks  []string
	List    bool
	Check   bool
	Extract string
	Stats   bool
	Save    bool
	Globals root.GlobalOptionsType
}

var Options ChunksOptionsType

func (opts *ChunksOptionsType) TestLog() {
	logger.TestLog(len(opts.Blocks) > 0, "Blocks: ", opts.Blocks)
	logger.TestLog(opts.List, "List: ", opts.List)
	logger.TestLog(opts.Check, "Check: ", opts.Check)
	logger.TestLog(len(opts.Extract) > 0, "Extract: ", opts.Extract)
	logger.TestLog(opts.Stats, "Stats: ", opts.Stats)
	logger.TestLog(opts.Save, "Save: ", opts.Save)
	opts.Globals.TestLog()
}

func FromRequest(r *http.Request) *ChunksOptionsType {
	opts := &ChunksOptionsType{}
	for key, value := range r.URL.Query() {
		switch key {
		case "blocks":
			opts.Blocks = append(opts.Blocks, value...)
		case "list":
			opts.List = true
		case "check":
			opts.Check = true
		case "extract":
			opts.Extract = value[0]
		case "stats":
			opts.Stats = true
		case "save":
			opts.Save = true
		}
	}
	opts.Globals = *root.FromRequest(r)

	return opts
}
