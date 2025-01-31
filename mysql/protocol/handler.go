// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

// CommandHandler represents a MySQL command handler.
type CommandHandler interface {
	// HandleQuery handles a query command.
	HandleQuery(Conn, *Query) (Response, error)
	// PrepareStatement prepares a statement.
	PrepareStatement(Conn, *StmtPrepare) (*StmtPrepareResponse, error)
	// ExecuteStatement executes a statement.
	ExecuteStatement(Conn, *StmtExecute) (Response, error)
	// CloseStatement closes a statement.
	CloseStatement(Conn, *StmtClose) (Response, error)
}
