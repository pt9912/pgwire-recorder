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
| `internal/testpki`, `internal/bootstrap/tlsproxy` | neu (Testhilfe) | Zertifikate zur Testzeit mit `crypto/x509`, nichts eingecheckt (`testpki`); der TLS-Proxy für die Integrationstests (`tlsproxy`). Der Proxy liegt unter `internal/bootstrap`, weil `.a-check.yml` `crypto/tls` nur in den beiden PGWire-Adaptern und der Verdrahtung zulässt und `test/integration` ihn nicht importieren darf; `testpki` kommt ohne `crypto/tls` aus. Nur Testcode importiert beide; sie gehören zu keiner Schicht und zählen nicht |

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

- **TLS-Version und Verfahren** — die Voreinstellung der Standardbibliothek, keine Zusage.

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
- **Zertifikatsspeicher des Systems im Integrationstest** — das Zertifikat des Tests liegt
  nicht im Speicher des Test-Containers, und `SSL_CERT_FILE` wirkt erst bei der ersten Ladung
  im Prozess; der Speicher wird im Test des Upstream-Adapters über ein austauschbares Feld
  geprüft (gefüllter Speicher nimmt, leerer und nicht ladbarer Speicher lehnen ab, nur
  `--upstream-ca` nimmt), nicht Ende zu Ende. Die Grenze des Images steht im zweiten Risiko unten.
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
