package repeater

import (
	"context"
	"sync"
	"testing"
	"time"
)

// Тестирование добавления и удаления потоков, а также работы репитера
func TestRepeater(t *testing.T) {
	// Создаем процессор
	repeater := Init("Main", 100*time.Millisecond, 3, Tick)

	// Запускаем процессор в отдельной горутине
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repeater.Start(ctx, wg)

	// Подождем немного, чтобы поток успел выполнить свою работу
	time.Sleep(500 * time.Millisecond)

	// Останавливаем процессор
	cancel()
	wg.Wait()
}
