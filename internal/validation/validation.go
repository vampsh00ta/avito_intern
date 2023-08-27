package validation

//	func Error(err error) error {
//		arrArr := err.(validator.ValidationErrors)
//		var errs []error
//		for _, e := range arrArr {
//			errs = append(errs,errors.New("field error:%s",))
//		}
//	}
func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}
