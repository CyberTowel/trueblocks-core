// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// TODO: Can this be concurrent?
func (opts *ChunksOptions) CheckSequential(fnArray, cacheArray, remoteArray []string, allowMissing bool, report *types.ReportCheck) error {
	if err := opts.checkSequential("disc", fnArray, allowMissing, report); err != nil {
		return err
	}

	report.MsgStrings = append(report.MsgStrings, "")
	if err := opts.checkSequential("cache", cacheArray, allowMissing, report); err != nil {
		return err
	}

	report.MsgStrings = append(report.MsgStrings, "")
	if err := opts.checkSequential("contract", remoteArray, allowMissing, report); err != nil {
		return err
	}

	return nil
}

func (opts *ChunksOptions) checkSequential(which string, array []string, allowMissing bool, report *types.ReportCheck) error {
	prev := cache.NotARange
	for _, item := range array {
		var fR cache.FileRange
		var err error
		if fR, err = cache.RangeFromFilename(item); err != nil {
			return err
		}

		w := utils.MakeFirstUpperCase(which)
		if prev != cache.NotARange {
			report.VisitedCnt++
			report.CheckedCnt++
			if prev != fR {
				if !prev.Preceeds(fR, !allowMissing) {
					// try to figure out why
					if prev.Intersects(fR) {
						report.MsgStrings = append(report.MsgStrings, fmt.Sprintf("%s: file ranges intersect %s:%s", w, prev, fR))
					} else if prev.Follows(fR, !allowMissing) {
						report.MsgStrings = append(report.MsgStrings, fmt.Sprintf("%s: previous range is after current %s:%s", w, prev, fR))
					} else if prev.Preceeds(fR, false) {
						report.MsgStrings = append(report.MsgStrings, fmt.Sprintf("%s: gap in sequence %d to %d skips %d", w, prev.Last, fR.First, (fR.First-prev.Last-1)))
					} else {
						report.MsgStrings = append(report.MsgStrings, fmt.Sprintf("%s: files not sequental %s:%s", w, prev, fR))
					}
				} else {
					report.PassedCnt++
				}
			} else {
				report.MsgStrings = append(report.MsgStrings, fmt.Sprintf("%s: duplicate at %s", w, fR))
			}
		}
		prev = fR
	}
	return nil
}
