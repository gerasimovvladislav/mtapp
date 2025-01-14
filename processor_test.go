package mtapp

import (
	"context"
	"sync"
	"testing"
	"time"
)

// Тестирование добавления и удаления потоков, а также работы процессора
func TestProcessor(t *testing.T) {
	// Создаем процессор
	processor := NewProcessor()

	// Создаем несколько потоков
	thread1 := &Thread{id: ThreadID("first"), interval: 100 * time.Millisecond, limit: 2}
	thread2 := &Thread{id: ThreadID("second"), interval: 200 * time.Millisecond, limit: 1}

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
	thread3 := &Thread{id: ThreadID("third"), interval: 50 * time.Millisecond, limit: 3}
	processor.AddThread(thread3)

	// Проверяем, что поток добавлен
	if len(processor.threads) != 1 {
		t.Errorf("expected 1 thread, got %d", len(processor.threads))
	}

	// Проверяем, что поток можно найти
	foundThread := processor.GetThread(thread3.ID())
	if foundThread == nil {
		t.Errorf("expected to find thread with ID %v", thread3.ID().String())
	}

	// Останавливаем процессор
	cancel()
	wg.Wait()

	// Проверяем, что все потоки завершены
	if len(processor.threads) != 0 {
		t.Errorf("expected 0 threads, got %d", len(processor.threads))
	}
}
