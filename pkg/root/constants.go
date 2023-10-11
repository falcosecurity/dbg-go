// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package root

const (
	configPathFmt = "%s/driverkit/config/%s/%s/%s" // Eg: repo-root/driverkit/config/5.0.1+driver/x86_64/centos_5.14.0-325.el9.x86_64_1.yaml
	outputPathFmt = "%s/driverkit/output/%s/%s/%s" // Eg: repo-root/driverkit/output/5.0.1+driver/x86_64/falco_centos_5.14.0-325.el9.x86_64_1.{ko,o}
)
