# Slice slice-v1-abschluss-binary-namen: Release-Binaries ohne Version im Dateinamen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität)

**Berührte Spec-Stellen:** `SPEC-055` (Punkt 2 und 3, [`harness/vertraege/`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md)) · `SPEC-067` (Punkt 4 und 19 (3), (4), (7), (9), [`harness/vertraege/`](../../../../harness/vertraege/SPEC-067-release.md)) · `LH-FA-01` · `LH-FA-19.a` · `LH-QA-03`

**Entscheidung:** [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) (Accepted, Entscheidung des Nutzers vom 2026-10-20: keine Version im Dateinamen der Binaries, Schema `pgwire-recorder-<os>-<arch>`)

**Verantwortlich:** pt9912
**Autor:** Architect. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make binaries`, `make test-container` und die Pipeline (`release.sh`) schreiben, erwarten und prüfen die Binaries als `pgwire-recorder-<os>-<arch>` (unter Windows `.exe`) ohne Version im Namen; das Image behält den Tag mit Version. Schema, Dateiliste, Zeilen von `SHA256SUMS`, Eingaben und Fälle entscheidet [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 2 und 3 und [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 4; warum so, [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md). Der Slice liegt vor `slice-v1-abschluss-homebrew`, dessen Vorlage und Fälle die Namen brauchen (Bindung unten).

**Sendung an `slice-v1-abschluss-homebrew`** (`AGENTS.md` §3.13, Bindung): Die Vorlage der Formel, die Prüfung der Summen (`SPEC-067` Punkt 19 (3)) und ihre Fälle `p19-*` rechnen mit den Namen `pgwire-recorder-<os>-<arch>`, die dieser Slice einführt; der Fall `p19-summen-altes-schema` (Namen des alten Schemas, Exit 1 `summen`) liegt dort. Beides steht dort in §1 und §6 mit der Kennung dieses Slice und als Start-Trigger.

**Sendung an `slice-erster-release-homebrew-nachweis`** (Bindung): Der Nachtrag der Formel für `v0.1.0` entfällt; die erste Formel entsteht durch den Job `homebrew` mit dem nächsten stabilen Tag, dessen Binaries das neue Schema tragen. Das steht dort in §1, DoD, §4 und §6 mit der Kennung dieses Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Vorlage der Formel, der Schritt und der Job `homebrew` — `slice-v1-abschluss-homebrew` (Klasse 1; dort §1 und §6 nehmen die Namen an). Grund: ein eigener Liefergegenstand mit eigener Gegenprobe; dieser Slice ändert nur Namen und ihre Prüfung.
- Ein neuer Tag und ein neuer Release (der erste mit dem neuen Schema) — Handlung des Nutzers mit neuem Ja zum Tag (Klasse 3, anderer Vorgang). Den Tag setzt der Nutzer; der Slice, der ihn vorbereitet und belegt, ist `slice-v1-abschluss-release-v0-1-1` (Welle `welle-v1-nachtrag`, Reihenfolge 4), und `slice-erster-release-homebrew-nachweis` wartet auf ihn (dort §4).
- Der Abschnitt „Binary vom GitHub-Release“ im Handbuch und die Zeile zu `v0.1.0` in `README.md` — Bestand bleibt bewusst stehen (Klasse 2). Grund: Sie beschreiben, was `v0.1.0` trägt (Namen mit Version, einmal gelaufen); das bleibt wahr, bis ein neuer Tag erscheint. `slice-v1-abschluss-release-v0-1-1` schreibt sie nach `SPEC-067` Punkt 18 um (dort §1 und DoD 2 nehmen die Sendung an); dieser Slice ändert nur den Abschnitt zu `make binaries`.
- Der Release `v0.1.0` selbst — unverändert (Entscheidung des Nutzers vom 2026-10-20); ein Tag wird nie verändert (`SPEC-067` Punkt 11).
- Das Image, seine Tags und die Registries — Bestand bleibt (Klasse 2); die Entscheidung gilt nur den Binaries.
- Änderung des Lastenhefts — es legt kein Namensschema fest (`LH-FA-01`, `LH-QA-03` nennen Plattformen); eine Anhebung der Version nach `harness/conventions.md` entfällt.
- Produkt-Code (`cmd/`, `internal/`) und `Dockerfile` — Schicht-Abgrenzung (Klasse 4): `VERSION` geht weiter in `-X main.version`; der Name entsteht allein in `binaries.sh`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung): `tools/harness/binaries.sh`, `run-container.sh` und `release.sh` (Dateiliste des Schritts `binaries`, Summe je Plattform im Schritt `image`) schreiben und erwarten `pgwire-recorder-<os>-<arch>` ohne Version, `.exe` allein bei `windows`; `make test-container-gegenprobe` und `make release-gegenprobe` tragen die Fälle und Mutationen aus `SPEC-055` Punkt 3 (`p3-namen`, `p3-exe`, `p3-name-ohne-version`, `p3-name-trenner`, `p3-sums-zeile-<os>-<arch>`, `p3-dist-andere-version`, `p3-dist-altes-schema`; die Mutation `p3-name` kehrt sich um) und `SPEC-067` Punkt 4 (`p4-altes-schema`, `p4-namen-je-plattform`) und führen jeden alten Namen mit Version nach; je Zusage Mutation und roter Fall im Bericht (`AGENTS.md` §3.10); die Zeile von `make binaries` in `harness/README.md` §Sensors nennt das neue Schema.
- [x] Benutzerhandbuch (Abschnitt zu `make binaries`: Namen, Beispiele mit `dev`, Aufruf des Binaries) und `docs/maintainer/releasing.md` (§2 *Was lokal läuft*, §3 *Versionierung*) beschreiben das Ausgeführte im Ist-Zustand ohne Chronik (`AGENTS.md` §3.11); der Abschnitt „Binary vom GitHub-Release“ und die Zeile zu `v0.1.0` in `README.md` bleiben (§1); die Beispiele des Handbuchs, die Namen nennen, liegen als Datei vor und laufen gegen das gebaute Binary, soweit sie es aufrufen (`dist/pgwire-recorder-linux-amd64 version`).
- [x] `make gates` grün; `make test-container-gegenprobe` und `make release-gegenprobe` am Endstand (`AGENTS.md` §6).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/binaries.sh` | update | Name `pgwire-recorder-${os}-${arch}` (Marke `P3name`), Meldung bei `DIST` unverändert |
| `tools/harness/run-container.sh` | update | `binaer()` nennt den neuen Namen (Punkt 4 und 6 von `SPEC-055`) |
| `tools/harness/release.sh` | update | `erwartete_namen` ohne Version; Summe im Schritt `image` aus `pgwire-recorder-linux-<arch>` |
| `tools/harness/container-gegenprobe.sh` | update | `bname`, `namen_soll` (ohne Version), Muster der Summen-Zeile, Meldung zu Punkt 4; neue Fälle `p3-name-trenner`, `p3-name-ohne-version`, `p3-sums-zeile-<os>-<arch>`, `p3-dist-andere-version`, `p3-dist-altes-schema`, `p3-exe` je Architektur, und ihre Mutationen (`.exe`, Name, Trenner, Versionsdatei, `DIST`) |
| `tools/harness/release-gegenprobe.sh` | update | Attrappe von `make` (Namen; Modi `altes-schema` und `datei-fehlt` mit `GP_FEHLT=<os>-<arch>`), `NAMEN`, `ALLE7`, Fälle zu Punkt 5 und 9; Fälle `p4-altes-schema`, `p4-namen-je-plattform-<os>-<arch>` und ihre Mutationen (`m-p4-altes-schema-angenommen`, `m-p4-liste-nicht-gezaehlt`, `m-p9-name-fehlt-<os>-<arch>` mit dem Fall der Plattform), Mutationen des Namens der Summe im Schritt `image` (`m-p5-summe-name-arch-fest`, `m-p5-summe-name-os`) |
| `harness/README.md` (§Sensors, Zeile `make binaries`) | update | nennt das Schema, `DIST` auch bei Dateien einer anderen `VERSION`, die Bindung [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) |
| `docs/user/benutzerhandbuch.md` (Abschnitt zu `make binaries`) | update | Namen und Beispiele im Ist-Zustand; `dist/` auch bei Dateien einer anderen `VERSION` nicht leer |
| `docs/maintainer/releasing.md` | update | §2, §3 (Namen tragen keine Version mehr), Version 0.11 und Zeile der Änderungshistorie |

Nicht in dieser Liste: `Dockerfile`, `harness/mk/container.mk` und `cmd/`, `internal/` (§1), der Abschnitt zu `v0.1.0` im Handbuch und `README.md` (§1), die Verträge `SPEC-055` und `SPEC-067` (liegen vor dem Code, Architect).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) ist `Accepted` (Entscheidung des Nutzers vom 2026-10-20), und `welle-erster-release` ist `done` (Start-Trigger von `welle-v1-nachtrag`). Der Slice steht in der Reihenfolge der Welle vor `slice-v1-abschluss-homebrew` und blockiert ihn (dessen Start-Trigger nennt ihn).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat 3d927615..<letzter Commit des Slice> -- tools harness/README.md docs/user docs/maintainer README.md` (Tests eingeschlossen, `AGENTS.md` §3.14) **über 800** liegt. Schätzung des Architect in dieser Metrik: Produktcode rund 30 (die drei Skripte, je ersetzte Zeile zwei), Doku rund 80 (Handbuch, `releasing.md`, `harness/README.md`), die Gegenprobe rund 250 bis 400 (viele ersetzte Namen, dazu neue Fälle und Mutationen; Faktor von Tabellen- und Gegenprobe-Code zum Produktcode rund zwei und mehr), zusammen 360 bis 510. Kein vorgemerkter Schnitt: Der Slice hat zwei Liefer-Punkte, und die Obergrenze liegt unter der Marke.
- `in-progress` → `open`: Der Nutzer nimmt [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) nicht an oder wählt ein anderes Schema — zurück in die Planung; die Verträge `SPEC-055` und `SPEC-067` ziehen mit der Entscheidung nach (Architect).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des Vertrags** (`AGENTS.md` §3.12; alle vor dem Code entschieden, jede an der genannten Stelle der Verträge):

| Randform | Entscheidung | Stelle |
|---|---|---|
| (a) Schema | `pgwire-recorder-<os>-<arch>`, Bindestriche, klein; `.exe` allein bei `windows`; sechs Namen | `SPEC-055` Punkt 3 |
| (b) Woher die Version, wenn der Name sie nicht trägt | aus `VERSION` wie bisher (`-X main.version`); nach dem Lauf steht sie nur im Binary (`version`), eine Versionsdatei in `DIST` entsteht nicht; die Pipeline bindet sie durch den Aufruf `VERSION=<Version>` und für `linux/amd64` durch Summe im Image und Smoke | `SPEC-055` Punkt 2 und 3; `SPEC-067` Punkt 4 |
| (c) Zeilen von `SHA256SUMS` | `<Summe>  <Name>`, sechs Zeilen, C-Locale sortiert, LF; Namen nach (a) | `SPEC-055` Punkt 3 |
| (d) Mehrere Versionen im selben `DIST` | `DIST` muss fehlen oder leer sein; Dateien einer anderen `VERSION` oder des alten Schemas sind ein nicht leeres `DIST`, rot vor dem ersten `docker` | `SPEC-055` Punkt 3 (`p3-dist-andere-version`, `p3-dist-altes-schema`) |
| (e) Eingaben und Fehler | das Muster von `VERSION` bleibt (Wert geht weiter in `-X` und in die Befehlszeile); `dev` am Ziel `test-container` bleibt rot; keine Lockerung eines Gates | `SPEC-055` Punkt 2 |
| (f) Release-Pipeline | Dateiliste der sechs und `SHA256SUMS`, der alte Name ist „zu viel“ und „fehlt“; sieben Anhänge bleiben; Summe im Image aus `pgwire-recorder-linux-<arch>` | `SPEC-067` Punkt 4, 5 und 9 |
| (g) Altes Schema (`v0.1.0`) | wird nicht gelesen und nicht übersetzt; der Schritt `homebrew` endet bei diesen Namen mit Exit 1 `summen` (Fall in `slice-v1-abschluss-homebrew`) | `SPEC-067` Punkt 19 (3), (7) |
| (h) Plattform | Namen aus Kleinbuchstaben, Ziffern, `-`, `.`: gültig auf Linux, macOS, Windows; erzeugt und geprüft nur unter Linux im Docker, macOS und Windows Grenze | `SPEC-055` Punkt 3 und 13 (4) |
| (i) Aufrufarten, Abbruch | unverändert; die Rezepte und die Fälle zu `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`, `SIGINT`, `SIGTERM` bleiben, nur ihre Namen ziehen nach | `SPEC-055` Punkt 9 und 11 |

- Fundstellen des alten Schemas, die dieser Plan nicht nennt (`grep -rn 'pgwire-recorder_' --exclude-dir=.git --exclude-dir=baseline`, ohne `docs/reviews/` und `done/`; Stand beim Schreiben: die fünf Skripte, die zwei Verträge, `harness/README.md`, Handbuch, `releasing.md`, `README.md`, `slice-v1-abschluss-homebrew`) — **Ausgang:** *entfallen*: Review (Schwerpunkt 1) und Verifikation lasen jeden Treffer am Endstand einzeln; keiner ist ein gültiger Name, jeder nennt das alte Schema als alt oder beschreibt `v0.1.0`, und in `binaries.sh`, `run-container.sh`, `release.sh`, `harness/README.md`, `releasing.md` und `README.md` bleibt keiner (Zahlen in §7).
- Der Abschnitt zu `v0.1.0` in Handbuch und `README.md` nennt Namen, die der nächste Tag ablöst; wer ihn umschreibt, ist `slice-v1-abschluss-release-v0-1-1` (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`) — **Ausgang:** *weiter offen*, übergeben an `slice-v1-abschluss-release-v0-1-1` (dort §1, DoD 2 und §8 nennen die Sendung mit der Kennung dieses Slice; `SPEC-067` Punkt 18) und ins Register eingetragen (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, Beleg `evidence/slice-v1-abschluss-binary-namen.md`).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

### Belege des Implementers

Klon im Dev-Container `pgr-dev2`, Zweig `impl-namen` ab `de41645c`. Code- und Doku-Endstand `9492d6c4`; danach ändert nur dieser Abschnitt den Zweig. Die Randformen aus §6 waren vollständig; eine Rückgabe an den Architect war nicht nötig.

**Läufe am Endstand `9492d6c4`**

| Lauf | Ergebnis |
|---|---|
| `make gates` | Exit 0 (alle Gates der Kette, darunter `docs-check`, `kopf-check`, `abdeckung-check`, `test-container` mit dem echten `Dockerfile`: „Punkt 4 grün — zwei Builds, dieselben Summen“, „Punkt 7 grün — Smoke auf linux/amd64“) |
| `make test-container-gegenprobe` | Exit 0, „unveränderte Kopie grün, je Mutation genau die benannten Fälle rot“ (`GP_ENTDECKE=1`, Mutationen mit den roten Fällen im Lauf-Log) |
| `make release-gegenprobe` | in fünf Teilen (Auswahl `GP_MUT`, Absatz darunter): Teile 1 bis 3 und 5 Exit 0, „unveränderte Kopie grün“; Teil 4 (`m-p16-*`) sechs Mutationen rot, allein gefahren grün, am Stand `de41645c` ebenso rot (Vorbefund) |
| Metrik (`AGENTS.md` §3.14) | `git diff --numstat de41645c..9492d6c4 -- tools harness/README.md docs/user docs/maintainer README.md`: 139 Einfügungen, 77 Löschungen, Summe 216 (Marke in §4: über 800); dieselbe Zahl mit der Basis `3d927615` aus §4, denn die Pfade der Liste wurden zwischen den Basen nicht geändert. Am Endstand `bb5d47ea` (Closure): 140 + 77 = **217**, die eine Zeile mehr ist die Mutation `p5-name-altes-schema` aus `c1326b0c`; die Basis `3d927615` liefert dieselbe Zahl |

Die Gegenprobe der Veröffentlichung lief nicht in einem Zug: Ihr Arbeitsordner wächst mit jeder der 535 Mutationen (5 GB am Ende), und der Datenträger des Containers hatte 7,5 GB frei, geteilt mit dem Host, bei einem Lauf voll (0 frei). Die Teile (`GP_ENTDECKE=1`, je `make release-gegenprobe` mit der Auswahl): `^(ok|m-p1-|m-p2-|m-p3-|m-p4-|m-p5-)` (289 Zeilen), `^m-p(6|7|8|9)-` (80), `^m-p1[0-5]-` (36), `^(rekursion|fall-ohne)`, alle Exit 0; `^m-p16-` (130 Mutationen) Exit 2 mit sechs roten Mutationen: `m-p16-marke-loescht-nach-auswahl`, `m-p16-auswahl-gleiche-variante-doppelt`, `m-p16-blobs-leer`, `m-p16-index-mehr-als-zwei-abgelehnt`, `m-p16-entpacken-je-plattform`, `m-p16-auswahl-variante-ausgeschlossen`. Allein gefahren (`GP_MUT` mit diesen sechs Namen) sind sie grün („grün (6 Läufe)“). Am Stand `de41645c` (`git worktree`, ohne diesen Slice) sind in diesem Container mit `GP_MUT='^m-p16-'` dieselben sechs rot (Exit 1): ein Vorbefund der Umgebung, nicht dieses Slice; der Schritt `scan` ist nicht berührt, und die Ursache ist nicht untersucht. Ein Lauf ohne Auswahl am Stand `24a1d1d1` (bis auf die in `9492d6c4` nachgezogenen Listen der Mutationen von Punkt 3 dieser Kopie, die drei Mutationen `m-p3-*` rot färbten) war auch bei `m-p16-*` grün, damals mit 7,5 GB frei. Hinweis für den Verifier: `make release-gegenprobe` in einem Zug braucht rund 6 GB freien Platz; bei Rot in `m-p16-*` zuerst die rote Mutation allein fahren.

**Zusagen und Mutationen** (je Zeile: Zusage · Mutation · roter Fall; „+“ nennt Fälle, die dieselbe Mutation mit färbt, weil sie denselben Namen lesen; Mutationen sind Textersatz in einer Kopie, jede einmal rot gesehen)

`SPEC-055` Punkt 3 (`make test-container-gegenprobe`, Skript `container-gegenprobe.sh`, Mutation an `binaries.sh`):

| Zusage | Mutation | roter Fall |
|---|---|---|
| der Name ist `pgwire-recorder-<os>-<arch>`, genau die sechs Namen und `SHA256SUMS` | `p3-name`: Zeile `# P3name` wieder `pgwire-recorder_${VERSION}_${os}_${arch}` (Umkehr) | `p3-namen`; + `p3-exe`, `p3-name-ohne-version`, `p3-name-trenner`, `p3-sums-format`, `p3-sums-je-binary`, `p3-sums-zeile-*` |
| der Name enthält keine `VERSION` (gleiche Namen bei zwei Läufen) | `p3-name-version`: `pgwire-recorder-${VERSION}-${os}-${arch}` | `p3-name-ohne-version` (der Name allein, ohne `_`); + `p3-namen`, `p3-exe`, `p3-sums-*`; `p3-name-trenner` bleibt grün |
| zwei Läufe mit verschiedener `VERSION` haben verschiedene Summen | `p3-build-arg`: `--build-arg "VERSION=$VERSION"` entfällt | `p3-name-ohne-version`; + `p3-version-im-build` |
| der Trenner ist `-`, nicht `_` | `p3-name-trenner` (beide `_`), `p3-name-trenner-vorn` (nur zwischen `pgwire-recorder` und `<os>`), `p3-name-trenner-hinten` (nur zwischen `<os>` und `<arch>`) | je `p3-name-trenner`; + `p3-namen`, `p3-exe`, `p3-sums-format`, `p3-sums-je-binary`, `p3-sums-zeile-*` |
| die Zeichenmenge (Kleinbuchstaben, Ziffern, `-`, `.`) | `p3-name-gross` (`${os^}`), `p3-name-ziffern` (`${arch//64/}`, entfernt die Ziffern), `p3-exe-ohne-punkt` (`_exe`), `p3-name-trenner` (`_`) | je `p3-namen`; + `p3-exe`, `p3-sums-je-binary`, `p3-sums-zeile-*`; `p3-sums-format` bleibt bei den Ziffern grün (Muster `[a-z0-9]+`), `p3-name-trenner` färbt `_exe` |
| `.exe` trägt allein `windows`, je Architektur | `p3-exe` (Zeile `# P3exe` entfällt), `p3-exe-nur-amd64`, `p3-exe-nur-arm64`, `p3-exe-auch-linux`, `p3-exe-auch-darwin` | je `p3-exe`; + `p3-namen`, `p3-sums-je-binary`, `p3-sums-zeile-` der betroffenen Plattform |
| die Dateiliste von `DIST` ist genau die sechs Binaries und `SHA256SUMS`, keine Versionsdatei | `p3-versionsdatei`: `echo "$VERSION" > VERSION` nach `sha256sum` | `p3-namen` allein |
| `SHA256SUMS` führt je Plattform die Zeile mit diesem Namen | `p3-sums-zeile-<os>/<arch>` (6): `grep -v -e "-<os>-<arch>"` hinter `sha256sum` | je `p3-sums-zeile-<os>-<arch>`; + `p3-sums-format`, `p3-sums-je-binary` |
| ein `DIST` mit den sechs Dateien eines Laufs anderer `VERSION` ist rot vor dem ersten `docker`, `DIST` unverändert | `p3-dist-andere-version`: die Prüfung nimmt `pgwire-recorder-*` und `SHA256SUMS` aus; `p3-dist-unveraendert`: `rm -rf "${dist:?}"/*` vor der Meldung | `p3-dist-andere-version`; + `p3-zweiter-lauf` (ein Name allein trennt „andere“ von „gleiche“ `VERSION` nicht); bei `p3-dist-unveraendert` `p3-dist-andere-version`, `p3-dist-altes-schema`, `p3-dist-nicht-leer` |
| ein `DIST` mit einer Datei des alten Namens ist rot, ebenso | `p3-dist-altes-schema`: die Prüfung nimmt `pgwire-recorder_*` aus | `p3-dist-altes-schema` allein |
| (Bestand, Eingaben, nun mit neuen Fällen in der Liste) `DIST` vor dem ersten `docker` | `p3-vor-docker`: `exit 2` entfällt; `p3-dist`: Zeile `# P3dist` entfällt; `p3-dist-fehlt`: `mkdir -p` entfällt | `p3-dist-nicht-leer`, `p3-zweiter-lauf` + `p3-dist-andere-version`, `p3-dist-altes-schema`; `p3-dist-fehlt` färbt zusätzlich die Fälle, die `out` brauchen |

`SPEC-055` Punkt 4 (Namen in der Meldung): Mutationen `p4-abweichung-<os>/<arch>` (6; `grep -v -e "-<os>-<arch>"` im Vergleich) → je `p4-abweichung-<os>-<arch>` (`linux-amd64` zusätzlich `p8-4-vor-6`), `p4-namen` (Meldung ohne die Namen aus `diff`) → die sechs `p4-abweichung-*` und `p4-namen`. Der Fall `p4-namen` erwartet die Zeile `die Summen zweier Builds weichen ab: pgwire-recorder-darwin-arm64 pgwire-recorder-linux-amd64`.

`SPEC-055` Punkte 3 bis 5 (Name im Prüfer, `binaer()` in `run-container.sh`; Befund F-1 der Prüfung): Die Zusage „der Prüfer erwartet `pgwire-recorder-<os>-<arch>` (Kopf, Summe, Smoke, Image), `.exe` allein bei `windows`“ fängt die Mutation `p5-name-altes-schema` (`local n="pgwire-recorder-$1-$2"` wird `pgwire-recorder_$1_$2`) mit den roten Fällen `p2-gueltig`, `p4-gruen`, `p5-gruen` und den zwölf `p5-arch-*` und `p5-format-*` (je Plattform); die Mutation läuft mit `nur="p2-gueltig p4-gruen p5"`, die übrigen Fälle, die denselben Namen lesen (`p6-*`, `p7-*`, `p9-*`), liegen außerhalb dieser Auswahl und färbt sie im ungewählten Lauf mit (ohne Auswahl rot gesehen: 54 Fälle, `GP_ENTDECKE=1`). Gesehen mit `GP_MUT='^p5-name-altes-schema$' bash tools/harness/container-gegenprobe.sh`: Exit 0, „je Mutation genau die benannten Fälle rot“ (unveränderte Kopie grün).

`SPEC-067` Punkt 4, 5 und 9 (`make release-gegenprobe`, Skript `release-gegenprobe.sh`, Mutation an `release.sh`):

| Zusage | Mutation | roter Fall |
|---|---|---|
| genau die sechs Dateien `pgwire-recorder-<os>-<arch>` und `SHA256SUMS`, sonst Exit 1 (Punkt 4) | `m-p4-liste-entfernt` (Zeile `# P4liste` entfällt) | `p4-datei-fehlt`, `p4-datei-zu-viel`, `p4-altes-schema`, `p4-namen-je-plattform-*` (6) |
| der alte Name `pgwire-recorder_<Version>_<os>_<arch>` wird nicht angenommen | `m-p4-altes-schema-angenommen`: die Liste nimmt auch Dateien mit `pgwire-recorder_` an | `p4-altes-schema` allein |
| die Liste prüft jeden der sechs Namen (je Plattform) | `m-p9-name-fehlt-<os>-<arch>` (6): `erwartete_namen` lässt einen Namen aus | `p4-namen-je-plattform-<os>-<arch>` der Plattform (bei `linux-arm64` auch `p4-datei-fehlt`); + die Fälle, die den Anhang prüfen (`p9-anhaenge`, `p4-aufruf` u. a.) |
| die Dateiliste wird gezählt (Überzähliges ist rot) | `m-p4-liste-nicht-gezaehlt`: die Liste prüft nur, dass die sechs Namen da sind | `p4-datei-zu-viel` allein |
| `VERSION` im Aufruf, `DIST` außerhalb des Baums | `m-p4-version-fehlt`, `m-p4-dist-im-baum`, `m-p4-dist-vorher-angelegt` | `p4-aufruf` (+ `p2-trockenlauf-ohne-aussenwirkung` bei `VERSION`); die Version steht nicht mehr im Namen, die übrigen Fälle der Dateiliste färben diese Mutation nicht mehr (`SPEC-067` Punkt 4: die Bindung ist der Aufruf) |
| `sha256sum -c` und Exit des Aufrufs | `m-p4-summe-entfernt`, `m-p4-make-fehler-ignoriert` | `p4-summe-falsch`, `p4-make-scheitert` (unverändert) |
| die Summe im Image kommt aus dem Eintrag `pgwire-recorder-linux-<arch>` (Punkt 5) | `m-p5-summe-name-arch-fest` (immer `linux-amd64`), `m-p5-summe-name-os` (`darwin` statt `linux`) | alle Fälle mit grünem Schritt `image` (53 und 55; Liste in der Mutation) |
| sieben Anhänge, die Namen der sechs und `SHA256SUMS` (Punkt 9) | `m-p9-achter-anhang`, `m-p9-name-fehlt-sha256sums` | `p9-anhaenge`, `p9-release-existiert-gleich`, `p11-wiederanlauf-release-*` (unverändert, Liste nun ohne Version) |

Grüne Mutanten: keine. Keine Mutation wird als äquivalent geführt; jede der Zeilen oben war rot (`GP_ENTDECKE=1`), und die Mengen der roten Fälle stehen als erwartete Menge im Skript (`mut` und `mutr` melden eine Abweichung als rot). Einen Mitläufer nennt die Tabelle mit „+“: Er liest denselben Namen (`p3-namen` prüft alle sechs, `p3-sums-zeile-*` je einen), er ist kein zweiter Beleg.

**Befundzeilen des alten Schemas** (`grep -rnI 'pgwire-recorder_' --exclude-dir=.git --exclude-dir=baseline .`, ohne `docs/reviews/` und `done/`)

- vorher (`de41645c`, `git grep -nI 'pgwire-recorder_'` ohne `.harness/baseline`, `docs/reviews/` und `done/`): 50 Treffer in 12 Dateien: `binaries.sh` 1, `run-container.sh` 1, `release.sh` 2, `container-gegenprobe.sh` 6, `release-gegenprobe.sh` 13, `SPEC-055` 1, `SPEC-067` 2, `harness/README.md` 1, Handbuch 18 (davon 10 im Abschnitt „Binary vom GitHub-Release“), `slice-v1-abschluss-release-v0-1-1` 1, dieser Plan 1, [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) 3.
- nachher (`9492d6c4`): Treffer nur noch, wo das alte Schema als alt genannt oder `v0.1.0` beschrieben wird: `container-gegenprobe.sh` 6 (Fall `p3-dist-altes-schema` und die Mutationen `p3-name*`; am Endstand `bb5d47ea` sieben, dazu die Mutation `p5-name-altes-schema` aus `c1326b0c`), `release-gegenprobe.sh` 2 (Modus `altes-schema` der Attrappe, Mutation `m-p4-altes-schema-angenommen`), Handbuch 10 (Abschnitt „Binary vom GitHub-Release“, `v0.1.0`, Bestand nach §1), `SPEC-055` 1 und `SPEC-067` 2 (Nennung des alten Schemas im Vertrag), [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) 3, `slice-v1-abschluss-release-v0-1-1` 1, dieser Plan; `binaries.sh`, `run-container.sh`, `release.sh` und `harness/README.md` ohne Treffer

**Handbuch** (`AGENTS.md` §3.11, Schritt 17): Die Proben liefen gegen das gebaute Binary im Dev-Container (`make binaries VERSION=1.2.3` in `dist/`, danach `make binaries` mit `DIST=/tmp/d2`, danach `make binaries` in `dist/`).

| Satz im Abschnitt „Binaries für alle Plattformen bauen“ | Probe, Ergebnis |
|---|---|
| In `dist/` liegen sechs Binaries und `SHA256SUMS` | `ls -A dist`: sieben Einträge |
| der Name ist `pgwire-recorder-<Betriebssystem>-<Architektur>`, unter Windows `.exe` | die sechs Namen der Liste im Handbuch, Zeichen für Zeichen |
| die Version steht nicht im Namen; die Namen bleiben dieselben | `VERSION=1.2.3` und `VERSION` ohne Angabe: `ls -A` beider Verzeichnisse gleich |
| `VERSION=1.2.3` gibt `pgwire-recorder 1.2.3` aus statt `pgwire-recorder dev` | `dist/pgwire-recorder-linux-amd64 version`: `pgwire-recorder 1.2.3`; ohne Angabe: `pgwire-recorder dev` |
| `SHA256SUMS` hat eine Zeile je Binary, jede Zeile der Ausgabe von `sha256sum -c` endet auf `OK` | sechs Zeilen `…: OK`; die Summen der beiden Läufe unterscheiden sich an allen sechs Zeilen |
| `dist/` muss leer sein oder fehlen, sonst bricht `make binaries` ab; das gilt auch für Dateien eines früheren Laufs mit anderer `VERSION` | zweiter Lauf mit `VERSION=2.0.0` in `dist/` des ersten: `binaries: Eingabe: DIST 'dist' ist nicht leer`, Exit 2, `dist/` unverändert (sieben Einträge) |
| `./dist/pgwire-recorder-linux-amd64 version` gibt `pgwire-recorder dev` aus | `rm -rf dist; make binaries` (Vorgabe `dev`), dann `dist/pgwire-recorder-linux-amd64 version`: `pgwire-recorder dev`; die Namen in `dist/` sind die sechs Zeilen des Handbuchs (`diff` gegen den Block „text“, leer) |

Quantoren-Durchsicht (Schritt 17): „jede der sechs Plattformen“ (Bestand, sechs Dateien in `dist/`), „Jede Zeile der Ausgabe“ (Bestand, sechs `OK`), „muss leer sein oder fehlen“ (Bestand; die Ausnahme „Dateien einer anderen `VERSION`“ steht nun im Satz, Probe oben), „die Namen bleiben dieselben“ (neu, Probe oben). Eine Ausnahme, die die Seite für diesen Gegenstand sonst aufzählt, gibt es nicht. Der Abschnitt „Binary vom GitHub-Release“ und die Zeile zu `v0.1.0` in `README.md` bleiben (§1).

**Kommentare** (Schritt 20): Die geänderten Kommentare in `container-gegenprobe.sh` (`namen_trenner`, `namen_gleich_summen_verschieden`, „Zwei Läufe mit verschiedener VERSION …“, „Ein DIST mit den Dateien …“, „Eine gültige Version geht in den Build ein“) und `release-gegenprobe.sh` (`GP_FEHLT`) sagen den Ist-Zustand und nennen je einen Fall, der sie prüft.

### Belege am Endstand (Review, Verifikation, Closure)

Zweig `ver-namen`, Endstand `bb5d47ea`; der Code ist seit `c1326b0c` unverändert, danach kamen nur Berichte unter `docs/reviews/`. Review-Report: [Review-Report](../../../reviews/2026-10-20-review-slice-v1-abschluss-binary-namen.md) (0 HIGH, 0 MEDIUM, F-1 LOW behoben in `c1326b0c`, F-2 und F-3 INFO). Verifikationsbericht: [Verifikationsbericht](../../../reviews/2026-10-20-verifikation-slice-v1-abschluss-binary-namen.md) (0 HIGH, 0 MEDIUM, F-V1 LOW, F-V2 bis F-V4 INFO).

| Lauf (Angabe der Koordination, Host) | Stand | Ergebnis |
|---|---|---|
| `make release-gegenprobe` in einem Zug, `TMPDIR=/Development/tmp-gp` (Arbeitsordner auf anderer Platte) | `bb5d47ea` | „release-gegenprobe: grün (538 Läufe)“, Exit 0; die sechs `m-p16-*` aus dem Lauf des Implementers sind grün |
| `make gates` | `b70e9c95` (Code gleich `bb5d47ea`, danach nur Berichte) | grün |
| `make test-container-gegenprobe` | `b70e9c95` (Verifikation, Auswahl) und `9492d6c4` (Implementer, voll) | grün; die unveränderte Kopie grün, je Mutation genau die benannten Fälle rot |

Die Ursache der sechs roten `m-p16-*` im Lauf des Implementers und im Lauf des Verifiers war der volle Datenträger, kein Fehler dieses Slice und kein Vorbefund des Schritts `scan`: Der Arbeitsordner wächst rund 40 MB je Mutation (über 5 GB bei 130 bis 535 Mutationen), und der Host-Root hatte 3,8 bis 4,8 GB frei. Auf dem Host mit anderer Platte für `TMPDIR` ist derselbe Stand grün.

### Closure-Notiz

- **Was hat funktioniert:** Die Randformen standen vor dem Code in den Verträgen (`SPEC-055` Punkt 3, `SPEC-067` Punkt 4) und in §6; der Implementer gab nichts an den Architect zurück. Review und Verifikation fanden kein HIGH und kein MEDIUM; die Mutationen der Zusagen waren in der Stichprobe je rot, die unveränderten Kopien grün, und das Handbuch stimmte mit einem echten Bau überein (Namensblock gleich `ls -A`). Endstand: Summe 217 in der Metrik von §4 (Marke 800, Schätzung 360 bis 510 unterschritten), zwei Liefer-Punkte, zwei Schichten.
- **Was ging anders als geplant:** (1) Das Review fand die Zusage zu `run-container.sh` ohne benannte Mutation (F-1); behoben mit `p5-name-altes-schema` in `c1326b0c`. (2) `make release-gegenprobe` lief weder beim Implementer noch beim Verifier in einem Zug: der Arbeitsordner füllte den Datenträger, sechs `m-p16-*` meldeten Rot, allein gefahren grün; der Implementer nannte es einen Vorbefund der Umgebung ohne Ursache, die Verifikation konnte es nicht bestätigen (F-V1, F-V2). Der Lauf in einem Zug auf dem Host mit anderer Platte (oben) beendet die Frage. (3) Zwei Zahlen waren am Endstand um eins zu niedrig (F-V3: sieben statt sechs Fundstellen in `container-gegenprobe.sh`; F-V4: 217 statt 216); in §7 nachgezogen.
- **Steering-Loop-Eintrag:** benannte Lücke: Eine Gegenprobe, die je Mutation Arbeitsplatz im Temp-Verzeichnis belegt und ihn erst am Ende des Laufs freigibt, scheitert auf einem Host mit wenig freiem Platz als scheinbarer Befund, und ein Teil-Lauf ersetzt den Lauf in einem Zug nicht; dazu ist die Laufzeit der beiden größten Gegenproben so lang, dass der Lauf in einem Zug gescheut wird. Die Lücke hat eine Adresse: `slice-harness-gegenproben-schneller` (`open/`; DoD 3 entfernt den Arbeitsordner je Mutation, DoD 1 und 2 senken die Laufzeit; §1 dort nennt die Sendung und die Vermutung des Nutzers, die Ursache ist nicht belegt). Kein `liegt in`: mit diesem Slice wurde nichts verkörpert. **Empfehlung des Planners, keine Entscheidung:** Ob §3.12 (Umgebung einer Gegenprobe) um „freier Platz“ ergänzt wird, entscheidet der Nutzer mit der Zuweisung von `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (9×).
- **Beobachtungs-Register (`../observations/`):** Belege `evidence/slice-v1-abschluss-binary-namen.md` in `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts/` (F-2, F-V1, F-V2; 9×, Größe freier Platz), `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht/` (Risiko 2, weiter offen; 2×) und `BEO-REPO/test-faellt-unter-cpu-last-aus/` (`TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` unter Last, allein 1,19 s; **3×, Schwelle erreicht**, Ausgang beim Nutzer und beim Lese-Schritt der Closure von `welle-v1-nachtrag`). Kein Auftreten und Gegenbelege ohne Evidence: `BEO-REPO/verhalten-nur-unter-linux-geprueft` (Randform (h) vor dem Code entschieden, Grenze an der Zusage in `SPEC-055` Punkt 13 (4)); `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (die zwei Sendungen und die Bindung stehen in `slice-v1-abschluss-homebrew` und `slice-erster-release-homebrew-nachweis`, `grep` findet die Kennung des Gebers dort je dreimal; der Nehmer des Abschnitts zu `v0.1.0`, `slice-v1-abschluss-release-v0-1-1`, nennt sie in §1, DoD 2 und §8).
- **Folge-Slices:** `slice-harness-gegenproben-schneller` (`open/`, ohne Welle; angelegt in diesem Slice, Lerneintrag). Er ist eine Adresse, die annimmt: §1 dort nennt Arbeitsordner und Laufzeit, kein Ausschluss trifft die Sendung, er liegt nicht in `done/`; die Zählung steht in seinem §8 (drei Liefer-Punkte, eine Schicht). Für den Abschnitt zu `v0.1.0` in Handbuch und `README.md`: `slice-v1-abschluss-release-v0-1-1` (nimmt in §1, DoD 2 und §8 an, nicht in `done/`).
- **Risiken aus §6:** *Fundstellen des alten Schemas* — entfallen (Begründung in §6: jeder Treffer am Endstand von Review und Verifikation gelesen, keiner ein gültiger Name). *Abschnitt zu `v0.1.0` in Handbuch und `README.md`* — weiter offen → `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (Beleg eingetragen, 2×), Adresse `slice-v1-abschluss-release-v0-1-1`.
- **Findings:** Review F-1 behoben (`c1326b0c`), F-2 (Ursache der `m-p16-*`: voller Datenträger, Beleg oben; Adresse `slice-harness-gegenproben-schneller`) und F-3 (Mitläufer der Mutationen offen mit „+“ geführt, kein Fall je Zeichenklasse nötig, weil `p3-namen` den Namen gegen eine feste Liste prüft) INFO ohne Handlung. Verifikation F-V1 gegenstandslos durch den Lauf in einem Zug (oben), F-V2 siehe Register und Folge-Slice, F-V3 und F-V4 in §7 nachgezogen.
- **Drei Paarungen:** Anker: Der Lerneintrag verkörpert nichts (`liegt in` entfällt, die Lücke ist benannt). Folge-Slice: `slice-harness-gegenproben-schneller` und `slice-v1-abschluss-release-v0-1-1` existieren, nehmen an (oben), nicht in `done/`. Register: drei Einträge fortgeschrieben, je mit nicht leerem `evidence/`.
- **Welle:** Der Slice trägt `Welle: welle-v1-nachtrag`; er bleibt als Volltext in `done/`, bis die Welle schließt und `make archive-welle` ihn archiviert; die Roadmap bleibt unverändert (keine Umplanung).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen. Treffer: `verhalten-nur-unter-linux-geprueft` (Plattform; hier als Randform (h) vor dem Code entschieden, `AGENTS.md` §3.12), `folge-slice-adresse-nimmt-nicht-an` (die zwei Sendungen in §1 stehen in den Nehmern, §3.13), `folge-slice-liefert-handbuch-teil-nicht` (der Abschnitt zu `v0.1.0` hat keinen Nehmer mit Kennung, §6 und Auftrag an den Planner). Der Stand der Zähler dieser drei Einträge ist beim Start zu lesen; erreicht einer mit diesem Slice drei, braucht er einen Folge-Slice.

**Zählung** (Architect, 2026-10-20, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 2** (1: Namen in den drei Skripten, ihren Gegenproben und der Zeile von `harness/README.md`; 2: Handbuch und `releasing.md`). **Schichten: zwei** — Harness-Werkzeug (Skripte, Gegenproben, Zeile in `harness/README.md`) und Nutzer- und Wartungs-Doku (`docs/user/`, `docs/maintainer/`); Produkt-Code zählt nicht (§1). Die Verträge liegen vor dem Code beim Architect und zählen nicht als Liefer-Punkt. Größe: Schätzung und Marke in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält Spezifikation und Entscheidungen vor dem Code; Verträge und [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) gehen den Skripten voraus).
