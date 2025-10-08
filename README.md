# TelegramBot
Что я умею:
🎯 Список дел (To-Do): Создавайте задачи, устанавливайте дедлайны и отслеживайте свой прогресс.
💰 Управление финансами: Контролируйте доходы и расходы, чтобы всегда быть в курсе своего бюджета.
✍️ Заметки: Быстро записывайте идеи, мысли и важную информацию, чтобы ничего не потерять.
🎲 Медиа-копилка: Сохраняйте интересные статьи, видео и ссылки, а я буду присылать их вам в случайном порядке для вдохновения или развлечения.

Начните прямо сейчас и наведите порядок в своих делах!

Для запуска 
go mod tidy
go build -o bot.exe
.\bot.exe -tg-bot-token "You token"

Либо запускать всегда так:
$env:TELEGRAM_TOKEN="ТОКЕН"; .\bot.exe

Либо сделать постоянную переменную один раз:
setx TELEGRAM_TOKEN "ТОКЕН"
(и перезапустить PowerShell, чтобы она подхватилась)

Или добавить поддержку флага -tg-bot-token, тогда запускать проще:

.\bot.exe -tg-bot-token "ТОКЕН"

What I can
do: 🎯 To-Do list: Create tasks, set deadlines, and track your progress.
💰 Financial management: Monitor income and expenses to always be aware of your budget.
✍️ Notes: Quickly write down ideas, thoughts, and important information so that you don't lose anything.
Media piggy bank: Save interesting articles, videos and links, and I will send them to you in random order for inspiration or entertainment.

Start right now and get your affairs in order!

To launch

go mod tidy
go build -o bot.exe
//.\bot.exe -tg-bot-token "You token"
