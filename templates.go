package main

import "text/template"

var rootTemplate = template.Must(template.New("root").Parse(
	`{{ if .Alive }}{{ .Emoji }}{{ else }}☠️{{ end }} {{ .Name }}

{{ .AgeString }}🕑
{{ .WeightString }}⚖️
{{ .MoodString }}💭
{{ .HealthString }}💗
{{ .FoodString }}🍽️
{{ .HappyString }}😶
`))

var feedTemplate = template.Must(template.New("feed").Parse(
	`{{ .FoodString }}
What do you prefer?
`))

var playTemplate = template.Must(template.New("feed").Parse(
	`{{ .HappyString }}
Let's play!
`))

var healTemplate = template.Must(template.New("heal").Parse(
	`{{ .HealthString }}
Heal me...
`))

var topTemplate = template.Must(template.New("top").Parse(
	`🏆 Top Survivors

{{ range . }}{{ .TopString }}
{{ end }}
`))
