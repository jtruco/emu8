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
	Name  string         // Machine model
	Ids   []string       // Model Ids
	Build func() Machine // Build builds the machine model
}

// Register register a machine model
func Register(model *Model) {
	id := strings.ToLower(model.Name)
	models[id] = model

	// register other model IDs
	for _, ids := range model.Ids {
		id = strings.ToLower(ids)
		models[id] = model
	}
}

// RegisterModels register a machine model list
func RegisterModels(models []Model) {
	for i := 0; i < len(models); i++ {
		Register(&models[i])
	}
}

// FindModel finds a model by Id
func FindModel(modelId string) *Model {
	modelId = strings.ToLower(modelId)
	model, ok := models[modelId]
	if ok {
		return model
	}
	return nil
}

// Create returns a machine from a model name
func Create(modelId string) (Machine, error) {
	model := FindModel(modelId)
	if model == nil {
		return nil, errors.New("Machine : unknown machine model")
	}
	machine := model.Build()
	machine.Config().Name = model.Name
	return machine, nil
}
