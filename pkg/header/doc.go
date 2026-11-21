// Copyright © 2026-present the longxi-go/golayout authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package header maintains repository-owned Go source notices independently of
// code layout. A plain-text template is authoritative; this package does not
// infer licensing, copyright ownership or calendar-year policy.
//
// Source preserves code, package documentation, build constraints and generated
// markers. Discovery, template loading and writes belong to callers, allowing
// generators and repository tooling to share the same in-memory operation.
package header
