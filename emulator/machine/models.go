package machine

import (
	"fmt"
	"strings"
)

// -----------------------------------------------------------------------------
// Machine model factory
// -----------------------------------------------------------------------------

var models = map[string]*Model{} // Registered models by ID

// Model contains a machine model description. Machine packages register
// their models at init time (blank import in the emulator package); the
// registry is read only afterwards.
type Model struct {
	Name  string                 // Machine model
	Ids   []string               // Model Ids
	Build func(Settings) Machine // Build builds the machine model
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

// Create returns a machine of the settings model
func Create(settings Settings) (Machine, error) {
	model := FindModel(settings.ModelID)
	if model == nil {
		return nil, fmt.Errorf("Machine : unknown machine model : %q", settings.ModelID)
	}
	machine := model.Build(settings)
	machine.Config().Name = model.Name
	return machine, nil
}
