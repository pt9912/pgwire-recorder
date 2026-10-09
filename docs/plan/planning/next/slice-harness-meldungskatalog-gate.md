# Slice slice-harness-meldungskatalog-gate: Gate für Code-Tabelle und Meldungskatalog

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein
Abnahmeszenario und kein Meilenstein hängt an ihm. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom 2026-10-07:
in der wellenlosen Reihe nach `slice-harness-integration-wait` und vor
`slice-harness-abdeckung-gate` (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler) (stabile Meldungscodes je Ursache). Bindung an Entscheidungen: [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) (Folgepflicht: ein Gate, das Code-Tabelle und Katalog abgleicht).

**Berührte Spec-Stellen:** `SPEC-034` (gelesen, Tabelle der Codes; geändert nur, wenn der Architect einen Bestands-Befund dort berichtigt) · `spezifikation.md §11` (neue Kennung für den Vertrag des Gates, vergeben vom Architect vor dem Code)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Gate meldet rot, wenn ein Meldungscode `PGR-…` im Go-Code nicht im
Meldungskatalog der Spezifikation (Tabelle in `SPEC-034`) steht, und umgekehrt, wenn
ein Code des Katalogs im Go-Code fehlt. Was als Code im Go-Code zählt, was als
Katalog-Zeile und wie „umgekehrt“ für Codes gilt, die der Code noch nicht erzeugt,
entscheidet der Architect vor dem ersten Code-Commit im Vertrag des Gates
(`spezifikation.md §11`); die Randformen dazu stehen in §6. Geliefert wie
`make kopf-check` und `make abdeckung-check`: ein Prüfskript unter `tools/harness/`
mit Gate-Ziel, eine Gegenprobe als eigenes Gate-Ziel, beide über ein Fragment unter
`harness/mk/` an `GATE_CHECKS`, eine Sensor-Datei unter `harness/sensors/` und je eine
Zeile in `harness/README.md` §Sensors.

**Herkunft:** [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) nennt als
Folgepflicht ein Gate, das Code-Tabelle und Katalog abgleicht, sobald beide
existieren; seit `slice-replay-semantik-meldungscodes` liegt die Code-Tabelle in
`internal/hexagon/model/fehler.go`. Die Closure von welle-replay-semantik fand, dass die
Folgepflicht keine Slice-Kennung hat (Nebenbefund 4 in
`welle-replay-semantik-results.md`): welle-v1-abschluss §6 und
`slice-v1-abschluss-container` §1 nannten das Gate als „Folge-Slice“ ohne Kennung.
Entscheidung des Nutzers vom 2026-10-07 (Audit nach [`AGENTS.md`](../../../../AGENTS.md)
§3.13, Fund 1): eigener wellenloser Slice, verglichen wird der Go-Code mit dem Katalog
der Spezifikation.

**Übernommen aus welle-v1-abschluss** (dort §6, Out-of-Scope) **und aus
`slice-v1-abschluss-container`** (dort §1, Abgrenzung): das Gate für Code-Tabelle und
Katalog; das ist das Ziel oben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Codes erzeugen, die der Katalog führt und der Code noch nicht — Schicht-Abgrenzung:
  kein Verhalten im Produkt-Code; an `internal/hexagon/model/fehler.go` ändert der
  Slice höchstens die Code-Tabelle, und nur nach §6 *Bestand*. Die Codes der offenen Anforderungen liefern die Slices von
  welle-v1-abschluss, die sie umsetzen; wie das Gate bis dahin mit ihnen umgeht, ist
  eine Randform (§6), keine Arbeit am Produkt.
- Klasse, Exit Code und Form des Fehlertexts eines Codes prüfen — Bestand bleibt
  stehen: Das prüfen die Tests aus `slice-replay-semantik-meldungscodes`
  (`fehler_test.go`, Kopf und Exit Code in `cli_test.go`, `bootstrap_test.go`); das
  Gate gleicht Mengen von Codes ab, kein Verhalten.
- Codes, die andere Abschnitte der Spezifikation, das Lastenheft oder Slice-Pläne in
  Prosa nennen — anderer Vorgang: Der Abgleich gilt der Code-Tabelle und dem Katalog,
  nicht jeder Nennung; eine Prosa-Nennung eines unbekannten Codes ist Urteil des
  Reviews.
- Den Katalog neu ordnen oder Codes neu vergeben — anderer Vorgang: Vergabe ist
  Spezifikationsarbeit eines Slice, der eine Ursache einführt (`SPEC-034`
  §Stabilität); dieser Slice berichtigt nur einen Bestands-Befund, den sein Vertrag
  rot macht, und nur nach Entscheidung des Architects (§6, *Bestand*).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler):
      Der Vertrag des Gates steht vor dem ersten Code-Commit in `spezifikation.md §11`
      unter einer neuen Kennung und entscheidet jede Randform aus §6; das Prüfskript
      (bash, ohne Docker) hängt mit seinem Ziel über ein Fragment unter `harness/mk/`
      an `GATE_CHECKS` und ist am Stand des Starts grün, ein Befund des Bestands ist nach
      Entscheidung des Architects berichtigt (§6, *Bestand*).
- [ ] Gegenprobe als eigenes Gate-Ziel in Temp-Bäumen: je Randform aus §6 ein Fall, der
      abgelehnt bzw. angenommen werden muss — mindestens ein Code nur im Go-Code, ein
      Code nur im Katalog, ein Code in falscher Form; dazu in einem Temp-Baum mit
      Makefile und Fragment, dass beide Ziele an `GATE_CHECKS` hängen und über `make`
      bei rotem Skript mit Exit ungleich 0 enden. Je Zusage des Skripts ist die Mutation
      vom Implementer rot gesehen ([`AGENTS.md`](../../../../AGENTS.md) §3.10); die
      Tabelle Zusage · Mutation · roter Test steht in §7.
- [ ] Sensor-Datei `harness/sensors/meldungskatalog-check.md` (Vertrag, Grenze, Ausgabe
      und Ausgänge, wie `harness/sensors/kopf-check.md`) und je eine Zeile für beide
      Ziele in `harness/README.md` §Sensors mit Bindung an
      [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) und den neuen Vertrag.
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
| `spec/spezifikation.md` §11 (neue Kennung, Historie) | update (Architect, vor Code) | Vertrag des Gates: Gegenstand, Fundort, Richtung, Form, Ausgabe, Grenze; je Randform aus §6 ein nummerierter Punkt |
| `docs/plan/adr/` (neue ADR, Index) | new, update (bedingt, Architect) | nur wenn die Richtung „umgekehrt“ (§6) zwischen Alternativen entschieden wird, die [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) nicht trägt; sonst genügt der Vertrag |
| `tools/harness/meldungskatalog-check.sh` | neu | liest die Codes aus dem Go-Code und aus der Tabelle in `SPEC-034` und vergleicht die Mengen |
| `tools/harness/meldungskatalog-check-gegenprobe.sh` | neu | Fehlformen in Temp-Bäumen je Punkt des Vertrags, wie `tools/harness/kopf-check-gegenprobe.sh`; dazu der Fall mit Makefile und Fragment |
| `harness/mk/meldungskatalog-check.mk` | neu | zwei Ziele, Prüfung und Gegenprobe, beide an `GATE_CHECKS` |
| `harness/sensors/meldungskatalog-check.md` | neu | Vertrag-Zeiger, Grenze, Ausgabe und Ausgänge, Sperren |
| `harness/README.md` §Sensors | update | zwei Zeilen mit Bindung |
| `internal/hexagon/model/fehler.go`, `spec/spezifikation.md` (`SPEC-034`) | update (bedingt) | nur ein Bestands-Befund, den der Vertrag rot macht, nach Entscheidung des Architects (§6, *Bestand*) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Slices von welle-v1-abschluss außer
`slice-v1-abschluss-container` und `slice-v1-abschluss-homebrew` liegen in `done/`
(WIP-Limit 1); dieser Slice ist Schritt 22 der Reihenfolge. Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness, M3 vor M4):
die Slices von welle-v1-abschluss, darin `slice-harness-integration-wait` direkt vor
`slice-v1-abschluss-herunterfahren` und `slice-harness-meldungskatalog-gate` direkt vor
`slice-v1-abschluss-container`, dann `slice-harness-abdeckung-gate` und
`slice-harness-coverage`, dann die Slices von welle-erster-release, danach
`slice-harness-commit-struktur-id`, `slice-tests-ueberlebende-mutanten`,
`slice-tests-ueberlebende-mutanten-driving`, `slice-harness-lint-warten-ohne-frist`,
`slice-harness-mutation`; die Reihenfolge
der Wellen-Slices steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md). Technisch hängt
dieser Slice an keinem davon; die Code-Tabelle und der Katalog liegen seit
welle-replay-semantik vor. Er steht direkt vor `slice-v1-abschluss-container`, dessen §1
das Gate als vorher geliefert nennt: Es misst die Codes, die die Slices der Welle
eingeführt haben, bevor die Betriebsdokumentation den Katalog schreibt, und die Frage, ob
es diesen Katalog mitprüft (§6, Punkt 8), ist vor dessen Entstehen entschieden. Dass die
Slices der Welle ihre Codes nicht unter dem Gate liefern, ist die Folge der Reihenfolge;
was sie am Katalog vorbei eingeführt haben, berichtigt dieser Slice nach §6 *Bestand*. Vor dem ersten
Code-Commit entscheidet der Architect die Randformen aus §6 im Vertrag und ob eine
ADR nötig ist (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Vertrag zieht den Katalog
  der Betriebsdokumentation als dritte Quelle hinzu (§6) und Skript, Gegenprobe und
  Bestands-Berichtigung sind nicht mehr in einer Review-Sitzung prüfbar; dann bleibt
  hier der Abgleich mit `SPEC-034`, die dritte Quelle geht in einen eigenen Slice.
- `in-progress` → `open` (blockiert — Carveout?): Die Richtung „umgekehrt“ (§6) ist
  nur als Gleichheit entscheidbar und färbt das Gate am Stand des Starts rot, ohne dass
  der Nutzer eine Stufung oder einen Carveout entschieden hat.

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

**Randformen des neuen Vertrags** ([`AGENTS.md`](../../../../AGENTS.md) §3.12). Keine
ist entschieden; der Architect entscheidet jede vor dem ersten Code-Commit im Vertrag
(`spezifikation.md §11`), eine Entscheidung des Nutzers wird dort festgehalten. Der
Implementer entscheidet keine still, sondern gibt eine nicht genannte zurück.

1. **Richtung „umgekehrt“ und Codes, die der Code noch nicht erzeugt.** Stand
   2026-10-07: Der Katalog in `SPEC-034` führt 34 Codes, die Code-Tabelle in
   `internal/hexagon/model/fehler.go` 19; 15 stehen nur im Katalog — die Rückfälle
   `PGR-E2000`, `PGR-E3000`, `PGR-E5000`, `PGR-E6000` und die Codes offener Slices von
   welle-v1-abschluss (`PGR-E2003` bis `PGR-E2007`, `PGR-E4004` bis `PGR-E4006`,
   `PGR-E5004`, `PGR-E6003`, `PGR-W3002`). Die Code-Tabelle führt nach
   `slice-replay-semantik-meldungscodes` „nur die Codes, die der Code erzeugt“. Eine
   Gleichheit beider Mengen wäre am Stand des Starts rot. Zu entscheiden: Gleichheit,
   Teilmenge (Go-Code im Katalog) mit einer Markierung im Katalog für vergebene, noch
   nicht erzeugte Codes, oder eine andere Form; und ob die Markierung mit dem Slice
   fällt, der den Code erzeugt.
2. **Rückfall-Codes und Klassen ohne eigenen Code.** `SPEC-034` sagt, jede Klasse hat
   einen Rückfall `…000`; die Code-Tabelle führt nur `PGR-E1000` und `PGR-E4000`. Ob ein
   Rückfall, den der Code nie erzeugt, im Code stehen muss, und ob das Gate prüft, dass
   jede Klasse ihren Rückfall im Katalog hat.
3. **Reservierte und zurückgezogene Codes.** Schwere `I` und Warnbereich 9 sind
   reserviert; ein entfallener Code bleibt vergeben (`SPEC-034` §Stabilität). Der
   Katalog hat heute keine Form für „zurückgezogen“; ob das Gate eine verlangt und wie
   ein zurückgezogener Code im Go-Code gewertet wird.
4. **Fundort im Code.** Nur die Konstanten der Code-Tabelle in `fehler.go`, oder jedes
   Literal `PGR-…` in Go-Dateien unter `internal/` und `cmd/` — heute stehen Literale in
   Kommentaren von Produkt-Dateien, etwa `PGR-E4006` in
   `internal/adapters/driving/pgwire/server.go`, das die Code-Tabelle nicht führt. Ob
   eine Konstante außerhalb von `fehler.go` als Code-Tabelle gilt.
5. **Codes in Tests und Kommentaren.** Ob `_test.go`-Dateien und `test/integration`
   zählen (ein Test, der einen nicht katalogisierten Code erwartet), und ob ein Code in
   einem Kommentar ein Befund ist oder nicht gelesen wird.
6. **Form.** Gilt nur die ERE aus `SPEC-034` (`PGR-[EWI][0-9]{4}`)? Ist eine Zeichenkette
   wie `PGR-E123`, `PGR-X0001` oder der Platzhalter `PGR-…` ein Befund „falsche Form“
   oder kein Code? Doppelte Konstante für denselben Code, doppelte Katalog-Zeile.
7. **Grenze des Katalogs.** Nur die erste Spalte der Tabelle in `SPEC-034`, oder auch
   Codes in Prosa von `SPEC-034` und anderen Kennungen; was gilt, wenn die Tabelle
   fehlt oder der Abschnitt nicht zu finden ist (Sperre mit eigenem Exit, wie bei
   `make kopf-check`).
8. **Katalog der Betriebsdokumentation.** [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md)
   und `SPEC-034` §Stabilität nennen den Katalog in der Betriebsdokumentation; das
   Benutzerhandbuch führt ihn heute in §7 *Fehlercodes*, mit Zeilen, die zwei Codes
   tragen (`PGR-E2000`, `PGR-E2001`). Die Entscheidung des Nutzers nennt die
   Spezifikation; ob das Handbuch als dritte Quelle in diesen Slice gehört oder eine
   eigene Adresse bekommt, ist zu entscheiden (Rückführung in §4).
9. **Prüfung der Klasse.** Ob das Gate nur Mengen abgleicht oder auch, dass die Spalte
   *Klasse / Bereich* zur ersten Ziffer passt (`SPEC-034` *Fehler*, *Warnungen*); ohne
   Entscheidung gilt die Abgrenzung in §1.
10. **Form der Meldung und Ausgänge.** Zeile je Befund auf stderr, etwa
    `meldungskatalog-check: <quelle>: <code>: <befund>`, Sortierung, Exit 0, 1 und 2,
    keine Ausgabe bei Grün — wie `harness/sensors/kopf-check.md`; zu entscheiden sind
    Quelle (Pfad mit Zeile?) und Befund-Texte.
11. **Gegenprobe.** Welche Fälle je Punkt des Vertrags abgelehnt bzw. angenommen werden
    müssen, welche akzeptierten Negative sie annimmt, und dass der Fall mit Makefile
    und Fragment beide Ziele an `GATE_CHECKS` zeigt.

**Risiken:**

- **Bestand.** Je nach Entscheidung zu 1, 2 und 4 ist das Gate am Stand des Starts rot.
  Berichtigt wird nur nach Entscheidung des Architects: in der Code-Tabelle (ohne
  Verhalten) oder im Katalog von `SPEC-034`; ein Rot, das erst die Slices von
  welle-v1-abschluss auflösen, braucht eine Stufung im Vertrag oder einen Carveout
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6), kein stilles Rot — **Ausgang:** offen bis
  Closure.
- **Folgelast für welle-v1-abschluss.** Jeder Slice, der einen Code einführt, muss ab
  diesem Slice Katalog und Code-Tabelle zusammen ändern; das Gate macht das sichtbar,
  die Pläne nennen es heute nicht — **Ausgang:** offen bis Closure.
- **Adresse aus `slice-v1-abschluss-container`.** Dieser Slice schließt vor dem
  Container-Slice; der verweist auf ihn als Liefernden, nicht als späteren Folge-Slice.
  Schreibt der Container-Slice den Katalog der Betriebsdokumentation fort, hält dieses
  Gate nur, was Punkt 8 entscheidet — **Ausgang:** offen bis Closure.
- **Lesart ohne Entscheidungsort.** Ein Werkzeug-Vertrag, dessen Lesarten erst im Code
  festgelegt werden, wiederholt `BEO-REPO/harness-lesart-ohne-entscheidungsort` (2×) —
  **Ausgang:** offen bis Closure.

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

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-07); alle Einträge liegen in `REPO`. Einschlägig:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — 13×, über der Schwelle, verkörpert ([`AGENTS.md`](../../../../AGENTS.md) §3.10); die Gegenprobe ist eigener Liefer-Punkt (§2).
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — 17×, verkörpert (§3.11); Sensor-Datei und Zeile in `harness/README.md` sagen nur zu, was Skript und Gegenprobe prüfen.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 11×, verkörpert (§3.12); die Randformen stehen in §6 und sind vor dem Code zu entscheiden.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` — 1×; der Slice ist wellenlos, §4 nennt die Entscheidung des Architects vor dem Code.
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` — 2×; mit diesem Slice das dritte Mal berührt, falls eine Lesart erst im Code entschieden wird (§6).
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` — 3×, verkörpert (§3.13); die Geber nennen diesen Slice, dieser nennt sie in §1.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` — 1×; ein Gate, das Dateien sucht, darf eine nicht lesbare Quelle nicht still übergehen (§6, Punkt 7).
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` — 2×; das Fragment unter `harness/mk/` muss so an `GATE_CHECKS` hängen, wie es sich liest (Gegenprobe mit Makefile, §2).

Die Einträge über der Schwelle bekommen ihren Ausgang bei der Closure ihrer Welle, nicht hier.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
