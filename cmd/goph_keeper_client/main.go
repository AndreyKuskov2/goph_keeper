package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"goph_keeper/internal/client"
)

var (
	version   = "1.0.0"
	buildTime = time.Now().Format("2006-01-02 15:04:05")
	gitCommit = "unknown"
)

func main() {
	// Показываем информацию о версии при запуске
	showVersionInfo()

	// Создаем клиент
	cli := client.NewCLIClient("http://localhost:9000")

	// Главное меню авторизации
	for {
		choice := showAuthMenu()

		switch choice {
		case 1:
			// Регистрация
			if err := handleRegistration(cli); err != nil {
				fmt.Printf("Ошибка регистрации: %v\n", err)
				pause()
				continue
			}
			fmt.Println("Регистрация успешна!")
			pause()

		case 2:
			// Авторизация
			if err := handleLogin(cli); err != nil {
				fmt.Printf("Ошибка авторизации: %v\n", err)
				pause()
				continue
			}
			fmt.Println("Авторизация успешна!")
			pause()

		case 3:
			// Информация о версии
			showVersionInfo()
			pause()
			continue

		case 4:
			// Выход
			fmt.Println("До свидания!")
			return

		default:
			fmt.Println("Неверный выбор. Попробуйте снова.")
			pause()
			continue
		}

		// Если дошли сюда, значит авторизация прошла успешно
		// Переходим в основное меню
		mainMenu(cli)
		break
	}
}

func showVersionInfo() {
	fmt.Println("GophKeeper CLI Client")
	fmt.Println("========================================")
	fmt.Printf("Версия: %s\n", version)
	fmt.Printf("Дата сборки: %s\n", buildTime)
	fmt.Printf("Git commit: %s\n", gitCommit)
	fmt.Printf("Go версия: %s\n", runtime.Version())
	fmt.Printf("Платформа: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("========================================")
	fmt.Println()
}

func showAuthMenu() int {
	fmt.Println("ГЛАВНОЕ МЕНЮ")
	fmt.Println("========================================")
	fmt.Println("1. Регистрация")
	fmt.Println("2. Авторизация")
	fmt.Println("3. Информация о версии")
	fmt.Println("4. Выход")
	fmt.Println("========================================")

	return getChoice(1, 4)
}

func handleRegistration(cli *client.CLIClient) error {
	fmt.Println("\nРЕГИСТРАЦИЯ")
	fmt.Println("========================================")

	login := getInput("Введите логин: ")
	if login == "" {
		return fmt.Errorf("логин не может быть пустым")
	}

	password := getPassword("Введите пароль: ")
	if password == "" {
		return fmt.Errorf("пароль не может быть пустым")
	}

	return cli.Register(login, password)
}

func handleLogin(cli *client.CLIClient) error {
	fmt.Println("\nАВТОРИЗАЦИЯ")
	fmt.Println("========================================")

	login := getInput("Введите логин: ")
	if login == "" {
		return fmt.Errorf("логин не может быть пустым")
	}

	password := getPassword("Введите пароль: ")
	if password == "" {
		return fmt.Errorf("пароль не может быть пустым")
	}

	token, err := cli.Login(login, password)
	if err != nil {
		return err
	}

	// Сохраняем токен
	if err := cli.SaveToken(token); err != nil {
		fmt.Printf("Предупреждение: не удалось сохранить токен: %v\n", err)
	}

	return nil
}

func mainMenu(cli *client.CLIClient) {
	for {
		fmt.Println("\nОСНОВНОЕ МЕНЮ")
		fmt.Println("========================================")
		fmt.Println("1. Учетные данные")
		fmt.Println("2. Текстовые данные")
		fmt.Println("3. Банковские карты")
		fmt.Println("4. Бинарные данные")
		fmt.Println("5. Информация о версии")
		fmt.Println("6. Выход")
		fmt.Println("========================================")

		choice := getChoice(1, 6)

		switch choice {
		case 1:
			credentialsMenu(cli)
		case 2:
			textDataMenu(cli)
		case 3:
			bankCardsMenu(cli)
		case 4:
			binariesMenu(cli)
		case 5:
			showVersionInfo()
			pause()
		case 6:
			fmt.Println("До свидания!")
			return
		}
	}
}

func credentialsMenu(cli *client.CLIClient) {
	for {
		fmt.Println("\nУЧЕТНЫЕ ДАННЫЕ")
		fmt.Println("========================================")
		fmt.Println("1. Создать")
		fmt.Println("2. Показать все")
		fmt.Println("3. Показать по ID")
		fmt.Println("4. Обновить")
		fmt.Println("5. Удалить")
		fmt.Println("6. Назад")
		fmt.Println("========================================")

		choice := getChoice(1, 6)

		switch choice {
		case 1:
			createCredentials(cli)
		case 2:
			listCredentials(cli)
		case 3:
			getCredentials(cli)
		case 4:
			updateCredentials(cli)
		case 5:
			deleteCredentials(cli)
		case 6:
			return
		}
	}
}

func textDataMenu(cli *client.CLIClient) {
	for {
		fmt.Println("\nТЕКСТОВЫЕ ДАННЫЕ")
		fmt.Println("========================================")
		fmt.Println("1. Создать")
		fmt.Println("2. Показать все")
		fmt.Println("3. Показать по ID")
		fmt.Println("4. Обновить")
		fmt.Println("5. Удалить")
		fmt.Println("6. Назад")
		fmt.Println("========================================")

		choice := getChoice(1, 6)

		switch choice {
		case 1:
			createTextData(cli)
		case 2:
			listTextData(cli)
		case 3:
			getTextData(cli)
		case 4:
			updateTextData(cli)
		case 5:
			deleteTextData(cli)
		case 6:
			return
		}
	}
}

func bankCardsMenu(cli *client.CLIClient) {
	for {
		fmt.Println("\nБАНКОВСКИЕ КАРТЫ")
		fmt.Println("========================================")
		fmt.Println("1. Создать")
		fmt.Println("2. Показать все")
		fmt.Println("3. Показать по ID")
		fmt.Println("4. Обновить")
		fmt.Println("5. Удалить")
		fmt.Println("6. Назад")
		fmt.Println("========================================")

		choice := getChoice(1, 6)

		switch choice {
		case 1:
			createBankCard(cli)
		case 2:
			listBankCards(cli)
		case 3:
			getBankCard(cli)
		case 4:
			updateBankCard(cli)
		case 5:
			deleteBankCard(cli)
		case 6:
			return
		}
	}
}

func binariesMenu(cli *client.CLIClient) {
	for {
		fmt.Println("\nБИНАРНЫЕ ДАННЫЕ")
		fmt.Println("========================================")
		fmt.Println("1. Загрузить файл")
		fmt.Println("2. Показать все")
		fmt.Println("3. Скачать файл")
		fmt.Println("4. Обновить файл")
		fmt.Println("5. Удалить файл")
		fmt.Println("6. Назад")
		fmt.Println("========================================")

		choice := getChoice(1, 6)

		switch choice {
		case 1:
			createBinaryData(cli)
		case 2:
			listBinaryData(cli)
		case 3:
			getBinaryData(cli)
		case 4:
			updateBinaryData(cli)
		case 5:
			deleteBinaryData(cli)
		case 6:
			return
		}
	}
}

// Функции для работы с учетными данными
func createCredentials(cli *client.CLIClient) {
	fmt.Println("\nСОЗДАНИЕ УЧЕТНЫХ ДАННЫХ")
	fmt.Println("========================================")

	login := getInput("Логин: ")
	password := getPassword("Пароль: ")
	meta := getInput("Описание (необязательно): ")

	if err := cli.CreateCredentials(login, password, meta); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Учетные данные созданы успешно!")
	}
	pause()
}

func listCredentials(cli *client.CLIClient) {
	fmt.Println("\nСПИСОК УЧЕТНЫХ ДАННЫХ")
	fmt.Println("========================================")

	if err := cli.ListCredentials(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	pause()
}

func getCredentials(cli *client.CLIClient) {
	fmt.Println("\nПОИСК УЧЕТНЫХ ДАННЫХ")
	fmt.Println("========================================")

	id := getInput("Введите ID: ")
	if id == "" {
		fmt.Println("ID не может быть пустым")
		pause()
		return
	}

	if err := cli.GetCredentials(id); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	pause()
}

func updateCredentials(cli *client.CLIClient) {
	fmt.Println("\nОБНОВЛЕНИЕ УЧЕТНЫХ ДАННЫХ")
	fmt.Println("========================================")

	id := getInput("Введите ID: ")
	login := getInput("Новый логин: ")
	password := getPassword("Новый пароль: ")
	meta := getInput("Новое описание (необязательно): ")

	if err := cli.UpdateCredentials(id, login, password, meta); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Учетные данные обновлены успешно!")
	}
	pause()
}

func deleteCredentials(cli *client.CLIClient) {
	fmt.Println("\nУДАЛЕНИЕ УЧЕТНЫХ ДАННЫХ")
	fmt.Println("========================================")

	id := getInput("Введите ID для удаления: ")
	if id == "" {
		fmt.Println("ID не может быть пустым")
		pause()
		return
	}

	if !confirmAction("Вы уверены, что хотите удалить эти данные?") {
		fmt.Println("Удаление отменено")
		pause()
		return
	}

	if err := cli.DeleteCredentials(id); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Учетные данные удалены успешно!")
	}
	pause()
}

// Функции для работы с текстовыми данными
func createTextData(cli *client.CLIClient) {
	fmt.Println("\nСОЗДАНИЕ ТЕКСТОВОЙ ЗАМЕТКИ")
	fmt.Println("========================================")

	text := getInput("Введите текст: ")
	meta := getInput("Описание (необязательно): ")

	if err := cli.CreateTextData(text, meta); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Текстовая заметка создана успешно!")
	}
	pause()
}

func listTextData(cli *client.CLIClient) {
	fmt.Println("\nСПИСОК ТЕКСТОВЫХ ДАННЫХ")
	fmt.Println("========================================")

	if err := cli.ListTextData(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	pause()
}

func getTextData(cli *client.CLIClient) {
	fmt.Println("\nПОИСК ТЕКСТОВЫХ ДАННЫХ")
	fmt.Println("========================================")

	id := getInput("Введите ID: ")
	if id == "" {
		fmt.Println("ID не может быть пустым")
		pause()
		return
	}

	if err := cli.GetTextData(id); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	pause()
}

func updateTextData(cli *client.CLIClient) {
	fmt.Println("\nОБНОВЛЕНИЕ ТЕКСТОВЫХ ДАННЫХ")
	fmt.Println("========================================")

	id := getInput("Введите ID: ")
	text := getInput("Новый текст: ")
	meta := getInput("Новое описание (необязательно): ")

	if err := cli.UpdateTextData(id, text, meta); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Текстовые данные обновлены успешно!")
	}
	pause()
}

func deleteTextData(cli *client.CLIClient) {
	fmt.Println("\nУДАЛЕНИЕ ТЕКСТОВЫХ ДАННЫХ")
	fmt.Println("========================================")

	id := getInput("Введите ID для удаления: ")
	if id == "" {
		fmt.Println("ID не может быть пустым")
		pause()
		return
	}

	if !confirmAction("Вы уверены, что хотите удалить эти данные?") {
		fmt.Println("Удаление отменено")
		pause()
		return
	}

	if err := cli.DeleteTextData(id); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Текстовые данные удалены успешно!")
	}
	pause()
}

// Функции для работы с банковскими картами
func createBankCard(cli *client.CLIClient) {
	fmt.Println("\nСОЗДАНИЕ БАНКОВСКОЙ КАРТЫ")
	fmt.Println("========================================")

	cardNumber := getInput("Номер карты: ")
	holder := getInput("Владелец карты: ")
	cvc := getPassword("CVC код: ")
	expirationDate := getInput("Дата истечения (YYYY-MM-DD): ")
	meta := getInput("Описание (необязательно): ")

	if err := cli.CreateBankCard(cardNumber, holder, cvc, expirationDate, meta); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Банковская карта создана успешно!")
	}
	pause()
}

func listBankCards(cli *client.CLIClient) {
	fmt.Println("\nСПИСОК БАНКОВСКИХ КАРТ")
	fmt.Println("========================================")

	if err := cli.ListBankCards(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	pause()
}

func getBankCard(cli *client.CLIClient) {
	fmt.Println("\nПОИСК БАНКОВСКОЙ КАРТЫ")
	fmt.Println("========================================")

	id := getInput("Введите ID: ")
	if id == "" {
		fmt.Println("ID не может быть пустым")
		pause()
		return
	}

	if err := cli.GetBankCard(id); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	pause()
}

func updateBankCard(cli *client.CLIClient) {
	fmt.Println("\nОБНОВЛЕНИЕ БАНКОВСКОЙ КАРТЫ")
	fmt.Println("========================================")

	id := getInput("Введите ID: ")
	cardNumber := getInput("Новый номер карты: ")
	holder := getInput("Новый владелец карты: ")
	cvc := getPassword("Новый CVC код: ")
	expirationDate := getInput("Новая дата истечения (YYYY-MM-DD): ")
	meta := getInput("Новое описание (необязательно): ")

	if err := cli.UpdateBankCard(id, cardNumber, holder, cvc, expirationDate, meta); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Банковская карта обновлена успешно!")
	}
	pause()
}

func deleteBankCard(cli *client.CLIClient) {
	fmt.Println("\nУДАЛЕНИЕ БАНКОВСКОЙ КАРТЫ")
	fmt.Println("========================================")

	id := getInput("Введите ID для удаления: ")
	if id == "" {
		fmt.Println("ID не может быть пустым")
		pause()
		return
	}

	if !confirmAction("Вы уверены, что хотите удалить эти данные?") {
		fmt.Println("Удаление отменено")
		pause()
		return
	}

	if err := cli.DeleteBankCard(id); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Банковская карта удалена успешно!")
	}
	pause()
}

// Функции для работы с бинарными данными
func createBinaryData(cli *client.CLIClient) {
	fmt.Println("\nЗАГРУЗКА ФАЙЛА")
	fmt.Println("========================================")

	filePath := getInput("Путь к файлу: ")
	meta := getInput("Описание (необязательно): ")

	if err := cli.CreateBinaryData(filePath, meta); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Файл загружен успешно!")
	}
	pause()
}

func listBinaryData(cli *client.CLIClient) {
	fmt.Println("\nСПИСОК БИНАРНЫХ ДАННЫХ")
	fmt.Println("========================================")

	if err := cli.ListBinaryData(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	pause()
}

func getBinaryData(cli *client.CLIClient) {
	fmt.Println("\nСКАЧИВАНИЕ ФАЙЛА")
	fmt.Println("========================================")

	id := getInput("Введите ID: ")
	if id == "" {
		fmt.Println("ID не может быть пустым")
		pause()
		return
	}

	outputFile := getInput("Имя файла для сохранения (необязательно): ")

	if err := cli.GetBinaryData(id, outputFile); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Файл скачан успешно!")
	}
	pause()
}

func updateBinaryData(cli *client.CLIClient) {
	fmt.Println("\nОБНОВЛЕНИЕ ФАЙЛА")
	fmt.Println("========================================")

	id := getInput("Введите ID: ")
	filePath := getInput("Путь к новому файлу: ")
	meta := getInput("Новое описание (необязательно): ")

	if err := cli.UpdateBinaryData(id, filePath, meta); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Файл обновлен успешно!")
	}
	pause()
}

func deleteBinaryData(cli *client.CLIClient) {
	fmt.Println("\nУДАЛЕНИЕ ФАЙЛА")
	fmt.Println("========================================")

	id := getInput("Введите ID для удаления: ")
	if id == "" {
		fmt.Println("ID не может быть пустым")
		pause()
		return
	}

	if !confirmAction("Вы уверены, что хотите удалить эти данные?") {
		fmt.Println("Удаление отменено")
		pause()
		return
	}

	if err := cli.DeleteBinaryData(id); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Файл удален успешно!")
	}
	pause()
}

// Вспомогательные функции
func getChoice(min, max int) int {
	for {
		fmt.Printf("Выберите опцию (%d-%d): ", min, max)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil || choice < min || choice > max {
			fmt.Printf("Неверный выбор. Введите число от %d до %d\n", min, max)
			continue
		}

		return choice
	}
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getPassword(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func confirmAction(message string) bool {
	for {
		fmt.Printf("%s (y/n): ", message)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "y", "yes", "да", "д":
			return true
		case "n", "no", "нет", "н":
			return false
		default:
			fmt.Println("Введите 'y' для подтверждения или 'n' для отмены")
		}
	}
}

func pause() {
	fmt.Println("\nНажмите Enter для продолжения...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}
