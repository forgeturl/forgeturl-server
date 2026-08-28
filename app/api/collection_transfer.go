package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/space"

	"github.com/bytedance/sonic"
)

const (
	collectionTransferCopy = "copy"
	collectionTransferMove = "move"
)

func decodeCollections(content string) ([]*space.Collections, error) {
	collections := make([]*space.Collections, 0)
	if content == "" {
		return collections, nil
	}
	if err := sonic.UnmarshalString(content, &collections); err != nil {
		return nil, err
	}
	return collections, nil
}

func transferCollection(
	sourceCollections []*space.Collections,
	targetCollections []*space.Collections,
	sourceIndex int,
	operation string,
) ([]*space.Collections, []*space.Collections, error) {
	if sourceIndex < 0 || sourceIndex >= len(sourceCollections) {
		return nil, nil, common.ErrBadRequest("source collection does not exist")
	}
	if operation != collectionTransferCopy && operation != collectionTransferMove {
		return nil, nil, common.ErrBadRequest("operation must be copy or move")
	}

	targetCollections = append(targetCollections, sourceCollections[sourceIndex])
	if operation == collectionTransferMove {
		sourceCollections = append(sourceCollections[:sourceIndex], sourceCollections[sourceIndex+1:]...)
	}
	return sourceCollections, targetCollections, nil
}
