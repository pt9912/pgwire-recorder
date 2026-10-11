# Slice slice-v1-abschluss-einspielen-tls: TLS zum Server beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md), [ADR-0014](../../adr/0014-konfigurationsdatei.md), [ADR-0016](../../adr/0016-einspielen-anmeldung-und-tls.md), [ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md)

**Berührte Spec-Stellen:** `LH-FA-20.a` · `LH-FA-17.a` · `SPEC-017` · `SPEC-022` · `SPEC-033` · `SPEC-034` · `ARC-005` · `ARC-007` · `ARC-009`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `pgwire-recorder play` verbindet sich auf Wunsch (`--upstream-tls`, `sslmode=require` der benutzten Verbindung) mit TLS zum Server, prüft dessen Zertifikat, auch gegen eine eigene Zertifizierungsstelle (`--upstream-ca`), und meldet jeden Fehlschlag von TLS und Zertifikat als `PGR-E4005`, eine nicht lesbare CA-Datei beim Start als `PGR-E2007`.

**Übernimmt:** `slice-v1-abschluss-einspielen` — dessen Teil *TLS* (Marke [T] in dessen §6)
nach dem Schnitt vom 2026-10-09 (Prüfung des Architect, dort §6 Risiko *Größe*; Entscheidung
des Nutzers; dort §1, *Abgegeben*). Im Einzelnen:

- die Optionen `--upstream-tls` und `--upstream-ca` am allgemeinen Leser, in der Reihenfolge
  der Optionstabelle, mit ihren Umgebungsvariablen und Schlüsseln im Abschnitt `play:`
  (`LH-FA-17.a`);
- TLS zum Server nach `sslmode` der benutzten Verbindung, wenn `--upstream-tls` nicht gesetzt
  ist (`LH-FA-17.a` *Wirkung einer URL*); der Punkt kam nach `slice-v1-abschluss-einspielen`
  aus `slice-v1-abschluss-upstream-verbinden` (dort §1, Abgrenzung);
- das Lesen der Datei aus `--upstream-ca` beim Start (`PGR-E2007`) und `PGR-E4005` von TLS
  und Zertifikat (`LH-FA-20.a` *Start*, *TLS*; [ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md));
- im Benutzerhandbuch §5 *Konfigurationsdatei* der TLS-Teil der Wirkung einer Verbindung bei
  `play` aus Befund F-533 (Review von `slice-v1-abschluss-upstream-verbinden`, über
  `slice-v1-abschluss-einspielen`): `sslmode=require` mit TLS und Prüfung des Zertifikats, ein
  gesetztes `--upstream-tls` vor `sslmode`; dieser Handbuch-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-einspielen-tls-doku`;
- die Randformen [T] aus §6 jenes Slice (§6 unten). Dieser Slice ist auch die Adresse der
  Abgrenzung *TLS* von `slice-v1-abschluss-antwortvergleich` (dort §1).

**Aufsetzen.** Der Slice setzt auf dem Kern `slice-v1-abschluss-einspielen` auf: Dort sind
das Kommando `play`, die übrigen Optionen, der Verbindungsaufbau ohne TLS und die Einstufung
einer Fehlerantwort im Aufbau geliefert. Bis zu diesem Slice sind `--upstream-tls` und
`--upstream-ca` bei `play` unbekannt (`PGR-E2001`), ihre Schlüssel im Abschnitt `play:`
unbekannt (`PGR-E2004`), ihre Umgebungsvariablen unbeachtet, und `sslmode=require` der
benutzten Verbindung ist bei `play` wie bei `record` `PGR-E2004` (Zwischenstand in §6 jenes
Slice); diesen Zwischenstand ersetzt dieser Slice.

**Aus `slice-v1-abschluss-einspielen-laufsteuerung`** (Review F-568 jenes Slice): die Bindung an
den Vertrag der Ports `Einspielziel` und `EinspielSession`, nach dem jeder Fehler seinen
Meldungscode trägt; auf ihr ruht das akzeptierte Negativ (e) jenes Slice, und jeder neue Fehler
von TLS und Zertifikat im Aufbau hält sie (§6, *Bindung an den Vertrag der Ports*).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-einspielen-tls-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Die Anmeldung mit Passwort — `slice-v1-abschluss-einspielen-anmeldung` (dort §1,
  *Übernimmt*); beide hängen nur am Kern, nicht aneinander. Getestet wird hier gegen einen
  Server ohne Passwort-Anmeldung.
- TLS zum Upstream bei `record` — nicht Teil des Produkts in dieser Welle
  ([welle-v1-abschluss](../welle-v1-abschluss.md) §6); `sslmode=require` bleibt bei `record`
  `PGR-E2004`.
- TLS zum Client bei `record` und `replay` — `slice-v1-abschluss-tls-client`; das ist die
  andere Seite der Verbindung und ein anderer Vorgang.
- Client-Zertifikate, weitere Werte von `sslmode` und ein Überspringen der Prüfung —
  außerhalb des Funktionsumfangs: `LH-FA-20` schließt das Überspringen aus, `LH-FA-17.a`
  kennt nur `disable` und `require`.
- Code im Kern (Play-Service, Ports) und im PGWire-Adapter — Schicht-Abgrenzung: TLS liegt im
  Upstream-Adapter, Optionen und CA-Datei liest der CLI-Adapter, der Bootstrap reicht sie
  weiter.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): Mit `--upstream-tls` oder mit `sslmode=require` der benutzten Verbindung
      baut `play` die Verbindung mit TLS auf und prüft das Serverzertifikat gegen den
      Zertifikatsspeicher des Systems, ergänzt um die Zertifikate aus `--upstream-ca`, und
      den Namen gegen den eingesetzten Host, eine IPv6-Adresse ohne ihre Zone gegen die
      IP-Adressen des Zertifikats; ein gesetztes `--upstream-tls`, auch ausdrücklich
      `false`, geht `sslmode=require` vor, aus jeder Quelle; ohne beides baut `play` keine
      TLS-Verbindung auf (Integrationstest gegen einen Server mit dem Zertifikat einer
      eigenen Zertifizierungsstelle).
- [ ] [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler): `N` auf das `SSLRequest`, ein Fehler der Aushandlung, ein abgelaufenes,
      ungültiges oder auf einen anderen Namen ausgestelltes Zertifikat und ein Server, der
      eine unverschlüsselte Verbindung ablehnt, sind `PGR-E4005`; ein anderes Byte, ein
      Verbindungsende vor der Antwort und Bytes nach `S` vor der Aushandlung sind
      `PGR-E4002`; `--upstream-ca` ohne TLS ist `PGR-E2001`, nach `--upstream` und vor den
      Variablen der Platzhalter; eine CA-Datei, die keine reguläre Datei ist (FIFO und
      Verzeichnis ohne Warten), keinen PEM-Block, einen Block eines anderen Typs oder kein
      lesbares X.509-Zertifikat enthält, ist beim Start vor dem Laden der Aufzeichnung
      `PGR-E2007`, die Meldung nennt weder Pfad noch Inhalt ([`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)) (Test). Beleg in
      §7 für Punkt 1 und 2: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10). Die Abdeckungstabellen sind über `make abdeckung` nachgezogen.
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
| `internal/adapters/driven/postgres` | update | `SSLRequest`, Aushandlung und Prüfung des Zertifikats im Aufbau von `play` (`einspielen_tls.go`, `Einspielziel.TLS` und `.CA`); Einstufung `PGR-E4005` und `PGR-E4002` |
| `internal/adapters/driving/cli` | update | Optionen `--upstream-tls` und `--upstream-ca` mit Umgebung und Schlüsseln im Abschnitt `play:`; „gesetzt“ vom Standardwert unterschieden (`gelesen.gesetzt`, `quellen`); TLS nach `sslmode`; `--upstream-ca` ohne TLS `PGR-E2001`; die Datei aus `--upstream-ca` als letzte Prüfung des Starts lesen und einstufen (`zertifizierungsstelle.go`, `PGR-E2007`), die Zertifikate als `Zertifikate` (Typ über `[]*x509.Certificate`, mit `Format`) in den Optionen von `play` ablegen; `PlayOptions` ist damit nicht mehr mit `==` vergleichbar |
| `internal/bootstrap` | update | Verdrahtung ohne Logik: `UpstreamTLS` und die Zertifikate aus den Optionen in den `Einspielziel` des Upstream-Adapters setzen (`postgres.Zertifikate(o.UpstreamCA)`); die Datei liest der CLI-Adapter (§6, *Ort der CA-Datei*) |
| `internal/hexagon/model` | update | die Kennung `PGR-E2007` als Konstante neben den übrigen Codes; keine Logik, zählt wie die Verdrahtung nicht als Schicht |
| `test/integration` | update | Happy/Negative nach LH-FA-20 über einen TLS-Proxy vor der Instanz: Zertifikat der eigenen Zertifizierungsstelle, abgelaufenes Zertifikat, falscher Name, Server, der unverschlüsselte Verbindungen ablehnt, Anmeldung über TLS |
| `internal/testpki`, `internal/bootstrap/tlsproxy` | neu (Testhilfe) | Zertifikate zur Testzeit mit `crypto/x509`, nichts eingecheckt (`testpki`); der TLS-Proxy für die Integrationstests (`tlsproxy`). Der Proxy liegt unter `internal/bootstrap`, weil `.a-check.yml` `crypto/tls` nur in den beiden PGWire-Adaptern und der Verdrahtung zulässt und `test/integration` ihn nicht importieren darf; `testpki` kommt ohne `crypto/tls` aus. Nur Testcode darf beide importieren (geprüft von einem Test, siehe Aufträge unten); sie gehören zu keiner Schicht und zählen nicht |

**Aufträge aus der Review vom 2026-10-11 (Architect)** — kein neuer Liefer-Punkt, sie gehören zu DoD-Punkt 2 (Zusagen mit Mutation, `AGENTS.md` §3.10):

- **F-602, Import der Testhilfen.** Ein Test im Paket `internal/bootstrap` (Unit-Test, Test-Image; kein neues Gate, keine ADR) ruft `go list -deps` über `./cmd/...` auf und hält die Menge gegen die Verbotsliste `testing`, `internal/testpki`, `internal/bootstrap/tlsproxy`. Mutation: Blank-Import von `tlsproxy` in `bootstrap.go`, der Test wird rot. Umgesetzt als `TestBinaryOhneTesthilfen`; `go list` läuft im Test-Image ohne Netz (`GOPROXY=off`), der Rückfall auf die engere Zusage war nicht nötig. Die Zusage steht in beiden `doc.go` nur in dem Maß, das der Test prüft: kein Paket unter `./cmd/...` importiert die Testhilfen; ob eine Testdatei außerhalb von `./cmd/...` sie importiert, prüft nichts. Die Gate-Alternative (Regel in `.a-check.yml` für `internal/bootstrap/**`) ist eine Gate-Änderung (`AGENTS.md` §3.6) und braucht eine ADR; sie ist hier nicht beauftragt.
- **F-603, Standardpfad des Zertifikatsspeichers.** Ein Integrationstest startet das Binary mit `SSL_CERT_FILE` auf die CA-Datei der Testzeit im Prozessumfeld des Unterprozesses, `sslmode=require` ohne `--upstream-ca`, gegen den TLS-Proxy: Verbindung gelingt. Gegenprobe im selben Test: ohne `SSL_CERT_FILE` `PGR-E4005`. Mutation: der Standardpfad ersetzt `x509.SystemCertPool` durch einen leeren Speicher (Mutant M14 der Review), der Test wird rot. Umgesetzt als `TestE2EPlayTLSZertifikatsspeicherDesSystems`. Beleg in §7.
- **F-604, Untergrenze der TLS-Version.** Keine Zusage (§6): Der Fall in `einspielen_tls_test.go` heißt „S, dann Gegenstelle nur TLS 1.0, Aushandlung scheitert“ und prüft, was er nennt (die Gegenstelle spricht nur TLS 1.0 und die Aushandlung scheitert), ohne Hinweis auf eine Untergrenze des Clients; kein Code.
- **F-606, `Hello` des Proxys.** Der hängende Proxy schließt `Hello` einmal (`sync.Once`), auch bei einer zweiten Verbindung; `TestHaengenderProxyZweiteVerbindung` im Paket `tlsproxy`.
- **F-601, F-607** sind in `spec/architecture.md` §6 und `LH-FA-20.a` *TLS* berichtigt; F-605, F-608 zur Kenntnis.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen` liegt in `done/`
(Kommando `play`, Verbindungsaufbau ohne TLS, Einstufung der Fehler im Aufbau). Schritt 14
der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt nach
`slice-v1-abschluss-einspielen-anmeldung`, ohne von ihm abzuhängen (Schnitt vom 2026-10-09).
Die Randformen aus §6 entschied der Architect am 2026-10-09 vor dem Code in `LH-FA-20.a`
und `LH-FA-17.a`; vor dem ersten Code-Commit prüft er die Liste gegen den gelieferten Kern und
legt fest, in welchem Paket die CA-Datei gelesen und geprüft wird, gegen die Regeln von
`make a-check` (`AGENTS.md` §3.12). Geprüft und entschieden am 2026-10-11 (§6, *Ort der
CA-Datei*): der CLI-Adapter, zwei Schichten.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar. Schnitt dann: TLS mit dem Zertifikatsspeicher des Systems hier,
  `--upstream-ca` mit `PGR-E2007` als eigener Slice; der erste ist allein lieferbar, weil
  `--upstream-ca` nur ergänzt.
- `in-progress` → `open` (blockiert — Carveout?): Die Zertifikate der Tests lassen sich im
  Testgeschirr nicht ohne Netz erzeugen oder nicht mit einem gepinnten Werkzeug; dann zuerst
  die Entscheidung über ihre Herkunft.

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

**Randformen** (`AGENTS.md` §3.12) — die mit Marke [T] aus §6 von
`slice-v1-abschluss-einspielen`, geprüft vom Architect am 2026-10-09 vor dem ersten
Code-Commit jenes Slice. Neu entschieden heißt: im Commit jener Prüfung in `LH-FA-20.a`, sonst
an der genannten Stelle. Offen ist keine.

- **Optionen** [T] — `--upstream-tls` und `--upstream-ca` am allgemeinen Leser, in der
  Reihenfolge der Optionstabelle, Abschnitt `play:`; bestätigt, `LH-FA-17.a`.
- **`--upstream-tls` ausdrücklich `false`** [T] — geht `sslmode=require` vor, aus jeder Quelle
  (Option, Umgebungsvariable, Schlüssel); neu entschieden in `LH-FA-17.a` *Wirkung einer URL*.
  *Hinweis an den Implementer:* Der Leser muss „gesetzt“ vom Standardwert unterscheiden
  (`gelesen.cliOk`, `envOk`, `datei.wert`).
- **`--upstream-ca` ohne TLS** [T] — `PGR-E2001`, in der Reihenfolge nach `--upstream` und vor
  den Variablen der Platzhalter, weil TLS vom `sslmode` der benutzten Verbindung abhängt; neu
  entschieden in `LH-FA-17.a` *Fehler*. Die Meldung nennt die Option, nicht den Pfad.
- **Datei aus `--upstream-ca`** [T] — gelesen nach den Prüfungen von `LH-FA-17.a`, vor der
  Aufzeichnung; nur eine reguläre Datei (Links gefolgt), FIFO, Verzeichnis, Gerät, Socket nicht
  lesbar, ohne Warten; mindestens ein PEM-Block, jeder `CERTIFICATE` mit lesbarem X.509, Text
  außerhalb unbeachtet; sonst `PGR-E2007`, Meldung ohne Pfad und Inhalt; neu entschieden in
  `LH-FA-20.a` *Start*. Ablauf eines Zertifikats erst beim Aufbau (`PGR-E4005`), bestätigt.
- **Antwort auf `SSLRequest`** [T] — `S` Aushandlung, `N` `PGR-E4005`, anderes Byte oder Ende
  davor `PGR-E4002`, Bytes nach `S` vor der Aushandlung `PGR-E4002`; jeder Fehler der
  Aushandlung `PGR-E4005`; neu entschieden in `LH-FA-20.a` *TLS*.
- **Name im Zertifikat** [T] — gegen den eingesetzten Host, IPv6 ohne Zone gegen die
  IP-Adressen; neu entschieden in `LH-FA-20.a` *TLS*.
- **Zertifikatsspeicher des Systems nicht ladbar** [T] — gilt als leer (Grenze); neu
  entschieden in `LH-FA-20.a` *TLS*.
- **Abbruch im Aufbau mit TLS** [T] (aus `slice-v1-abschluss-einspielen`, §6,
  Randform-Rückgabe R3, entschieden vom Architect am 2026-10-09 vor dem ersten Code-Commit
  jenes Slice) — nach einem Fehler im Aufbau keine weitere Nachricht, kein `Terminate`; ein
  Alarm der TLS-Schicht (auch `close_notify` oder ein Alarm der gescheiterten Aushandlung)
  ist keine Nachricht; neu entschieden in `LH-FA-20.a` *Abbruch im Aufbau*. Nachgezählt beim
  Eintragen (`AGENTS.md` §3.13): drei Liefer-Punkte, zwei Schichten (CLI-Adapter,
  Upstream-Adapter); der Punkt liegt im Upstream-Adapter, ohne eigene Zusage in der DoD.
- **Ort der CA-Datei** (geprüft vom Architect am 2026-10-11, Nutzerentscheidung Option A) —
  der CLI-Adapter liest die Datei, stuft `PGR-E2007` ein und legt die Zertifikate als
  `cli.Zertifikate` (benannter Typ über `[]*x509.Certificate`, dessen `Format` kein Zertifikat ausgibt) in `cli.PlayOptions`; der Bootstrap setzt sie als `postgres.Zertifikate` in den `Einspielziel`,
  der Upstream-Adapter baut daraus mit dem Systemspeicher das `tls.Config`. `crypto/x509` und
  `encoding/pem` sind Standardbibliothek; `.a-check.yml` schränkt nur `crypto/tls` (auf die
  beiden PGWire-Adapter) ein, kein Gate ändert sich (`AGENTS.md` §3.6), keine ADR. Kein
  `x509`-Typ erreicht Model, Services oder Ports (der Core kennt keine Zertifikate, `spec/architecture.md`
  §4.4). Der Ort steht in `spec/architecture.md` §4.4. Die Reihenfolge stimmt mit
  `LH-FA-17.a` und `LH-FA-20.a` *Start*: `--upstream-ca` ohne TLS (`PGR-E2001`) im Hook von
  `--upstream` zwischen `pruefeUpstream` und `zielPlay`, die Datei danach als letzte Prüfung,
  vor dem Laden der Aufzeichnung (das der Bootstrap erst nach `Parse` macht); ein Fehler nennt
  die Option und die Ursache des Betriebssystems ohne Pfad (Muster `nichtLesbar`). Neu
  entschieden in `LH-FA-17.a` *Fehler* (Stelle der Datei).
- **Rückgabe des Implementers (geprüft vom Architect am 2026-10-11, nach dem Code)** — vier
  Punkte, die §6 nicht nannte, und ein Rand der Antwort auf `SSLRequest`:
  (1) *Ort der Testhilfen*: `internal/bootstrap/tlsproxy` und `internal/testpki` bleiben; die
  Composition Root ist von den Schichtregeln ausgenommen und darf `crypto/tls`; `go list -deps
  ./cmd/...` zeigt keines der beiden Pakete und kein `testing` im Binary; kein Gate ändert sich
  (`AGENTS.md` §3.6), keine ADR; Satz in `spec/architecture.md` beim Architektur-Gate.
  (2) *Typ der Zertifikate*: `cli.Zertifikate` und `postgres.Zertifikate` sind Typen der
  Adapter, die Umsetzung macht die Composition Root; akzeptiert, dass `PlayOptions` und
  `Command` nicht mehr mit `==` vergleichbar sind. (3) *Senden des `SSLRequest` scheitert*:
  `PGR-E4002`; die Aushandlung beginnt mit dem `S`; Satz in `LH-FA-20.a` *TLS*. (4)
  *Hook-Signatur im Leser* (`option.zuletzt` mit `quellen`, `lies` geteilt): bestätigt,
  CLI-intern, kein Vertrag. *Akzeptiertes Negativ:* Bytes nach `S` werden nur erkannt, wenn
  sie mit dem `S` eintreffen; später eintreffende stören die Aushandlung (`PGR-E4005`).
  Grund: Die Segmentierung von TCP ist im Test nicht steuerbar, und für den Client ist es
  dasselbe Ereignis; steht als Grenze in `LH-FA-20.a` *TLS*.
- **Host als IP-Adresse, Aufbau nach der Aushandlung, `config show`** — IPv4 wie IPv6 gegen
  die IP-Adressen des Zertifikats; nach der Aushandlung laufen Aufbau und Anmeldung wie ohne
  TLS auf der verschlüsselten Verbindung, ein Klartext-Passwort ist über TLS zulässig;
  `config show` liest die Datei nicht; neu entschieden in `LH-FA-20.a` *TLS*.
- **Abbruchsignal in der Aushandlung** — der Aufbau schließt die Aushandlung ein; es gilt
  *Abbruchsignal* (das erste lässt ihn zu Ende laufen, das zweite schließt ohne Fehler),
  bestätigt, keine neue Regel; geprüft mit einem Test, der in der Aushandlung das zweite
  Signal gibt.
- **Bindung an den Vertrag der Ports** (aus `slice-v1-abschluss-einspielen-laufsteuerung`, dort
  §6 akzeptiertes Negativ (e), Review F-568) — jeder neue Fehler im Aufbau mit TLS (Antwort auf
  `SSLRequest`, Aushandlung, Zertifikat, Ablehnung ohne TLS) verlässt den Upstream-Adapter mit
  seinem Meldungscode (`PGR-E4005` oder `PGR-E4002`), keiner ohne Code; sonst hinge er nach
  `SPEC-034` *Ausgabe* als Ursache an eine frühere Meldung, mit Exit-Code 4 statt 1. Entschieden
  vom Architect am 2026-10-10 dort: keine eigene Regel, der Ausschluss ruht auf dem Vertrag der
  Ports (Kommentar an `Verbinde`). Geprüft mit DoD-Punkt 2, der für jeden Fehlschlag den Code
  zusagt; kein neuer Liefer-Punkt. `PGR-E2007` entsteht beim Start vor jeder Interaktion und
  hat keinen früheren Fehler.
- **Server lehnt unverschlüsselte Verbindung ab** — `PGR-E4005`; bestätigt, Tabelle
  *Fehlerregeln beim Einspielen* in `LH-FA-20.a`; hier geprüft, weil der Fall einen Server mit
  TLS-Pflicht braucht.
- **Ablösung des Zwischenstands** — bis zu diesem Slice sind `--upstream-tls` und
  `--upstream-ca` bei `play` unbekannt, und `sslmode=require` ist `PGR-E2004` (§6 von
  `slice-v1-abschluss-einspielen`, *Zwischenstand*). Die Tests des Kerns, die das prüfen,
  ändert dieser Slice.

*Akzeptiertes Negativ der Prüfung vom 2026-10-09* (aus §6 von `slice-v1-abschluss-einspielen`):

- **TLS-Version und Verfahren** — die Voreinstellung der Standardbibliothek, keine Zusage
  (`LH-FA-20.a` *TLS*); der Client setzt keine Untergrenze, und kein Test sagt eine zu
  (Review F-604, 2026-10-11).

*Akzeptierte Negative der Prüfung vom 2026-10-11:*

- **PostgreSQL mit `ssl=on` im Integrationstest** — der gepinnte Server (`postgres:17-alpine`)
  enthält kein `openssl`, und das Runner-Skript erzeugte Zertifikate nur mit zusätzlichem
  Werkzeug; der Slice prüft Aushandlung und Zertifikat deshalb gegen einen **TLS-Proxy im
  Test** (Zertifikate zur Testzeit mit `crypto/x509` erzeugt, nichts eingecheckt): Er antwortet
  `S`, handelt aus und reicht die entschlüsselte Verbindung an den echten PostgreSQL weiter.
  Fehlerfälle (`N`, anderes Byte, abgelaufen, anderer Name, Server lehnt unverschlüsselt ab)
  laufen gegen eine Test-Gegenstelle im Paket des Upstream-Adapters. Nicht geprüft: das
  Zusammenspiel mit dem TLS-Server von OpenSSL; Grund: Client ist die Standardbibliothek, die
  Gegenstelle ein Standard-TLS.
- **Zertifikatsspeicher des Systems** — Das Zertifikat des Tests liegt nicht im Speicher des
  Test-Containers. Der Speicher wird im Test des Upstream-Adapters über ein austauschbares Feld
  geprüft (gefüllter Speicher nimmt, leerer und nicht ladbarer Speicher lehnen ab, nur
  `--upstream-ca` nimmt); der **Standardpfad** (`x509.SystemCertPool`) wird zusätzlich Ende zu
  Ende geprüft (§3, Auftrag F-603 der Review vom 2026-10-11): Das Binary läuft als frischer
  Prozess, `SSL_CERT_FILE` wirkt dort bei der ersten Ladung. Die Grenze des Images steht im
  zweiten Risiko unten.
- **SNI, Punkt am Ende des Hosts, Größe der CA-Datei, Zertifikat mit Kopfzeilen im PEM-Block**
  — Verhalten der Standardbibliothek, keine Zusage; kein Fall des Produkts.

**Risiken:**

- Das Testgschirr braucht Zertifikate einer eigenen Zertifizierungsstelle, ein abgelaufenes
  und eines auf einen anderen Namen — **Ausgang:** entschieden am 2026-10-11: zur Testzeit im
  Test erzeugt, nicht eingecheckt (siehe Negative oben); offen bis Closure, ob das Geschirr
  dadurch unvertretbar wächst.
- Im Produkt-Image kann der Zertifikatsspeicher des Systems leer sein; dann prüft `play` nur
  gegen `--upstream-ca` (Grenze in `LH-FA-20.a` *TLS*) — **Ausgang:** offen bis Closure.

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

### Belege des Implementers

**Stand und Läufe.** Alles Docker-only. `make gates` lief grün (Exit 0) auf dem Inhalt von
Commit `147afea` (Produktcode, Tests, Plan §3 und §6, Abdeckungstabellen) und erneut auf dem
Inhalt von Commit `058c677` (Nacharbeit zur Review vom 2026-10-11, F-602 bis F-604 und F-606,
Plan §3, Abdeckungstabellen); in der Kette: `make test` (gofmt, vet, Unit-Tests),
`make test-integration` (mit den `TestE2EPlayTLS*`, `TestE2EPlayKlartextAbgelehnt`,
`TestE2EPlayUpstreamCAFehler`, seit `058c677` auch
`TestE2EPlayTLSZertifikatsspeicherDesSystems`), `make lint` (0 issues), `make a-check`
(0 Befunde), `make a-check-negativ` und `make lint-gegenprobe` grün. `make abdeckung` ist
nachgezogen. Der Stand dieses Absatzes steht in einem eigenen Commit nach `058c677`, der nur
§7 ändert; danach lief `make gates` ein weiteres Mal auf dem sauberen Baum, damit die Images
den Baum zeigen (dieser Lauf steht hier nicht, er würde einen weiteren Commit verlangen).

**Größe.** `3d0bd0c..147afea`: 27 Dateien, +2288 −130 Zeilen. Davon Produktcode +349 −55
(Upstream-Adapter 96, CLI-Adapter ohne Tests 209 netto, Bootstrap 6, Model 1), Tests und
Testhilfen etwa +1850 (Unit 1100, Bootstrap 132, Integration 359 −43, Testhilfen `testpki` und
`tlsproxy` 246), Abdeckungstabellen +63, Plan 6. Das liegt rund 14 % über den ~2000 Zeilen
mit Tests, die eine Review-Sitzung tragen soll; ich habe nicht angehalten, weil der Produktcode
bei 350 Zeilen liegt und der Überhang in den Tests je Zusage und Quelle steht (§3.10).

**Mutationen (`AGENTS.md` §3.10).** Weg: je Mutant eine frische Kopie des Arbeitsbaums
(`cp -r` ohne `-p`, ohne `.git`), die Datei ersetzt, `gofmt -w` auf ihr, `gofmt -l` leer; die
Unit-Mutanten laufen mit `go test -run` im gepinnten Test-Image mit Bind-Mount der Kopie (kein
Build-Kontext, daher kein mtime-Fall), die Integrations-Mutanten mit `make test-integration`
in der Kopie. Jeder Mutant der Tabelle ist **einmal rot gesehen**, an dem hier genannten Test und
aus dem Grund der Zusage (Meldung gelesen). Bis `147afea`: 62 Unit-Läufe, 8 Integrations-Läufe,
alle rot. **Berichtigt:** Die frühere Aussage „kein grüner Mutant“ stimmte nicht. Der Mutant
M14 der Review (Standardpfad `x509.SystemCertPool` durch einen leeren Speicher) war bis
F-603 grün; er ist seit `058c677` rot (Zeile *Standardpfad*). Nacharbeit zur Review: fünf
Unit-Läufe (drei rot zu F-602, einer rot zu F-606, einer grün zu F-604) und zwei
Integrations-Läufe (beide rot zu F-603); insgesamt 67 Unit- und 10 Integrations-Läufe. Die
Kopien sind entfernt; danach lief `make gates` neu, damit die Images den Baum zeigen.

**Ein Mutant bleibt grün (F-604).** `MinVersion: tls.VersionTLS10` im Client
(`tlsKonfiguration`) lässt `go test ./internal/adapters/driven/postgres` grün. Er ändert das
Verhalten (der Client nähme einen Server mit TLS 1.0 an), ist also über die Schnittstelle
fangbar; fangen würde ihn nur ein Server mit `MinVersion` und `MaxVersion` TLS 1.0. Dafür gibt
es hier keinen Test und keine Test-Idee mit Adresse, weil die Untergrenze keine Zusage ist (§6,
*TLS-Version und Verfahren*): Wer sie zusagt, schreibt diesen Test mit der Zusage. Der Fall
„S, dann Gegenstelle nur TLS 1.0, Aushandlung scheitert“ prüft nur, dass eine Aushandlung, die
scheitert (hier an der Untergrenze der Gegenstelle), PGR-E4005 ist.

**Ein zweiter Mutant bleibt grün (V-167).** Der Mutant „Fehler des Ladens nimmt die
Systemwurzeln“ (im Fehlerzweig `pool, _ = x509.SystemCertPool()`, danach die CA ergänzt) lässt
alle Pakete grün (Upstream-Adapter, CLI-Adapter, Bootstrap; die Verifikation fuhr ihn als V32).
Eine frühere Fassung dieses Abschnitts sagte, `…/nicht_ladbarer_Speicher,_in_CA` werde rot und
*leer* und *nicht ladbar* würden erst mit CA getrennt; das stimmte nicht. Der Grund: Die
Systemwurzeln des Test-Containers enthalten die Test-CA nicht, und mit CA gelingt der Aufbau in
beiden Fassungen. Der Unterschied *leer* gegen *Systemwurzeln* ist nur mit einem Serverzertifikat
einer vertrauten Wurzel zu beobachten, und das kann die Testumgebung nicht erzeugen. Der Mutant
gilt deshalb als in der Testumgebung äquivalent; die Zusage steht in `LH-FA-20.a` *TLS* als
„(Grenze)“. Was der Hook prüft, ist *Fehler und `nil` gelten als leer*, und dort ist der Mutant
`nil` ohne Fehler rot.

| Zusage | Mutation | roter Test |
|---|---|---|
| Mit TLS: SSLRequest, Aushandlung, Startup nur verschlüsselt | TLS-Zweig im Aufbau aus; SSLRequest-Nummer um 1 verändert; Startup nach der Aushandlung auf der Klartext-Verbindung | `TestEinspielTLS` (alle drei); der dritte auch `TestEinspielTLSKlartextPasswort`, `TestEinspielTLSAbbruchImAufbau` |
| Ohne TLS kein SSLRequest, auch mit gesetztem CA | TLS-Zweig immer an; Bootstrap setzt `TLS: true` | `TestEinspielOhneTLS`; `TestRunPlayTLS/ohne_TLS`, `TestRunPlayPasswort`, `TestRunPlay` |
| Das Serverzertifikat wird geprüft, ein Überspringen gibt es nicht | `InsecureSkipVerify` | `TestEinspielTLSZertifikatFehler` (alle 6 Zeilen), `TestEinspielTLSName`, `TestEinspielTLSZertifikatsspeicher`; Ende zu Ende `TestE2EPlayTLSZertifikatFehler` (alle 5 Zeilen) |
| Der Name wird geprüft | Kette geprüft, Name nicht | `TestEinspielTLSZertifikatFehler/andere_IP-Adresse`, `/nur_DNS-Name_localhost`, `/IP-Adresse_als_DNS-Name`, `TestEinspielTLSName/Name_gegen_anderen_DNS-Namen`, `/Name_gegen_IP-Adresse`; Ende zu Ende `TestE2EPlayTLSZertifikatFehler/andere_IP-Adresse`, `/Name_statt_IP-Adresse` |
| … gegen den eingesetzten Host | Name fest `127.0.0.1`; Name mit Port | `TestEinspielTLSName/Name_gegen_DNS-Namen`, `/Name_gegen_IP-Adresse`; `TestEinspielTLS`, `TestEinspielTLSName`, `TestEinspielTLSIPv6MitZone` |
| … eine IP-Adresse nur gegen IP-Adressen, nie gegen DNS-Namen | Namen einer IP-Adresse mit Punkt am Ende (dann DNS-Pfad) | `TestEinspielTLSZertifikatFehler/IP-Adresse_als_DNS-Name`, `TestEinspielTLSName/IPv4_gegen_IP-Adresse`, `TestEinspielTLSServerName` |
| … IPv6 ohne Zone | **kein eigener Mutant fahrbar**: Kein Code des Produkts streift die Zone ab; die Standardbibliothek (Go 1.27, `VerifyHostname` über `netip`) gleicht eine IPv6-Adresse mit und ohne Zone gegen die IP-Adressen ab | `TestEinspielTLSServerName` hält das Verhalten der Standardbibliothek fest, `TestEinspielTLSIPv6MitZone` fährt `[::1%lo]` Ende zu Ende (endet ohne Prüfung, wenn die Schleife kein IPv6 hat; hier lief er) |
| Zertifikatsspeicher des Systems **und** CA, nicht eines statt des anderen | CA ersetzt den Speicher; CA nicht aufgenommen | `TestEinspielTLSZertifikatsspeicher/im_Speicher,_andere_in_CA` und `/im_Speicher_neben_anderen,_eine_eigene`; `/nur_in_CA`, `/andere_im_Speicher,_in_CA`, `TestEinspielTLS` |
| Speicher nicht ladbar gilt als leer | Speicher `nil` ohne Fehler; Fehler des Ladens nimmt die Systemwurzeln: **grün**, siehe *Ein zweiter Mutant bleibt grün (V-167)* | `TestEinspielTLSZertifikatsspeicher/Speicher_nil_ohne_Fehler,_in_CA` (als Absturz); für den Mutanten der Systemwurzeln **kein roter Test** |
| `S` mit Bytes im selben Lesen: PGR-E4002 | Prüfung der Bytes nach `S` entfernt | `TestEinspielTLSAntwort/S_mit_einem_Byte_dahinter`, `/S_mit_TLS-Record_dahinter`, `TestEinspielTLSSendenUndLesen` |
| `N`: PGR-E4005 | `N` als PGR-E4002 | `TestEinspielTLSAntwort/N`, `/N_mit_Bytes_dahinter` |
| Anderes Byte: PGR-E4002 | als PGR-E4005 | `TestEinspielTLSAntwort/E`, `/Fehlerpaket`, `/x`, `/s`, `/Byte_eines_TLS-Records`, `/Byte_0` |
| Ende vor der Antwort: PGR-E4002 | als PGR-E4005 | `TestEinspielTLSAntwort/Ende_davor`, `TestEinspielTLSSendenUndLesen` |
| Senden des SSLRequest scheitert: PGR-E4002 | als PGR-E4005 | `TestEinspielTLSSendenUndLesen` |
| Fehler der Aushandlung, auch Ende darin, Nicht-TLS, gescheiterte Aushandlung gegen eine Gegenstelle mit nur TLS 1.0: PGR-E4005 | Aushandlung als PGR-E4002 | `TestEinspielTLSAntwort/S,_dann_Ende`, `/S,_dann_kein_TLS`, `/S,_dann_Gegenstelle_nur_TLS_1.0,_Aushandlung_scheitert`, `TestEinspielTLSZertifikatFehler` (6), `TestEinspielTLSName`, `TestEinspielTLSZertifikatsspeicher` |
| Server lehnt Klartext ab: PGR-E4005 | Klasse 28 nicht mehr PGR-E4005 | `TestEinspielKlartextAbgelehnt` |
| Nach der Aushandlung läuft die Anmeldung auf TLS | Passwort im TLS-Aufbau weggelassen | `TestEinspielTLSKlartextPasswort`; Ende zu Ende `TestE2EPlayTLSAnmeldung` (scram, md5, password) lief grün, ohne eigene Mutation |
| Abbruch im Aufbau: kein Terminate | `Terminate` nach dem Fehler | `TestEinspielTLSAbbruchImAufbau` |
| Das zweite Signal in der Aushandlung schließt | Aushandlung vom Abbruch ausgenommen (`HandshakeContext` ohne `ctx` und `AfterFunc` nach dem Aufbau; beide Schutzmittel zugleich, denn jedes allein genügt) | `TestEinspielTLSAbbruchInDerAushandlung` (5-s-Frist), `TestE2EPlayTLSAbbruch` (10-s-Frist) |
| Das erste Signal lässt den Aufbau weiterlaufen | Der Play-Service übergibt `Verbinde` das Kontext des ersten Signals (Änderung am Kern, nur zur Probe) | `TestE2EPlayTLSAbbruch` (play endet nach dem ersten Signal) |
| `--upstream-tls` wirkt | Option setzt nichts | `TestPlayUpstreamTLS` (10 Zeilen), `TestLeserAlleOptionen` |
| `sslmode=require` ⇒ TLS; `disable`, ohne Angabe und `host:port` ⇒ kein TLS | `sslmode` wirkungslos; jede benannte Verbindung ⇒ TLS | `TestPlayUpstreamTLS/Verbindung_require`; `/Verbindung_ohne_sslmode`, `/Verbindung_disable`; Ende zu Ende `TestE2EPlayTLS/sslmode=require`, `TestE2EPlayTLSAbgelehnt/sslmode=require` |
| Ein gesetztes `--upstream-tls`, auch `false`, geht `sslmode` vor — je Quelle | Option, Umgebungsvariable, Schlüssel zählt je einzeln nicht als gesetzt; `sslmode` setzt immer | Option: `TestPlayUpstreamTLS/require,_Option_false`; Umgebung: `/require,_Umgebungsvariable_false`; Schlüssel: `/require,_Schlüssel_false`; gemeinsam `/require,_Option_false_vor_Umgebungsvariable_true`; Ende zu Ende `TestE2EPlayTLSAbgelehnt/require,_Schlüssel_false`, `TestE2EPlayTLS/Schlüssel_mit_relativem_Pfad` |
| Kommandozeile vor Umgebung vor Schlüssel | Umgebung vor Kommandozeile | `TestPlayUpstreamTLS/Option_false_vor_Umgebungsvariable_true` und `/Option_true_vor_Umgebungsvariable_false`; `TestPlayUpstreamCAQuellen` |
| `sslmode=require` ist bei `play` keine Fehlerverwendung mehr | `zielPlay` lehnt es wieder mit PGR-E2004 ab | `TestPlayVariablen`, `TestPlayUpstreamTLS` (6 Zeilen) |
| Bootstrap reicht TLS und CA durch | `TLS` weggelassen; `CA` weggelassen | `TestRunPlayTLS/TLS_mit_Zertifizierungsstelle` und `/TLS_ohne_Zertifizierungsstelle`; `/TLS_mit_Zertifizierungsstelle`; Ende zu Ende `TestE2EPlayTLS` (5 Zeilen), `TestE2EPlayTLSAnmeldung` (3) |
| `--upstream-ca` ohne TLS: PGR-E2001 (Bedingung *Option gesetzt*) | Fehler bei jeder Verwendung ohne TLS, auch ohne Option | `TestLeserAlleOptionen`, `TestParsePlay*`, `TestPlayVerbindung` |
| … (Bedingung *ohne TLS*) | Fehler auch mit TLS | `TestPlayUpstreamCAOhneTLS/host:port,_Option_TLS`, `/host:port,_Umgebungsvariable_TLS`, `/host:port,_Schlüssel_TLS`, `/require` |
| … je Weise ohne TLS und je Quelle | Prüfung entfernt | `TestPlayUpstreamCAOhneTLS` (8 Zeilen: host:port × Option/Umgebung/Schlüssel, Verbindung ohne und mit `disable`, `require` mit `false` aus Option, Umgebung, Schlüssel) |
| … nach `--upstream`, vor den Variablen | Prüfung nach `zielPlay`; Prüfung vor `pruefeUpstream` | `TestPlayUpstreamCAOhneTLSReihenfolge` (beide) |
| Die Datei aus `--upstream-ca` wird als letzte Prüfung gelesen | Lesen vor den Variablen der Platzhalter | `TestPlayUpstreamCAOhneTLSReihenfolge` |
| Quelle von `--upstream-ca`: Option, Umgebung, Schlüssel | Umgebungsvariable ignoriert; Schlüssel ignoriert; Schlüssel auf der obersten Ebene | `TestPlayUpstreamCAQuellen`, `TestPlayUpstreamCAOhneTLS/host:port,_Umgebungsvariable` und `/host:port,_Schlüssel` |
| Relativer Pfad ab dem aktuellen Verzeichnis, auch in der Datei | relativer Pfad gegen ein anderes Verzeichnis | `TestPlayUpstreamCAQuellen` |
| E2007: reguläre Datei | Prüfung entfernt | `TestPlayUpstreamCAKeineRegulaere` (FIFO endet nicht binnen 10 s), `TestPlayUpstreamCADateiFehler/Verzeichnis`; Ende zu Ende `TestE2EPlayUpstreamCAFehler/Verzeichnis` |
| E2007: mindestens ein Block | Prüfung entfernt | `TestPlayUpstreamCADateiFehler` (6 Zeilen: `leer`, `nur_Text`, `nur_Zeilenumbrüche`, `Block_ohne_Ende`, `Block_mit_ungültigem_base64`, `abgeschnittenes_Zertifikat`) |
| E2007: jeder Block hat den Typ CERTIFICATE | Typ nicht geprüft | `TestPlayUpstreamCADateiFehler/Zertifikat_im_Block_TRUSTED_CERTIFICATE`, `/Zertifikat,_dann_Zertifikat_im_Block_X509_CERTIFICATE`, `/Zertifikat_im_Block_X509_CERTIFICATE,_dann_Zertifikat` |
| … **jeder** Block, nicht nur der erste | nur der erste Block gelesen | `TestPlayUpstreamCADatei/zwei_Zertifikate`, `/zwei_Zertifikate_andersherum`, `TestPlayUpstreamCADateiFehler/Zertifikat,_dann_Block_CERTIFICATE_ohne_Zertifikat` |
| E2007: lesbares X.509-Zertifikat | Fehler des Lesens ignoriert | `TestPlayUpstreamCADateiFehler/Block_CERTIFICATE_ohne_Zertifikat` (und zwei Zeilen mit gültigem Nachbarn) |
| Text außerhalb der Blöcke bleibt unbeachtet | Text vor dem ersten Block abgelehnt; Text nach dem letzten abgelehnt | `TestPlayUpstreamCADatei/Text_davor`; `/Text_danach` |
| Die Meldung nennt weder Pfad noch Inhalt | Pfad bei Fehler des Betriebssystems; Pfad bei fehlender Datei; Inhalt des Blocks | `TestPlayUpstreamCADateiFehler/Pfad_unterhalb_einer_Datei`; `/fehlende_Datei`; `/Block_CERTIFICATE_ohne_Zertifikat` (3 Zeilen) |
| Der Code ist PGR-E2007 (Exit 2) | Wert der Konstante verändert | `TestPlayUpstreamCADateiFehler` (alle Zeilen) |
| Die Datei vor dem Laden der Aufzeichnung | **keine Mutation gefahren**: Das Lesen steht im Parser der Kommandozeile, das Laden im Bootstrap danach; die Reihenfolge folgt aus dem Aufbau | `TestE2EPlayUpstreamCAFehler` (die Aufzeichnung fehlt, die Meldung ist trotzdem PGR-E2007 oder PGR-E2001, nie PGR-E3001) lief grün |
| `config show` liest die Datei nicht | `config show` liest sie | `TestConfigShowLiestCANicht` |
| Eine Formatierung gibt kein Zertifikat aus | `Format` der Zertifikate entfernt (CLI und Upstream-Adapter); `Format` gibt den Namen aus | `TestPlayOptionenOhneZertifikat`; `TestEinspielzielOhneZertifikat` (beide) |
| Standardpfad: ohne `--upstream-ca` gilt der Zertifikatsspeicher des Systems, den das Binary als frischer Prozess lädt | Standardpfad `x509.SystemCertPool` durch einen leeren Speicher (M14 der Review) | `TestE2EPlayTLSZertifikatsspeicherDesSystems/SSL_CERT_FILE_nennt_die_Zertifizierungsstelle` (PGR-E4005, unknown authority) |
| … ohne `SSL_CERT_FILE` ist dieselbe Verbindung PGR-E4005 (verneinend: nichts eingespielt, nichts weitergereicht) | `InsecureSkipVerify` | `TestE2EPlayTLSZertifikatsspeicherDesSystems/ohne_SSL_CERT_FILE` (als einziger Teilfall rot) |
| Kein Paket unter `./cmd/...` importiert `testing`, den TLS-Proxy oder die Zertifikatserzeugung (verneinend; je Eintrag der Verbotsliste) | Blank-Import von `tlsproxy` in `bootstrap.go`; von `testing`; von `testpki` | `TestBinaryOhneTesthilfen`: der erste nennt alle drei Pakete, der zweite `testing`, der dritte `testpki` und `testing` (`testpki` importiert `testing`; `testpki` allein ist deshalb nicht von `testing` zu trennen) |
| Der hängende Proxy schließt `Hello` einmal, auch bei einer zweiten Verbindung (Testhilfe) | `sync.Once` entfernt, `close(p.Hello)` | `TestHaengenderProxyZweiteVerbindung` (Absturz „close of closed channel“) |
| Der Hilfetext nennt die Optionen | Zeile von `--upstream-ca` verändert | `TestPlayHilfeTLS` |

**Grenzen und nicht Geprüftes.**

- *Bytes nach `S`*: erkannt werden Bytes, die mit `S` im selben Lesen eintreffen; was erst
  später eintrifft, stört die Aushandlung und ist PGR-E4005 (`LH-FA-20.a` *TLS* nennt
  PGR-E4002 für jedes Byte vor der Aushandlung). Die Tests erzeugen den Fall mit einem
  Schreiben; ein Server, der `S` und sein Rauschen in zwei Segmenten sendet, ist ungeprüft.
- *IPv6 ohne Zone*: wie in der Tabelle; das Verhalten ist das der Standardbibliothek.
- *Zertifikatsspeicher des Systems*: Der Standardpfad ist Ende zu Ende geprüft
  (`TestE2EPlayTLSZertifikatsspeicherDesSystems`, `SSL_CERT_FILE` im Prozessumfeld des
  Binaries); die Fälle *leerer*, *nicht ladbarer* und *gefüllter* Speicher prüft der Unit-Test
  über das austauschbare Feld. Ohne CA bleiben die Zeilen *leerer* und *nicht ladbarer Speicher*
  gegen einen Mutanten, der die Systemwurzeln des Containers einsetzt, grün (sie enthalten die
  Test-CA nicht), und auch mit CA bleibt er grün (siehe *Ein zweiter Mutant bleibt grün
  (V-167)*): Kein Test trennt den leeren Speicher von den Systemwurzeln des Containers. Nicht
  geprüft: der Speicher des Produkt-Images im Betrieb. Gemessen: Das Image
  `pgwire-recorder:dev` enthält `/etc/ssl/certs/ca-certificates.crt`; der Integrationstest
  fährt das Binary des Test-Images, nicht das des Produkt-Images (Risiko 2 in §6).
- *Testhilfen im Binary*: `TestBinaryOhneTesthilfen` prüft die Pakete unter `./cmd/...`; ob eine
  Testdatei außerhalb von `./cmd/...` `tlsproxy` oder `testpki` importiert, prüft nichts
  (`make a-check` nimmt `internal/bootstrap/**` aus, ein Gate dafür wäre eine ADR).
- *MD5 und SCRAM über TLS* sind nur Ende zu Ende geprüft (`TestE2EPlayTLSAnmeldung`), ohne
  eigene Mutation; im Unit-Test nur das Klartext-Passwort.
- *`localhost` und IPv6*: `TestEinspielTLSName` setzt voraus, dass `localhost`
  auflöst (sonst endet er ohne Prüfung); `TestEinspielTLSIPv6MitZone` setzt IPv6 auf der
  Schleife voraus. Beide liefen hier.
- Der Mutant *Aushandlung vom Abbruch ausgenommen* nimmt zwei Schutzmittel zugleich, weil
  jedes allein genügt (Redundanz im Produkt, kein Mangel des Tests).

**Namensabgleich.** Jeder Test und jeder Teilfall in diesem Abschnitt wurde nach der
Nacharbeit, vor der Übergabe, mit einem Skript gegen das Repo abgeglichen (Funktion
`func <Name>(` für 45 Testnamen; für 70 Teilfälle der Name mit Unterstrichen als Leerzeichen
als Zeichenkette irgendwo in den Go-Quellen, was nicht an den Test gebunden ist). Nicht
gefunden: `TestE2EPlayZwischenstand` (als entfernt genannt) und zwei Teilfälle von
`TestE2EPlayTLSZertifikatsspeicherDesSystems`, deren Name den Unterstrich von `SSL_CERT_FILE`
selbst trägt (das Skript ersetzt ihn); beide stehen als `PASS` im Lauf von `make gates`. Der
Abgleich des Abschnitts *Hinweise* und der Pfade ist nicht Teil des Skripts.

**Hinweise an Architect, Planner und Review.**

1. *Testhilfen und `.a-check.yml`.* `make a-check` meldet `crypto/tls` und `pgproto3` in
   `test/integration` und in einem Hilfspaket außerhalb der PGWire-Adapter als `tech-leak`.
   Ohne Gate zu ändern (`AGENTS.md` §3.6) liegt der TLS-Proxy der Integrationstests deshalb in
   `internal/bootstrap/tlsproxy` (Verdrahtung, dort ist `crypto/tls` zulässig) und die
   Zertifikatserzeugung in `internal/testpki` (ohne `crypto/tls`). Beides gehört zu keiner
   Schicht; der a-check-Hinweis *gescannte Dateien ohne Schicht* wächst um die beiden Dateien
   von `testpki`. Der Architect hat den Ort am 2026-10-11 bestätigt (§6, *Rückgabe des
   Implementers*); dass kein Binary-Paket sie importiert, prüft `TestBinaryOhneTesthilfen`.
2. *Typ der Zertifikate.* `PlayOptions` trägt `cli.Zertifikate` und der `Einspielziel`
   `postgres.Zertifikate`, benannte Typen über `[]*x509.Certificate` mit eigenem `Format`
   (Auftrag: kein Zertifikatsinhalt in Ausgaben); §6 *Ort der CA-Datei* nennt sie jetzt so.
   `PlayOptions` und `Command` sind damit nicht mehr mit `==` vergleichbar; die Tests
   vergleichen mit `reflect.DeepEqual`.
3. *Hook-Signatur im Leser.* `option.zuletzt` bekommt zusätzlich `quellen` (Stand und Wert je
   Option), weil `--upstream-tls` „gesetzt“ von „Standardwert“ unterscheiden und `--upstream-ca`
   erst zuletzt gelesen werden muss; `lies` ist dafür in `liesUmgebung` und `zusammenfuehren`
   geteilt (Komplexität des Lint-Profils).
4. *`ohneTLS` bei `play`.* `zielPlay` prüft `sslmode=require` nicht mehr; `ohneTLS` gilt nur
   noch für `record` und nennt `record` fest. Die Tests des Zwischenstands (`TestPlayVariablen`,
   `TestParsePlayFremdeUmgebung`, `TestDateiUngueltig`, `TestE2EPlayZwischenstand`) sind
   angepasst beziehungsweise entfernt.
5. *SSLRequest-Senden.* Ein gescheitertes Senden des SSLRequest habe ich als *gescheitertes
   Senden im Aufbau* gelesen (PGR-E4002), nicht als Fehler der Aushandlung (PGR-E4005); die
   Aushandlung beginnt in `LH-FA-20.a` *TLS* mit `S`. Wenn der Architect es anders meint, ist
   es eine Zeile.
6. *Handbuch und README* sind nicht angefasst (Doku-Folge-Slice); das Handbuch nennt
   `sslmode=require` und `--upstream-tls` bei `play` noch als nicht verfügbar.

7. *Befunde der Review vom 2026-10-11.* F-602, F-603 und F-604 sind mit den Zeilen oben
   behoben (Tests, Mutanten, Namen). **F-605** (vier Randformen im Code vor §6): Der Ablauf ist
   gelaufen, wie der Befund ihn beschreibt; die Randformen sind inzwischen in §6 entschieden, an
   der Reihenfolge ändert sich nichts mehr, kein Code. **F-606** (`close(p.Hello)`): behoben mit
   `sync.Once` und einem Test, weil der Aufwand klein war (vier Zeilen und ein Test von rund
   40 Zeilen) und der Absturz den ganzen Testprozess trifft, nicht nur einen Test; Zeile
   *Hello* in der Tabelle. **F-607** (Grenze am falschen Satzteil in `LH-FA-20.a`): Sache der
   Spezifikation, nicht dieser Nacharbeit; der Architect hat sie berichtigt. **F-608** (Größe
   rund 14 % über dem Maß): wie in *Größe* oben benannt; das Maß gehört dem Planner, die
   Nacharbeit fügt rund 190 Zeilen Go hinzu (Tests und Testhilfe), keinen Produktcode.

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

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `f7c9c79` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft` (1×) — die TLS-Aushandlung steht
  vor dem Aufbau in `Open`; nach `S` liest der Adapter über eine neue Schicht, und Bytes vor
  der Aushandlung sind `PGR-E4002` (§6). Die Adresse für den Test des Lesepuffers bleibt
  `slice-v1-abschluss-anmeldung`; hier keine Zusage.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×) — der Zertifikatsspeicher des Systems
  hängt an der Plattform; geprüft wird nur unter Linux im Container. Er betrifft Dateisystem
  und Plattform, nicht dieselbe Eigenschaft; kein Beleg vor dem Code.
- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×), `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel`
  (1×) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (1×) — betreffen den Schnitt des
  Gebers und werden dort gezählt (§8 von `slice-v1-abschluss-einspielen`); dieser Slice
  übernimmt nur TLS mit drei Liefer-Punkten.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (18×, `AGENTS.md` §3.12),
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (6×),
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (20×, §3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (26×, §3.11),
  `BEO-REPO/plan-folgt-korrektur-nicht` (21×, §3.9) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (4×, §3.13) — verkörpert; die Randformen in §6
  sind vor dem Code entschieden, je Zusage eine Mutation mit Beleg in §7, die Geber
  `slice-v1-abschluss-einspielen` und `slice-v1-abschluss-antwortvergleich` zeigen im selben
  Commit hierher.

Nachgezählt beim Eintragen von F-568 aus `slice-v1-abschluss-einspielen-laufsteuerung`
(2026-10-10, `AGENTS.md` §3.13): Die Bindung ist eine Randform in DoD-Punkt 2; drei
Liefer-Punkte, zwei Schichten nach der Teilung dieses Plans (CLI-Adapter, Upstream-Adapter; der
Bootstrap reicht weiter, §1). Der Geber zählt den Bootstrap als eigene Schicht (dort §1,
Schicht-Abgrenzung); nach dessen Teilung berührt dieser Plan schon ohne die Sendung drei
Schichten (`BEO-REPO/schichtteilung-je-plan-verschieden`, 2×). Die Teilung prüft der Architect
vor dem ersten Code-Commit (§4 *Start*).

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-einspielen-tls-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 2. Schichten: zwei Schichten nach der Teilung dieses Plans (Upstream-Adapter, CLI-Adapter; der Bootstrap reicht weiter, §1; die Zeile oben zählt den Bootstrap noch als eigene Schicht, überholt durch den Absatz unten).

Nachgezählt nach der Entscheidung des Nutzers vom 2026-10-10 zum Bootstrap (überholt durch den Absatz *Berichtigt* unten) (`AGENTS.md` §3.13, seit slice-v1-abschluss-einspielen-anmeldung): Der Bootstrap zählt nur mit Logik. Hier liest er die Datei aus `--upstream-ca` beim Start und stuft den Fehler ein (`PGR-E2007`, §3), das ist Logik; er zählt also. Schichten: drei (Upstream-Adapter, CLI-Adapter, Bootstrap), der Plan liegt über der Grenze. Nicht geschnitten; der Planner entscheidet den Schnitt vor dem Anspruch (`next` → `in-progress`), der Architect prüft ihn vor dem ersten Code-Commit (§4 *Start*). Die Zeilen oben zählten den Bootstrap als Verdrahtung und sind überholt.

**Berichtigt am 2026-10-11 (Prüfung des Architect, Nutzerentscheidung Option A; `AGENTS.md` §3.13):** Die Datei aus `--upstream-ca` liest der CLI-Adapter (§6, *Ort der CA-Datei*); der Bootstrap hat keine Logik, er setzt `UpstreamTLS` und die Zertifikate in den `Einspielziel` und zählt nicht. Die Konstante `PGR-E2007` im Model ist eine Kennung ohne Logik und zählt wie die Verdrahtung nicht. Schichten: **zwei** (Upstream-Adapter, CLI-Adapter), Liefer-Punkte: **zwei** (DoD 1 und 2). Der Plan liegt innerhalb der Grenze; er wird nicht geschnitten, die Rückführung in §4 bleibt die Rückfallebene. Die Absätze oben, die drei Schichten nennen, sind überholt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
