package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestScoreBot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ScoreBot Suite")
}
