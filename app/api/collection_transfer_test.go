package api

import (
	"testing"

	"forgeturl-server/api/space"
)

func TestTransferCollectionCopy(t *testing.T) {
	source := []*space.Collections{{Title: "source"}, {Title: "keep"}}
	target := []*space.Collections{{Title: "target"}}

	newSource, newTarget, err := transferCollection(source, target, 0, collectionTransferCopy)
	if err != nil {
		t.Fatalf("transferCollection() error = %v", err)
	}
	if len(newSource) != 2 || len(newTarget) != 2 || newTarget[1].Title != "source" {
		t.Fatalf("unexpected copy result: source=%v target=%v", newSource, newTarget)
	}
}

func TestTransferCollectionMove(t *testing.T) {
	source := []*space.Collections{{Title: "move"}, {Title: "keep"}}
	target := []*space.Collections{{Title: "target"}}

	newSource, newTarget, err := transferCollection(source, target, 0, collectionTransferMove)
	if err != nil {
		t.Fatalf("transferCollection() error = %v", err)
	}
	if len(newSource) != 1 || newSource[0].Title != "keep" || len(newTarget) != 2 || newTarget[1].Title != "move" {
		t.Fatalf("unexpected move result: source=%v target=%v", newSource, newTarget)
	}
}

func TestTransferCollectionRejectsInvalidInput(t *testing.T) {
	_, _, err := transferCollection(nil, nil, 0, collectionTransferCopy)
	if err == nil {
		t.Fatal("expected missing collection error")
	}

	_, _, err = transferCollection([]*space.Collections{{Title: "source"}}, nil, 0, "delete")
	if err == nil {
		t.Fatal("expected invalid operation error")
	}
}
