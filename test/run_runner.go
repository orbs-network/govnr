// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"
)

// uses the go test runner "go test" to run a test with an identical name
// in the _supervised_in_test directory and takes expectations regarding output
func executeGoTestRunner(t *testing.T, expectedLogs []string, unexpectedLogs []string) {
	out, _ := exec.Command(
		path.Join(runtime.GOROOT(), "bin", "go"),
		"test",
		"../_supervised_in_test/",
		"-v",
		"-run",
		"^("+t.Name()+")$").CombinedOutput()

	goTestOutput := string(out)
	ribbon := "------------------ EXTERNAL TEST OUTPUT (_supervised_in_test." + t.Name() + ")------------------"
	debugMsgOutput := fmt.Sprintln(ribbon, "\n", goTestOutput, "\n", ribbon)

	for _, logLine := range expectedLogs {
		require.Truef(t, strings.Contains(goTestOutput, logLine), "log should contain: '%s'\n\n%s", logLine, debugMsgOutput)
	}
	for _, logLine := range unexpectedLogs {
		require.Falsef(t, strings.Contains(goTestOutput, logLine), "log should not contain: '%s'\n\n%s", logLine, debugMsgOutput)
	}
}
