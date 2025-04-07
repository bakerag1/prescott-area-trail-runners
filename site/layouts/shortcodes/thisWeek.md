{{ $postdate := (time.AsTime .Page.Date)}}
{{ range (where .Site.RegularPages.ByDate "Section" "events") -}}
  {{ $start := (time.AsTime .Params.startdate) }}
  {{ if and (lt $start ($postdate.AddDate 0 0 7)) (gt $start ($postdate.AddDate 0 0 -1)) }}
  <li style="padding:5;list-style-type:none">
    <span>
      <time class="pull-right post-list" datetime="2006-01-02T15:04:05Z">{{ $start.Format "Monday, Jan 2 3:04pm"}}</time> 
      <a href="{{ .Permalink }}">{{ .Title }}</a> 
    </span>
  </li>
  {{ end }}
{{- end }}