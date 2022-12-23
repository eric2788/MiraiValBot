package modules

import (
	"fmt"
	"testing"
)

type (
	IAnimal interface {
		Alive()
	}

	IDog interface {
		//IAnimal
		Walk()
	}

	Dogge struct {
	}
)

func (d *Dogge) Alive() {
	fmt.Println("doggle is alive!!")
}

func (d *Dogge) Walk() {
	fmt.Println("doggle is walking!!")
}

func TestModulePolymorphism(t *testing.T) {
	var animal IAnimal = &Dogge{}
	animal.Alive()
	if dog, ok := animal.(IDog); ok {
		dog.Walk()
	}
}
