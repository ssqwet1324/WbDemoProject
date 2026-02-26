package main

import "fmt"

// Logger Новый интерфейс, который ожидает фреймворк
type Logger interface {
	Log(message string)
}

// LegacyPrinter старый тип, который нельзя менять из допустим старой библиотеки
type LegacyPrinter struct{}

func (p *LegacyPrinter) Print(msg string) {
	fmt.Println("[Legacy]:", msg)
}

// PrinterAdapter делаем возможность совместить с новым интерфейсом
type PrinterAdapter struct {
	legacy *LegacyPrinter
}

func (a *PrinterAdapter) Log(message string) {
	a.legacy.Print(message)
}

// Process код, который работает только с Logger
func Process(logger Logger) {
	logger.Log("Обработка данных...")
	logger.Log("Готово!")
}

func main() {
	old := &LegacyPrinter{}
	adapter := &PrinterAdapter{legacy: old}

	// старый логгер можно использовать там, где требуется новый интерфейс
	Process(adapter)
}
