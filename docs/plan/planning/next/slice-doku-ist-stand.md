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
  Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` (und diesen Plan). Zeigt das
  Binary ein Verhalten, das der Spezifikation widerspricht, beschreibt das Handbuch das
  Verhalten des Binaries nicht; der Fund geht als Befund an den Planner, keine Korrektur am
  Code hier.
- Die Abschnitte des Handbuchs, die erst die Folge-Slices liefern (Vergleich, Zeitangaben,
  Extended beim Einspielen, Anmeldung, TLS, SQLite-Format und die übrigen) — ein Folge-Slice
  übernimmt sie: Jeder trägt in seiner DoD, dass er seinen Teil von Handbuch und README im
  Ist-Zustand liefert (seit slice-v1-abschluss-einspielen-laufsteuerung, je §2 und §3). Dieser
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

- [ ] Das Benutzerhandbuch nennt nur Kommandos, Optionen, Umgebungsvariablen, Schlüssel der
      Konfigurationsdatei, Meldungscodes und Exit-Codes, die das gebaute Binary kennt: je
      Kommando sind die Optionen gegen `--help` des Binaries gehalten (Liste in §7), jeder
      Hinweis, der eine unbekannte Option beschreibt (darunter `--compare-responses`,
      `--keep-timing`, `--upstream-tls`, `--upstream-ca` und das Passwort bei `play`), ist
      entfernt oder auf Geliefertes gekürzt; jede Codezeile in §7 *Fehlercodes* und
      *Warnungen* ist im Binary erreichbar oder entfernt (Beleg in §7 je Abschnitt).
- [ ] Jedes Beispiel des Handbuchs (Aufrufe, Konfigurationsdatei, Ausgaben) läuft als Datei
      gegen das gebaute Binary mit dem beschriebenen Ergebnis (`AGENTS.md` §3.11); die
      Kopfzeilen (*Software-Version*, *Stand*, *Gültigkeitsbereich*) und §11
      *Änderungshistorie* nennen den Ist-Zustand ohne Chronik; das Handbuch enthält keinen
      Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch
      nicht“, „kommt“, „geplant“ über das Produkt (`grep` in §7).
- [ ] `README.md` beschreibt in Einleitung und *Was kann ich heute tun?* nur Geliefertes, ohne
      Chronik („steht am Beginn der Umsetzung“, „noch nicht“) und ohne „vergleicht auf
      Wunsch“; Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | alle Abschnitte gegen das gebaute Binary: Optionen gegen `--help`, Beispiele als Datei, Codes, Kopfzeilen, §11; Unbekanntes entfernt oder gekürzt |
| `README.md` | update | Einleitung und *Was kann ich heute tun?* auf Geliefertes, ohne Chronik |

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

**Randformen** (`AGENTS.md` §3.12) — **offen, der Architect entscheidet sie vor dem ersten
Commit** (§4). Der Slice legt keinen Vertrag an; Ort der Entscheidung ist dieser Abschnitt.

- **Software-Version** — die Zeile sagt heute „noch nicht veröffentlicht“; zu entscheiden, was
  im Ist-Zustand dort steht (etwa die Ausgabe von `pgwire-recorder version` des gebauten
  Binaries oder die Versionsdatei), solange kein Release vorliegt, und ob §11 dann leer bleibt.
- **Teilweise geliefert** — ein Kommando, das nur einen Teil kann (etwa `play` ohne
  vorbereitete Anweisungen, ohne Passwort-Anmeldung und ohne TLS): ob das Handbuch die Grenze
  als Eigenschaft des Binaries nennt („spielt einfache Anfragen ein“) oder schweigt; „noch
  nicht“ ist ausgeschlossen.
- **Verhalten neben der Spezifikation** — zeigt das Binary ein Verhalten, das die
  Spezifikation anders regelt: Das Handbuch beschreibt das Binary oder lässt die Stelle weg;
  welches von beiden, entscheidet der Architect je Fund, der Fund geht an den Planner (§1).
- **Verweise im README** — welche Verweise auf `spec/`, `docs/plan/` und `docs/reviews/`
  bleiben und in welcher Form (§1, `AGENTS.md` §3.11).

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (etwa macOS,
  Homebrew, ein Treiber außerhalb des Containers); dann gilt es nicht als belegt
  (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** offen bis Closure.
- Das Handbuch schrumpft so stark, dass Abschnitte des Standards
  (`benutzerhandbuch-standard.md`) leer stehen; jede Folge-Lieferung muss sie wieder füllen —
  **Ausgang:** offen bis Closure.
- Ein Folge-Slice liefert seinen Handbuch-Teil nicht, weil seine DoD die Zeile nur im ersten
  Liefer-Punkt trägt und das Review sie übersieht — **Ausgang:** offen bis Closure.

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

Liefer-Punkte: drei. Schichten: eine (Dokumentation), kein Code.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
