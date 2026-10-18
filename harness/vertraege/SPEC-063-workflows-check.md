# SPEC-063 — Gate der Workflows (`workflows-check`)

`make workflows-check` prüft jede Workflow-Datei unter `.github/workflows/`, bevor sich ein Mensch oder `make gates`
auf sie verlässt, in drei Gruppen: (a) Syntax und Ausdrücke durch `actionlint` in einem Image, das per Digest gepinnt
ist; (b) Rechte und Zeitgrenzen durch das Skript `tools/harness/workflows-check.sh`: oben `permissions: {}`, in jedem
Job eine Zeitgrenze; (c) die Bindung der Pin-Prüfung: dass `make docs-check` die Pins der Dateien prüft. Die Pins
(voller Commit-SHA mit Tag-Kommentar je fremde `uses:`-Referenz) prüft **allein** `make docs-check` über das Modul
`workflows` von d-check; dieses Gate prüft sie kein zweites Mal, sondern prüft, dass die Konfiguration dafür
steht (Punkt 8). Das Ziel ist ein Gate: es steht in `GATE_CHECKS` und nicht in der Ausnahmeliste der schnellen
Kette ([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 4); es braucht Docker wie `docs-check`. Seine Gegenprobe
`make workflows-check-gegenprobe` steht in `GEGENPROBEN_CONTAINER` ([`SPEC-060`](SPEC-060-gegenproben.md) Punkte 1 bis 3).
Warum die CI auf GitHub Actions läuft und die Dateien ein Gate bestehen, steht in
[ADR-0050](../../docs/plan/adr/0050-ci-auf-github-actions.md); diese Kennung wiederholt es nicht.

**Stand der Entscheidung.** Der Nutzer hat [ADR-0050](../../docs/plan/adr/0050-ci-auf-github-actions.md) am 2026-10-18 angenommen; ändert er diese Kennung, ändert sich nur sie mit
Fragment und Skripten. Auftrag des Nutzers vom 2026-10-18 („Erst das Release voranbringen“): die CI vor die Veröffentlichung.

**Schnittstelle.** Das Fragment `harness/mk/workflows.mk` definiert `workflows-check` und `workflows-check-gegenprobe`
(beide `.PHONY`). Das Rezept von `workflows-check` ist `@bash tools/harness/workflows-check.sh`, ohne `$(MAKE)` und ohne
`+`; die Variable `DOCKER` kommt aus `harness/mk/build.mk`, die Variable `WORKFLOWS_ACTIONLINT` (Image mit Digest,
Herkunft Docker Hub, `rhysd/actionlint@sha256:<Index-Digest>`, Tag als Kommentar daneben; `ghcr.io/rhysd/actionlint` ist anonym nicht ziehbar, gemessen am 2026-10-18) steht im Fragment und wird dem Skript als Umgebung
übergeben. Version und Digest stehen nicht in dieser Kennung (sie altern); der Implementer ermittelt sie mit Netz beim
Anlegen des Fragments (Punkt 3, Element Digest). Der Digest ist der des **Index** (Multi-Arch-Manifest), nicht der einer Plattform: er gilt
auf `amd64` wie auf `arm64`, die Attrappe der Gegenprobe und der Lauf der CI lesen dieselbe Zeichenkette, und Docker wählt die Plattform des Servers.
Das Limit anonymer Pulls des Docker Hub gilt für dieses Image wie für die übrigen Images der Gates ([ADR-0050](../../docs/plan/adr/0050-ci-auf-github-actions.md)). Das Skript nimmt als erstes Argument eine Wurzel; ohne Argument ist die
Wurzel das Verzeichnis zwei Ebenen über dem Skript (die Gegenprobe ruft es mit Wurzeln von Temp-Bäumen).

**Plattform** (`AGENTS.md` §3.12, Plattform von [`LH-QA-03`](../../spec/lastenheft.md#lh-qa-03--portabilität)). Zugesagt ist
ein Linux-Host oder -Container mit bash ≥ 4, GNU `make`, GNU-Werkzeugen (`awk`, `find`, `sort`, `cksum`) und Docker mit
Linux-Containern, `linux/amd64`. Das Image von `actionlint` liegt für `linux/amd64` und `linux/arm64` vor; getestet ist das
erste. macOS und Windows als Host sind nicht zugesagt (akzeptiertes Negativ; Grund: kein solcher Host steht zur Messung bereit,
und die Gates dieses Repos setzen Linux-Semantik voraus).

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes, den er prüft. Je Zusage und je Element einer
Aufzählung gibt es einen Fall und eine Mutation des Skripts, die nur die Fälle ihrer Gruppe rot färbt (`AGENTS.md` §3.10);
die unveränderte Kopie ist grün. Wo „jede andere Form“ steht, sind die genannten Formen Fälle derselben Zusage und teilen
eine Mutation.

**Randform → Punkt** (Randformen aus §6 des Plans):

| Randform | Punkt |
|---|---|
| 1 Kein Verzeichnis, leeres Verzeichnis, keine Datei | 2 |
| 2 Endungen und Ort | 1 |
| 3 Weitere Orte (Actions, Dependabot) | 11 |
| 4 Formen von `permissions` | 6 |
| 5 Einstufung der Befunde | 4 |
| 6 `timeout-minutes` | 7 |
| 7 Reihenfolge der Gründe | 9 |
| 8 Aufrufarten des Treibers | 10 |
| 9 Unterverzeichnis, Quelle read-only, Symlinks | 1, 3, 10 |
| 10 Parallelität | 3, 10 |
| 11 Version, Digest, Herkunft, fehlendes Image | 3 |
| 12 Zusammenspiel mit `docs-check` | 8 |
| Gültigkeit und Aktualität der Pins, Lauffähigkeit | 11 |

1. **Prüfbereich.** Geprüft werden die regulären Dateien unmittelbar in `<Wurzel>/.github/workflows/` mit der Endung
   `.yml` (E1) oder `.yaml` (E2), klein geschrieben und als Ende des Namens, in der Reihenfolge der Namen (Locale `C`).
   Eine Datei in einem Unterverzeichnis und eine Datei mit anderer Endung (`.YML`, `.yml.bak`, `README.md`) werden nicht geprüft;
   jede steht auf stdout als Zeile `hinweis: <Pfad> wird nicht geprüft (<Grund>)` und ist kein Befund. Ein Symlink mit **anderer** Endung ist wie jede Datei mit anderer Endung ein Hinweis `(Endung)`, kein Befund und nicht gefolgt. Ein Symlink mit einer der
   beiden Endungen ist ein Befund `symlink` (Zeile 1), wird nicht gefolgt und nicht an `actionlint` gegeben, auch bei einem Ziel
   innerhalb des Baums. Fälle: `p1-yml`, `p1-yaml` (je eine Datei mit Befund aus Punkt 6, rot, Datei genannt),
   `p1-andere-endung` (Dateien `x.YML`, `x.yml.bak` und `x.txt` mit Befund: grün, drei Hinweise), `p1-unterverzeichnis`,
   `p1-symlink`, `p1-symlink-andere-endung` (`x.txt` als Symlink: grün, Hinweis mit `(Endung)`, kein Befund). Mutationen: Prüfbereich nur `.yml`; nur `.yaml`; Endung als Teilstring (`*.yml*`); rekursive Suche; Symlink gefolgt; Symlink mit anderer Endung wird Befund oder ohne Hinweis übergangen (je eine, färben `p1-symlink-andere-endung`).
2. **Kein Prüfbereich.** Fehlt das Verzeichnis (F1), ist es leer (F2) oder enthält es nur Dateien, die Punkt 1 weder prüft noch als Symlink beanstandet (F3), meldet das
   Skript auf stdout `workflows-check: keine Workflow-Dateien unter .github/workflows` und endet mit Exit 0, **ohne** `docker` zu rufen
   und ohne die Prüfungen der Punkte 4 und 8. Das ist grün über leerem Prüfbereich und gilt bis zur ersten Datei (Grenze, Punkt 11); die
   Meldung ist deshalb Pflicht. Fälle: `p2-kein-verzeichnis`, `p2-leer`, `p2-nur-andere` (Exit 0, Meldung, Attrappe von `docker` nicht
   gerufen). Mutationen je Form: F1 endet mit Exit 2; F2 endet still; F3 ruft `docker`. **Symlinks gehen vor:** Besteht das Verzeichnis nur aus Symlinks mit einer Endung
   von Punkt 1 (F4), gibt es keinen Kandidaten, aber einen Befund: Der Befund `symlink` gilt (Punkt 1), nicht die Meldung „keine Workflow-Dateien“; das Skript endet mit Exit 1 und der Schlusszeile
   von Punkt 9, ruft `docker` nicht (keine Dateiliste) und prüft weder Punkt 4 noch Punkt 8 (beide brauchen einen Kandidaten bzw. `actionlint`). Ein Symlink neben einem regulären Kandidaten
   ändert nichts am Lauf der übrigen Punkte; sein Befund kommt hinzu (Punkt 9). Fall `p2-nur-symlink` (Exit 1, Befundzeile `symlink`, Summe 1, Attrappe nicht gerufen, die Meldung
   „keine Workflow-Dateien“ fehlt); Mutationen: ein Verzeichnis nur mit Symlinks endet mit Exit 0 und der Meldung von F3; das Skript ruft bei F4 `docker`.
3. **Aufruf von `actionlint`.** Das Skript ruft genau einmal `"$DOCKER" run --rm --network none -v "<Wurzel>:/repo:ro" -w /repo
   "$WORKFLOWS_ACTIONLINT" -oneline <Dateien>`; `<Dateien>` sind die Kandidaten aus Punkt 1 ohne Symlinks als Pfade
   `.github/workflows/<Name>` in der Reihenfolge der Namen. Elemente: `--rm` (E1), `--network none` (E2), Mount mit `:ro` (E3), das Image
   aus `WORKFLOWS_ACTIONLINT` (E4), `-oneline` (E5), die Dateiliste (E6), kein `--name` (E7; Punkt 10), kein `-ignore`, `-shellcheck`
   und `-pyflakes` (E8; Punkt 4). Fälle `p3-e1` bis `p3-e8` an einer Attrappe von `docker`, die ihre Argumente aufzeichnet; Mutation je Element:
   das Element fehlt (E1 bis E6) bzw. kommt hinzu (E7, E8). **Exit-Code von `actionlint`:** 0 heißt keine Befunde dieser Gruppe; 1 heißt Befunde,
   deren Zeilen (stdout der Attrappe) unverändert in die Ausgabe kommen (Punkt 9); jeder andere Wert (2, 3, 125 beim Pull-Fehler, 127), ein
   fehlendes `docker` und ein Pull-Fehler enden mit Exit 2 und der Meldung `workflows-check: actionlint nicht lauffähig (Exit <n>)` samt
   stderr des Aufrufs; ein Fehler ist nie „keine Befunde“. Das Image läuft als Nicht-root: Exit 3 mit `permission denied` im stderr heißt, dass Wurzel oder `.github/workflows` für andere nicht lesbar sind (etwa ein Verzeichnis mit Modus 700); das Skript hängt dann die Zeile `workflows-check: Hinweis: der Container liest als Nicht-root; Wurzel und .github/workflows brauchen Leserechte für alle` an. Das Gate kopiert den Baum nicht und ändert keine Rechte (Punkt 10, es schreibt nichts); ein Baum mit Modus 700 ist Exit 2, nicht grün. Die Gegenprobe legt ihre Temp-Bäume für die echten Fälle mit `chmod -R a+rX` lesbar an (Punkt 12). Fall `p3-exit3-rechte` (Attrappe: Exit 3, stderr `permission denied`: Exit 2 und Hinweiszeile); Mutation: der Hinweis fehlt. Fälle `p3-exit0`, `p3-exit1`, `p3-exit2`, `p3-exit3` (fehlende Datei), `p3-exit125`, `p3-kein-docker`; Mutationen: Exit 1
   gilt als 0; jeder Exit ungleich 0 gilt als 1; ein fehlendes `docker` gilt als keine Befunde. **Digest** (E4): das Fragment trägt für
   `WORKFLOWS_ACTIONLINT` die Form `<Repository>@sha256:<64 Hexziffern>` mit `<Repository>` = `rhysd/actionlint`; Fall `p3-digest-form` an Kopien des Fragments; Mutation: ein Tag
   statt des Digest. Das Skript bricht ab mit Exit 2, wenn `WORKFLOWS_ACTIONLINT` leer ist (Fall `p3-image-leer`; Mutation: leeres Image
   wird gerufen). Die Gegenprobe fährt daneben drei echte Fälle (Punkt 12). Gemessen am echten Image (2026-10-18): `-oneline`, `--network none` und `:ro` gehen, Exit 1 bei Befunden, Zeilenform `<Pfad>:<Zeile>:<Spalte>: …`, Exit 3 bei fehlender Datei.
4. **Einstufung.** Jede Zeile von `actionlint` ist ein Befund, auch die der eingebauten Prüfungen von `shellcheck` und `pyflakes`; es gibt keine
   Ausnahmeliste und keine Option dafür. Existiert `.github/actionlint.yaml` (E1) oder `.github/actionlint.yml` (E2), ist das ein Befund
   `konfiguration` (Zeile 1 der Datei): `actionlint` läse sie und könnte Befunde ausblenden. Eine Ausnahme ist eine Änderung dieser Kennung, kein
   Eintrag in einer Datei. Fälle `p4-yaml`, `p4-yml` (Datei vorhanden, ein Workflow ohne sonstigen Befund: rot); Mutationen je Endung: nicht
   geprüft. Der Fall `p3-e8` hält, dass das Skript keine Option zum Ausblenden übergibt.
5. **Lesegrammatik der Regeln aus 6 und 7.** Die Regeln lesen jede Datei zeilenweise in der Blockform, die dieses Repo schreibt, und raten nicht.
   Leerzeilen und Zeilen, deren erstes Zeichen außer Leerzeichen `#` ist, sind ohne Bedeutung. Ein Schlüssel der obersten Ebene beginnt in Spalte 1;
   `jobs:` steht allein oder mit Kommentar; die Schlüssel der Jobs stehen mit Einzug von genau zwei Leerzeichen, deren Eigenschaften (`runs-on:`,
   `uses:`, `timeout-minutes:`, `permissions:`) mit genau vier. Ein Tabulator im Einzug (E1), ein Job-Schlüssel mit anderem Einzug als zwei (E2) und
   `jobs:` mit einem Wert in der Zeile, etwa als Flow-Mapping (E3), sind Befunde `form`; für diese Datei laufen die Regeln der Punkte 6 und 7 nicht
   (ihre Befunde wären Folgefehler). Fälle `p5-tab`, `p5-einzug`, `p5-jobs-flow`; Mutation je Element: die Form wird still angenommen und die
   Regeln laufen. Eine Eigenschaft in einem Schritt (Einzug über vier) ist keine Eigenschaft des Jobs: Fall `p5-schritt-ist-kein-job`
   (`timeout-minutes:` nur in einem Schritt, Job ohne eigene: Befund `timeout-fehlt`); Mutation: Suche auf jedem Einzug.
6. **Rechte.** Oben (Spalte 1) steht genau eine Zeile `permissions:` mit dem Wert `{}`, danach höchstens Leerraum und ein Kommentar. Fehlt sie, ist
   das ein Befund `permissions-fehlt`; jede andere Form (`read-all`, `write-all`, ein Block mit Einträgen, `{contents: read}`) ist ein Befund
   `permissions-form` mit dem gelesenen Wert im Text. In einem Job (Einzug vier) ist `permissions:` zulässig mit dem Wert `{}` oder als Block,
   dessen Einträge `<Scope>: read`, `<Scope>: write` oder `<Scope>: none` lauten; ein Skalar (`read-all`, `write-all`, ein anderer Wert) oder
   ein Flow-Mapping mit Einträgen ist ein Befund `job-permissions-form`, ein anderer Wert im Block ein Befund `job-permissions-wert`.
   Fehlt `permissions:` am Job, ist das kein Befund. Das Gate urteilt nicht über die Zahl der Scopes. Fälle: `p6-fehlt`, `p6-leer-ok`,
   `p6-leer-kommentar-ok` (Kommentar hinter `{}` und Kommentarzeile davor), `p6-read-all`, `p6-write-all`, `p6-block`, `p6-flow-eintrag`;
   `p6-job-block-ok`, `p6-job-leer-ok`, `p6-job-ohne-ok`, `p6-job-write-all`, `p6-job-read-all`, `p6-job-flow-eintrag`, `p6-job-wert`
   (`contents: admin`). Mutationen: oben fehlendes `permissions:` angenommen; jeder Wert oben angenommen (färbt `p6-read-all`, `p6-write-all`,
   `p6-block`, `p6-flow-eintrag`); ein Skalar am Job angenommen (färbt `p6-job-write-all`, `p6-job-read-all`); ein Flow-Mapping am Job angenommen;
   ein unbekannter Wert im Block angenommen; `permissions:` am Job verlangt.
7. **Zeitgrenzen.** Jeder Job mit der Eigenschaft `runs-on:` (Einzug vier) trägt `timeout-minutes:` mit einer dezimalen Ganzzahl in Normalform (Ziffern ohne Vorzeichen und ohne führende Null) von 1 bis 60, danach
   höchstens Leerraum und ein Kommentar. Fehlt sie: Befund `timeout-fehlt`; jeder andere Wert (0, 61, `010`, `0060`, `00`, `+5`, `5.0`, ein Ausdruck `${{ … }}`, `30.5`, ein
   Wort): Befund `timeout-wert` mit dem gelesenen Wert. Die Obergrenze 60 liegt unter dem Standard von GitHub (360); sie zu heben ist eine
   Änderung dieser Kennung. Ein Job ohne `runs-on:` (der Aufruf eines wiederverwendbaren Workflows mit `uses:`) wird nicht geprüft, weil GitHub
   dort keine `timeout-minutes` zulässt; `actionlint` fängt sie dort (Punkt 3). Fälle: `p7-fehlt`, `p7-1-ok`, `p7-60-ok`, `p7-61`, `p7-0`,
   `p7-ausdruck`, `p7-nicht-ganzzahl`, `p7-fuehrende-null` (`010`, `0060`, `00`: je rot, `timeout-wert` mit dem Wert), `p7-vorzeichen` (`+5`), `p7-dezimal` (`5.0`), `p7-uses-job-ok`. Mutationen: Pflicht entfällt; Obergrenze 61; Obergrenze 360; 0 angenommen; Ausdruck
   angenommen; Nicht-Ganzzahl angenommen; führende Null angenommen (färbt `p7-fuehrende-null`); Vorzeichen angenommen (färbt `p7-vorzeichen`); Dezimalform angenommen (färbt `p7-dezimal`); Job ohne `runs-on:` geprüft.
8. **Bindung der Pin-Prüfung.** Gibt es mindestens eine Datei nach Punkt 1, muss `<Wurzel>/.d-check.yml` (a) in der Zeile `modules:` in der Flow-Form
   `[…]` das Element `workflows` führen, (b) eine Zeile `workflows:` in Spalte 1 führen und (c) darunter, vor dem nächsten Schlüssel in Spalte 1,
   eine Zeile `  dir:` mit dem Wert `.github/workflows` (mit oder ohne Anführungszeichen, ohne schließenden Schrägstrich). Sonst ist es ein Befund
   `pins-ungeprueft` (Zeile 1 der Workflow-Datei, die zuerst sortiert) mit dem fehlenden Teil im Text; fehlt `.d-check.yml`, ebenso. Das Gate
   urteilt nicht über die Pins selbst (Randform 12): eine Datei mit `uses: actions/checkout@v4` und gebundener Konfiguration ist für dieses Gate
   grün; `make docs-check` meldet sie (Punkt 12, Fall `real-pins`). Fälle: `p8-ok`, `p8-modules` (Element fehlt), `p8-block` (Zeile fehlt), `p8-dir`
   (anderer Wert), `p8-ohne-datei`, `p8-pins-nicht-hier` (unpinnte Referenz, Konfiguration gebunden: grün). Mutationen: (a) nicht geprüft; (b) nicht
   geprüft; (c) nicht geprüft; fehlende Datei angenommen; das Skript meldet unpinnte Referenzen selbst.
9. **Ausgabe, Exit-Codes, Reihenfolge der Gründe.** Reihenfolge: (1) kein Prüfbereich und kein Symlink: Exit 0 (Punkt 2); nur Symlinks: Exit 1 ohne Aufruf von `docker` (Punkt 2, F4); (2) `docker` oder `actionlint` nicht lauffähig:
   Exit 2 sofort, **ohne** die Zeilen der Gruppen (b) und (c), weil eine unvollständige Liste Befunde vortäuscht (Punkt 3); (3) sonst laufen alle Gruppen in
   einem Lauf und nennen alle Befunde: zuerst die Zeilen von `actionlint` in ihrer Reihenfolge, dann je Datei in der Reihenfolge der Namen die
   eigenen Befunde nach Zeile, zuletzt die Befunde ohne einzelne Datei (`konfiguration`, `pins-ungeprueft`). Eine Befundzeile des Skripts lautet
   `<Pfad>:<Zeile>: <Regel>: <Text>` auf stderr (Zeile 1 für einen Befund der ganzen Datei); die Zeilen von `actionlint` stehen unverändert ebenfalls auf
   stderr. Bei mindestens einem Befund folgt als letzte Zeile `workflows-check: <n> Befund(e)` (Summe beider Quellen) und der Exit ist 1. Ohne
   Befund steht auf stdout `workflows-check: <m> Datei(en) geprüft, kein Befund` (m = Zahl der Kandidaten ohne Symlinks) und der Exit ist 0. Exit 2
   steht nur für Punkt 3, eine nicht lesbare Datei und ein nicht lesbares Verzeichnis. Fälle: `p9-reihenfolge` (je ein Befund aus `actionlint`, aus
   Punkt 6 in der zweiten Datei, aus Punkt 7 in der ersten und aus Punkt 8: Ausgabe in der genannten Reihenfolge), `p9-alle-befunde` (eine Datei mit drei
   Befunden: alle drei genannt, Summe 3), `p9-umgebung-vor-befund` (`docker` fehlt und ein Befund aus Punkt 6: Exit 2, keine Zeile aus Punkt 6),
   `p9-gruen` (Meldung mit m), `p9-nicht-lesbar` (Datei ohne Leserecht: Exit 2; als root nicht herstellbar, dann laut übersprungen). Mutationen: die
   Gruppen laufen in anderer Reihenfolge; das Skript endet beim ersten Befund; die Summe zählt nur eine Quelle; Exit 1 bei Umgebungsfehler; die Gruppen
   (b) und (c) laufen vor dem Umgebungsfehler.
10. **Aufrufarten, Unterverzeichnis, Schreiben, Parallelität.** Das Rezept ruft das Skript genau einmal und trägt weder `$(MAKE)` noch `+`; darum
    führt `make` es unter `-n`, `-q` und `-t` nicht aus. Fälle an einer Kopie von Makefile und Fragment mit einer Attrappe von `docker`
    (Aufzeichnung der Aufrufe): `p10-n` (Exit 0, Rezept gedruckt, kein Aufruf), `p10-q` (Exit 1, kein Aufruf), `p10-t` (Exit 0, kein Aufruf), `p10-i`
    (Baum mit Befund: das Skript läuft, `make -i` endet mit Exit 0; den Nachweis der Kette schützt [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 21, nicht dieses
    Gate), `p10-makeflags` (`MAKEFLAGS=n` in der Umgebung wie `-n`). Mutation für `p10-n`, `p10-q`, `p10-t`: ein `+` vor dem Rezept (eine Mutation für
    drei Fälle, weil `make` die drei Flags gleich behandelt). Das Skript urteilt über die Wurzel, die es aus seinem Pfad oder seinem Argument
    bestimmt, nicht über sein Arbeitsverzeichnis: Fall `p10-unterverzeichnis` (Aufruf aus einem Unterverzeichnis, gleicher Befund wie vom Wurzelverzeichnis);
    Mutation: die Wurzel ist das Arbeitsverzeichnis. Das Skript schreibt nichts in den Baum der Wurzel: Fall `p10-schreibt-nichts` (Prüfsumme der Dateiliste
    und Inhalte vor und nach einem roten und einem grünen Lauf gleich); Mutation: das Skript legt eine Datei im Baum an. Gleichzeitige Läufe sind
    zugelassen (`make -j` fährt das Gate neben anderen): der Aufruf von `docker` trägt weder `--name` noch ein Hilfsverzeichnis mit festem Namen
    (Fall `p3-e7`); die Aufruf-Voraussetzung von [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 8 braucht das Gate nicht.
11. **Grenzen** (akzeptierte Negative, jeweils mit Grund). Ein grünes Gate sagt nicht, dass ein Workflow auf GitHub läuft (Runner, Rechte der Organisation,
    `startup_failure`); Grund: nur dort messbar. Die Pins werden auf Form geprüft, nicht auf Gültigkeit oder Aktualität des SHA; Grund: Netz, außerhalb
    der Gates (dieselbe Grenze nennt das Modul `workflows`). Das Gate liest `.github/actions/` (zusammengesetzte Actions), wiederverwendbare Workflows
    außerhalb von `.github/workflows/` und `.github/dependabot.yml` nicht; Grund: das Repo führt keine, und ein Zusatz ist eine Änderung dieser Kennung.
    Die Regeln lesen nur die Blockform von Punkt 5; ein Anker oder `<<` in der Datei wird nicht aufgelöst, und `permissions` oder `timeout-minutes`, die ein
    Anker einträgt, gelten als fehlend; Grund: ein zweiter Parser stünde neben dem von `actionlint` und von d-check. Grün über leerem Prüfbereich
    (Punkt 2). Zeilenende `\r\n` (etwa nach einem Checkout unter Windows mit `autocrlf`): nicht zugesagt; die Regeln lesen das `\r` als Teil des Werts, der Befund ist falsch-positiv und damit sicher, und seine Meldung kann einen Wert zeigen, der richtig aussieht (`gelesen: {}`); Grund: die Plattform ist für die Gates dieses Repos nicht zugesagt, und `actionlint` selbst nimmt solche Dateien an ([`SPEC-064`](SPEC-064-ci.md) Punkt 17). `shellcheck` in `actionlint` kennt nur die Shells, die es bedient; ein `run:` mit `shell: pwsh` ist für diese Prüfung nicht gelesen.
12. **Gegenprobe.** `make workflows-check-gegenprobe` fährt das Skript an Temp-Bäumen mit einer Attrappe von `docker` (über `DOCKER`), die ihre Argumente
    aufzeichnet und Ausgabe und Exit-Code je Fall vorgibt; die Fälle der Punkte 1 bis 10 laufen ohne Netz. Dazu drei **echte** Fälle gegen das Image aus
    dem Fragment (die Temp-Bäume der echten Fälle sind mit `chmod -R a+rX` lesbar, weil das Image als Nicht-root liest) und das Image von d-check aus `DCHECK_REF` (`make` übergibt beide): `real-gruen` (eine gültige Datei, Exit 0), `real-ausdruck` (eine Datei
    mit einer Eigenschaft, die `github.event` nicht kennt: Exit 1, die Zeile von `actionlint` mit Dateiname) und `real-pins` (in einem Temp-Baum mit
    `.d-check.yml` aus Punkt 8: `uses: actions/checkout@v4` ⇒ Befund `uses-pin-missing`, ein voller SHA ohne Kommentar ⇒ `uses-pin-untagged`, ein voller
    SHA mit Tag-Kommentar ⇒ grün; je ein Aufruf des Images von d-check). Mutationen der echten Fälle: das Image des Fragments wird durch eine Attrappe ersetzt, die
    stets Exit 0 liefert (`real-ausdruck` rot); die Konfiguration des Temp-Baums nennt `workflows` nicht in `modules:` (`real-pins` rot). Das Skript der
    Gegenprobe trägt keine Abdeckungs-Deklaration (es gibt keine Anforderung des Lastenhefts dafür; `make abdeckung-check` nimmt eine Gegenprobe ohne Deklaration
    ohne Zeile an, [`SPEC-048`](SPEC-048-abdeckung.md)). Voraussetzung der echten Fälle: Docker und die beiden Images (der erste Lauf zieht sie, danach läuft die Gegenprobe
    ohne Netz); sie steht deshalb in `GEGENPROBEN_CONTAINER` und nicht in `GATE_CHECKS`.
