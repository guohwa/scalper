package web

var auth map[string]map[string]bool = map[string]map[string]bool{
	"Public": {
		"account": true,
		"captcha": true,
	},
	"Admin": {
		"app":      true,
		"param":    true,
		"user":     true,
		"service":  true,
		"account":  true,
		"captcha":  true,
		"home":     true,
		"order":    true,
		"password": true,
		"profile":  true,
		"customer": true,
	},
	"User": {
		"account":  true,
		"captcha":  true,
		"home":     true,
		"order":    true,
		"password": true,
		"profile":  true,
		"customer": true,
	},
}
