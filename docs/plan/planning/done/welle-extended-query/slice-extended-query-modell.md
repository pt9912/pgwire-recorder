# Slice slice-extended-query-modell: Spezifikation und Domain-Modell

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-extended-query.

**Bezug:** [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [ADR-0006](../../adr/0006-kanonisches-domain-model.md), [ADR-0007](../../adr/0007-strict-replay.md), [ADR-0012](../../adr/0012-extended-query-gruppen.md)

**Berührte Spec-Stellen:** `LH-FA-18.a` · `SPEC-041` · `SPEC-001` · `SPEC-002` · `SPEC-039` · `ARC-001` · `ARC-008`

**Verantwortlich:** —
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Domain-Modell trägt die Typen des Extended Query Protocol (Client-Nachrichten, Server-Nachrichten, Gruppen) und die Formregeln einer Interaktion (`Interaction.Validate`); die Verfeinerung in der Spezifikation (`LH-FA-18.a`, `SPEC-041`) ist mit dem Modell abgeglichen, und die Entscheidung dazu ist als ADR angenommen. Der Recording-Adapter (Format `yaml`) schreibt und liest `type: extended` nach `SPEC-041`.

Der Leser gehört in diesen Slice, weil `SPEC-001` Extended-Interaktionen schon zur Formatversion 1 zählt: ohne ihn lehnt jeder Lauf eine gültige Aufzeichnung als beschädigt ab, und das Modell hätte keinen Pfad, auf dem seine Formregeln greifen. Der Schreiber gehört dazu, weil `toDTO` eine Extended-Interaktion sonst still ohne ihre Gruppen schreibt. Damit liefert dieser Slice den Format-Punkt, den `slice-extended-query-record` bisher in seiner DoD trug.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Umsetzung in den PGWire-Adaptern und im Record-Service (Gruppieren nach `LH-FA-18.a`) — `slice-extended-query-record`.
- Replay und Matcher, einschließlich des Feldvergleichs einer Client-Nachricht — `slice-extended-query-replay`. Der Replay-Service bleibt im Code unverändert: eine Extended-Interaktion am Cursor ist für eine einfache Anfrage eine Abweichung (`PGR-E5001`); ein Test hält diese Prüfung fest, weil der Leser die Interaktion jetzt liefert.
- SQLite-Format (`SPEC-043`) — `slice-v1-abschluss-sqlite-format`; es gibt noch keinen SQLite-Adapter.
- Eine Abdeckungs-Deklaration für `LH-FA-18`: kein Akzeptanzkriterium von `LH-FA-18` ist ohne Record und Replay belegbar; die Domain-Tests tragen keine. Die Lese-Negativfälle zählen zu `LH-FA-07`/Negative.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): `LH-FA-18.a` und `SPEC-041` stimmen mit dem Domain-Modell überein; die Entscheidung dazu steht in [ADR-0012](../../adr/0012-extended-query-gruppen.md).
- [x] Das Domain-Modell trägt Request- und Response-Typen des Extended Query Protocol und ihre Formregeln (Domain-Tests, je Regel ein Negativfall).
- [x] [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings): Der Recording-Adapter schreibt und liest `type: extended` nach `SPEC-041` (Roundtrip, Beispiel aus `SPEC-041`); eine abgeschnittene oder fremd geformte Extended-Interaktion ist beschädigt (`PGR-E3003`, Unit-Tests).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der Closure von `welle-extended-query` geprüft.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` (`LH-FA-18.a`, `SPEC-001`, `SPEC-041`) | update | Abgleich mit dem Domain-Modell: `ParameterStatus` fehlte unter den Server-Nachrichten (`LH-FA-24.a` setzt es voraus); Felder je Client-Nachricht und `param_types` der `parameter_description` festgelegt. Nach Review: Stelle von `type` je Art (`SPEC-001`, `SPEC-041`); fehlendes oder mit `null` belegtes Feld, `param_types` nur und stets an `parameter_description`, Form-Schlüssel mit `null` sind beschädigt. Nach Verifikation: leerer oder `null`-`type`, `param_types: null` an anderen Antworten, `client`/`server` je Gruppe Pflicht. Nach Folge-Review: Anker, Aliase und Merge-Keys sind beschädigt (`SPEC-001`). Bei Closure (F-300, Entscheidung des Nutzers): Format-Entwicklung streng — jedes neue Feld erhöht `version`, ein Leser lehnt unbekannte Felder ab und liest jede `version` von 1 bis zu seiner eigenen (`SPEC-001`); der Leser tat das für Version 1 schon (`KnownFields`, Schlüssel je Typ, `version` ungleich 1 ist `PGR-E3002`). Historienzeilen |
| `spec/architecture.md` (`ARC-001`) | update | Komponentenzeile nennt Extended-Gruppe, Client-Nachricht und die Formregeln |
| `internal/hexagon/model` | update | Neue Typen (`ClientMessage`, `Group`, Server-Nachrichten, `RequestExtended`), `Interaction.Validate` mit den Formregeln beider Arten; die Menge der Antworttypen einer einfachen Anfrage zieht aus dem Leser hierher |
| `internal/adapters/driven/recording` | update | DTOs für `type: extended`; Schreiber mit genau den Feldern je Client-Nachricht und mit `param_types` an jeder `parameter_description`; Leser verlangt die Schlüssel je Typ mit Wert, lehnt fremde ab, lehnt am YAML-Baum `request`/`responses`/`groups` mit `null`, leeren oder `null`-`type`, fehlende oder `null`-`client`/`server` einer Gruppe und `param_types` an anderen Antworten (auch `null`) ab und ruft `Validate`; lehnt als ersten Schritt Anker, Aliase und Merge-Keys ab; Meldungen der Vorprüfung nennen Session und Interaktion |
| `internal/hexagon/services/replay_test.go` | update | Test für die Typprüfung des Matchers: Extended-Interaktion am Cursor gegen einfache Anfrage, auch die leere |
| `docs/user/abdeckung-*.md` | update | `make abdeckung` nach der neuen Deklaration (`LH-FA-07`/Negative) |
| `docs/plan/planning/done/slice-extended-query-record.md` | update | Format-Punkt ist hier geliefert (`AGENTS.md` §3.9); DoD-Punkt ersetzt, Risiken zur Feldbelegung und zum Schreiben ohne `Validate` übergeben |
| `docs/plan/planning/done/slice-extended-query-replay.md` | update | Risiko nil gegenüber leerer Liste übergeben (`AGENTS.md` §3.9) |
| `docs/plan/planning/open/slice-v1-abschluss-zeitangaben.md` | update | Rückführungs-Trigger folgt F-300: schon ein neues optionales Feld erhöht `version` (`AGENTS.md` §3.9) |
| `docs/reviews/2026-10-04-review-slice-extended-query-modell.md`, `docs/reviews/2026-10-04-verifikation-slice-extended-query-modell.md`, `docs/reviews/2026-10-04-folge-review-slice-extended-query-modell.md`, `docs/reviews/2026-10-04-mutationen-slice-extended-query-modell.md` | neu | Review-Report, Verifikationsbericht, Folge-Review; Mutationstabelle für den Verifier |
| `docs/plan/planning/observations/BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/evidence/slice-extended-query-modell.md`, `docs/plan/planning/observations/BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/evidence/slice-extended-query-modell.md` | neu | Beleg je Klasse aus dem Review |
| `docs/plan/planning/observations/BEO-REPO/spec-randform-erst-im-review-entschieden/` | neu | Beobachtung aus Review, Verifikation und Folge-Review (Closure) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-walking-skeleton` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: die Spezifikation verlangt eine Änderung am Lastenheft — Change Request.
- `in-progress` → `open`: Eine Bibliotheksgrenze verhindert die Abbildung — Carveout.


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

- Das Domain-Modell zeigt eine Lücke in `SPEC-041` (Format Version 1 umfasst Extended, `SPEC-001`) — Stand: zwei Lücken gefunden und in der Spezifikation geschlossen: `ParameterStatus` fehlte unter den Server-Nachrichten von `LH-FA-18.a`/`SPEC-041`, und `SPEC-041` legte die Felder je Client-Nachricht und die der `parameter_description` nicht fest. Nach Review zwei weitere: der Wortlaut zur Stelle von `type` je Art und die Behandlung fehlender Felder (F-285, F-286); nach Verifikation drei Randformen: `param_types: null`, `client`/`server` je Gruppe, leerer oder `null`-`type` (V-16 bis V-18); nach Folge-Review die Lücke, dass die Vorprüfung am YAML-Baum Anker, Aliase und Merge-Keys nicht auflöst (F-294, F-295): geschlossen, indem die Spezifikation sie aus dem Format nimmt und der Leser sie als beschädigt ablehnt; bei Closure die Spannung zwischen `version` und unbekannten Feldern (F-300): jedes neue Feld erhöht `version`. **Ausgang:** eingetreten: jede Lücke in `SPEC-001`/`SPEC-041` geschlossen (Historienzeilen 2026-10-04 und 2026-10-05), der Leser folgt jeder mit Negativfall.

- Der Leser des Recording-Adapters lehnt `type: extended` noch als beschädigt ab; dieser Slice muss Extended-Interaktionen lesen — Stand: in §1 aufgenommen, Leser und Schreiber liefern `type: extended`. **Ausgang:** eingetreten: Leser und Schreiber liefern `type: extended` (`TestExtendedRoundtrip`, `TestUnmarshalSpec041`, `TestUnmarshalExtendedFehler`).

- Leere Liste und nil sind im Modell gleichbedeutend; der Leser liefert nil, ein PGWire-Adapter kann leere Listen liefern. Ein Matcher, der Client-Nachrichten mit `reflect.DeepEqual` vergleicht, meldet dann eine Abweichung, die keine ist — betrifft `slice-extended-query-replay`. **Ausgang:** weiter offen → `slice-extended-query-replay` §6.

- `Validate` prüft je Client-Nachricht Typ, Stellung und Zielart, nicht, ob nur die Felder ihres Typs belegt sind; das prüft allein der YAML-Leser über die Schlüssel (fremde abgelehnt, fehlende und `null` ebenso). Der Schreiber gibt nur die Felder des Typs aus, ein fremd belegtes Feld einer Modell-Nachricht fiele beim Schreiben still weg; ebenso `Response.ParamTypes` an jeder Antwort außer `parameter_description` (Verifikation V-21) — betrifft das Mapping in `slice-extended-query-record`. **Ausgang:** weiter offen → `slice-extended-query-record` §6.

- Der Schreiber prüft nicht mit `Validate`: eine Gruppe mit zwei `sync` wird geschrieben und erst beim Laden als beschädigt erkannt (Review F-292) — betrifft das Gruppieren in `slice-extended-query-record`, dessen DoD `Validate` je aufgezeichneter Interaktion verlangt. **Ausgang:** weiter offen → `slice-extended-query-record` §6 und DoD-Punkt 2.

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

- **Was hat funktioniert:** Das Domain-Modell trägt Client-Nachrichten, Gruppen und Server-Nachrichten des Extended Query Protocol mit `Interaction.Validate`; der Recording-Adapter schreibt und liest `type: extended` nach `SPEC-041`, das Beispiel der Spezifikation wird wörtlich gelesen, und die Ausgabe einfacher Anfragen blieb byteidentisch. Jede Formregel und jede Leser-Regel hat einen Negativfall, dessen rot färbende Mutation gesehen ist (Mutationstabelle, von Verifikation und Folge-Review in Kopien nachgefahren).
- **Was ging anders als geplant:** Leser und Schreiber kamen aus `slice-extended-query-record` in diesen Slice (§1). Die Spezifikation trug die Randformen des Formats nicht; vier Runden (Review, Verifikation, Folge-Review, Closure) entschieden je weitere — fehlende und `null`-Felder, Stelle und Wert von `type`, `client`/`server` je Gruppe, Anker, Aliase und Merge-Keys, zuletzt die strenge Format-Entwicklung (F-300, Entscheidung des Nutzers). Die Vorprüfung am YAML-Baum, mit der `d2efe37` V-16 bis V-18 schloss, brachte selbst eine Regression (F-294, F-295), weil sie Aliase nicht auflöste; `5c6a53d` nahm Aliase aus dem Format, statt sie aufzulösen.
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke, geschlossen: Die Randformen des Recording-Formats (`SPEC-001`, `SPEC-041`) sind entschieden, Beleg die Historienzeilen 2026-10-04 und 2026-10-05 der Spezifikation. Kein `liegt in`: Die Spezifikation trägt keinen Herkunfts-Anker auf einen Slice, und ein Sensor oder eine Regel des Harness entstand nicht — die Vorprüfung am YAML-Baum und `ohneVerweise` sind Produkt-Code mit Unit-Tests, die Mutationstabelle ist ein einmaliger Beleg, kein laufender Sensor. Die Frage aus Review und Folge-Review an den Architect (Mutationslauf als Gate, Regel „Prüfung je Zusage“) entscheidet die Closure von `welle-extended-query` mit den beiden 3×-Einträgen.
- **Beobachtungs-Register (`../observations/`):** `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` je Beleg `evidence/slice-extended-query-modell.md` — Zähler je 3×, beide gehen in die Closure von `welle-extended-query`, die über die Verkörperung entscheidet; `state.md` bleibt bis dahin `offen`. Neu: `BEO-REPO/spec-randform-erst-im-review-entschieden/` — 1×.
- **Folge-Slices:** `slice-extended-query-record` (Feldbelegung im Mapping, `Validate` vor dem Übernehmen; F-292, V-21), `slice-extended-query-replay` (nil gegenüber leerer Liste) — beide Dateien in `open/`. Die Pflicht aus `SPEC-001`, dass ein Leser jede `version` von 1 bis zu seiner eigenen liest, trifft den ersten Slice, der `version` erhöht; für Version 1 hält der Leser sie (`version` ungleich 1 ist `PGR-E3002`).
- **Risiken aus §6:** jedes mit genau einem Ausgang — zwei eingetreten, drei weiter offen mit Folge-Slice; siehe §6.
- **Drei Paarungen:** Repo mit Wellen-Betrieb — geprüft von der Closure von `welle-extended-query`.
- **Belege:** Review `docs/reviews/2026-10-04-review-slice-extended-query-modell.md` (`cd03e09`), Verifikation `docs/reviews/2026-10-04-verifikation-slice-extended-query-modell.md` (`6b4967d`; DoD 1, 2 und `make gates` bestätigt, DoD 3 im Kern mit Rest V-16), Folge-Review `docs/reviews/2026-10-04-folge-review-slice-extended-query-modell.md` (`d2efe37`; V-16 für ausgeschriebene Werte behoben, Rest über Merge-Keys F-295), Mutationstabelle `docs/reviews/2026-10-04-mutationen-slice-extended-query-modell.md`. Die Korrektur `5c6a53d` (F-294 bis F-299, schließt den Rest von V-16) belegen `TestUnmarshalVerweise`, `TestVorpruefungNenntOrt`, die Mutationen P1 bis P4 der Implementer-Rolle und `make gates`; ein weiteres Review gibt es nicht. Die Spezifikations-Änderung zu F-300 belegt `make gates`; der Code war unverändert.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `../observations/BEO-REPO/` durchgegangen; einschlägig sind:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (offen; mit dem Beleg dieses Slice 3×): der Slice führt einen neuen Lese-Vertrag ein. Jede Formregel von `Validate` und jede Leser-Regel hat einen Negativfall, dessen Meldungstext die gemeinte Regel nennt; die rot färbende Mutation je Regel ist einmal gesehen worden; die Tabelle steht in `docs/reviews/2026-10-04-mutationen-slice-extended-query-modell.md`. Das Review fand dennoch zwei ungehaltene Zusagen (F-283, F-284); der Beleg dazu liegt im Register.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (offen; mit dem Beleg dieses Slice 3×): die Kommentare von `Validate`, `fromDTO` und `clientFelder` nennen nur, was geprüft wird; die Grenze von `Validate` (keine Feldprüfung je Typ) steht in §6. Review, Verifikation und Folge-Review fanden Zusagen weiter als die Prüfung (F-289, `request: null`; V-19, Testkommentar und Abdeckungszeile; F-296, Kommentare zur Vorprüfung über Aliase); behoben, Beleg im Register.
- `BEO-REPO/plan-folgt-korrektur-nicht` (verkörpert in `AGENTS.md` §3.9): §1, §3, §6 und der Kopf folgen dem Umschnitt; `slice-extended-query-record` nennt, was hier geliefert ist.

Nicht einschlägig: `gate-konfiguration-wirkt-anders-als-gelesen`, `gate-regel-ersetzt-statt-ergaenzt`, `kern-fremdimporte-nur-ueber-tech-liste` (keine Gate- oder Architektur-Konfiguration berührt), `record-fehlerantwort-im-aufbau-ungeregelt` (Record unberührt), `werkzeug-commit-ohne-zugelassene-kennung` (behoben).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (deklariert in `harness/conventions.md`; der Slice ergänzt Typen und Format ohne Bestand, der konvergieren müsste).
