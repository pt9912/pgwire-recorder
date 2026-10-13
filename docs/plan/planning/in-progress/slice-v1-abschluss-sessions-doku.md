# Slice slice-v1-abschluss-sessions-doku: Handbuch und README zu: Mehrere Sessions und Verbindungsfehler

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-sessions`)

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

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-sessions` liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-sessions`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-sessions` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Übernommen aus `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, offene Frage zu V-3; Planner, 2026-10-12, `AGENTS.md` §3.13): Das Handbuch nennt `--session-assignment` und `--record-empty-sessions` erst mit `slice-v1-abschluss-sessions` und nennt dann zu jeder Option ihre Umgebungsvariable mit den Kommandos, bei denen sie gilt; die Variable beim anderen Kommando ist `PGR-E2001` (Satz zur Meldung `PGR-E2001` im Handbuch, §5 *Einstellungen*). Das Handbuch beschreibt die Variable, nachdem das gebaute Binary sie liest.

**Übernommen aus `slice-v1-abschluss-sessions`** (Architect, 2026-10-12, `AGENTS.md` §3.11, §3.13; der Wortlaut ist der Ist-Zustand nach der Lieferung jenes Slice und wird am gebauten Binary geprüft, nicht an diesem Absatz): Das Handbuch sagt in der Sprache des Nutzers, ohne Verweis auf Spezifikation oder Slices:

- `record` zeichnet jede Verbindung mit mindestens einer abgeschlossenen Interaktion als eigene Session auf, mehrere Verbindungen gleichzeitig; die `id` zählt ab 1 in der Reihenfolge, in der die Sessions enden. Mit `--record-empty-sessions` (Umgebungsvariable `PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS`, Schlüssel `record_empty_sessions` im Abschnitt `record:`) wird auch eine Verbindung ohne Anfrage eine Session ohne Interaktionen, und das Recording trägt `empty_sessions: true` (YAML) beziehungsweise die Zeile `empty_sessions` in `meta` (SQLite), auch wenn keine Session entsteht. Keine Session ist eine Verbindung, deren Aufbau zum Server scheitert, und eine mit einer nicht unterstützten Interaktion.
- `replay --session-assignment first-request` (Default; Umgebungsvariable `PGWIRE_RECORDER_SESSION_ASSIGNMENT`, Schlüssel `session_assignment` im Abschnitt `replay:`): Die n-te Verbindung, die eine Anfrage stellt, erhält die n-te Session mit Interaktion; Sessions ohne Interaktion werden übersprungen. `connection`: Die n-te Verbindung erhält die Session mit der `id` n, auch wenn sie nichts anfragt; das Recording muss mit `--record-empty-sessions` entstanden sein (`empty_sessions: true`), sonst endet `replay` beim Start mit `PGR-E2001`. Andere Werte, auch `Connection` und der leere, sind `PGR-E2001`.
- Grenzen, die das Binary zeigt: Mehr Verbindungen als Sessions: eine Anfrage ist `PGR-E5003`, eine Lebendprüfung (nur Leerraum und Kommentare) wird beantwortet. Eine Anfrage an eine Session ohne Interaktion ist `PGR-E5001`. Die Zuordnung ist nur zugesichert, wenn die Clients ihre Verbindungen nacheinander aufbauen und nutzen. Ein Verbindungsfehler beendet nur diese Verbindung; `record` und `replay` laufen weiter, und der Exit-Code ist der des ersten gemerkten Fehlers.
- Verbindungen, die keine Session sind (Randformen aus `slice-v1-abschluss-sessions`, am gebauten Binary geprüft, nicht an diesem Absatz): Eine Verbindung zählt ab der gelesenen Startnachricht; `SSLRequest`, `GSSENCRequest`, `CancelRequest` und ein Ende vor der Startnachricht erzeugen im Record keine Session, auch mit `--record-empty-sessions`, und zählen im Replay nicht als Verbindung (Probe: `TestE2ERecordKeineVerbindung`, `TestE2EReplayZaehlungDerVerbindungen`). Eine Startnachricht mit vollständigem Kopf, deren Rest fehlt, ist `PGR-E6001` mit Exit-Code 6 und keine Session; weniger als 8 Bytes endet still (Probe: `TestE2ERecordHalbeStartnachricht`). Das Zwangsende im Aufbau schreibt keine Session, auch mit der Option (Probe: `TestE2ERecordZwangsendeImAufbau`: Aufbau zum Server, der nie endet). Ob das Handbuch diese Fälle nennt, entscheidet der Architect vor dem ersten Commit; was es nennt, nennt es nur nach einer Probe am Binary.
- `PGR-E5003` und `PGR-E4006` stehen im Meldungskatalog des Handbuchs bereits; geprüft wird ihr Text gegen den Katalog, ihre Beschreibung bei Sessions ergänzt nur, was das Binary zeigt.

Der Satz zu `--format sqlite` mit `--record-empty-sessions` (Probe: schreiben, mit `replay` lesen) gehört dazu, wenn das Handbuch SQLite beschreibt (Abschnitt zum Format).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-sessions` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-sessions` liefert, im Ist-Zustand des gebauten Binaries (zu `--session-assignment` und `--record-empty-sessions` mit Umgebungsvariable und dem Verhalten beim anderen Kommando, aus `slice-v1-abschluss-umgebung-anderes-kommando`; Probe am gebauten Binary): Die Abschnitte, die Kommandos, Optionen, Meldungscodes, Beispiele und Grenzen dieses Slice betreffen, stehen im Handbuch, und die Grenz-Sätze, die er überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [ ] `README.md` nennt, was `slice-v1-abschluss-sessions` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
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
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-sessions` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt |
| `README.md` | update | Sätze, die `slice-v1-abschluss-sessions` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-sessions` liegt in `done/` (Voraussetzung; Schritt 26 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

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

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect prüft vor dem ersten Commit, ob `slice-v1-abschluss-sessions` eine weitere Randform liefert, die dort nicht steht.

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

**Stand.** Geändert sind `docs/user/benutzerhandbuch.md` (Version 0.6: neuer Abschnitt *Mehrere Verbindungen aufzeichnen und wiedergeben*; Ergebnis und Hinweise beim Aufzeichnen, Voraussetzung und Hinweise beim Wiedergeben, Hinweise beim Einspielen, Satz zu `--fail-on-unconsumed` unter *Herunterfahren mit Frist*, Absatz und Tabelle in §5, Exit-Code 5 und der Absatz darunter, Zeilen `PGR-E2001`, `PGR-E2004`, `PGR-E3004`, `PGR-E5001`, `PGR-E5002`, `PGR-E5003`, `PGR-W2001`, Glossar), `README.md` und dieser Plan (§7). Produkt-Code, Tests, Gates, Hilfetexte und `docs/user/abdeckung-*.md` sind unverändert (`git diff cd57943 -- . ':!docs' ':!README.md'` ist leer); `make abdeckung-check` endet mit Exit 0. Commit `1fe572a` (Handbuch, README).

**Probe-Aufbau.** `make build` am Stand `cd57943` (Image `pgwire-recorder:dev`, `sha256:deb48422487d…`; Produkt-Code seit dort unverändert). Proben in einem Scratch-Verzeichnis außerhalb des Repos; Docker-Netz `impl-sd-net` mit dem gepinnten PostgreSQL-Image aus `harness/mk/integration.mk` (`postgres:17-alpine@sha256:b0f9560a…`, 17.11, Anmeldung `trust`, Netzalias `postgres`, damit die Beispiele des Handbuchs wörtlich laufen; `log_connections=on` für Zeile 44); eine zweite Instanz mit Passwort (`pgpw`, SCRAM) für Zeile 31; ein stummer Upstream (`nc -l` im Image von `psql`, Alias `stumm`) für Zeile 29. `record` und `replay` laufen als Container (`docker run -d`, Alias `rec` beziehungsweise `rep`, Beenden mit `docker stop`, Exit-Code aus `docker wait`); die Clients sind `psql` aus dem Image der Datenbank, eine Verbindung ohne Anfrage ist `psql` mit leerer Standardeingabe, gleichzeitig offene Verbindungen sind `(echo …; sleep n; echo …) | psql` im Hintergrund, die Verbindungen ohne Startnachricht, mit Anfrage nach Verschlüsselung, mit Abbruchwunsch und mit abgebrochener Startnachricht sendet `printf … | nc` aus demselben Image. Die Beispiele des neuen Abschnitts sind aus dem Handbuch in Dateien extrahiert (je Codeblock eine Datei) und laufen wörtlich, `pgwire-recorder …` als `docker run … pgwire-recorder:dev …`. Container und Netz danach entfernt; die Dateien liegen im Scratch-Verzeichnis außerhalb des Repos. Aufzeichnungen: `f` (Session 1 `eins`, 2 ohne Anfrage, 3 `drei`, mit `--record-empty-sessions`, als YAML und als SQLite), `g` (Session 1 `q1`, `q2`; Session 2 `r1`, `r2`), `b` und `c` (leer, `select 1`, leer; YAML und SQLite), `d` (nur leere Sessions), `a` (ohne Option, ohne Kennzeichnung).

**Probe-Tabelle** (Handbuch = `docs/user/benutzerhandbuch.md`, *MV* = der neue Abschnitt; `Q` = der Satz im Handbuch, Probe = Aufruf, Ergebnis = beobachtet):

| Nr. | Aussage (Handbuch/README) | Probe | Ergebnis |
|---|---|---|---|
| 1 | MV *Vorgehen* 1, Beispiel wörtlich; *Ergebnis*: die ersten Zeilen der YAML-Datei | Beispieldatei als `record …`, drei Verbindungen (ohne Anfrage, `select 1`, ohne Anfrage), `docker stop` | Exit 0; `head -5` der Datei gleich dem Block des Handbuchs (`diff` leer); drei Sessions |
| 2 | MV 1: Variable `PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS=true` genügt | `record` ohne Option, Variable gesetzt, eine Verbindung ohne Anfrage | Exit 0; `empty_sessions: true`, `- id: 1` |
| 3 | MV 1: Schlüssel `record_empty_sessions: true` im Abschnitt `record:`, Beispieldatei | `config show --config ./bsp2.yaml`; `record --config ./bsp2.yaml`, eine Verbindung ohne Anfrage | `config show` Exit 0, Inhalt wie geschrieben; Aufzeichnung mit `empty_sessions: true`, `- id: 1` |
| 4 | MV *Aufzeichnen*: Option geht der Variablen vor, diese der Datei | Variable `true` mit `--record-empty-sessions=false`; Datei `true` mit `--record-empty-sessions=false`; je eine Verbindung ohne Anfrage | je `sessions: []`, keine Kennzeichnung |
| 5 | MV *Aufzeichnen*: Option ohne Wert, `=true`, `=false`; `=1`, `=True`, `=` ungültig (`PGR-E2001`) | je eine Verbindung ohne Anfrage; die drei ungültigen Formen | ohne Wert und `=true`: Kennzeichnung und `- id: 1`; `=false`: nichts davon; `=1`, `=True`, `=`: `PGR-E2001` (`invalid boolean value … erlaubt sind true und false`), Exit 2 |
| 6 | MV: Variable nimmt `true`, `false`; leer gilt als nicht gesetzt; andere Werte `PGR-E2001`; Schlüssel andere Werte `PGR-E2004` | Variable `vielleicht`, `1`, `True`; Variable leer; Schlüssel `yes`, `1` | `PGR-E2001` (`Umgebungsvariable PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS: erlaubt sind true und false`), Exit 2; leer: `record gestartet`; Schlüssel: `PGR-E2004` (`record.record_empty_sessions: erlaubt sind true und false`), Exit 2 |
| 7 | MV *Aufzeichnen*: `id` zählt ab 1 in der Reihenfolge der Enden, die später begonnene Verbindung kann `1` haben | Verbindung A beginnt zuerst und endet nach 4 s, Verbindung B beginnt 1 s später und endet nach 1 s | `- id: 1` mit den Anfragen von B, `- id: 2` mit denen von A |
| 8 | Handbuch *Aufzeichnen*, Hinweis 1; MV: mehrere gleichzeitig offene Verbindungen sind eigene Sitzungen | wie Zeile 7: beide Verbindungen gleichzeitig offen | zwei Sessions mit je den eigenen Anfragen |
| 9 | Ohne Option ist eine Verbindung ohne abgeschlossene Anfrage keine Sitzung | eine Verbindung ohne Anfrage, eine mit `select 1`, ohne Option | eine Session (`- id: 1`), keine Kennzeichnung |
| 10 | Mit Option: Sitzung ohne Interaktionen; YAML `interactions: []`, SQLite Zeile in `session` ohne Zeilen in `interaction` | leer, `select 1`, leer, mit Option, als YAML und als SQLite | YAML: `- id: 1` und `- id: 3` mit `interactions: []`, `- id: 2` mit einer Interaktion; SQLite: `session` 1, 2, 3, `interaction` nur `2|1|select 1 as a;` |
| 11 | MV: auch eine Verbindung, die beim Zwangsende ihre erste Anfrage noch ausführte (`PGR-E4006`) oder keine gestellt hatte | `select pg_sleep(20)`, `--shutdown-timeout 500ms`, `docker stop`, ohne und mit Option; eine offene Verbindung ohne Anfrage beim `docker stop`, ohne und mit Option | pg_sleep: ohne Option `sessions: []`, mit Option Session 1 geschrieben, beide `PGR-E4006`, Exit 4; ohne Anfrage: ohne Option `sessions: []`, mit Option `- id: 1` mit `interactions: []`, Exit 0 |
| 12 | MV: ohne Option legt das Ende einer Verbindung ohne Anfrage die Zieldatei nicht an, mit Option schon; Handbuch *Aufzeichnen*, *Ergebnis*: aktualisiert nach dem Ende jeder Sitzung | neue Zieldatei, eine Verbindung ohne Anfrage, nach 1 s `ls`; danach eine mit Anfrage; YAML und SQLite, ohne und mit Option | ohne Option: Zieldatei fehlt (YAML und SQLite), nach der Verbindung mit Anfrage eine Session; mit Option: Zieldatei da, eine Session (SQLite: `meta` mit `empty_sessions`, `session` 1, `interaction` 0) |
| 13 | MV: Kennzeichnung auch ohne Sitzung, fehlt ohne Option | `record` ohne Verbindung, mit und ohne Option, YAML und SQLite | mit Option: YAML `empty_sessions: true`, `sessions: []`; SQLite `meta` mit `empty_sessions|true`, 0 Sessions; ohne Option: keine Kennzeichnung |
| 14 | MV *Ergebnis*: `select key, value from meta` liefert `empty_sessions` mit `true` neben `format`, `version`; §SQLite-Tabelle `meta` | Beispieldatei gegen eine mit Option aufgezeichnete SQLite-Datei; ohne Option | `format|pgwire-recorder`, `version|1`, `empty_sessions|true`; ohne Option nur die ersten zwei |
| 15 | MV: keine Sitzung, auch mit Option: Aufbau scheitert (`PGR-E4002`, Exit 4) | `--upstream 127.0.0.1:1`, eine Verbindung, mit und ohne Option | `PGR-E4002` an den Client, Exit 4, `sessions: []` |
| 16 | MV: `COPY` (`PGR-E6001`, Exit 6); Verbindungsfehler: die andere Verbindung wird aufgezeichnet | Verbindung A `COPY (SELECT 1) TO STDOUT`, Verbindung B offen mit zwei Anfragen, mit und ohne Option | `PGR-E6001` auf A; B vollständig; Aufzeichnung mit nur B (`- id: 1`); Exit 6 |
| 17 | MV: Zwangsende im Aufbau zur Datenbank: keine Sitzung, ohne Meldung, Exit 0 | stummer Upstream, `--shutdown-timeout 300ms`, eine Verbindung im Aufbau, `docker stop`, mit und ohne Option | Exit 0, kein `PGR-E4006`, `sessions: []` |
| 18 | Handbuch *Aufzeichnen*, Passwort: die Verbindung wird nicht aufgezeichnet, auch nicht mit Option | Datenbank mit Passwort (SCRAM), mit Option, eine Verbindung | `PGR-E6001` (Anmeldeverfahren), Exit 6, `sessions: []` mit Kennzeichnung |
| 19 | MV: eine Verbindung zählt erst mit der Startnachricht; ohne Daten, Anfrage nach Verschlüsselung (SSL und GSSAPI), Abbruchwunsch sind keine Sitzung | `printf` der Bytes (`SSLRequest`, `GSSENCRequest`, `CancelRequest`) und leere Verbindung an `record --record-empty-sessions` | Antwort `N` auf beide Anfragen nach Verschlüsselung; `PGR-W3001` für den Abbruchwunsch; Exit 0; `empty_sessions: true`, `sessions: []` |
| 20 | MV: abgebrochene Startnachricht ab 8 Bytes `PGR-E6001`, Exit 6, keine Sitzung; weniger als 8 Bytes still | `printf` von 1, 4, 7 Bytes; von 8 Bytes (Kopf mit Länge 40); von 20 Bytes; 8 Bytes ohne Option | 1, 4, 7: Exit 0, keine Meldung, `sessions: []`; 8 und 20: Antwort `PGR-E6001 … Startnachricht nicht lesbar`, Exit 6, `sessions: []`; ohne Option ebenso, ohne Kennzeichnung |
| 21 | MV: mit `--format sqlite` dasselbe, `replay` liest die Datei wie dieselbe Aufzeichnung als YAML, auch mit `connection` | `f` als SQLite, `replay --session-assignment connection`, Verbindungen wie Zeile 24 | wie Zeile 24: Exit 0 mit `--fail-on-unconsumed`; zweite Verbindung mit Anfrage `PGR-E5001` für Session 2 |
| 22 | Handbuch *Wiedergeben*, Hinweis zu `first-request`: n-te Verbindung mit Anfrage erhält n-te Sitzung mit Interaktionen, leere übersprungen, Verbindungen ohne Anfrage zählen nicht | `f` (Session 2 leer): zwei Verbindungen mit `eins`, `drei`, mit `--fail-on-unconsumed`; dieselben zwei und eine dritte ohne Anfrage, ohne die Option | beide Anfragen beantwortet (`eins`, `drei`), Exit 0, keine Meldung |
| 23 | derselbe Hinweis: eine Anfrage über die Sitzungen mit Interaktionen hinaus `PGR-E5003` | `f`, drei Verbindungen mit Anfrage | dritte: `PGR-E5003` (`Verbindung 3: keine aufgezeichnete Session mehr frei`), Exit 5 |
| 24 | MV *Wiedergeben*: `connection`: n-te Verbindung die Sitzung n, auch ohne Anfrage und ohne Interaktionen; `--fail-on-unconsumed` ohne Meldung | `f`, `--session-assignment connection --fail-on-unconsumed`, Verbindungen: `eins`, ohne Anfrage, `drei` | beide Antworten, Exit 0, keine Meldung |
| 25 | MV: Anfrage an eine Sitzung ohne Interaktionen `PGR-E5001`; E5001-Zeile: „erwartet die Aufzeichnung keine, nennt sie das“ | `f`, `connection`, zweite Verbindung `select 'zwei'` | `PGR-E5001: Session 2: die Aufzeichnung erwartet keine Interaktion; empfangen "select 'zwei';"`; die dritte Verbindung beantwortet; Exit 5 |
| 26 | MV: keine Sitzung mit der `id` der Verbindung `PGR-E5003`; Lebendprüfung beantwortet | `f`, `connection`: vierte Verbindung `select 'vier'`; vierte Verbindung `-- ping` | `PGR-E5003: Verbindung 4: keine aufgezeichnete Session mit der Kennung 4`, Exit 5; `-- ping` ohne Fehler (`psql` Exit 0), Exit 0 |
| 27 | MV: `connection` verlangt die Kennzeichnung, sonst `PGR-E2001` beim Start (Exit 2); `first-request` startet mit derselben Aufzeichnung; E2001-Zeile | `a` (ohne Option aufgezeichnet) und `c2` (SQLite) mit `connection`; `a` mit `--session-assignment first-request` | `PGR-E2001: … trägt keine Kennzeichnung empty_sessions, die --session-assignment connection verlangt`, Exit 2; mit `first-request` Verbindung beantwortet |
| 28 | Handbuch *Wiedergeben* und MV: Sitzung ohne Interaktionen weder nicht verbraucht noch nie zugeordnet; eine Sitzung mit Interaktionen ohne Verbindung ist nie zugeordnet (`PGR-W2001`, `PGR-E5002`); Zeilen `PGR-E5002`, `PGR-W2001`, Exit-Code 5, Herunterfahren-Satz, Absatz unter den Exit-Codes | `f` mit `first-request`: zwei Verbindungen (Exit 0, keine Meldung, auch mit `--fail-on-unconsumed`), eine (W2001 beziehungsweise E5002, Kennung 3); `connection`: nur Verbindung 1 ohne und mit Option | Exit 0 ohne Meldung; `1 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 3` als `PGR-W2001` (Exit 0) beziehungsweise `PGR-E5002` (Exit 5); Session 2 (leer) nie genannt |
| 29 | MV: Verbindung zählt ab der Startnachricht; Anfrage nach Verschlüsselung, Abbruchwunsch, leere Verbindung zählen nicht; nach der Anfrage nach Verschlüsselung gesendete Startnachricht zählt einmal | `connection` und `first-request`: vor den Verbindungen `SSLRequest`, `GSSENCRequest`, `CancelRequest`, leere Verbindung (`printf | nc`); ferner `PGSSLMODE=prefer` | Antwort `N`, `PGR-W3001`; die folgenden Verbindungen erhalten Session 1, 2, 3 beziehungsweise `eins`, `drei`; Exit 0 mit `--fail-on-unconsumed`; mit `prefer` zählt die Verbindung einmal (`eins`, `drei`, Exit 0) |
| 30 | MV und Zeile `PGR-E3004`: nur Sitzungen ohne Interaktionen, `PGR-E3004` mit beiden Zuordnungen; Voraussetzung *Wiedergeben* | `d` (zwei leere Sessions), `replay` mit `first-request` und `connection`; `replay` einer Aufzeichnung ohne Session mit Kennzeichnung (YAML, SQLite) | je `PGR-E3004: … enthält keine Session mit Interaktion`, Exit 3 |
| 31 | MV: `--session-assignment` nur `first-request`, `connection` in Kleinbuchstaben; ungültig `PGR-E2001`, Datei `PGR-E2004`, Variable nennt den Namen | `Connection`, `first_request`, leer, `GEHEIM` als Option; Variable `Connection`; Schlüssel `Connection`, `""` | Option und Variable `PGR-E2001` (`erlaubt sind first-request und connection`, bei der Variable mit ihrem Namen), Datei `PGR-E2004` (`replay.session_assignment`), je Exit 2 |
| 32 | MV: Variable gilt wie die Option, leer gilt als nicht gesetzt, Option vor Variable vor Datei; Beispiel 4 als Datei und Beispiel 3 wörtlich | zweite Verbindung `drei` in `f`: Datei `connection`; Datei mit Option/Variable `first-request`; Variable `connection` (mit Option `first-request`); Variable leer mit Datei `connection`; Standard; `config show` der Beispieldatei; Beispiel 3 wörtlich gegen die Aufzeichnung aus Zeile 1 | `connection` (Datei, Variable, Option, Variable leer): `PGR-E5001` für Session 2; `first-request` (Option vor Variable und Datei, Variable vor Datei, Standard): `drei`; `config show` Exit 0, Inhalt wie geschrieben; Beispiel 3 und 4: Verbindungen beantwortet, Exit 0 |
| 33 | MV *Verbindungsfehler*: Abweichung auf einer Verbindung, die andere läuft durch, Exit 5; `COPY` auf einer, die andere aufgezeichnet, Exit 6; ein späterer Fehler anderer Klasse ändert den Exit-Code nicht | `g`: A (`q1`, nach 3 s `q2`) offen, B `select 'falsch'`; zu 16: `COPY`, danach eine laufende `pg_sleep(20)` bis zum Zwangsende (`--shutdown-timeout 500ms`) | A beantwortet beide, B `PGR-E5001`, Exit 5; Zeile 16 Exit 6 mit `PGR-E6001`, dann `PGR-E4006` (Exit bleibt 6) |
| 34 | MV *Ungültige Aufrufe*: Option beim anderen Kommando, Variable beim anderen Kommando (auch `false`), Schlüssel im anderen Abschnitt; Absatz in §5 | `--session-assignment` bei `record`, `play`; `--record-empty-sessions` bei `replay`, `play`; die Variablen bei den anderen Kommandos (auch `RECORD_EMPTY_SESSIONS=false`); Schlüssel `record:`/`play:` `session_assignment`, `replay:`/`play:` `record_empty_sessions` | Option: `PGR-E2001 flag provided but not defined`; Variable: `PGR-E2001 … gehört zu einer Option, die <Kommando> nicht kennt`; Schlüssel: `PGR-E2004 … unbekannter Schlüssel`; je Exit 2 |
| 35 | Tabelle §5: Betriebsart und Standard | `record --help`, `replay --help`, `play --help`, `config show --help`, `grep` | `--record-empty-sessions[=true\|false]` bei `record`, `--session-assignment first-request\|connection` mit Standard `first-request` bei `replay`, bei `play` und `config show` keine der beiden Optionen |
| 36 | Handbuch *Einspielen*: Sitzung mit Interaktionen läuft über eine Verbindung, für eine ohne öffnet `play` keine | `play` von `b` (leer, `select 1`, leer) und von `d` (zwei leere); Zahl der Zeilen `connection authorized` im Log der Datenbank | `b`: 1 Verbindung, `eingespielt=1`; `d`: 0 Verbindungen, `eingespielt=0`, Exit 0 |
| 37 | Zeile `PGR-E2001`/`PGR-E2004`: neue Fälle | Zeilen 5, 6, 31, 27 | wie dort |
| 38 | README: `record` zeichnet Verbindungen als Sitzungen auf, mit `--record-empty-sessions` auch ohne Anfrage; `first-request`, `connection` (verlangt Kennzeichnung, sonst `PGR-E2001`); `PGR-E5003`, `PGR-E5001` | Zeilen 1, 7, 9, 10, 22 bis 27 | wie dort |
| 39 | Glossar: Sitzung ohne Interaktionen entsteht nur mit `--record-empty-sessions` | Zeilen 9, 10 | wie dort |
| 40 | Hinweis *Aufzeichnen* (Passwort) und Zeilen `PGR-E3004`, `PGR-E5003`: neue Teilsätze | Zeilen 18, 30, 23, 26 | wie dort |

**Stand der Gates.** Wird nach dem Lauf eingetragen.

**Neue und geänderte Sätze neben Bestand** (Schritt 17; Merkmal aus den Aussagen des neuen Abschnitts, die eine Allgemeinregel des Bestands einschränken: „jede Verbindung ohne Anfrage ist keine Sitzung“, „jede aufgezeichnete Sitzung hat Anfragen“, „eine beendete Verbindung schreibt die Datei“, „jede Sitzung läuft über eine Verbindung“, „der Exit-Code ist der des ersten Fehlers“; `grep` des Handbuchs nach `Sitzung`, `ohne Anfrage`, `beendete Verbindung`, `nie verwendet`, `Exit-Code`):

| Stelle | Änderung | Beleg |
|---|---|---|
| *Aufzeichnen*, Ergebnis; Hinweis zu `--force`; SQLite-Abschnitt (Hinweise, Tabelle `meta`/`session`) | „nach dem Ende jeder Verbindung“ und „beendete Verbindung“ zu „Sitzung“ (eine Verbindung ohne Anfrage legt die Datei nicht an); Tabelle nennt `empty_sessions` | Zeilen 12, 13, 14 |
| *Aufzeichnen*, Hinweis zu Verbindungen ohne Anfrage und zum Passwort | Verweis auf `--record-empty-sessions`; Passwort: auch mit der Option nicht aufgezeichnet | Zeilen 9, 10, 18 |
| *Aufzeichnen*, Hinweis zum Beenden („Eine Verbindung ohne begonnene Anfrage, auch eine noch im Aufbau, endet ohne Meldung“; „Die abgeschlossenen Anfragen … bleiben“) | unverändert: bleibt wahr (keine Meldung, Exit 0 auch mit Option, Zeilen 11, 17) | Zeilen 11, 17 |
| *Wiedergeben*, Voraussetzung („mit mindestens einer Sitzung mit einer Anfrage“) | enger gefasst: nur leere Sitzungen sind `PGR-E3004` | Zeile 30 |
| *Wiedergeben*, Hinweis zur Zuordnung; Hinweis zu `PGR-W2001`/`PGR-E5002` (Sitzungen ohne Interaktionen gehören nicht dazu) | `first-request` genannt, Sitzungen mit Interaktionen | Zeilen 22, 23, 28 |
| *Einspielen*, Ergebnis („Jede Sitzung mit Interaktionen läuft über eine eigene Verbindung“); TLS-Hinweis | enger gefasst beziehungsweise umformuliert („jede Verbindung, die `play` öffnet“): keine neue Zusage über TLS, keine neue Probe | Zeile 36 (Verbindungen); die TLS-Probe liegt in `slice-v1-abschluss-einspielen-tls-doku` |
| *Herunterfahren mit Frist* (`--fail-on-unconsumed`: „eine aufgezeichnete Sitzung mit Anfragen nie verwendet“); Exit-Code 5; Absatz unter den Exit-Codes („Nie verwendete Sitzungen mit Anfragen“) | enger gefasst | Zeile 28 |
| Absatz „Ein Fehler, der nur eine Verbindung betrifft … Exit-Code des ersten aufgetretenen Fehlers“ unter den Exit-Codes | unverändert: der Satz gilt auch für die Fälle des neuen Abschnitts | Zeile 33 |
| Zeilen `PGR-E5001`, `PGR-E5002`, `PGR-E5003`, `PGR-W2001`, `PGR-E3004`, `PGR-E2001`, `PGR-E2004` | Teilsätze zu Sitzungen ohne Interaktionen und zu den beiden Optionen | Zeilen 5 bis 6, 23 bis 31 |

**Ungenannt geblieben** (keine Probe am gebauten Binary oder bewusst nicht im Handbuch): das aufgezeichnete Ende des Servers (`FATAL`, `PANIC`; Lücke vom Nutzer akzeptiert, Entscheidung vom 2026-10-12, `slice-v1-abschluss-verbindungsende-*`); eine Verbindung, die nur Lebendprüfungen stellt, beim Aufzeichnen; die Zuordnung bei gleichzeitigem Aufbau (der Satz „nur zugesichert, wenn nacheinander“ sagt zu, was das Binary tut, nichts: keine Probe); die Kennzeichnung ohne leere Session bei `connection`; das Verhalten von `replay` bei einer Aufzeichnung, die die Kennzeichnung trägt und von Hand geändert wurde; die Meldungstexte der Fälle jenseits der Codes und Exit-Codes; die Grenze „weniger als 8 Bytes“ ist an 1, 4 und 7 Bytes gezeigt.

**Befund an den Planner.** Keiner. Das Binary zeigt kein Verhalten, das die Spezifikation anders regelt, soweit ich sie für die Zeilen oben gelesen habe (Sitzung bei unvollständiger erster Anfrage mit Option wird geschrieben: `PGR-E4006` nennt es selbst, Zeile 11).

**Befehlsfolge aus `slice-doku-ist-stand` §7** (Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD`), am Stand `1fe572a`:

| Prüfung | Ausgabe |
|---|---|
| Optionen des Handbuchs, die `--help` nicht kennt | `--h`, `--help`, `--version` (Nicht-Optionen im Handbuch), `--name`, `--rm` (`docker run`) |
| Variablen ohne Option, ohne `PGWIRE_RECORDER_PASSWORD` | leer |
| Codes außerhalb des Katalogs | leer |

`--record-empty-sessions`, `--session-assignment` und ihre Variablen stehen nicht in der Ausgabe. Die Befehlsfolge ist nicht geändert, kein Wächter eingeführt; ihre Mutanten stehen in `done/slice-v1-abschluss-einspielen-anmeldung-doku` §7 und werden nicht wiederholt. `grep -n -i -E "spec/|lastenheft|slice|welle|review|LH-|SPEC-|ADR|noch nicht|kommt|geplant|künftig"` über das Handbuch findet nur Bestand ohne Bezug (`noch nicht` in §3 „Zieldatei existiert noch nicht“, `kommt` in der Ruhefrist von `record` und in §4 zu TLS, `Adresse` als Treffer auf `ADR`); über die Zeilen, die dieser Diff in `README.md` hinzufügt, findet er nichts. `make docs-check`: 669 Dateien, 0 Befunde.

**Mutation (§3.10)** entfällt: Der Slice liefert keinen neuen Vertrag, nur Beschreibung; die Probe-Tabelle ersetzt sie (rot wäre ein Satz ohne bestandene Probe).

**Kommentare (Schritt 20)** entfallen: kein Code, kein Skript und keine Konfiguration geändert.

**Rückgaben an den Architect:** keine; keine Randform entschieden. Zur Frage aus §1 (nennt das Handbuch die Verbindungen, die keine Session sind): Es nennt sie, soweit Zeilen 17, 19 und 20 sie am Binary zeigen: ohne Daten, Anfrage nach Verschlüsselung (beide Arten unter diesem Wort), Abbruchwunsch, abgebrochene Startnachricht (1, 4, 7, 8 und 20 Bytes), Zwangsende im Aufbau. Nicht genannt ist `GSSENCRequest` unter eigenem Namen. **Nachzählen:** Liefer-Punkte 2, Schichten 1 (Nutzerdokumentation), wie in §8.

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
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-sessions`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code. Nachgezählt beim Eintragen der Übernahme aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12): Sie gehört zu DoD-Punkt 1 (Handbuch), kein neuer Liefer-Punkt; Liefer-Punkte weiter 2, Schichten weiter eine.

Nachgezählt beim Eintragen der Übernahme aus `slice-v1-abschluss-sessions` (Architect, 2026-10-12): Sie gehört zu DoD-Punkt 1 (Handbuch) und DoD-Punkt 2 (README, soweit es die Optionen nennt); kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Dokumentation).

Nachgezählt beim Eintragen der Verbindungen, die keine Session sind (Planner, 2026-10-12, `AGENTS.md` §3.13, Befund V-308 der Verifikation von `slice-v1-abschluss-sessions`): Sie gehört zu DoD-Punkt 1 (Handbuch); kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Dokumentation).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
