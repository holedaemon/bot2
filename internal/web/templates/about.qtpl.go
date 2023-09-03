// Code generated by qtc from "about.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package templates

import "github.com/holedaemon/bot2/internal/version"

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

type AboutPage struct {
	BasePage
}

func (p *AboutPage) StreamTitle(qw422016 *qt422016.Writer) {
	qw422016.N().S(`About`)
}

func (p *AboutPage) WriteTitle(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamTitle(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *AboutPage) Title() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteTitle(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *AboutPage) StreamBody(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    <section class="section">
    <div class="container">
      <div class="columns is-mobile is-centered">
        <div class="column is-half has-background-success box">
          <h1 class="title">DILF</h1>
          <p class="subtitle">...is a Discord bot created and maintained by <a class="has-text-warning"
              href="https://holedaemon.net" target="_blank">holedaemon</a>. It's primarily written in Go, but uses some
            other technologies. You can find the source code on <a class="has-text-warning"
              href="https://github.com/holedaemon/bot2" target="_blank">GitHub</a>. Check out the <a
              class="has-text-warning" href="/docs">docs</a> to find out more about what the bot can do.</p>

          <p class="subtitle">If you'd like to add DILF to your server, please contact <a class="has-text-warning"
              href="https://discord.com/users/67803413605253120" target="_blank">@holedaemon</a> on Discord.</p>

          <p class="subtitle">This site is running version <code>`)
	qw422016.E().S(version.Version())
	qw422016.N().S(`</code>.</p>
        </div>
      </div>
    </div>
  </section>
`)
}

func (p *AboutPage) WriteBody(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamBody(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *AboutPage) Body() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteBody(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
