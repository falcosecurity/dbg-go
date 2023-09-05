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
