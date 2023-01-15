package web

var auth map[string]map[string]bool = map[string]map[string]bool{
	"Public": {
		"account": true,
		"captcha": true,
	},
	"Admin": {
		"user":     true,
		"account":  true,
		"captcha":  true,
		"home":     true,
		"order":    true,
		"follow":   true,
		"service":  true,
		"password": true,
		"profile":  true,
		"customer": true,
	},
	"User": {
		"account":  true,
		"captcha":  true,
		"home":     true,
		"order":    true,
		"follow":   true,
		"service":  true,
		"password": true,
		"profile":  true,
		"customer": true,
	},
}
