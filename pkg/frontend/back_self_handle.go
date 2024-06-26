// Copyright 2024 Matrix Origin
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

package frontend

import (
	"context"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/tree"
)

func execInFrontendInBack(requestCtx context.Context,
	backSes *backSession,
	execCtx *ExecCtx) (err error) {
	//check transaction states
	switch st := execCtx.stmt.(type) {
	case *tree.BeginTransaction:
		err = backSes.GetTxnHandler().TxnBegin()
		if err != nil {
			return
		}
	case *tree.CommitTransaction:
		err = backSes.GetTxnHandler().TxnCommit()
		if err != nil {
			return
		}
	case *tree.RollbackTransaction:
		err = backSes.GetTxnHandler().TxnRollback()
		if err != nil {
			return
		}
	case *tree.Use:
		err = handleChangeDB(requestCtx, backSes, st.Name.Compare())
		if err != nil {
			return
		}
	case *tree.CreateDatabase:
		err = inputNameIsInvalid(execCtx.proc.Ctx, string(st.Name))
		if err != nil {
			return
		}
		if st.SubscriptionOption != nil && backSes.tenant != nil && !backSes.tenant.IsAdminRole() {
			err = moerr.NewInternalError(execCtx.proc.Ctx, "only admin can create subscription")
			return
		}
		st.Sql = execCtx.sqlOfStmt
	case *tree.DropDatabase:
		err = inputNameIsInvalid(execCtx.proc.Ctx, string(st.Name))
		if err != nil {
			return
		}
		// if the droped database is the same as the one in use, database must be reseted to empty.
		if string(st.Name) == backSes.GetDatabaseName() {
			backSes.SetDatabaseName("")
		}
	case *tree.Grant:
		switch st.Typ {
		case tree.GrantTypeRole:
			if err = handleGrantRole(requestCtx, backSes, st.GrantRole); err != nil {
				return
			}
		case tree.GrantTypePrivilege:
			if err = handleGrantPrivilege(requestCtx, backSes, st.GrantPrivilege); err != nil {
				return
			}
		}
	case *tree.Revoke:
		switch st.Typ {
		case tree.RevokeTypeRole:
			if err = handleRevokeRole(requestCtx, backSes, st.RevokeRole); err != nil {
				return
			}
		case tree.RevokeTypePrivilege:
			if err = handleRevokePrivilege(requestCtx, backSes, st.RevokePrivilege); err != nil {
				return
			}
		}
	case *tree.EmptyStmt:
		if err = handleEmptyStmt(requestCtx, backSes, st); err != nil {
			return
		}
	default:
		return moerr.NewInternalError(requestCtx, "backExec does not support %s", execCtx.sqlOfStmt)
	}
	return
}
