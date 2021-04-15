package machine

import (
	"errors"
	"strings"
)

// -----------------------------------------------------------------------------
// Machine model factory
// -----------------------------------------------------------------------------

var models = map[string]*Model{} // Registered models by ID

// Description machine model description
type Model struct {
	Id       string         // Machine model Id
	OtherIds []string       // Other models Id
	Build    func() Machine // Build builds the machine model
}

// Register register a machine model
func Register(model *Model) {
	models[model.Id] = model

	// register other model IDs
	for _, otherId := range model.OtherIds {
		models[otherId] = model
	}
}

// RegisterModels register a machine model list
func RegisterModels(models []Model) {
	for _, model := range models {
		Register(&model)
	}
}

// FindModel finds a model by ID
func FindModel(modelID string) *Model {
	modelID = strings.ToLower(modelID)
	model, ok := models[modelID]
	if ok {
		return model
	}
	return nil
}

// Create returns a machine from a model name
func Create(modelID string) (Machine, error) {
	model := FindModel(modelID)
	if model == nil {
		return nil, errors.New("Machine : unknown machine model")
	}
	return model.Build(), nil
}
