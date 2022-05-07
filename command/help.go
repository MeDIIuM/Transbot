package command

const HelpMsg = `Правила пользования ботом:

		Это бот телеграмм для работы с Ethereum:
		/balance - Поможет проверить балланс счета
		/transfer - Поможет перевести деньги
		/help - отобразить это сообщение
	`

func Help() string {
	return HelpMsg
}
