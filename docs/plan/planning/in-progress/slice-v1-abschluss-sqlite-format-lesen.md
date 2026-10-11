# Slice slice-v1-abschluss-sqlite-format-lesen: SQLite-Aufzeichnungen lesen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), [ADR-0013](../../adr/0013-sqlite-recording-backend.md)

**Berührte Spec-Stellen:** `LH-FA-22.a` · `SPEC-001` · `SPEC-043` · `SPEC-044` · `ARC-008` · `ARC-014`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-11.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `replay` und `play` erkennen das Format der Eingabedatei an ihren ersten 16 Bytes und lesen eine SQLite-Aufzeichnung wie dieselbe Aufzeichnung als YAML (Roundtrip-Gleichheit, Abnahmeszenario 14); eine Datei mit unbekannter Version oder ohne gültige Aufzeichnung wird in beiden Formaten erkannt.

**Übernimmt:** `slice-v1-abschluss-sqlite-format` — dessen Lese-Hälfte nach dem Schnitt vom 2026-10-11 (Größe über einer Review-Sitzung, geschätzt 2500 bis 3000 Zeilen; Entscheidung des Nutzers vom 2026-10-11: Schnitt *Schreiben | Lesen*; dort §1, *Abgegeben*, und §4). Im Einzelnen:

- die Formaterkennung an den ersten 16 Bytes für `replay` und `play` und das Lesen der Datei in einer Lesetransaktion (`LH-FA-22.a` *Wahl und Erkennung*, *Lesen von `sqlite`*), mit `mode=ro` und der Sperrfrist von fünf Sekunden;
- die Umkehr des `base64:`-Escapes beim Lesen (`SPEC-043`, Regel `base64:`) und das Lesen von `{}` in `startup` und `server_parameters` wie das fehlende Feld in `yaml`;
- die Fehlerklassen des Lesens: `PGR-E3002` für jede kanonische Version außer 1, `PGR-E3003` für alles, was die Datei als Aufzeichnung ungültig macht (Typ-Prüfung der Spaltenwerte, verwaiste Zeilen, doppelter Schlüssel, `fields` fehlt oder ist kein Objekt, Zeile gegen die Regeln des Abschnitts), `PGR-E3001` für die Menge der primären Ergebniscodes aus `LH-FA-22.a` *Lesen von `sqlite`*, abgeschnittene Datei `PGR-E3003`;
- das Lesen von `meta.version` (nur kanonischer Dezimaltext) und `meta.empty_sessions` (nur `true`);
- das Ersetzen des Platzhalters `Load` des SQLite-Adapters (aus `slice-v1-abschluss-sqlite-format`, dort §3: liefert bis dahin `PGR-E1000`) durch den Leser, mit der Bildung der `file:`-URI aus dem Schreiben (kein zweiter Weg);
- der Roundtrip über das Binary: `record --format sqlite` (aus `slice-v1-abschluss-sqlite-format`), dann `replay` und `play`, gleich dem Verhalten bei YAML, und die gleiche Ablehnung beider Formate bei unbekanntem Schlüssel, `null` und doppeltem Schlüssel.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Schreiben von `sqlite`, die Fabrik `recording.New(format)`, die Option `--format` und alles Schema — `slice-v1-abschluss-sqlite-format` (dort §1 und DoD), das vor diesem Slice schließt (§4): Dieser Slice liest, was jener schreibt.
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-sqlite-format-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten.
- Die Zeitangaben in der SQLite-Datei (`interaction.offset_ms` ist hier nur Spalte, deren Typ geprüft wird) — `slice-v1-abschluss-zeitangaben` (dort DoD): Das Verhalten mit Zeitangaben ist ein anderer Liefer-Punkt.
- Umwandlung zwischen den Formaten und weitere Formate — Out-of-Scope von LH-FA-22.
- Das Standardformat YAML — bleibt unverändert aus `slice-walking-skeleton-record`; dieser Slice ändert sein Lesen nicht, nur die Wahl des Lesers.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat): Eine Aufzeichnung mit `--format sqlite` liefert im Replay und beim Einspielen dasselbe Verhalten wie dieselbe Aufzeichnung als YAML (Roundtrip-Gleichheit, Abnahmeszenario 14); `replay` und `play` erkennen das Format der Datei selbst (Test über das Binary). Die Datei wird nur lesend geöffnet und nie verändert (Test).
- [ ] [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats): Eine Datei mit unbekannter Version (`PGR-E3002`) oder ohne gültige Aufzeichnung (`PGR-E3003`) wird in beiden Formaten erkannt, dieselbe Aufzeichnung aus `yaml` und aus `sqlite` bekommt dieselbe Ablehnung (Test).
- [ ] Die Einstufung der Fehler des Lesens folgt `LH-FA-22.a` *Lesen von `sqlite`*: `PGR-E3001` für die Menge der primären Ergebniscodes, sonst `PGR-E3003`, eine abgeschnittene Datei `PGR-E3003` (Test an der Seitengrenze und mitten in einer Seite).
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
| `internal/adapters/driven/recording` | update | Leser für `sqlite` hinter `RecordingRepository`: Formaterkennung an den ersten 16 Bytes, Lesetransaktion, Prüfung der Spalten und der Regeln aus `SPEC-043`, Umkehr des `base64:`-Escapes; Fehlerklassen nach dem primären Ergebniscode |
| `internal/bootstrap` | update | `replay` und `play` reichen die Wahl des Lesers durch statt `recording.YAML{}`; reine Verdrahtung, der Bootstrap kennt kein Format und keinen Adapter (wie die Fabrik in `slice-v1-abschluss-sqlite-format`) |
| `test/integration` | update | Roundtrip-Gleichheit beider Formate über das Binary (`record --format sqlite`, dann `replay` und `play`); ohne `sqlite3`-Binary und ohne Änderung an `tools/test/run-integration-tests.sh` |

**Aufträge an den Implementer** (vom Architect vor dem Code zu setzen und hier zu prüfen; Quelle der Regeln: `LH-FA-22.a`, `SPEC-043`; die Aufträge zum Lesen standen vor dem Schnitt in `slice-v1-abschluss-sqlite-format` §3):

- Der Codec bleibt ein Pfad (Rückführung §4): JSON-Baum ↔ `yaml.Node` ↔ die DTOs von `yaml`. Die Umwandlung kennt nur JSON-Typen (Zahl → `!!int`, String → `!!str` in Anführungszeichen, Wahrheitswert, `null`), keine Nachrichtenart; braucht sie ein Wissen über ein Feld, ist das die Rückführung. `Node.Decode` kennt kein `KnownFields`: der unbekannte Schlüssel und die Vorprüfung der Form gehen denselben Weg wie bei `yaml`. Ein Test belegt, dass dieselbe Aufzeichnung aus `yaml` und aus `sqlite` dieselbe Ablehnung bekommt (unbekannter Schlüssel, `null`, doppelter Schlüssel).
- Strings ohne gültiges UTF-8 beim Lesen (`SPEC-043`, Regel `base64:` mit Umschreibung): Tests mindestens für das Byte `0xE9` in einem Startup-Wert, einem Fehlerfeld und einem Schlüssel, für einen gültigen String, der mit `\` oder `base64:` beginnt, und für ungültiges Base64 (`PGR-E3003`); `interaction.sql` mit `0xE9` unverändert hin und zurück. Liefert die Bibliothek die Bytes von `interaction.sql` nicht unverändert zurück, gibt der Implementer es dem Architect zurück.
- Datei als `file:`-URI mit `mode=ro` zum Lesen; Fehlerklasse über `Error.Code()&0xff` und die Menge in `LH-FA-22.a` *Lesen von `sqlite`*; Test mit einer abgeschnittenen Datei an der Seitengrenze und mitten in einer Seite, jeweils `PGR-E3003`. Der Pfad geht als korrekt escapte `file:`-URI an die Bibliothek (Test mit `?`, `#`, `%` und Leerzeichen im Pfad, wie beim Schreiben in `slice-v1-abschluss-sqlite-format`).
- Je Zusage dieses Vertrags ein Test und eine selbst gefahrene Mutation, im Bericht Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
- Die Doc-Kommentare nachziehen, die der Slice ändert oder die ihm widersprechen (`AGENTS.md` §3.7, §3.11).
- Erkennung vor dem Lesen (`LH-FA-22.a` *Wahl und Erkennung*): Datei öffnen, bis zu 16 Bytes lesen; Fehler dabei `PGR-E3001` mit der Ursache, weniger als 16 Bytes ist `yaml`; der Eingabepfad gilt wörtlich. Tests: fehlende Datei, Verzeichnis, leere und 15-Byte-Datei, Pfad mit Sonderzeichen.
- Eine Randform, die §6 nicht nennt, entscheidet er nicht, sondern gibt sie dem Architect zurück (`AGENTS.md` §3.12).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-sqlite-format` liegt in `done/` (Voraussetzung: das Schreiben, der Schema-Text und die Fabrik; Schritt 17 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1), `slice-v1-abschluss-einspielen` ist `done` (Kommando `play`). Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar; die Schätzung liegt bei 900 bis 1000 Zeilen. Schnitt dann: DoD-Punkt 3 (Einstufung der Fehler des Lesens) als eigener Slice.
- `in-progress` → `next`: Der Slice braucht mehr als einen Codec-Pfad (das Lesen von `sqlite` trägt das Modell nur mit einer zweiten, eigenen Abbildung der Nachrichten neben `yaml`) — zurück zur Zerlegung.
- `in-progress` → `open`: Die Bibliothek liefert die Bytes eines Strings oder die Ergebniscodes nicht so, wie `LH-FA-22.a` und `SPEC-043` sie voraussetzen — Entscheidung klären.

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

- **Risiko: Zwischenstand** (aus `slice-v1-abschluss-sqlite-format`, dort §6): Bis zum Schluss dieses Slice liest `replay` eine Datei von `record --format sqlite` als YAML und meldet sie als beschädigt. Das ist wellenintern und ohne Release dazwischen akzeptiert (Entscheidung des Nutzers vom 2026-10-11) — **Ausgang: entfallen**, sobald dieser Slice schließt: Erkennung und Leser liegen dann vor (DoD-Punkt 1).
- Eine abgeschnittene Datei meldet die Bibliothek je nach Stelle der Kürzung (Seitengrenze, mitten in einer Seite) mit anderen Ergebniscodes; `LH-FA-22.a` legt `PGR-E3003` fest — **Ausgang:** offen bis Closure (Test an beiden Stellen, DoD-Punkt 3).
- Das Lesen ist nur unter Linux im Container geprüft (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×); ob Sperre und Lesetransaktion unter macOS, Windows und auf einem Netzwerk-Dateisystem halten, sagt `LH-FA-22.a` *Grenzen* nur für das Schreiben — **Ausgang:** offen bis Closure (Register).

**Randformen des Lesens** (`AGENTS.md` §3.12) — je Randform der Ort der Entscheidung; die Spezifikation steht vor dem ersten Code-Commit. Die Randformen des Lesens trug bis zum Schnitt vom 2026-10-11 `slice-v1-abschluss-sqlite-format` §6; die Entscheidungen des Architect vom 2026-10-11 (Commits `36cfbf0`, `cc3ad1a`, `bb67108`) gelten unverändert. Eine Randform, die nur das Schreiben trägt, steht dort.

| Randform | Entschieden in |
|---|---|
| Typ-Prüfung der Spaltenwerte beim Lesen (`PGR-E3003`) | `LH-FA-22.a` *Lesen von `sqlite`* |
| `fields` fehlt oder ist kein JSON-Objekt; kein `bytes` (`PGR-E3003`) | `SPEC-043` |
| Erkennung an den ersten 16 Bytes; leere oder kürzere Datei und jede andere Datei als YAML | `LH-FA-22.a` *Wahl und Erkennung* |
| Lesetransaktion: ein Lesen sieht nur bei Beginn vollständige Sessions; Sperre fünf Sekunden; nur lesend geöffnet, auch auf einem schreibgeschützten Mount | `LH-FA-22.a` *Lesen von `sqlite`* |
| `meta.version`: nur kanonischer Dezimaltext, sonst `PGR-E3003`; kanonische Zahl ≠ 1 `PGR-E3002`; `meta.empty_sessions` nur `true`; Reihenfolge Formatkennung, Version, Inhalt | `SPEC-043`, `LH-FA-22.a` *Lesen von `sqlite`* |
| Zusätzliche Tabellen, Spalten, Indizes und Trigger ignoriert; verwaiste Zeilen, doppelter Schlüssel, `NULL`, Zeile gegen die Regeln des Abschnitts, `{}` in `startup`/`server_parameters` | `LH-FA-22.a` *Lesen von `sqlite`*, `SPEC-043` |
| Fehlerklasse beim Lesen nach primärem SQLite-Ergebniscode; abgeschnittene Datei `PGR-E3003` | `LH-FA-22.a` *Lesen von `sqlite`* |
| Strings ohne gültiges UTF-8 beim Lesen (Umkehr von `base64:` und `\`; ungültiges Base64 `PGR-E3003`); `interaction.sql` roh | `SPEC-043` |
| Rollback-Journal neben der Datei beim nur lesenden Öffnen (nicht zugesagt) | `LH-FA-22.a` *Lesen von `sqlite`*, *Grenze* |
| Gleiche Ablehnung beider Formate bei unbekanntem Schlüssel, `null`, doppeltem Schlüssel | `LH-FA-22.a` *Lesen von `sqlite`* (Forderung an den Leser); Test hier |
| Leser und Erkennung im Recording-Adapter, der Bootstrap reicht nur durch (wie die Fabrik für das Schreiben) | `ARC-008` (Adapter-Grenze); hier §3 — der Architect prüft sie vor dem Code |
| Eingabepfad mit `?`, `#`, `%` und Leerzeichen gilt wörtlich, für die Erkennung wie für beide Leser | `LH-FA-22.a` *Wahl und Erkennung* |
| Fehler beim Öffnen oder Lesen der ersten 16 Bytes (Datei fehlt, nicht lesbar, Verzeichnis, Ein- oder Ausgabefehler) `PGR-E3001` wie bei `yaml`; weniger als 16 Bytes ist `yaml`; Reihenfolge: Erkennung, dann Leser | `LH-FA-22.a` *Wahl und Erkennung* |
| `Load` des SQLite-Adapters aus `slice-v1-abschluss-sqlite-format` (Platzhalter mit `PGR-E1000`) wird hier durch den Leser ersetzt | `slice-v1-abschluss-sqlite-format` §3 (Geber); hier §1 |

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) — der Slice liefert den Leser eines neuen Vertrags; §3 fordert je Zusage die Mutation, §6 nennt die Randformen.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist Nehmer des Lese-Teils von `slice-v1-abschluss-sqlite-format`; §1 nennt ihn mit der Kennung des Gebers, im selben Commit wie der Geber ihn abgibt.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) und `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` (2×, offen) — der Schnitt vom 2026-10-11 teilt auf rund 1800 und 900 bis 1000 Zeilen, keine Hälfte liegt an der Grenze (kein dritter Beleg); der Geber lief über `next/` zurück (kein dritter Beleg).
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — das Lesen ist nur unter Linux geprüft (§6).

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-sqlite-format` (`AGENTS.md` §3.13): Liefer-Punkte: 3 (Lesen und Roundtrip, Ablehnung beider Formate, Einstufung der Fehler). Schichten: eine (Recording-Adapter mit seinen Tests); der Bootstrap reicht den Leser nur durch (reine Verdrahtung, zählt nicht), der Integrationstest läuft über das Binary und trägt keine eigene Schicht; die Dokumentation liegt im Doku-Folge-Slice `slice-v1-abschluss-sqlite-format-doku`. Geschätzt 900 bis 1000 Zeilen mit Tests.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält die Aufzeichnung als Code unter Spezifikation und Entscheidungen).
