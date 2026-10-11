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

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand 2026-10-11): `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/spec-randform-erst-im-review-entschieden` (beide verkörpert, `AGENTS.md` §3.10 und §3.12) treffen den neuen Vertrag des Schreibens; §3 fordert die Mutation je Zusage, §6 nennt die Randformen. `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) trifft die Grenzen in §6. Keiner erreicht mit diesem Plan 3× neu.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-sqlite-format-doku`; die Dokumentation zählt hier nicht mit. Nachgezählt beim Schnitt vom 2026-10-11 (Abgabe des Lesens an `slice-v1-abschluss-sqlite-format-lesen`): Liefer-Punkte: 3 (Schreiben mit Fabrik und Option, Schema mit Ersatztest, Architektur-Gate mit Erlaubnisliste). Schichten: zwei (Recording-Adapter, CLI-Adapter); der Bootstrap reicht nur `o.Format` an die Fabrik (Verdrahtung, keine Logik), der Kommentar am Port ist keine Schicht, `schema.yaml`, die Regel in `.a-check.yml` und das Lint-Profil sind Nachzug am Rand. Geschätzt rund 1800 Zeilen mit Tests. Gesichtet beim Schnitt: `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) und `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` (2×) treffen nicht, keine Hälfte liegt an der Grenze und der Slice läuft über `next/` zurück.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
