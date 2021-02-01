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
	startMessage = `Здравствуйте. Вас ждёт анкета из нескольких вопросов. Указанные данные будут переданы только в чат команды психологов.

В разделе "как обращаться" вы можете указать никнейм и остаться анонимной/ым до контакта с психологом.

Для запуска анкеты отправьте боту /go`

	cancelMessage = "Заполнение анкеты отменено"

	finalMessage = `Ваши ответы отправлены. Заявка обрабатывается. Психолог/психотерапевт свяжется с вами после получения этой информации.

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
		ID: chatID,
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

	text := fmt.Sprintf("Новая заявка от %d!\n", chat.ID)
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
