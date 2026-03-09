package markdown_test

import (
"testing"

"github.com/chifamba/canonical-corpus/pkg/markdown"
)

func TestBuildFilename(t *testing.T) {
tests := []struct {
lang, translationID, want string
}{
{"en", "kjv", "en-kjv.md"},
{"en", "web", "en-web.md"},
{"en", "asv", "en-asv.md"},
{"el", "lxx", "el-lxx.md"},
{"he", "hmt", "he-hmt.md"},
{"sn", "bdb", "sn-bdb.md"},
{"en", "", "en.md"},
{"", "", "en.md"},
}
for _, tt := range tests {
got := markdown.BuildFilename(tt.lang, tt.translationID)
if got != tt.want {
t.Errorf("BuildFilename(%q, %q) = %q, want %q", tt.lang, tt.translationID, got, tt.want)
}
}
}
