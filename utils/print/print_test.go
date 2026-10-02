package print_test

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"patrol_build_flutter/utils/print"
)

func captureOutput(f func(), t *testing.T) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe: %v", err)
	}
	os.Stdout = stdout
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}
	return buf.String()
}

func expectedColor(color string) string {
	if runtime.GOOS == "windows" {
		return ""
	}
	return color
}

func TestError_PrintsRed(t *testing.T) {
	msg := "error message"
	out := captureOutput(func() { print.Error(msg) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Red), msg, expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestSuccess_PrintsGreen(t *testing.T) {
	msg := "success"
	out := captureOutput(func() { print.Success(msg) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Green), msg, expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestWarning_PrintsYellow(t *testing.T) {
	msg := "warn"
	out := captureOutput(func() { print.Warning(msg) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Yellow), msg, expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestAction_PrintsBlue(t *testing.T) {
	msg := "action"
	out := captureOutput(func() { print.Action(msg) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Blue), msg, expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestStepCompleted_PrintsPurple(t *testing.T) {
	msg := "done"
	out := captureOutput(func() { print.StepCompleted(msg) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Purple), msg, expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestStepInitiated_PrintsCyan(t *testing.T) {
	msg := "init"
	out := captureOutput(func() { print.StepInitiated(msg) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Cyan), msg, expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestInfo_PrintsCyan(t *testing.T) {
	msg := "info message"
	out := captureOutput(func() { print.Info(msg) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Cyan), msg, expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestErrorf_PrintsRedFormatted(t *testing.T) {
	out := captureOutput(func() { print.Errorf("error %d: %s", 404, "not found") }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Red), "error 404: not found", expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestSuccessf_PrintsGreenFormatted(t *testing.T) {
	out := captureOutput(func() { print.Successf("done in %d seconds", 5) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Green), "done in 5 seconds", expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestWarningf_PrintsYellowFormatted(t *testing.T) {
	out := captureOutput(func() { print.Warningf("warning: %s", "low memory") }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Yellow), "warning: low memory", expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestInfof_PrintsCyanFormatted(t *testing.T) {
	out := captureOutput(func() { print.Infof("info: %s v%s", "app", "1.0") }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Cyan), "info: app v1.0", expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestActionf_PrintsBlueFormatted(t *testing.T) {
	out := captureOutput(func() { print.Actionf("executing %s", "build") }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Blue), "executing build", expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestStepCompletedf_PrintsPurpleFormatted(t *testing.T) {
	out := captureOutput(func() { print.StepCompletedf("step %d completed", 1) }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Purple), "step 1 completed", expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestStepInitiatedf_PrintsCyanFormatted(t *testing.T) {
	out := captureOutput(func() { print.StepInitiatedf("starting step %s", "test") }, t)
	want := fmt.Sprintf("%s%s%s\n", expectedColor(print.Cyan), "starting step test", expectedColor(print.Reset))
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestVanilla_PrintsPlain(t *testing.T) {
	msg := "plain"
	out := captureOutput(func() { print.Vanilla(msg) }, t)
	if strings.TrimSpace(out) != msg {
		t.Errorf("expected %q, got %q", msg, out)
	}
}

func TestVanillaf_PrintsPlainFormatted(t *testing.T) {
	out := captureOutput(func() { print.Vanillaf("plain %d", 42) }, t)
	if strings.TrimSpace(out) != "plain 42" {
		t.Errorf("expected %q, got %q", "plain 42", out)
	}
}

func TestWindows_PrintColors(t *testing.T) {
	// Emulate Windows environment
	print.SetColorsForOS("windows")

	if print.Reset != "" {
		t.Errorf("Reset should be empty on windows")
	}
	if print.Red != "" {
		t.Errorf("Red should be empty on windows")
	}
	if print.Green != "" {
		t.Errorf("Green should be empty on windows")
	}
	if print.Yellow != "" {
		t.Errorf("Yellow should be empty on windows")
	}
	if print.Blue != "" {
		t.Errorf("Blue should be empty on windows")
	}
	if print.Purple != "" {
		t.Errorf("Purple should be empty on windows")
	}
	if print.Cyan != "" {
		t.Errorf("Cyan should be empty on windows")
	}
}
