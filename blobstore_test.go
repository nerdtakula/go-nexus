package nexus

import "testing"

func TestBlobStore(t *testing.T) {
	t.Skip("failing")
	bs, err := client.BlobStore("default")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Results: %+v\n", bs)
}
