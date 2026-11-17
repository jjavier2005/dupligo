package recycle

import (
	"fmt"

	"github.com/mono0x/trash"
)

func Send(paths []string) (int, error) {
	moved := 0

	for _, path := range paths {
		if err := trash.Move(path); err != nil {
			return moved, fmt.Errorf("Moviendo el archivo %q a la papelera de reciclaje: %w", path, err)
		}
		moved++
	}
	return moved, nil
}
