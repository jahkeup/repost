// This is a generated file, do not edit.
package emailfuncs


// Name returns the Name from the parsed Email address provided.
func Name(str string) (string, error) {
	a, err := Parse(str)
	if err != nil {
		return "", err
	}
	return a.Name()
}


// Tag returns the Tag from the parsed Email address provided.
func Tag(str string) (string, error) {
	a, err := Parse(str)
	if err != nil {
		return "", err
	}
	return a.Tag()
}


// User returns the User from the parsed Email address provided.
func User(str string) (string, error) {
	a, err := Parse(str)
	if err != nil {
		return "", err
	}
	return a.User()
}


// Domain returns the Domain from the parsed Email address provided.
func Domain(str string) (string, error) {
	a, err := Parse(str)
	if err != nil {
		return "", err
	}
	return a.Domain()
}

