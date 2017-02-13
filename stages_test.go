package appix

import (
	"fmt"
	"log"
	"testing"
)

func TestCreateStagePoolSuccessCase(t *testing.T) {
	var cmds []Stage

	// execute "go version". Supposed to work across all-platform
	for i := 0; i < 10; i++ {
		cmds = append(cmds, Stage{
			Cmd:  "go version",
			Name: fmt.Sprintf("command #%d", i),
		})
	}

	failed := <-CreateStagePool(cmds)

	if failed {
		t.Errorf("An error happened while running the commands")
		t.FailNow()
	} else {
		log.Printf("The test passed successfully")
	}
}

func TestCreateStagePoolFailCase(t *testing.T) {
	var cmds []Stage

	// append fake command to generate a fail
	cmds = append(cmds, Stage{
		Cmd:  "fakeCommand version",
		Name: "failing command",
	})

	for i := 0; i < 10; i++ {
		cmds = append(cmds, Stage{
			Cmd:  "go version",
			Name: fmt.Sprintf("command #%d", i),
		})
	}

	failed := <-CreateStagePool(cmds)

	if failed {
		log.Printf("The test passed successfully")
	} else {
		t.Errorf("An error happened while running the commands")
		t.FailNow()
	}
}

func BenchmarkStagePool(b *testing.B) {
	var cmds []Stage

	// execute "go version". Supposed to work across all-platform
	for i := 0; i < 10; i++ {
		cmds = append(cmds, Stage{
			Cmd:  "go version",
			Name: fmt.Sprintf("command #%d", i),
		})
	}

	for i := 0; i < 10; i++ {
		failed := <-CreateStagePool(cmds)

		if failed {
			b.Errorf("An error happened while running the commands")
			b.FailNow()
		} else {
			log.Printf("The test passed successfully")
		}
	}
}
