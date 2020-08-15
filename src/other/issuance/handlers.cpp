/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2018, 2019 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
#include "options.h"

//--------------------------------------------------------------
class CAccumulator {
public:
    CReconciliation rec;
    period_t period;
    time_q prevDate;
    string_q fmt;
    bool discrete;
    CAccumulator(void) {
        period = BY_NOTHING;
        discrete = false;
    }
};

const string_q uncleBlocks = "data/uncle_blocks.csv";
const string_q resultsFile = "data/results.csv";
//--------------------------------------------------------------
bool COptions::model_issuance(void) {
    if (!fileExists(uncleBlocks))
        return usage("Cannot open list of uncle blocks. Run with --uncles option first.");

    blknum_t prev = 0;
    CUintArray blocks;

    string_q contents = asciiFileToString(uncleBlocks);
    explode(blocks, contents, '\n');

    for (auto block : blocks) {
        for (blknum_t bn = prev; bn < block; bn++) {
            if (true) {
                CReconciliation rec;
                rec.blockNum = bn;
                rec.timestamp = getBlockTimestamp(bn);
                if (bn == 0) {
                    cout << STR_HEADER_EXPORT << endl;
                    for (auto item : prefundWeiMap)
                        rec.minerRewardInflow2 = (rec.minerRewardInflow2 + item.second);
                } else {
                    rec.minerRewardInflow2 = getBlockReward(bn, false);
                }
                if (bn != 0 && bn == prev) {
                    // This block has an uncle
                    rec.minerAddInflow = getBlockReward(bn, true) - rec.minerRewardInflow2;
                    uint64_t count = getUncleCount(bn);
                    for (size_t i = 0; i < count; i++) {
                        CBlock uncle;
                        getUncle(uncle, bn, i);
                        rec.uncleRewardInflow = (rec.uncleRewardInflow + getUncleReward(bn, uncle.blockNumber));
                    }
                }
                cout << rec.Format(STR_DISPLAY_EXPORT) << endl;
            }
        }
        prev = block;
    }

    return true;
}

CStringArray header;
////--------------------------------------------------------------
//bool auditLine(const char* str, void* data) {
//    if (header.empty()) {
//        string_q headers = "blockNum,timestamp,month,day,minerRewardInflow2,minerAddInflow,minerIssuance,uncleRewardInflow,issuance";
//        explode(header, headers, ',');
//        return true;
//    }
//
//    string_q line = str;
//    CReconciliation rec;
//    rec.parseCSV(header, line);
//    if (rec.blockNum == 0)
//        return true;
//    // END LAST TIME: 459690
//    if (rec.blockNum < 1145042) //(firstTransactionBlock - 10))
//        return true;
//
//    CTransaction trans;
//    trans.loadTransAsBlockReward(rec.blockNum, 99999, block.miner);
//    expContext().accountedFor = block.miner;
//    if (rec.reconcileIssuance(rec.blockNum, &trans)) {
//        cout << "Block " << rec.blockNum << " balances" << endl;
//    } else {
//        cerr << "Block " << rec.blockNum << " does not balances. Hit enter to coninue..." << endl;
//        cerr << rec << endl;
//        getchar();
//    }
//    return true;
//}

//--------------------------------------------------------------
int sortByFirstAppearance(const void *v1, const void *v2) {
    CAddressRecord_base *vv1 = (CAddressRecord_base*)v1;
    CAddressRecord_base *vv2 = (CAddressRecord_base*)v2;
    if (vv1->offset > vv2->offset)
        return 1;
    else if (vv1->offset < vv2->offset)
        return -1;
    return 0;
}

//--------------------------------------------------------------
int sortReverseByTxid(const void *v1, const void *v2) {
    CAppearance_base *vv1 = (CAppearance_base*)v1;
    CAppearance_base *vv2 = (CAppearance_base*)v2;
    if (vv2->txid > vv1->txid)
        return 1;
    else if (vv2->txid < vv1->txid)
        return -1;
    return 0;
}

//--------------------------------------------------------------
bool reconcileIssuance(const CAppearance& app) {
    if (app.bn == 0)
        return true;
    if (app.addr == "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead")
        return true;

    bigint_t minerReward, minerTxFee, uncleReward1, uncleReward2, netTx;
    size_t d1=0, d2=0;
    uint64_t uncleCount = getUncleCount(app.bn);
    if (app.tx != 99998) {
        minerReward = getBlockReward(app.bn);
        minerTxFee = getTransFees(app.bn);
    }

    for (size_t i = 0 ; i < uncleCount ; i++) {
        CBlock uncle;
        getUncle(uncle, app.bn, i);
        if (uncle.miner == app.addr) {
            if (uncleReward1 == 0) {
                uncleReward1 = getUncleReward(app.bn, uncle.blockNumber);
                d1 = app.bn - uncle.blockNumber;
            } else {
                uncleReward2 = getUncleReward(app.bn, uncle.blockNumber);
                d2 = app.bn - uncle.blockNumber;
            }
        }
    }

    bigint_t begBal = getBalanceAt(app.addr, app.bn - 1);
    bigint_t endBal = getBalanceAt(app.addr, app.bn);
    bigint_t expected = begBal + minerReward + minerTxFee + uncleReward1 + uncleReward2;

    if (expected == endBal) {
        return true;
    }

    CBlock block;
    getBlock(block, app.bn);
    for (auto trans : block.transactions) {
        if (trans.from == app.addr) {
            netTx -= trans.isError ? 0 : trans.value;
            netTx -= str_2_BigInt(trans.getValueByName("gasCost"));
        }
        if (trans.to == app.addr) {
            netTx += trans.isError ? 0 : trans.value;
        }
        expected = begBal + minerReward + minerTxFee + uncleReward1 + uncleReward2 + netTx;
        cerr << expected << " : " << endBal << " : " << (expected - endBal) << string_q(90, ' ') << "\r"; cerr.flush();
        if (endBal == expected) return true;
    }

    if (endBal != expected) {
        cout << endl;
        cout << app << endl;
        cout << "hash:         " << block.blockNumber << endl;
        cout << "begBal:       " << padLeft(bni_2_Str(begBal), 20) << endl;
        cout << "endBal:       " << padLeft(bni_2_Str(endBal), 20) << endl;
        cout << "minerReward:  " << padLeft(bni_2_Str(minerReward), 20) << endl;
        cout << "minerTxFee:   " << padLeft(bni_2_Str(minerTxFee), 20) << endl;
        cout << "d1:           " << padLeft(int_2_Str(d1), 20) << endl;
        cout << "uncleReward1: " << padLeft(bni_2_Str(uncleReward1), 20) << endl;
        cout << "d2:           " << padLeft(int_2_Str(d2), 20) << endl;
        cout << "uncleReward2: " << padLeft(bni_2_Str(uncleReward2), 20) << endl;
        cout << "netTx:        " << padLeft(bni_2_Str(netTx), 20) << endl;
        cout << "expected:     " << padLeft(bni_2_Str(expected), 20) << endl;
        cout << "diff:         " << padLeft(bni_2_Str((expected - endBal)), 20) << endl;
    }

    return (endBal == expected);
}

//--------------------------------------------------------------
bool visitIndexChunk(CIndexArchive& chunk, void* data) {
    static bool skip = true;
    string_q fn = substitute(chunk.getFilename(), indexFolder_finalized, "");
    if (contains(fn, "006963699-006965754"))
        skip = false;
    if (skip) return true;

    size_t nMiners = 0;
    size_t nUncles = 0;
    CAppearanceArray checkList;
    for (size_t i = 0 ; i < chunk.nAddrs ; i++) {
        CAddressRecord_base *rec = &chunk.addresses[i];
        CAppearance_base *start = &chunk.appearances[rec->offset];
        qsort(start, rec->cnt, sizeof(CAppearance_base), sortReverseByTxid);
        if (start[0].blk != 0 && (start[0].txid == 99997 || start[0].txid == 99998 || start[0].txid == 99999)) {
            CAppearance_base *appPtr = &chunk.appearances[rec->offset];
            CAppearance app;
            app.bn = appPtr->blk;
            app.tx = appPtr->txid;
            app.addr = bytes_2_Addr(rec->bytes);
            nMiners += (appPtr->txid == 99997 || appPtr->txid == 99999);
            nUncles += (appPtr->txid == 99998);
            checkList.push_back(app);
        }
    }

    cout << fn << " nMiners: " << nMiners << " nUncles: " << nUncles << string_q(120, ' ') << endl;
    for (auto app : checkList) {
        if (reconcileIssuance(app)) {
            cerr << fn << ": Block reward for  " << app.addr << " at block " << app.bn << " balances" << "    \r"; cerr.flush();
        } else {
            cout << fn << ": Block reward for  " << app.addr << " at block " << app.bn << bRed << " does not balance" << cOff << endl;
        }
    }
    return true;
}

namespace qblocks {
extern bool forEveryIndexChunk(INDEXCHUNKFUNC func, void* data);
}
//--------------------------------------------------------------
bool COptions::audit_issuance(void) {
    //forEveryLineInAsciiFile(resultsFile, auditLine, NULL);
    qblocks::forEveryIndexChunk(visitIndexChunk, NULL);
    return true;
}

//--------------------------------------------------------------
bool visitLine(const char* str, void* data) {
    CAccumulator *acc = (CAccumulator*)data;

    if (header.empty()) {
        string_q headers = "blockNum,timestamp,month,day,minerRewardInflow2,minerAddInflow,minerIssuance,uncleRewardInflow,issuance";
        explode(header, headers, ',');
        return true;
    }

    string_q line = str;
    CReconciliation rec;
    rec.parseCSV(header, line);

    time_q curDate = str_2_Date(rec.Format(per_2_Str(acc->period)));
    bool same = isSamePeriod(acc->period, acc->prevDate, curDate);
    if (!same) {
        CReconciliationOutput out(acc->rec);
        cout << out.Format(acc->fmt) << endl;
        acc->prevDate = curDate;
        if (acc->discrete) {
            CReconciliation reset;
            acc->rec = reset;
        }
        acc->rec.blockNum = rec.blockNum;
        acc->rec.timestamp = rec.timestamp;
    } else {
        if (!(rec.timestamp % 23))
            cerr << "Processing " << ts_2_Date(rec.timestamp).Format(FMT_JSON) << "\r"; cerr.flush();
    }
    acc->rec = (acc->rec + rec);
//    cerr << "A: " << acc->rec.Format(acc->fmt) << endl;
//    cerr << "B: " << rec.Format(acc->fmt) << endl;
    return true;
}

//--------------------------------------------------------------
bool COptions::show_uncle_blocks(void) {
    blknum_t latest = getLatestBlock_client();
    for (size_t bn = 0 ; bn < latest ; bn++) {
        if (getUncleCount(bn))
            cout << bn << endl;
    }
    return true;
}

//--------------------------------------------------------------
bool COptions::summary_by_period(void) {

    CAccumulator accumulator;
    accumulator.period = by_period;
    accumulator.fmt = substitute(STR_DISPLAY_EXPORT, "[{MONTH}],[{DAY}]", per_2_Str(by_period));
    accumulator.rec.timestamp = blockZeroTs;
    accumulator.prevDate = ts_2_Date(blockZeroTs);
    accumulator.discrete = discrete;
    forEveryLineInAsciiFile(resultsFile, visitLine, &accumulator);
    CReconciliationOutput out(accumulator.rec);
    cout << out.Format(accumulator.fmt) << endl;

    return true;
}

//--------------------------------------------------------------------------------
string_q STR_DISPLAY_EXPORT =
"[{BLOCKNUM}],"\
"[{TIMESTAMP}],"\
"[{MONTH}],"\
"[{DAY}],"\
"[{MINERREWARDINFLOW2}],"\
"[{MINERADDINFLOW}],"\
"[{MINERISSUANCE}],"\
"[{UNCLEREWARDINFLOW}],"\
"[{ISSUANCE}]";


string_q STR_HEADER_EXPORT =
"blockNum,timestamp,month,day,minerReward,minerAdd,minerIssuance,uncleReward,issuance";
