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
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-einspielen-anmeldung` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt |
| `README.md` | update | Sätze, die `slice-v1-abschluss-einspielen-anmeldung` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe grün.

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
