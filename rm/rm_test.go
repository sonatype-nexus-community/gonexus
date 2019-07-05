package nexusrm

func ExampleListComponents() {
	rm, err := nexusrm.New("http://localhost:8081", "user", "password")
	if err != nil {
		panic(err)
	}

	items, err := rm.ListComponents("maven-central")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", items)
}
