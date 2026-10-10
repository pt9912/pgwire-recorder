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
- [x] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-einspielen-extended` liefert, im Ist-Zustand des gebauten Binaries: Aus `slice-doku-ist-stand` (dort §6, *Teilweise geliefert*): Der Slice ersetzt im Handbuch die Grenz-Sätze zur Extended-Interaktion bei `play` — §4 *Eine Aufzeichnung in eine Datenbank einspielen* (Voraussetzung „mit einfachen Anfragen“ und der Hinweis, dass eine Folge des erweiterten Protokolls mit `PGR-E6001` endet) und §7 Zeile `PGR-E6001` (Ursache „die Aufzeichnung für `play` enthält eine Folge des erweiterten Protokolls“) — und im README die Sätze „spielt die einfachen Anfragen … ein“ und „Eine Aufzeichnung mit vorbereiteten Anweisungen lehnt `play` ab“. Dazu (Nachtrag der Closure von `slice-v1-abschluss-einspielen-extended`, `AGENTS.md` §3.13): im Handbuch §4 *Hinweise* die Sätze zu `--continue-on-error` und `--allow-recorded-errors` (bei Extended wartet `play` nach einer Fehlerantwort nur noch auf deren `ReadyForQuery`, Exit-Code 4, die Datenbank führt den Rest der Interaktion bis zum `Sync` nicht aus (Prüfweg und Grenze in §6); ein Fehler gilt als erwartet, wenn die Aufzeichnung in irgendeiner Gruppe der Interaktion eine Fehlerantwort trägt, `FATAL` bleibt `PGR-E4003`) und zum Abbruchsignal (bei Extended endet `play` nach dem `ReadyForQuery` der laufenden Interaktion, nicht nach der laufenden Anfrage). Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [x] `README.md` nennt, was `slice-v1-abschluss-einspielen-extended` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt (auch der Untertitel in Zeile 3, der `play` auf „einfache Anfragen“ beschränkt; Nachtrag der Closure von `slice-v1-abschluss-einspielen-extended`), Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
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

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** weiter offen (Register, `BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×; kein Auftreten, keine neue Datei: der Diff trägt keinen neuen Codeblock, die Sätze beschreiben das Verhalten von `play` und der Datenbank, nicht eine Eigenschaft des Dateisystems; alle Proben, auch die der Verifikation, liefen nur unter Linux im Container).
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** entfallen (trat nicht ein: die Stellen aus DoD-Punkt 1 und 2 sind ersetzt, die Suche nach „einfache Anfragen“, „Folge des erweiterten Protokolls“ als Ablehnungsgrund und „lehnt“ findet nichts Verbliebenes, Verifikation Punkt 1 und 2; der Eintrag trägt den Vermerk „Teil geliefert“ ohne Zähler).

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
| 6 | Der Satz zu `FATAL` im Hinweis zu `--allow-recorded-errors` (`PGR-E4003`, Abbruch auch mit der Option) gilt auch für eine Extended-Interaktion | Nacharbeit zu F-587, `make build` am Stand `bec21ad`, dasselbe Image und Server (PostgreSQL 17, `trust`). Eine Aufzeichnung mit `record` und `psql \bind` von `select pg_terminate_backend(pg_backend_pid())` scheitert (`record` endet die Sitzung mit `PGR-E6001`, die Aufzeichnung trägt die Interaktion nicht); deshalb Aufzeichnung von Hand aus der von `select 1 \bind \g`, SQL durch den `pg_terminate_backend`-Aufruf ersetzt. 6a: Aufzeichnung mit `error_response` (`S: FATAL`, `57P01`) in der Interaktion; `play` ohne Option, mit `--continue-on-error`, mit `--allow-recorded-errors`, mit beiden. 6b: Aufzeichnung ohne `error_response`, mit `--allow-recorded-errors` | 6a: je Lauf Exit 4, ein `PGR-E4003` („Session 1, Interaktion 1“, 57P01), kein `PGR-E4004`. 6b: Exit 4, `PGR-E4003`. Der Satz gilt damit für Extended, auch mit `--continue-on-error`, und unabhängig davon, ob die Aufzeichnung einen Fehler trägt; das Handbuch bleibt unverändert, nennt keine Serverversion. Container, Netz und Scratch-Dateien (`impl2-extdoku-`) entfernt |

**Zur Kenntnis (Review, Nacharbeit).** F-588 (Begriff „Interaktion“ in der Meldung gegen „Nummer der Anfrage“ im Handbuch) und F-589 (README-Wort „Slices“ in der Zeile zu `make doc-trace`, kein Verstoß) gehen an den Planner; F-590 (Bestätigung) geht an den Verifier. Die Entscheidung liegt nicht beim Implementer.

**Neue und entfernte Aussagen.** Handbuch §4 *Voraussetzung* („mit einfachen Anfragen“) entfernt (Probe 1). *Hinweise*: neu der Satz zu Fehler und `--continue-on-error` bei Extended (2a, 2b), der Zusatz zu `--allow-recorded-errors` (3a bis 3d), der Satz zum Abbruchsignal bei Extended (4a bis 4c); der Satz „Enthält die Aufzeichnung eine Folge des erweiterten Protokolls … spielt nichts ein“ ersetzt durch den Wortlaut aus §6 (5a, 5b); er war auch falsch („spielt nichts ein“: Folgen vor dem Copy laufen, 5a). §7 Zeile `PGR-E6001`: Ursache „Aufzeichnung enthält eine Folge des erweiterten Protokolls“ ersetzt durch „die Datenbank antwortet bei `play` mit einem COPY-Datenstrom“. README: Untertitel „einfachen“ entfernt, Satz zu `play` nennt einfache wie vorbereitete Anweisungen, der Satz „lehnt `play` ab“ ersetzt durch den Copy-Satz.

**Ungenannt geblieben (entschieden in §6):** Warten je Antwort, Gegendruck, `CopyBothResponse`, nicht lesbare Serverantwort, Aufzeichnung endet mitten in einer Interaktion; Passwort (`PGR-E4005`) und `sslmode=require` (`PGR-E2004`) unverändert; keine Serverversion genannt. `FATAL` bei Extended: der Satz bleibt unverändert; die E2E-Tests tragen `FATAL` nur für einfache Anfragen (`TestE2EPlayFortsetzungAbbruch`, `q(...)`), den Beleg für Extended trägt die Probe 6 (Nacharbeit zu F-587).

**Optionen, Variablen, Codes** (Befehlsfolge aus `slice-doku-ist-stand` §7, erneut gefahren): bei den Optionen nur `--h`, `--help`, `--version` (Nicht-Optionen im Handbuch) und `--rm`, `--name` (`docker run`); Variablen und Codes ohne Ausgabe, grün. `grep -n -i -E "spec/|lastenheft|slice|welle|review|LH-|SPEC-|noch nicht|kommt|geplant|künftig"` über das Handbuch findet nur Bestand ohne Bezug (`noch nicht` in §3 „Zieldatei existiert noch nicht“, „kommt“ in der Ruhefrist von `record`); `make docs-check`: 0 Befunde. Das Handbuch trägt keinen neuen Codeblock; das Beispiel in §4 (`play --upstream … --input …`) hat die Form der Probe 1.

**Mutation (§3.10)** entfällt (kein neuer Vertrag); die Probe-Matrix ersetzt sie.

**Gate-Lauf.** `make gates` am Stand `2cb0519` auf sauberem Baum: Exit 0 (docs-check 0 Befunde, test, test-integration, a-check, lint und ihre Gegenproben grün).

### Closure

**Review und Verifikation.** Review `bec21ad` (F-587 LOW, F-588 bis F-590 INFO; keine HIGH, keine MEDIUM), F-587 behoben in `2cb0519` und `8892421` (Probe 6 in §7, Handbuch unverändert). Verifikation `bb0dfd3`: DoD-Punkte 1 bis 3 bestätigt, V-159 bis V-161 INFO.

**Entscheidungen zu den Übergaben des Planners.**

- **F-588 und V-160** (Begriff „Interaktion“ in der Meldung, „Nummer der Anfrage“ im Handbuch): Ins Register, kein Nehmer. Das Handbuch ist am Binary nicht falsch; eine Bereinigung wäre eine Wortwahl in einem Slice, dessen Gegenstand sie nicht ist. Der einzige Doku-Slice, der dieselben Abschnitte noch anfasst (`slice-v1-abschluss-einspielen-anmeldung-doku`), steht mit drei Liefer-Punkten an der Grenze und nähme sie nicht an (`AGENTS.md` §3.13). Welche Seite den Begriff ändert, Handbuch oder Meldung, entscheidet die Spezifikation (`LH-FA-20.a` *Interaktion*), nicht dieser Slice. `BEO-REPO/begriff-der-meldung-weicht-vom-handbuch-ab`, neu, 1×, offen.
- **V-161** (`FATAL` bei Extended ruhe nur auf einer Probe, kein Test): Die Feststellung stimmt nur halb. `TestPlayExtendedFatal` in `internal/hexagon/services/play_extended_test.go` trägt die Zusage für eine Extended-Interaktion (`FATAL`, Fehlerantwort aufgezeichnet, beide Optionen: `PGR-E4003`, Exit 4, keine weitere Gruppe) und ist in `abdeckung-unit.md` deklariert; der Test lief in `slice-v1-abschluss-einspielen-extended` mit Mutationen rot. Offen ist nur ein Lauf gegen eine echte Datenbank, und den trägt Probe 6. Keine Sendung (kein Nehmer nimmt an: `slice-v1-abschluss-einspielen-anmeldung` läge mit einem weiteren Test nach der Zählung des Hexagons über zwei Schichten, siehe dessen §8, und der Punkt gehört nicht zu ihrem Gegenstand) und kein eigener Eintrag; der Fund steht im Beleg von `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. Die Verifikation bleibt als Lauf-Beleg unverändert.
- **F-589** (README „Slices“ in der Zeile zu `make doc-trace`): Bestand bleibt stehen. `AGENTS.md` §3.11 verbietet dem README Chronik und Zielstand, nicht das Wort; die Zeile erklärt die Spalten der Ausgabe eines Targets und sagt nichts über einen Stand.
- **V-159** (Umbruch zweier Absätze nach dem Ersetzen): Hinweis, kein Auftrag; der gerenderte Text ist richtig, kein Gate betroffen. Nachgezogen wird beim nächsten Slice, der diese Absätze ohnehin ändert.
- **F-590**: Bestätigung, keine Handlung.

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge):

| Risiko | Ausgang | Begründung |
|---|---|---|
| Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut | weiter offen | Kein Auftreten und keine neue Datei: Der Diff trägt keinen neuen Codeblock, die Sätze beschreiben `play` und die Datenbank. Alle Proben liefen nur unter Linux im Container; der Eintrag `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) führt die Klasse. |
| Der Slice liefert seinen Teil von Handbuch und README nicht vollständig | entfallen | Trat nicht ein (Verifikation Punkt 1 und 2). Der Eintrag `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` bleibt 1×, offen, und trägt den Vermerk „Teil geliefert“ ohne Zähler; neun Doku-Folge-Slices stehen noch aus. |

**Register** (Zähler = Dateien unter `evidence/`, gelesen am Stand dieser Closure):

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — neue Datei `evidence/slice-v1-abschluss-einspielen-extended-doku.md` (F-587): 32 → 33. Verkörpert (`AGENTS.md` §3.11, Schritt 17, 19 und 20 von `implement-slice`); Sensor geplant (`slice-harness-mutation`).
- `BEO-REPO/begriff-der-meldung-weicht-vom-handbuch-ab` — neu, `evidence/slice-v1-abschluss-einspielen-extended-doku.md` (F-588, V-160): 1×, offen.
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` — kein Auftreten; Vermerk „Teil geliefert“ in `state.md`, bleibt 1×, offen.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` — kein Auftreten (siehe Risiko); bleibt 1×, offen.
- Kein nicht verkörperter Eintrag erreicht mit diesem Slice 3×. Keine Datei für V-159 (Form, einmalig) und F-589 (Einordnung).

**Lerneintrag — geschärfte Regel.** Die Probe-Matrix hat den Doku-Slice getragen, ohne Mutation: Jede neue Aussage hatte eine Zeile, und was sich nicht belegen ließ, blieb ungenannt; das ist die Form aus `slice-doku-ist-stand` §6 und braucht keine neue Regel. Die Lücke lag anderswo: F-587 war ein Satz, der *nicht* geändert wurde. Der Architect hatte ihn in §6 als „bleibt, wie er ist“ entschieden, weil er Extended nicht nennt; er stand aber unmittelbar hinter dem neuen Satz über dieselbe Option und las sich damit für Extended mit. Die Matrix prüfte neue Sätze, nicht die Reichweite der Nachbarn. Geschärft ist Schritt 17 von `implement-slice`: Ein unverändert bleibender Satz neben einem neuen Satz über denselben Gegenstand bekommt eine Probe oder einen Test für dessen Fall oder wird enger gefasst. `liegt in` `.claude/commands/implement-slice.md` Schritt 17 · seit slice-v1-abschluss-einspielen-extended-doku. **Kein neuer Sensor:** Ein Sensor, der Handbuch gegen Binary prüft, ist in §1 ausgeschlossen; die Regel steht in dem Schritt, den der Implementer beim Schreiben der Sätze liest.

**Technik ohne Verkörperung.** Die Aufzeichnung für eine Probe, die `record` nicht erzeugen kann (`FATAL` beendet die Sitzung, `record` endet mit `PGR-E6001`), entsteht von Hand aus der Aufzeichnung einer ähnlichen Anfrage, SQL ersetzt (Probe 6). Ein Weg, keine Regel. Kein Eintrag.

**Paarungen.** (a) Anker: `state.md` von `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` nennt `.claude/commands/implement-slice.md` Schritt 17; die Datei trägt dort `seit slice-v1-abschluss-einspielen-extended-doku`. (b) Folge-Slice: Diese Closure weist keinem Slice etwas Neues zu (V-160 ins Register, V-161 keine Sendung). Die bestehenden Adressen für Passwort und `sslmode=require`, `slice-v1-abschluss-einspielen-anmeldung-doku` und `slice-v1-abschluss-einspielen-tls-doku`, liegen in `next/` und nennen ihren Geber `slice-doku-ist-stand` (`grep -c` je 6); ihre Sätze blieben in diesem Slice unangefasst. (c) Register: die vier genannten Kennungen existieren als Verzeichnis mit nicht leerem `evidence/`.

**Nachzählen.** Eingetragen wurde nur der Stand der Schichtzählung in §8 von `slice-v1-abschluss-einspielen-anmeldung` (drei Schichten nach der Zählung des Hexagons, zwei nach der Teilung des Plans, ohne Schnitt; der Architect entscheidet vor dem Code). Keine Sendung, daher kein Nachzählen eines Nehmers.

**Stand der Gates.** Die Läufe stehen oben (`2cb0519`); den Lauf auf dem Stand der Closure nennt der Bericht des Planners.

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
