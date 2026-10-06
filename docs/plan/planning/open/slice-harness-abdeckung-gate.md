# Slice slice-harness-abdeckung-gate: Nachweisart Gate und geteilte Messung in der Abdeckung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Er trägt zum
Abnahmeszenario 17 (M3) bei: Mit ihm steht die Anforderung in den Abdeckungstabellen
*teilweise*, nachweisbar wird das Szenario mit `slice-harness-coverage`. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-06: nach `slice-harness-blackbox-einstieg` und vor `slice-harness-coverage`
(WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Nachweis der Messmethoden 1 und 3). Bindung an Entscheidungen: [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Nachweisart Gate, geteilte Messung; Folgepflicht dieser ADR), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Deklaration je Anforderung und Pfad, gilt sonst unverändert).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `tools/test/abdeckung.sh` liest nach
[ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) neben den Tests auch
Abdeckungs-Deklarationen im Kopf von Gegenproben (Nachweisart *Gate*, nur `LH-QA` und
`LH-RB`) und den Pfad `Messung-<i>-von-<n>` für eine geteilte Messung; eine geteilte
Messung ist vollständig, wenn alle Teile 1 bis n belegt sind. `make abdeckung-gegenprobe`
lehnt die Fehlformen aus §6 ab. Die Gegenprobe des Lint-Gates deklariert Teil 1
(statische Analyse) und Teil 3 (Lage der Unit-Tests) von LH-QA-07; die Anforderung
steht damit in den Abdeckungstabellen als *teilweise*. `harness/README.md` §Sensors
nennt die Nachweisart beim Vertrag von `make abdeckung-check`.

**Herkunft:** Bei `slice-lastenheft-pruefbarkeit` hat das Team entschieden, dass
LH-QA-07 in den Abdeckungstabellen nicht als unbelegt stehen bleibt (2026-10-06). Den
Nachweis-Weg hat der Architect vorgeschlagen; der Nutzer hat am 2026-10-06
[ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) angenommen, deren Folgepflicht einen eigenen Harness-Slice verlangt, und diesen
Slice mit seinem Platz in der Reihe entschieden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Teil 2 (Anweisungsabdeckung) von LH-QA-07 — die Deklaration gehört an die Gegenprobe
  des Coverage-Gates, die es erst mit `slice-harness-coverage` gibt; der deklariert ihn
  dort, und erst dann ist die Anforderung vollständig.
- Die Gegenprobe des Lint-Gates und ihre Fälle — liefern `slice-harness-lint` und die
  vier Umstellungs-Slices (der Fall eines White-Box-Tests je Paketgruppe). Dieser Slice
  ergänzt an ihr nur die Deklaration, keinen Fall; deckt ein Fall die deklarierte
  Messmethode nicht, ist das die Rückführung aus §4.
- Prüfung, ob eine Gegenprobe an `GATE_CHECKS` hängt — Grenze nach
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Konsequenzen,
  Negativ): Sie steht als Grenze im Skriptkopf und bleibt Review.
- Deklarationen an Tests und der ungeteilte Pfad `Messung` — Bestand bleibt bewusst
  stehen ([ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)); die heute
  mit ungeteilter Messung belegten Qualitätsanforderungen werden nicht umgeschrieben.
- Lastenheft, Produkt-Code und die Spezifikation außerhalb ihres Abschnitts für
  Harness-Werkzeuge — Schicht-Abgrenzung: Der Slice ändert das Abdeckungs-Skript,
  seine Gegenprobe, den Kopf der Lint-Gegenprobe, die erzeugten Tabellen unter
  `docs/user/` und `harness/README.md`, in der Spezifikation nur den Abdeckungs-Vertrag
  in jenem Abschnitt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Werkzeug und Vertrag: `tools/test/abdeckung.sh` liest Deklarationen im Kopf von
      Gegenproben (Nachweisart Gate) und den Pfad `Messung-<i>-von-<n>` nach den
      Randformen aus §6, schreibt die Tabelle der Nachweisart Gate und zählt eine
      geteilte Messung erst mit allen Teilen als vollständig; die Deklarationen an Tests
      und der Pfad `Messung` bleiben gültig, die übrigen Tabellen unverändert bis auf die
      neue Nachweisart. Die Randformen stehen am Ort aus §6 (Abschnitt für Harness-Werkzeuge der Spezifikation, angelegt von `slice-harness-vertraege-spezifikation`), und
      `harness/README.md` §Sensors nennt beim Vertrag von `make abdeckung-check` die
      Nachweisart Gate und die geteilte Messung, mit Bindung an [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md);
      beide sagen nur zu, was die Gegenprobe prüft (`AGENTS.md` §3.11).
- [ ] Gegenprobe: `make abdeckung-gegenprobe` führt je Fehlform aus §6 einen Fall, der
      abgelehnt wird, und je Gutform einen, der angenommen wird; die vorhandenen Fälle
      bleiben; je Zusage ist die Mutation gesehen (`AGENTS.md` §3.10). Ihre Zeile in
      `harness/README.md` §Sensors nennt die neuen Fälle.
- [ ] Erste Gate-Deklarationen: Die Gegenprobe des Lint-Gates trägt im Kopf die
      Deklaration von LH-QA-07, Teil 1 und Teil 3 von 3; die Tabellen sind mit
      `make abdeckung` neu geschrieben, LH-QA-07 steht in der Gesamtsicht als
      *teilweise* und nicht in `abdeckung-vollstaendig.md`; `make abdeckung-check` ist grün.
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
| `tools/test/abdeckung.sh` | update | Nachweisart Gate (Kopf von `tools/**/*-gegenprobe.sh`), Pfad `Messung-<i>-von-<n>`, Tabelle Gate, Vollständigkeit einer geteilten Messung; Grenze „ob die Gegenprobe an `GATE_CHECKS` hängt, prüft das Skript nicht“ im Skriptkopf |
| `tools/test/abdeckung-gegenprobe.sh` | update | je Fehlform und je Gutform aus §6 ein Fall in einem Temp-Baum |
| Gegenprobe des Lint-Gates (aus `slice-harness-lint`) | update | Kopf: Deklaration LH-QA-07, Teil 1 und 3 von 3; kein neuer Fall |
| `docs/user/abdeckung-*.md` | update | mit `make abdeckung` neu geschrieben, dazu die Tabelle der Nachweisart Gate |
| `harness/README.md` | update | §Sensors: Vertrag von `make abdeckung-check` und `make abdeckung-gegenprobe` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-blackbox-einstieg` liegt in `done/`
(WIP-Limit 1). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-06:
`slice-harness-lint`, die vier Umstellungs-Slices, dieser Slice,
`slice-harness-coverage`, `slice-harness-mutation`. Grund: Teil 3 gilt erst, wenn
`testpackage` überall scharf ist, und das ist es mit `slice-harness-blackbox-einstieg`.
Vor dem ersten Code-Commit bestätigt der Architect die Randformen aus §6 und
entscheidet die offenen (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Erweiterung ändert die
  Form der bestehenden Tabellen oder der Test-Deklarationen so weit, dass Skript,
  Gegenprobe und Tabellen nicht in einer Review-Sitzung prüfbar sind; dann zuerst die
  Erweiterung ohne Deklaration, die Deklaration in einem eigenen Slice.
- `in-progress` → `open` (blockiert — Carveout?): Die Gegenprobe des Lint-Gates liegt
  nicht unter `tools/**/*-gegenprobe.sh` oder führt für Messmethode 1 oder 3 keinen
  Fall, der rot wird; dann deklarierte sie, was sie nicht prüft (`AGENTS.md` §3.11),
  und zuerst entscheidet der Architect Ort oder Fall.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit `abdeckung-check` und `abdeckung-gegenprobe`,
LH-QA-07 in der Gesamtsicht *teilweise*, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des Vertrags** (`AGENTS.md` §3.12). Vorgegeben durch
[ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) und §6 *Nachweis* von
`slice-lastenheft-pruefbarkeit` (Architect, 2026-10-06); vor dem Code zu bestätigen.
**Ort** (F-406): der Abschnitt für Harness-Werkzeuge der Spezifikation (Technik-Stratum), angelegt von `slice-harness-vertraege-spezifikation`; dort steht der Vertrag von `abdeckung`, den dieser Slice erweitert, und die Kennungen der neuen Randformen vergibt dieser Slice. Die ADR trägt Entscheidung und Gründe. Jede dort zugesagte Randform bekommt einen Fall in
`make abdeckung-gegenprobe` (`AGENTS.md` §3.10, §3.11). Bis zum Code stehen sie hier:

- **Pfad-Schreibweise** — `Messung-<i>-von-<n>` mit 2 ≤ n und 1 ≤ i ≤ n.
- **Abgelehnt** — i > n; n < 2; verschiedene n für dieselbe Anforderung;
  `Messung` neben einer geteilten Messung derselben Anforderung.
- **Nachweisart Gate nur für `LH-QA` und `LH-RB`** — eine Gate-Deklaration an einer
  `LH-FA` belegt kein Produktverhalten
  ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), Entscheidung).
- **Gelesen werden** — Dateien `tools/**/*-gegenprobe.sh`, die Deklaration als
  `#`-Zeilen im Kopf.
- **Nachweis-Spalte** — der Skriptpfad.
- **Tabelle** — eine eigene Tabelle der Nachweisart Gate neben E2E und Unit.
- **Vollständigkeit** — eine geteilte Messung ist vollständig, wenn alle Teile 1 bis n
  belegt sind; sonst *teilweise*. `abdeckung-vollstaendig.md` (die RTM-Quelle) führt sie
  erst dann.
- **Messmethode (4) von LH-QA-07** — kein eigener Teil; die Nachweisart Gate erfüllt
  sie ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)). Die Gesamtzahl
  ist darum 3.
- **Anschluss an `GATE_CHECKS`** — prüft das Skript nicht; Grenze im Skriptkopf, Review.

**Offen** — vor dem Code vom Architect zu entscheiden und am selben Ort festzuhalten
(Ort oben); was dort nicht steht, entscheidet der Implementer nicht, er gibt es
zurück:

- **Form des Kopfes** — wo der Kopf einer Gegenprobe endet (nach der Shebang-Zeile bis
  zur ersten Zeile ohne `#`?), ob Folgezeilen wie bei Tests fortsetzen, ob der Trenner
  „ — “ Pflicht ist, und ob eine eingerückte oder eine Deklaration nach dem Kopf
  abgelehnt wird.
- **Derselbe Teil zweimal** — zwei Nachweise deklarieren denselben Teil i derselben
  Anforderung: zulässig (wie zwei Tests am selben Pfad) oder Fehlform.
- **Darstellung der Teile in der Gesamtsicht** — wie die Spalte *Messung* eine geteilte
  Messung zeigt (belegte Teile, Nachweisart je Teil).
- **Dateiname der neuen Tabelle** und ob der Kopftext von `abdeckung-gesamt.md` die
  Nachweisart Gate nennt.
- **Label der RTM-Quelle** — `trace.coverage` in `.d-check.yml` trägt das Label `Tests`;
  ob es mit Gate-Nachweisen in `abdeckung-vollstaendig.md` angepasst wird.
- **Geteilte Messung an Go-Tests** — ob `Messung-<i>-von-<n>` auch in einer
  Deklaration an einem Test zulässig ist ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) sagt „jeder Nachweis“) oder nur
  an Gegenproben.
- **Deklaration außerhalb von `tools/**/*-gegenprobe.sh`** — was mit einer
  Gate-Deklaration in einem Skript geschieht, das nicht unter das Muster fällt (etwa
  eine Gegenprobe an der Gate-Kette mit anderem Namen): still überlesen oder Fehler,
  wie bei der eingerückten Deklaration an Tests.
- **Zahlform von i und n** — führende Null, Leerraum, mehrstellige Zahlen.

**Risiken:**

- **Falsche Gesamtzahl n** — Gesamtzahl und Zuordnung eines Teils zur Messmethode sind
  Urteil dessen, der deklariert; das Skript fängt nur verschiedene n für dieselbe
  Anforderung, eine einheitlich falsche Zahl fängt nur das Review
  ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), Re-Evaluierungs-
  Trigger). — **Ausgang:** — (bei Closure)
- **Deklaration sagt mehr zu, als die Gegenprobe prüft** — Teil 3 an der
  Lint-Gegenprobe trägt nur, wenn sie für beide Bedingungen von Messmethode 3 einen
  Rot-Fall führt: je Paketgruppe einen White-Box-Test (der Test liegt in der Einheit)
  und die Fälle zur Brücke aus `slice-harness-lint` (Zugriff auf Internes an der
  Brücke vorbei, Test in der Brückendatei)
  (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, 11×). Vor dem Code gegen die
  Fälle der Lint-Gegenprobe abgleichen. — **Ausgang:** — (bei Closure)
- **Bestehende Prüfung fällt weg** — die Erweiterung tritt neben die Prüfung der
  Test-Deklarationen, nicht an ihre Stelle; die vorhandenen Fälle der Gegenprobe
  bleiben (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). — **Ausgang:** — (bei
  Closure)
- **Andere Gegenproben werden gelesen** — `tools/**/*-gegenprobe.sh` trifft auch die
  vorhandenen Gegenproben ohne Deklaration; sie dürfen weder Fehler noch Zeilen
  erzeugen. — **Ausgang:** — (bei Closure)

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `daef83b` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (10×, verkörpert in `AGENTS.md`
  §3.10) — die neue Deklarationsform ist ein neuer Vertrag; daher je Fehlform ein Fall
  der Gegenprobe (erster Liefer-Punkt).
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (11×, verkörpert in §3.11) — die
  Gate-Deklaration ist eine Zusage; ein Risiko in §6.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (8×, verkörpert in §3.12) und
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — darum stehen die
  offenen Randformen in §6 und werden vor dem Code entschieden; ein drittes Auftreten
  in diesem Slice brächte den zweiten Eintrag auf die Schwelle.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — §4 *Start* nennt den
  Architect-Schritt vor dem Code.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — ein Risiko in §6.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
