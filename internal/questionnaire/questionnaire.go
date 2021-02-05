package questionnaire

import (
	"errors"
	"fmt"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Andrianov/psychoHelpBOT/internal/config"
	"github.com/Andrianov/psychoHelpBOT/internal/models"
	"github.com/Andrianov/psychoHelpBOT/internal/storage"
)

type Questionnaire struct {
	cfg     *config.Config
	bot     *tgApi.BotAPI
	storage storage.Storage
}

const (
	_ = `Внимание, важная информация: наш чат бот временно не работает. 

Наш аккаунт, принимающий заявки, временно заблокировали. Мы уже подали апелляцию и делаем все возможное, чтобы восстановить работу. 

В профиле будет апдейт, когда мы сможем добиться результата. 

Причиной блокировки стали множественные жалобы на спам. При этом единственная деятельность, осуществляемая на аккаунте - ответы на заявки...
Мы не хотим делать выводы, не имея достоверных данных, так что просто скажем, что это прискорбно. 

Наши психологи сейчас работают с теми заявками, что удалось обработать.

Хотим оставить для вас информацию о тех центрах и проектах, которые безвозмездно оказывают психологическую помощь сейчас: 

Центр "Насилию.нет" открыли чат бот для сбора заявок. 

"Открытое пространство" оказывает психологическую помощь для правозащитников и активистов. 

Также центр «Сестры» оказывает психологическую поддержку тем, «кто остро переживает происходящее сейчас».`

	startMessage = `Данный чат-бот служит сбору заявок от людей, пострадавших во время мирных митингов РФ, а также от тех, кто остро проживает актуальные события. 

Команда психологов, психотерапевтов и психиатров, объединенная движением «Психология за Права Человека» в проекте помощи людям, которые переживают насилие в Беларуси, работает с людьми из России. Мы убеждены в том, что все, несмотря на политические взгляды и выбор участвовать или нет в уличных акциях - имеют право на то, чтобы получить бесплатную качественную психологическую помощь в ситуации социальной и политической неопределенности, острого нарушения безопасности, насилия со стороны силовых структур. 

Сводки по проекту каждый вечер публикуются в [инстаграме](https://instagram.com/_eto_normalno).

Кураторками/ами проекта по оказанию психологической помощи являются: 
[Анна Край](https://instagram.com/_to_the_edge_)
[Игнат Пименов](https://instagram.com/_eto_normalno)
[Ольга Размахова](https://instagram.com/za_900_let)`

	cancelMessage = "Заполнение анкеты отменено"

	finalMessage = `Ваши ответы отправлены. Заявка обрабатывается. Вам будут высланы контакты психолога/психоерапевта для связи.

Помощь уже рядом.`

	sorryMessage = "Не удалось обработать ответ, разбираемся..."
)

func New(cfg *config.Config, bot *tgApi.BotAPI, storage storage.Storage) *Questionnaire {
	return &Questionnaire{cfg, bot, storage}
}

func (q *Questionnaire) Intro(update tgApi.Update) error {
	if update.Message == nil {
		return errors.New("message is nil")
	}

	msg := tgApi.NewMessage(update.Message.Chat.ID, startMessage)
	msg.ParseMode = "markdown"
	_, err := q.bot.Send(msg)
	return err
}

func (q *Questionnaire) Cancel(update tgApi.Update) error {
	if update.Message == nil {
		return errors.New("message is nil")
	}

	chatID := update.Message.Chat.ID

	err := q.storage.Delete(chatID)
	if err != nil {
		return err
	}

	msg := tgApi.NewMessage(chatID, cancelMessage)
	_, err = q.bot.Send(msg)
	return err
}

func (q *Questionnaire) Start(update tgApi.Update) error {
	if update.Message == nil {
		return errors.New("message is nil")
	}

	chatID := update.Message.Chat.ID
	userName := update.Message.Chat.UserName

	err := q.storage.Delete(chatID)
	if err != nil {
		return err
	}

	steps := make([]*models.Step, 0, len(FlowSteps))
	for _, step := range FlowSteps {
		step := step
		steps = append(steps, &step)
	}

	chat := models.Chat{
		ID:       chatID,
		UserName: userName,
		Flow: &models.Flow{
			Steps: steps,
		},
	}
	err = q.storage.Save(chat)
	if err != nil {
		return err
	}

	return q.next(chat)
}

func (q *Questionnaire) Continue(update tgApi.Update) error {
	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		return errors.New("failed to parse chatID")
	}

	chat, err := q.storage.Get(chatID)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			fmt.Println("try to continue chat that not found", chatID, update)
			return nil
		} else {
			return err
		}
	}

	// continue finished chat
	if chat.Flow.IsFinished() {
		return nil
	}

	err = q.saveAnswer(&chat, update)
	if err != nil {
		return err
	}

	if chat.Flow.IsFinished() {
		return q.finish(chat)
	}

	return q.next(chat)
}

func (q *Questionnaire) next(chat models.Chat) error {
	step := chat.Flow.NextStep()
	if step == nil {
		return errors.New("failed to find next step")
	}

	msg := tgApi.NewMessage(chat.ID, step.Question)

	if len(step.Options) != 0 {
		keyboard := tgApi.InlineKeyboardMarkup{}
		for _, option := range step.Options {
			btn := tgApi.NewInlineKeyboardButtonData(option, option)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgApi.InlineKeyboardButton{btn})
		}

		msg.ReplyMarkup = keyboard
	}

	_, err := q.bot.Send(msg)
	return err
}

func (q *Questionnaire) saveAnswer(chat *models.Chat, update tgApi.Update) error {
	step := chat.Flow.NextStep()
	if step == nil {
		return errors.New("failed to find next step")
	}

	if len(step.Options) != 0 {
		if update.CallbackQuery != nil {
			step.Answer = update.CallbackQuery.Data
		} else {
			return q.sorry(*chat)
		}
	} else {
		if update.Message != nil {
			step.Answer = update.Message.Text
		} else {
			return q.sorry(*chat)
		}
	}

	return q.storage.Save(*chat)
}

func (q *Questionnaire) finish(chat models.Chat) error {
	msg := tgApi.NewMessage(chat.ID, finalMessage)
	_, err := q.bot.Send(msg)
	if err != nil {
		return err
	}

	text := fmt.Sprintf("Новая заявка от @%s!\n", chat.UserName)
	for _, step := range chat.Flow.Steps {
		text += fmt.Sprintf("*%s*: %s\n", step.Name, step.Answer)
	}

	msg = tgApi.NewMessage(q.cfg.MainChatID, text)
	msg.ParseMode = "markdown"
	_, err = q.bot.Send(msg)
	if err != nil {
		return err
	}

	msg = tgApi.NewMessage(q.cfg.TechChatID, text)
	msg.ParseMode = "markdown"
	_, err = q.bot.Send(msg)
	return err
}

func (q *Questionnaire) sorry(chat models.Chat) error {
	msg := tgApi.NewMessage(chat.ID, sorryMessage)
	_, err := q.bot.Send(msg)
	return err
}
