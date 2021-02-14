package questionnaire

import "github.com/Andrianov/psychoHelpBOT/internal/models"

var (
	nameStep = models.Step{
		Name:     "ФИО",
		Question: "Как к вам можно обращаться?",
	}

	contactStep = models.Step{
		Name:     "Контакты",
		Question: "Никнейм телеграма, чтобы куратор мог с вами связаться",
	}

	ageStep = models.Step{
		Name:     "Совершеннолетний?",
		Question: "Являетесь ли вы совершеннолетним?",
		Options:  []string{"да", "нет"},
	}

	symptomsStep = models.Step{
		Name:     "Проявления",
		Question: "Перечислите проявления из следующего списка, если вы замечаете их у себя. Если их нет, поставьте прочерк\n\nНарушение сна, нарушение концентрации внимания, притупленность эмоций, приступы ярости",
	}

	stateStep = models.Step{
		Name:     "Состояние",
		Question: "Опишите свое состояние, а также то, что с вами произошло",
	}

	placeStep = models.Step{
		Name:     "Страна/Город",
		Question: "В какой стране и в каком городе вы проживаете?",
	}
)

var FlowSteps = []models.Step{
	nameStep,
	ageStep,
	symptomsStep,
	stateStep,
	placeStep,
}

var AnonymousFlowSteps = []models.Step{
	nameStep,
	contactStep,
	ageStep,
	symptomsStep,
	stateStep,
	placeStep,
}
