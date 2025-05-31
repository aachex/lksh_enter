# Отбор в ЛКШ, параллель P
## Переменные окружения
Ниже приведены следующие переменные окружения, используемые в приложении:
- API_HOST - хост сервера, предоставляющего API
- API_TOKEN - токен доступа к API
- SRV_PORT - порт, на котором будет запущено приложение

## Общее описание папок
В репозитории хранятся три папки:
- base - базовая часть
- advanced - продвинутая часть
- general - общая часть, здесь хранится код, используемый как base, так и advanced

## Подробное описание
### Пакет general
Общий код для базовой и продвинутой частей. В данном пакете хранятся следующие сущности:
- Client - является расширением структуры http.Client. Включает в себя методы, работающие с предоставленным API.
- Player - структура игрока
- Team - структура команды
- Match - структура матча
- Goal - структура гола, забитого игроком

### Пакет base
Базовая часть приложения. Логика работы следующая:<br><br>
Сначала создаётся клиент для взаимодействия с API, который затем получает отсортированный список игроков и выводит его:
```go
client := general.Client{}

// вывод игроков
players, err := client.PlayerNamesSorted()
if err != nil {
	panic(err)
}
for _, p := range players {
    fmt.Println(p)
}
```
Далее в бесконечном цикле пользователь вводит запрос на получение статистики команды stats? или versus? на получение количества матчей, где игрок 1 играл против игрока 2
```go
var s string
in := bufio.NewReader(os.Stdin)
for {
	fmt.Scan(&s)
	switch s {
	case "stats?":
		teamName, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}
		teamName = teamName[1 : len(teamName)-3] // убираем кавычки
		wins, defeats, scored, missed := client.GetStats(client.TeamId(teamName))
		fmt.Println(wins, defeats, scored-missed)

	case "versus?":
		var id1, id2 int
		fmt.Scan(&id1, &id2)
		fmt.Println(client.Versus(client.PlayerTeam(id1), client.PlayerTeam(id2)))
	}
}
```

### Пакет advanced
Представляет собой веб-сервер, включающим функционал base-пакета, но расширенный
#### Логгирование
Логгирование реализовано с помощью пакета slog:
```go
logOpts := slog.HandlerOptions{
	Level: slog.LevelDebug,
}
logger := slog.New(slog.NewJSONHandler(os.Stdout, &logOpts))
```
Все приходящие на сервер запросы логируются с помощью функции Middleware пакета logging:
```go
func Middleware(next http.HandlerFunc, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(
			"New request",
			"endpoint", r.URL.String(),
			"method", r.Method,
			"user-agent", r.UserAgent(),
		)

        // logReponseWriter - структура, созданная на основе http.ResponseWriter. Она необходима для доступа к коду статуса после обработки запроса
		lrw := &logReponseWriter{w, http.StatusOK}
		next(lrw, r)

		logger.Debug(
			"Request processed",
			"statusCode", lrw.statusCode,
		)
	}
}
```

#### HTML представления
HTML представления реализованы с помощью пакета *html/template*. Эндпоинты, имеющие постфикс **Html**, возвращают html разметку, но не json

## Запуск приложения
Пожалуйста, запускайте приложение из корня:    
- ```go run base\main.go``` для запуска базовой части<br>
- ```go run advanced\main.go``` для запуска продвинутой части

