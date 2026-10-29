package flattenbson_test

import (
	"fmt"

	flat "github.com/SagarSingh2003/flatten-bson"
)

func ExampleFlatten() {
	type X struct {
		Quantity int `bson:"quantity"`
	}

	type Human struct {
		Eyes   *X `bson:"eyes,omitempty"`
		Tongue X  `bson:"tongue"`
	}

	fmt.Println(flat.Flatten(Human{Eyes: &X{Quantity: 2}, Tongue: X{}}))
	// Output: map[eyes.quantity:2 tongue.quantity:0] <nil>
}
