// Copyright 2022 Tatris Project Authors. Licensed under Apache-2.0.

package metadata

import (
	"github.com/tatris-io/tatris/internal/common/log/logger"
	"github.com/tatris-io/tatris/internal/meta/metadata/storage"
	"github.com/tatris-io/tatris/internal/meta/metadata/storage/boltdb"
	"go.uber.org/zap"
)

const AliasPath = "/_alias/"
const IndexPath = "/_index/"
const IndexTemplatePath = "/_index_template/"

var MStore storage.MetaStore

func init() {
	var err error
	MStore, err = boltdb.Open()
	if err != nil {
		logger.Panic("init metastore failed", zap.Error(err))
	}

	if err := LoadIndexes(); err != nil {
		logger.Panic("load indexes failed", zap.Error(err))
	}

	if err := LoadAliases(); err != nil {
		logger.Panic("load aliases failed", zap.Error(err))
	}

	if err := LoadIndexTemplates(); err != nil {
		logger.Panic("load index templates failed", zap.Error(err))
	}
}