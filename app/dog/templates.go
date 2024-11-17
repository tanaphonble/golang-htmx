package dog

import "html/template"

var dogListTableTemplate = template.Must(template.New("table").Parse(`
{{range .}}
<tr>
  <td>{{.Name}}</td>
  <td>{{.Breed}}</td>
</tr>
{{end}}
`))
