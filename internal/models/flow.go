package models

type Flow struct {
	Steps []*Step `json:"steps"`
}

func (f *Flow) IsFinished() bool {
	return f.NextStep() == nil
}

func (f *Flow) NextStep() *Step {
	for _, step := range f.Steps {
		if !step.IsAnswered() {
			return step
		}
	}

	return nil
}