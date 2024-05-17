package easysyntax

type HandleErrorFunc func(error) bool

func Simple() HandleErrorFunc {
	return func(err error) bool {
		if err != nil {
			return false
		}
		return true
	}
}

func HandleMultiError(f HandleErrorFunc, errs ...error) error {
	for _, err := range errs {
		if !f(err) {
			return err
		}
	}
	return nil
}
