package slinguist

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TSuite struct{}

var _ = Suite(&TSuite{})
