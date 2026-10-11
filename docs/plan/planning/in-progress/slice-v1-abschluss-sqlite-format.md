# Slice slice-v1-abschluss-sqlite-format: SQLite als Aufzeichnungsformat

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [ADR-0005](../../adr/0005-recording-store-ist-driven-adapter.md), [ADR-0013](../../adr/0013-sqlite-recording-backend.md)

**Berührte Spec-Stellen:** `LH-FA-22.a` · `LH-FA-07.a` · `SPEC-001` · `SPEC-043` · `SPEC-044` · `ARC-008` · `ARC-014`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `record --format sqlite` speichert die Aufzeichnung je Session in einer Transaktion in einer SQLite-Datei, die zu jedem Zeitpunkt nur vollständige Sessions enthält.

**Abgegeben** (Schnitt vom 2026-10-11, §4): Das Lesen von `sqlite` für `replay` und `play` ist in `slice-v1-abschluss-sqlite-format-lesen` (welle-v1-abschluss, `next/`, direkt hinter diesem Slice) aufgegangen — Grund: Größe über einer Review-Sitzung, auf 2500 bis 3000 Zeilen geschätzt (Entscheidung des Nutzers vom 2026-10-11: Schnitt *Schreiben | Lesen*). Dorthin gehen, jeweils mit der Kennung dieses Slice als Geber in dessen §1 und DoD:

- die Formaterkennung an den ersten 16 Bytes für `replay` und `play` und das Lesen in einer Lesetransaktion (`LH-FA-22.a` *Wahl und Erkennung*, *Lesen von `sqlite`*);
- die Umkehr des `base64:`-Escapes beim Lesen und das Lesen von `{}`;
- `PGR-E3002` und `PGR-E3003` beim Lesen, `meta.version`, `meta.empty_sessions` beim Lesen, Waisen, die Einstufung der Fehler des Lesens nach dem Ergebniscode;
- die Roundtrip-Gleichheit beider Formate (Abnahmeszenario 14) und die gleiche Ablehnung beider Formate.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Lesen von `sqlite` (Liste oben) — `slice-v1-abschluss-sqlite-format-lesen` (dort §1 und DoD, nimmt es an), das direkt hinter diesem Slice steht: Mit dem Lesen läge dieser Plan über der Grenze einer Review-Sitzung. Bis zu jenem Slice liest `replay` eine Datei von `record --format sqlite` als YAML und meldet sie als beschädigt; das ist wellenintern und ohne Release dazwischen akzeptiert (Entscheidung des Nutzers vom 2026-10-11).
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-sqlite-format-doku` (dort §1 und DoD), hinter `slice-v1-abschluss-sqlite-format-lesen` in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Das Handbuch beschreibt erst, was das Lesen geliefert hat; zwischen den Slices hinkt es dem Binary hinterher, ein Release liegt nicht dazwischen.
- Umwandlung zwischen den Formaten und weitere Formate — Out-of-Scope von LH-FA-22.
- Das Standardformat YAML — bleibt unverändert aus `slice-walking-skeleton-record`.

**Übernommen von `slice-v1-abschluss-schreiben`:** Was von `LH-FA-07.a` *Temporäre Datei*
(Name, Rechte, Fehlschlag, übrig gebliebene Datei) und von der Probedatei aus *Zielpfad beim
Start* für `sqlite` gilt, nennt §6 dieses Slice als Randformen und entscheidet der Architect
in `LH-FA-22.a` *Schreiben von `sqlite`*, vor dem ersten Code-Commit (`AGENTS.md` §3.12).
Dort gilt heute nur: Zielpfad wie in `LH-FA-07.a`, temporäre Datei im Verzeichnis der
Zieldatei, atomar verschoben.
Dazu, mit der Kennung `slice-v1-abschluss-schreiben` (Entscheidung des Nutzers vom
2026-10-09): die Grenze der Prüfung für `sqlite` — ob für das atomare Verschieben in
`LH-FA-22.a` *Schreiben von `sqlite`* dieselbe Grenze gilt wie in `LH-FA-07.a` Schritt 2 (geprüft
nur unter Linux im selben Dateisystem, macOS und Windows nicht) und ob sie für das Ergänzen
einer Session in einer Transaktion eine eigene braucht. §6 nennt sie als Randform, der Architect
dieses Slice entscheidet sie in `LH-FA-22.a`, vor dem ersten Code-Commit (`AGENTS.md` §3.12).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat): `record --format sqlite` schreibt die Aufzeichnung je Session in einer Transaktion; die Datei enthält zu jedem Zeitpunkt nur vollständige Sessions, wird atomar angelegt und nie als Ganzes neu geschrieben, ein Rollback-Journal neben dem Zielpfad ist `PGR-E3001`; `--format` hat die Wertemenge `yaml`/`sqlite`, und `replay` und `play` kennen es nicht (Test über das Binary; Inhalt der Datei gegen die Liste in `SPEC-043` mit rohem SQL, ohne `sqlite3`).
- [ ] Das Schema liegt in `tools/schema/schema.yaml` (`message.bytes` entfällt, `fields` ist `required`, `session.server_parameters` kommt hinzu), das eingebettete SQL ist daraus erzeugt, und ein Test legt es in einer In-Memory-Datenbank an und prüft Tabellen, Spalten sowie Primär- und Fremdschlüssel gegen `SPEC-043`.
- [ ] Die SQLite-Bibliothek ist im Architektur-Gate (`.a-check.yml`, `tech`-Regel) auf den Recording-Adapter begrenzt; eine absichtliche Verletzung lässt `make a-check` fehlschlagen (Test), und die Erlaubnisliste des Lint-Gates führt sie.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/recording` | update | zweiter Adapter hinter `RecordingRepository`: Schreiben von SQLite, Transaktion je Session (`sqlite.go`, `sqlite_zeilen.go`, `sqlite_json.go`); das eingebettete SQL liegt als `schema.sql` neben dem Code (`go:embed`); Erkennung und Lesen liegen in `slice-v1-abschluss-sqlite-format-lesen` |
| `internal/adapters/driving/cli` | update | Option `--format` von `record` (Wertemenge, Hilfetext, `RecordOptions`); `replay` und `play` kennen sie nicht; `Formate()` nennt die Wertemenge für den Test der gemeinsamen Zusage der Fabrik (§6) |
| `tools/schema/schema.yaml` | update | Das Schema liegt für diesen Slice NICHT vor: `message.bytes` entfällt, `message.fields` wird `required`, `session.server_parameters` kommt hinzu (`SPEC-043`; Entscheidung des Nutzers vom 2026-10-11: `bytes` gestrichen). Das SQL für SQLite wird daraus erzeugt und im Adapter eingebettet |
| `internal/hexagon/ports/driven/recording.go`, `internal/hexagon/services/record.go` | update | Der Kommentar von `Write` sagt „als Ganzes“; das stimmt für `sqlite` nicht (Ergänzen je Session) und wird nachgezogen (`AGENTS.md` §3.7, §3.11); ebenso der Kommentar am Typ `RecordService` („schreibt die Aufzeichnung als Ganzes“ wird „übergibt sie als Ganzes“). Nur Kommentare, keine Schicht |
| `internal/adapters/driven/recording`, `internal/bootstrap` | update | Fabrik `recording.New(format)` (`format.go`) für das Schreiben; der Bootstrap reicht nur `o.Format` durch und kennt keinen Adapter (reine Verdrahtung, `replay` und `play` behalten bis `slice-v1-abschluss-sqlite-format-lesen` den Leser von `yaml`); der Test der gemeinsamen Zusage und der Lauf von `record` mit beiden Formaten liegen in `internal/bootstrap/format_test.go` |
| `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk` (Kopf), `harness/README.md` (Zeile `make a-check-negativ`) | update | Die Gegenprobe nennt `modernc.org/sqlite`: im Recording-Adapter zugelassen, im Domain Model, im Postgres-Adapter, im PGWire-Adapter und im CLI-Adapter abgelehnt; Skript, Kopfkommentar und Gate-Zeile sagen dasselbe (`AGENTS.md` §3.11) |
| `go.mod`, `go.sum` | update | per `make go-mod-tidy`; `modernc.org/libc` exakt in der Version, die die `go.mod` der Bibliothek festlegt (v1.77.1, [ADR-0024](../../adr/0024-sqlite-bibliothek.md)); `golang.org/x/sync` steigt indirekt von v0.17.0 auf v0.23.0, weil `modernc.org/libc` es verlangt |
| `.a-check.yml` (`tech`-Regel) | prüfen | Die Regel für `modernc.org/sqlite` liegt vor; der Slice belegt sie mit der Gegenprobe oben |
| `.golangci.yml`, `tools/harness/lint-gegenprobe.sh` | update | `gomodguard_v2.allowed` führt `modernc.org/sqlite` im selben Commit wie `go.mod` ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Punkt 3, `SPEC-049`); die Gegenprobe bekommt den Fall `p5-gomodguard-frei-sqlite` neben `p5-gomodguard-frei-yaml`. Das ist der vorgesehene Nachzug einer Erlaubnisliste, keine Lockerung (`AGENTS.md` §3.6) |
| `test/integration` | update | `sqlite_format_e2e_test.go`: `record --format sqlite` über das Binary: die Datei beginnt mit `SQLite format 3` und einem NUL-Byte, daneben bleibt kein `-journal`, ein vorhandener Journal-Rest ist `PGR-E3001`, `--format` mit falschem Wert `PGR-E2001`, `replay --format` eine unbekannte Option; ohne `sqlite3`-Binary und ohne Änderung an `tools/test/run-integration-tests.sh`. Der Roundtrip über `replay` und `play` gehört `slice-v1-abschluss-sqlite-format-lesen` |
| `docs/user/abdeckung-*.md` | update | über `make abdeckung` aus den Deklarationen der neuen Tests (erzeugt, kein Handbuch) |

**Aufträge an den Implementer** (vom Architect vor dem Code gesetzt; Quelle der Regeln: `LH-FA-22.a`, `SPEC-043`):

- `tools/schema/schema.yaml` ändern: Spalte `bytes` entfernen, `fields` `required`, `session.server_parameters` ergänzen; das eingebettete SQL daraus neu erzeugen (`make schema-generate`).
- Ersatztest für das eingebettete SQL (Entscheidung des Nutzers vom 2026-10-11: Drift zwischen `schema.yaml` und eingebettetem SQL ist akzeptiert, der Test ersetzt den Abgleich): Er legt das eingebettete SQL in einer In-Memory-Datenbank an und prüft Tabellen, Spalten sowie Primär- und Fremdschlüssel gegen die Liste in `SPEC-043`.
- Transaktion je Session mit `BEGIN IMMEDIATE`; der Pfad der Datei geht als korrekt escapte `file:`-URI an die Bibliothek (Test mit `?`, `#`, `%` und Leerzeichen im Pfad).
- Je Zusage dieses Vertrags ein Test und eine selbst gefahrene Mutation, im Bericht Zusage · Mutation · roter Test (`AGENTS.md` §3.10); dazu die gemeinsame Zusage der Fabrik (siehe §6).
- Die Doc-Kommentare nachziehen, die der Slice ändert oder die ihm widersprechen (`AGENTS.md` §3.7, §3.11), insbesondere den Port-Kommentar von `Write`.
- Strings ohne gültiges UTF-8 (`SPEC-043`, Regel `base64:` mit Umschreibung), Schreib-Richtung: ein Codec für Werte und Schlüssel an der Stelle, an der der `yaml.Node` zum JSON-Baum wird (Tag `!!binary` → `base64:`); Tests mindestens für das Byte `0xE9` in einem Startup-Wert, einem Fehlerfeld und einem Schlüssel, für einen gültigen String, der mit `\` oder `base64:` beginnt (steht mit einem weiteren `\`), gelesen mit rohem SQL; `interaction.sql` mit `0xE9` steht unverändert in der Datei. Die Umkehr beim Lesen, ungültiges Base64 (`PGR-E3003`) und das Hin und Zurück liefert `slice-v1-abschluss-sqlite-format-lesen`. Liefert die Bibliothek die Bytes von `interaction.sql` nicht unverändert zurück, gibt der Implementer es dem Architect zurück.
- Der Codec bleibt ein Pfad (Rückführung §4): die DTOs von `yaml` ↔ `yaml.Node` ↔ JSON-Baum. Die Umwandlung kennt nur JSON-Typen (Zahl → `!!int`, String → `!!str` in Anführungszeichen, Wahrheitswert, `null`), keine Nachrichtenart; braucht sie ein Wissen über ein Feld, ist das die Rückführung. Die Gegenrichtung (Lesen) und der Test auf gleiche Ablehnung beider Formate liegen in `slice-v1-abschluss-sqlite-format-lesen`; die Schreib-Richtung legt den Baum so an, dass jene Richtung ihn umkehrt, ohne ihn zu ändern.
- Datei als `file:`-URI mit `mode=rw` zum Ergänzen (ohne `mode` legt SQLite still an), `_txlock=immediate` nur für schreibende Transaktionen, Sperrfrist fünf Sekunden; `context.WithoutCancel` für die Transaktion. Jeder Fehler der Bibliothek beim Schreiben ist `PGR-E3001` (`LH-FA-22.a` *Schreiben von `sqlite`*). `mode=ro`, die Fehlerklasse beim Lesen und die abgeschnittene Datei liefert `slice-v1-abschluss-sqlite-format-lesen`.
- Ob `sqlite3` die Datei liest, prüft dieser Slice nicht (kein Hilfscontainer): die Unabhängigkeit vom Lesecode des Adapters liefert ein Test im Paket des Adapters, der die Datei mit rohem SQL gegen die Liste in `SPEC-043` liest. Dieser Slice hat keinen Leser; was der Test braucht, liest er mit der Bibliothek selbst.
- `Load` des SQLite-Adapters (Entscheidung des Architect vom 2026-10-11, vor dem ersten Code-Commit): Der Port bleibt unverändert; `Load` ist ein Platzhalter und liefert immer `model.CodeInternal` (`PGR-E1000`) mit dem Text „Lesen von sqlite ist in diesem Stand nicht verfügbar“, ohne Pfad und ohne Wert, kein Pfad des Produkts ruft es. Ein Test belegt Code und Text (Mutation: Platzhalter liefert nil → rot); `slice-v1-abschluss-sqlite-format-lesen` ersetzt ihn (dort §1 und §6). Die Spezifikation nennt ihn nicht.
- `.golangci.yml` und `lint-gegenprobe.sh` nachziehen (siehe Tabelle).
- Nummern der Sessions (Entscheidung des Architect vom 2026-10-11, `LH-FA-22.a` *Schreiben von `sqlite`*): ein Test, dass ein `Write` mit einer doppelten Session-Nummer `PGR-E3001` liefert und die Datei auf dem Stand vor dem Aufruf lässt, und ein Test, dass eine Lücke (etwa `1` und `3`) spätestens beim nächsten `Write` `PGR-E3001` liefert; je eine Mutation (Primärschlüssel von `session` entfernt, Prüfung der Nummern vor dem Ergänzen entfernt) und im Bericht Zusage · Mutation · roter Test. Dazu ein Test für das Journal der temporären Datei (`record` rührt eine fremde `.<Name>.<Hex>.tmp-journal` nicht an) nach `LH-FA-07.a` *Übrig gebliebene Datei*: eine Datei dieses Namens im Verzeichnis bleibt nach `Prepare` und `Write` unverändert.
- Eine Randform, die §6 nicht nennt, entscheidet er nicht, sondern gibt sie dem Architect zurück (`AGENTS.md` §3.12).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-walking-skeleton-record` ist `done` (YAML-Format liegt vor), und [ADR-0013](../../adr/0013-sqlite-recording-backend.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: die SQLite-Bibliothek verlangt native Abhängigkeiten für die Zielplattformen — zurück zur Zerlegung.
- `in-progress` → `next`: Der Slice braucht mehr als einen Codec-Pfad (das Schreiben von `sqlite` trägt das Modell nur mit einer zweiten, eigenen Abbildung der Nachrichten neben `yaml`; für das Lesen gilt dieselbe Bedingung in `slice-v1-abschluss-sqlite-format-lesen` §4) — zurück zur Zerlegung.
- `in-progress` → `next` (zu groß): Der Diff des Schreibens ist nicht in einer Review-Sitzung prüfbar; die Schätzung liegt bei rund 1800 Zeilen. Schnitt dann: DoD-Punkt 3 (Gate, Gegenprobe, Erlaubnisliste) als eigener Slice.
- `in-progress` → `open`: Beide Formate tragen das Modell nicht gleich — Entscheidung klären.

**Eingetreten** (Rückführung `in-progress` → `next`, zu groß; 2026-10-11, vor dem ersten Code-Commit des Implementers): Nach den Randformen des Architect (`bb67108`) schätzte der Architect den Diff des ganzen Slice auf 2500 bis 3000 Zeilen mit Tests, über einer Review-Sitzung, obwohl der Plan drei Liefer-Punkte und zwei Schichten zählte (die Größenregel zählt Punkte, nicht Zeilen). Entscheidung des Nutzers vom 2026-10-11: Schnitt nach Option 1, *Schreiben | Lesen*; der Zwischenstand „sqlite schreibbar, aber nicht wiedergebbar“ ist wellenintern akzeptiert, ein Release liegt nicht dazwischen. Dieser Slice behält das Schreiben (rund 1800 Zeilen), das Lesen (900 bis 1000 Zeilen) ist `slice-v1-abschluss-sqlite-format-lesen` in `next/`. Der Slice lief über `next/` zurück und startet von dort neu (Anspruch auf dem Hauptzweig vor der Arbeit, Modul 5). Beobachtung für die Closure: `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2. Beleg, 1. Beleg `slice-v1-abschluss-einspielen`).

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

- **Risiko: Grenze der Prüfung für atomares Verschieben und Transaktion je Session** (aus `slice-v1-abschluss-schreiben`, §1 *Übernommen*) — **Ausgang: entfallen**, entschieden in `LH-FA-22.a` *Schreiben von `sqlite`*, *Grenzen* (Verschieben wie `LH-FA-07.a` Schritt 2, Transaktionen nur unter Linux auf einem lokalen Dateisystem geprüft, Stromausfall nicht belegt).
- Die gewählte Bibliothek (reines Go, siehe [ADR-0024](../../adr/0024-sqlite-bibliothek.md)) ist für die sechs Zielplattformen erst durch den Build im Slice belegt — **Ausgang:** offen bis Closure (Register).
- **Risiko: Zwischenstand nach dem Schnitt vom 2026-10-11:** Bis `slice-v1-abschluss-sqlite-format-lesen` liest `replay` eine Datei von `record --format sqlite` als YAML und meldet sie als beschädigt (`PGR-E3003`); die Spezifikation beschreibt den Endstand und nennt keinen Zwischenstand. Wellenintern und ohne Release dazwischen akzeptiert (Entscheidung des Nutzers vom 2026-10-11) — **Ausgang: eingetreten:** `slice-v1-abschluss-sqlite-format-lesen`.

**Randformen des neuen Vertrags** (`AGENTS.md` §3.12) — je Randform der Ort der Entscheidung; die Spezifikation steht vor dem ersten Code-Commit. Die Randformen des Lesens sind mit dem Schnitt vom 2026-10-11 nach `slice-v1-abschluss-sqlite-format-lesen` §6 gegangen (Typ-Prüfung der Spaltenwerte, `fields` beim Lesen, Lesetransaktion, `meta.version` und `empty_sessions` beim Lesen, Waisen und Zeilen gegen die Regeln, Fehlerklasse beim Lesen, Umkehr des `base64:`-Escapes); die Entscheidungen des Architect gelten dort unverändert. Entscheidungen des Nutzers vom 2026-10-11: `message.bytes` entfällt (`fields` ist die einzige Darstellung), und der Ersatztest statt des Abgleichs mit dem eingebetteten SQL (§3).

| Randform | Entschieden in |
|---|---|
| Rollback-Journal `<Zielpfad>-journal` beim Start (`PGR-E3001`, auch mit `--force` und ohne Zieldatei) | `LH-FA-22.a` *Schreiben von `sqlite`* |
| Nachholen nach fehlgeschlagener Transaktion (nächster Schreibvorgang ergänzt, was fehlt) | `LH-FA-22.a` *Schreiben von `sqlite`* |
| Verschwundene oder ersetzte Datei beim Ergänzen (`PGR-E3001`, Datei unberührt) | `LH-FA-22.a` *Schreiben von `sqlite`* |
| `fields` steht in jeder geschriebenen Zeile als JSON-Objekt (bei einer Nachricht ohne Felder `{}`), kein `bytes` | `SPEC-043` |
| Hilfe-Text und `RecordOptions` zu `--format`; Wertemenge, leerer Wert, `SQLite`, `sqlite3` (`PGR-E2001`); `replay`/`play` kennen `--format` nicht; die Erkennung an den ersten 16 Bytes liegt in `slice-v1-abschluss-sqlite-format-lesen` | `LH-FA-22.a` *Wahl und Erkennung* |
| Abbruchsignal unterbricht keine Transaktion | `LH-FA-22.a` *Schreiben von `sqlite`* |
| Pfad mit Sonderzeichen (`?`, `#`, `%`, Leerzeichen) gilt wörtlich | `LH-FA-22.a` *Schreiben von `sqlite`* |
| Werte nur als Parameter vorbereiteter Anweisungen (SQL-Injection) | `LH-FA-22.a` *Schreiben von `sqlite`* |
| `session.server_parameters` im Schema | `SPEC-043`; Nachzug in `tools/schema/schema.yaml` (Auftrag in §3) |
| Strings ohne gültiges UTF-8 in `startup`, `server_parameters`, `fields` (Wert und Schlüssel) beim Schreiben als `base64:`; `interaction.sql` roh | `SPEC-043` |
| Reihenfolge Pfadprüfung vor Journal; Prüfung der Datei auch ohne Ergänzung; Fehler beim Schließen nach der Transaktion; jeder Fehler des Schreibens `PGR-E3001` | `LH-FA-22.a` *Schreiben von `sqlite`* |
| `modernc.org/sqlite` in der Erlaubnisliste des Lint-Gates | `SPEC-049`; Nachzug in `.golangci.yml` (Auftrag in §3) |
| `meta.version` und `meta.empty_sessions` beim Schreiben (`empty_sessions` nur mit `--record-empty-sessions`, dann `true`) | `SPEC-043` |
| Fabrik `recording.New(format)` im Adapter; der Bootstrap reicht nur `o.Format` durch | `ARC-008` (Adapter-Grenze); hier §3 |
| Gemeinsame Zusage: jeder Wert, den die CLI annimmt, nimmt `New` an | Zusage dieses Slice; Mutation: ein Wert der CLI-Wertemenge aus `New` streichen, der Test über alle CLI-Werte wird rot |
| „Vorhanden“ beim Journal `<Zielpfad>-journal` wie beim Zielpfad (`Stat`: Verknüpfung nach ihrem Ziel, Verzeichnis vorhanden, Verknüpfung ins Leere nicht, nicht feststellbar `PGR-E3001`) | `LH-FA-22.a` *Schreiben von `sqlite`* (Entscheidung des Architect vom 2026-10-11) |
| Aufruf mit weniger Sessions, als der Recorder zuletzt ergänzt hat: `PGR-E1000`, Datei unberührt | `LH-FA-22.a` *Schreiben von `sqlite`* (Entscheidung des Architect vom 2026-10-11) |
| Nummer der Session, die nicht die nächste ist: doppelt `PGR-E3001` (Primärschlüssel), Lücke `PGR-E3001` spätestens beim nächsten Schreibvorgang | `LH-FA-22.a` *Schreiben von `sqlite`* (Entscheidung des Architect vom 2026-10-11); Test siehe §3 |
| Rollback-Journal der temporären Datei nach einem Absturz mitten im Anlegen: übrig gebliebene Datei, `record` sucht und entfernt es nicht | `LH-FA-22.a` *Schreiben von `sqlite`*, mit `LH-FA-07.a` *Übrig gebliebene Datei* (Entscheidung des Architect vom 2026-10-11) |
| `Load` des SQLite-Adapters bis `slice-v1-abschluss-sqlite-format-lesen`: Platzhalter mit `PGR-E1000`, Port unverändert; kein Pfad des Produkts ruft ihn; die Spezifikation beschreibt den Endstand und nennt ihn nicht | Entscheidung des Architect vom 2026-10-11, hier §3; Übernahme durch `slice-v1-abschluss-sqlite-format-lesen` §1 |

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

**Stand.** Geliefert sind die drei Liefer-Punkte der DoD (Schreiben mit Fabrik und Option, Schema mit
Ersatztest, Architektur-Gate mit Erlaubnisliste) in vier Commits seit `4914a4a`: `f052e63`
(Recording-Adapter, Schema, `go.mod`, Erlaubnisliste, Kommentare an Port und `RecordService`, Plan §3),
`fc59441` (`--format`, Bootstrap, Test der gemeinsamen Zusage, Tests über das Binary), `07d34c4`
(Gegenprobe des Architektur-Gates), `9ebab6c` (die Kopfzeile im Test über das Binary als Konstante: der
erste `make gates` meldete `gochecknoglobals`, weil `make lint` vor dem Anlegen der Datei lief) und der
Commit dieser Belege. Neu sind `internal/adapters/driven/recording/sqlite.go`, `sqlite_zeilen.go`,
`sqlite_json.go`, `format.go`, `schema.sql`, `sqlite_test.go`, `sqlite_schreiben_test.go`,
`sqlite_einstellungen_test.go`, `sqlite_json_test.go`, `sqlite_hilfen_test.go`,
`internal/bootstrap/format_test.go` und `test/integration/sqlite_format_e2e_test.go`; geändert
`export_test.go`, `cli.go`, `bootstrap.go`, der Port-Kommentar, `go.mod`, `go.sum`, `.golangci.yml`,
`tools/harness/lint-gegenprobe.sh`, `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk`,
`harness/README.md`, `tools/schema/schema.yaml` und `docs/user/abdeckung-*.md`. Port `Write`,
Handbuch und `README.md` sind unverändert. **Größe:** `git diff --numstat 4914a4a..9ebab6c` zählt
2794 Zeilen hinzu, davon 676 Go-Code und 1975 Go-Tests; die Schätzung in §4 war rund 1800. Die
Rückführung *zu groß* aus §4 hängt an der Prüfbarkeit in einer Review-Sitzung, und die entscheidet das
Review; die Tests sind überwiegend Tabellen mit rohem SQL.

**Zusage · Mutation · roter Test.** Jede Mutation lief in einer frischen Kopie des Baums (`cp -r` ohne `-p`,
`go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`), mit einer Ersetzung, `gofmt -l` der Datei leer und
`go test -count=1 -run <Tests>` im gepinnten Go-Image mit Bind-Mount (der Mutant kommt ohne Build-Kontext an);
`E1` bis `E5` liefen mit `make test-integration` in einer Kopie des ganzen Baums. Die Kennung in Klammern
ist die des Laufs.

| Zusage | Mutation | roter Test |
|---|---|---|
| `--format` nimmt genau `yaml` und `sqlite`, sonst `PGR-E2001` | `Set` vergleicht ohne Groß-/Kleinschreibung (C1, E5); `Set` nimmt jeden Wert an (C6) | `TestParseRecordFormat`, `TestLeserAlleOptionen`; am Binary `TestE2EFormatUngueltig` |
| `sqlite` ist ein Wert von `--format` | `Formate()` ohne `sqlite` (C2) | `TestFormate`, `TestParseRecordFormat`, `TestLeserAlleOptionen` |
| Standard ist `yaml` | Standard `sqlite` (C3) | `TestParseRecordFormat`, `TestLeserAlleOptionen`, `TestParseRecord` |
| der Wert gelangt nach `RecordOptions.Format` | `setze` schreibt immer `yaml` (C7) | `TestParseRecordFormat`, `TestLeserAlleOptionen` |
| `replay` und `play` kennen `--format` nicht | `replay` bekommt die Option (C4, E4) | `TestLeserOptionen`; am Binary `TestE2EFormatUngueltig` |
| die Hilfe von `record` nennt `--format yaml\|sqlite` und die Umgebungsvariable | Hilfetext entfernt (C5) | `TestHilfeFormat` |
| **gemeinsame Zusage:** jeder Wert der Kommandozeile nimmt `New` an | `sqlite` aus `New` gestrichen (F1, F1b) | `TestFormateDerKommandozeileKenntDieFabrik` (Bootstrap), `TestNewFormat` |
| `New` lehnt jeden anderen Wert mit `PGR-E2001` ab | unbekannter Wert liefert YAML (F2); Code `PGR-E1000` (F3) | `TestNewFormat` |
| der Bootstrap reicht `o.Format` durch | `recording.New("yaml")` (B1, E1) | `TestRunRecordFormat`; am Binary `TestE2ERecordFormatSQLite`, `TestE2ERecordFormatQuellen`, `TestE2ERecordFormatSQLiteJournal` |
| `<Zielpfad>-journal` beim Start `PGR-E3001` | anderer Name geprüft (P1); Suffix `.journal` (P7); nur ohne `--force` (P3); kein Zweig für „nicht prüfbar“ (P5); Prüfung weg (E2) | `TestSQLitePrepareJournal`, `TestSQLitePrepareJournalReihenfolge`; am Binary `TestE2ERecordFormatSQLiteJournal` |
| das Journal wird nach dem Zielpfad geprüft | Journal zuerst (P2b) | `TestSQLitePrepareJournal`, `TestSQLitePrepareJournalReihenfolge` |
| „vorhanden“ wie beim Zielpfad (Verknüpfung nach ihrem Ziel) | `Lstat` statt `Stat` (P4) | `TestSQLitePrepareJournalReihenfolge` |
| der Zielpfad wird wie bei `yaml` geprüft | `pruefe` aus `Prepare` entfernt (P6) | `TestSQLitePrepareZielpfad`, `TestSQLitePrepareJournal`, `TestSQLitePrepareJournalReihenfolge` |
| der erste Schreibvorgang verschiebt die temporäre Datei atomar | `os.Rename` am Hook vorbei (W32) | `TestSQLiteNameDerTempDatei`, `TestSQLiteErsterSchreibvorgangScheitert` |
| die Datei wird nie als Ganzes neu geschrieben | jeder Aufruf legt die Datei neu an (W31) | `TestSQLiteErgaenzt` (dieselbe Datei) |
| jeder Schreibvorgang ergänzt genau die fehlenden Sessions | alle Sessions statt der fehlenden (W36); Zähler nicht erhöht (W20); umgekehrte Reihenfolge (W23) | `TestSQLiteErgaenzt`, `TestSQLiteHalbeSessionNie`, `TestSQLiteSchliessenScheitert` |
| die Datei enthält zu jedem Zeitpunkt nur vollständige Sessions, je Session eine Transaktion | `Commit` statt `Rollback` (Z8); eine Transaktion für alle Sessions (W35); nach einem Fehler weiter (W22) | `TestSQLiteHalbeSessionNie`, `TestSQLiteTransaktionJeSession` |
| der nächste Schreibvorgang ergänzt, was fehlt | `angelegt` vor dem Erfolg gesetzt (W25); nach dem Anlegen zurückgesetzt (W37) | `TestSQLiteErsterSchreibvorgangScheitert`, `TestSQLiteErgaenzt` |
| die Datei wird vor jedem Ergänzen geprüft: fehlt sie, ist das `PGR-E3001` und sie wird nicht angelegt | `mode=rw` entfernt (W4) | `TestSQLiteDateiGeprueft` |
| … trägt `session` nicht genau die Nummern 1 bis zur zuletzt ergänzten | Prüfung entfernt (W16); Prüfung nur bei neuen Sessions (W17); nur die Zahl (W19); nur die Nummern (W18) | `TestSQLiteDateiGeprueft` (acht Fälle) |
| Schließen nach der Transaktion: `PGR-E3001`, nicht doppelt eingetragen | Fehler verschluckt (W15); Zähler nach dem Schließen zurückgenommen (W21) | `TestSQLiteSchliessenScheitert` |
| ein Abbruchsignal unterbricht keine Transaktion | `ctx` statt `context.WithoutCancel(ctx)` (W1) | `TestSQLiteAbbruchUnterbrichtNicht` |
| der Zielpfad gilt wörtlich (`?`, `#`, `%`, Leerzeichen) | `%`, `?` oder `#` nicht maskiert (W11, W12, W13) | `TestSQLiteSonderzeichenImPfad` |
| jeder Wert ist ein Parameter einer Anweisung | SQL-Text, Startup oder Feld in den Anweisungstext gesetzt (Z16, Z17, Z18) | `TestSQLiteWerteAlsParameter` |
| Rollback-Journal, `foreign_keys`, `synchronous=FULL`, Sperrfrist 5000 ms | `WAL` (W6); `NORMAL` (W7); `OFF` (W8) | `TestSQLiteEinstellungen`, `TestSQLiteFremdschluessel` |
| jede Transaktion beginnt mit `BEGIN IMMEDIATE` | `_txlock` entfernt (W9) | `TestSQLiteTransaktionImmediate` |
| die Sperrfrist beträgt fünf Sekunden | 60 s (W5b); 1 s (W5c) | `TestSQLiteSperrfrist`, `TestSQLiteEinstellungen` |
| Name, Anlegen und Namenswahl der temporären Datei | Präfix `.` entfernt (T2); elf Versuche (T1); `O_TRUNC` statt `O_EXCL` (T3); jeder Fehler zieht einen neuen Namen (T5) | `TestSQLiteNameDerTempDatei`, `TestSQLiteAnlegenScheitert`, `TestSQLiteBelegterName` |
| Rechte der temporären Datei (0666 nach umask; ersetzte Datei behält ihre) | `0o600` (T4); Rechte nicht gesetzt (W24) | `TestSQLiteRechte` |
| Fehlschlag beim ersten Schreibvorgang: `PGR-E3001`, temporäre Datei entfernt, Entfernen als Ursache | Entfernen nicht aufgerufen (W26, W27); Fehler des Entfernens nicht verbunden (W28); Schließen ohne Fehler (W29) | `TestSQLiteErsterSchreibvorgangScheitert`, `TestSQLiteEntfernenScheitert` |
| `meta`: Formatkennung, Version, `empty_sessions` nur mit leeren Sessions | `empty_sessions` immer (W30b), nie (W30), mit anderem Text (W30e); Version `1.0` (W30c); andere Kennung (W30d) | `TestSQLiteInhalt`, `TestSQLiteEmptySessions` |
| Inhalt der Tabellen nach der Liste: `group_no` ab 1, `position` über beide Richtungen, Client vor Server in der Gruppe, `group_no` und `sql` nur bei der jeweiligen Art, `offset_ms` `NULL` oder Zahl, `direction`, Spaltenfolge | 0-basiert (Z1); `position` je Richtung (Z2); Server vor Client (Z3); `group_no` bei `query` (Z4); `sql` bei `query` fehlt (Z5), bei `extended` gesetzt (Z5b), leerer Text als `NULL` (Z5c); 0 als `NULL` (Z6); jede Zahl 1 (Z6b); `server_parameters` aus `startup` (Z7b); Spalten vertauscht (Z10, Z11); Richtung `client` (Z12); `Server` (Z15) | `TestSQLiteInhalt`, `TestSQLiteLeereWerte` |
| `fields` ohne `type`, mit den Schlüsseln und Werten der YAML-Darstellung, nie `NULL` | `type` bleibt (Z13, Z14, J6); Zahlen als String (J7); `{}` als `null` (J11) | `TestSQLiteInhalt`, `TestSQLiteFieldsGleichYAML`, `TestSQLiteLeereWerte` |
| die Umwandlung kennt Abbildung, Liste, String, Zahl, Wahrheitswert, `null` und maskiert `<`, `>`, `&` nicht | `null` als String (J8); Wahrheitswert immer `true` (J9); HTML maskiert (J10); `!!null`-Schlüssel nicht erkannt (J1) | `TestSQLiteJSONTypen` |
| ein String ohne gültiges UTF-8 steht als `base64:` mit seinen Bytes, als Wert und als Schlüssel | Schlüssel nicht umgeschrieben (J4); Bytes ohne Präfix (J5) | `TestSQLiteStringOhneUTF8`, `TestSQLiteUmschreibungKehrtZurueck`, `TestSQLiteFluchtzeichen` |
| ein gültiger String mit führendem `\` oder `base64:` steht mit einem weiteren `\` | nur `base64:` maskiert (J2); nur `\` maskiert (J3); nichts maskiert (J13) | `TestSQLiteFluchtzeichen`, `TestSQLiteUmschreibungKehrtZurueck` |
| `interaction.sql` steht unverändert | Text durch die Umschreibung geschickt, hier um ein Zeichen verändert (Z9) | `TestSQLiteSQLUnveraendert`, `TestSQLiteInhalt` |
| `Load` ist ein Platzhalter mit `PGR-E1000`, ohne Pfad | gibt `nil` zurück (W2); Code `PGR-E3003` (W3); Pfad in der Meldung (W3b) | `TestSQLiteLoadPlatzhalter` |
| Aufruf mit weniger Sessions, als die Datei trägt, ist `PGR-E1000` und ändert nichts | Prüfung entfernt (W33) | `TestSQLiteWenigerSessionsAlsDieDatei` |
| das Schema: Tabellen, Spalten, Typen, `NOT NULL`, Primär- und Fremdschlüssel nach der Liste | Spalte `server_parameters` entfernt (S1); Primärschlüssel von `message` verkürzt (S2); `ON DELETE CASCADE` entfernt (S3); Spalte `bytes` (S4); `fields` nullable (S5); `offset_ms` `NOT NULL` (S6); `meta` ohne Primärschlüssel (S7); `type` als Zahl (S8); Fremdschlüssel von `interaction` entfernt (S9) | `TestSQLiteSchema` |
| `modernc.org/sqlite` im Architektur-Gate nur im Recording-Adapter | Regel entfernt (A1); auch im CLI-Adapter zugelassen (A2); nur im PGWire-Adapter zugelassen (A3); auch im Postgres-Adapter zugelassen (A4) | `make a-check-negativ`: A1 meldet die Fälle `model-sqlite`, `postgres-sqlite`, `pgwire-sqlite`, `cli-sqlite`; A2 `cli-sqlite`; A3 `recording-sqlite` (und die Fälle, deren Kopie nun selbst rot ist); A4 `postgres-sqlite` |
| die Erlaubnisliste des Lint-Gates führt die Bibliothek | Eintrag aus `.golangci.yml` entfernt | `make lint`: `import of package modernc.org/sqlite is blocked` an `sqlite.go:18` (gomodguard_v2) |
| `record --format sqlite` ergänzt nach jeder beendeten Session | `CloseSession` schreibt nicht mehr (E3) | `TestE2ERecordFormatSQLite` |

Summe: 110 Mutanten der Unit-Tests (einschließlich `F1b`, `Z16` bis `Z18`, `T1` bis `T5`, `U1`), fünf am Binary
(`E1` bis `E5`), vier der Gegenprobe (`A1` bis `A4`) und einer des Lint-Gates; alle rot bis auf drei äquivalente
(nächster Absatz). Die Läufe sind Lauf-Beleg und nicht im Repo; die Tabelle trägt Mutation und Test.

**Nachtrag zu den Aufträgen des Architect (`75874ba`).**

| Zusage | Mutation | roter Test |
|---|---|---|
| doppelte Session-Nummer ist `PGR-E3001`, die Datei bleibt auf dem Stand vor dem Aufruf | Primärschlüssel von `session` entfernt (N1b) | `TestSQLiteDoppelteSessionNummer` |
| eine Lücke in den Nummern ist spätestens beim nächsten `Write` `PGR-E3001`, Datei unberührt | Prüfung vor dem Ergänzen entfernt (N2); Prüfung der einzelnen Nummer entfernt (N2b) | `TestSQLiteLueckeInDenNummern` |
| eine fremde `.<Name>.<Hex>.tmp` und `.tmp-journal` bleiben nach `Prepare` und `Write` unverändert | `Write` entfernt das Journal einer temporären Datei (N3) | `TestSQLiteFremdesTempJournal` |

Der Code blieb unverändert, die Prüfung griff. Die vier Mutanten liefen wie oben in frischen Kopien. Namensabgleich:
die drei Tests existieren (`grep -n "^func TestSQLite…("` in `sqlite_schreiben_test.go`).

**Grüne Mutanten, eingeordnet.** Drei bleiben grün und ändern das Verhalten nicht:
(B2) Der Bootstrap ignoriert den Fehler von `recording.New`: Die Kommandozeile lehnt vorher jeden Wert außer `yaml`
und `sqlite` ab (`TestParseRecordFormat`, `TestRunRecordFormatUngueltig`), der Zweig ist durch kein Programm
erreichbar. (W10) `filepath.Abs` entfällt: SQLite löst einen relativen Pfad gegen das Arbeitsverzeichnis auf,
`TestSQLiteSonderzeichenImPfad` fährt den relativen Fall unter Linux; die Form `/C:/…` für Windows ist ungeprüft
(Grenze Plattform, `BEO-REPO/verhalten-nur-unter-linux-geprueft`). (W14) Das Leerzeichen wird nicht als `%20`
maskiert: SQLite nimmt es in einer `file:`-URI roh an. Im ersten Lauf waren außerdem grün und sind geändert:
ein totes `nichtNil` für leere Abbildungen (die Bibliothek schreibt `nil` als `{}`), ein toter Zweig für ungültiges
UTF-8 in `maskiere` (die Bibliothek schreibt solche Strings als `!!binary`) und das Entfernen von Leerraum vor
`DecodeString` (das Zeilenumbrüche überspringt) — alle drei gelöscht; `null`, Wahrheitswert und die
HTML-Maskierung ohne Test — jetzt `TestSQLiteJSONTypen`; der Mutant E5 am Binary — der Test prüft die Meldung der
Kommandozeile (`erlaubt sind yaml und sqlite`), nicht nur das Wort `format`, das auch `Aufzeichnungsformat` der
Fabrik enthält; ein Mutant `W25`, der den Fehler nicht baute, wurde neu gebaut; die Meldung des Journals
(`vorhanden`, P3) prüft jetzt der Test.

Nach dem ersten Lauf kam `U1` hinzu (`PRAGMA user_version = 1` im Schema geschrieben): rot an `TestSQLiteInhalt`, der
`PRAGMA user_version` liest und `0` verlangt.

**Läufe.** Alle im Docker-Netz oder im Image, nichts auf dem Host außer `make`, `docker` und `bash`.

| Lauf | Stand | Ergebnis |
|---|---|---|
| `make lint` | vor `f052e63`: 10 Befunde `errcheck` in den neuen Tests, behoben; im ersten `make gates` vor `9ebab6c`: 1 Befund `gochecknoglobals` im Test über das Binary, behoben; danach | 0 Befunde |
| `make a-check-negativ` | vor `07d34c4` | grün, siebzehn Fälle |
| `make test-integration` | vor `fc59441` und im zweiten `make gates` | grün; legt Netz, Volume und Container unter dem Namen `pgr-it-<pid>` an und entfernt sie danach |
| `make docs-check`, `make kopf-check`, `make abdeckung-check` | nach dem Nachziehen von §3 und der Abdeckung | `d-check: 539 Datei(en) geprüft, 0 Befund(e)`; die beiden anderen Exit 0 |
| `make schema-validate` | nach der Änderung von `schema.yaml` | `Validation passed` |
| `make schema-generate` gegen `schema.sql` | vor `f052e63` | gleich bis auf die Kopfzeilen mit dem Zeitstempel |
| `make gates` | auf `9ebab6c` mit den Abdeckungstabellen als Änderung | Exit 0, 5 min 26 s; `lint-gegenprobe: gruen`, `a-check-negativ: gruen — alle siebzehn Faelle`, `baseline-verify: v6.18.0 OK` |
| Build des Binaries für sechs Plattformen | auf `9ebab6c` | `linux`, `darwin`, `windows` je `amd64` und `arm64`: `docker build --target build --build-arg TARGETOS=… --build-arg TARGETARCH=…` endet jeweils mit Exit 0; die Images danach per Name entfernt |

Der Lauf für sechs Plattformen belegt das Bauen der reinen Go-Bibliothek, nicht ihr Verhalten; ausgeführt ist nur
`linux/amd64` (Risiko in §6: *erst durch den Build im Slice belegt* — der Build ist belegt, das Verhalten auf den
anderen Plattformen bleibt ungeprüft, `BEO-REPO/verhalten-nur-unter-linux-geprueft`).

**Namensabgleich.** Ein Skript zog aus §3, §6 und diesem Abschnitt jeden Namen `Test…` und jeden Dateinamen in
Backticks und suchte `^func <Name>(` in `internal/**/*_test.go` und `test/integration/*_test.go` sowie die Dateien
im Baum: 47 Testnamen und 38 Dateien gefunden, nichts fehlt (`interaction.sql` ist eine Spalte, keine Datei).

**Entscheidungen des Architect vom 2026-10-11 zu den Hinweisen des Implementers** (alle vier in `LH-FA-22.a` *Schreiben von `sqlite`* festgehalten, Tabelle in §6):
1. *Journal vorhanden.* Gewollt: `Stat` wie beim Zielpfad (`LH-FA-07.a` *Vorhanden*, *Nicht prüfbar*). Code und `TestSQLitePrepareJournalReihenfolge` bleiben.
2. *Weniger Sessions als zuletzt ergänzt.* `PGR-E1000` (Fehler des Aufrufers, Datei unberührt) ist richtig. Code und `TestSQLiteWenigerSessionsAlsDieDatei` bleiben.
3. *Nummer der Session.* Eine Zusage: doppelte Nummer `PGR-E3001`, Lücke `PGR-E3001` spätestens beim nächsten Schreibvorgang. Auftrag an den Implementer in §3 (zwei Tests mit Mutation), der Code bleibt, wenn die Tests grün sind.
4. *Journal der temporären Datei.* Eine übrig gebliebene Datei im Sinn von `LH-FA-07.a`; `record` sucht und entfernt es nicht. Der Code stimmt; Auftrag: ein Test (§3). Akzeptiertes Negativ: Bleibt das Journal im selben Lauf nach einem Fehlschlag liegen, weil die temporäre Datei entfernt wird, räumt `record` es nicht auf; ein Name mit demselben Zufallsteil kommt nicht wieder vor.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand 2026-10-11): `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/spec-randform-erst-im-review-entschieden` (beide verkörpert, `AGENTS.md` §3.10 und §3.12) treffen den neuen Vertrag des Schreibens; §3 fordert die Mutation je Zusage, §6 nennt die Randformen. `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) trifft die Grenzen in §6. Keiner erreicht mit diesem Plan 3× neu.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-sqlite-format-doku`; die Dokumentation zählt hier nicht mit. Nachgezählt beim Schnitt vom 2026-10-11 (Abgabe des Lesens an `slice-v1-abschluss-sqlite-format-lesen`): Liefer-Punkte: 3 (Schreiben mit Fabrik und Option, Schema mit Ersatztest, Architektur-Gate mit Erlaubnisliste). Schichten: zwei (Recording-Adapter, CLI-Adapter); der Bootstrap reicht nur `o.Format` an die Fabrik (Verdrahtung, keine Logik), der Kommentar am Port ist keine Schicht, `schema.yaml`, die Regel in `.a-check.yml` und das Lint-Profil sind Nachzug am Rand. Geschätzt rund 1800 Zeilen mit Tests. Nachgezählt bei der Umsetzung (Implementer): Liefer-Punkte 3, Schichten zwei (Recording-Adapter, CLI-Adapter); der Bootstrap tauscht `recording.YAML{}` gegen `recording.New(o.Format)` und gibt den Fehler der Fabrik zurück wie den von `NewRecordService` (Verdrahtung, keine Logik, `internal/bootstrap/format_test.go` mit ihr), die Kommentare an Port und `RecordService` sind keine Schicht; gemessen 2794 Zeilen hinzu, davon 676 Go-Code und 1975 Go-Tests (§7, *Größe*); die drei Tests des Nachtrags (§7) kommen mit rund 70 Zeilen Go-Tests hinzu, ohne Code und ohne neue Schicht. Gesichtet beim Schnitt: `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) und `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` (2×) treffen nicht, keine Hälfte liegt an der Grenze und der Slice läuft über `next/` zurück.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
