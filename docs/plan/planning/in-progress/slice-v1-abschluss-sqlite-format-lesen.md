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
- die Fehlerklassen des Lesens: `PGR-E3002` für jede kanonische Version außer 1, `PGR-E3003` für alles, was die Datei als Aufzeichnung ungültig macht (Typ-Prüfung der Spaltenwerte, verwaiste Zeilen, doppelter Schlüssel, `fields` fehlt oder ist kein Objekt, Zeile gegen die Regeln des Abschnitts), `PGR-E3001` für die Menge der primären Ergebniscodes aus `LH-FA-22.a` *Lesen von `sqlite`*, abgeschnittene Datei `PGR-E3003`; dazu aus der Rückgabe des Implementers (Entscheidung des Architect vom 2026-10-11): ein maskiertes einzelnes Surrogat im JSON-Text, ein doppelter Schlüssel in `meta` und ein Verstoß gegen die Zählung und Reihenfolge der Gruppen `PGR-E3003`, der Abbruch des Laufs `PGR-E3001`, und eine Meldung zu einem Fehler des Inhalts ohne Wert und ohne Zeilennummer;
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
| `internal/adapters/driven/recording` | update | Leser für `sqlite` hinter `RecordingRepository`: `SQLite.Load` (`sqlite_lesen.go`: Lesetransaktion, Prüfung der Spalten und der Regeln aus `SPEC-043`, Fehlerklassen nach dem primären Ergebniscode; `sqlite_lesen_baum.go`: Interaktionen und Gruppen; `sqlite_json_lesen.go`: JSON-Baum und Umkehr des `base64:`-Escapes) und der Typ `Leser` (`leser.go`: Formaterkennung an den ersten 16 Bytes); `uriVon` in `sqlite.go` bildet die `file:`-URI für Schreiben und Lesen |
| `internal/bootstrap` | update | `replay` und `play` setzen statt `recording.YAML{}` `recording.Leser{}` ein; reine Verdrahtung, der Bootstrap kennt kein Format und unterscheidet keinen Fehler (Entscheidung vom 2026-10-11) |
| `test/integration` | update | Roundtrip-Gleichheit beider Formate über das Binary (`sqlite_lesen_e2e_test.go`: `record --format sqlite`, dann `replay` und `play`; die Aufzeichner `aufnehmenExtended` und `nimmAufEinspielAblaufAls` wählen das Format nach der Endung `.sqlite`); ohne `sqlite3`-Binary und ohne Änderung an `tools/test/run-integration-tests.sh` |

**Aufträge an den Implementer** (vom Architect vor dem Code zu setzen und hier zu prüfen; Quelle der Regeln: `LH-FA-22.a`, `SPEC-043`; die Aufträge zum Lesen standen vor dem Schnitt in `slice-v1-abschluss-sqlite-format` §3):

- Der Codec bleibt ein Pfad (Rückführung §4): JSON-Baum ↔ `yaml.Node` ↔ die DTOs von `yaml`. Die Umwandlung kennt nur JSON-Typen (Zahl → `!!int`, String → `!!str` in Anführungszeichen, Wahrheitswert, `null`), keine Nachrichtenart; braucht sie ein Wissen über ein Feld, ist das die Rückführung. `Node.Decode` kennt kein `KnownFields`: der unbekannte Schlüssel und die Vorprüfung der Form gehen denselben Weg wie bei `yaml`. Ein Test belegt, dass dieselbe Aufzeichnung aus `yaml` und aus `sqlite` dieselbe Ablehnung bekommt (unbekannter Schlüssel, `null`, doppelter Schlüssel).
- Strings ohne gültiges UTF-8 beim Lesen (`SPEC-043`, Regel `base64:` mit Umschreibung): Tests mindestens für das Byte `0xE9` in einem Startup-Wert, einem Fehlerfeld und einem Schlüssel, für einen gültigen String, der mit `\` oder `base64:` beginnt, und für ungültiges Base64 (`PGR-E3003`); `interaction.sql` mit `0xE9` unverändert hin und zurück. Liefert die Bibliothek die Bytes von `interaction.sql` nicht unverändert zurück, gibt der Implementer es dem Architect zurück.
- Datei als `file:`-URI mit `mode=ro` zum Lesen; Fehlerklasse über `Error.Code()&0xff` und die Menge in `LH-FA-22.a` *Lesen von `sqlite`*; Test mit einer abgeschnittenen Datei an der Seitengrenze und mitten in einer Seite, jeweils `PGR-E3003`. Der Pfad geht als korrekt escapte `file:`-URI an die Bibliothek (Test mit `?`, `#`, `%` und Leerzeichen im Pfad, wie beim Schreiben in `slice-v1-abschluss-sqlite-format`).
- Je Zusage dieses Vertrags ein Test und eine selbst gefahrene Mutation, im Bericht Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
- Die Doc-Kommentare nachziehen, die der Slice ändert oder die ihm widersprechen (`AGENTS.md` §3.7, §3.11).
- Erkennung vor dem Lesen (`LH-FA-22.a` *Wahl und Erkennung*): Datei öffnen, bis zu 16 Bytes lesen; Fehler dabei `PGR-E3001` mit der Ursache, weniger als 16 Bytes ist `yaml`; der Eingabepfad gilt wörtlich. Tests: fehlende Datei, Verzeichnis, leere und 15-Byte-Datei, Pfad mit Sonderzeichen.
- Der Leser mit Erkennung ist ein Typ des Adapters, der den Port `RecordingRepository` erfüllt, damit `NewReplayService` und `NewPlayService` unverändert bleiben (`Load` erkennt und ruft `YAML.Load` oder den Leser von `SQLite`; `Prepare` und `Write` liefern `PGR-E1000`, ein Aufruf ist ein Programmierfehler). `New(format)` bleibt die Fabrik des Schreibens und ändert sich nicht. Ein Test belegt, dass `Load` mit beiden Formaten dieselbe Aufzeichnung liefert.
- JSON-Zahl im Lesen: nur ganze, kanonische Zahlen im Bereich von 64 Bit mit Vorzeichen (`SPEC-043`); die Bibliothek darf sie nicht über `float64` führen (Test mit `9007199254740993`, `1.0`, `1e2`, `-0`, `9223372036854775808`).
- Zeilen ordnet der Leser nach den Nummern, nicht nach der Reihenfolge in der Datei (Test mit vertauscht eingefügten Zeilen); `PRAGMA user_version` und Fremdschlüssel liest und erzwingt er nicht.
- Maskiertes einzelnes Surrogat (`SPEC-043`): `jsonKnoten` lehnt es mit `PGR-E3003` ab, in Wert und Schlüssel, auch tief verschachtelt; die Bibliothek liefert dafür U+FFFD, also prüft der Leser den Text des Escapes selbst (kein Vergleich des Ergebnisses mit U+FFFD, denn das Zeichen selbst ist gültig). Tests in `sqlite_json_lesen`-Nähe (`TestLesenJSONVerletzt`): `"\ud800"`, `"\udc00"`, `"\ud800x"`, `"\ud800\u0041"`, `"\udc00\ud800"` abgelehnt; `"\ud83d\ude00"` (Paar) und ein wörtliches U+FFFD angenommen. Mutation: Prüfung entfällt; Paar wird abgelehnt; U+FFFD wird abgelehnt.
- Meldung zu einem Fehler des Inhalts (`LH-FA-22.a` *Lesen von `sqlite`*): ohne Wert aus der Datei und ohne Zeilennummer. Die Ursache der YAML-Bibliothek aus `Unmarshal` (Text mit `line N` und eingebettetem Wert) steht für `sqlite` nicht in der Meldung; der Leser nennt Datei und Art, bei eigenen Prüfungen zusätzlich Tabelle und Nummern. Der Implementer prüft dazu jede Meldung der Dateien `sqlite_lesen*.go` und `sqlite_json_lesen.go` auf `%q`/`%v`/`%s` eines Werts aus der Datei. Test: je ein Fall pro Weg (eigene Prüfung, `Unmarshal`-Weg, JSON) mit einem Wert wie `geheim-wert` in der Datei; die Meldung enthält ihn nicht und passt nicht auf `[Zz]eile \d+|line \d+`. Mutation: Ursache der Bibliothek wieder durchgereicht; Wert in eine eigene Meldung. Die Meldungen des Lesens von `yaml` bleiben unverändert.
- Schließen nach dem Lesen: der Fehler bleibt ignoriert, kein Test (kein herstellbarer Aufbau); der Doc-Kommentar sagt es im Indikativ, ohne Zusage.
- Abbruch von `ctx` beim Lesen: `PGR-E3001` mit der Ursache; Test mit einem abgebrochenen `ctx` (Mutation: `PGR-E3003`).
- Doppelter Schlüssel in `meta` `PGR-E3003` und Verstoß gegen `group_no` und die Reihenfolge der Gruppen `PGR-E3003`: sind beide schon Fälle von `TestLesenMeta` und `TestLesenInhaltVerletzt` (M14, B4 bis B8), kein neuer Code; der Test nennt die Stelle der Spezifikation.
- Schlüssel `<<` (F-630): ein Test, dass eine `sqlite`-Datei mit `{"<<":"x"}` in `startup`, `server_parameters` und in einem Objekt von `fields` `PGR-E3003` liefert (Datei per rohem SQL, `Write` schreibt ihn nicht). Mutation: `c.Value == "<<" ||` in `ohneVerweise` (`yaml.go`) entfällt; der Schlüssel steht in der Datei als String in Anführungszeichen, `yaml`-Tests mit `<<` ohne Anführungszeichen fangen das nicht. Kein Produktcode.
- Eine Randform, die §6 nicht nennt, entscheidet er nicht, sondern gibt sie dem Architect zurück (`AGENTS.md` §3.12).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-sqlite-format` liegt in `done/` (Voraussetzung: das Schreiben, der Schema-Text und die Fabrik; Schritt 17 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1), `slice-v1-abschluss-einspielen` ist `done` (Kommando `play`). Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar; die Schätzung des Architect vom 2026-10-11 liegt bei 1300 bis 1700 Zeilen mit Tests; Rückführung, wenn der Diff rund 2000 Zeilen übersteigt. Schnitt dann: DoD-Punkt 3 (Einstufung der Fehler des Lesens) als eigener Slice.
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
| Zahl im JSON-Text: Bruch, Exponent, außerhalb 64 Bit `PGR-E3003` | `SPEC-043` |
| Reihenfolge der Zeilen ohne Bedeutung; Fremdschlüssel und `user_version` beim Lesen nicht erzwungen und nicht gelesen | `LH-FA-22.a` *Lesen von `sqlite`*, `SPEC-043` |
| Write-Ahead-Log-Modus einer fremden Datei; Ersetzen der Datei zwischen Erkennung und Lesen (nicht zugesagt, akzeptiertes Negativ: harmlos und ohne Wiederkehr in diesem Werkzeug, das nur Journal-Modus schreibt) | `LH-FA-22.a` *Lesen von `sqlite`*, *Grenze* |
| Maskiertes einzelnes Surrogat im JSON-Text (`\ud800` allein, Paar angenommen, U+FFFD selbst gültig) `PGR-E3003` | `SPEC-043` (Regel `base64:`) |
| Meldung zu einem Fehler des Inhalts: ohne Wert aus der Datei, ohne Zeilennummer; Tabelle und Nummern, wo der Leser selbst feststellt | `LH-FA-22.a` *Lesen von `sqlite`* |
| Schließen der Datenbank nach dem Lesen scheitert: nicht gemeldet, nicht zugesagt (akzeptiertes Negativ: die Aufzeichnung ist ganz gelesen, ein Aufbau mit dem Fehler ist nicht herstellbar, ohne Wiederkehr) | `LH-FA-22.a` *Lesen von `sqlite`* |
| Abbruch von `ctx` und jeder andere Fehler ohne Typ der Bibliothek `PGR-E3001` | `LH-FA-22.a` *Lesen von `sqlite`* |
| Doppelter Schlüssel in `meta` `PGR-E3003` | `LH-FA-22.a` *Lesen von `sqlite`* |
| `group_no` ab 1 lückenlos in der Reihenfolge der Gruppen, Client-Nachricht nach Server-Nachricht derselben Gruppe `PGR-E3003` | `SPEC-043`, `LH-FA-22.a` *Lesen von `sqlite`* |
| Speicher des Lesers: Aufzeichnung als Ganzes, Baum und Text, ein Mehrfaches der Dateigröße, ohne Obergrenze (akzeptiertes Negativ, Entscheidung des Nutzers vom 2026-10-11: „große Aufzeichnungen“ in `LH-FA-22` gilt nur für das Schreiben; der Leser darf wie der von `yaml` die Aufzeichnung im Speicher halten, ein Mehrfaches der Dateigröße, ohne zugesagte Obergrenze; keine Änderung des Lastenhefts, kein streamendes Lesen) | `LH-FA-22.a` *Grenzen* (Schreiben von `sqlite`) |
| Schlüssel `<<` in `startup`, `server_parameters` oder `fields` (F-630): der Leser lehnt ihn wie `yaml` als Merge-Key ab (`PGR-E3003`), weil der Weg über `Unmarshal` ihn am Wert erkennt; das Schreiben (`SQLite.Write`, `PGR-E1000`) ist nicht zugesagt (akzeptiertes Negativ: Postgres kennt keinen Parameternamen `<<`, ohne Wiederkehr) | `LH-FA-22.a` *Lesen von `sqlite`* |
| Leser mit Erkennung erfüllt den Port, `Prepare` und `Write` `PGR-E1000`; der Bootstrap tauscht nur den Wert | hier §3 (Adapter-Interna, `ARC-008`) |
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

### Belege des Implementers

**Stand.** Geliefert sind die drei Liefer-Punkte der DoD (Lesen und Roundtrip, Ablehnung beider Formate, Einstufung
der Fehler des Lesens) in einem Code-Commit, `45af81d` seit `97f3ca6`, und dem Commit dieser Belege. Neu sind
`internal/adapters/driven/recording/leser.go`, `sqlite_lesen.go`, `sqlite_lesen_baum.go`, `sqlite_json_lesen.go`,
`leser_test.go`, `sqlite_lesen_test.go`, `sqlite_lesen_inhalt_test.go` und `test/integration/sqlite_lesen_e2e_test.go`;
geändert `sqlite.go` (`uriVon` für Schreiben und Lesen, der Platzhalter `Load` entfällt), `export_test.go`,
`sqlite_einstellungen_test.go` (der Test des Platzhalters entfällt), `internal/bootstrap/bootstrap.go` (`replay` und
`play` setzen `recording.Leser{}` ein), `extended_replay_e2e_test.go`, `play_extended_e2e_test.go`, §3 dieses Plans und
`docs/user/abdeckung-*.md`. `go.mod`, `go.sum`, Port, Handbuch und `README.md` sind unverändert; die Bibliothek
`modernc.org/sqlite/lib` gehört zum Modul des Treibers.
**Größe:** `git diff --numstat 97f3ca6..45af81d` zählt 2137 Zeilen hinzu und 37 weg; davon 799 Go-Code, 1291 Go-Tests
(1056 Unit-Tests in drei Dateien, 24 Zeilen Eingriffe in `export_test.go`, 211 am Binary) und 47 Doku und Tabellen.
Die Schätzung in §8 war 1300 bis 1700, die Schwelle der Rückführung in §4 rund 2000: **überschritten**, mit 2090 Zeilen
Go (Code und Tests). Die Entscheidung über Rückführung oder Schnitt (§4: DoD-Punkt 3 als eigener Slice) liegt beim
Planner; die Tests sind überwiegend Tabellen mit rohem SQL.

**Zusage · Mutation · roter Test.** Jede Mutation lief in einer frischen Kopie des Baums (`cp -r` ohne `-p`, `go.mod`,
`go.sum`, `cmd/`, `internal/`, `test/`), mit Ersetzungen in einer Datei, `gofmt -l` der Datei leer und
`go test -count=1 -run <Tests>` im gepinnten Go-Image mit Bind-Mount (der Mutant kommt ohne Build-Kontext an);
`K1` und `K2` am Bootstrap liefen mit `make test-integration` in einer Kopie des ganzen Baums. Die Kennung in
Klammern ist die des Laufs.

| Zusage | Mutation | roter Test |
|---|---|---|
| nur eine Datei, die mit `SQLite format 3` und NUL beginnt, ist `sqlite`; jede andere, auch eine leere oder kürzere, wird als `yaml` gelesen | nur 15 Bytes verglichen (E1); ohne das NUL, kürzere Datei als Präfix angenommen (E2); kurze Datei ist ein Fehler (E3) | `TestLeserErkennungAnDenErsten16Bytes` |
| die Erkennung folgt dem Inhalt, nicht der Endung | immer `yaml` (E6); immer `sqlite` (E7) | `TestLeserErkennungOhneEndung`, `TestLeserErkennungAnDenErsten16Bytes`, `TestLesenGleicheAblehnungBeiderFormate` |
| Öffnen der Datei scheitert: `PGR-E3001` mit Pfad | Code `PGR-E3003` (E4); beim Lesen der 16 Bytes (E5) | `TestLeserNichtLesbar` |
| `Prepare` und `Write` des Lesers sind `PGR-E1000` ohne Pfad | liefern `nil` (E8, E9); die Meldung nennt den Pfad (E10) | `TestLeserPrepareUndWrite` |
| der Eingabepfad gilt wörtlich (`?`, `#`, `%`, Leerzeichen) | `file:`-URI ohne Maskierung (O4) | `TestLesenSonderzeichenImPfad` |
| die Datei wird nur lesend geöffnet und nie verändert | `mode=rw` (O1) | `TestLesenNurLesend`; am Binary `TestE2EReplaySQLite`, `TestE2EPlaySQLite` (Bytes und Zeit unverändert) |
| Sperrfrist fünf Sekunden | ohne Frist (O2); 20 s (O3) | `TestLesenSperrfrist` |
| eine einzige Lesetransaktion | Anfragen außerhalb der Transaktion (O6); `immutable=1` (O7) | `TestLesenEineLesetransaktion` (O7 auch `TestLesenSperrfrist`) |
| der Meldungscode folgt dem primären Ergebniscode: zehn Codes `PGR-E3001`, jeder andere `PGR-E3003` | je einer der zehn aus der Menge genommen (S PERM bis S PROTOCOL); ohne Maske auf den primären Code (S12); Rückfall `PGR-E3001` (S13) | `TestLesenStufeNachErgebniscode`; S13 auch `TestLesenFehlendeTabelleOderSpalte`, `TestLesenAbgeschnitten` |
| ein abgeschnittenes Lesen (`IOERR_SHORT_READ`) ist `PGR-E3003` | Sonderfall entfernt (S11) | `TestLesenStufeNachErgebniscode` |
| Fehler ohne Typ der Bibliothek (Abbruch von `ctx`) `PGR-E3001`; Fehler der Bibliothek nach Stufe | `PGR-E3003` (S14); jeder Fehler der Bibliothek `PGR-E3001` (S15) | `TestLesenNichtLesbarUndKopfOhneInhalt`, `TestLesenFehlendeTabelleOderSpalte`, `TestLesenAbgeschnitten` |
| abgeschnittene Datei an der Seitengrenze und mitten in einer Seite: `PGR-E3003` | Rückfall `PGR-E3001` (S13); Fehler ohne Stufung (S15) | `TestLesenAbgeschnitten` (dreiundzwanzig Längen); am Binary `TestE2ELesenAbgelehnt` |
| Reihenfolge: Formatkennung, Version, Inhalt | Version vor Formatkennung (M9); unbekannter Schlüssel vor Version (M10) | `TestLesenMeta` |
| `meta`: `format` verlangt und gleich `pgwire-recorder`; `version` verlangt | Formatkennung ungeprüft (M1); fehlendes `format` angenommen (M2); fehlende `version` angenommen (M3) | `TestLesenMeta` |
| `version`: nur kanonischer Dezimaltext (`PGR-E3003`), jede kanonische Zahl außer 1 `PGR-E3002` | Text ungeprüft (M4, M11); Code `PGR-E3003` für Version ≠ 1 (M5) | `TestLesenMeta` (dreizehn Texte, sieben Zahlen) |
| `meta`: unbekannter Schlüssel, doppelter Schlüssel, kein TEXT: `PGR-E3003`; `empty_sessions` nur `true` | Schlüssel ignoriert (M6); doppelter ignoriert (M14); Wert ungeprüft (M7); `empty_sessions` ignoriert (M8); BLOB als Text (M12) | `TestLesenMeta`; M8 auch `TestLesenInhaltVerletzt`, `TestLesenRoundtrip` |
| die Art der gespeicherten Werte ist die des Spaltentyps | Zahl nimmt `float64` (M13); `offset_ms` ungeprüft (I1); `sql` bei `query` ungeprüft (I2); `sql` bei `extended` ungeprüft (I3); unbekannter `type` (I4); `direction` ungeprüft (I8) | `TestLesenInhaltVerletzt` |
| keine verwaiste Zeile | Interaktion ohne Session ignoriert (I5); Nachricht ohne Interaktion ignoriert (I6) | `TestLesenInhaltVerletzt` |
| `position` zählt ab 1 lückenlos | Position ungeprüft (I7) | `TestLesenInhaltVerletzt` |
| `group_no`: ab 1 lückenlos in der Reihenfolge, nur bei `extended`; Client vor Server in der Gruppe; `query` nur Server-Nachrichten | Richtung bei `query` ungeprüft (B1); `group_no` bei `query` ungeprüft (B2); `NULL` bei `extended` angenommen (B3); Beginn bei beliebiger Zahl (B4); Lücke (B5); zurück (B6); Client nach Server (B7); Gruppe 0 (B8) | `TestLesenInhaltVerletzt` |
| die Reihenfolge der Zeilen hat keine Bedeutung | Sessions (I15), Interaktionen (I16) oder Nachrichten (I17) nicht nach Nummern geordnet | `TestLesenReihenfolgeDerZeilen` |
| zusätzliche Tabellen, Spalten, Indizes und Trigger ignoriert | `SELECT *` für `message` (I18) oder `session` (I19) | `TestLesenIgnoriertZusaetze` |
| `{}` in `startup` und `server_parameters` ist das fehlende Feld; beide sind Objekte mit String-Werten | `{}` bleibt leere Abbildung (I14); Nicht-Objekt angenommen (I12); Wert kein String angenommen (I13); `fields` kein Objekt angenommen (I11) | `TestLesenLeereObjekteSindKeineEintraege`, `TestLesenRoundtrip`, `TestLesenJSONVerletzt` |
| derselbe Inhalt wie in `yaml` kommt zurück | Startup (B9), Server-Parameter (B10), `offset_ms` (B11) fehlen; `meta` nicht gelesen (I20) | `TestLesenRoundtrip` |
| ungültiger JSON-Text und Text ohne gültiges UTF-8 `PGR-E3003`, auch tief verschachtelt | Syntax nicht vorab geprüft (J5); UTF-8 ungeprüft (J4) | `TestLesenJSONVerletzt`, `TestLesenJSONTief` |
| JSON-Zahl: ganz, kanonisch, 64 Bit mit Vorzeichen, nicht über `float64` | über `float64` (J1); nicht kanonisch angenommen (J2); außerhalb 64 Bit angenommen (J3) | `TestLesenJSONZahlen` |
| Umkehr der Umschreibung: `\` entfällt (eines), `base64:` dekodiert (Standard-Alphabet mit Auffüllung, sonst `PGR-E3003`), sonst wörtlich; Bytes ohne UTF-8 kommen zurück | `\` bleibt (J6); zwei Zeichen entfernt (J7); nicht dekodiert (J8); URL-Alphabet (J9); ohne Auffüllung (J10); Fehler ignoriert (J11); Bytes als Text (J12); Präfix ohne Beachtung der Schreibweise (J15) | `TestLesenStringOhneUTF8`, `TestLesenUmkehrDerUmschreibung` |
| dieselbe Ablehnung in beiden Formaten (unbekannter Schlüssel, `null`, doppelter Schlüssel, Version) | immer `yaml` (E6); immer `sqlite` (E7) | `TestLesenGleicheAblehnungBeiderFormate` |
| `replay` erkennt das Format selbst | `recording.YAML{}` statt `recording.Leser{}` in `replay` (K1) | am Binary `TestE2EReplaySQLite`, `TestE2EReplayExtendedSQLite`, `TestE2ELesenAbgelehnt` |
| `play` erkennt das Format selbst | `recording.YAML{}` statt `recording.Leser{}` in `play` (K2) | am Binary `TestE2EPlaySQLite`, `TestE2ELesenAbgelehnt` |

Summe: 92 Mutanten der Unit-Tests und 2 am Binary; rot 89, grün 5 (nächster Absatz; O5 ist nach dem Review rot, siehe *Nacharbeit zum Review*). Die Läufe sind Lauf-Beleg und
nicht im Repo; die Tabelle trägt Mutation und Test.

**Grüne Mutanten, eingeordnet.** Fünf blieben im ersten Lauf grün; vier ändern das Verhalten nicht, (O5) war nicht
äquivalent (Berichtigung im Abschnitt *Nacharbeit zum Review*; die Form `/C:/…` für Windows bleibt ungeprüft, Grenze
Plattform, `BEO-REPO/verhalten-nur-unter-linux-geprueft`). (I9) `group_no` ohne Prüfung der Art: Ein Wert, der keine ganze Zahl
ist, wird zu 0, und 0 lehnen die Regeln der Gruppen an jeder Stelle ab; die Prüfung gibt die Meldung, ändert nicht das
Ergebnis. (J13) Strings als `!!str` ohne Anführungszeichen: Der Encoder von `yaml.v3` setzt sie für jeden String, der
sonst anders gelesen würde, selbst; der Stil macht es unabhängig davon. (E11) Jeder Fehler beim Lesen der 16 Bytes gilt
als kurze Datei: `YAML.Load` liest die Datei danach selbst und meldet denselben `PGR-E3001` mit derselben Meldung; der
Fehler ist nur bei einer Datei beobachtbar, die sich zwischen den beiden Lesen ändert (nicht zugesagt). (E12) Der Fehler
des Schließens der Datei der Erkennung entfällt: Das Schließen einer nur gelesenen Datei scheitert in keinem Aufbau,
den ein Test herstellen kann.
Im ersten Lauf waren außerdem grün und sind geändert: die eigene Prüfung auf doppelte Schlüssel in einem JSON-Objekt, die
eigene Prüfung des Schlüssels `type` in `fields` und die auf negatives `offset_ms` — alle drei gelöscht, weil derselbe
Fehler über den Weg von `yaml` (`Unmarshal`) schon abgelehnt wird (die Fälle bleiben als Tests, `TestLesenJSONVerletzt`,
`TestLesenInhaltVerletzt`); Nachrichten ohne Interaktion (I6), `direction` bei `extended` (I8) und `NULL` in `group_no`
mitten in einer Gruppe (B3) ohne Fall — jetzt Fälle in `TestLesenInhaltVerletzt`.
**Nicht durch eine echte Datei erreicht:** `IOERR_SHORT_READ` liefern die abgeschnittenen Dateien nicht (die Bibliothek
meldet `SQLITE_NOTADB` unter 512 Bytes und sonst `SQLITE_CORRUPT`, weil die Seitenzahl der Kopfzeile größer ist als die
der Datei); die Zuordnung prüft `TestLesenStufeNachErgebniscode` mit der Zahl. Eine Datei, die mitten in einer
ungenutzten Seite endet (das Ende einer Indexseite, die das Lesen nicht berührt), liest als gültige Aufzeichnung; der
Test am Binary schneidet deshalb bei der halben Dateigröße.

**Läufe.** Alle im Docker-Netz oder im Image, nichts auf dem Host außer `make`, `docker` und `bash`.

| Lauf | Stand | Ergebnis |
|---|---|---|
| `make lint` | vor dem Aufteilen von `liesMeta` und `liesInteraktionen`: 4 Befunde (`containedctx`, zweimal `cyclop`, `gocognit`), behoben; danach | 0 Befunde |
| `make test-integration` | erster Lauf vor der Änderung des Schnitts in `TestE2ELesenAbgelehnt`: rot (der Fall schnitt in eine Indexseite); danach | grün (im Lauf von `make gates` 87 Tests, darunter die vier neuen); legt Netz, Volume und Container unter dem Namen `pgr-it-<pid>` an und entfernt sie danach |
| `make docs-check`, `make kopf-check`, `make abdeckung-check` | nach `make abdeckung` und dem Nachziehen von §3 | `d-check: 548 Datei(en) geprüft, 0 Befund(e)`; die beiden anderen Exit 0 |
| `make gates` | auf `45af81d`, Baum sauber | Exit 0, 6 min 11 s; `lint-gegenprobe: gruen`, `a-check-negativ: gruen — alle siebzehn Faelle`, `baseline-verify: v6.18.0 OK`, `kopf-check-gegenprobe: gruen`, `abdeckung-gegenprobe: gruen` |

Der Commit dieser Belege ändert nur diesen Plan; ein `make gates` auf ihm ist im Bericht der Übergabe genannt, nicht hier.

**Namensabgleich.** Ein Skript zog aus §3, §6 und diesem Abschnitt jeden Namen `Test…` und jeden Dateinamen in Backticks
und suchte `^func <Name>(` in `internal/**/*_test.go` und `test/integration/*_test.go` sowie die Dateien im Baum:
28 Testnamen und 24 Dateinamen gefunden, nichts fehlt (drei Dateien liegen im vendored Regelwerk, `interaction.sql` ist
eine Spalte, keine Datei).

**Hinweise an den Architect** (Randformen, die §6 nicht nannte) — **entschieden vom Architect am 2026-10-11**, Ort: `LH-FA-22.a` *Lesen von `sqlite`*, *Grenzen* und `SPEC-043`; die Aufträge stehen in §3:
1. *Maskiertes Surrogat im JSON-Text:* abgelehnt, `PGR-E3003`. Der Code liest heute U+FFFD; Auftrag in §3.
2. *Stelle in der Meldung:* keine Zeile, kein Wert; Tabelle und Nummern bei eigenen Prüfungen. Der Code gibt die Ursache der YAML-Bibliothek weiter; Auftrag in §3. Ob dieselbe Regel auch für `yaml` gilt, ist **entschieden vom Nutzer am 2026-10-11: nicht jetzt**; die Regel gilt vorerst nur für das Lesen von `sqlite`, das Thema geht als Beobachtung ins Register (der Planner trägt sie bei der Closure ein).
3. *Schließen nach dem Lesen:* nicht gemeldet, akzeptiertes Negativ (Code stimmt überein).
4. *Abbruch von `ctx`:* `PGR-E3001` (Code stimmt überein; Test fehlt, Auftrag in §3).
5. *Doppelter Schlüssel in `meta`:* `PGR-E3003` (Code stimmt überein).
6. *Gruppen:* Verstoß `PGR-E3003` (Code stimmt überein).
7. *Speicher:* ohne Obergrenze, in `LH-FA-22.a` *Grenzen* als Beschreibung festgehalten; **Entscheidung des Nutzers vom 2026-10-11:** „große Aufzeichnungen“ in `LH-FA-22` gilt nur für das Schreiben; der Leser darf wie der von `yaml` die Aufzeichnung im Speicher halten, ein Mehrfaches der Dateigröße, ohne zugesagte Obergrenze. Keine Änderung des Lastenhefts, kein streamendes Lesen. Die Frage an den Nutzer ist geschlossen.
8. *Größe:* 2090 Zeilen Go, über der Schwelle von §4. **Empfehlung** des Architect an den Planner (Schnitt entscheidet er nicht): hinnehmen. Die Überschreitung ist rund 5 % über „rund 2000“, 1291 der Zeilen sind Tabellen-Tests, 799 Code in einer Schicht mit drei Liefer-Punkten; der in §4 genannte Schnitt (DoD-Punkt 3) ließe sich nicht liefern, denn `stufeCode` und seine Tests hängen am Leser, den DoD-Punkt 1 liefert, und der Code liegt gate-grün vor. Ob der Diff in einer Sitzung prüfbar ist, beurteilt am Ende das Review.


### Nacharbeit zu den Entscheidungen des Architect (`ea1ef40`)

Ein Commit, `18695d4`: Surrogat, Meldungen, Abbruch von `ctx`, Doc-Kommentar zum Schließen. Mutationen wie oben (frische
Kopie, `gofmt -l` leer, Bind-Mount); die Kennung ist die des Laufs.

| Zusage | Mutation | roter Test |
|---|---|---|
| ein maskiertes einzelnes Surrogat ist `PGR-E3003`, als Wert und Schlüssel, in `startup` und tief in `fields` | Prüfung entfällt (SG1); tiefes allein angenommen (SG5); Obergrenze des tiefen (SG6) oder Untergrenze (SG7) nicht geprüft; nur `\ud800` als hohes erkannt (SG8) | `TestLesenSurrogate` |
| ein Paar, `\ufffd` und das wörtliche U+FFFD sind gültig | Paar abgelehnt (SG2); `\ufffd` abgelehnt (SG3); wörtliches U+FFFD abgelehnt (SG3b) | `TestLesenSurrogate` |
| ein maskierter Backslash vor `ud800` ist Text | Backslash nicht übersprungen (SG4) | `TestLesenSurrogate` |
| die Meldung nennt keinen Wert aus der Datei und keine Zeilennummer | Ursache der YAML-Bibliothek durchgereicht (ME1); Wert in eigener Meldung: Formatkennung (ME2), Schlüssel (ME3), `type` (ME4), `direction` (ME8), Zahl (ME5), Version (ME6); Zeilennummer (ME7) | `TestLesenMeldungOhneWertUndZeile` (achtzehn Fälle auf den drei Wegen: eigene Prüfung, `Unmarshal`, JSON) |
| Abbruch von `ctx`: `PGR-E3001` mit der Ursache | `PGR-E3003` (CX1) | `TestLesenNichtLesbarUndKopfOhneInhalt` |

Grün bleibt ein Mutant (ME9): Fällt `ladeInhalt` bei einem Fehler ohne `*model.Error` auf `PGR-E3001` statt `PGR-E3003`, ändert
sich nichts, denn `Unmarshal` liefert nur `*model.Error`; der Rückfall ist nicht erreichbar. Der doppelte Schlüssel in
`meta` und die Gruppen (Punkt 4 der Entscheidung) brauchten keinen Code; ihre Tests (`TestLesenMeta`,
`TestLesenInhaltVerletzt`) nennen die Spezifikationsstelle in ihrer Deklaration. Das Schließen der Datenbank nach dem Lesen
steht in der Doc von `leseAlsYAML`. **Größe jetzt:** 863 Zeilen Go-Code und 1375 Zeilen Go-Tests seit `97f3ca6`, zusammen 2238.
Die Meldungen der Dateien `sqlite_lesen*.go` und `sqlite_json_lesen.go` sind auf `%q`, `%v` und `%s` geprüft: Nach den
Änderungen steht nur noch Tabelle, Spalte, Nummer, Pfad und der Name der Konstante darin; die Ursache der Bibliothek
(`sqlite`-Fehler, Base64) hängt an Meldungen zum Lesen und zu JSON, nicht an `Unmarshal`.

### Nacharbeit zum Review (`415b725`)

Ein Commit, `415b725`: F-623 bis F-626; F-627 ist dieser Plan. Mutationen wie oben (frische Kopie, `gofmt -l` leer,
Bind-Mount).

| Zusage | Mutation | roter Test |
|---|---|---|
| jede Meldung zu einem Fehler des Inhalts nennt den Pfad der Datei: bei eigenen Prüfungen des Lesers und bei der unbekannten Version | Pfad entfällt bei eigenen Fehlern (MP1); nur bei `PGR-E3003` (MP5); Version-Fehler ohne Pfad (MP4) | `TestLesenMeldungOhneWertUndZeile` (Pfad je Fall) |
| … und auf dem Weg über `Unmarshal` Pfad und Art (der Text des Fehlers von `Unmarshal`, etwa `Session 2 ohne Interaktion`) | Pfad entfällt (MP2); Art entfällt (MP3) | `TestLesenMeldungOhneWertUndZeile` |
| … und bei eigenen Prüfungen die Art | nur der Pfad (MP6) | `TestLesenMeldungOhneWertUndZeile` (`enthaelt` je Fall) |
| der Eingabepfad gilt wörtlich auch mit führendem `//`, beim Schreiben, Ergänzen und Lesen | `filepath.Abs` entfällt (O5); `file:`-URI des Lesens aus dem rohen Pfad (O5b) | `TestLesenPfadMitFuehrendemDoppelstrich` |
| Abbruch von `ctx` zwischen den Tabellen: `PGR-E3001` mit Ursache und Pfad | `PGR-E3003` (AB1) | `TestLesenAbbruchWaehrendDesLesens` |

**Berichtigung zu O5 (F-624).** `filepath.Abs` ist nicht äquivalent: Ein Pfad mit führendem `//` wird ohne `Abs` als
`file://tmp/…` mit Authority gelesen und abgelehnt. Der Test deckt Lesen und Schreiben (`uriVon` bedient beide); der Mutant
`W10` des Schreib-Slice ist dieselbe Ersetzung und mit diesem Test ebenfalls rot, `slice-v1-abschluss-sqlite-format` bleibt
unverändert.
**U7 (F-625).** `rows.Err()` in `zeilen` entfällt: grün. Ein Abbruch oder Fehler mitten in einer Tabelle fällt in
`QueryContext` oder `Scan`, nicht in `rows.Err()`; ich habe den Zweig mit abgebrochenem `ctx` zwischen den Tabellen, mit
beschädigten Seiten und mit Dateien aus 60 Sessions nicht erreicht. Der Zweig bleibt als Absicherung der Bibliothek
(Vertrag von `database/sql`) und ist durch keinen Test belegt; die Zusage der Spezifikation „während des Lesens“ trägt der
Test zwischen den Tabellen.
**Meldung auf dem Weg über `Unmarshal` (F-623).** Die Art ist dort so fein wie der Text, den `Unmarshal` ohne Wert und
Zeile trägt: `Aufzeichnung beschädigt` (Form, unbekannter Schlüssel) oder Session und Interaktion; eine feinere Art nähme die
Ursache der YAML-Bibliothek mit Wert und Zeile. Frage an den Architect, ob das die Zusage „die Art des Fehlers“ trägt.
**Kommentare (F-626).** `Load` nennt, was der Leser selbst prüft und was `Unmarshal`; `JSONKnoten` nennt `!!map` und `!!seq`
mit leerem Wert.

### Entscheidung des Architect zu F-630

Der Befund liest den Leser falsch. Der Leser baut aus der Datei einen YAML-Text und liest ihn mit `Unmarshal`; `ohneVerweise`
vergleicht den Wert des Schlüssels (`<<`), nicht seine Schreibweise, und lehnt auch den Schlüssel in Anführungszeichen ab. Gemessen
an einer Kopie des Arbeitsbaums: eine Datei mit `startup = {"<<":"x"}` liest `Leser.Load` mit `PGR-E3003`, `SQLite.Write` lehnt denselben Namen
mit `PGR-E1000` ab (`Tag !!merge ist kein String`), `YAML.Marshal` schreibt `<<: x`, das `Unmarshal` mit `PGR-E3003` ablehnt.
**Verdikt:** akzeptiertes Negativ für das Schreiben (kein Postgres-Parametername heißt `<<`, ohne Wiederkehr; die Ablehnung eines
Schlüssels, den `yaml` nicht liest, ist konsistent), der Leser bleibt bei `PGR-E3003`. Festgehalten in `LH-FA-22.a` *Lesen von `sqlite`*
und §6. Folge: ein Test ohne Produktcode, Auftrag in §3 letzter Punkt vor der Randform-Rückgabe. Die Spezifikation sagt nur das Lesen zu;
dass `YAML.Write` den Namen schreibt und nicht lesen kann, ist Bestand von `slice-walking-skeleton-record` und bleibt stehen
(Postgres-Namen ohne `<<`, harmlos).

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

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-sqlite-format` (`AGENTS.md` §3.13): Liefer-Punkte: 3 (Lesen und Roundtrip, Ablehnung beider Formate, Einstufung der Fehler). Schichten: eine (Recording-Adapter mit seinen Tests); der Bootstrap reicht den Leser nur durch (reine Verdrahtung, zählt nicht), der Integrationstest läuft über das Binary und trägt keine eigene Schicht; die Dokumentation liegt im Doku-Folge-Slice `slice-v1-abschluss-sqlite-format-doku`. Geschätzt 1300 bis 1700 Zeilen mit Tests (Nachzählung des Architect vom 2026-10-11: der Geber lag bei rund 2300 mit Tests für Schreiben, JSON und Zeilen; der Leser trägt weniger Code, aber die Tests der Randformen; unter einer Review-Sitzung von rund 2000). Gemessen nach der Nacharbeit zum Review (`git diff --numstat 97f3ca6..415b725 -- '*.go'`): 2311 Zeilen Go, davon 886 Code und 1425 Tests; Liefer-Punkte und Schichten wie gezählt (drei, eine). Der erste Stand (`45af81d`) lag bei 2090 (799 und 1291), §7 führt die Stufen.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält die Aufzeichnung als Code unter Spezifikation und Entscheidungen).
