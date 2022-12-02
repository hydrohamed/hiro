package hiro

func Get(source string) {
	Add(source, "main")
	Start("main", true)
}
