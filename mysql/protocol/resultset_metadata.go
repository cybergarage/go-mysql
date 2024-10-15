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

// MySQL :: MySQL 8.0 C API Developer Guide :: 5.4.67 mysql_result_metadata()
// https://dev.mysql.com/doc/c-api/8.0/en/mysql-result-metadata.html

type ResultsetMetadata = uint8

const (
	// ResultsetMetadataNone represents the MYSQL_RESULT_METADATA_NONE.
	ResultsetMetadataNone ResultsetMetadata = 0
	// ResultsetMetadataFull represents the MYSQL_RESULT_METADATA_FULL.
	ResultsetMetadataFull ResultsetMetadata = 1
)
