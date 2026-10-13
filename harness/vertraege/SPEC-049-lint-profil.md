# SPEC-049 — Lint-Profil (`lint`)

`make lint` prüft den Go-Code des Moduls statisch nach dem Profil in `.golangci.yml`
und mit drei eigenen Prüfungen; es schreibt nichts in den Arbeitsbaum.

1. **Gegenstand.** Geprüft werden alle Pakete des Moduls (`./...`) mit dem Build-Tag
   `integration`, also `cmd/`, `internal/` und `test/` einschließlich der
   Integrationstests und aller Testdateien. Das Profil ist die Datei `.golangci.yml`
   an der Wurzel des Repos; Pfade in ihr gelten relativ zu ihr
   (`run.relative-path-mode: cfg`) und beginnen mit `^`.
2. **Werkzeug und Umgebung.** golangci-lint `v2.14.0` aus dem Image
   `golangci/golangci-lint`, per Digest gepinnt, als Stufe `lint` des `Dockerfile`; der
   Pin steht nur dort. Die Stufe läuft auf der Plattform des Bau-Hosts, mit den Modulen
   der Stufe `deps`, ohne Netz, mit `GOTOOLCHAIN=local` und `GOFLAGS=-mod=readonly`.
   Trägt die Go-Version des Images die `go`-Zeile von `go.mod` nicht, scheitert die
   Stufe.
3. **Linter.** Aktiv sind genau diese (`linters.default: none`): `errcheck`, `govet`,
   `ineffassign`, `staticcheck`, `unused`, `containedctx`, `contextcheck`, `cyclop`,
   `dupl`, `fatcontext`, `forbidigo`, `funlen`, `gochecknoglobals`, `gochecknoinits`,
   `gocognit`, `gocyclo`, `gomodguard_v2`, `iface`, `inamedparam`, `interfacebloat`,
   `ireturn`, `maintidx`, `nestif`, `noctx`, `reassign`, `revive`, `testpackage`,
   `unparam`.
4. **Schwellen.** `cyclop` 15, `funlen` 100 Zeilen und 60 Anweisungen, `gocognit` 20,
   `gocyclo` 15, `nestif` 5, `maintidx` unter 20, `dupl` 150, `interfacebloat` 10.
5. **Einstellungen.**
   - `errcheck`: Ein ungeprüfter Fehler ist nur bei `fmt.Fprint`, `fmt.Fprintf`,
     `fmt.Fprintln`, `(net.Conn).Close` und `(net.Listener).Close` zulässig.
   - `ireturn`: Eine Funktion darf ein Interface nur liefern, wenn es `error`, leer,
     anonym, aus der Standardbibliothek oder generisch ist, unter
     `internal/hexagon/ports/` liegt oder aus dem Paket `pgproto3` stammt.
   - `forbidigo`: Verboten sind `fmt.Print`, `fmt.Printf`, `fmt.Println`, die
     eingebauten `print` und `println` sowie `os.Stdout` und `os.Stderr`; die beiden
     letzten sind unter `cmd/` und `test/` erlaubt.
   - `gomodguard_v2`: Erlaubt sind Importe aus der Standardbibliothek, dem Modul selbst
     und den Modulen `github.com/jackc/pgx/v5`, `go.yaml.in/yaml/v3` und `modernc.org/sqlite`. Jedes andere
     Modul ist ein Befund, auch eines, das `go.mod` als indirekt führt.
   - `revive`: die Regeln `blank-imports`, `context-as-argument`, `context-keys-type`,
     `dot-imports`, `empty-block`, `error-naming`, `error-return`, `error-strings`,
     `errorf`, `exported`, `if-return`, `increment-decrement`, `indent-error-flow`,
     `package-comments`, `range`, `receiver-naming`, `redefines-builtin-id`,
     `superfluous-else`, `time-naming`, `unexported-return`, `unreachable-code`,
     `unused-parameter`, `var-declaration`, `var-naming` und `unused-receiver`, auf jedem
     Pfad, auch unter `internal/`. `exported` prüft die Form des Doc-Kommentars, nicht
     seine Sprache.
   - `testpackage`: Übersprungen wird nur eine Datei namens `export_test.go`
     (`skip-regexp` `(^|/)export_test\.go$`).
   - `contextcheck`: ohne Einstellung und ohne Ausnahme, auch in Testdateien. Meldet
     es in einem Test einen neuen Kontext, kommt der Kontext aus `t.Context()` des
     laufenden Tests oder Subtests; in einer Funktion für `t.Cleanup`, die nach dessen
     Abbruch läuft, aus `context.WithoutCancel(t.Context())`. `context.Background()` an
     einer Stelle ohne Befund bleibt zulässig.
   - `gochecknoglobals`: ohne Einstellung; die Variable `version`, die Leerstelle `_`
     und Fehlerwerte mit Präfix `Err` lässt der Linter selbst zu.
6. **Kein `//nolint`.** Eine eigene Prüfung meldet jede Zeile einer Datei `*.go` unter
   `cmd/`, `internal/` und `test/`, in der auf `//` oder `/*`, gefolgt von beliebig
   vielen Leerzeichen, Tabs oder `/`, das Wort `nolint` folgt — gleich in welcher
   Schreibweise, an welcher Stelle der Zeile und mit welchem Zusatz (`:linter`,
   Begründung). Maßgeblich ist jedes `//` und jedes `/*` der Zeile, nicht nur das, mit
   dem der Kommentar beginnt: `// x //nolint` ist ein Befund. Folgt `nolint` erst nach
   einem anderen Wort auf das letzte Kommentarzeichen davor (`// siehe nolint`), ist
   das kein Befund.
7. **Export-Test-Brücke.** Unit-Tests liegen im Paket `<name>_test`. Auf Unexportiertes
   greifen sie nur über die Datei `export_test.go` im Verzeichnis des Pakets zu; sie
   gehört zum Paket `<name>` und ist dort die einzige Testdatei. Sie enthält nur
   Typ-Aliase, Konstanten und Funktionen oder Methoden, die an Unexportiertes
   weiterreichen. Zustand setzt sie nur an einem Wert, den sie übergeben bekommt, nie
   auf Paketebene. Einen Wert eines unexportierten Typs legt nur der Produkt-Code an:
   Weder die Brücke noch ein Test erzeugt ihn selbst, auch nicht als lokalen Wert, als
   Literal oder über einen Typ-Alias. Die Brücke bekommt ihn übergeben oder reicht an
   eine Funktion des Pakets weiter, die ihn erzeugt; so prüft ein Test nur Zustände,
   die das Produkt erzeugt. Diese Funktion ruft auch der Produkt-Code, um denselben
   Wert zu erzeugen; eine Funktion, die nur die Brücke ruft, gibt es nicht. Werte
   exportierter Typen, die sie annimmt, etwa eine Verbindung, stellt der Test. Was nur
   an einem Wert eines unexportierten Typs zu sehen ist, den die Brücke weder
   übergeben bekommt noch über eine solche Funktion erzeugen lässt, prüft der Test
   über die exportierte Schnittstelle. Jede andere Testdatei im Paket `<name>` ist ein
   Befund von `testpackage`; eine Variable auf Paketebene in der Brücke ist ein Befund von
   `gochecknoglobals`; eine Funktion, deren Name mit `Test`, `Benchmark`, `Example`
   oder `Fuzz` beginnt, meldet in der Brücke eine eigene Prüfung.
8. **Ausnahmen.** Eine Ausnahme steht nur in `.golangci.yml`, als Einstellung nach
   Punkt 5 oder als Regel unter `linters.exclusions.rules`. Unmittelbar über jeder
   Regel steht ein Kommentarblock, dessen erste Zeile mit `# Why:` beginnt; fehlt er,
   meldet das eine eigene Prüfung. Das Profil hat dafür eine feste Form: `linters:`
   ohne Einzug, darunter `exclusions:` mit zwei Leerzeichen, darunter `rules:` mit vier,
   jede Regel als Eintrag `- ` mit sechs, alle drei Schlüssel in Blockform (nach dem
   Doppelpunkt höchstens Leerraum und ein Kommentar). Ein Eintrag ist auch ein `-` allein
   auf der Zeile, dessen Inhalt auf den Fortsetzungszeilen steht; über ihm steht der
   Kommentarblock wie über jedem Eintrag. Ein Schlüssel `exclusions` an anderer Stelle
   oder in Flussform und ein Schlüssel `rules` unter `exclusions` mit anderem Einzug
   oder in Flussform ist ein Befund `Form nicht erkannt` mit seiner Zeile. Ebenso ist
   unter `rules` jede Zeile ein solcher Befund, die weder Kommentar noch Leerzeile ist,
   weder ein Eintrag `- ` mit sechs Leerzeichen noch dessen Fortsetzung mit mindestens
   acht unter einem solchen Eintrag; ein Eintrag `- ` mit anderem Einzug als sechs ist
   es immer. Unter `exclusions` stehen nur `warn-unused: true` und `rules`. Zulässig sind nur diese Regeln, jede dauerhaft und
   mit einem Grund, der auch für neuen Code gilt: für Testdateien `cyclop`,
   `gocognit`, `gocyclo`, `nestif`, `funlen`, `noctx`, `unparam` und `revive` mit
   `unused-parameter` und `unused-receiver`; `staticcheck` mit `ST1005`, weil die
   Fehlertexte deutsch sind und mit einem Substantiv beginnen dürfen;
   `gochecknoglobals` für eine benannte Nachschlage-Tabelle oder einen benannten
   Sentinel-Wert, die nach der Initialisierung nur gelesen werden, je Regel mit Datei
   und Namen; `forbidigo` für `os.Stdout` und `os.Stderr` nach Punkt 5. Eine Regel,
   die einen Befund nur deshalb ausblendet, weil der Bestand ihn trägt, gibt es nicht.

   Eine Regel, die im Lauf keinen Befund ausblendet (`exclusions.warn-unused: true`), ist
   ein Befund.
9. **Ausgabe und Ausgang.** golangci-lint liest das Profil nur aus `.golangci.yml`
   (`-c`), nie aus einer Default-Suche. Die Stufe prüft in dieser Reihenfolge:
   - **Profil vorhanden.** Fehlt `.golangci.yml`, schreibt sie
     `lint: .golangci.yml: fehlt`; Schema-Prüfung und golangci-lint laufen dann nicht,
     die eigenen Prüfungen nach Punkt 6 und 7 schon.
   - **Schema.** `golangci-lint config verify` prüft das Profil gegen das eingebettete
     Schema, ohne Netz. Lehnt es das Profil ab, etwa wegen eines unbekannten Schlüssels,
     schreibt die Stufe `lint: .golangci.yml: von golangci-lint abgelehnt` und dazu die
     Meldung von `config verify`; golangci-lint läuft trotzdem.
   - **Befunde.** golangci-lint gibt seine Befunde im Textformat ungekürzt aus
     (`max-issues-per-linter: 0`, `max-same-issues: 0`, `uniq-by-line: false`). Jede
     Meldung ist ein Befund, auch mehrere auf derselben Zeile.
   - **Ungenutzte Regel.** Je Warnung von `warn-unused` schreibt die Stufe eine Zeile
     `lint: .golangci.yml: Regel ohne Befund: <feld>: <wert>, …` ohne Zeilennummer. Die
     Felder stehen in dieser Reihenfolge und nur, wenn die Warnung sie nennt: `Linter`,
     `Pfad`, `Pfad außer`, `Text`, `Quelle`, je mit dem Wert aus der Warnung. Nennt sie
     keines davon, lautet die Zeile `lint: .golangci.yml: Regel ohne Befund: Felder
     nicht erkannt`; die Warnung selbst steht daneben auf der Ausgabe.

   Die eigenen Prüfungen nach Punkt 6 bis 8 schreiben je Befund eine Zeile
   `lint: <pfad>:<zeile>: <befund>`. Jede Prüfung läuft, auch wenn eine andere einen
   Befund hat. Die Stufe endet mit einem Ausgang ungleich 0 bei mindestens einer
   `lint:`-Zeile oder einem Befund von golangci-lint.
10. **Werkzeug, dann Gate.** Solange der Bestand außerhalb der Regeln nach Punkt 8
    Befunde hat, ist `make lint` ein Werkzeug ohne Gate: Es meldet alle Befunde des
    Moduls mit Pfad und endet nach Punkt 9, hängt aber nicht an der Gate-Kette von
    `make gates`. Ohne Befund im Bestand hängt es an der Gate-Kette und endet `make
    gates` mit Fehlerstatus, wenn es rot ist.

**Grenze.** Ob ein `Why:` zutrifft, ob die Brücke nur weiterreicht und ob ein Wert
eines unexportierten Typs nur aus dem Produkt-Code stammt, prüft das Werkzeug nicht. Dateien mit der Markierung für generierten Code nimmt golangci-lint nach
seinem Default aus; das Modul hat keine. Ebenso wenig prüft es, ob eine Einstellung
nach Punkt 5 ihren Grund als Kommentar trägt und ob unter `exclusions` nur die Schlüssel
nach Punkt 8 stehen. Ob eine Regel unter `rules` zu den zulässigen nach Punkt 8 gehört
und ob sie mehr als Bestand ausblendet, prüft es ebenfalls nicht: Eine neue Regel mit
`# Why:`, die einen Befund ausblendet, lässt die Stufe grün; das ist Urteil des Review. Eine Regel, die `config verify` annimmt und golangci-lint erst
beim Laden ablehnt, macht die Stufe nach Punkt 9 rot, aber ohne `lint:`-Zeile; zu
sehen ist dann nur die Meldung von golangci-lint. Ein `nolint` nach Punkt 6 in einem String-Literal ist ebenfalls ein Befund. Testdateien im Paket `main` lässt
`testpackage` zu; unter `cmd/` gibt es keine.
Ist das Profil durch einen Eintrag unter `rules` mit falschem Einzug kein gültiges YAML
mehr, sagt Punkt 8 nur den Befund `Form nicht erkannt` an der ersten solchen
Eintragszeile zu; wie die Prüfung die Zeilen danach einordnet, ist nicht zugesagt. Rot
ist die Stufe durch diesen Befund nach Punkt 9. Sind Image, Module und Inhalt des
Build-Kontexts dieselben wie bei einem früheren Lauf mit Ausgang 0, nimmt der Build das
Ergebnis der Stufe aus dem Cache, ohne neu zu prüfen; ein Lauf mit Ausgang ungleich 0
liegt nie im Cache. Eine Datei, deren Größe und Änderungszeit seit dem letzten Lauf
gleich blieben, überträgt der Build nicht neu; ihre Änderung sieht der Lauf dann nicht.
