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

package validate

import "fmt"

type WrongConfigNameErr struct {
	configName         string
	expectedConfigName string
}

func (w *WrongConfigNameErr) Error() string {
	return fmt.Sprintf("config filename is wrong (%s); should be %s", w.configName, w.expectedConfigName)
}

type WrongArchInConfigErr struct {
	configPath string
	arch       string
}

func (w *WrongArchInConfigErr) Error() string {
	return fmt.Sprintf("wrong architecture in config file %s: %s", w.configPath, w.arch)
}

type WrongOutputProbeNameErr struct {
	outputProbeName         string
	expectedOutputProbeName string
}

func (w *WrongOutputProbeNameErr) Error() string {
	return fmt.Sprintf("output probe filename is wrong (%s); expected: %s.o", w.outputProbeName, w.expectedOutputProbeName)
}

type WrongOutputProbeArchErr struct {
	probe string
	arch  string
}

func (w *WrongOutputProbeArchErr) Error() string {
	return fmt.Sprintf("output probe filename has wrong architecture in its path (%s); expected %s", w.probe, w.arch)
}

type WrongOutputModuleNameErr struct {
	outputModuleName         string
	expectedOutputModuleName string
}

func (w *WrongOutputModuleNameErr) Error() string {
	return fmt.Sprintf("output module filename is wrong (%s); expected: %s.o", w.outputModuleName, w.expectedOutputModuleName)
}

type WrongOutputModuleArchErr struct {
	module string
	arch   string
}

func (w *WrongOutputModuleArchErr) Error() string {
	return fmt.Sprintf("output module filename has wrong architecture in its path (%s); expected %s", w.module, w.arch)
}

type KernelConfigDataNotBase64Err struct{}

func (k *KernelConfigDataNotBase64Err) Error() string {
	return fmt.Sprintf("kernelconfigdata must be a base64 encoded string")
}
