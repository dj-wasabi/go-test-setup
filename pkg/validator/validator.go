package validator

type Validator struct {
	Errors map[string]string
}

// type logging struct {
// 	logging *slog.Logger
// }

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// func (v *Validator) PrintErrors(app any, in map[string]string) {
// 	logger := logging.Initialize()
// 	for k, v := range in {
// 		logger.Error(k, "value", v)
// 	}
// }
