package utils

import "testing"

func TestExec(t *testing.T) {
	out, err := Exec("/usr/bin/git", false)
	if err != nil {
		t.Error(err)
	}
	t.Logf("输出: %v", out)
}

func TestShellout(t *testing.T) {
	err, out, e := Shellout("git clone --bare https://educoder:xinhu1ji2qu3@bdgit.educoder.net/forge01/listnodedemo.git")
	t.Log(err, out, e)
}
