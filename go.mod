module godev90/slicer

go 1.23.4

replace (
	godev90/validator => ../validator
)

require (
	godev90/validator v0.0.0-00010101000000-000000000000
	gorm.io/gorm v1.30.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
