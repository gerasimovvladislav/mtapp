package mtapp

import (
	"context"
	"sync"
	"testing"
	"time"
)

// Тестирование добавления и удаления потоков, а также работы процессора
func TestProcessor(t *testing.T) {
	mu := &sync.Mutex{}
	numTicks := 0

	// Создаем несколько потоков
	thread1 := NewThread("first", NewProcess(func(ctx context.Context) (cancelFunc context.CancelFunc) {
		ctx, cancelFunc = context.WithCancel(ctx)

		// Логируем, чтобы отслеживать, что происходит
		mu.Lock()
		numTicks++
		mu.Unlock()
		t.Log("First thread tick")

		return
	}), false, 100*time.Millisecond, 1)

	thread2 := NewThread("second", NewProcess(func(ctx context.Context) (cancelFunc context.CancelFunc) {
		ctx, cancelFunc = context.WithCancel(ctx)

		// Логируем, чтобы отслеживать, что происходит
		mu.Lock()
		numTicks++
		mu.Unlock()
		t.Log("Second thread tick")

		return
	}), false, 100*time.Millisecond, 1)

	// Добавляем новый поток и проверяем его наличие
	thread3 := NewThread("third", NewProcess(func(ctx context.Context) (cancelFunc context.CancelFunc) {
		ctx, cancelFunc = context.WithCancel(ctx)

		// Логируем, чтобы отслеживать, что происходит
		mu.Lock()
		numTicks++
		mu.Unlock()
		t.Log("Third thread tick")

		return
	}), true, 100*time.Millisecond, 1)

	// Создаем процессор с потоками
	processor := NewProcessor(thread1, thread2, thread3)

	// Проверяем, что потоки были добавлены
	if len(processor.threads) != 3 {
		t.Errorf("expected 3 threads, got %d", len(processor.threads))
	}

	// Запускаем процессор в отдельной горутине
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor.Start(ctx, wg)

	// Подождем немного, чтобы потоки успели выполнить свою работу
	time.Sleep(400 * time.Millisecond)

	// Проверяем, что два потока были запущены
	if numTicks != 2 {
		t.Errorf("expected 2 ticks, got %d", numTicks)
	}

	// Снимаем с паузы третий поток
	thread3.Start(ctx, wg)

	// Подождем немного, чтобы потоки успели выполнить свою работу
	time.Sleep(200 * time.Millisecond)

	// Проверяем, что поток отработал
	if numTicks != 3 {
		t.Errorf("expected 3 ticks, got %d", numTicks)
	}

	// Останавливаем процессор
	cancel()
	wg.Wait()

	// Проверяем, что все потоки завершены
	if len(processor.threads) != 3 {
		t.Errorf("expected 3 threads, got %d", len(processor.threads))
	}
}
