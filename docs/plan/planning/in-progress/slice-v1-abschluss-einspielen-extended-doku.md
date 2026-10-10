# Slice slice-v1-abschluss-einspielen-extended-doku: Handbuch und README zu: Extended-Interaktionen einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-einspielen-extended`)

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

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-einspielen-extended` liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-einspielen-extended`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-einspielen-extended` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Nachtrag der Closure von `slice-v1-abschluss-einspielen-extended` (2026-10-10):** `play --help` und die Optionen blieben unverändert; geliefert ist das Einspielen von Extended-Interaktionen mit Warten je Gruppe, Abbruch, Fortsetzung, erwartetem Fehler und Abbruchsignal innerhalb der Interaktion. Die Sätze dazu stehen in der DoD; Kennung des Gebers: `slice-v1-abschluss-einspielen-extended`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-einspielen-extended` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-einspielen-extended` liefert, im Ist-Zustand des gebauten Binaries: Aus `slice-doku-ist-stand` (dort §6, *Teilweise geliefert*): Der Slice ersetzt im Handbuch die Grenz-Sätze zur Extended-Interaktion bei `play` — §4 *Eine Aufzeichnung in eine Datenbank einspielen* (Voraussetzung „mit einfachen Anfragen“ und der Hinweis, dass eine Folge des erweiterten Protokolls mit `PGR-E6001` endet) und §7 Zeile `PGR-E6001` (Ursache „die Aufzeichnung für `play` enthält eine Folge des erweiterten Protokolls“) — und im README die Sätze „spielt die einfachen Anfragen … ein“ und „Eine Aufzeichnung mit vorbereiteten Anweisungen lehnt `play` ab“. Dazu (Nachtrag der Closure von `slice-v1-abschluss-einspielen-extended`, `AGENTS.md` §3.13): im Handbuch §4 *Hinweise* die Sätze zu `--continue-on-error` und `--allow-recorded-errors` (bei Extended wartet `play` nach einer Fehlerantwort nur noch auf deren `ReadyForQuery`, Exit-Code 4, die Datenbank führt den Rest der Interaktion bis zum `Sync` nicht aus (Prüfweg und Grenze in §6); ein Fehler gilt als erwartet, wenn die Aufzeichnung in irgendeiner Gruppe der Interaktion eine Fehlerantwort trägt, `FATAL` bleibt `PGR-E4003`) und zum Abbruchsignal (bei Extended endet `play` nach dem `ReadyForQuery` der laufenden Interaktion, nicht nach der laufenden Anfrage). Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [ ] `README.md` nennt, was `slice-v1-abschluss-einspielen-extended` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt (auch der Untertitel in Zeile 3, der `play` auf „einfache Anfragen“ beschränkt; Nachtrag der Closure von `slice-v1-abschluss-einspielen-extended`), Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
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
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-einspielen-extended` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt |
| `README.md` | update | Sätze, die `slice-v1-abschluss-einspielen-extended` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe grün.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen-extended` liegt in `done/` (Voraussetzung; Schritt 11 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Handbuch und README getrennt.
- `in-progress` → `open` (blockiert — Carveout?): Eine Aussage lässt sich nur mit einer Änderung am Binary wahr machen, und der Nutzer will sie nicht streichen; dann zuerst die Entscheidung des Nutzers, welcher Slice das Verhalten liefert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, die Proben und `grep`-Läufe in §7, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect prüft vor dem ersten Commit, ob `slice-v1-abschluss-einspielen-extended` eine weitere Randform liefert, die dort nicht steht.

**Randformen dieses Slice** — vom Architect vor dem ersten Commit entschieden (2026-10-10). **Ort ist dieser Abschnitt, nicht die Spezifikation:** Der Slice legt keinen Vertrag an; die Verhaltensregeln stehen schon in `LH-FA-20.a` (*Gruppen*, *Interaktion*, *Abbruchsignal*) und in §6 von `slice-v1-abschluss-einspielen`. Das Handbuch sagt nur, was am Binary eine Probe zeigt. Probe-Client ist `psql` aus dem gepinnten PostgreSQL-Image (`\bind` und `\parse` erzeugen das erweiterte Protokoll), Server das gepinnte Image; Aufzeichnungen, die `record` nicht erzeugen kann (Copy), schreibt die Probe als Datei.

- **Probe-Matrix** (jede Zeile ein Handbuch-Satz; ohne bestandene Probe entfällt der Satz):
  1. Aufzeichnung mit Extended-Interaktion (`record`, `psql \bind`), `play` gegen die Datenbank: Wirkung in der Tabelle, Exit 0. Ersetzt Voraussetzung „mit einfachen Anfragen“ in §4 und den README-Satz.
  2. Fehler in einer Extended-Interaktion (Verletzung eines Schlüssels in der ersten Gruppe, eine spätere Gruppe derselben Interaktion schreibt eine Zeile): `PGR-E4004`, Exit 4; mit `--continue-on-error` läuft `play` mit der nächsten Interaktion weiter, Exit 4 am Ende; die Zeile der späteren Gruppe der fehlerhaften Interaktion fehlt in der Tabelle. Der Handbuch-Satz sagt das Beobachtbare („die Datenbank führt den Rest der Interaktion bis zum `Sync` nicht aus; die nächste Interaktion läuft“), nicht, wie `play` sendet.
  3. `--allow-recorded-errors` bei Extended: Aufzeichnung trägt die Fehlerantwort in einer anderen Gruppe als die, in der der Fehler jetzt auftritt (handgeschrieben): Exit 0 ohne Fehlerzeile; ohne Fehlerantwort in der Aufzeichnung `PGR-E4004`. Satz: „irgendwo in der Interaktion“, nicht SQLSTATE, nicht Stelle.
  4. Abbruchsignal bei Extended: Interaktion mit `pg_sleep` in `Execute`, `SIGTERM` währenddessen: `play` endet nach dem Ergebnis der Interaktion, die Zeile aus der Interaktion steht in der Tabelle, Exit 0. Mit `--finish-session-on-interrupt` nach der Sitzung; zweites Signal bleibt wie beim einfachen Hinweis (nur gesagt, wenn Probe es für Extended zeigt, sonst Satz nicht um Extended erweitert).
  5. Copy: handgeschriebene Aufzeichnung mit `COPY … FROM STDIN` und mit `COPY … TO STDOUT` im erweiterten Protokoll: `PGR-E6001`, Exit 6, Meldung nennt Sitzung und Nummer. Ersetzt die Ursache der Zeile in §7 und den Hinweis in §4.
- **Ungenannt, weil nicht belegbar** (akzeptiertes Negativ, Grund je Punkt; die nächste Runde liest sie als entschieden): Warten je Antwort nach `Flush` und bis `ReadyForQuery` nach `Sync`, „übrige Gruppen ohne Warten“, Gegendruck (nur an der Ausführungszeit oder gar nicht sichtbar; ein Satz dazu prüfte kein Test des Handbuchs); `CopyBothResponse` und eine nicht lesbare Serverantwort (keine Eingabe gegen das gepinnte Image löst sie aus; Unit-Test, wie `PGR-E1000` in `slice-doku-ist-stand` §7); Warten ohne eigene Frist bei einer Aufzeichnung, die mitten in einer Interaktion endet (zeigt die Probe etwas Stabiles, geht der Fund als Zeile in §7 an den Planner, der Slice beschreibt es nicht).
- **`FATAL` bei Extended** — der Satz zu `FATAL` und `PGR-E4003` im Hinweis zu `--allow-recorded-errors` bleibt, wie er ist. Er nennt Extended nicht; erweitert wird er nur, wenn der Implementer eine Probe oder einen Test für Extended benennt (`test/integration/play_laufsteuerung_e2e_test.go` trägt `FATAL` und Copy; er prüft vor dem Schreiben, ob das Extended einschließt).
- **Serverversionen** — das Gate fährt PostgreSQL 17 (gepinntes Image). Das Handbuch nennt keine Version und kein „ab“; die Sätze beschreiben „die Datenbank“, die Proben laufen gegen das gepinnte Image (§7 nennt es).
- **Grenzen Anmeldung und TLS** — bleiben unverändert stehen: Passwort bei `play` (`PGR-E4005`) und `sslmode=require` (`PGR-E2004`) sind Bestand im Präsens und gehören `slice-v1-abschluss-einspielen-anmeldung-doku` und `slice-v1-abschluss-einspielen-tls-doku` (beide in `next/`, nehmen die Sendung per `slice-doku-ist-stand` an). Dieser Slice fasst die Sätze zu Passwort und `sslmode` in §1, §4, §5, §7 nicht an; eine neue Zuweisung entsteht nicht (`AGENTS.md` §3.13).
- **Wortlaut der Grenzen** — im Präsens, als Eigenschaft des Binaries, nur mit beobachtbarer Reaktion (Meldungscode, Exit-Code): „Antwortet die Datenbank mit einem COPY-Datenstrom (`COPY … FROM STDIN`, `COPY … TO STDOUT`), kann `play` ihn nicht verarbeiten und endet mit `PGR-E6001` (Exit-Code 6).“ Kein Satz sagt, dass etwas kommt; „Folge des erweiterten Protokolls“ als Ablehnungsgrund entfällt überall (§4, §7, README, Untertitel Zeile 3).
- **Mutation (§3.10)** — entfällt: kein neuer Vertrag, nur Beschreibung; die Prüfung ist die Probe-Matrix, rot, wenn ein Satz ohne bestandene Probe stehen bliebe.
- **Größe** — zwei Liefer-Punkte (Handbuch, README), eine Schicht (Dokumentation); die Rückführung aus §4 greift nicht.

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** offen bis Closure.
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** offen bis Closure.

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

**Probe-Aufbau.** `make build` am Stand `2e010b5` (Image `pgwire-recorder:dev`, `sha256:b8abba88772c…`). Proben in einem Scratch-Verzeichnis außerhalb des Repos; Netz `impl-extdoku-net` mit dem gepinnten PostgreSQL-Image aus `harness/mk/integration.mk` (`postgres:17-alpine@sha256:b0f9560a…`, Anmeldung `trust`), je Probe eine eigene Datenbank mit `t (n int PRIMARY KEY, s text)`. Probe 1 zeichnet mit `record` und `psql` (`\bind`, Image wie der Server) auf; Aufzeichnungen der Proben 2 bis 5 sind von Hand geschrieben (`psql` 17 kennt keine Pipeline, `record` kann Copy nicht erzeugen). Container, Netz und Scratch-Dateien sind danach entfernt.

| Nr. | Aussage (Handbuch/README) | Probe | Ergebnis |
|---|---|---|---|
| 1 | `play` spielt Extended-Interaktionen ein; „mit einfachen Anfragen“ entfällt (§4 Voraussetzung, README-Satz, Untertitel Zeile 3) | `record` + `psql`: zwei `INSERT … \bind … \g`; `play --database p1` | Exit 0, Zeilen `1:eins`, `2:zwei` |
| 2a | Fehler in einer Extended-Folge: Rest bis `Sync` nicht ausgeführt, Exit 4 | 1. Folge: `INSERT (1)` (Schlüssel vorhanden) mit Flush, danach `INSERT (2)` mit Sync; 2. Folge `INSERT (3)`; ohne Option | Exit 4, ein `PGR-E4004` („Interaktion 1“, 23505), Tabelle `1:vorhanden` |
| 2b | `--continue-on-error`: nächste Folge läuft, Exit 4 | dieselbe Aufzeichnung | Exit 4, ein `PGR-E4004`, Tabelle `1:vorhanden,3:naechste` (Zeile 2 fehlt) |
| 3a | `--allow-recorded-errors`: Fehlerantwort irgendwo in der aufgezeichneten Folge genügt | wie 2, Fehlerantwort (23505) nur in der Sync-Gruppe aufgezeichnet, Fehler tritt in der Flush-Gruppe auf | Exit 0, keine Zeile `error`, Tabelle `1:vorhanden,3:naechste` |
| 3b | ohne die Option wirkt die aufgezeichnete Fehlerantwort nicht | dieselbe Aufzeichnung ohne Option | Exit 4, `PGR-E4004` |
| 3c | ohne Fehlerantwort in der Aufzeichnung bleibt es ein Fehler, auch mit der Option | Aufzeichnung aus 2 mit `--allow-recorded-errors` | Exit 4, `PGR-E4004` |
| 3d | Gleichheit der SQLSTATE ist nicht verlangt | Aufzeichnung aus 3a mit `XX000` statt `23505` | Exit 0 |
| 4a | Signal: `play` endet nach der laufenden Folge | Folge `INSERT … FROM (SELECT pg_sleep(4))`, zweite Folge `INSERT (11)`, zweite Sitzung `INSERT (12)`; `SIGTERM` nach 1,5 s | Exit 0, Log `Abbruchsignal, play endet vorzeitig`, 2,1 s nach dem Signal beendet, Tabelle `10:warten` |
| 4b | mehrere Gruppen: endet nach dem `Sync` der Folge | Folge mit Flush-Gruppe `pg_sleep(4)` und Sync-Gruppe `INSERT (20)`, zweite Folge `INSERT (21)`; `SIGTERM` nach 1,5 s | Exit 0, Tabelle `10:warten,20:gruppe-b` |
| 4c | `--finish-session-on-interrupt`: nach der laufenden Sitzung | wie 4a mit der Option | Exit 0, Tabelle `10:warten,11:zweite` (kein `12`) |
| 4d | zweites Signal: Verbindung sofort geschlossen, kein Fehler, auch mit der Option | wie 4a und 4c, zweites `SIGTERM` 0,5 s nach dem ersten | Exit 0, beendet 0,02 s nach dem zweiten Signal, keine Zeile `error`; der bestehende Satz gilt damit auch für Extended |
| 5a | COPY-Datenstrom von der Datenbank: `PGR-E6001`, Exit 6 | Folge `COPY t FROM STDIN` nach einer Folge `INSERT (5)`, danach `INSERT (6)`; auch mit `--continue-on-error` | Exit 6, `PGR-E6001` („Session 1, Interaktion 2“, `CopyInResponse`), Tabelle `5:vorher` |
| 5b | wie 5a für `COPY … TO STDOUT` | Folge `COPY t TO STDOUT` | Exit 6, `PGR-E6001` (`CopyOutResponse`), Tabelle `5:vorher` |

**Neue und entfernte Aussagen.** Handbuch §4 *Voraussetzung* („mit einfachen Anfragen“) entfernt (Probe 1). *Hinweise*: neu der Satz zu Fehler und `--continue-on-error` bei Extended (2a, 2b), der Zusatz zu `--allow-recorded-errors` (3a bis 3d), der Satz zum Abbruchsignal bei Extended (4a bis 4c); der Satz „Enthält die Aufzeichnung eine Folge des erweiterten Protokolls … spielt nichts ein“ ersetzt durch den Wortlaut aus §6 (5a, 5b); er war auch falsch („spielt nichts ein“: Folgen vor dem Copy laufen, 5a). §7 Zeile `PGR-E6001`: Ursache „Aufzeichnung enthält eine Folge des erweiterten Protokolls“ ersetzt durch „die Datenbank antwortet bei `play` mit einem COPY-Datenstrom“. README: Untertitel „einfachen“ entfernt, Satz zu `play` nennt einfache wie vorbereitete Anweisungen, der Satz „lehnt `play` ab“ ersetzt durch den Copy-Satz.

**Ungenannt geblieben (entschieden in §6):** Warten je Antwort, Gegendruck, `CopyBothResponse`, nicht lesbare Serverantwort, Aufzeichnung endet mitten in einer Interaktion; Passwort (`PGR-E4005`) und `sslmode=require` (`PGR-E2004`) unverändert; keine Serverversion genannt. `FATAL` bei Extended: der Satz bleibt unverändert, `test/integration/play_laufsteuerung_e2e_test.go` und `play_extended_e2e_test.go` tragen `FATAL` nur für einfache Anfragen (`TestE2EPlayFortsetzungAbbruch`, `q(...)`), eine Extended-Probe ist nicht gefahren.

**Optionen, Variablen, Codes** (Befehlsfolge aus `slice-doku-ist-stand` §7, erneut gefahren): bei den Optionen nur `--h`, `--help`, `--version` (Nicht-Optionen im Handbuch) und `--rm`, `--name` (`docker run`); Variablen und Codes ohne Ausgabe, grün. `grep -n -i -E "spec/|lastenheft|slice|welle|review|LH-|SPEC-|noch nicht|kommt|geplant|künftig"` über das Handbuch findet nur Bestand ohne Bezug (`noch nicht` in §3 „Zieldatei existiert noch nicht“, „kommt“ in der Ruhefrist von `record`); `make docs-check`: 0 Befunde. Das Handbuch trägt keinen neuen Codeblock; das Beispiel in §4 (`play --upstream … --input …`) hat die Form der Probe 1.

**Mutation (§3.10)** entfällt (kein neuer Vertrag); die Probe-Matrix ersetzt sie.

**Gate-Lauf.** `make gates` am Stand `880522e` auf sauberem Baum: Exit 0 (docs-check 0 Befunde, test, test-integration, a-check, lint und ihre Gegenproben grün).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — genau der Gegenstand dieses Slice: Das Handbuch sagt nur zu, was das Binary zeigt.
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-einspielen-extended`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code. Nachgezählt beim Eintragen der Sendung aus der Closure von `slice-v1-abschluss-einspielen-extended` (2026-10-10): weiter zwei Liefer-Punkte (Handbuch, README) und eine Schicht; die Sätze zu den Optionen und zum Abbruchsignal liegen im Handbuch-Punkt, der Untertitel im README-Punkt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
