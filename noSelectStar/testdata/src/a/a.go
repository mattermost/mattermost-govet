package a

func Select(args ...string)  {}
func Columns(args ...string) {}

func example() {
	// These should trigger warnings
	_ = "SELECT * FROM users" // want `do not use SELECT \*: explicitly select the needed columns instead`
	_ = "select * from table" // want `do not use SELECT \*: explicitly select the needed columns instead`

	// These should not trigger warnings
	_ = "SELECT id, name FROM users"
	_ = "Just a * by itself"
	_ = "SELECT"

	// These should trigger warnings for Select function
	Select("*")     // want `do not use Select with \*: explicitly select the needed columns instead`
	Select("id, *") // want `do not use Select with \*: explicitly select the needed columns instead`
	Select("id", "*", "name") // want `do not use Select with \*: explicitly select the needed columns instead`
	Select("id", "name", "*") // want `do not use Select with \*: explicitly select the needed columns instead`

	// These should not trigger warnings for Select function
	Select("id, name")
	Select("")
	Select("id", "name", "email")

	// These should trigger warnings for Columns function
	Columns("*")     // want `do not use Columns with \*: explicitly select the needed columns instead`
	Columns("id, *") // want `do not use Columns with \*: explicitly select the needed columns instead`
	Columns("id", "*", "name") // want `do not use Columns with \*: explicitly select the needed columns instead`
	Columns("id", "name", "*") // want `do not use Columns with \*: explicitly select the needed columns instead`

	// These should not trigger warnings for Columns function
	Columns("id, name")
	Columns("")
	Columns("id", "name", "email")
}
