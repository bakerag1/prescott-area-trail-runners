{{ range (where .Site.RegularPages.ByDate.Reverse "Section" "post") -}}
  {{ if ( in .Params.tags "race-report" ) }}
    <article class="post">
      <h1 class="post-title">
        <a href="{{ .Permalink }}">{{ .Title }}</a>
      </h1>
      <time datetime="{{ .Date.Format "2006-01-02T15:04:05Z0700" }}" class="post-date">{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
      {{ .Summary }}
      {{ if .Truncated }}
      <div class="read-more-link">
        <a href="{{ .RelPermalink }}">Read More…</a>
      </div>
      {{ end }}
    </article>
  {{ end }}
{{- end }}