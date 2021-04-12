package easyerrors

func HandleMultiError(f func(error) bool, errs ...error) error {
	for _, err := range errs {
		if err != nil && !f(err) {
			return err
		}
	}
	return nil
}
