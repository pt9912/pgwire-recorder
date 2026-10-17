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
2026-10-06: nach `slice-harness-lint` und vor `slice-harness-coverage`; seit der
Entscheidung vom 2026-10-07 stehen `slice-harness-commit-struktur-id`,
`slice-harness-integration-wait` und `slice-harness-meldungskatalog-gate` davor (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Nachweis der Messmethoden 1 und 3). Bindung an Entscheidungen: [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Nachweisart Gate, geteilte Messung; Folgepflicht dieser ADR), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Deklaration je Anforderung und Pfad; gilt unverändert, auch für die RTM: Empfehlung des Architect vom 2026-10-17, dass keine ergänzende ADR nötig ist, siehe §1; der Nutzer hatte am 2026-10-06 eine ADR für die Darstellung verlangt und entscheidet, ob er darauf verzichtet).

**Berührte Spec-Stellen:** `SPEC-048` (Vertrag von `abdeckung` in `harness/vertraege/SPEC-048-abdeckung.md`, fortgeschrieben)

**Verantwortlich:** pt9912

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

Die RTM (`make doc-trace`) bleibt bei einer Quelle, `docs/user/abdeckung-vollstaendig.md`, und
zählt eine Anforderung weiterhin nur, wenn alle Pfade gedeckt sind, gleich aus welcher Nachweisart
([ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md));
die Nachweisarten und die gedeckten Pfade je Anforderung stehen in `abdeckung-vollstaendig.md` und
`abdeckung-gesamt.md`. Geändert wird nur das Label der Quelle in `.d-check.yml` von `Tests` auf
`Abdeckung`, weil ab diesem Slice eine Anforderung auch aus Gate-Nachweisen vollständig sein kann.
**Empfehlung des Architect (2026-10-17), Entscheidung beim Nutzer:** Damit entfällt die neue ADR, die
der Nutzer am 2026-10-06 für eine reichere Darstellung verlangt hat (Nachweisart und Pfad je
Anforderung in der RTM selbst): Die einfachste Form, die kein stilles Grün erlaubt, ist die vorhandene
Quelle, die nur vollständig belegte Anforderungen nennt; eine Darstellung je Nachweisart in der RTM
bräuchte mehrere Quellen, die d-check getrennt als belegt zählt, und gäbe einer teilweise belegten
Anforderung genau das Grün, das [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) ausschließt. Verlangt der Nutzer sie
trotzdem, schreibt der Architect die ADR (`Proposed`, vor dem Anspruch), und der Planner legt dafür den
Slice `slice-harness-rtm-nachweisarten` an; dieser Slice behält dann Liefer-Punkt 3 ohne die RTM.

**Bindung aus `slice-harness-integration-volumes`** (Geber; `AGENTS.md` §3.13, Sendung an diesen
Slice): `tools/test/integration-volumes-gegenprobe.sh` liegt unter `tools/**/*-gegenprobe.sh` und
trägt keine Abdeckungs-Deklaration; das Skript liest sie mit und erzeugt dort weder Fehler noch
Zeilen. Der Fall steht in der Gegenprobe dieses Slice (DoD, Liefer-Punkt 2) und in §6 (Risiko
*Andere Gegenproben werden gelesen*).

**Herkunft:** Bei `slice-lastenheft-pruefbarkeit` hat das Team entschieden, dass
LH-QA-07 in den Abdeckungstabellen nicht als unbelegt stehen bleibt (2026-10-06). Den
Nachweis-Weg hat der Architect vorgeschlagen; der Nutzer hat am 2026-10-06
[ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) angenommen, deren Folgepflicht einen eigenen Harness-Slice verlangt, und diesen
Slice mit seinem Platz in der Reihe entschieden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eingriff in die vorhandenen Gegenproben unter `tools/` außer dem Kopf der Lint-Gegenprobe —
  Bestand bleibt bewusst stehen: Sie tragen keine Deklaration und werden vom Skript nur mitgelesen;
  die Bindung an `tools/test/integration-volumes-gegenprobe.sh` (Sendung aus
  `slice-harness-integration-volumes`, §1: kein Kopf in Deklarationsform) hält dieser Slice mit
  einem Fall seiner Gegenprobe fest, ändert die Datei aber nicht.
- Teil 2 (Anweisungsabdeckung) von LH-QA-07 — die Deklaration gehört an die Gegenprobe
  des Coverage-Gates, die es erst mit `slice-harness-coverage` gibt; der deklariert ihn
  dort, und erst dann ist die Anforderung vollständig.
- Die Gegenprobe des Lint-Gates und ihre Fälle — liefert `slice-harness-lint`, auch den
  Fall eines White-Box-Tests je Paketgruppe der vier Umstellungs-Slices. Dieser Slice
  ergänzt an ihr nur die Deklaration, keinen Fall; deckt ein Fall die deklarierte
  Messmethode nicht, ist das die Rückführung aus §4.
- Prüfung, ob eine Gegenprobe an `GATE_CHECKS` hängt — Grenze nach
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Konsequenzen,
  Negativ): Sie steht als Grenze im Skriptkopf und bleibt Review.
- Deklarationen an Tests und der ungeteilte Pfad `Messung` — Bestand bleibt bewusst
  stehen ([ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)); die heute
  mit ungeteilter Messung belegten Qualitätsanforderungen werden nicht umgeschrieben.
- Lastenheft, Produkt-Code und die Spezifikation — Schicht-Abgrenzung: Der Slice ändert das Abdeckungs-Skript,
  seine Gegenprobe, den Kopf der Lint-Gegenprobe, die erzeugten Tabellen unter
  `docs/user/` und `harness/README.md`, in `harness/vertraege/` nur den Abdeckungs-Vertrag
  `SPEC-048`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Werkzeug und Vertrag: `tools/test/abdeckung.sh` liest Deklarationen im Kopf von
      Gegenproben (Nachweisart Gate) und den Pfad `Messung-<i>-von-<n>` nach den
      Randformen aus §6, schreibt die Tabelle der Nachweisart Gate und zählt eine
      geteilte Messung erst mit allen Teilen als vollständig; die Deklarationen an Tests
      und der Pfad `Messung` bleiben gültig, die übrigen Tabellen unverändert bis auf die
      neue Nachweisart. Die Randformen stehen am Ort aus §6 (`harness/vertraege/SPEC-048-abdeckung.md`), und
      `harness/README.md` §Sensors nennt beim Vertrag von `make abdeckung-check` die
      Nachweisart Gate und die geteilte Messung, mit Bindung an [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md);
      beide sagen nur zu, was die Gegenprobe prüft (`AGENTS.md` §3.11).
- [x] Gegenprobe: `make abdeckung-gegenprobe` führt je Fehlform aus §6 einen Fall, der
      abgelehnt wird, und je Gutform einen, der angenommen wird; die vorhandenen Fälle
      bleiben, dazu ein Fall, der eine Gegenprobe ohne Deklaration im Muster (nach dem Vorbild
      von `tools/test/integration-volumes-gegenprobe.sh`, Sendung aus
      `slice-harness-integration-volumes`) ohne Fehler und ohne Zeile annimmt; je Zusage ist die
      Mutation gesehen (`AGENTS.md` §3.10). Ihre Zeile in
      `harness/README.md` §Sensors nennt die neuen Fälle.
- [x] Erste Gate-Deklarationen und RTM: Die Gegenprobe des Lint-Gates trägt im Kopf die
      Deklaration von LH-QA-07, Teil 1 und Teil 3 von 3; die Tabellen sind mit
      `make abdeckung` neu geschrieben, LH-QA-07 steht in der Gesamtsicht als
      *teilweise* (Messung `1/3 Gate; 2/3 —; 3/3 Gate`) und nicht in `abdeckung-vollstaendig.md`;
      `make abdeckung-check` ist grün. Das Label von `trace.coverage` in `.d-check.yml` heißt
      `Abdeckung` (siehe §1); `make doc-trace` zeigt LH-QA-07 nicht als belegt, eine vollständig
      belegte Anforderung mit dem neuen Label (Beleg in §7).
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
| `tools/test/abdeckung.sh` | update | Nachweisart Gate (Kopf von `tools/**/*-gegenprobe.sh`), Pfad `Messung-<i>-von-<n>`, Tabelle Gate, Vollständigkeit einer geteilten Messung; Grenze „ob die Gegenprobe an `GATE_CHECKS` hängt, prüft das Skript nicht“ im Skriptkopf |
| `tools/test/abdeckung-gegenprobe.sh` | update | je Fehlform und je Gutform aus §6 ein Fall in einem Temp-Baum |
| Gegenprobe des Lint-Gates (aus `slice-harness-lint`) | update | Kopf: Deklaration LH-QA-07, Teil 1 und 3 von 3; kein neuer Fall |
| `docs/user/abdeckung-*.md` | update | mit `make abdeckung` neu geschrieben (nie von Hand), dazu die neue Tabelle `abdeckung-gate.md` |
| `harness/sensors/abdeckung-check.md` | update | Ausgänge: Exit 2 auch bei Lastenheft ohne Anforderung und fehlendem `docs/user/`; Tabellen fünf statt vier |
| `harness/vertraege/SPEC-048-abdeckung.md` | geschrieben (Architect) | Punkte 2, 3, 4, 6 fortgeschrieben, Punkte 7 und 8 neu, Grenze; liegt vor dem Anspruch; weicht der Code ab, zieht der Implementer den Vertrag im selben Commit nach (`AGENTS.md` §3.9) |
| `harness/README.md` | update | §Sensors: Vertrag von `make abdeckung-check` und `make abdeckung-gegenprobe` |
| `harness/mk/abdeckung.mk` | update | Kommentar und Hilfetexte nennen auch die Gegenproben als Quelle der Deklarationen (`AGENTS.md` §3.7); keine Änderung am Rezept |
| `.d-check.yml` | update | `trace.coverage`: Label `Tests` → `Abdeckung`; Quelle und Status „belegt nur bei vollständiger Belegung“ bleiben |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint`,
`slice-harness-integration-wait`, `slice-harness-meldungskatalog-gate`, die Slices von
welle-v1-abschluss und `slice-harness-integration-volumes` liegen in `done/` (WIP-Limit 1). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness, M3 vor M4):
die drei Slices der Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; Entscheidung des Nutzers vom 2026-10-13, vor die Release-Slices gezogen), dann die Slices von welle-v1-abschluss, darin `slice-harness-integration-wait` direkt vor
`slice-v1-abschluss-herunterfahren` und `slice-harness-meldungskatalog-gate` direkt vor
`slice-v1-abschluss-container`, dann `slice-harness-integration-volumes`,
`slice-harness-abdeckung-gate` und `slice-harness-coverage`, dann die Slices von welle-erster-release und welle-v1-nachtrag, danach
`slice-harness-commit-struktur-id`, `slice-tests-ueberlebende-mutanten`,
`slice-tests-ueberlebende-mutanten-driving`, `slice-harness-lint-warten-ohne-frist`,
`slice-harness-mutation`; die Reihenfolge
der Wellen-Slices steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md). Die Slices zwischen `slice-harness-lint` und diesem sind
keine technische Abhängigkeit, sondern die Reihenfolge des Nutzers. Grund für den Platz nach `slice-harness-lint`: Teil 1 und
Teil 3 gelten erst, wenn das Lint-Gate steht und `testpackage` überall scharf ist, und
das ist es mit `slice-harness-lint`; dessen Gegenprobe trägt dann alle Fälle, die die
Deklaration deckt.
Die Randformen aus §6 hat der Architect am 2026-10-17 bestätigt und die offenen in
[`SPEC-048`](../../../../harness/vertraege/SPEC-048-abdeckung.md) entschieden
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Größe** (`AGENTS.md` §3.14): Metrik ist die Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools/test tools/harness/lint-gegenprobe.sh
harness/README.md harness/sensors/abdeckung-check.md harness/vertraege/SPEC-048-abdeckung.md docs/user .d-check.yml`,
Tests eingeschlossen. Basis: `cbeb3381`, der Elternteil des Anspruchs-Commits (`next` → `in-progress`, `1ab4f9d6`). Schätzung (vom Architect am 2026-10-17 auf 450 bis 700 korrigiert, Mitte 550, Obergrenze
700): `abdeckung.sh` rund 170 Zeilen (Gate-Kopf, Teile und Querprüfungen, Tabelle, Darstellung der Teile,
Ausgang 2), die Gegenprobe mit dem Doppelten davon (rund 300 bei rund 30 Fällen), Tabellen unter
`docs/user`, Gate-Index, Sensor-Datei, Kopf der Lint-Gegenprobe und `.d-check.yml` zusammen rund 50;
`SPEC-048` liegt mit dem Anspruch im Bestand; ihr Pfad steht in der Metrik, die Entscheidungen des Architect nach der Umsetzung (32 Zeilen) zählen mit, eine ADR entfällt (§1). Meldeschwelle
und Marke bleiben. **Meldeschwelle 800:** Wer im Lauf
darüber liegt oder liegen wird, meldet es dem Planner, bevor er weiterbaut. **Marke 1000:**
darüber gilt die Rückführung unten, ohne Ermessen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe der Metrik liegt über 1000,
  oder die Erweiterung ändert die Form der bestehenden Tabellen oder der Test-Deklarationen so weit,
  dass Skript, Gegenprobe und Tabellen nicht in einer Review-Sitzung prüfbar sind; dann zuerst die
  Erweiterung ohne Deklaration, die Deklaration in einem eigenen Slice. Der Zerlegungs-Vorschlag
  `slice-harness-rtm-nachweisarten` ist gestrichen (Architect, 2026-10-17): Er trüge nur die
  reichere RTM-Darstellung, die §1 nicht empfiehlt; verlangt sie der Nutzer, legt der Planner ihn an.
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

**Randformen des Vertrags** (`AGENTS.md` §3.12). **Entschieden vom Architect am 2026-10-17 in
[`SPEC-048`](../../../../harness/vertraege/SPEC-048-abdeckung.md)** (Ort nach `MR-002`; eine neue Kennung ist nicht
nötig, `SPEC-057` ist an `slice-harness-integration-volumes` vergeben, `SPEC-056` an
`slice-harness-coverage` vorgemerkt); die ADR trägt Entscheidung und Gründe ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)).
Jede zugesagte Randform bekommt einen Fall in `make abdeckung-gegenprobe` und eine Mutation, die nur ihn
bricht (`AGENTS.md` §3.10, §3.11, je Element einer Aufzählung ein Fall). Eine nicht genannte Randform
entscheidet der Implementer nicht, er gibt sie zurück.

| Randform | Entscheidung | Stelle in `SPEC-048` |
|---|---|---|
| Pfad-Schreibweise, Zahlform | `Messung-<i>-von-<n>`, i und n je `[1-9]`, i ≤ n, 2 ≤ n; mehrstellig, `01`, Vorzeichen, Leerraum sind Fehlformen (je ein Fall); n ist höchstens 9, damit die Prüfung der Vollständigkeit endlich und lesbar bleibt | 2, 8 |
| Fehlformen der Teile | i > n; n = 1; verschiedene n für dieselbe Anforderung; `Messung` neben geteilter Messung derselben Anforderung, gleich in welcher Datei | 8 |
| Derselbe Teil zweimal | zulässig, zwei Zeilen wie zwei Tests am selben Pfad | 8 |
| Reihenfolge mehrerer Gründe | Form, Schreibweise einschließlich Paarung, Wertebereich, Anforderung unbekannt, Widerspruch zwischen Deklarationen; ein Element nennt einen Grund; ein Widerspruch steht an der späteren Fundstelle (Dateipfad, Zeile) und nennt die frühere | 6 |
| Nachweisart Gate nur für `LH-QA`, `LH-RB` | an einer `LH-FA` eine Fehlform | 7 |
| Gelesen werden | Dateien auf `-gegenprobe.sh` in jeder Tiefe unter `tools/` an der Wurzel; Kopf = zusammenhängende `#`-Zeilen ab Zeile 1 | 7 |
| Form des Kopfes | `# Abdeckung: <Anforderung>/<Pfad>, … — <Text>` am Zeilenanfang; Trenner „ — “ Pflicht; Fortsetzung mit `# `, Ende an `#` allein, nächster Deklaration oder Kopfende; mehrere Deklarationen zulässig; eingerückt, ohne Trenner und `# Abdeckung: ` nach dem Kopf sind Fehlformen (je ein Fall) | 7 |
| Deklaration außerhalb des Musters | überlesen, die Anforderung steht dann höchstens *teilweise* da, nie grün (akzeptiertes Negativ, Grund dort) | 7 |
| Nachweis-Spalte, Tabelle, Dateiname | Skriptpfad; eigene Tabelle `docs/user/abdeckung-gate.md` ohne Spalte *Ort*; Kopftexte nicht zugesagt | 4 |
| Darstellung der Teile in der Gesamtsicht | alle Teile 1 bis n, aufsteigend, `<i>/<n> <Arten>` oder `<i>/<n> —`, durch „; “ getrennt; Nachweisarten in der Folge `E2E, Gate, Unit` | 4 |
| Vollständigkeit | alle Teile in mindestens einer Deklaration, gleich welcher Art; `abdeckung-vollstaendig.md` führt sie erst dann | 4 |
| Geteilte Messung an Go-Tests | zulässig ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) „jeder Nachweis“), die Tabellen E2E und Unit zeigen den Pfad wie deklariert; an einer `LH-FA` eine Fehlform | 2, 8 |
| Darstellung in der RTM | eine Quelle, Label `Abdeckung`; keine neue ADR (Empfehlung, §1) | — (`.d-check.yml`) |
| Fehlendes `docs/user/` | Exit 2 mit Meldung, nichts geprüft, nichts angelegt | 6 |
| `test/integration/` nur an der Wurzel | so im Vertrag; ein tieferes Verzeichnis gleichen Namens ergibt Unit (ein Fall mit Deklaration) | 3 |
| Lastenheft ohne Anforderung | Exit 2 mit Meldung statt Exit 1 ohne Zeile; der Fall `lastenheft-ohne-anforderung` und die Datei `harness/sensors/abdeckung-check.md` ziehen nach | 6 |
| Messmethode (4) von LH-QA-07 | kein eigener Teil, Gesamtzahl 3 ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)); Urteil dessen, der deklariert | Grenze |
| Anschluss an `GATE_CHECKS` | prüft das Skript nicht; Grenze im Skriptkopf, Review | Grenze |
| Aufrufarten von `make` | nicht neu: Das Ziel schreibt Tabellen und glaubt keinen Stand; kein Eintrag | — |

**Hinweise an den Implementer** (kein Vertrag):

- *Fixtures der Gegenprobe:* `tools/test/abdeckung-gegenprobe.sh` liegt selbst unter dem gelesenen
  Muster. Eine Zeile `# Abdeckung: …` am Zeilenanfang dieser Datei außerhalb ihres Kopfes (etwa in einem
  Heredoc) wäre eine Fehlform (Punkt 7); die Fixtures entstehen wie heute mit `printf` in einer Zeile.
  Der Kopf dieser Datei trägt keine Deklaration (sie belegt keine Anforderung).
- *Struktur:* Die Datei wächst um rund 30 Fälle; je Zusage ein Fall mit der Mutation, die nur ihn bricht
  (Aufzählungen: jede Zahlform, jede Fehlform des Kopfes, jede Fehlform der Teile ein eigener Fall; die
  Reihenfolge der Gründe je Paar ein Fall). Die Fehlformen laufen über die vorhandene Funktion
  `abgelehnt`, die Tabellenzeilen über `zeilen`. Die Mutationen gehen gegen `tools/test/abdeckung.sh` in einer
  Kopie; trifft eine Ersetzung nichts oder mehr als eine Stelle, ist die Mutation ungültig.
- *Kosten:* Alle Fälle laufen in Temp-Bäumen ohne Docker und ohne Netz, je Fall ein Aufruf von `bash`;
  die Laufzeit des Ziels bleibt im Sekundenbereich.
- *Tabellen:* `docs/user/abdeckung-*.md` entstehen mit `make abdeckung`, nie von Hand.
  `make abdeckung-check` ist danach grün. Vor dem Eintrag der Deklaration an der Lint-Gegenprobe
  (`tools/harness/lint-gegenprobe.sh`) steht der Abgleich aus dem Risiko *Deklaration sagt mehr zu*:
  vom Architect am 2026-10-17 gemacht, Messmethode 1 deckt `p0-grundlauf`, `p3-*`, `p4-*`, `p5-*` und
  das Ausnahmen-Profil (`p8-*`, `p9-*`), Messmethode 3 beide Bedingungen (`p7-whitebox-kern`,
  `-driven`, `-pgwire`, `-einstieg`; `p7-internal-test`, `p7-test-in-bruecke-*`); der Implementer
  prüft die Namen gegen die Datei, wenn er den Kopf schreibt.
- *Aufrufarten von `make`:* Kein Stand, kein Hook; `abdeckung-check` und `abdeckung-gegenprobe` hängen
  wie bisher an `GATE_CHECKS`.

**Prüfung des Architect nach der Umsetzung** (2026-10-17, Rückgabe des Implementers aus §7; entschieden in
[`SPEC-048`](../../../../harness/vertraege/SPEC-048-abdeckung.md) Punkt 6 bis 8, Code unverändert wörtlich richtig):

| Randform | Entscheidung | Rückgabe an den Implementer |
|---|---|---|
| Lastenheft nur mit `### LH-XX-01` | Ausgang 2 prüft die Zeile `### LH-`, nicht die Form; Deklarationen scheitern an (4) mit Exit 1; akzeptiertes Negativ (ein Lastenheft ohne Anforderung fällt beim ersten Test auf) | Fall `lastenheft-fremde-kennung` in `tools/test/abdeckung-gegenprobe.sh`: Exit 1 mit „steht nicht im Lastenheft“, nicht 2 |
| Widerspruch gegen mehrere frühere | Der Fehler nennt die erste festlegende Deklaration; ein Widersprechender legt nichts fest; weitere mit dem widersprechenden n stehen ebenfalls im Fehler | Fall: zwei frühere Deklarationen mit n = 3, danach n = 2, die Meldung nennt die erste (Mutant „zuletzt gelesene“ wird rot); ebenso `Messung` zweimal, dann ein Teil; Satz im Kopf der Gegenprobe („Ebenso nicht gezeigt …“) streichen |
| Zeile im Fehler einer Gate-Deklaration | Zeile des Anfangs `# Abdeckung: …`, bei eingerückt und nach dem Kopf die Zeile selbst; gepinnt ist nur `datei:zeile`, nicht der Wortlaut | Fälle `gate-zeile` genügen; je ein Fall mit mehrzeiliger Deklaration (Anfang, nicht Fortsetzung) und eingerückt/nach dem Kopf prüfen, sonst ergänzen |
| Eingerückt an jeder Stelle | Fehlform in jeder Stelle der Datei, Leerzeichen und Tab; anderer Leerraum nicht geprüft | Fall `gate-eingerueckt-nach-kopf` (Leerzeichen und Tab nach dem Kopf); Mutant „nur im Kopf“ wird rot |
| `# ` ohne Text | beendet die Fortsetzung wie `#` allein (Zeile aus `#` und höchstens Leerzeichen); jede andere Kopfzeile ohne `# ` beendet sie ebenfalls | Fälle sind in `gate-kopf` vorhanden (§7, Mutationen); keine |
| Symlink auf `-gegenprobe.sh` | nicht gelesen (reguläre Dateien); `-type f` ist damit nicht äquivalent | Fall `gate-nicht-gelesen` mit Symlink ergänzen |

**Äquivalente Mutanten aus §7, geführt ohne Zusage** (`AGENTS.md` §3.11; der Vertrag sagt nichts, was sie brechen):
Reihenfolge von i > n und n = 1 (Vertrag: Meldung nennt einen der beiden, nicht welchen, Punkt 8); Ende der Fortsetzung am
Kopfende (die Zusage ist das Ergebnis, nicht der Zweig); `LH-FA`-Zweig allein (Wortlaut der Meldung nicht zugesagt, Punkt 6);
Wahl der früheren Deklaration (nicht äquivalent, seit heute zugesagt, siehe oben); `-type f` bei einem Verzeichnis (kein Fall des
Vertrags, bei einem Symlink nicht äquivalent, siehe oben). `LC_ALL=C` entfernen bleibt grün, weil die Umgebung des Laufs C ist:
akzeptiertes Negativ, die Zusage „bytegenau“ trägt der Mutant mit anderer Sortierung, das Skript setzt das Locale selbst.

**Risiken:**

- **Falsche Gesamtzahl n** — Gesamtzahl und Zuordnung eines Teils zur Messmethode sind
  Urteil dessen, der deklariert; das Skript fängt nur verschiedene n für dieselbe
  Anforderung, eine einheitlich falsche Zahl fängt nur das Review
  ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), Re-Evaluierungs-
  Trigger). — **Ausgang (Planner, Closure):** weiter offen → Register
  `BEO-REPO/gesamtzahl-geteilter-messung-einheitlich-falsch` (1×, Beleg
  `evidence/slice-harness-abdeckung-gate.md`); kein Auftreten, Teil 2 von `LH-QA-07` kommt mit
  `slice-harness-coverage`.
- **Deklaration sagt mehr zu, als die Gegenprobe prüft** — Teil 3 an der
  Lint-Gegenprobe trägt nur, wenn sie für beide Bedingungen von Messmethode 3 einen
  Rot-Fall führt: je Paketgruppe einen White-Box-Test (der Test liegt in der Einheit)
  und die Fälle zur Brücke aus `slice-harness-lint` (Zugriff auf Internes an der
  Brücke vorbei, Test in der Brückendatei)
  (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, 11×). Der Abgleich gegen die
  Fälle der Lint-Gegenprobe ist vom Architect gemacht (§6, Hinweise); der Implementer prüft die Namen
  beim Schreiben des Kopfs. — **Ausgang (Planner, Closure):** entfallen: Jedes Präfix des Kopfs der
  Lint-Gegenprobe hat Fälle in der Datei (§7 *Abgleich der Namen*, von der Verifikation nachgemessen),
  Teil 1 nennt bewusst `p6-*` statt `p9-*` (§7), und Teil 2 steht nicht dort; Review und Verifikation
  fanden keine Überzusage.
- **Bestehende Prüfung fällt weg** — die Erweiterung tritt neben die Prüfung der
  Test-Deklarationen, nicht an ihre Stelle; die vorhandenen Fälle der Gegenprobe
  bleiben (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). — **Ausgang (Planner, Closure):** entfallen:
  Die vorhandenen Fälle bleiben, nur `lastenheft-ohne-anforderung` wurde auf Exit 2 gedreht (Vertrag
  Punkt 6); die Gegenprobe wuchs von 13 auf 56 Aufrufe, und der Mutant `exit 0` ist rot (Verifikation).
- **RTM zählt eine Teilabdeckung als belegt** — die einzige Quelle bleibt
  `abdeckung-vollstaendig.md`, in der eine geteilte Messung erst mit allen Teilen steht. Bricht das,
  sieht man es nur im Lauf: Der Implementer prüft mit `make doc-trace`, dass LH-QA-07 mit zwei von
  drei Teilen nicht belegt erscheint und eine vollständige Anforderung mit dem Label `Abdeckung`
  belegt (Beleg in §7). — **Ausgang (Planner, Closure):** entfallen: Am Repo steht LH-QA-07 mit `—`;
  im Gegenbeweis der Verifikation (Kopie mit allen drei Teilen, mit Teil 1 entfernt) zeigt
  `make doc-trace` sie nur bei Vollständigkeit mit `Abdeckung`, und die Gegenprobe fängt den Mutanten
  `t-voll-irgendein-teil`.
- **Andere Gegenproben werden gelesen** — `tools/**/*-gegenprobe.sh` trifft auch die
  vorhandenen Gegenproben ohne Deklaration (`tools/harness/*-gegenprobe.sh`,
  `tools/hook/commit-msg-gegenprobe.sh`, `tools/test/postgres-versionen-gegenprobe.sh` und, aus
  `slice-harness-integration-volumes`, `tools/test/integration-volumes-gegenprobe.sh`); sie dürfen
  weder Fehler noch Zeilen erzeugen, und ein Fall der Gegenprobe zeigt es an einer Datei ohne
  Deklaration im Muster. — **Ausgang (Planner, Closure):** entfallen: Der Fall
  `gegenprobe-ohne-deklaration` nimmt eine Datei nach dem Vorbild von
  `tools/test/integration-volumes-gegenprobe.sh` ohne Fehler und ohne Zeile an, und `make
  abdeckung-check` ist am Repo grün.

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

**Größe** (§4, Metrik `git diff --numstat cbeb3381..<letzter Commit> -- tools/test tools/harness/lint-gegenprobe.sh harness/README.md harness/sensors/abdeckung-check.md harness/vertraege/SPEC-048-abdeckung.md docs/user .d-check.yml`): Summe aus Einfügungen und Löschungen **909** am Stand `74446efb` (Skript 273, Gegenprobe 526, `harness/sensors/abdeckung-check.md` 32, `SPEC-048` 32 aus der Entscheidung des Architect, Tabellen unter `docs/user` 26, `.d-check.yml` 11, `harness/README.md` 6, Kopf der Lint-Gegenprobe 3; `harness/mk/abdeckung.mk` liegt außerhalb der Pfade). Über der Meldeschwelle 800 (gemeldet), unter der Marke 1000; die Schätzung 450 bis 700 ist überschritten. Schichten zwei, Liefer-Punkte drei; die E2E- und die Unit-Tabelle sind byte-gleich geblieben.

**Vorher und nachher** (`make abdeckung`): `abdeckung-gesamt.md` führt LH-QA-07 vorher nicht, nachher `| LH-QA-07 | n/a | n/a | n/a | 1/3 Gate; 2/3 —; 3/3 Gate | teilweise |`; `abdeckung-vollstaendig.md` führt sie nicht; neu `abdeckung-gate.md` mit zwei Zeilen (`Messung-1-von-3`, `Messung-3-von-3`, Nachweis `tools/harness/lint-gegenprobe.sh`). RTM (`make doc-trace`, `Coverage`-Spalte): LH-QA-07 `—`, LH-QA-02, LH-QA-06 und LH-RB-01 `Abdeckung`. Mit einem Mutanten, der eine geteilte Messung schon bei einem belegten Teil als vollständig zählt (Kopie des Repos, `return 1` → `return 0` in `vollstaendig`), schreibt `make abdeckung` LH-QA-07 in `abdeckung-vollstaendig.md`, und `d-check --trace` zeigt sie mit `Abdeckung`: Das Risiko *RTM zählt eine Teilabdeckung als belegt* ist damit sichtbar, und die Gegenprobe fängt den Mutanten (`t-voll-irgendein-teil`, unten).

**Abgleich der Namen** (Risiko *Deklaration sagt mehr zu*): Der Kopf von `tools/harness/lint-gegenprobe.sh` nennt `p0-grundlauf`, `p3-*`, `p4-*`, `p5-*`, `p6-*`, `p8-*` (Teil 1) und `p7-whitebox-{kern,driven,pgwire,einstieg}`, `p7-internal-test`, `p7-endet-auf-export-test`, `p7-test-in-bruecke-*` (Teil 3); jedes Präfix hat Fälle in der Datei (`grep` je Präfix: 1, 17, 18, 36, 12, 48 und 4, 1, 1, 1). Teil 2 steht bewusst nicht dort (§1).

**Mutationen** (Weg: je Mutant eine Kopie von `tools/test/` im Scratchpad, genau eine Ersetzung, die Zahl der Treffer geprüft, die unveränderte Gegenprobe gegen den Mutanten; die Kopie ist ein frischer Pfad, kein Build-Kontext). Die Spalte „roter Test“ nennt die Fälle, die rot wurden; unverändert sind alle grün.

| Zusage | Mutation | roter Test |
|---|---|---|
| Zahlform: mehrstellig, `01`, Vorzeichen, Leerraum, je an i und n, `0` (Punkt 8) | Muster mit `[1-9][0-9]*`, `0?[1-9]`, `[+]?[1-9]`, ` ?[1-9]` an i bzw. n; `[0-9]`; Anker `$` entfernt | `teil-zahlform-<Form>` und `gate-teil-zahlform-<Form>` der jeweiligen Form, bei `01` und `0` dazu weitere Fälle derselben Form (Mutanten `zf2-*`, je Form einer) |
| Zahlform an `LH-RB` | Muster nur für `LH-QA` (`zf-nur-qa-teil`) | `teil-i-ueber-n-rb` und die Bäume mit `LH-RB-01` |
| geteilte Messung an `LH-FA` Fehlform | an einer `LH-FA` zugelassen | `teil-fa` |
| i > n | `a > b + 1`; Prüfung entfernt; `continue` entfernt | `teil-i-ueber-n`, `gate-teil-i-ueber-n`, `grund-3-vor-4` |
| n = 1 | `b < 1`; Prüfung entfernt; `continue` entfernt | `teil-n-eins`, `gate-teil-n-eins`, `grund-5-ohne-3` |
| i = n zulässig | `a >= b` | `gate-kopf`, `teil-gut`, `teile`, … (viele) |
| verschiedene n für dieselbe Anforderung | Prüfung entfernt; n global statt je Anforderung | `widerspruch-n`, `widerspruch-bytegenau`, `widerspruch-zeile`, `widerspruch-gate-*` bzw. `teil-gut`, `teile*` |
| `Messung` neben Teil, Reihenfolge Messung zuerst | Prüfung entfernt | `widerspruch-messung-dann-teil`, `widerspruch-messung-gate` |
| `Messung` neben Teil, Teil zuerst | Prüfung entfernt | `widerspruch-teil-dann-messung` |
| Widerspruch an der späteren Fundstelle, nennt die frühere | Orte vertauscht | `widerspruch-n`, `-bytegenau`, `-zeile`, `-gate-frueher`, `-gate-spaeter`, `-gate-gate` |
| Ordnung nach Dateipfad | Schlüssel entfernt; locale-abhängig (`f`) | `widerspruch-gate-spaeter`, `widerspruch-teil-dann-messung`; `widerspruch-bytegenau` |
| Ordnung nach Zeile numerisch | lexikalisch | `widerspruch-zeile` |
| Ordnung der Elemente einer Deklaration | Schlüssel entfernt; rückwärts | `widerspruch-element-rueckwaerts`; `widerspruch-element` |
| Reihenfolge der Gründe: Form vor Elementen | Elemente nach der Formprüfung noch gelesen; Formprüfung entfernt | `grund-1-vor-2-bis-4`, `grund-5-ohne-1` |
| (4) Anforderung unbekannt nach (3) | Prüfung (4) vor Wertebereich | `grund-3-vor-4` |
| (5) betrachtet nur Elemente ohne Grund aus (1) bis (4) | `continue` nach (3) bzw. (4) entfernt, ungültiges Element weitergereicht | `grund-5-ohne-3`, `grund-5-ohne-4`, `grund-3-vor-4` |
| Gate-Deklaration: nur `LH-QA`, `LH-RB` | `LH-FA` zugelassen (Zweig entfernt und Muster geweitet) | `gate-fa`, `gate-fa-messung`, `gate-fa-teil`, `gate-pfad-happy` |
| gelesen werden Dateien auf `-gegenprobe.sh` | Muster ohne Bindestrich; mit Nachsilbe; `tools/*.sh` | `gate-nicht-gelesen tools/gegenprobe.sh`, `… tools/x-gegenprobe.sh.bak`; die dritte Mutation fällt auf `gate-nicht-gelesen tools/x.sh` |
| nur unter `tools/` an der Wurzel | Suche ab der Wurzel des Repos | `gate-nicht-gelesen harness/…`, `… sub/tools/…` |
| in jeder Tiefe | `-maxdepth 1` | `gate-tiefer` |
| Kopf endet an der ersten Zeile ohne `#` | Kopf endet erst an einer Leerzeile | `gate-nach-dem-kopf`, `gate-nach-dem-kopf-code` |
| Shebang-Zeile gehört zum Kopf; Deklaration in Zeile 1 und 2 | Shebang Pflicht; Zeile 1 übersprungen | `gate-tiefer` |
| eingerückt (Leerzeichen, Tab) | Prüfung nur im Kopf; nur Leerzeichen | `gate-eingerueckt`, `gate-eingerueckt-tab` |
| Deklaration ohne „ — “ | Fehlermeldung nur für Tests | `gate-ohne-trenner` |
| `Abdeckung:` am Zeilenanfang nach dem Kopf | Zweig entfernt | `gate-nach-dem-kopf`, `gate-nach-leerzeile`, `gate-nach-dem-kopf-code` |
| Fortsetzung `# ` + Text, Ende an `#` allein, an `# ` ohne Text, an `#` ohne Leerzeichen, an der nächsten Deklaration | Fortsetzung auch ohne Leerzeichen (`^#`); auch ohne Text; `^#` ohne Ausnahme für `#` allein; Ende an der nächsten Deklaration entfernt | `gate-kopf` (jede der drei Fortsetzungs-Mutationen); `gate-kopf`, `teile-gate`, `teile-gesamt` (nächste Deklaration) |
| Zeile der Fehlermeldung ist die Zeile der Deklaration | `NR - 1` | `gate-zeile`, `gate-zwei-dateien`, `widerspruch-gate-*` |
| Fortsetzung bei Tests | ohne Leerzeichen angefügt | `gueltig-maskierung` |
| `|` im Text maskiert | `gsub` entfernt | `gueltig-maskierung`, `gate-gutform` |
| Tabelle Gate: eigene Datei ohne Spalte *Ort*, nur Gate-Zeilen | Spalte *Ort*; Unit-Zeilen in die Gate-Tabelle | `inhalt-gate-leer`, `gegenprobe-ohne-deklaration`, `teile-gate` |
| Tabelle Gate sortiert nach Anforderung, Pfad, Nachweis | Nachweis aus dem Schlüssel; Pfad aus dem Schlüssel | `teile-gate` (beide) |
| Teile in der Gesamtsicht: alle 1 bis n, aufsteigend, `<i>/<n> <Arten>` oder `<i>/<n> —`, durch „; “ | Schleife bis n − 1; `—` entfernt; Trenner „, “ | `teile-gesamt`, `teile-voll-gesamt`, `teile-reihenfolge-gesamt`, `teil-fehlt-*-gesamt` |
| Nachweisarten in der Folge `E2E, Gate, Unit` | `sort -ru` | `inhalt-gesamt`, `teile-gesamt`, `teile-voll-gesamt` |
| Vollständig: jeder Teil 1 bis n (erster, mittlerer, letzter, neunter) | Schleife bis n − 1; irgendein Teil genügt (`t-voll-irgendein-teil`); immer Teil 1; Pfad `Messung` als einziger Maßstab | `teil-fehlt-erster`, `-mittlerer`, `-letzter`, `-neunter`, `-erster-von-neun`, `teile-vollstaendig`, `teil-neun-voll` |
| Gesamtzahl wird aus den Deklarationen gelesen | fest 2 | `teile-gesamt`, `teil-fehlt-*` |
| E2E-Verzeichnis nur an der Wurzel | Muster `*test/integration/*` | `e2e-tiefer-e2e`, `e2e-tiefer-unit` |
| Ausgang 2: Lastenheft fehlt; ohne Überschrift; ohne `docs/user/` | Exit 1 an jeder der drei Stellen; Prüfung der Überschrift entfernt | `ohne-lastenheft`, `lastenheft-ohne-anforderung`, `ohne-docs-user` (ohne und mit `--check`), `ausgang2-ueberschrift-vor-docs` |
| Reihenfolge der drei Gründe | `docs/user/` vor dem Lastenheft; vor der Überschrift | `ausgang2-lastenheft-zuerst`, `ausgang2-ueberschrift-vor-docs` |
| `make abdeckung` legt `docs/user/` nicht an | `mkdir -p` statt Abbruch | `ohne-docs-user`, `ohne-docs-user --check` |
| Fünf Tabellen geschrieben und geprüft | Gate-Tabelle aus der Schleife genommen | `schreiben-fehlend`, `inhalt-gate-leer`, `check-gate-veraendert`, `check-gate-fehlt` |
| Ausgang 1 bei Fehlform; bei veralteter Tabelle (Verdrahtung des Aufrufers im Skript) | `exit 0` statt `exit 1` nach den Meldungen; `exit 0` am Ende | fast jeder Fall von `abgelehnt`; `check-veraltet`, `check-zwei-veraltet`, `check-gate-*` |

**Äquivalente und nicht unterscheidbare Mutanten** (grün, eingeordnet; die Grenze trägt der Kopf der Gegenprobe):

- Reihenfolge von i > n und n = 1 untereinander (beide Gründe (3), die Meldung nennt einen davon): Der Vertrag ordnet sie nicht.
- Das Ende der Fortsetzung am Kopfende (`abschluss()` am Kopfende entfernt, oder die Fortsetzung auch nach dem Kopf zugelassen, einzeln und zusammen): Die erste Zeile nach dem Kopf ist keine `#`-Zeile und beendet die Fortsetzung selbst; der Fall `gate-kopfende` zeigt das Ergebnis, nicht den Zweig.
- Der Zweig für `LH-FA` in der Gate-Deklaration allein (ohne das Muster zu weiten): `LH-FA` fällt dann am Muster, mit einer anderen Meldung; der Vertrag nennt die Meldung nicht.
- (entfallen) Welche frühere Deklaration ein Widerspruch nennt: seit der Entscheidung des Architect (SPEC-048 Punkt 6) fangen `widerspruch-erste-n`, `-erste-messung`, `-erster-teil` den Mutanten „zuletzt gelesene“.
- `LC_ALL=C` (Punkt 6, bytegenau): Die Umgebung des Laufs setzt kein anderes Locale; der Mutant mit `en_US.UTF-8` im Sortieren wird rot (`widerspruch-bytegenau`), das Entfernen von `export LC_ALL=C` nicht, weil die Umgebung des Laufs C ist.

**Randformen nach der Entscheidung des Architect** (Commit `cacaf26a`, SPEC-048 Punkte 6 bis 8): Nachgearbeitet nur in der Gegenprobe. Zusage · Mutation · roter Test:

| Zusage | Mutation | roter Test |
|---|---|---|
| Lastenheft nur mit `### LH-XX-01`: Exit 1 mit „steht nicht im Lastenheft“, nicht 2 | Ausgang 2 prüft `### LH-(FA|QA|RB)-` | `lastenheft-fremde-kennung` |
| Widerspruch nennt die erste Deklaration, die n festlegt | n und Ort bei jedem gültigen Element neu gesetzt | `widerspruch-erste-n`, `widerspruch-erster-teil` |
| bei `Messung` neben einem Teil die erste Deklaration des anderen Pfads | Ort von `Messung` bei jeder Deklaration neu gesetzt | `widerspruch-erste-messung` |
| eingerückt (Leerzeichen, Tab) an jeder Stelle, Zeile selbst | Prüfung nur im Kopf; nur Leerzeichen; Zeile − 1 | `gate-eingerueckt-nach-kopf`, `-tab`, `gate-eingerueckt-zeile` (Leerzeichen und Tab: `…-tab`) |
| Zeile einer mehrzeiligen Deklaration ist ihr Anfang | Zeile bei jeder Fortsetzungszeile neu gesetzt | `gate-zeile-mehrzeilig` |
| Symlink auf eine Gegenprobe wird nicht gelesen | `-type f` entfernt | `gate-nicht-gelesen-symlink` |

Review F-1 und F-2 nachgezogen (Zusage · Mutation · roter Test): eingerückte Test-Deklaration mit Leerzeichen · Muster auf `[\t]+` · `eingerueckt-leerzeichen`; Vorzeichen `-` an i, an n · `[-]?` am Ziffernmuster · `teil-zahlform-Messung--1-von-2`, `…-1-von--2`; Tab als Leerraum an i, an n · `\t?` · `…-\t1-von-2`, `…-1-von-\t2`; nachgestellter Leerraum (Leerzeichen, Tab) hinter n · `[ ]?` bzw. `\t?` vor `$` · `…-1-von-2 `, `…-1-von-2\t` (je auch als Gate). Dabei gefunden und behoben: `cut -f2` kürzte die Meldung an einem Tab in der Deklaration; jetzt `sed`.

Neu fallen auch die Fälle `gate-nach-kopf-zeile` (Zeile 3). Die Fragen 1 bis 5 sind beantwortet; die Äquivalenz-Liste oben und der Kopf der Gegenprobe führen die verbliebenen Fälle.

**Bewusste Auswahl der Deklaration:** Teil 1 der Lint-Gegenprobe nennt `p6-*` (kein `//nolint`, Ausnahmen zentral) und nicht `p9-*` wie der Hinweis des Architect in §6; Messmethode 1 verlangt „Ausnahmen stehen zentral und begründet, nicht im Code“, das trägt `p6-*`, die Fälle `p9-*` betreffen die Ausgabe des Gates.

**Akzeptiertes Negativ (Verdrahtung):** Kein Fall der Gegenprobe zeigt, dass das Rezept von `abdeckung-check` in `harness/mk/abdeckung.mk` das Skript mit `--check` ruft und dass beide Ziele an `GATE_CHECKS` hängen (Bestand; in diesem Slice berührt nur der Kommentar die Datei). Grund: Beides ist Verdrahtung des Treibers `make`, die Gegenprobe prüft das Skript; [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Konsequenzen) nennt den Anschluss als Review. `harness/mk/abdeckung.mk` und `harness/sensors/abdeckung-check.md` sagen das so.

**Läufe:** `make abdeckung-check`, `make abdeckung-gegenprobe` (rund 6 s), `make docs-check` (992 Dateien, 0 Befunde), `make kopf-check`, `make doc-trace` am Stand dieses Commits grün; `make lint-gegenprobe` und `make lint` laufen im vollen `make gates`. `make gates` am Endstand (Stand der Arbeitsbaum-Datei `gates-passed.dateiliste` vom 2026-10-17 16:00): Exit 0, `abdeckung-check`, `abdeckung-gegenprobe`, `lint`, `lint-gegenprobe` und alle übrigen Gates grün; ein erster Lauf als Hintergrundjob der Shell scheiterte an `container-gegenprobe` (SIGINT wird von einem so gestarteten Prozess geerbt ignoriert, Fälle `p9-*INT`), der zweite, über das Werkzeug gestartet, lief grün; keine Flake (`p7-int-treiber`, `p9-ungenutzt-allein`) trat auf.

### Closure-Notiz

Stand: Planner am 2026-10-17, vor dem `git mv` nach `done/`. Review (2 MEDIUM, 2 LOW, 1 INFO) und Verifikation (bestanden, 2 LOW) liegen unter `docs/reviews/` (`2026-10-17-review-slice-harness-abdeckung-gate.md`, `2026-10-17-verify-slice-harness-abdeckung-gate.md`). `make gates` am Endstand `63619f98`: Exit 0, vom Hauptagenten über das Werkzeug gefahren; seither ändert sich nur dieser Plan und `docs/plan/planning/` (inert).

- **Was hat funktioniert:** Die Randformen standen vor dem Code im Vertrag, die Rückgabe des Implementers erreichte den Architect (`cacaf26a`), und die Gegenprobe wuchs von 13 auf 56 Aufrufe mit Mutation je Zusage. LH-QA-07 steht in der Gesamtsicht als *teilweise* (`1/3 Gate; 2/3 —; 3/3 Gate`), nie grün; die Tabellen E2E und Unit blieben byte-gleich.
- **Was ging anders als geplant:** Die Schätzung 450 bis 700 lag unter der Messung 909 (Meldeschwelle 800 erreicht, Marke 1000 nicht). Das Review fand zwei Lücken, die weder Implementer noch Verifier fanden: zwei Aufzählungen der Zusage (Leerzeichen und Tab bei Test-Deklarationen; Vorzeichen und Leerraum in der Zahlform) ohne Fall je Element. Die neuen Fälle fanden einen Code-Defekt, den kein Fall vorher zeigte: `cut -f2` kürzte die Fehlermeldung an einem Tab in der Deklaration (jetzt `sed`).
- **Ausgänge der Findings:** F-1 und F-2 (MEDIUM) behoben in `74446efb`: je ein Fall und eine Mutation je Element. F-3 (Größe in §7 veraltet) behoben in `33ab3e4b` (909 am Endstand). F-4 (Verdrahtung `--check` und `GATE_CHECKS` ohne Fall) akzeptiertes Negativ, Fundort §7 *Akzeptiertes Negativ (Verdrahtung)* und `harness/sensors/abdeckung-check.md`: Beides ist Verdrahtung des Treibers `make`, [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Konsequenzen) nennt den Anschluss als Review. F-5 (INFO, Kommentar nennt `stat`) behoben in `74446efb`. Verifikation F-1 (Größe) behoben wie Review F-3, F-2 (`p6-*` statt `p9-*`) in §7 *Bewusste Auswahl* begründet.
- **Steering-Loop-Eintrag:** Neuer Sensor: `make abdeckung-check` liest Gate-Deklarationen im Kopf von Gegenproben (Nachweisart *Gate*) und den Pfad `Messung-<i>-von-<n>`; `make abdeckung-gegenprobe` führt 56 Aufrufe mit Mutation je Zusage — liegt in `tools/test/abdeckung.sh` und `tools/test/abdeckung-gegenprobe.sh`, Vertrag `SPEC-048` Punkte 6 bis 8. Keine geschärfte Regel von diesem Lauf: Die Lücken F-1 und F-2 trafen §3.10 (Aufzählung je Element), die schon gilt; sie gehen als Wiederauftreten ins Register, `AGENTS.md` bleibt unberührt (Entscheidung des Nutzers). Die Klassen erreichten alle die Schwelle schon vorher und stehen auf *verkörpert* beziehungsweise *geplant*.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-harness-abdeckung-gate.md` ergänzt in `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` (49×), `BEO-REPO/plan-folgt-korrektur-nicht/` (62×), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (74×), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests/` (12×, Stand `geplant` bleibt) und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben/` (21×); `BEO-REPO/randform-wellenlos-ohne-architect-vor-code/` Kein Auftreten (1×); neu angelegt `BEO-REPO/gesamtzahl-geteilter-messung-einheitlich-falsch/` (1×, offen) und `BEO-REPO/lauf-ignoriert-sigint-nach-nohup/` (1×, offen). Kein Eintrag wechselte seinen Ausgang.
- **Parallel-Lauf (Befund, keine Regeländerung):** Dieser Slice lief auf dem Hauptzweig `main`, zugleich mit `slice-harness-devcontainer-austausch-repo` im Worktree `../pgwire-recorder-austausch` auf dem Zweig `slice/devcontainer-austausch-repo`. `Verantwortlich:` blieb `pt9912` ohne Zweig, weil der Lauf auf `main` stand und die Form `<Person> (Zweig <branch-name>)` einen Zweig braucht, den es für diesen Lauf nicht gab; der Plan des anderen Slice trägt die Zweigform. Die Ablage unterschied die beiden Läufe über die Verzeichnis-Position und die Commits; ob der Hauptzweig als Zweig zu benennen ist, entscheidet der Nutzer.
- **Folge-Slices:** keiner.
- **Marke (`AGENTS.md` §3.14):** Summe aus Einfügungen und Löschungen, Metrik aus §4: **909** am Stand `74446efb` (Meldeschwelle 800 erreicht und gemeldet, Marke 1000 nicht).
- **Risiken aus §6:** Gesamtzahl weiter offen (Register); die übrigen vier entfallen mit Beleg; Begründungen in §6.
- **`make gates` am Endstand:** Exit 0 am Stand `63619f98`. Der erste Lauf des Implementers, mit `nohup` gestartet, war ohne Codeänderung rot (`container-gegenprobe`, `SIGINT` wird von einem so gestarteten Prozess geerbt ignoriert); der Lauf über `run_in_background` war grün. Als Beobachtung ins Register (`BEO-REPO/lauf-ignoriert-sigint-nach-nohup`, 1×, offen).
- **Drei Paarungen:** Anker (kein `liegt in`: keine Regel verkörpert; der neue Sensor liegt im Repo, `SPEC-048` und `harness/README.md` §Sensors tragen ihn mit `· seit slice-harness-abdeckung-gate`) · Folge-Slice (keiner zugewiesen; Teil 2 von LH-QA-07 gehört `slice-harness-coverage`, dessen §1 und DoD die Kennung dieses Slice nennen, `grep` findet sie dort) · Register (fünf Belege ergänzt, ein Kein-Auftreten, zwei Einträge neu).

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

**Nachgezählt nach der Vorprüfung des Architect (2026-10-17):** Eine ADR und die RTM-Umstellung
entfallen (§1), `harness/sensors/abdeckung-check.md` kommt dazu; Liefer-Punkte weiter 3, Schichten
weiter 2: `harness/README.md` und `harness/sensors/abdeckung-check.md` zählen mit `docs/user/` zur
Doku-Schicht, nicht zu einer dritten (`AGENTS.md` §3.13, die Wartungs-Doku ist eine Schicht).

**Nachgezählt** (Stand `91a1bd02`, `AGENTS.md` §3.13) mit der Sendung aus
`slice-harness-integration-volumes` (Kopf der Gegenprobe ohne Deklaration, Fall in der
Gegenprobe) und dem Schätzwert aus §4: Liefer-Punkte 3 (Werkzeug und Vertrag · Gegenprobe · erste
Gate-Deklarationen und RTM; die Sendung wächst in Punkt 2 hinein, ein Fall mehr). Schichten 2:
Harness-Werkzeug (`tools/test`, `tools/harness/lint-gegenprobe.sh`, `.d-check.yml`,
`harness/README.md`, `harness/vertraege/`) und Nutzer-Doku (`docs/user/abdeckung-*.md`, erzeugte
Tabellen); die ADR ist Planung. Kein `gates-schnell.mk`-Eintrag: das Ziel `abdeckung-gegenprobe`
steht schon in der Ausnahmeliste (`SPEC-052`), und es kommt kein neues Gate hinzu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
