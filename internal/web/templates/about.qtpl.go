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
    <div class="container m-auto grid grid-cols-1 xl:grid-cols-3 place-content-center">
        <div class="xl:col-start-2 xl:col-end-2 bg-slate-50 dark:bg-stone-900 text-gray-900 dark:text-white sm:rounded px-5 py-5">
            <h1 class="text-gray-900 dark:text-white text-3xl">DILF</h1>
            <p class="text-gray-900 dark:text-white mt-3">...is a Discord bot created and maintained by <a class="text-blue-300 hover:text-blue-400 dark:text-slate-300 dark:hover:text-slate-400" href="https://holedaemon.net" target="_blank">holedaemon</a>. It's primarily written in Go, but uses some other technologies. You can find the source code on <a class="text-blue-300 hover:text-blue-400 dark:text-slate-300 dark:hover:text-slate-400" href="https://github.com/holedaemon/bot2" target="_blank">GitHub</a>. Check out the <a class="text-blue-300 hover:text-blue-400 dark:text-slate-300 dark:hover:text-slate-400" href="/docs">docs</a> to find out more about what the bot can do.</p>

            <p class="text-gray-900 dark:text-white mt-3">
            If you'd like to add DILF to your server, please contact <a class="text-blue-300 hover:text-blue-400 dark:text-slate-300 dark:hover:text-slate-400" href="https://discord.com/users/67803413605253120" target=_blank">@holedaemon</a> on Discord.
            </p>

            <p class="text-gray-900 dark:text-white mt-3">
            This site is running version <code class="text-blue-500 dark:text-green-500">`)
	qw422016.E().S(version.Version())
	qw422016.N().S(`</code>.
            </p>
        </div>
    </div>
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