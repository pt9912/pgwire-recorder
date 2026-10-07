# `make lint` — prüft den Go-Code nach dem Lint-Profil `.golangci.yml` und den eigenen Prüfungen

---

## Vertrag

Rot heißt: Der Go-Code unter `cmd/`, `internal/` und `test/` verletzt
[`SPEC-049`](../../spec/spezifikation.md#spec-049--lint-profil-lint): golangci-lint
meldet einen Befund nach dem Profil (Linter, Schwellen und Einstellungen, Punkte 1 bis
5), eine Zeile trägt eine Direktive `nolint` (Punkt 6), die Brücke `export_test.go`
trägt eine Testfunktion (Punkt 7), oder das Profil selbst ist nicht in Ordnung: Es
fehlt, das Schema lehnt es ab, eine Regel unter `exclusions.rules` hat keinen
Kommentarblock `# Why:`, blendet nichts aus oder steht nicht in der festen Form
(Punkte 8 und 9). Als Gate hängt `make lint` an der Gate-Kette von `make gates`
(Punkt 10). Was als Befund, als Ausnahme und als feste Form gilt, sagt der Vertrag;
diese Datei entscheidet nichts davon.

## Grenze — was das Grün nicht abdeckt

Die Grenze von `SPEC-049`; jede ist ein akzeptiertes Negativ der Bindung und
permanent, solange sie gilt.

1. **Inhalt der Ausnahmen** — ob ein `Why:` zutrifft, ob eine Einstellung nach Punkt 5
   ihren Grund als Kommentar trägt und ob unter `exclusions` nur `warn-unused` und
   `rules` stehen, prüft der Lauf nicht; das bleibt Urteil des Reviews. Ob eine
   Regel unter `rules` zu den zulässigen nach Punkt 8 gehört und ob sie mehr als
   Bestand ausblendet, prüft es ebenfalls nicht: Eine neue Regel mit `# Why:`, die
   einen Befund ausblendet, lässt die Stufe grün; das ist Urteil des Review.
2. **Inhalt der Brücke** — ob `export_test.go` nur weiterreicht und ob ein Wert eines
   unexportierten Typs nur aus dem Produkt-Code stammt, prüft der Lauf nicht.
3. **Ausgenommener Code** — Dateien mit der Markierung für generierten Code nimmt
   golangci-lint nach seinem Default aus, Testdateien im Paket `main` lässt
   `testpackage` zu.
4. **Ergebnis aus dem Cache** — sind Image, Module und Inhalt des Build-Kontexts
   dieselben wie bei einem früheren grünen Lauf, nimmt der Build das Ergebnis aus dem
   Cache, ohne neu zu prüfen. Eine Datei, deren Größe und Änderungszeit gleich
   blieben, überträgt der Build nicht neu; eine Änderung mit `cp -p` oder `touch -r`
   sieht der Lauf dann nicht. Heilbar je Lauf durch ein `touch` auf die Datei.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`find cmd internal test -name '*.go' | wc -l`.

Der Lauf trägt keine Vollständigkeits-Zeile: `0 issues.` am Ende der Ausgabe von
golangci-lint sagt, dass golangci-lint nichts meldet, nicht wie viele Dateien es
geprüft hat, und nichts über die `lint:`-Zeilen.

## Ausgabe und Ausgänge

Die Ausgabe ist das Log der Stufe `lint` im Build (`--progress=plain`):

- je Befund einer eigenen Prüfung eine Zeile `lint: <pfad>:<zeile>: <befund>`, für
  das Profil `lint: .golangci.yml: <befund>` ohne Zeilennummer (fehlt, vom Schema
  abgelehnt, `Regel ohne Befund: …`), bei Ablehnung durch das Schema dazu die Meldung
  von `config verify`;
- die Befunde von golangci-lint im Textformat, `<pfad>:<zeile>:<spalte>: <text>
  (<linter>)`, ungekürzt.

| Exit | Bedeutung |
|---|---|
| 0 | kein Befund; keine `lint:`-Zeile |
| ≠ 0 | die Stufe endete mit Ausgang 1 (mindestens eine `lint:`-Zeile oder ein Befund von golangci-lint), oder der Build scheiterte vor ihr |

`make` meldet jeden Ausgang ungleich 0 als Fehler. Auseinander hält die beiden
Ursachen das Build-Log: Hat die Stufe geprüft, endet es mit `process "/bin/sh -c bash
tools/harness/lint.sh" did not complete successfully: exit code: 1`; scheiterte der
Build davor, etwa beim Laden der Module in der Stufe `deps` ohne Netz, nennt es den
Schritt, an dem er scheiterte. Lehnt golangci-lint eine Regel erst beim Laden ab
(`can't load config`), ist die Stufe rot ohne `lint:`-Zeile; zu sehen ist nur die
Meldung von golangci-lint. Ist das Profil durch einen Eintrag mit falschem Einzug kein
gültiges YAML mehr, ist nur die erste Zeile `Form nicht erkannt` zugesagt.

Ist `make lint` am Bestand rot, wird auch `make lint-gegenprobe` rot: Ihr erster Fall
`p0-grundlauf` nennt die Zeile des Bestands, die übrigen roten Fälle sind Folgen
davon. Gelesen wird deshalb zuerst `make lint`.

Aus dem Rot führt, den Befund zu beheben; eine Ausnahme ist eine Regel in
`.golangci.yml` mit `# Why:` und einem Grund, der auch für neuen Code gilt (Punkt 8),
nie eine Direktive im Code (`AGENTS.md` §3.2).

## Sperren

Keine: Die Stufe hat keine benannte Abbruch-Meldung. Jede Prüfung läuft, auch wenn
eine andere einen Befund hat; fehlt das Profil, laufen nur Schema-Prüfung und
golangci-lint nicht.

## Bindung

[ADR-0034](../../docs/plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) · Vertrag
[`SPEC-049`](../../spec/spezifikation.md#spec-049--lint-profil-lint) · Image
`golangci/golangci-lint` per Digest im `Dockerfile` (Stufe `lint`)

An der Gate-Kette von `make gates` hängt neben `make lint` auch seine Gegenprobe
`make lint-gegenprobe`; sie prüft, ob das Werkzeug richtig ist, und ihr Vertrag steht
in ihrer Zeile in [`../README.md` §Sensors](../README.md#sensors-feedback-gates).
