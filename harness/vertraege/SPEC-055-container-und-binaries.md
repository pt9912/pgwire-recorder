# SPEC-055 — Binaries je Zielplattform, Container-Image und Smoke (`binaries`, `test-container`)

`make binaries` baut das Binary für jede Zielplattform und schreibt es mit `SHA256SUMS`;
`make test-container` ist das Gate dazu: Es belegt, dass zwei Builds desselben Quellstands
dieselben Summen liefern, dass das Image je Plattform das Binary dieser Plattform trägt, und
dass das Image Aufzeichnen und Wiedergeben besteht. `make test-container-gegenprobe` prüft das
Werkzeug. Diese Kennung legt fest, was das Gate prüft und woran es rot wird. Welche Plattformen das
Produkt zusagt, steht nicht hier, sondern in der Produktspezifikation (`SPEC-035`); dieser Vertrag
sagt nur, was das Gate belegt. Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes.

Die Gegenprobe fährt in Temp-Bäumen mit einer Attrappe für `docker` (über `DOCKER`), ohne Docker
und ohne Netz; die unveränderte Kopie ist grün, und je Fall bricht eine Mutation genau ihn. Was
Punkt 4 über den echten Build zeigt, belegt der Lauf des Gates selbst (Punkt 12).

1. **Plattformen.** Die Listen stehen an einer Stelle, im Fragment `harness/mk/container.mk`:
   die Plattformen der Binaries (`linux`, `darwin`, `windows`, je `amd64` und `arm64`) und die der
   Images (`linux/amd64`, `linux/arm64`), Form `<os>/<arch>`. Rot vor dem ersten Aufruf von
   `docker`, mit Nennung des Eintrags: (a) leere Liste, (b) Eintrag ohne `/` oder mit leerem Teil,
   (c) Betriebssystem oder Architektur außerhalb dieser Mengen, (d) Eintrag doppelt, (e) eine der
   sechs Plattformen der Binaries oder eine der zwei der Images fehlt, (f) Image-Plattform, die
   nicht in der Liste der Binaries steht. Die Listen sind keine Eingabe des Aufrufers; ein
   überschriebener Wert fällt unter (a) bis (f). Jede Form hat einen Fall und eine Mutation, die nur sie bricht.
2. **Version.** `VERSION` (Vorgabe `dev`; am Ziel `test-container` `0.0.0-test`) geht in
   `-X main.version` des Builds ein und in den Dateinamen. Am Ziel `test-container` ist `dev` als
   `VERSION` rot vor dem ersten Aufruf von `docker`, mit Nennung des Werts, denn `dev` ist die
   Vorgabe des Codes: Fehlte `-X main.version=` im `Dockerfile`, gäbe `version` dasselbe aus, und
   der Smoke (Punkt 7a) belegte nichts. Ein Aufruf mit `VERSION` auf der Befehlszeile gilt vor
   dem Vorgabewert; eine Variable der Umgebung gilt nicht. Ob das echte `Dockerfile` den Wert
   einsetzt, zeigt der Smoke des Gates mit einem Wert, der nicht `dev` ist; die Gegenprobe prüft an
   der Attrappe die Vorgabe, die Ablehnung von `dev` und, dass das `Dockerfile` die Zeile
   `-X main.version=$VERSION` enthält. Sie passt auf `^[0-9A-Za-z][0-9A-Za-z.+_-]*$`, sonst rot vor dem ersten Aufruf von
   `docker`: leer, mit Leerzeichen, mit `$`, mit Anführungszeichen, mit führendem `-`, je ein Fall
   und eine Mutation. Zwei Builds derselben `VERSION` sind das, was Punkt 4 vergleicht; andere
   `VERSION` ergibt andere Summen, und das ist keine Verletzung.
3. **Binaries (`make binaries`).** Je Plattform ein Build der Stufe `dist` des `Dockerfile`; die
   Datei heißt `pgwire-recorder_<VERSION>_<os>_<arch>`, bei `windows` mit Endung `.exe`, und
   liegt in `DIST` (Vorgabe `dist`). `DIST` fehlt (wird angelegt) oder ist leer, sonst rot vor
   dem ersten Aufruf von `docker`; ein leerer Wert von `DIST` ist ungültige Eingabe, rot mit Meldung vor dem ersten Aufruf von `docker`; ein Lauf mischt sich nie mit einem früheren. Gebaut wird in ein Stage-Verzeichnis außerhalb von `DIST`; erst nach dem letzten Build gehen Binaries und `SHA256SUMS` nach `DIST`, bei einem Fehler bleibt `DIST` leer. `SHA256SUMS` hat
   das Format von `sha256sum` (Summe, zwei Leerzeichen, Name), je Binary eine Zeile, nach Name in Byte-Reihenfolge
   (C-Locale) sortiert, mit LF, und `sha256sum -c` darin ist grün. Scheitert ein Build, ist das Ziel rot und
   nennt die Plattform; `SHA256SUMS` entsteht dann nicht (keine Teilmenge). Das Ziel ist ein
   Werkzeug, kein Gate: Es fährt allein, ohne Punkt 4 bis 7.
4. **Zwei Builds, dieselbe Summe.** Das Gate fährt Punkt 3 zweimal in zwei Verzeichnissen
   außerhalb des Arbeitsbaums; der zweite Lauf umgeht den Cache der Stufen `build` und `dist`
   (`--no-cache-filter`), sodass er neu kompiliert, statt den ersten zu wiederholen. Die zwei
   Dateien `SHA256SUMS` sind byte-gleich, sonst rot mit den Namen der Binaries, deren Summen
   abweichen. Fälle: eine Summe weicht ab, je Plattform einer (sechs); ein zweiter Lauf, der
   dank Cache dieselben Bytes liefert, wo ein frischer abwiche (die Attrappe liefert ohne das
   Flag die zwischengespeicherte Fassung, mit ihm eine abweichende), und die Mutation ist der
   Lauf ohne das Flag. Je Plattform eine Mutation, die nur ihre Zeile vom Vergleich ausnimmt.
   Zugesagt ist die Gleichheit für denselben Quellstand, dieselbe `VERSION` und dieselben gepinnten
   Images des `Dockerfile` auf demselben Docker-Server; Punkt 13 nennt die Grenzen.
5. **Kopf der Binaries.** Jedes Binary beginnt mit dem Kopf seines Formats und trägt die Architektur
   seiner Plattform: ELF (Magic `7f 45 4c 46`, Maschine `3e 00` für `amd64`, `b7 00` für `arm64`),
   Mach-O 64 little-endian (Magic `cf fa ed fe`, CPU-Typ `07 00 00 01`, `0c 00 00 01`), PE (`MZ`,
   Maschine `64 86`, `64 aa`, gelesen am Offset aus Byte 0x3c). Sonst rot mit Plattform und dem
   gelesenen Wert. Je Format ein Fall mit der Architektur der anderen Plattform des Formats und ein
   Fall mit dem Kopf eines anderen Formats; je Format eine Mutation, die nur seine Prüfung bricht.
6. **Image je Plattform.** Je Image-Plattform baut das Gate die Stufe `runtime` mit
   `--platform`, unter einem Tag mit der Kennung des Laufs. Rot, mit Plattform: das Image
   nennt in `docker image inspect` ein anderes Betriebssystem oder eine andere Architektur; die
   Datei `/pgwire-recorder` im Image (aus einem erzeugten, nie gestarteten Container
   kopiert) hat eine andere SHA-256 als das Binary derselben Plattform aus Punkt 4. Je Plattform
   und Form ein Fall und eine Mutation. Eine Manifestliste entsteht nicht (Punkt 13). Der
   Bau des Images läuft ohne Emulation: Die Stufe `runtime` führt keinen `RUN` aus.
7. **Smoke im Container.** Er läuft mit dem Image der Plattform des Docker-Servers
   (`docker version --format '{{.Server.Arch}}'`; ein Wert außerhalb von `amd64` und `arm64` ist
   rot und nennt ihn). Das Netz ist `--internal`; kein Container hat ein Terminal; der Benutzer ist
   der des Images, ohne `--user`, ohne Volume und ohne Bind-Mount.
   - a. `version`: Ausgabe genau `pgwire-recorder <VERSION>` und LF, Exit-Code 0.
   - b. Aufzeichnen: PostgreSQL (`POSTGRES_IMAGE` aus `harness/mk/integration.mk`) und ein
     Container mit `record` (Upstream PostgreSQL, Listen `0.0.0.0`, Ausgabe unter `/tmp`). Ein
     Client (`psql` aus demselben PostgreSQL-Image) sendet `SELECT 1 AS eins`; seine Ausgabe (Sicht A) ist nicht leer und
     nennt `eins`. Der Container endet auf `SIGTERM` mit Exit-Code 0, und die Aufzeichnung ist
     nicht leer und per `docker cp` entnommen.
   - c. Wiedergeben: Der PostgreSQL-Container ist entfernt, **bevor** der Container mit `replay`
     startet; die Aufzeichnung kommt per `docker cp` in den erzeugten Container. Derselbe Client
     sendet dieselbe Anfrage; seine Ausgabe (Sicht B) ist gleich Sicht A. Der Container endet auf
     `SIGTERM` mit Exit-Code 0.
   Rot, mit Schritt und Wert: eine Abweichung von a, b oder c einzeln, eine leere Sicht A (zwei
   leere Sichten sind gleich), Sicht A ungleich Sicht B, ein Exit-Code ungleich 0 je Container, ein
   PostgreSQL-Container, der beim Start des Wiedergebens noch lebt, ein Aufruf mit `--user`, `-v`
   oder `--mount`. Je Form ein Fall und eine Mutation, die nur sie bricht.
8. **Reihenfolge der Gründe.** Das Gate prüft in der Ordnung der Punkte 1, 2, 4 (mit 3), 5, 6, 7 und endet
   beim ersten roten Punkt mit Exit-Code ungleich 0; ein Fehler der Eingabe (Punkt 1, 2) endet vor dem
   ersten `docker`. Ein Fall prüft, dass bei roten Punkten 1 und 5 die Meldung zu Punkt 1 kommt, einer, dass ein roter
   Punkt 4 den Bau der Images (Punkt 6) verhindert.
9. **Aufräumen.** Container, Netz und Tags der Images des Laufs sind nach Ende entfernt, auch nach
   einem roten Lauf und nach `SIGINT` und `SIGTERM`; je Fall des Endes ein Fall. Die Namen tragen die
   Kennung des Laufs. Die temporären Verzeichnisse der Punkte 4 und 7 und das Stage-Verzeichnis von `make binaries` (Punkt 3) sind ebenfalls entfernt, auch bei Fehler, `SIGINT` und `SIGTERM`.
   **Volumes** (seit slice-harness-integration-volumes): Kein Volume, das ein Container des Laufs angelegt hat, bleibt
   zurück. Das PostgreSQL-Image deklariert ein `VOLUME`; der Container `$pg` legt dafür ein anonymes Volume an, das
   `docker rm -f` ohne `-v` auf dem Host lässt. Das Gate entfernt deshalb jeden Container des Laufs mit `rm -fv`, auch
   `$pg` vor dem Wiedergeben (Punkt 7) und im Aufräumen am Ende; es gilt für grünen und roten Lauf, `SIGINT` und
   `SIGTERM` wie für Container, Netz und Tags (je Ende ein Fall `p9-*-volume`). Die Attrappe von `docker` legt für
   `$pg` ein Volume an und entfernt es nur mit `-v`; die Mutation ist `rm -f` statt `rm -fv` an beiden Stellen, einzeln.
   Das Gate liest keine Namen und prüft nicht nach (anders als `SPEC-057` Punkt 4; akzeptiertes Negativ; Grund: Nur
   `$pg` trägt ein `VOLUME`, und die Gegenprobe fängt das Weglassen von `-v`). Am echten Docker belegt der Slice es
   einmal von Hand, nach Namen vor und nach einem Lauf. Volumes, die der Lauf nicht angelegt hat, rührt das Gate nicht an.
10. **Arbeitsbaum.** Das Gate schreibt nichts in den Arbeitsbaum (die Gegenprobe vergleicht die Liste der Dateinamen vor und nach einem Lauf; das Überschreiben einer vorhandenen Datei fängt sie nicht) und liest nur, was der
    Build-Kontext (`.dockerignore`) zulässt, dazu `harness/mk/container.mk` und das Fragment der
    Integration. Dass es nur dies liest, hat keinen Fall. Es liest nichts unter `docs/plan/planning/`, `docs/reviews/`, `docs/user/`, `docs/developer/` und nicht
    die `README.md` (SPEC-052 Punkt 13).
11. **Aufrufarten des Treibers.** Das Gate schreibt keinen Stand; geglaubt wird allein sein
    Exit-Code in der Kette von `make gates` ([SPEC-052](SPEC-052-gate-kette.md) Punkt 1). Das Rezept
    jedes der drei Ziele (`binaries`, `test-container`, `test-container-gegenprobe`) besteht aus einer Zeile ohne
    `$(MAKE)` und ohne `+`; unter `-n`, `-q` und `-t` fährt es kein Docker und schreibt nichts, und
    `MAKEFLAGS` aus der Umgebung gilt wie das Flag. Unter `-i` überliest `make` den Exit-Code des
    Skripts: Ein Aufruf mit `-i` belegt nichts, und `make -i gates` schreibt keinen Nachweis
    (SPEC-052 Punkt 6 und 21). Bei Kombination von `-n`, `-q` oder `-t` mit `-i` gilt `-n`, `-q`,
    `-t`: kein Docker. Je Flag `-n`, `-q`, `-t`, `-i`, je `MAKEFLAGS` aus der Umgebung ein Fall am
    Ziel `test-container` und einer am Ziel `binaries`; die Mutation ist ein Rezept mit `$(MAKE)`.
12. **Ausnahmeliste und Gegenprobe.** `test-container` und `test-container-gegenprobe` hängen an
    `GATE_CHECKS` und stehen in der Ausnahmeliste der schnellen Kette (SPEC-052 Punkt 2 und 4):
    das erste startet Container und baut zwölf Binaries, die zweite hängt am Namen. Die
    Gegenprobe prüft die Punkte 1 bis 11 an der Attrappe; ob das `Dockerfile` selbst
    deterministisch baut, zeigt jeder Lauf des Gates, nicht die Gegenprobe (Punkt 13).
13. **Grenzen.** (1) Reproduzierbar ist das Binary, nicht das Image; das Image wird über seinen
    Digest identifiziert (`LH-QA-01`, Messmethode „Bereitstellung“ nennt nur das Binary). (2) Das Gate
    vergleicht zwei Builds auf demselben Docker-Server mit denselben Eingangs-Images; es prüft nicht, dass die
    Eingangs-Images (Go-Toolchain, Basis-Image) per Digest gepinnt sind, und baut nicht auf einem
    zweiten Host oder mit einer anderen Architektur des Bau-Hosts (Pflege der Pins ist eine Handlung
    nach [ADR-0026](../../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md); eine angehobene Toolchain ändert die Summen, und das ist keine Verletzung).
    (3) Die Manifestliste für `linux/amd64` und `linux/arm64` in einer Referenz entsteht bei der
    Veröffentlichung, nicht in diesem Gate; der Smoke läuft nur mit dem Image der Plattform des
    Docker-Servers, das andere ist nur nach Punkt 6 belegt (kein Emulator vorausgesetzt).
    (4) Ob die Binaries für `darwin` und `windows` laufen, prüft kein Gate; belegt sind
    Bau, Kopf und Summe. (5) Gleichzeitige Läufe, auch `make -j gates`, sind nicht zugesagt; die
    Namen tragen die Kennung des Laufs, sodass sie sich nicht stören, aber der Cache des Builds ist geteilt. (6)
    `SIGKILL` und der Ausfall des Docker-Dienstes sind nicht zugesagt; ein Zeitlimit des ganzen Laufs gibt es nicht.
    (7) Die Registry, der Tag `latest` und die Veröffentlichung sind nicht Gegenstand
    (`slice-erster-release-veroeffentlichung`).
