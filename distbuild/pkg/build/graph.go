package build

import (
	"crypto/sha1"
	"encoding"
	"encoding/hex"
	"fmt"
	"path/filepath"
)

type ID [sha1.Size]byte

var (
	_ = encoding.TextMarshaler(ID{})
	_ = encoding.TextUnmarshaler(&ID{})
)

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

func (id ID) Path() string {
	return filepath.Join(hex.EncodeToString(id[:1]), hex.EncodeToString(id[:]))
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(id[:])), nil
}

func (id *ID) UnmarshalText(b []byte) error {
	raw, err := hex.DecodeString(string(b))
	if err != nil {
		return err
	}

	if len(raw) != len(id) {
		return fmt.Errorf("invalid id size: %q", b)
	}

	copy(id[:], raw)
	return nil
}

// Job описывает одну вершину графа сборки.
type Job struct {
	// ID задаёт уникальный идентификатор джоба.
	//
	// ID вычисляется как хеш от всех входных файлов, команд запуска и хешей зависимых джобов.
	//
	// Выход джоба целиком опеределяется его ID. Это важное свойство позволяет кешировать
	// результаты сборки.
	ID ID

	// Name задаёт человекочитаемое имя джоба.
	//
	// Например:
	//   build gitlab.com/slon/disbuild/pkg/b
	//   vet gitlab.com/slon/disbuild/pkg/a
	//   test gitlab.com/slon/disbuild/pkg/test
	Name string

	// Inputs задаёт список файлов из директории с исходным кодом,
	// которые нужны для работы этого джоба.
	//
	// В типичном случае, тут будут перечислены все .go файлы одного пакета.
	Inputs []string

	// Deps задаёт список джобов, выходы которых нужны для работы этого джоба.
	Deps []ID

	// Cmds описывает список команд, которые нужно выполнить в рамках этого джоба.
	Cmds []Cmd
}

// Cmd описывает одну команду сборки.
//
// Есть несколько видов команд. Все виды команд описываются одной структурой.
// Реальный тип определяется тем, какие поля структуры заполнены.
//
//   exec - выполняет произвольную команду
//   cat  - записывает строку в файл
//
// Все строки в описании команды могут содержать в себе на переменные. Перед выполнением
// реальной команды, переменные заменяются на их реальные значения.
//
//   {{OUTPUT_DIR}} - абсолютный путь до выходной директории джоба.
//   {{SOURCE_DIR}} - абсолютный путь до директории с исходными файлами.
//   {{DEP:f374b81d81f641c8c3d5d5468081ef83b2c7dae9}} - абсолютный путь до директории,
//   содержащей выход джоба с id f374b81d81f641c8c3d5d5468081ef83b2c7dae9.
type Cmd struct {
	// Exec описывает команду, которую нужно выполнить.
	Exec []string

	// Environ описывает переменные окружения, которые необходимы для работы команды из Exec.
	Environ []string

	// WorkingDirectory задаёт рабочую директорию для команды из Exec.
	WorkingDirectory string

	// CatTemplate задаёт шаблон строки, которую нужно записать в файл.
	CatTemplate string

	// CatOutput задаёт выходной файл для команды типа cat.
	CatOutput string
}

func (cmd Cmd) Render(outputDir, sourceDir string, deps map[ID]string) Cmd {
	panic("implement me")
}

type Graph struct {
	SourceFiles map[ID]string

	Jobs []Job
}
