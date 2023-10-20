// Code generated by qtc from "base.qtpl". DO NOT EDIT.
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

type Page interface {
	Title() string
	StreamTitle(qw422016 *qt422016.Writer)
	WriteTitle(qq422016 qtio422016.Writer)
	Meta() string
	StreamMeta(qw422016 *qt422016.Writer)
	WriteMeta(qq422016 qtio422016.Writer)
	Navbar() string
	StreamNavbar(qw422016 *qt422016.Writer)
	WriteNavbar(qq422016 qtio422016.Writer)
	Scripts() string
	StreamScripts(qw422016 *qt422016.Writer)
	WriteScripts(qq422016 qtio422016.Writer)
	Body() string
	StreamBody(qw422016 *qt422016.Writer)
	WriteBody(qq422016 qtio422016.Writer)
}

func StreamPageTemplate(qw422016 *qt422016.Writer, p Page) {
	qw422016.N().S(`
<!DOCTYPE html>
<html>
    <head>
      <title>`)
	p.StreamTitle(qw422016)
	qw422016.N().S(` &mdash; DILF</title>

      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <link rel="stylesheet" href="/static/bot2.css">

      <link rel="apple-touch-icon" sizes="180x180" href="/static/apple-touch-icon.png?v=1234">
      <link rel="icon" type="image/png" sizes="32x32" href="/static/favicon-32x32.png?v=1234">
      <link rel="icon" type="image/png" sizes="16x16" href="/static/favicon-16x16.png?v=1234">
      <link rel="manifest" href="/static/site.webmanifest?v=1234">
      <link rel="mask-icon" href="/static/safari-pinned-tab.svg?v=1234" color="#5bbad5">
      <link rel="shortcut icon" href="/static/favicon.ico?v=1234">
      <meta name="msapplication-TileColor" content="#00aba9">
      <meta name="msapplication-config" content="/static/browserconfig.xml?v=1234">
      <meta name="theme-color" content="#ffffff">

      <meta property="og:title" content="DILF">
      <meta property="og:type" content="website">
      <meta property="og:image" content="https://bot.holedaemon.net/static/dilf.jpg">
      <meta property="og:url" content="https://bot.holedaemon.net">

      `)
	p.StreamMeta(qw422016)
	qw422016.N().S(`
    </head>
    <body>
        `)
	p.StreamNavbar(qw422016)
	qw422016.N().S(`

        `)
	p.StreamBody(qw422016)
	qw422016.N().S(`
    </body>

    `)
	p.StreamScripts(qw422016)
	qw422016.N().S(`

    <script>
      document.addEventListener("DOMContentLoaded", () => {
        const burgers = Array.prototype.slice.call(document.querySelectorAll(".navbar-burger"), 0);

        burgers.forEach(el => {
          el.addEventListener("click", () => {
            const target = el.dataset.target;
            const targetElement = document.getElementById(target);

            el.classList.toggle("is-active");
            targetElement.classList.toggle("is-active");
          });
        });
      });
    </script>
</html>
`)
}

func WritePageTemplate(qq422016 qtio422016.Writer, p Page) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	StreamPageTemplate(qw422016, p)
	qt422016.ReleaseWriter(qw422016)
}

func PageTemplate(p Page) string {
	qb422016 := qt422016.AcquireByteBuffer()
	WritePageTemplate(qb422016, p)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

type BasePage struct {
	Username string
}

func (p *BasePage) StreamTitle(qw422016 *qt422016.Writer) {
}

func (p *BasePage) WriteTitle(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamTitle(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *BasePage) Title() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteTitle(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *BasePage) StreamMeta(qw422016 *qt422016.Writer) {
}

func (p *BasePage) WriteMeta(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamMeta(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *BasePage) Meta() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteMeta(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *BasePage) StreamScripts(qw422016 *qt422016.Writer) {
}

func (p *BasePage) WriteScripts(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamScripts(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *BasePage) Scripts() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteScripts(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *BasePage) StreamBody(qw422016 *qt422016.Writer) {
}

func (p *BasePage) WriteBody(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamBody(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *BasePage) Body() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteBody(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *BasePage) StreamNavbar(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
  <div class="container">
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a class="navbar-item is-size-3" href="/">
          DILF
        </a>

        <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" data-target="bot2-main-navbar">
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>

      <div id="bot2-main-navbar" class="navbar-menu">
        <div class="navbar-end">
          <a href="/docs" class="navbar-item">Docs</a>
          <a href="/about" class="navbar-item">About</a>

          `)
	if p.Username == "" {
		qw422016.N().S(`
            <a href="/login" class="navbar-item">Log in</a>
          `)
	} else {
		qw422016.N().S(`
            <div class="navbar-item has-dropdown is-hoverable">
              <a class="navbar-link">
                Sup, `)
		qw422016.E().S(p.Username)
		qw422016.N().S(`?
              </a>

              <div class="navbar-dropdown">
                <a href="/guilds" class="navbar-item">
                  Guilds
                </a>
                <a href="/profile" class="navbar-item">
                  Profile
                </a>
                <a href="/logout" class="navbar-item">
                  Log out
                </a>
              </div>
            </div>
          `)
	}
	qw422016.N().S(`
        </div>
      </div>
    </nav>
  </div>
`)
}

func (p *BasePage) WriteNavbar(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamNavbar(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *BasePage) Navbar() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteNavbar(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
