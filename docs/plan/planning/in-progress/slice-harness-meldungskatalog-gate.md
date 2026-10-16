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

**Berührte Spec-Stellen:** `SPEC-034` (gelesen, Tabelle der Codes; geändert nur, wenn der Architect einen Bestands-Befund dort berichtigt) · [`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) (Vertrag des Gates, vom Architect vor dem Code geschrieben; die nächste freie Nummer ist danach 055) · `SPEC-052` Punkt 13 (Beleg des Handbuchs in der schnellen Kette, vom Architect nachgezogen)

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
ein Code des Katalogs im Go-Code fehlt, es sei denn, er ist ein Rückfall `…000`; und
wenn die Codes der Fehlercode- und Warnungstabellen des Handbuchs nicht genau die der
Code-Tabelle sind. Was als Quelle, Code, Befund und Sperre zählt, entscheidet
[`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) (Architect,
vor dem ersten Code-Commit); die Randformen dazu stehen in §6 mit ihrem Punkt. Geliefert wie
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
der Spezifikation. Das Handbuch kommt als dritte Quelle dazu (Empfehlung des Architect
vom 2026-10-17, §6 Punkt 8; der Nutzer hat es nicht ausgeschlossen): Es kostet eine
Funktion im Skript und schützt `slice-v1-abschluss-container`, bevor der den Katalog der
Betriebsdokumentation fortschreibt.

**Übernommen aus welle-v1-abschluss** (dort §6, Out-of-Scope) **und aus
`slice-v1-abschluss-container`** (dort §1, Abgrenzung): das Gate für Code-Tabelle und
Katalog; das ist das Ziel oben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Codes erzeugen, die der Katalog führt und der Code noch nicht — Schicht-Abgrenzung:
  kein Verhalten im Produkt-Code und keine Änderung an `internal/hexagon/model/fehler.go`
  oder `SPEC-034`; das Gate ist am Stand des Starts grün (§6 *Bestand*), die vier Codes des
  Katalogs ohne Konstante sind Rückfälle, die `SPEC-054` Punkt 4b zulässt.
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
      Der Vertrag `SPEC-054` steht vor dem ersten Code-Commit in `harness/vertraege/`
      und entscheidet jede Randform aus §6; das Prüfskript (bash und awk, ohne Docker)
      hängt mit seinem Ziel über ein Fragment unter `harness/mk/` an `GATE_CHECKS`, ist
      am Stand des Starts grün, ohne dass `fehler.go` oder `SPEC-034` sich ändern, und
      steht in der schnellen Kette; die Gegenprobe steht in `GATE_SCHNELL_AUSNAHME`
      (`harness/mk/gates-schnell.mk`).
- [ ] Gegenprobe als eigenes Gate-Ziel in Temp-Bäumen: je Punkt von `SPEC-054` und je
      Element einer Aufzählung darin (`AGENTS.md` §3.10) ein Fall, der abgelehnt bzw.
      angenommen werden muss — darunter ein Code nur im Go-Code, nur im Katalog, nur im
      Handbuch, in falscher Form; dazu in einem Temp-Baum mit
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
| `harness/vertraege/SPEC-054-meldungskatalog.md`, `README.md` des Ordners, `SPEC-052-gate-kette.md` Punkt 13 | neu, update (Architect, vor Code, geliefert) | Vertrag des Gates; Beleg des Handbuchs in der schnellen Kette |
| `docs/plan/adr/` | keine ADR | [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) trägt die Folgepflicht, die Form folgt dem Muster der Sensoren; die Gründe (Handbuch als dritte Quelle, Rückfälle zugelassen, Go-Test unmöglich wegen des Build-Kontexts) stehen im Vertrag |
| `tools/harness/meldungskatalog-check.sh` | neu | liest die Codes aus `fehler.go`, aus der Tabelle in `SPEC-034` und aus den Tabellen des Handbuchs und vergleicht die Mengen (`SPEC-054`) |
| `tools/harness/meldungskatalog-check-gegenprobe.sh` | neu | Fehlformen in Temp-Bäumen je Punkt des Vertrags, wie `tools/harness/kopf-check-gegenprobe.sh`; dazu der Fall mit Makefile und Fragment |
| `harness/mk/meldungskatalog-check.mk` | neu | zwei Ziele, Prüfung und Gegenprobe, beide an `GATE_CHECKS` |
| `harness/mk/gates-schnell.mk` | update | `meldungskatalog-check-gegenprobe` in `GATE_SCHNELL_AUSNAHME` |
| `harness/sensors/meldungskatalog-check.md` | neu | Vertrag-Zeiger, Grenze, Ausgabe und Ausgänge, Sperren |
| `harness/README.md` §Sensors | update | zwei Zeilen mit Bindung |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Slices von welle-v1-abschluss außer
`slice-v1-abschluss-container` liegen in `done/`
(WIP-Limit 1); dieser Slice ist Schritt 43 der Reihenfolge. Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness, M3 vor M4):
die drei Slices der Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; Entscheidung des Nutzers vom 2026-10-13, vor die Release-Slices gezogen), dann die Slices von welle-v1-abschluss, darin `slice-harness-integration-wait` direkt vor
`slice-v1-abschluss-herunterfahren` und `slice-harness-meldungskatalog-gate` direkt vor
`slice-v1-abschluss-container`, dann `slice-harness-abdeckung-gate` und
`slice-harness-coverage`, dann die Slices von welle-erster-release und welle-v1-nachtrag, danach
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

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und
  Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness`
  liegt über 900; Basis ist `8d14a918` (HEAD vor dem Anspruch; er enthält Vertrag und Plan, die Move-Commits berühren `tools` und `harness` nicht). Die Schätzung des
  Architect (2026-10-17) in dieser Metrik: Spanne 450 bis 750, Mitte 580, Obergrenze 750
  (Skript 110 bis 150; Gegenprobe 300 bis 450, das Zwei- bis Dreifache des Skripts nach dem
  Faktor aus `AGENTS.md` §3.14; Fragment, Sensor-Datei und Zeilen in `harness/README.md` und
  `gates-schnell.mk` rund 80); die Schätzung liegt unter der Marke, ein Schnitt vor dem Code
  ist nicht geboten. Meldeschwelle: Stehen bei 750 noch Fälle der Gegenprobe aus, meldet der
  Implementer es dem Architect, bevor er weiterschreibt.
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

**Randformen des neuen Vertrags** ([`AGENTS.md`](../../../../AGENTS.md) §3.12). Der Architect
hat jede am 2026-10-17 im Vertrag
[`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) entschieden; die Nummer
nennt den Punkt. Der Implementer entscheidet keine still, sondern gibt eine nicht genannte
zurück. Stand am 2026-10-17: C (Konstanten) 30 Codes, K (Katalog) 34, H (Handbuch) 30; die
Zahlen „34 gegen 19“ des ersten Plans sind durch die Slices von welle-v1-abschluss überholt.

1. **Richtung „umgekehrt“, nicht erzeugte Codes** — Punkt 4b: K ohne C ist ein Befund, außer
   bei Endung `000`; eine Markierung im Katalog gibt es nicht, und am Stand des Starts ist
   nichts mehr offen.
2. **Rückfälle** — Punkt 4b und 5: Ein nicht erzeugter Rückfall muss nicht im Code stehen; je
   Klasse 1 bis 6 steht `PGR-E<k>000` in K.
3. **Reserviert, zurückgezogen** — Punkt 7: keine Form; der Slice, der zurückzieht, ändert den
   Vertrag.
4. **Fundort im Code** — Punkt 1 (nur die Konstanten `Code…` in `fehler.go`) und 6 (ein
   Literal `"PGR-` in einer anderen Nicht-Test-Datei unter `cmd/` und `internal/` ist ein
   Befund).
5. **Tests und Kommentare** — Grenze: nicht gelesen.
6. **Form** — Punkte 2 und 3: ERE `PGR-[EWI][0-9]{4}` ganz; doppelte Konstante, Katalog- und
   Handbuch-Zeile sind `doppelt`.
7. **Grenze des Katalogs; fehlende Tabelle** — Punkte 1 und 8: nur die erste Zelle der
   Tabellenzeilen im Abschnitt `### SPEC-034 `; fehlende Datei, Abschnitt oder Zeilen sind
   eine Sperre (Exit 2).
8. **Katalog der Betriebsdokumentation** — Punkt 4c und 4d: Das Handbuch ist dritte Quelle und
   führt genau die Codes der Code-Tabelle. Empfehlung des Architect, nicht Entscheidung des
   Nutzers; sie gilt, bis er widerspricht. Die Rückführung in §4 bindet nicht an die Quelle,
   sondern an die Größe: Die Schätzung liegt unter der Marke.
9. **Klasse** — Grenze: Die Spalten werden nicht gelesen.
10. **Form der Meldung und Ausgänge** — Punkte 9 und 8: Zeile, Sortierung, Exit 0, 1 und 2.
11. **Gegenprobe** — Punkt 11 und DoD: je Punkt ein Fall; der Anschluss-Fall mit Makefile und
    Fragment.
12. **Aufrufarten** — Punkt 10: kein Stand und kein Ergebnis außer dem Ausgang; keine
    Sonderbehandlung von `-n`, `-q`, `-t`, `-i`.
13. **Wahl der Kette** — Punkt 11: `meldungskatalog-check` in der schnellen Kette (`SPEC-052`
    Punkt 13), die Gegenprobe in der Ausnahmeliste.

**Risiken:**

- **Bestand.** Am Stand des Starts (2026-10-17) ist das Gate nach `SPEC-054` grün: C ist
  Teilmenge von K, K ohne C sind vier Rückfälle, H gleich C, kein Literal außerhalb von
  `fehler.go`. Ändert sich der Stand bis zum Start, ist ein Rot ein Befund für den
  Architect, kein stilles Rot und keine Stufung ([`AGENTS.md`](../../../../AGENTS.md) §3.6)
  — **Ausgang:** offen bis Closure.
- **Folgelast für welle-v1-abschluss.** Jeder Slice, der einen Code einführt, muss ab
  diesem Slice Katalog und Code-Tabelle zusammen ändern; das Gate macht das sichtbar,
  die Pläne nennen es heute nicht — **Ausgang:** offen bis Closure.
- **Adresse aus `slice-v1-abschluss-container`.** Dieser Slice schließt vor dem
  Container-Slice; der verweist auf ihn als Liefernden, nicht als späteren Folge-Slice.
  Schreibt der Container-Slice den Katalog der Betriebsdokumentation fort, hält dieses
  Gate die Tabellen `### Fehlercodes` und `### Warnungen` des Handbuchs (`SPEC-054` Punkt 1),
  nicht andere Betriebsdokumente — **Ausgang:** offen bis Closure.
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

**Nachgezählt** mit der Sendung aus `slice-harness-vertraege-umzug-gate-kette` (Ort neuer
Kennungen: `harness/vertraege/`, nächste freie Nummer 053; `AGENTS.md` §3.13): drei
Liefer-Punkte (§2: Prüfung, Gegenprobe, Sensor-Datei, unverändert). Schichten ohne die bedingte
Zeile von §3 (Bestands-Befund in `internal/hexagon/model/fehler.go` und `SPEC-034`): eine,
Harness (Vertragsdatei, `tools/`, `harness/mk/`, `harness/sensors/`, `harness/README.md`). Mit
der bedingten Zeile kämen Produkt-Code und Spezifikation dazu, also drei; das war vor der
Sendung ebenso (die Vertragsdatei lag in der Spezifikation) und ist Sache der Entscheidung des
Architects zu *Bestand* in §6, die Sendung ändert daran nichts.

**Nachgezählt** mit der Vorprüfung des Architect (2026-10-17, `AGENTS.md` §3.13): drei
Liefer-Punkte (§2: Prüfung mit Anschluss, Gegenprobe, Sensor-Datei mit Zeilen in
`harness/README.md`). Schichten: eine, Harness; die bedingte Zeile für `fehler.go` und
`SPEC-034` ist entfallen (Bestand grün). Das Handbuch wird gelesen, nicht geändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
