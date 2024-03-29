// Code generated by qtc from "helpers.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package templates

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

func streamimportFontAwesome(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.2/css/all.min.css" integrity="sha512-z3gLpd7yknf1YoNbCzqRKc4qyor8gaKU1qmn+CShxbuBusANI9QpRohGBreCFkKxLhei6S9CQXFEbbKuqLg0DA==" crossorigin="anonymous" referrerpolicy="no-referrer" />
`)
}

func writeimportFontAwesome(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamimportFontAwesome(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func importFontAwesome() string {
	qb422016 := qt422016.AcquireByteBuffer()
	writeimportFontAwesome(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamimportJquery(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.7.1/jquery.min.js" integrity="sha512-v2CJ7UaYy4JwqLDIrZUI/4hqeoQieOmAZNXBeQyjo21dadnwR+8ZaIJVT8EE2iyI61OV8e6M8PP2/4hpQINQ/g==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
`)
}

func writeimportJquery(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamimportJquery(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func importJquery() string {
	qb422016 := qt422016.AcquireByteBuffer()
	writeimportJquery(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamimportBootstrapTableBulmaCSS(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
<link rel="stylesheet" href="https://unpkg.com/bootstrap-table@1.22.1/dist/themes/bulma/bootstrap-table-bulma.min.css">
`)
}

func writeimportBootstrapTableBulmaCSS(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamimportBootstrapTableBulmaCSS(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func importBootstrapTableBulmaCSS() string {
	qb422016 := qt422016.AcquireByteBuffer()
	writeimportBootstrapTableBulmaCSS(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamimportBootstrapTableJS(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
<script src="https://unpkg.com/bootstrap-table@1.22.1/dist/bootstrap-table.min.js"></script>
<script src="https://unpkg.com/bootstrap-table@1.22.1/dist/themes/bulma/bootstrap-table-bulma.min.js"></script>
`)
}

func writeimportBootstrapTableJS(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamimportBootstrapTableJS(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func importBootstrapTableJS() string {
	qb422016 := qt422016.AcquireByteBuffer()
	writeimportBootstrapTableJS(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
