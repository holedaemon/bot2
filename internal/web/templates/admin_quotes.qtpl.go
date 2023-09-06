// Code generated by qtc from "admin_quotes.qtpl". DO NOT EDIT.
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

type AdminQuotesPage struct {
	BasePage
}

func (p *AdminQuotesPage) StreamTitle(qw422016 *qt422016.Writer) {
	qw422016.N().S(`Import Quotes`)
}

func (p *AdminQuotesPage) WriteTitle(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamTitle(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *AdminQuotesPage) Title() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteTitle(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *AdminQuotesPage) StreamBody(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    <section class="section">
        <div class="container">
            <h1 class="title has-text-centered">
                Import quotes
            </h1>

            <div class="columns is-centered">
                <div class="column is-one-third has-text-centered">
                    <form id="bot2-quote-import-form" method="POST" action="/admin/quotes" autocomplete="off">
                        <div class="field">
                            <div class="control">
                                <div id="bot2-quote-import-file" class="file has-name is-centered">
                                    <label class="file-label">
                                        <input class="file-input" type="file" name="quotes">
                                        <span class="file-cta">
                                            <span class="file-icon">
                                                <i class="fas fa-upload"></i>
                                            </span>
                                            <span class="file-label">
                                                Pick a file, boss
                                            </span>
                                        </span>
                                        <span class="file-name">
                                            Placeholder!!!!
                                        </span>
                                    </label>
                                </div>
                            </div>
                        </div>
                        <div class="field">
                            <div class="control">
                                <button class="button is-success">Import</button>
                            </div>
                        </div>
                    </form>
                    <br>
                    <div id="bot2-output" class="has-text-left"></div>
                </div>
            </div>
        </div>
    </section>
    `)
	streamimportJquery(qw422016)
	qw422016.N().S(`
    <script>
        $(document).ready(function () {
            const fileInputs = $("#bot2-quote-import-file input[type=file]");
            const fileInput = fileInputs[0];
            fileInput.onchange = () => {
                if (fileInput.files.length > 0) {
                    const fileNames = $("#bot2-quote-import-file .file-name");
                    fileNames[0].textContent = fileInput.files[0].name;
                }
            }

            $("#bot2-quote-import-form").submit(function(event) {
                event.preventDefault();

                const formData = new FormData();
                if (fileInput.files.length > 0) {
                    formData.append("quotes", fileInput.files[0]);
                } else {
                    $("#bot2-output").addClass("box").text("You didn't select a file, numbnuts");
                    return;
                }

                $.ajax({
                    method: "POST",
                    url: "/admin/quotes",
                    data: formData,
                    contentType: false,
                    processData: false,
                }).done(function(res) {
                    $("#bot2-output").addClass("box").text(res);
                }).fail(function(req) {
                    $("#bot2-output").addClass("box").text(req.responseText);
                });
            })
        });
    </script>
`)
}

func (p *AdminQuotesPage) WriteBody(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamBody(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *AdminQuotesPage) Body() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteBody(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *AdminQuotesPage) StreamMeta(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	streamimportFontAwesome(qw422016)
	qw422016.N().S(`
`)
}

func (p *AdminQuotesPage) WriteMeta(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamMeta(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *AdminQuotesPage) Meta() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteMeta(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *AdminQuotesPage) StreamScripts(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	streamimportJquery(qw422016)
	qw422016.N().S(`
`)
}

func (p *AdminQuotesPage) WriteScripts(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamScripts(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *AdminQuotesPage) Scripts() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteScripts(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
