package ciloop_test

import (
	"testing"

	"github.com/fall-out-bug/sdp/internal/ciloop"
)

func TestClassifyGoTest(t *testing.T) {
	got := ciloop.Classify("go-test")
	if got != ciloop.ClassAutoFixable {
		t.Errorf("expected AutoFixable for go-test, got %q", got)
	}
}

func TestClassifyGoBuild(t *testing.T) {
	got := ciloop.Classify("go-build")
	if got != ciloop.ClassAutoFixable {
		t.Errorf("expected AutoFixable for go-build, got %q", got)
	}
}

func TestClassifyK8sValidate(t *testing.T) {
	got := ciloop.Classify("k8s-validate")
	if got != ciloop.ClassAutoFixable {
		t.Errorf("expected AutoFixable for k8s-validate, got %q", got)
	}
}

func TestClassifySecrets(t *testing.T) {
	got := ciloop.Classify("secrets-scan")
	if got != ciloop.ClassEscalate {
		t.Errorf("expected Escalate for secrets-scan, got %q", got)
	}
}

func TestClassifyFlaky(t *testing.T) {
	got := ciloop.Classify("flaky-detector")
	if got != ciloop.ClassEscalate {
		t.Errorf("expected Escalate for flaky-detector, got %q", got)
	}
}

func TestClassifyUnknownEscalates(t *testing.T) {
	got := ciloop.Classify("some-unknown-check")
	if got != ciloop.ClassEscalate {
		t.Errorf("expected Escalate for unknown check, got %q", got)
	}
}

func TestClassifyGoTestCaseInsensitive(t *testing.T) {
	cases := []string{"Go-Test", "GO-BUILD", "K8S-VALIDATE"}
	for _, c := range cases {
		got := ciloop.Classify(c)
		if got != ciloop.ClassAutoFixable {
			t.Errorf("expected AutoFixable for %q (case-insensitive), got %q", c, got)
		}
	}
}
