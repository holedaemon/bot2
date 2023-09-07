// Code generated by qtc from "guild.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package templates

import "github.com/holedaemon/bot2/internal/db/models"

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

func streamisActive(qw422016 *qt422016.Writer, item string, class string) {
	qw422016.N().S(`
    `)
	if item == class {
		qw422016.N().S(`
        is-active
    `)
	}
	qw422016.N().S(`
`)
}

func writeisActive(qq422016 qtio422016.Writer, item string, class string) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamisActive(qw422016, item, class)
	qt422016.ReleaseWriter(qw422016)
}

func isActive(item string, class string) string {
	qb422016 := qt422016.AcquireByteBuffer()
	writeisActive(qb422016, item, class)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamquoteLink(qw422016 *qt422016.Writer, q *models.Quote) {
	qw422016.N().S(`
    discord://discord.com/channels/`)
	qw422016.E().S(q.GuildID)
	qw422016.N().S(`/`)
	qw422016.E().S(q.ChannelID)
	qw422016.N().S(`/`)
	qw422016.E().S(q.MessageID)
	qw422016.N().S(`
`)
}

func writequoteLink(qq422016 qtio422016.Writer, q *models.Quote) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamquoteLink(qw422016, q)
	qt422016.ReleaseWriter(qw422016)
}

func quoteLink(q *models.Quote) string {
	qb422016 := qt422016.AcquireByteBuffer()
	writequoteLink(qb422016, q)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

type GuildPage struct {
	BasePage
	Guild *models.Guild
}

func (p *GuildPage) StreamNavigation(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
<div class="container">
    <nav class="navbar is-hidden-tablet is-dark" role="navigation" aria-label="main navigation">
        <div class="navbar-brand">
            <div class="navbar-item">
                Navigation
            </div>

            <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false"
                data-target="bot2-guild-navbar">
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
            </a>
        </div>

        <div id="bot2-guild-navbar" class="navbar-menu">
            <div class="navbar-end">
                <a href="/guild/`)
	qw422016.E().S(p.Guild.GuildID)
	qw422016.N().S(`" class="navbar-item">General</a>
                <a href="/guild/`)
	qw422016.E().S(p.Guild.GuildID)
	qw422016.N().S(`/roles" class="navbar-item">Roles</a>
                <a href="/guild/`)
	qw422016.E().S(p.Guild.GuildID)
	qw422016.N().S(`/quotes" class="navbar-item">Quotes</a>
            </div>
        </div>
    </nav>
`)
}

func (p *GuildPage) WriteNavigation(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamNavigation(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildPage) Navigation() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteNavigation(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildPage) StreamSidebar(qw422016 *qt422016.Writer, item string) {
	qw422016.N().S(`
    <div class="column is-hidden-mobile is-narrow">
        <aside class="menu">
            <ul class="menu-list">
                <li><a href="/guild/`)
	qw422016.E().S(p.Guild.GuildID)
	qw422016.N().S(`" class='`)
	streamisActive(qw422016, item, "general")
	qw422016.N().S(`'>General</a></li>
                <li><a href="/guild/`)
	qw422016.E().S(p.Guild.GuildID)
	qw422016.N().S(`/roles" class='`)
	streamisActive(qw422016, item, "roles")
	qw422016.N().S(`'>Roles</a></li>
                <li><a href="/guild/`)
	qw422016.E().S(p.Guild.GuildID)
	qw422016.N().S(`/quotes" class='`)
	streamisActive(qw422016, item, "quotes")
	qw422016.N().S(`'>Quotes</a></li>
            </ul>
        </aside>
    </div>
`)
}

func (p *GuildPage) WriteSidebar(qq422016 qtio422016.Writer, item string) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamSidebar(qw422016, item)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildPage) Sidebar(item string) string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteSidebar(qb422016, item)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildPage) StreamTitle(qw422016 *qt422016.Writer) {
	qw422016.E().S(p.Guild.GuildName)
}

func (p *GuildPage) WriteTitle(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamTitle(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildPage) Title() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteTitle(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildPage) StreamBody(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	p.StreamNavigation(qw422016)
	qw422016.N().S(`
    <div class="container mt-6">
        <div class="columns is-variable is-8">
            `)
	p.StreamSidebar(qw422016, "general")
	qw422016.N().S(`
            <div class="column mx-6">
                <h1 class="title">`)
	qw422016.E().S(p.Guild.GuildName)
	qw422016.N().S(`</h1>
                <hr>
                <p class="subtitle">Quotes: `)
	qw422016.E().V(p.Guild.DoQuotes)
	qw422016.N().S(`</p>
            </div>
        </div>
    </div>
`)
}

func (p *GuildPage) WriteBody(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamBody(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildPage) Body() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteBody(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

type GuildQuotesPage struct {
	GuildPage
	Quotes models.QuoteSlice
}

func (p *GuildQuotesPage) StreamTitle(qw422016 *qt422016.Writer) {
	qw422016.N().S(`Quotes &mdash; `)
	qw422016.E().S(p.Guild.GuildName)
}

func (p *GuildQuotesPage) WriteTitle(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamTitle(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildQuotesPage) Title() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteTitle(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildQuotesPage) StreamMeta(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	streamimportBootstrapTableBulmaCSS(qw422016)
	qw422016.N().S(`
    `)
	streamimportFontAwesome(qw422016)
	qw422016.N().S(`
`)
}

func (p *GuildQuotesPage) WriteMeta(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamMeta(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildQuotesPage) Meta() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteMeta(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildQuotesPage) StreamScripts(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	streamimportJquery(qw422016)
	qw422016.N().S(`
    `)
	streamimportBootstrapTableJS(qw422016)
	qw422016.N().S(`
    <script>
        $("#bot2-quotes-table").on("post-body.bs.table", function (e) {
            $("#bot2-quotes-table").removeClass("is-hidden");
        });
    </script>
`)
}

func (p *GuildQuotesPage) WriteScripts(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamScripts(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildQuotesPage) Scripts() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteScripts(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildQuotesPage) StreamBody(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	p.StreamNavigation(qw422016)
	qw422016.N().S(`
    <div class="container mt-6">
        <div class="columns is-variable is-8">
            `)
	p.StreamSidebar(qw422016, "quotes")
	qw422016.N().S(`
            <div class="column mx-6">
                <h1 class="title">Quotes &mdash; `)
	qw422016.E().S(p.Guild.GuildName)
	qw422016.N().S(`</h1>
                <hr>
                `)
	if len(p.Quotes) == 0 {
		qw422016.N().S(`
                    <h1 class="title">No quotes!!!</h1>
                `)
	} else {
		qw422016.N().S(`
                    <table id="bot2-quotes-table" class="table is-hidden" data-toggle="table" data-pagination="true" data-search="true">
                        <thead>
                            <tr>
                                <th data-sortable="true">#</th>
                                <th data-sortable="true" data-width="900">Quote</th>
                                <th data-sortable="true">User</th>
                                <th data-sortable="true">Link</th>
                            </tr>
                        </thead>
                        <tbody>
                            `)
		for _, q := range p.Quotes {
			qw422016.N().S(`
                                <tr>
                                    <td>`)
			qw422016.N().D(q.Num)
			qw422016.N().S(`</td>
                                    <td class="bot2-break">`)
			qw422016.E().S(q.Quote)
			qw422016.N().S(`</td>
                                    <td>`)
			qw422016.E().S(q.QuotedUsername)
			qw422016.N().S(`</td>
                                    <td><a class="button is-success" href="`)
			streamquoteLink(qw422016, q)
			qw422016.N().S(`">Link</a></td>
                                </tr>
                            `)
		}
		qw422016.N().S(`
                        </tbody>
                    </table>
                `)
	}
	qw422016.N().S(`
            </div>
        </div>
    </div>
`)
}

func (p *GuildQuotesPage) WriteBody(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamBody(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildQuotesPage) Body() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteBody(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

type GuildRolesPage struct {
	GuildPage
	Roles models.RoleSlice
}

func (p *GuildRolesPage) StreamTitle(qw422016 *qt422016.Writer) {
	qw422016.N().S(`Roles &mdash; `)
	qw422016.E().S(p.Guild.GuildName)
}

func (p *GuildRolesPage) WriteTitle(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamTitle(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildRolesPage) Title() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteTitle(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildRolesPage) StreamMeta(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	streamimportBootstrapTableBulmaCSS(qw422016)
	qw422016.N().S(`
    `)
	streamimportFontAwesome(qw422016)
	qw422016.N().S(`
`)
}

func (p *GuildRolesPage) WriteMeta(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamMeta(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildRolesPage) Meta() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteMeta(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildRolesPage) StreamScripts(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	streamimportJquery(qw422016)
	qw422016.N().S(`
    `)
	streamimportBootstrapTableJS(qw422016)
	qw422016.N().S(`
    <script>
        $("#bot2-roles-table").on("post-body.bs.table", function (e) {
            $("#bot2-roles-table").removeClass("is-hidden");
        });
    </script>
`)
}

func (p *GuildRolesPage) WriteScripts(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamScripts(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildRolesPage) Scripts() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteScripts(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func (p *GuildRolesPage) StreamBody(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
    `)
	p.StreamNavigation(qw422016)
	qw422016.N().S(`
    <div class="container mt-6">
        <div class="columns is-variable is-8">
            `)
	p.StreamSidebar(qw422016, "roles")
	qw422016.N().S(`
            <div class="column mx-6">
                <h1 class="title">Roles &mdash; `)
	qw422016.E().S(p.Guild.GuildName)
	qw422016.N().S(`</h1>
                <hr>
                `)
	if len(p.Roles) == 0 {
		qw422016.N().S(`
                    <h1 class="title">No roles!!!</h1>
                `)
	} else {
		qw422016.N().S(`
                    <table id="bot2-roles-table" data-toggle="table" data-pagination="true" data-search="true">
                        <thead>
                            <tr>
                               <th data-sortable="true">Name</th>
                            </tr>
                        </thead>
                        <tbody>
                            `)
		for _, r := range p.Roles {
			qw422016.N().S(`
                                <tr>
                                    `)
			if r.RoleName.IsZero() || r.RoleName.String == "" {
				qw422016.N().S(`
                                        <td>N/A</td>
                                    `)
			} else {
				qw422016.N().S(`
                                        <td>`)
				qw422016.E().S(r.RoleName.String)
				qw422016.N().S(`</td>
                                    `)
			}
			qw422016.N().S(`
                                </tr>
                            `)
		}
		qw422016.N().S(`
                        </tbody>
                    </table>
                `)
	}
	qw422016.N().S(`
            </div>
        </div>
    </div>
`)
}

func (p *GuildRolesPage) WriteBody(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamBody(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *GuildRolesPage) Body() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteBody(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
