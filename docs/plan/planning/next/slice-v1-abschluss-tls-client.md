# Slice slice-v1-abschluss-tls-client: TLS zum Client bei record und replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-23`](../../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client), [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [ADR-0018](../../adr/0018-tls-zum-client.md)

**Berührte Spec-Stellen:** `LH-FA-23.a` · `LH-FA-17.a` · `LH-FA-13.a` · `LH-FA-05.c` · `LH-FA-05.e` · `SPEC-018` · `SPEC-019` · `SPEC-020` · `SPEC-034` · `SPEC-033` · `ARC-006` · `ARC-005`

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

**Ziel:** `record` und `replay` nehmen auf Wunsch TLS von Clients an (`--tls-cert`, `--tls-key`), weisen unverschlüsselte Clients dann ab, sofern sie nicht zugelassen sind (`--allow-plaintext`), und verhalten sich ohne Konfiguration wie bisher.

**Übernommen aus `slice-v1-abschluss-upstream-verbinden`** (Hinweis F-532 aus dessen Review, zugewiesen bei dessen Closure): Die Prüfungen nach der Zusammenführung laufen im CLI-Adapter je Option in der Reihenfolge der Optionstabelle (`zuletzt`); `LH-FA-17.a` *Fehler* verlangt am Ende Pflichtoptionen, Kombinationen, `--upstream`. Die Kombinationen dieses Slice (`--tls-key` allein, `--allow-plaintext` ohne `--tls-cert`, `PGR-E2001`) kommen bei `record` vor jedem Fehler von `--upstream`, gleich, wo ihre Optionen in der Tabelle stehen (Test: ungültige Kombination neben ungültigem `--upstream` meldet die Kombination; gehört zu DoD-Punkt 2).

**Gebunden durch `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, Paarungen, offene Frage zu V-3; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): `--tls-cert`, `--tls-key` und `--allow-plaintext` liest der Stand noch nicht; `nichtGelesen()` in `internal/adapters/driving/cli/fremd_test.go` führt sie. Dieser Slice liefert den Leser, kürzt `nichtGelesen()` um diese Option(en) und bindet ihn in `liesUmgebung` (`internal/adapters/driving/cli/cli.go`) an; die Zuordnung `zuordnung()` nennt sie für `record` und `replay` (Spalte *Kommando* der Optionstabelle, `LH-FA-17.a`). Sonst wird `TestLeserUmgebungNichtGelesenerOptionen` rot, weil die Variable beim eigenen Kommando dann geprüft wird (`TestLeserZuordnungGleichTabelle` bleibt grün). Beim anderen Kommando (`play`) ist die Variable schon jetzt `PGR-E2001`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-tls-client-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- TLS zum Upstream im Record-Modus — Out-of-Scope von LH-FA-23; die Verbindung zum Server bleibt unverschlüsselt.
- Prüfung von Client-Zertifikaten und Zertifikatsverwaltung — Out-of-Scope von LH-FA-23.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-23`](../../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client): Ein Client verbindet sich verschlüsselt, `record` zeichnet auf und `replay` liefert die Aufzeichnung aus, jeweils wie bei einer unverschlüsselten Verbindung; dieselbe Aufzeichnung ist über beide Verbindungsarten gleich (Abnahmeszenario 15; Integrationstest).
- [ ] Ohne Konfiguration wird `SSLRequest` mit `N` beantwortet; mit Konfiguration wird ein unverschlüsselter Client (auch nach `GSSENCRequest`, und nur mit einer `StartupMessage` der Version 3.0) mit `PGR-E6003` abgewiesen, mit `--allow-plaintext` zugelassen, ein Klartext-`CancelRequest` bleibt `PGR-W3001`; mehrere Verbindungen nacheinander funktionieren; `--tls-cert` oder `--tls-key` allein und `--allow-plaintext=true` ohne `--tls-cert` sind `PGR-E2001`, vor jedem Fehler von `--upstream` (Test; die Quellen Kommandozeile, Umgebung und Datei, auch gemischt); die Leser der drei Optionen (auch ihrer Umgebungsvariablen) sind angebunden und `nichtGelesen()` in `fremd_test.go` führt sie nicht mehr (aus `slice-v1-abschluss-umgebung-anderes-kommando`).
- [ ] Nicht lesbare, ungültige, abgelaufene, verschlüsselte oder nicht passende Dateien beenden den Start mit `PGR-E2007` und Exit-Code 2, vor Ausgabedatei, Aufzeichnung und Listen; ein gescheiterter Handshake (auch Bytes vor dem `S`) endet mit `PGR-W3002` ohne Wirkung auf den Exit-Code; Herunterfahren und Zwangsende wirken auf eine verschlüsselte Verbindung wie auf eine unverschlüsselte; Schlüssel erscheinen weder in Logs noch in Meldungen noch in `config show` (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/pgwire` | update | TLS-Terminierung, Aushandlung nach `SSLRequest`, Abweisung unverschlüsselter Clients (`PGR-E6003`), `PGR-W3002`; baut die TLS-Konfiguration aus den PEM-Bytes (`crypto/tls` nur hier, Fitness Function von [ADR-0018](../../adr/0018-tls-zum-client.md)) |
| `internal/adapters/driving/cli` | update | Optionen `--tls-cert`, `--tls-key`, `--allow-plaintext` (Tabelle, Schlüssel der Datei, Hilfetext), Kombinationen, Laden und Prüfen beider Dateien nach `LH-FA-23.a` *Dateien* ohne `crypto/tls` (PEM, Gültigkeit und Passung über `crypto/x509`; die Regel des Architektur-Gates lässt `crypto/tls` im CLI-Adapter nicht zu und wird nicht gelockert, `AGENTS.md` §3.6), Fehlerabbildung auf `PGR-E2001` und `PGR-E2007`; der Schlüssel liegt in einem Typ, der ihn in keiner Textform ausgibt |
| `internal/hexagon/model` | update | Meldungscodes `PGR-E6003` und `PGR-W3002` (Katalog-Konstanten, keine Logik; zählt nicht als Schicht, siehe §8) |
| `internal/bootstrap` | update | reine Verdrahtung: reicht Zertifikat und Schlüssel (PEM-Bytes) an den PGWire-Server; ein Fehler der Konstruktion geht unverändert zurück; verzweigt oder stuft der Bootstrap ein, zählt er (`AGENTS.md` §3.13) |
| `internal/testpki` | update | Testhilfe: Serverzertifikat mit Gültigkeit und Schlüssel als PEM, Schlüsselvarianten; erzeugt im Testlauf, nichts eingecheckt |
| `test/integration` | update | Happy/Boundary/Negative nach LH-FA-23 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`, `slice-v1-abschluss-sessions` ist `done`, und [ADR-0018](../../adr/0018-tls-zum-client.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Handshake und Abweisung der Klartext-Verbindungen sprengen den Slice — zurück zur Zerlegung. **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 2300** geht der Slice zurück. Schnitt dann: (a) Optionen, Laden und Prüfen der Dateien, `PGR-E2001`/`PGR-E2007` im CLI-Adapter; (b) Aushandlung, Abweisung und `PGR-W3002` im PGWire-Adapter. **Schätzung** (Architect, 2026-10-13, in dieser Metrik): Produktcode rund 330 Zeilen (PGWire-Adapter 130, CLI-Adapter 170, Bootstrap und Katalog 30); Tests nach dem Faktor 2 des Produktcodes rund 660, nach dem Stand des Vorgängers mit gleich vielen Zusagen (`slice-v1-abschluss-einspielen-tls`: Tests rund 5-fach, Diff 2443, in einer Sitzung geprüft) bis rund 1650; mit der Review-Runde (15 bis 25 %) **1400 bis 1800**, die Marke liegt darüber.
- `in-progress` → `open`: Die TLS-Bibliothek verlangt eine Abhängigkeit außerhalb des PGWire-Adapters, oder die Prüfung von Gültigkeit und Passung gelingt im CLI-Adapter nicht ohne `crypto/tls` — Carveout beziehungsweise Entscheidung des Nutzers über die Regel des Architektur-Gates (`AGENTS.md` §3.6).

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

- Test-Zertifikate müssen im Testlauf erzeugt werden; ein fest eingecheckter Schlüssel wäre ein Geheimnis im Repository — **Ausgang:** offen bis Closure.

**Randformen** (`AGENTS.md` §3.12), entschieden vom Architect am 2026-10-13 in `LH-FA-23.a` (Technik-Stratum), `LH-FA-17.a` (Stelle des Ladens in der Reihenfolge des Starts); [ADR-0018](../../adr/0018-tls-zum-client.md) trägt nur die Entscheidung. Je Randform ein Fall im Test und, nach `AGENTS.md` §3.10, eine Mutation, die nur ihn bricht; bei Aufzählungen je Element:

- *Variable beim eigenen und beim anderen Kommando* (Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`, §1): beim eigenen Kommando liest der Leser die Variable, beim anderen (`play`) bleibt sie `PGR-E2001` (`LH-FA-17.a`; Test: `TestLeserUmgebungNichtGelesenerOptionen` mit gekürztem `nichtGelesen()`, `TestLeserUmgebungAndererKommandos`). Eine Variable mit ungültigem Wert (`PGWIRE_RECORDER_ALLOW_PLAINTEXT=ja`) ist `PGR-E2001`, auch wenn die Kommandozeile die Option setzt; eine leere Variable gilt als nicht gesetzt.
- *Quelle und Schlüssel der Datei*: `tls_cert`, `tls_key`, `allow_plaintext` in `record:` und `replay:`; in `play:` `PGR-E2004`; Quellen dürfen gemischt werden (Fall: Zertifikat aus der Umgebung, Schlüssel aus der Datei).
- *Fehlender, leerer, unbekannter Wert*: `--tls-key` allein, `--tls-cert` allein, `--allow-plaintext=true` ohne `--tls-cert` je `PGR-E2001`; `--allow-plaintext=false` ohne `--tls-cert` gültig; leerer Pfad ungültiger Wert; bei mehreren Verstößen nennt die Meldung die erste Option in der Reihenfolge `--tls-cert`, `--tls-key`, `--allow-plaintext`.
- *Aufrufart mehrfach*: letzte Angabe gilt (`LH-FA-17.a`); die Kombination wird nach der Zusammenführung geprüft.
- *Alias*: keiner (nur die drei Optionen, ihre Variablen und Schlüssel).
- *Stelle in der Reihenfolge des Starts*: Kombinationen vor jedem Fehler von `--upstream` (Übernahme aus `slice-v1-abschluss-upstream-verbinden`, §1); Laden von Zertifikat, dann Schlüssel, als letzte Prüfung der Kommandozeile und der Konfiguration, vor Ausgabedatei, Aufzeichnung und Listen (Fall: ungültiges Zertifikat und vorhandene Ausgabedatei meldet `PGR-E2007`); `config show` liest die Dateien nicht.
- *Form der Dateien*: reguläre Datei (FIFO, Verzeichnis nicht lesbar); Zertifikatsdatei nur `CERTIFICATE`-Blöcke (Schlüsselblock darin ungültig), erster Block der Server, Ablauf nur dort, Grenzen eingeschlossen, gegen die Systemuhr beim Start; Schlüsseldatei genau ein Block der drei Typen, `ENCRYPTED PRIVATE KEY` und `Proc-Type: 4,ENCRYPTED` verschlüsselt, anderer Typ und zweiter Block ungültig; Passung über den öffentlichen Schlüssel; eine Datei für beide nicht unterstützt; Namen und Verwendung ungeprüft; Meldung nennt Option und Ursache, nie Pfad oder Inhalt.
- *Handshake*: Mindestversion 1.2 (Client mit höchstens 1.1 scheitert, `PGR-W3002`); Bytes nach dem `SSLRequest` vor dem `S`, Abbruch des Clients nach dem `S` und Klartext nach dem `S` je `PGR-W3002` ohne Wirkung auf den Exit-Code; kein Client-Zertifikat angefordert; kein Zeitlimit; `PGR-W3002` trägt `code` und `remote`, nie Zertifikat oder Schlüssel.
- *Abbruch und Herunterfahren*: Herunterfahren und Zwangsende während der Aushandlung beenden sie ohne `PGR-W3002`; eine offene verschlüsselte Session endet bei Herunterfahren und Zwangsende wie eine unverschlüsselte (Fälle bei `record` und `replay`; das Zwangsende wartet nicht auf den Client, Fall: Client liest nicht mehr); Ende ohne `close_notify` gleich dem Ende einer unverschlüsselten Verbindung, gleich aufgezeichnet.
- *Rangfolge am Verbindungsaufbau*: `PGR-E6003` nur für eine `StartupMessage` 3.0; andere Version `PGR-E6002`, Sonderanfrage `PGR-E6001`, keine Startnachricht `PGR-W3003`; Klartext-`CancelRequest` `PGR-W3001`; auf der verschlüsselten Verbindung ein weiterer `SSLRequest` und ein `GSSENCRequest` je `PGR-E6001`; direktes TLS (ClientHello als erste Bytes) `PGR-W3003`, auch mit `--allow-plaintext` (Grenze, kein Test über den Fall hinaus).
- *Gleichzeitigkeit*: ein hängender Handshake hält eine zweite Verbindung nicht auf (Fall mit zwei Verbindungen).
- *Verhalten je Serverversion*: keines; TLS betrifft nur die Seite des Clients. **Akzeptiertes Negativ:** Die Tests laufen gegen die gepinnte Version des Integrationstests, nicht je Version; Grund: `replay` hat keinen Server, `record` gibt die Bytes des Servers unverändert weiter, und die Matrix liegt bei `slice-v1-abschluss-postgres-versionen` (keine Bindung dorthin, die Szenarien dieses Slice gehören nicht in deren DoD).
- *Geheimnis*: der Schlüssel erscheint in keiner Log-Zeile (auch `debug`), in keiner Meldung, nicht in `config show` und in keiner Textform des `Command` (Fall: `%v` und `%+v`).

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-tls-client-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 3. Schichten: zwei Schichten (PGWire-Adapter, CLI-Adapter).

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 2 (Leser der Option und Kürzung von `nichtGelesen()` im Test des CLI-Adapters) und fügt keinen Liefer-Punkt hinzu: Liefer-Punkte 3. Schichten: zwei (PGWire-Adapter, CLI-Adapter); der CLI-Adapter war schon gezählt, `fremd_test.go` liegt in ihm.

Nachgezählt bei der Prüfung des Architect (2026-10-13, `AGENTS.md` §3.13): Die Randformen aus §6 gehören zu den drei Liefer-Punkten (Aushandlung und Abweisung; Optionen und Kombinationen; Dateien und Fehler) und fügen keinen hinzu. Schichten: zwei. Die Meldungscodes in `internal/hexagon/model` sind Konstanten des Katalogs und zählen nicht als Schicht (wie `PGR-E2007` im Vorgänger `slice-v1-abschluss-einspielen-tls`, der mit CLI- und Upstream-Adapter bei zwei blieb); der Bootstrap ist Verdrahtung (§3). Zählte der Katalog doch, läge der Slice bei drei Schichten, und der Schnitt aus §4 greift. `internal/testpki` ist Testhilfe.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
