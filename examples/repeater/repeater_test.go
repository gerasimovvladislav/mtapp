package repeater

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gerasimovvladislav/mtapp"
	"github.com/stretchr/testify/assert"
)

// Тестирование добавления и удаления потоков, а также работы репитера
func TestRepeater(t *testing.T) {
	// Создаем процессор
	repeater := Init("Main", 100*time.Millisecond, 3, Tick)
	defaultThread := repeater.Thread("Main")

	// Проверяем, что поток был добавлен в процессор
	assert.NotNil(t, defaultThread)

	// Запускаем процессор в отдельной горутине
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repeater.Start(ctx, wg)

	// Подождем немного, чтобы поток успел выполнить свою работу
	time.Sleep(500 * time.Millisecond)

	// Проверяем, что потоки удалились по завершению
	assert.Nil(t, repeater.Thread("Main"))

	// Добавляем новый поток и проверяем его наличие
	otherThread := mtapp.NewThread("Other", mtapp.NewProcess(Tick), 100*time.Millisecond, 1)
	repeater.AddThread(otherThread)

	// Проверяем, что поток был добавлен в процессор
	assert.NotNil(t, repeater.Thread("Other"))

	// Останавливаем процессор
	cancel()
	wg.Wait()
}
