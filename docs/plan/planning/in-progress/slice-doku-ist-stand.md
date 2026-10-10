# Slice slice-doku-ist-stand: Benutzerhandbuch und README im Ist-Zustand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er
von der nächsten Welle-Closure. Vorgezogen nach Entscheidung des Nutzers vom 2026-10-10: direkt
nach `slice-v1-abschluss-einspielen-laufsteuerung` und vor `slice-harness-upgrade-v6-18`
(WIP-Limit 1); welle-v1-abschluss geht danach in der Reihenfolge ihres §5 weiter.

**Bezug:** [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) (beschriebenes Verhalten, nicht geändert). Keine ADR: Die Regel ist eine Entscheidung des Nutzers vom 2026-10-10 und steht in `AGENTS.md` §3.11.

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte
Verhalten, das Zielbild bleibt in `spec/`)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben nur das Verhalten, das
das gebaute Binary heute zeigt: jede Option, jeder Abschnitt und jedes Beispiel, das das Binary
nicht kennt, ist entfernt oder auf Geliefertes gekürzt, ohne Chronik, ohne Zielstand und im
Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews
(`AGENTS.md` §3.11, *Handbuch und README beschreiben den Ist-Zustand*, seit
slice-v1-abschluss-einspielen-laufsteuerung).

**Herkunft:** Entscheidung des Nutzers vom 2026-10-10 bei der Closure von
`slice-v1-abschluss-einspielen-laufsteuerung`: „Im Handbuch immer nur den Ist-Zustand
beschreiben. Keine Chronik oder Verweise in die Spec, ADRs, Slices, Reviews etc.“ und „Bitte
auch im README.md den Ist-Zustand beschreiben. In ./spec liegt das Zielbild.“ Anlass ist
Verifikation V-147 jenes Slice: Der Hinweis zu `--compare-responses` in Handbuch §4 sagt die
Wirkung zweier Optionen mit Vergleich zu, `--compare-responses` ist bei `play` unbekannt
(`PGR-E2001`). Dieselbe Lage haben nach V-147 im selben Abschnitt die Hinweise zu
`--keep-timing`, `--upstream-tls`, `--upstream-ca` und das Passwort in *Vorgehen*. Das Urteil
zu F-553 in `slice-v1-abschluss-einspielen` (dort §7: Zielstand in Handbuch §4 bis zum Ende der
Welle hinnehmbar) ist damit überholt. Im README beschreibt die Einleitung `play` mit „vergleicht
auf Wunsch“.

**Aus `slice-v1-abschluss-einspielen-laufsteuerung`:** dessen Befund V-147 (Hinweis zu
`--compare-responses` in Handbuch §4) wird hier entfernt, nicht vom Vergleich nachgeliefert
(dort §1 und §7).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der
  Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md`, die Sendungen und
  Schichtzählungen in den Plänen unter `next/` (§7, §3.13), die zehn Doku-Folge-Slices dort mit
  `welle-v1-abschluss` und der Roadmap (§7, Entscheidung des Nutzers vom 2026-10-10), `AGENTS.md`
  §3.11 und §3.13, das Register und diesen Plan. Zeigt das
  Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary
  (§6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner, keine
  Korrektur am Code hier.
- Die Abschnitte des Handbuchs, die erst die Folge-Slices liefern (Vergleich, Zeitangaben,
  Extended beim Einspielen, Anmeldung, TLS, SQLite-Format und die übrigen) — ein Folge-Slice
  übernimmt sie, der Code-Slice in seiner DoD oder, wo die Schichtzählung das nicht zulässt, sein
  Doku-Folge-Slice `<Kennung>-doku` direkt dahinter (§7; `AGENTS.md` §3.11, §3.13). Dieser
  Slice entfernt nur, was heute nicht stimmt.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md`
  — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des
  Binaries; die Tabellen schreibt `make abdeckung`.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft (etwa die Optionen gegen
  `--help`) — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der
  Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.
- Die Abschnitte von `README.md` für Mitwirkende (Verträge, Auditierbarkeit,
  Rückverfolgbarkeit) — Bestand bleibt bewusst stehen, soweit er ohne Chronik ist: `README.md`
  darf auf `spec/`, `docs/plan/` und `docs/reviews/` zeigen (`AGENTS.md` §3.11; Rang 7 zeigt
  nach oben, `harness/README.md` §Source precedence). Gekürzt wird dort nur, was einen Stand
  der Umsetzung erzählt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Das Benutzerhandbuch nennt nur Kommandos, Optionen, Umgebungsvariablen, Schlüssel der
      Konfigurationsdatei, Meldungscodes und Exit-Codes, die das gebaute Binary kennt: je
      Kommando sind die Optionen gegen `--help` des Binaries gehalten (Liste in §7), jeder
      Hinweis, der eine unbekannte Option beschreibt (darunter `--compare-responses`,
      `--keep-timing`, `--upstream-tls`, `--upstream-ca` und das Passwort bei `play`), ist
      entfernt oder auf Geliefertes gekürzt; jede Codezeile in §7 *Fehlercodes* und
      *Warnungen* ist durch einen Test oder eine Probe belegt oder entfernt (Beleg in §7 je Abschnitt).
- [x] Jedes Beispiel des Handbuchs (Aufrufe, Konfigurationsdatei, Ausgaben) läuft als Datei
      gegen das gebaute Binary mit dem beschriebenen Ergebnis (`AGENTS.md` §3.11); die
      Kopfzeilen (*Software-Version*, *Stand*, *Gültigkeitsbereich*) und §11
      *Änderungshistorie* nennen den Ist-Zustand ohne Chronik, §1 *Voraussetzungen* und §2
      *Installation* nur, was ein Build aus dem Repository liefert (§6); das Handbuch enthält keinen
      Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch
      nicht“, „kommt“, „geplant“ über das Produkt (`grep` in §7).
- [x] `README.md` beschreibt in Leitsatz, Einleitung, *Was kann ich heute tun?* und
      *Kerngedanke* nur Geliefertes, ohne Chronik („steht am Beginn der Umsetzung“, „noch
      nicht“, „erste Version“), ohne „vergleicht auf Wunsch“ und ohne TLS zum Client; Verweise
      auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (§6).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | alle Abschnitte gegen das gebaute Binary: Optionen gegen `--help`, Beispiele als Datei, Codes, Kopfzeilen, §11; Unbekanntes entfernt oder gekürzt |
| `README.md` | update | Einleitung und *Was kann ich heute tun?* auf Geliefertes, ohne Chronik |
| Pläne unter `docs/plan/planning/next/` (Nehmer, Liste in §7) | update | Sendung der Ist-Zustand-Stellen in §1, DoD und §8 des Nehmers, mit der Kennung dieses Slice (`AGENTS.md` §3.13); §8 zählt die Schichten nach der einheitlichen Zählung |
| zehn neue Pläne `slice-v1-abschluss-<Kennung>-doku` unter `next/`, `docs/plan/planning/welle-v1-abschluss.md` (§4, §5), Roadmap (Drift-Log) | neu / update | Doku-Folge-Slices direkt hinter ihrem Code-Slice (Entscheidung des Nutzers vom 2026-10-10 zu V-151) |
| `AGENTS.md` §3.11, §3.13 | update | Schichtzählung für alle Pläne; Handbuch- und README-Teil im selben Slice oder im Doku-Folge-Slice |
| `docs/plan/planning/observations/BEO-REPO/` | update | Belege der Klassen dieses Slice, neuer Eintrag `folge-slice-liefert-handbuch-teil-nicht` |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos
  (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image
  aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach
  entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Abschnitt für Abschnitt: was bleibt, was gekürzt, was entfernt wird, in §7 je Abschnitt eine
  Zeile.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen-laufsteuerung` liegt in
`done/` (WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die
Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar. Schnitt dann: README und Handbuch §1 bis §4 hier, §5 bis §11 als
  eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): Eine Aussage lässt sich nur mit einer
  Änderung am Binary wahr machen, und der Nutzer will sie nicht streichen; dann zuerst die
  Entscheidung des Nutzers, welcher Slice das Verhalten liefert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, die Proben und `grep`-Läufe in §7, Closure-Notiz mit
Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — vom Architect vor dem ersten Commit entschieden
(2026-10-10). **Ort ist dieser Abschnitt, nicht die Spezifikation:** Der Slice legt keinen
Vertrag an (keine Option, kein Format, kein Gate); er wendet die Regel aus `AGENTS.md` §3.11 auf
zwei Dokumente der Ränge 6 und 7 an. Eine Spezifikationsstelle über das Handbuch zeigte aus dem
Technik-Stratum nach unten und machte eine Schreibregel zu Produktverhalten. Akzeptiertes
Negativ: Die Entscheidungen gehen mit diesem Plan nach `done/`; wiederkehrend ist nur die Form
der Grenz-Sätze, und die steht danach im Handbuch selbst, wo der liefernde Folge-Slice sie
findet und ersetzt.

- **Software-Version** — die Zeile nennt die Ausgabe von `pgwire-recorder version` des gebauten
  Binaries, heute `dev`; *Gültigkeitsbereich* nennt dasselbe Binary, gebaut aus dem Repository.
  *Stand* ist das Datum der Änderung, *Version* (Handbuch) steigt auf `0.2`. §11 trägt den Satz
  „Es gibt keine veröffentlichte Version.“ — Ist-Aussage ohne „noch“; je veröffentlichter
  Version kommt die Zeile mit dem Release, das sie liefert.
- **Installation und Voraussetzungen** (in der Liste nachgetragen, Handbuch §1 und §2 nennen
  Releases, Homebrew-Tap und Registry-Images, die es nicht gibt) — §2 beschreibt den einzigen
  Weg, den es gibt: Bauen aus dem Repository mit `make build` (Docker und GNU `make`), Aufruf
  über das Image `pgwire-recorder:dev` und, als Linux-Programm für die Architektur des bauenden
  Rechners, das aus dem Image kopierte Binary. Binaries, Homebrew, Registry-Images und das
  Compose-Beispiel mit `ghcr.io` entfallen; ein Compose-Beispiel, falls es bleibt, nutzt
  `pgwire-recorder:dev` und läuft als Datei. §1 *Voraussetzungen* nennt keine Plattform, die die
  Probe nicht zeigt (kein macOS, kein Windows). Die Abschnitte liefern
  `slice-v1-abschluss-homebrew`, `slice-v1-abschluss-container`,
  `slice-erster-release-veroeffentlichung` und `slice-erster-release-homebrew-nachweis`, je mit
  ihrer Zeile zum Ist-Zustand in DoD und §3.
- **Teilweise geliefert** — das Handbuch nennt eine Grenze als Eigenschaft des Binaries im
  Präsens, wenn das Binary an ihr beobachtbar reagiert (Meldungscode, Exit-Code, Ablehnung) und
  die Probe die Reaktion zeigt, etwa „Verlangt die Datenbank ein Passwort, endet `play` mit
  `PGR-E4005`“ (so sagt es `play --help` selbst) oder „`record` und `replay` lehnen eine
  Anfrage nach Verschlüsselung ab; ein Treiber mit `sslmode=prefer` verbindet sich
  unverschlüsselt“. Eine Grenze ohne beobachtbare Reaktion bleibt ungenannt. Kein Satz sagt,
  dass etwas kommt; der liefernde Slice ersetzt den Grenz-Satz.
- **Verhalten neben der Spezifikation** — das Handbuch beschreibt das Binary, belegt durch die
  Probe; Weglassen verschwiege dem Leser ein Verhalten, dem er begegnet. Je Fund eine Zeile in §7
  (Stelle im Handbuch · Verhalten des Binaries · Stelle der Spezifikation), die an den Planner
  geht; der Implementer hält dafür nicht an. Ist unklar, ob das Binary überhaupt ein stabiles
  Verhalten zeigt (etwa wechselnd zwischen Läufen), beschreibt das Handbuch die Stelle nicht und
  der Fund geht ebenso an den Planner.
- **Verweise im README** — die Verweise auf `spec/`, `docs/plan/`, `docs/reviews/`, `AGENTS.md`
  und `harness/README.md` bleiben, ohne Aussage über einen Stand: „Die Anforderungen stehen im
  Lastenheft“ statt „Die vollständige Beschreibung steht im Lastenheft“. Der Punkt *Gates* zählt
  die Ziele nicht auf, sondern zeigt auf `harness/README.md` §Sensors (der Gate-Index steht einmal dort, `AGENTS.md` §4; eine zweite Liste im README
  driftete gegen ihn). Leitsatz und *Kerngedanke* fallen unter den Ist-Zustand wie die Einleitung
  (DoD-Punkt 3).
- **Wie geprüft wird** — Optionen je Kommando gegen `--help` des gebauten Binaries (`record`,
  `replay`, `play`, `config show`); Umgebungsvariablen und Schlüssel der Konfigurationsdatei
  gegen die Regel, die `--help` nennt, und je Beispieldatei mit `config show`; jedes Beispiel als
  Datei gegen das Binary, `record` und `play` gegen das gepinnte PostgreSQL-Image (§3). Eine
  Zeile in §7 *Fehlercodes* und *Warnungen* bleibt, wenn der Code im Katalog des Binaries steht
  (`internal/hexagon/model/fehler.go`) **und** ein Test oder eine Probe ihn auslöst; löst ihn
  keine Eingabe aus und trägt ihn nur ein Unit-Test (`PGR-E1000`), bleibt die Zeile mit einem
  Ist-Satz, der das sagt, und die Zeile *1 · sonstiger Fehler* der Exit-Codes nennt ihn; ein Sammelcode der Klasse (`PGR-E2000`, `PGR-E3000`, …), den der
  Katalog nicht führt, entfällt aus der Zeile. Beleg je Abschnitt in §7.
- **Abdeckungstabellen** `docs/user/abdeckung-*.md` — nicht Teil des Handbuch-Begriffs aus
  §3.11: Sie sind Abdeckungs-Deklarationen der Tests, von `make abdeckung` erzeugt, und ihre
  Verweise auf das Lastenheft sind ihr Zweck. Unberührt (§1).
- **`docs/user/benutzerhandbuch-standard.md`** — unberührt (§1). Ein Abschnitt, den der Standard
  verlangt und für den das Binary nichts liefert, behält seine Überschrift mit einem Ist-Satz
  (etwa „Das Werkzeug kennt keine Rollen.“), statt zu entfallen.

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (etwa macOS,
  Homebrew, ein Treiber außerhalb des Containers); dann gilt es nicht als belegt
  (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** entfallen: Das Handbuch nennt nur Linux, und jedes Beispiel lief als Datei gegen das Binary (§7); der Eintrag bleibt bei 1×.
- Das Handbuch schrumpft so stark, dass Abschnitte des Standards
  (`benutzerhandbuch-standard.md`) leer stehen; jede Folge-Lieferung muss sie wieder füllen —
  **Ausgang:** entfallen: *Rollen und Rechte* (§6) und *FAQ* (§8) tragen Ist-Sätze, keine Überschrift des Standards steht leer (Verifikation, Abschnitt 1).
- Ein Folge-Slice liefert seinen Handbuch-Teil nicht, weil seine DoD die Zeile nur im ersten
  Liefer-Punkt trägt und das Review sie übersieht — **Ausgang:** weiter offen: → `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` im Register (1×; die Nehmer und ihre Doku-Folge-Slices schließen nach diesem Slice).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

**Probe-Aufbau.** `make build` am Stand `248a248` (Image `pgwire-recorder:dev`,
`sha256:b45c7f918488…`, `linux/amd64`, Benutzer `nonroot`). Proben in einem Scratch-Verzeichnis
außerhalb des Repos; Netz `impl-doku-net` mit dem gepinnten PostgreSQL-Image aus
`harness/mk/integration.mk` zweimal: `impl-doku-pg` (Alias `postgres`, Anmeldung `trust`) und
`impl-doku-pgpw` (Passwort, SCRAM). Client ist `psql` aus demselben Image in einem eigenen
Container; die Proben im Handbuch-Wortlaut („Anwendung auf `localhost:15432`“) laufen damit über
den Containernamen statt `localhost`. Container, Netze, Compose-Projekte und Scratch-Dateien sind
danach entfernt.

**Optionen je Kommando gegen `--help`** (Handbuch §5, Tabelle; §4):

| Kommando | `--help` nennt | Handbuch nennt (nachher) | entfernt (vorher im Handbuch, Probe `PGR-E2001` „flag provided but not defined“, Exit 2) |
|---|---|---|---|
| `record` | `--listen`, `--upstream`, `--output`, `--force`, `--shutdown-timeout`, `--log-level`, `--config` | dieselben | `--format`, `--record-timing`, `--record-empty-sessions`, `--tls-cert`, `--tls-key`, `--allow-plaintext` |
| `replay` | `--listen`, `--input`, `--fail-on-unconsumed`, `--shutdown-timeout`, `--log-level`, `--config` | dieselben | `--session-assignment`, `--tls-cert`, `--tls-key`, `--allow-plaintext` |
| `play` | `--upstream`, `--input`, `--user`, `--database`, `--continue-on-error`, `--allow-recorded-errors`, `--finish-session-on-interrupt`, `--log-level`, `--config` | dieselben | `--upstream-tls`, `--upstream-ca`, `--compare-responses`, `--keep-timing`, `--timing-mode`, `--timing-reference` |
| `config show` | `--config` | dasselbe | — |

Umgebungsvariablen: die Regel aus `--help` (`PGWIRE_RECORDER_` und Name der Option); das Handbuch
nennt 14, je eine zu einer Option oben. `PGWIRE_RECORDER_PASSWORD` entfernt (keine Option; `play`
mit gesetzter Variable gegen `impl-doku-pgpw`: `PGR-E4005`, Exit 4). Probe
`PGWIRE_RECORDER_FORCE=true` ersetzt, `=1` ist `PGR-E2001`. Schlüssel der Datei: `record.format`,
`record.record_timing`, `record.tls_cert`, `replay.session_assignment`, `play.compare_responses`,
`play.upstream_tls`, `play.password` je `PGR-E2004` „unbekannter Schlüssel“, Exit 2; `force: true`
ersetzt.

**Prüfung als Befehlsfolge** (aus dem Wurzelverzeichnis des Repos nach `make build`; `W` ist ein
Verzeichnis außerhalb des Repos; vergleicht jede Option des Handbuchs mit `--help`, jede Variable
mit einer Option, jeden Code mit dem Katalog `internal/hexagon/model/fehler.go`):

```bash
{ docker run --rm pgwire-recorder:dev --help
  for k in record replay play "config show"; do docker run --rm pgwire-recorder:dev $k --help; done; } > $W/help.txt 2>&1
grep -oE -- '--[a-z][a-z-]*' docs/user/benutzerhandbuch.md | sort -u > $W/o-hb
grep -oE -- '--[a-z][a-z-]*' $W/help.txt | sort -u > $W/o-help
comm -23 $W/o-hb $W/o-help          # Optionen des Handbuchs, die --help nicht kennt
grep -oE 'PGWIRE_RECORDER_[A-Z_]+' docs/user/benutzerhandbuch.md | sort -u > $W/v-hb
sed 's/^--//' $W/o-help | tr 'a-z-' 'A-Z_' | sed 's/^/PGWIRE_RECORDER_/' | sort -u > $W/v-help
comm -23 $W/v-hb $W/v-help          # Variablen ohne Option
grep -oE 'PGR-[EW][0-9]{4}' docs/user/benutzerhandbuch.md | sort -u > $W/c-hb
grep -oE 'PGR-[EW][0-9]{4}' internal/hexagon/model/fehler.go | sort -u > $W/c-kat
comm -23 $W/c-hb $W/c-kat           # Codes außerhalb des Katalogs
```

| Stand | Ergebnis |
|---|---|
| Handbuch `HEAD` (`248a248`) | rot: 13 unbekannte Optionen, 14 Variablen ohne Option, 9 Codes nicht im Katalog (`PGR-E2000`, `E2003`, `E2007`, `E3000`, `E5000`, `E5004`, `E6000`, `E6003`, `W3002`) |
| Handbuch nachher (Stand der Nacharbeit) | grün: bei den Optionen nur `--h`, `--help`, `--version` (das Handbuch nennt sie als Nicht-Optionen) und `--rm`, `--name` (`docker run`); bei den Variablen und Codes keine Ausgabe; das Handbuch nennt 14 Variablen |
| Mutante: Satz mit `--compare-responses` zurück | rot, `--compare-responses` in der ersten Ausgabe |
| Mutante: Zeile `PGWIRE_RECORDER_PASSWORD` zurück | rot, die Variable in der zweiten Ausgabe |
| Mutante: Zeile `PGR-E5004` zurück | rot, der Code in der dritten Ausgabe |

Die Befehlsfolge prüft Namen, nicht die Zuordnung zu Kommandos, Standardwerte und Exit-Codes; die
14 Datenzeilen der Tabelle in §5 sind von Hand gegen `--help` gehalten, die Exit-Codes durch die
Proben unten. Kein Optionspaar lehnt das Binary ab: `--help` nennt keine Kombinationsregel, und
`grep -rniE 'kombination|zusammen mit|nicht zusammen' internal --include=*.go` findet im Code
nichts dergleichen; darum nennt die Zeile `PGR-E2001` in Handbuch §7 keine Kombination
(Review F-571). Zu Plan §6 *Verweise im README* (F-572): `harness/README.md` §Sensors führt alle
14 Ziele von `GATE_CHECKS` (`make -pn gates`), je eine Tabellenzeile; der Satz im README
(„welche es sind und was jedes prüft, steht in …“) stimmt.

Ein Vertrag (§3.10) entsteht nicht; die Mutationen belegen nur, dass die Prüfung oben greift.

**Meldungscodes** (Handbuch §7; bleibt nur mit Eintrag im Katalog und Auslösung über das Binary):

| Code | Auslösung | Ergebnis |
|---|---|---|
| `PGR-E1000` | neun Fehlformen (Listen-Adresse, Verzeichnis, `/dev/null`, `/dev/full`, unbekannter Upstream, …) lösen ihn nicht aus; Unit-Tests `TestFailJeKlasse`, `TestFehlerFremd`, `TestFehlerGleichrangig` | bleibt als Ist-Satz „Ein Defekt des Werkzeugs; keine Eingabe löst ihn aus“; Exit-Code 1 nennt ihn (Verifikation V-149) |
| `PGR-E2001` | Probe (unbekannte Option, `version x`, `--log-level=bogus`) · `play_e2e_test.go` | bleibt, ohne Sammelcode `PGR-E2000` |
| `PGR-E2002` | Probe (vorhandene Zieldatei) · `schreiben_e2e_test.go` | bleibt |
| `PGR-E2004` | Probe (fehlende `--config`, unbekannter Schlüssel, `sslmode=require` bei `record` und `play`) · `konfiguration_e2e_test.go` | bleibt |
| `PGR-E2005` | Probe (`play --upstream ci` ohne `CI_DB_HOST`) | bleibt |
| `PGR-E2006` | Probe (Klartext-Passwort) · `konfiguration_e2e_test.go` | bleibt |
| `PGR-E3001` | Probe (fehlende Aufzeichnung, Verzeichnis als `--output`, nicht beschreibbares Verzeichnis im Compose-Beispiel: Exit 3) · `schreiben_e2e_test.go` | bleibt, ohne `PGR-E3000` |
| `PGR-E3002` | Probe (`version: 99`) | bleibt |
| `PGR-E3003` | `replay_e2e_test.go` | bleibt |
| `PGR-E3004` | Probe (`replay` mit `sessions: []`, Exit 3; `play` endet damit mit Exit 0, wie die Spezifikation es verlangt) | bleibt, „`replay` meldet den Fehler beim Start“ ergänzt |
| `PGR-E4000` | Probe (`--ulimit nofile=20:20`, 40 Verbindungen: `accept4: too many open files`, Exit 4) | bleibt |
| `PGR-E4001` | `unverbraucht_e2e_test.go` | bleibt |
| `PGR-E4002` | `verbindung_e2e_test.go`, `play_e2e_test.go` | bleibt |
| `PGR-E4003` | `extended_e2e_test.go`, `play_laufsteuerung_e2e_test.go` | bleibt |
| `PGR-E4004` | `play_e2e_test.go` | bleibt, ohne `--compare-responses` |
| `PGR-E4005` | Probe (`play` gegen `impl-doku-pgpw`: „verlangt ein Anmeldeverfahren (Code 10), das play nicht unterstützt“; `--user niemand`: SQLSTATE `28000`) · `play_e2e_test.go` | bleibt, Ursache auf Passwort und abgelehnte Anmeldung gekürzt |
| `PGR-E4006` | `herunterfahren_e2e_test.go` | bleibt |
| `PGR-E5001` | Probe (anderes SQL, Log-Zeile) · `replay_e2e_test.go` | bleibt, ohne `PGR-E5000` |
| `PGR-E5002` | `unverbraucht_e2e_test.go` | bleibt |
| `PGR-E5003` | Probe (dritte Verbindung bei zwei Sessions) | bleibt |
| `PGR-E6001` | Probe (`record` gegen `impl-doku-pgpw`: Client erhält `FATAL … [PGR-E6001]`, Aufzeichnung `sessions: []`, Exit 6) · `record_e2e_test.go` (`COPY`) · `play_e2e_test.go` (Extended in `play`, Exit 6) | bleibt, ohne `PGR-E6000`, Ursachen ergänzt |
| `PGR-E6002` | Probe (Startnachricht Protokoll 2.0 per `nc`) | bleibt |
| `PGR-W2001` | Probe (`replay` beendet, Sessions nie zugeordnet) · `replay_e2e_test.go` | bleibt |
| `PGR-W3001` | Probe (CancelRequest per `nc`) | bleibt |
| `PGR-W3003` | Probe (HTTP-Anfrage per `nc`) · `record_e2e_test.go` | bleibt |
| `PGR-E2003`, `PGR-E2007`, `PGR-E5004`, `PGR-E6003`, `PGR-W3002` | nicht im Katalog | Zeile entfernt |

**Beispiele als Datei** (aus dem geänderten Handbuch extrahiert, je Codeblock eine Datei, Aufruf
`pgwire-recorder …` als `docker run … pgwire-recorder:dev …`):

| Beispiel | Probe | Ergebnis |
|---|---|---|
| §2 `docker run --rm pgwire-recorder:dev version` | als Skript | `pgwire-recorder dev`, Exit 0 |
| §2 Binary aus dem Image | Skript in leerem Verzeichnis | kopiert, `./pgwire-recorder version` gibt `pgwire-recorder dev` aus; `file`: statisch gebunden, x86-64; läuft auch im Alpine-Container |
| §2 Compose (Wiedergabe) | `docker compose up -d` mit `recordings/test.yaml`, `psql` gegen `recorder:5432` | aufgezeichnete Antworten, `stop`: Exit 0 |
| §3 `record` | gegen `impl-doku-pg` (Alias `postgres`), `psql`, `docker stop` | `users.yaml` geschrieben, Exit 0 |
| §3 `replay` | Datenbank gestoppt, derselbe Ablauf | dieselben Antworten; eine weitere Verbindung `PGR-E5003` |
| §4 Compose (Aufzeichnen mit `stop_grace_period`) | mit Override-Datei, die den Dienst `postgres` ergänzt; `recordings` mit Modus 777 | `test.yaml` mit `SELECT 42`, Exit 0; mit Modus 755 Exit 3 (Satz zum beschreibbaren Verzeichnis) |
| §4 `play` | gegen `postgres:5432` | Exit 0 |
| §5 Konfigurationsdatei | als `.pgwire-recorder.yaml`: `config show` (Exit 0, Inhalt wie geschrieben); `record` ohne Optionen im Netz der Datenbank (`localhost:5432`); `play` ohne Optionen | `play` legt als `dev` in `myapp` an (Tabelle mit Eigentümer `dev`); `play --upstream ci` mit `CI_DB_HOST=postgres` Exit 0, ohne Variable `PGR-E2005` |
| §5 Log-Zeilen | Binary aus dem Image, `TZ=Europe/Berlin` | `time=…+02:00 level=WARN msg="…" code=PGR-W2001` und `level=ERROR msg=Fehler code=PGR-E5001 error="Replay [PGR-E5001]: …"`; im Image ebenso mit `TZ` |

**Abschnitte** (Aussage · Probe · Ergebnis):

| Abschnitt | Aussage nachher | Probe | Ergebnis |
|---|---|---|---|
| Kopf | Version 0.2, Software-Version `dev`, Stand 10.10.2026, Gültigkeit für das aus dem Repository gebaute Binary | `version` | `pgwire-recorder dev` |
| §1 Voraussetzungen | Datenbank ohne Passwort; Bauen mit Docker und `make`, Ausführen auf Linux der bauenden Architektur oder im Container; `sslmode=prefer` verbindet, `require` bricht ab | `psql sslmode=require` an `record` und `replay`: „server does not support SSL“; `prefer` an `replay` verbindet; Passwort: `PGR-E6001`/`PGR-E4005` | wie beschrieben; macOS, Windows, `arm64` entfernt |
| §2 Installation | Bauen, Binary aus dem Image, Container mit `nonroot` | Beispiele oben; `docker image inspect`: Benutzer `nonroot`, Einstiegspunkt `/pgwire-recorder` | Releases, Homebrew, Registry-Images entfernt |
| §3 Erste Schritte | unverändert | Beispiele oben | wie beschrieben |
| §4 Aufzeichnen | ohne `--format`, `--record-timing`, `--record-empty-sessions`; Passwort-Grenze `PGR-E6001`; Compose mit `pgwire-recorder:dev` | Proben oben | wie beschrieben |
| §4 Wiedergeben | ohne `--session-assignment` | Probe `PGR-E5003` | wie beschrieben |
| §4 Verschlüsselte Verbindungen | Abschnitt entfernt | `--tls-cert`, `--tls-key`, `--allow-plaintext`: `PGR-E2001` | — |
| §4 Einspielen | ohne Passwort-Schritt, ohne Zeit-, Vergleichs- und TLS-Hinweise; Grenzen Passwort (`PGR-E4005`), `sslmode=require` (`PGR-E2004`), Extended (`PGR-E6001`) | Proben oben, `play_e2e_test.go` (`TestE2EPlayZwischenstand`) | Abschlusszeile mit Zählern entfernt (Probe: `play beendet` ohne Zähler) |
| §4 Treiber | `record` und `replay` mit beiden Protokollen, ohne Verschlüsselung | `psql sslmode=prefer` | wie beschrieben |
| §5 Einstellungen | Tabelle mit 14 Zeilen (Datenzeilen, eine je Option) wie `--help`; `sslmode=require` bei `record` und `play` ungültig, unbenutzt gültig; Passwort in der URL ohne Anmeldung | Proben oben (`config show` mit `require` Exit 0) | wie beschrieben |
| §5 Exit-Codes | Zeile 5 ohne Vergleich | — | Zeilen 1–6 bleiben |
| §5 Konfigurationsdatei, leere Datei | leere Datei und Datei nur mit Kommentaren gültig und ohne Wirkung; `record:` ohne Inhalt ungültig | `config show --config` mit leerer Datei, mit einer Zeile `# nur Kommentar`, mit `record:` allein (Scratch `impl2-doku-w`, Image `sha256:b45c7f918488…`) | leer und Kommentar: Ausgabe nur der Pfad, Exit 0; `record:`: `PGR-E2004` „record: erwartet eine Abbildung“, Exit 2 (Review V-153: vorher „ebenfalls nicht“, mehrdeutig) |
| §6 Rollen | `record` ohne Passwort; `replay` nimmt jede Anmeldung an | `psql user=wer dbname=anders` mit `PGPASSWORD` an `replay`: Antworten wie aufgezeichnet | wie beschrieben |
| §7 Fehlerbehebung | Codes oben; *Die Anwendung kann sich nicht verbinden* ohne TLS-Optionen | Proben oben | wie beschrieben |
| §8 FAQ | ohne SQLite | — | — |
| §11 | „Es gibt keine veröffentlichte Version.“ | — | — |

**`grep` im Handbuch** (`-i`): `spezifikation|lastenheft|\bADR-|slice|welle|review` — vorher 0,
nachher 0 Treffer; `noch nicht|noch keine|kommt mit|geplant|veröffentlichte Version|ghcr|homebrew|macOS|windows`
— vorher 13, nachher 2: „die Zieldatei existiert noch nicht“ (Zustand der Datei, nicht des
Produkts) und der Satz in §11 aus §6.

**README** (`grep -iE "noch|erste Version|Beginn der Umsetzung|auf Wunsch|TLS"`): vorher 8 Zeilen
(„steht am Beginn der Umsetzung“, „noch nicht“, „fehlt noch“, „vergleicht auf Wunsch“ zweimal, TLS zum
Client), nachher 0. Leitsatz und *Kerngedanke* ohne Vergleich; „wartet `record` ohne Frist“
durch `--shutdown-timeout` (Standard 5 Sekunden, aus `--help`) ersetzt; „Das Benutzerhandbuch
beschreibt das Verhalten der ersten Version“ und „liegen vor“ entfernt; *Gates* zeigt auf
`harness/README.md` §Sensors statt einer Liste.

**Funde für den Planner** (Stelle im Handbuch · Verhalten des Binaries · Stelle der Spezifikation):

- §1 *Voraussetzungen*, §4 *Eine Anwendung aufzeichnen*, §6, §7 `PGR-E6001` und *Die Anwendung kann
  sich nicht verbinden* · `record` beendet die Verbindung des Clients mit `PGR-E6001`, wenn die
  Datenbank ein Anmeldeverfahren verlangt (SCRAM), die Session wird nicht aufgezeichnet, Exit 6 ·
  `spec/spezifikation.md` LH-FA-05.b (*Record*: vermittelt die Authentifizierung transparent) und
  LH-FA-17.a (*record* „vermittelt die Anmeldung des Clients“).

Keine weitere Abweichung gefunden; das Fehlen der Optionen oben ist Zielstand der Folge-Slices, keine
Abweichung.

**Sendungen an die Folge-Slices** (Review F-570, Verifikation V-151; §3.13): Der Handbuch- und
README-Teil von zehn Code-Slices liegt nach der Entscheidung des Nutzers vom 2026-10-10 in je
einem Doku-Folge-Slice `<Kennung>-doku` direkt hinter ihm (siehe *Doku-Folge-Slices* unten). Die
Sendungen gehen an den Doku-Folge-Slice; der Code-Slice nennt die Abgabe in §1 *Ausdrücklich
NICHT* und zählt seine Liefer-Punkte und Schichten ohne die Dokumentation. `grep -l
slice-doku-ist-stand docs/plan/planning/next/*.md` findet 28 Dateien: die zehn Doku-Folge-Slices
(§1 *Herkunft*, bei fünf davon auch die Sendung in der DoD), ihre zehn Code-Slices (Zeile
*Nachgezählt*), sieben weitere Nehmer (`slice-v1-abschluss-container`,
`slice-v1-abschluss-homebrew`, `slice-v1-abschluss-cancel-ohne-schluessel`,
`slice-v1-abschluss-postgres-versionen`, `slice-erster-release-veroeffentlichung`,
`slice-erster-release-homebrew-nachweis`, `slice-erster-release-freigabe`) und
`slice-harness-upgrade-v6-18`, das die Kennung nur im Satz zur Reihenfolge nennt (kein Nehmer).

| Sendung | Nehmer | Stelle |
|---|---|---|
| Installation | `slice-v1-abschluss-container` | Handbuch §2 *Container*, *Das Binary aus dem Image*, §1 Plattform-Zeile |
| Installation | `slice-erster-release-veroeffentlichung` | Handbuch §2 Release-Binaries und Registry-Images, §1 Plattform-Zeile |
| Installation | `slice-v1-abschluss-homebrew` | Handbuch §2 Homebrew-Abschnitt, soweit der Probe-Tap ihn belegt |
| Installation | `slice-erster-release-homebrew-nachweis` | Handbuch §2 Homebrew-Befehle, §1 Plattform-Zeile (macOS) |
| Software-Version | `slice-erster-release-freigabe` | Handbuch Kopf *Software-Version*, *Gültigkeitsbereich* und §11 (Punkt 4 des Auftrags; DoD-Punkt 3, §1, §3) |
| Grenz-Satz Passwort `record` | `slice-v1-abschluss-anmeldung-doku` | Handbuch §1, §4 (Aufzeichnen), §6, §7 `PGR-E6001` und *Die Anwendung kann sich nicht verbinden*; README |
| Meldungstext „das dieser Stand nicht vermittelt“ (F-573, V-152) | `slice-v1-abschluss-anmeldung` | `internal/adapters/driven/postgres/upstream.go`; §1, DoD-Punkt 2, §3 |
| Grenz-Satz Passwort `play` | `slice-v1-abschluss-einspielen-anmeldung-doku` | Handbuch §1, §4 (Einspielen), §5, §7 `PGR-E4005`; README |
| Grenz-Satz Extended `play` | `slice-v1-abschluss-einspielen-extended-doku` | Handbuch §4 (Einspielen), §7 `PGR-E6001`; README |
| Grenz-Satz TLS `play` | `slice-v1-abschluss-einspielen-tls-doku` | Handbuch §4 (Einspielen), §5 `sslmode`, §1; README |
| Grenz-Satz TLS zum Client | `slice-v1-abschluss-tls-client-doku` | Handbuch §1, §4 (Treiber), §7 *Die Anwendung kann sich nicht verbinden*; README |
| Handbuch- und README-Teil ohne Grenz-Satz | `slice-v1-abschluss-sessions-doku`, `-protokollrand-doku`, `-sqlite-format-doku`, `-zeitangaben-doku`, `-antwortvergleich-doku` | die Abschnitte, die der Code-Slice liefert |

Gelesen vor dem Eintragen: §1 *Ausdrücklich NICHT* und DoD jedes Nehmers; kein Ausschluss trifft
die Sendung, und die Doku-Folge-Slices sind neu angelegt (kein Ausschluss, `grep -n
slice-doku-ist-stand` findet die Kennung in §1 *Herkunft*, bei fünf in der DoD). Annahme geprüft
für `slice-v1-abschluss-anmeldung` (Meldungstext in §1, DoD-Punkt 2, §3) und
`slice-erster-release-freigabe` (§1, DoD-Punkt 3, §3).

**Doku-Folge-Slices** (V-151; Entscheidung des Nutzers vom 2026-10-10: „Dafür haben wir die Welle
– man kann dafür einen weiteren Slice anlegen. So bleibt alles zusammen: Code-Erweiterung und
Handbuch.“): Die Schichtzählung ist eine für alle Pläne (`AGENTS.md` §3.13, Teilung aus
`slice-harness-d-check-v0-85` §8): Nutzer- und Wartungs-Doku eine Schicht, Planung keine. Mit ihr
liegen zehn Slices über zwei Schichten, sobald ihr Handbuch- und README-Teil dazukommt:
`slice-v1-abschluss-sessions` (Kern, PGWire-Adapter), `-protokollrand` (PGWire-, Upstream-Adapter),
`-anmeldung` (Kern, PGWire-, Upstream-Adapter), `-einspielen-extended` (Kern, Upstream-Adapter),
`-einspielen-anmeldung` und `-einspielen-tls` (Upstream-, CLI-Adapter, dort schon durch einen
DoD-Punkt zum Handbuch), `-sqlite-format` (Recording-, CLI-Adapter), `-zeitangaben` (Kern,
PGWire-, Recording-, CLI-Adapter), `-antwortvergleich` (Kern, CLI-Adapter) und `-tls-client`
(PGWire-, CLI-Adapter). Jeder bekommt `<Kennung>-doku` in `next/` (aus der Vorlage, nur
Handbuch und README, höchstens drei Liefer-Punkte, Voraussetzung ist der Code-Slice, Prüfweg und
Randformen wie in §6 dieses Plans), in `welle-v1-abschluss` §4 und §5 direkt hinter dem Code-Slice
(Reihenfolge jetzt 34 Schritte), im Drift-Log der Roadmap mit Datum und Grund. `AGENTS.md` §3.11
lässt den Teil im selben Slice oder im Doku-Folge-Slice zu: Die Welle schließt erst mit beiden,
ein Release entsteht erst danach. **Zwischen Code-Slice und Doku-Slice hinkt das Handbuch dem
Binary hinterher; das ist akzeptiert, ein Release liegt nicht dazwischen.** Die Zeilen
*Nachgezählt* aller Pläne in `next/` zählen jetzt nach der einheitlichen Zählung.
`slice-v1-abschluss-container`, `-homebrew`, `-cancel-ohne-schluessel` und die drei Pläne von
welle-erster-release liegen mit der Dokumentation bei zwei Schichten und bleiben ungeschnitten.

**Funde der Zählung für den Planner** (nicht Folge der Dokumentation, schon vorher über zwei
Schichten): `slice-v1-abschluss-anmeldung` (Kern, PGWire-Adapter, Upstream-Adapter),
`slice-v1-abschluss-zeitangaben` (vier Code-Schichten) und `slice-v1-abschluss-postgres-versionen`
(Spezifikation, Entscheidungen, Harness, Dokumentation). Jeder Plan nennt es in seiner Zeile
*Nachgezählt*; der Architect prüft die Teilung vor dem ersten Code-Commit, der Planner entscheidet
über den Schnitt. Der Handbuch-Teil, den `slice-v1-abschluss-einspielen-anmeldung` und
`-einspielen-tls` aus Befund F-533 trugen, liegt jetzt in ihrem Doku-Folge-Slice: Die Kennung
bleibt Adresse einen Hop länger (Code-Slice → Doku-Folge-Slice, jener nennt sie in §1).

**Nachzug zu Review und Verifikation:** F-570 bis F-572, F-575 (`5f4dfa0`), V-148, V-150, V-153
(`33de11d`) erledigt. V-149 / F-574 (`a35ea51`): Die Zeile `PGR-E1000` steht wieder in Handbuch
§7, als Ist-Satz („Ein Defekt des Werkzeugs; keine Eingabe löst ihn aus“), die Exit-Code-Zeile
nennt ihn; die Regel in §6 lautet „Test oder Probe“ (neun Fehlformen lösen ihn nicht aus,
`TestFailJeKlasse`, `TestFehlerFremd`, `TestFehlerGleichrangig` tragen ihn). V-151: siehe
*Doku-Folge-Slices*. V-152 / F-573 (bekannter Teilstand, keine neue Abweichung): Der Hinweis zum
Meldungstext ist bei `slice-v1-abschluss-anmeldung` eingetragen (`ae19ada`); Handbuch §11 und der
Kopf *Software-Version* gehören `slice-erster-release-freigabe`, der sie annimmt.

**Größe:** Diff des Handbuchs +102/−214 Zeilen, README +26/−33; überwiegend Streichungen, in einer
Review-Sitzung prüfbar — die Rückführung aus §4 greift nicht.

**Gates:** `make gates` grün am Stand `94f26b0` auf sauberem Baum (Exit 0; `docs-check` 0 Befunde,
`run-integration-tests`, `a-check-negativ`, `lint-gegenprobe`, `kopf-check-gegenprobe`,
`abdeckung-gegenprobe`, `commit-msg-gegenprobe` je grün, `baseline-verify` OK). Nach der Nacharbeit zum Review (F-570 bis F-575): `make gates` grün am Stand `5f4dfa0` auf
sauberem Baum (Exit 0, alle Gates wie oben). Nach der Nacharbeit zur Verifikation (V-148, V-150, V-153): `make gates` grün am Stand `33de11d` auf
sauberem Baum (Exit 0). Kein weiterer
Sensor betroffen: Der Slice ändert weder Code noch Tests noch Gates.

**Ergebnis der Closure**

- **Gegenstand:** geliefert (Handbuch und README im Ist-Zustand), kein Übergabe-Fall.
- **Was hat funktioniert:** Die Befehlsfolge (Optionen gegen `--help`, Variablen gegen Optionen,
  Codes gegen den Katalog) fand 13 unbekannte Optionen, 14 Variablen ohne Option und 9 Codes
  außerhalb des Katalogs im Handbuch; jede Mutante machte sie rot. Die Verifikation fand mit
  eigener Probe keine falsche Aussage mehr.
- **Was ging anders als geplant:** Die Regel „Verhalten im selben Zug mit Handbuch und README“
  (`AGENTS.md` §3.11) stieß an die Größenregel: Mit der Dokumentation als Schicht lagen zehn Slices
  über zwei Schichten, und die Pläne zählten verschieden (V-151, `BEO-REPO/schichtteilung-je-plan-verschieden`).
  Der Nutzer entschied den Doku-Folge-Slice (siehe oben). Die Sendungen mussten nach dem Review
  (F-570) und nach der Entscheidung zweimal umgehängt werden.
- **Steering-Loop-Eintrag (geschärfte Regel):** Das Handbuch eines Slice liegt in ihm oder im
  Doku-Folge-Slice direkt dahinter, und alle Pläne zählen die Dokumentation als eine Schicht und
  die Planung gar nicht — liegt in `AGENTS.md §3.13` und `AGENTS.md §3.11`, Herkunfts-Anker
  `seit slice-doku-ist-stand`. Auslöser: `BEO-REPO/schichtteilung-je-plan-verschieden`
  (slice-harness-d-check-v0-85, slice-v1-abschluss-einspielen-laufsteuerung, slice-doku-ist-stand
  — 3×). **Akzeptiertes Negativ, ausdrücklich:** Zwischen Code-Slice und Doku-Folge-Slice hinkt das
  Handbuch dem Binary hinterher; die Welle schließt erst mit beiden, ein Release liegt nicht
  dazwischen.
- **Beobachtungs-Register:** `evidence/slice-doku-ist-stand.md` in
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (7×, verkörpert; F-570),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (30×, verkörpert; F-571, F-572, V-149,
  V-150), `BEO-REPO/plan-folgt-korrektur-nicht` (24×, verkörpert; V-148) und
  `BEO-REPO/schichtteilung-je-plan-verschieden` (3×, Ausgang verkörpert, Zielort
  `AGENTS.md §3.13`; V-151); neu `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen;
  Risiko 3). `BEO-REPO/verhalten-nur-unter-linux-geprueft` bleibt bei 1×. Kein nicht verkörperter
  Eintrag erreicht 3×.
- **Folge-Slices:** die zehn Doku-Folge-Slices `slice-v1-abschluss-sessions-doku`,
  `-protokollrand-doku`, `-anmeldung-doku`, `-einspielen-extended-doku`,
  `-einspielen-anmeldung-doku`, `-einspielen-tls-doku`, `-sqlite-format-doku`,
  `-zeitangaben-doku`, `-antwortvergleich-doku`, `-tls-client-doku` (alle in `next/`, `grep -n
  slice-doku-ist-stand` findet die Kennung in jedem), dazu die Nehmer der Sendungen oben.
- **Risiken aus §6:** Risiko 1 und 2 entfallen (mit Grund in §6), Risiko 3 weiter offen →
  `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`.
- **Drei Paarungen:** (a) Anker: `AGENTS.md` §3.11 und §3.13 tragen `seit slice-doku-ist-stand`.
  (b) Folge-Slices: jede genannte Datei liegt in `next/`, und die Kennung dieses Slice steht in
  §1 oder DoD des Nehmers (`grep -n`). (c) Register: jede genannte Kennung ist ein Verzeichnis mit
  nicht leerem `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `d652cda` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (29×, `AGENTS.md` §3.11) — verkörpert;
  genau der Gegenstand dieses Slice, für das Handbuch seit
  slice-v1-abschluss-upstream-verbinden (Beispiel als Datei), für den Ist-Zustand seit
  slice-v1-abschluss-einspielen-laufsteuerung.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×) — Proben laufen nur unter Linux im
  Container (§6 Risiko); bleibt unter der Schwelle.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (6×, §3.13) — verkörpert; die Folge-Slices
  nennen die Regel mit der Kennung von `slice-v1-abschluss-einspielen-laufsteuerung`.

- `BEO-REPO/schichtteilung-je-plan-verschieden` (2×, mit diesem Slice 3×, Ausgang verkörpert in
  `AGENTS.md` §3.13) und `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (neu, 1×, offen) —
  beide aus V-151 und Risiko 3, Ausgänge in §7.

Liefer-Punkte: drei. Schichten nach der Zählung in `AGENTS.md` §3.13: zwei (Nutzer-Doku:
Handbuch und README; Harness: `AGENTS.md`), kein Produkt-Code; Pläne, Welle, Roadmap und Register
zählen nicht.

Außer `BEO-REPO/schichtteilung-je-plan-verschieden` erreicht keiner der Einträge mit diesem Plan die Schwelle 3× neu; jener ist mit dieser Closure verkörpert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
