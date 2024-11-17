package dog

import "html/template"

var dogListTableTemplate = template.Must(template.New("table").Parse(`
{{range .}}
<tr>
  <td>{{.Name}}</td>
  <td>{{.Breed}}</td>
  <td class="buttons">
	<button
		class="show-on-hover"
		hx-delete="/dogs/{{.ID}}"
		hx-confirm="Are you sure?"
		hx-target="closest tr"
		hx-on:htmx:after-request="htmx.trigger('#dog-table', 'refresh-list')"
	>
		x
	</button>
</tr>
{{end}}
`))
