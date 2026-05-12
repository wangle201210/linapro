// This file defines the CMS controller constructor and dependencies.

package cms

import (
	"lina-plugin-cms/backend/api/cms"
	cmssvc "lina-plugin-cms/backend/internal/service/cms"
)

// ControllerV1 is the CMS plugin controller.
type ControllerV1 struct {
	cmsSvc cmssvc.Service // cmsSvc owns CMS business operations.
}

// NewV1 creates and returns the CMS API controller interface.
func NewV1() cms.ICmsV1 {
	return NewControllerV1()
}

// NewControllerV1 creates and returns the concrete CMS controller.
func NewControllerV1() *ControllerV1 {
	return &ControllerV1{cmsSvc: cmssvc.New()}
}
