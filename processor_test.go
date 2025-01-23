package mtapp

import (
	"context"
	"sync"
	"testing"
	"time"
)

// Тестирование добавления и удаления потоков, а также работы процессора
func TestProcessor(t *testing.T) {
	numTicks := 0

	// Создаем процессор
	processor := NewProcessor()

	// Создаем несколько потоков
	thread1 := NewThread("first", NewProcess(func(ctx context.Context) (cancelFunc context.CancelFunc) {
		ctx, cancelFunc = context.WithCancel(ctx)

		// Логируем, чтобы отслеживать, что происходит
		numTicks++
		t.Log("First thread tick")

		return
	}), 100*time.Millisecond, 1)

	thread2 := NewThread("second", NewProcess(func(ctx context.Context) (cancelFunc context.CancelFunc) {
		ctx, cancelFunc = context.WithCancel(ctx)

		// Логируем, чтобы отслеживать, что происходит
		numTicks++
		t.Log("Second thread tick")

		return
	}), 100*time.Millisecond, 1)

	// Добавляем потоки в процессор
	processor.AddThread(thread1)
	processor.AddThread(thread2)

	// Проверяем, что потоки были добавлены
	if len(processor.threads) != 2 {
		t.Errorf("expected 2 threads, got %d", len(processor.threads))
	}

	// Запускаем процессор в отдельной горутине
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor.Start(ctx, wg)

	// Подождем немного, чтобы потоки успели выполнить свою работу
	time.Sleep(300 * time.Millisecond)

	// Проверяем, что потоки удалились по завершению
	if len(processor.threads) != 0 {
		t.Errorf("expected 0 threads, got %d", len(processor.threads))
	}

	// Добавляем новый поток и проверяем его наличие
	thread3 := NewThread("third", NewProcess(func(ctx context.Context) (cancelFunc context.CancelFunc) {
		ctx, cancelFunc = context.WithCancel(ctx)

		// Логируем, чтобы отслеживать, что происходит
		numTicks++
		t.Log("Third thread tick")

		return
	}), 100*time.Millisecond, 0)
	processor.AddThread(thread3)

	// Проверяем, что поток добавлен
	if len(processor.threads) != 1 {
		t.Errorf("expected 1 thread, got %d", len(processor.threads))
	}

	// Проверяем, что поток можно найти
	foundThread := processor.Thread(thread3.ID())
	if foundThread == nil {
		t.Errorf("expected to find thread with ID %v", thread3.ID().String())
	}

	// Добавляем дополнительную задержку для наблюдения
	time.Sleep(100 * time.Millisecond)

	// Удаляем поток
	processor.DeleteThread(thread3.ID())

	// Подождем немного, чтобы потоки успели выполнить свою работу
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что поток удален
	if len(processor.threads) != 0 {
		t.Errorf("expected 0 threads, got %d", len(processor.threads))
	}

	// Останавливаем процессор
	cancel()
	wg.Wait()

	// Проверяем, что все потоки завершены
	if len(processor.threads) != 0 {
		t.Errorf("expected 0 threads, got %d", len(processor.threads))
	}

	// Проверяем, что счетчик сработал
	if numTicks != 4 {
		t.Errorf("expected 3 ticks, got %d", numTicks)
	}
}
