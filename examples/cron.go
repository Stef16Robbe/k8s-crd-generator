package examples

import (
	"fmt"
)

func example() {
	ct := CronTab{
		Spec: Spec{
			CronSpec: "",
			Image:    "",
			Replicas: 2,
		},
	}

	fmt.Printf("%+v\n", ct)
}
