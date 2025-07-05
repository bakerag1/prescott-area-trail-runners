<h2>PATR Events</h2>
<i>Group runs or mini-races set up by group members.</i>
{{ range where .Site.Pages "Section" "events" }}
    {{ range where (sort .Pages ".Params.Startdate" "asc" ) "Params.Patr" true}}
        {{ $start := (time.AsTime .Params.Startdate) }}
        {{ $start = time.AsTime (time.Format "2006-01-02" $start) }}
        {{ if and (ge $start.Unix (add $.Page.Date.Unix -86400)) (le $start.Unix (add $.Page.Date.Unix 5184000)) }}
            <div class="eventContainer">
                <div class="eventDate"> {{$start.Format "Jan 02"}}
                </div>
                <div class="event">
                    <summary><a href="{{.Permalink}}">{{.Title}}</a>
                        <br>
                        <details class="collapsible">
                            <p>
                                {{ .Content }}
                            </p>
                        </details>
                </div>
            </div>
        {{ end }}
    {{end}}
{{end}}
    <h2>Non-PATR Events</h2>
    <i>events not too far away, not directly PATR affiliated, but likely to be attended by PATRs</i>
    {{ range where .Site.Pages "Section" "events" }}
    {{ range where (sort .Pages ".Params.Startdate" "asc" ) "Params.Patr" false}}
        {{ $start := (time.AsTime .Params.Startdate) }}
        {{ $start = time.AsTime (time.Format "2006-01-02" $start) }}
        {{ if and (ge $start.Unix (add $.Page.Date.Unix -86400)) (le $start.Unix (add $.Page.Date.Unix 5184000)) }}
            <div class="eventContainer">
                <div class="eventDate"> {{$start.Format "Jan 02"}}
                </div>
                <div class="event">
                    <summary><a href="{{.Permalink}}">{{.Title}}</a>
                        <br>
                        <details class="collapsible">
                            <p>
                                {{ .Content }}
                            </p>
                        </details>
                </div>
            </div>
        {{ end }}
    {{end}}
{{end}}
<script>
    var coll = document.getElementsByClassName("collapsible");
    var i;

    for (i = 0; i < coll.length; i++) {
        coll[i].addEventListener("click", function () {
            this.classList.toggle("active");
            var content = this.nextElementSibling;
            if (content.style.maxHeight) {
                content.style.maxHeight = null;
            } else {
                content.style.maxHeight = content.scrollHeight + "px";
            }
        });
    }
</script>