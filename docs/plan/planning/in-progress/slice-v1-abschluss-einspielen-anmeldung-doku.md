# Slice slice-v1-abschluss-einspielen-anmeldung-doku: Handbuch und README zu: Anmeldung beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-einspielen-anmeldung`)

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

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-einspielen-anmeldung` liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-einspielen-anmeldung`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-einspielen-anmeldung` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-einspielen-anmeldung` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-einspielen-anmeldung` liefert, im Ist-Zustand des gebauten Binaries: Aus `slice-doku-ist-stand` (dort §6, *Teilweise geliefert*): Der Slice ersetzt im Handbuch die Grenz-Sätze zum Passwort bei `play` — §1 *Voraussetzungen*, §4 *Eine Aufzeichnung in eine Datenbank einspielen* (Voraussetzung und Hinweise: „ohne Passwort“), §5 *Konfigurationsdatei* („Ein Passwort in der URL setzt `play` ein, meldet sich damit aber nicht an“, dort mit DoD-Punkt 3) und §7 Zeile `PGR-E4005` („`play` meldet sich nur ohne Passwort an“) — und im README den Satz, dass die Datenbank den Benutzer ohne Passwort anmelden muss, soweit er `play` betrifft. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [ ] Das Benutzerhandbuch beschreibt in §5 *Konfigurationsdatei* das Passwort einer Verbindung bei `play` wie geliefert (aus dem Platzhalter, sonst `PGWIRE_RECORDER_PASSWORD`, wenn die Verbindung kein Passwort schreibt; Befund F-533) und nennt die Grenzen *SASLprep* und *Klartext ohne TLS*.
- [ ] `README.md` nennt, was `slice-v1-abschluss-einspielen-anmeldung` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
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
| `docs/user/benutzerhandbuch.md` | update | §1 *Voraussetzungen*, §4 *Eine Aufzeichnung in eine Datenbank einspielen* (Voraussetzung, Hinweise), §5 (Absatz hinter der Optionstabelle, Beispiel, Platzhalter-Absatz), §6 *Rollen und Rechte* (Nachbarsatz zur Anmeldung, der für `play` nicht mehr gilt) und §7 (Zeilen `PGR-E4005`, `PGR-E6001`) zu dem, was `slice-v1-abschluss-einspielen-anmeldung` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt |
| `README.md` | update | Sätze, die `slice-v1-abschluss-einspielen-anmeldung` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5 (§6), meldet die Zeile `comm -23 $W/v-hb $W/v-help` also bei jedem Lauf. In diesem Plan lautet die Variablen-Zeile `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`; mit ihr ist die Ausgabe der Variablen leer (Optionen: nur `--h`, `--help`, `--version`, `--name`, `--rm`; Codes: leer). Die Datei `done/slice-doku-ist-stand.md` bleibt unverändert. Der Mutant „Zeile `PGWIRE_RECORDER_PASSWORD` zurück“ entfällt für diese Variable und wird durch eine andere Variable ersetzt (§7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen-anmeldung` liegt in `done/` (Voraussetzung; Schritt 13 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

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

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect prüft vor dem ersten Commit, ob `slice-v1-abschluss-einspielen-anmeldung` eine weitere Randform liefert, die dort nicht steht.

**Randformen dieses Slice** (Architect, vor dem ersten Commit, 2026-10-10; Ort der Entscheidung ist dieser Plan, weil der Slice keinen Vertrag anlegt). Belegt sind sie durch Proben gegen das gebaute Binary und das gepinnte PostgreSQL 17 mit Servermodi `scram-sha-256`, `md5`, `password` und `trust` (Benutzer wie in `tools/test/run-integration-tests.sh`) oder, wo ein Server sie nicht auslöst, durch einen benannten Test. Der Implementer wiederholt die Proben und trägt sie in §7 ein; weicht ein Ergebnis ab, gilt das Ergebnis und die Zeile hier geht als Befund an den Planner.

| Aussage im Handbuch | Entscheidung | Beleg | Ort |
|---|---|---|---|
| Passwort aus Platzhalter, sonst `PGWIRE_RECORDER_PASSWORD` | genannt; schreibt die Verbindung einen Passwortteil, gilt nur dessen eingesetzter Wert, die Variable bleibt unbeachtet; ohne Passwortteil und bei `host:port` gilt die Variable | Probe: Platzhalter richtig und Variable falsch meldet an; Verbindung ohne Passwortteil mit und ohne Variable; `host:port` mit Variable | §5, ein Absatz hinter der Optionstabelle (keine Option, daher keine Zeile der Tabelle); §4 *Voraussetzung* verweist dorthin |
| Klartext, MD5, SCRAM-SHA-256 | genannt, je Verfahren eine Probe mit richtigem Passwort | drei Benutzer, Anmeldung gelingt (der Lauf erreicht die Anfragen) | §4 *Voraussetzung* |
| leere Variable | genannt: gilt als nicht gesetzt, wie bei Platzhaltern | Probe: leere Variable gegen `scram-sha-256` endet mit `PGR-E4005` | §5, im Absatz zur Passwortquelle |
| fehlendes Passwort | genannt: `PGR-E4005`, die Meldung sagt, dass `play` keines hat; **nicht** genannt: „ohne etwas zu senden“, das zeigt keine Probe von außen | Probe ohne Variable gegen `scram-sha-256` | §7 Zeile `PGR-E4005` |
| falsches Passwort, unbekannter Benutzer | genannt: `PGR-E4005`, die Meldung nennt SQLSTATE und Meldung der Datenbank (wie bei `PGR-E4004`); der Wortlaut der Datenbank steht **nicht** im Handbuch, er ist Sache des Servers | Probe je Verfahren mit falschem Passwort (`28P01`), unbekannter Benutzer (`28000`) | §7 Zeile `PGR-E4005` |
| Server ohne Passwortverlangen (`trust`) | genannt: eine gesetzte Variable ist ohne Wirkung, der Lauf gelingt; **nicht** genannt: „sendet kein Passwort“ (nicht beobachtbar) | Probe `trust` mit und ohne Variable | §5, im Absatz zur Passwortquelle |
| Passwort erscheint nicht in Ausgaben | genannt: keine Meldung, keine Log-Zeile und kein `config show` nennt es; `config show` nennt den Namen der Variable, nie den Wert | Probe: Marker `GEHEIM` im Passwort, `grep` über Ausgabe von richtigem und falschem Passwort, Log-Stufe `debug`, `config show`; dazu der E2E-Test des Code-Slice | §5 |
| Klartext ohne TLS | genannt als Sicherheitshinweis in der Form, die vor und nach dem TLS-Slice wahr ist: „Verlangt die Datenbank ein Klartext-Passwort, sendet `play` es über die Verbindung, wie sie ist, also ohne TLS unverschlüsselt.“ Kein „solange“, kein „noch nicht“ | Probe: Klartext-Benutzer meldet sich über die unverschlüsselte Verbindung an | §4 *Hinweise*, neben dem Satz zur unverschlüsselten Verbindung (den ersetzt `slice-v1-abschluss-einspielen-tls-doku`); MD5 und SCRAM nennen den Hinweis nicht |
| SCRAM ohne Channel Binding | genannt als eine Zeile: bietet die Datenbank nur `SCRAM-SHA-256-PLUS` an, endet `play` mit `PGR-E4005`. Nur wenn der Implementer einen Test findet, der genau diesen Fall prüft (`internal/adapters/driven/postgres/einspielen_anmeldung_test.go` nennt `SCRAM-SHA-256-PLUS`); sonst gestrichen | Test, kein Server löst es aus | §7 Zeile `PGR-E4005` |
| Verfahren, das `play` nicht kann (Kerberos, GSSAPI, SSPI) | wie die vorige Zeile: genannt, wenn ein Test den Fall prüft, sonst gestrichen; die Liste der Verfahren, die `play` kann, steht in der Zeile | Test | §7 Zeile `PGR-E4005` |
| SASLprep | genannt als Grenze nur nach Probe: ein Passwort mit einem Zeichen, das SASLprep ändert (etwa U+00A0), gegen `scram-sha-256` endet mit `PGR-E4005`; dann lautet der Satz „`play` sendet das Passwort unverändert; ein Passwort mit Zeichen, die SASLprep ändert, kann bei SCRAM scheitern.“ Scheitert die Probe nicht, steht keine Grenze im Handbuch, und das geht als Befund an den Planner (der Satz der Spezifikation behauptet mehr, als das Binary zeigt) | Probe | §5, im Absatz zur Passwortquelle |
| Frist, nicht unterbrechbare Berechnung, Obergrenze der Iterationen | **ungenannt**, akzeptiertes Negativ: Ein Server mit mehr als 10 000 000 Iterationen kommt nur von einem fehlerhaften oder feindlichen Server (PostgreSQL wählt standardmäßig 4096), der Nutzer hat dagegen keinen Handgriff, und eine Zahl im Handbuch wäre eine Zusage mehr, die jede Änderung der Grenze nachziehen müsste. Die Meldung `PGR-E4005` der Anmeldung deckt den Fall, ohne dass er einen Satz braucht | — | — |
| `PGR-E4005` ohne Passwort-Aussage „nur ohne Passwort“ | die Zeile in §7, §1 *Voraussetzungen*, §4 *Voraussetzung* und *Hinweise*, §5 und der Satz zu `record` in §7 *Die Anwendung kann sich nicht verbinden* bleiben nur dort, wo sie `record` betreffen; `record` wird nicht geändert (Passwort verlangt: `PGR-E6001`) | `grep -n "ohne Passwort"` über das Handbuch, jeder Treffer entweder `record` oder ersetzt | §1, §4, §5, §7 |
| Beispiel mit Passwort-Platzhalter | genannt: eine Verbindung `postgresql://benutzer:${PW}@host:5432/db` im Beispiel von §5, als Datei mit `config show` und `play` gegen den Server in der Probe; eine Variable mit Sonderzeichen (`$`, Leerzeichen) im Wert, weil das Handbuch zusagt, dass der Wert unverändert bleibt | Probe | §5 |
| Nachbarsatz in §6 *Rollen und Rechte* (Entscheidung des Implementers, vom Reviewer als tragend bestätigt) | der Satz „Das Werkzeug verwaltet keine Benutzer und Rollen und prüft keine Zugangsdaten“ gilt für `record` und `replay` wie vorher und steht neben dem neuen Passwort-Verhalten von `play`; er nennt für `play` die Anmeldung. Wortlaut: „Beim Aufzeichnen leitet es Benutzer und Datenbank der Anwendung an die Datenbank weiter, ohne Passwort; beim Wiedergeben nimmt es jede Anmeldung an, gleich mit welchem Benutzer, welcher Datenbank und welchem Passwort; beim Einspielen meldet es sich mit Benutzer, Datenbank und, wenn die Datenbank eines verlangt, dem Passwort an der Datenbank an.“ | Probe: Einspielen mit Passwort gegen SCRAM, MD5 und Klartext (Zeilen oben in §7); `record` und `replay` sind unverändert | §6 des Handbuchs, erster Absatz |
| Befehlsfolge aus `slice-doku-ist-stand` §7 und die Variable `PGWIRE_RECORDER_PASSWORD` | benannte Ausnahme: die Variable ist keine Option, die Variablen-Zeile der Befehlsfolge zieht genau sie mit `grep -vx` ab (§3); jede andere Variable ohne Option bleibt rot | Lauf mit Mutanten (§7) | §3, §7 |

*Eine Aufzeichnung für die Proben* enthält nur Anfragen, die der Probe-Benutzer darf (`SELECT 1`); eine Anweisung auf einer Tabelle des Superusers endet nach gelungener Anmeldung mit `PGR-E4004` und lässt jede Probe wie einen Fehlschlag aussehen.

*Zuweisungen* (`AGENTS.md` §3.13): keine. Der Klartext-Satz ist so gefasst, dass er nach `slice-v1-abschluss-einspielen-tls` wahr bleibt; der Satz „verbindet sich unverschlüsselt“ gehört schon in dessen DoD.

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

**Probe-Aufbau.** `make build` am Stand `36480b8` (Image `pgwire-recorder:dev`, `sha256:c2871ccf3ceb…`, Produktcode unverändert). Proben in einem Scratch-Verzeichnis außerhalb des Repos; Netz `impl-anmdoku-net`, ein PostgreSQL 17 aus dem gepinnten Image von `harness/mk/integration.mk` (`impl-anmdoku-pg`, Alias `postgres`), Aufzeichnung und `play` als `docker run … pgwire-recorder:dev` im selben Netz. Benutzer wie in `tools/test/run-integration-tests.sh` (`play_scram` mit Passwort `GEHEIM scram $%41` — Leerzeichen, `$`, `%` —, `play_md5`, `play_pw`), dazu `play_trust` (`trust`), `play_nbsp` (SCRAM, Passwort `GEHEIM` U+00A0 `x`), `play_sonder` (SCRAM, Passwort `a@b:c/d?e#f`). Die Aufzeichnung der Proben besteht aus `SELECT 1` (aufgezeichnet mit `record` und `psql` als `postgres`). Container, Netz, Images und Scratch-Dateien sind danach entfernt. Nur Linux im Container (`BEO-REPO/verhalten-nur-unter-linux-geprueft`).

**Proben je Aussage** (`PL` = `play --upstream postgres:5432 --input sel.yaml`; „Variable“ = `PGWIRE_RECORDER_PASSWORD`; „ok“ = `play beendet`, Exit 0):

| Aussage im Handbuch | Probe | Ergebnis |
|---|---|---|
| Klartext, MD5, SCRAM-SHA-256 melden an | `--user play_pw` / `play_md5` / `play_scram`, Variable richtig | je ok; SCRAM mit Leerzeichen und `$` im Wert |
| Variable gilt bei `host:port` | dieselben Läufe (`host:port`) | ok |
| Platzhalter vor Variable | Verbindung `mit` (`play_scram:${PW}@…`), Platzhalter richtig, Variable falsch | ok |
| dasselbe umgekehrt | Platzhalter falsch, Variable richtig | `PGR-E4005`, `28P01`, Exit 4 |
| Platzhalter ohne Variable | Platzhalter richtig, Variable ungesetzt | ok |
| Verbindung ohne Passwortteil: Variable gilt | Verbindung `ohne` mit Variable richtig | ok |
| Verbindung ohne Passwortteil, ohne Variable | `ohne`, Variable ungesetzt | `PGR-E4005` „der Server verlangt ein Passwort, play hat keines“, Exit 4 |
| leere Variable gilt als nicht gesetzt | `PL` mit `PGWIRE_RECORDER_PASSWORD=` gegen `play_scram`, und `ohne` mit leerer Variable | je `PGR-E4005` „… play hat keines“ |
| leerer Platzhalter-Wert / ungesetzter Platzhalter | `PW=` bzw. ungesetzt, Verbindung `mit` | `PGR-E2005`, Exit 2 (Bestand, unverändert) |
| `trust`: Variable ohne Wirkung | `--user play_trust` mit und ohne Variable | beide ok |
| falsches Passwort nennt SQLSTATE und Meldung der Datenbank | falsches Passwort gegen SCRAM, MD5, Klartext | je `PGR-E4005`, „Fehlerantwort im Aufbau 28P01 „password authentication failed for user …““, Exit 4 |
| unbekannter Benutzer | `--user niemand` | `PGR-E4005`, `28000` „role … does not exist“, Exit 4 |
| fehlendes Passwort: nur `PGR-E4005` | `play_scram` ohne Variable | `PGR-E4005` „der Server verlangt ein Passwort, play hat keines“ (Handbuch nennt kein „ohne etwas zu senden“) |
| Platzhalter-Wert unverändert, auch Sonderzeichen | `play_sonder`, `PW=a@b:c/d?e#f` | ok (zu Zeichen `@ : / ? #` im Bestandssatz) |
| SASLprep-Grenze | `play_nbsp` mit `GEHEIM` U+00A0 `x` über die Variable; zum Vergleich `psql` mit demselben Wert | `play`: `PGR-E4005`, `28P01`, Exit 4; `psql`: meldet an (SASLprep wirkt). Die Grenze steht im Handbuch |
| Passwort erscheint nicht in Ausgaben | Marker `GEHEIM`, `--log-level debug`: richtiges und falsches Passwort je Verfahren (SCRAM, MD5, Klartext), Platzhalter richtig und falsch, Verbindung ohne Passwortteil; `config show` mit Platzhalter und gesetzter Variable; `PGR-E2006` (Klartext-Passwort `GEHEIMklar` in der Datei) mit `config show` und `play --log-level debug` | `grep -c GEHEIM` über jede Ausgabe: 0. `config show` endet mit der Zeile `PGWIRE_RECORDER_PASSWORD` (Name, nie Wert). Gegenprobe des `grep`: Verbindungsname `GEHEIMname` in der Datei, `config show`: 1 Treffer |
| Beispiel §5 als Datei | `ex.yaml` aus dem Handbuch extrahiert; `config show --config` mit `CI_DB_HOST=x` | Ausgabe stimmt mit der Datei überein, Exit 0 |
| Beispiel §5 gegen den Server | Zeile `test` mit `sed` auf `play_scram`, `postgres:5432`, Datenbank `postgres` gesetzt (nur diese drei Werte), `DB_PASSWORD=GEHEIM scram $%41`, `--upstream test`, ohne `--user` | ok, Benutzer aus der URL; ohne `DB_PASSWORD`: `PGR-E2005`, Exit 2 |
| Beispiel §4 | `play --upstream postgres:5432 --input ./recordings/users.yaml` (Benutzer `postgres`, Server `trust`) | ok |
| Klartext-Satz | `play_pw` über die unverschlüsselte Verbindung | ok (Probe 3); der Satz sagt nicht mehr als das |

**Verfahren, die `play` nicht kann (Zeile `PGR-E4005`).** Kein Server löst sie aus; belegt sind sie durch `TestAnmeldungNichtUnterstuetzt` (`internal/adapters/driven/postgres/einspielen_anmeldung_test.go`) mit den Fällen `nur PLUS` (`SCRAM-SHA-256-PLUS` allein), `Kerberos` (Code 2), `GSS` (7), `SSPI` (9): Ergebnis `PGR-E4005`, nichts gesendet. Der Fall `SCM` (6) steht im Test, nicht im Handbuch. Mutanten in einer Kopie des Baums (`git archive`, `docker build --target test`; der Arbeitsbaum blieb unberührt): (1) `slices.Contains(verfahren, scramName)` durch einen Präfix-Vergleich ersetzt → rot `TestAnmeldungNichtUnterstuetzt/nur_PLUS` und `/mit_Anhängsel`; (2) Codes 2, 7, 9 in `anforderung` angenommen (`return nil, nil`) → rot `TestAnmeldungNichtUnterstuetzt/Kerberos`, `/GSS`, `/SSPI`, außerdem `TestAnmeldungWeitereAnforderungArt` und `TestEinspielAufbauFehler`. Nicht gefahren: ein Mutant für die Kerberos-Zeile außerhalb von Code 2, 7, 9.

**Prüfung als Befehlsfolge** (die aus `done/slice-doku-ist-stand.md` §7, am Binary `36480b8`): Optionen des Handbuchs gegen `--help` aller vier Kommandos, Variablen gegen die Optionen, Codes gegen den Katalog. Ergebnis: Optionen nur `--h`, `--help`, `--version` (Handbuch nennt sie als Nicht-Optionen) und `--name`, `--rm` (`docker run`); Codes keine Ausgabe; Variablen genau `PGWIRE_RECORDER_PASSWORD` — **bekannte Ausnahme**: Die Variable ist keine Option und hat keine Tabellenzeile (§6), `--help` von `play` nennt sie (`internal/adapters/driving/cli/cli.go`, `envPassword`). Mutanten auf einer Kopie des Handbuchs: Satz mit `--upstream-tls` → rot (Optionen); `PGWIRE_RECORDER_PASSWORT` → rot (Variablen); `PGR-E4007` → rot (Codes). Die Befehlsfolge prüft Namen, nicht Zuordnung; die Zeilen der Tabelle in §5 sind unverändert (14), die Zuordnung der Sätze zum Verhalten belegt die Tabelle oben.

**`grep`** (Handbuch, README): `grep -n "ohne Passwort"` → README Zeile 29 (`record`), Handbuch §6 (Aufzeichnen), §7 `PGR-E6001` (Aufzeichnen), §7 *Die Anwendung kann sich nicht verbinden* (Aufzeichnen); jeder Treffer betrifft `record`. Kein Treffer für `Spezifikation|Lastenheft|ADR|slice|welle|review|noch nicht|geplant|solange|künftig|bisher|zunächst` im Handbuch außer dem Alltagswort „noch nicht“ in §3 („die Zieldatei existiert noch nicht“, unverändert, bezieht sich auf eine Datei); in der README keiner.

**Geänderte Aussagen.** Entfernt: „Zum Aufzeichnen und Einspielen … ohne Passwort … `play` endet mit `PGR-E4005`“ (§1), „meldet den Benutzer ohne Passwort an“ (§4 *Voraussetzung*), „verbindet sich unverschlüsselt und ohne Passwort“ und „Verlangt sie ein Passwort … `PGR-E4005`“ (§4 *Hinweise*), „Ein Passwort in der URL setzt `play` ein, meldet sich damit aber nicht an“ (§5), „`play` meldet sich nur ohne Passwort an … lassen Sie die Anmeldung ohne Passwort zu“ (§7 `PGR-E4005`), im README „die Datenbank muss den Benutzer ohne Passwort anmelden“. Neu: Passwortquellen und Grenzen (§5 Absatz hinter der Tabelle), Beispielverbindung `test`, Satz zum Platzhalter bei `play`, Klartext-Satz (§4), `PGR-E4005` mit Ursachen, Verfahren und Meldungsinhalt (§7), `PGR-E6001` auf `record` beschränkt (§7), README. Nachbarsatz neben den neuen Sätzen (Schritt 17): §6 *Rollen und Rechte* sagte „kennt keine Benutzer, Rollen oder Anmeldung … prüft keine Zugangsdaten“ für alle Betriebsarten; er steht jetzt für `record` und `replay` wie vorher und nennt für `play`, dass es sich mit Benutzer, Datenbank und gegebenenfalls Passwort anmeldet (Probe: Einspielen mit Passwort oben).

**Nicht gefahren.** Die übrigen Codeblöcke des Handbuchs (Aufzeichnen, Wiedergeben, Compose, Docker-Aufrufe) sind in diesem Lauf unberührt und nicht erneut als Datei gelaufen; die Nachbarsätze zu `record`/`replay` ohne Passwort-Aussage (§4 Aufzeichnen, Zeile `PGR-E6001` bei `record`) wurden nicht erneut geprobt. Iterationsobergrenze, Frist, unterbrechbare Berechnung: im Handbuch ungenannt (akzeptiertes Negativ, §6). Review, Verifikation, Closure: nicht Teil dieses Laufs.

**Nacharbeit zum Review (F-597 bis F-600).** *F-597:* Befehlsfolge als im Plan formulierte Folge mit der benannten Ausnahme (§3, §6), am Binary `pgwire-recorder:dev` (Produktcode unverändert) und am Handbuch nach der Nacharbeit: Optionen nur `--h`, `--help`, `--name`, `--rm`, `--version`; Variablen roh `PGWIRE_RECORDER_PASSWORD`, **mit Abzug leer**; Codes leer. Mutanten auf Kopien des Handbuchs (je mit angehängtem Satz): `PGWIRE_RECORDER_PASSWORT` → Variablen mit Abzug `PGWIRE_RECORDER_PASSWORT`, rot; `PGWIRE_RECORDER_UPSTREAM_TLS_X` und `PGWIRE_RECORDER_FORCE_X` → beide gemeldet, rot; `PGWIRE_RECORDER_PASSWORD_FILE` → gemeldet, rot (der Abzug mit `grep -vx` greift nur auf den vollen Namen). Der Mutant „Zeile `PGWIRE_RECORDER_PASSWORD` zurück“ ist für diese Variable nicht mehr trennscharf und entfällt (die Variable ist die Ausnahme); er bleibt für jede andere Variable erhalten. Merkmal des Befunds: eine Variable ohne Option, die die Folge als Normalfall liest; die Ausprägungen oben (Tippfehler, andere Variablen, Namensverlängerung) sind gefahren. *F-598:* Randform „Nachbarsatz in §6 *Rollen und Rechte*“ mit Wortlaut und Beleg in §6 eingetragen; das Handbuch ist dafür unverändert. *F-599:* Die Verweise in §1 *Voraussetzungen* und §4 *Voraussetzung* zeigen auf den Abschnitt *Einstellungen* des Handbuchs (Anker `#5-einstellungen`) (Absatz hinter der Optionstabelle, in dem Passwortquellen und Vorrang stehen) statt auf *Konfigurationsdatei*; der Absatz ist dort eindeutig, `make docs-check` löst den Anker auf (Gate-Lauf unten). *F-600:* Die SASLprep-Grenze hat als Beleg nur die Probe (Zeile `SASLprep-Grenze` oben), keinen Test; das Handbuch fasst sie mit „kann“ und sagt nicht mehr zu, als die Probe zeigt. Ein Test dafür ist nicht Teil dieses Slice (kein Produkt-Test, §1).

**Gate-Lauf.** `make gates` auf sauberem Baum am Stand `8cc0a48` (Handbuch, README und Plan wie im Diff dieses Slice): Exit 0, einschließlich `docs-check` und `lint-gegenprobe`. Die Zeile selbst folgt in einem eigenen Commit; sie ändert nur diesen Plan.

**Funde für den Planner.**

1. Die Befehlsfolge aus `slice-doku-ist-stand` §7 meldet `PGWIRE_RECORDER_PASSWORD` bei jedem Lauf (Variable ohne Option); sie braucht die benannte Ausnahme oder einen Abzug dieser einen Variable, sonst liest sich jeder Lauf rot.
2. §6 des Plans nennt die Randform „Nachbarsatz in §6 *Rollen und Rechte*“ nicht; dieser Lauf hat ihn nach Schritt 17 angepasst (Wortlaut in der Zeile oben). Architect und Review prüfen, ob er trägt.
3. `config show` listet `PGWIRE_RECORDER_PASSWORD` unter den gesetzten Variablen, obwohl die Variable keine Option ist; das Handbuch sagt es in der Zeile zu `config show` (§4) bereits für alle `PGWIRE_RECORDER_*`, der neue Absatz nennt es für die Passwort-Variable ausdrücklich.
4. SASLprep: die Probe scheiterte (`PGR-E4005`), die Grenze steht deshalb im Handbuch; die Spezifikation und das Binary stimmen überein, kein Befund.

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
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-einspielen-anmeldung`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 3. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
