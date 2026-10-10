# Slice slice-v1-abschluss-einspielen-anmeldung: Anmeldung beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md), [ADR-0014](../../adr/0014-konfigurationsdatei.md), [ADR-0016](../../adr/0016-einspielen-anmeldung-und-tls.md)

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

**Ziel:** `pgwire-recorder play` meldet sich beim Server mit Klartext-Passwort, MD5 oder SCRAM-SHA-256 an, mit dem Passwort aus dem Platzhalter der benutzten Verbindung oder aus `PGWIRE_RECORDER_PASSWORD`, und meldet jeden Fehlschlag der Anmeldung als `PGR-E4005`.

**Übernimmt:** `slice-v1-abschluss-einspielen` — dessen Teil *Anmeldung* (Marke [A] in dessen
§6) nach dem Schnitt vom 2026-10-09 (Prüfung des Architect, dort §6 Risiko *Größe*;
Entscheidung des Nutzers; dort §1, *Abgegeben*). Im Einzelnen:

- die Passwortquellen: das Passwort aus dem eingesetzten Platzhalter der benutzten
  Verbindung, sonst aus `PGWIRE_RECORDER_PASSWORD`, eine leere Variable gilt als nicht
  gesetzt (`LH-FA-17.a` *Wirkung einer URL*, `LH-FA-20.a` *Anmeldung*); der Punkt kam nach
  `slice-v1-abschluss-einspielen` aus `slice-v1-abschluss-upstream-verbinden` (dort §1,
  Abgrenzung);
- die Verfahren Klartext, MD5 und SCRAM-SHA-256 ohne Channel Binding und `PGR-E4005` der
  Anmeldung; der Punkt kam nach `slice-v1-abschluss-einspielen` aus
  `slice-v1-abschluss-anmeldung` (dort §1, Abgrenzung *Anmeldung beim Einspielen*), und
  dieser Slice ist jetzt die Adresse jener Abgrenzung und der Abgrenzung *Anmeldung* von
  `slice-v1-abschluss-antwortvergleich` (dort §1);
- im Benutzerhandbuch §5 *Konfigurationsdatei* der Passwort-Teil der Wirkung einer
  Verbindung bei `play` aus Befund F-533 (Review von `slice-v1-abschluss-upstream-verbinden`,
  über `slice-v1-abschluss-einspielen`): das Passwort aus dem Platzhalter,
  `PGWIRE_RECORDER_PASSWORD`, wenn die Verbindung kein Passwort schreibt; dazu die beiden
  Grenzen *SASLprep* und *Klartext ohne TLS*; dieser Handbuch-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-einspielen-anmeldung-doku`;
- die Randformen [A] aus §6 jenes Slice (§6 unten).

**Aufsetzen.** Der Slice setzt auf dem Kern `slice-v1-abschluss-einspielen` auf: Dort sind
das Kommando `play`, die Optionen, der Verbindungsaufbau ohne Passwort und ohne TLS, die
Einstufung einer Fehlerantwort im Aufbau (SQLSTATE-Klasse 28 `PGR-E4005`, sonst
`PGR-E4002`) und `PGR-E2005` für eine Variable im Passwortteil geliefert. Bis zu diesem Slice
ist jede Passwort-Anforderung des Servers bei `play` ein nicht unterstütztes Verfahren
(`PGR-E4005`, nichts gesendet), und `PGWIRE_RECORDER_PASSWORD` bleibt unbeachtet
(Zwischenstand in §6 jenes Slice); diesen Zwischenstand ersetzt dieser Slice.

**Aus `slice-v1-abschluss-einspielen-laufsteuerung`** (Review F-568 und F-569 jenes Slice):

- die Bindung an den Vertrag der Ports `Einspielziel` und `EinspielSession`, nach dem jeder
  Fehler seinen Meldungscode trägt; auf ihr ruht das akzeptierte Negativ (e) jenes Slice, und
  jeder neue Fehler der Anmeldung hält sie (§6, *Bindung an den Vertrag der Ports*);
- der Kommentar von `play` im Bootstrap nennt die Kopplung, dass der Exit-Code der der ersten
  Meldung des gelieferten Fehlers ist und der Play-Service den abbrechenden Fehler darin zuerst
  stellt (`AGENTS.md` §3.7, Kopplung); der Slice ändert den Bootstrap ohnehin (§3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-einspielen-anmeldung-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- TLS zum Server, `--upstream-tls`, `--upstream-ca` und `sslmode=require` —
  `slice-v1-abschluss-einspielen-tls` (dort §1, *Übernimmt*); beide hängen nur am Kern, nicht
  aneinander. Dass `play` ein Klartext-Passwort auch ohne TLS sendet, ist hier eine Grenze
  (§6), keine Zusage über TLS.
- Die Anmeldung im Record-Modus — `slice-v1-abschluss-anmeldung`; `record` vermittelt den
  Austausch des Clients, `play` rechnet ihn selbst, das ist ein anderer Vorgang.
- Kerberos, GSSAPI, SSPI und Channel Binding — außerhalb des Funktionsumfangs von v1;
  `LH-FA-20.a` *Anmeldung* stuft sie als nicht unterstützt ein (`PGR-E4005`), und das prüft
  dieser Slice.
- Die Anmeldung gegen jede unterstützte PostgreSQL-Version — Bestand bleibt stehen: Der
  Integrationstest läuft gegen die Referenzversion 17 wie alle übrigen;
  `slice-v1-abschluss-postgres-versionen` fährt nur den Fall aus
  `slice-v1-abschluss-anmeldung` (dort §1) und ist keine Adresse für diesen (§8,
  `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`).
- Code im Kern (Play-Service, Ports) und im PGWire-Adapter — Schicht-Abgrenzung: Die
  Anmeldung liegt im Upstream-Adapter, das Passwort liest der CLI-Adapter, der Bootstrap
  reicht es weiter.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): `play` meldet sich an einem Server mit `password`, `md5` und
      `scram-sha-256` in `pg_hba.conf` an und spielt eine Aufzeichnung ein; das Passwort
      kommt aus dem Platzhalter der benutzten Verbindung, sonst aus
      `PGWIRE_RECORDER_PASSWORD`, eine leere Variable gilt als nicht gesetzt; verlangt der
      Server kein Passwort, sendet `play` keines (Integrationstest).
- [ ] [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler): Ein falsches Passwort, ein fehlendes Passwort, wenn der Server eines
      verlangt, ein nicht unterstütztes Verfahren (Kerberos, GSSAPI, SSPI, SASL ohne
      `SCRAM-SHA-256`) und jeder Fehler im SCRAM-Austausch (eine Nachricht des Servers, die
      nicht passt oder sich nicht lesen lässt, eine falsche Serversignatur) sind
      `PGR-E4005` mit Exit-Code 4 und brechen sofort ab; bei fehlendem Passwort und nicht
      unterstütztem Verfahren sendet `play` nichts; eine Fortsetzung im SCRAM-Austausch, die
      nicht passt, und eine nicht lesbare Anforderung eines unterstützten Verfahrens sind
      `PGR-E4005`, ein Verbindungsende während der Anmeldung `PGR-E4002`, nach dem Fehler
      ohne `Terminate` (aus `slice-v1-abschluss-einspielen`, §6); ein Signal in der Anmeldung folgt
      den Regeln des Aufbaus; das Passwort steht in keiner Meldung und
      keiner Log-Zeile, keiner Ursachenkette und in keiner Formatierung der Optionen ([`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)) (Test). Beleg in §7 für Punkt 1 und 2: je Zusage
      Zusage · Mutation · roter Test (`AGENTS.md` §3.10). Die Abdeckungstabellen sind über `make abdeckung` nachgezogen.
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
| `internal/adapters/driven/postgres` | update | Anmeldung als Client im Aufbau von `play`: Klartext, MD5, SCRAM-SHA-256 ohne Channel Binding (neue Datei `anmeldung.go`, Standardbibliothek); nicht unterstützte Verfahren und Fehler im SCRAM-Austausch als `PGR-E4005`; ersetzt den Zwischenstand des Kerns; `Einspielziel` trägt das Passwort als `Passwort`, das bei keiner Formatierung ausgegeben wird |
| `internal/adapters/driving/cli` | update | Passwort aus dem eingesetzten Platzhalter der benutzten Verbindung, sonst aus `PGWIRE_RECORDER_PASSWORD` (leer gilt als nicht gesetzt) in `PlayOptions.Passwort` (Typ `Passwort`, bei keiner Formatierung ausgegeben); Hilfetext von `play`; nie in Meldung oder Log |
| `internal/bootstrap` | update | das Passwort an den Upstream-Adapter von `play` reichen (reine Verdrahtung, keine Schicht, Entscheidung des Nutzers vom 2026-10-10); der Kommentar von `play` nennt die Kopplung an die erste Meldung des gelieferten Fehlers (F-569 aus `slice-v1-abschluss-einspielen-laufsteuerung`) |
| `test/integration`, `tools/test/run-integration-tests.sh` | update | der Runner legt auf der Instanz drei Benutzer mit `scram-sha-256`, `md5` und `password` in `pg_hba.conf` an und gibt ihre Passwörter an den Testlauf; neue Datei `play_anmeldung_e2e_test.go`: Happy/Negative nach LH-FA-20 gegen die reale Instanz |
| `internal/adapters/driven/postgres` (Tests) | update | neue Datei `einspielen_anmeldung_test.go`: Fehler im SCRAM-Austausch und nicht unterstützte Verfahren gegen einen Testserver, der die Nachrichten vorgibt; Vektoren aus RFC 7677; der Fake-Server `einspielServer` merkt die Client-Nachrichten |
| `internal/adapters/driving/cli`, `internal/bootstrap` (Tests) | update | Passwortquellen und Vorrang, Formatierung der Optionen; neue Datei `internal/bootstrap/play_anmeldung_test.go`: Passwort durch die Verdrahtung, Signal in der Anmeldung; `TestParsePlayFremdeUmgebung` und `TestPlayVerbindung` ändern sich (Ablösung des Zwischenstands) |
| `docs/user/abdeckung-e2e.md`, `docs/user/abdeckung-unit.md` | erzeugt | `make abdeckung` aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen` liegt in `done/`
(Kommando `play`, Verbindungsaufbau ohne Passwort, Einstufung der Fehler im Aufbau). Schritt 12
der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) (Schnitt vom
2026-10-09). Die Randformen aus §6 entschied der Architect am 2026-10-09 vor dem Code in
`LH-FA-20.a` *Anmeldung*; vor dem ersten Code-Commit prüft er die Liste gegen den gelieferten
Kern, besonders gegen die Form, in der der Kern den Aufbau in `Open` einstuft
(`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar. Schnitt dann: Klartext und MD5 hier, SCRAM-SHA-256 als eigener
  Slice; jedes Verfahren wirkt für sich, beide Teile sind einzeln lieferbar. Der Schnitt senkt die Zahl der Schichten nicht (§8).
  Ebenso zurück, sobald Logik in `internal/bootstrap` hineinkommt (mehr als das Reichen des Passworts und der Kommentar): Der Bootstrap zählt dann als dritte Schicht (§8, `AGENTS.md` §3.13).
- `in-progress` → `open` (blockiert — Carveout?): SCRAM-SHA-256 lässt sich mit
  `pgproto3` und der Standardbibliothek nicht bauen und verlangt eine weitere Bibliothek im
  Upstream-Adapter; dann zuerst die Entscheidung (ADR) und die Regel in `.a-check.yml`. Geprüft vom
  Architect am 2026-10-10: trägt mit der Standardbibliothek, die Bedingung ist nicht eingetreten (§6, *Bibliothek für SCRAM*).

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

**Randformen** (`AGENTS.md` §3.12) — die mit Marke [A] aus §6 von
`slice-v1-abschluss-einspielen`, geprüft vom Architect am 2026-10-09 vor dem ersten
Code-Commit jenes Slice. Neu entschieden heißt: im Commit jener Prüfung in `LH-FA-20.a`, sonst
an der genannten Stelle. Offen ist keine.

- **Passwortquellen** [A] — das Passwort aus dem eingesetzten Platzhalter der benutzten
  Verbindung, sonst aus `PGWIRE_RECORDER_PASSWORD`; bestätigt, `LH-FA-17.a` *Wirkung einer
  URL* und `LH-FA-20.a` *Anmeldung*. Eine nicht gesetzte oder leere Variable im Passwortteil
  ist `PGR-E2005` und liegt beim Kern (U8), auch wenn der Server kein Passwort verlangt.
- **Leere `PGWIRE_RECORDER_PASSWORD`** [A] — gilt als nicht gesetzt; neu entschieden in
  `LH-FA-20.a` *Anmeldung*.
- **Passwort fehlt, Server verlangt eines** [A] — `PGR-E4005`, nichts gesendet; **Server
  verlangt keines** — keines gesendet; neu entschieden in `LH-FA-20.a` *Anmeldung*.
- **Nicht unterstütztes Verfahren** (Kerberos, GSSAPI, SSPI, SASL ohne `SCRAM-SHA-256`) [A] —
  `PGR-E4005`, nichts gesendet; SCRAM ohne Channel Binding; neu entschieden in `LH-FA-20.a`
  *Anmeldung*.
- **Fehler im SCRAM-Austausch** (unpassende oder nicht lesbare Nachricht, falsche
  Serversignatur) [A] — `PGR-E4005`; neu entschieden in `LH-FA-20.a` *Anmeldung*.
- **SASLprep** [A] — nicht angewandt (Grenze); **Klartext ohne TLS** — gesendet, wenn
  verlangt (Grenze); neu entschieden in `LH-FA-20.a` *Anmeldung*; das Handbuch nennt beide.
- **Anmelde-Nachrichten im Austausch** [A] (aus `slice-v1-abschluss-einspielen`, §6,
  Randform-Rückgabe R2, entschieden vom Architect am 2026-10-09 vor dem ersten Code-Commit
  jenes Slice) — die Art einer Nachricht `R` nach ihrem Code, den der Adapter selbst liest;
  eine Fortsetzung (11, 12) im laufenden SCRAM-Austausch, die nicht passt, ist ein Fehler im
  Austausch (`PGR-E4005`), außerhalb eines Austauschs `PGR-E4002`; lässt sich der Rest der
  Anforderung eines unterstützten Verfahrens nicht lesen (etwa MD5 ohne vollständiges Salz,
  SASL mit nicht lesbarer Liste der Verfahren), ist das `PGR-E4005`; neu entschieden in
  `LH-FA-20.a` *Anmelde-Nachrichten*. Die Einstufung der Codes ohne Austausch (nicht
  unterstützte Anforderungen, Fortsetzung ohne Austausch, `R` nach `AuthenticationOk`)
  liefert und prüft der Kern [K].
- **Verbindungsende und Senden während der Anmeldung** [A] (aus
  `slice-v1-abschluss-einspielen`, §6, Lesart des Implementers bestätigt) — `PGR-E4002`, nicht
  `PGR-E4005`, auch im SCRAM-Austausch und nach dem Senden eines Passworts; nach einem Fehler
  keine weitere Nachricht, kein `Terminate`; neu entschieden in `LH-FA-20.a` *Aufbau* und
  *Abbruch im Aufbau* (die Regel liefert der Kern, hier geprüft mit einem Passwort).
  Nachgezählt beim Eintragen beider Punkte (`AGENTS.md` §3.13): drei Liefer-Punkte, zwei
  Schichten (CLI-Adapter, Upstream-Adapter); beide liegen im Upstream-Adapter und in
  DoD-Punkt 2.
- **Bindung an den Vertrag der Ports** (aus `slice-v1-abschluss-einspielen-laufsteuerung`, dort
  §6 akzeptiertes Negativ (e), Review F-568) — jeder neue Fehler der Anmeldung (fehlendes
  Passwort, nicht unterstütztes Verfahren, Fehler im SCRAM-Austausch, Verbindungsende) verlässt
  den Upstream-Adapter mit seinem Meldungscode (`PGR-E4005` oder `PGR-E4002`), keiner ohne Code;
  sonst hinge er nach `SPEC-034` *Ausgabe* als Ursache an eine frühere Meldung, mit Exit-Code 4
  statt 1. Entschieden vom Architect am 2026-10-10 dort: keine eigene Regel, der Ausschluss ruht
  auf dem Vertrag der Ports (Kommentar an `Verbinde`). Geprüft mit DoD-Punkt 2, der für jeden
  Fehlschlag den Code zusagt; kein neuer Liefer-Punkt.
- **Weitere Randformen, vom Architect am 2026-10-10 vor dem ersten Code-Commit entschieden** [A]
  — alle neu in `LH-FA-20.a` *Passwort*, *Verfahren* und *SCRAM-Austausch*:
  - Quelle und Vorrang: Passwortteil der benutzten Verbindung, sonst
    `PGWIRE_RECORDER_PASSWORD`, auch bei `--upstream` als `host:port`; beides schließt sich
    aus (die Variable bleibt bei einem Passwortteil unbeachtet); unverändert, ohne Kürzen;
    ein nicht verlangtes Passwort ist kein Fehler; NUL im Passwort kein Fall.
  - Aufbau der Anforderungen: Klartext (Rumpf nur der Code), MD5 (Code und vier Byte Salz),
    SASL (Namen mit NUL, leerer Name als Ende, nichts dahinter); jede Abweichung ist
    `PGR-E4005`; MD5 mit dem `user` der gesendeten Startup-Daten, ohne ihn leer (akzeptiertes
    Negativ: Ein Server fordert ohne `user` kein MD5; der Fall trägt keine eigene Zusage).
  - Wahl von `SCRAM-SHA-256` (genaue Schreibweise); `SCRAM-SHA-256-PLUS` allein, leere Liste
    und fehlendes `SCRAM-SHA-256` sind `PGR-E4005` ohne Senden; kein Wechsel des Verfahrens.
  - Nachrichten des Austauschs: Aufbau und Reihenfolge von Server-erste-Nachricht
    (`r=`, `s=`, `i=`, nichts sonst), Servernonce (Präfix der eigenen Nonce, länger), Salz
    (Base64), Iterationszahl (1 bis 10 000 000, Grenze mit Grund in der Spezifikation;
    vom Nutzer am 2026-10-10 bestätigt, bleibt),
    Abschluss (`v=` allein); `e=`, falsche Signatur, jede unvorgesehene Nachricht im
    Austausch `PGR-E4005`; weitere Fortsetzung nach dem Abschluss `PGR-E4002`; weitere
    Anforderung nach einer Antwort außerhalb des Austauschs `PGR-E4002`.
  - Kein Passwort und nichts Abgeleitetes in Meldung, Log und Ursachenkette; die Optionen
    von `play` geben es bei einer Formatierung nicht aus (Test mit `%+v`).
  - Abbruch in der Anmeldung: erstes Signal lässt den Aufbau zu Ende laufen, zweites schließt
    ohne `Terminate`; die Regel liefert der Kern (*Abbruchsignal*), hier mit einem
    Server geprüft, der nach der Passwort-Antwort schweigt (zu DoD-Punkt 2, kein neuer
    Liefer-Punkt).
  - `config show` ändert dieser Slice nicht: Platzhalter bleiben unaufgelöst, die Variable
    steht nur mit Namen (*Anzeige*, Test des Kerns der Konfiguration); akzeptiertes Negativ
    ohne neuen Test, weil die Regel für jede `PGWIRE_RECORDER_*`-Variable gilt.
  - Frist: keine eigene beim Warten auf die Anmeldung (*Interaktion*, *Aufbau*); die
    Berechnung von SCRAM ist nicht unterbrechbar und durch die Obergrenze der Iterationen auf
    Sekunden begrenzt (akzeptiertes Negativ, Grund in der Spezifikation).
- **Randformen aus der Rückgabe des Implementers (§7, Fragen 1 bis 5)**, vom Architect am
  2026-10-10 vor dem Review entschieden [A], neu in `LH-FA-20.a` *Verfahren* und
  *SCRAM-Austausch*:
  - Zweite Anforderung, auch nicht lesbar (Code 5 mit drei Byte Salz nach einer Antwort):
    `PGR-E4002`, der Code bestimmt die Art; der Code stimmt. **Auftrag:** ein Test, der es
    festlegt; Mutation: Lesbarkeit vor `beantwortet` prüfen.
  - Zweite Anforderung eines nicht unterstützten Verfahrens (Code 7, unbekannter Code) nach
    einer Antwort: `PGR-E4005`; der Code stimmt. **Auftrag:** ein Test; Mutation:
    `beantwortet` vor dem Verfahren prüfen.
  - Leeres Salz (`s=`): gültig. Der Code stimmt. **Auftrag:** ein Test, der einen Austausch
    mit leerem Salz gelingen lässt (Berechnung mit leerem Salz gegen einen Vektor aus
    `hashlib`, nicht aus dem Code); Mutation: leeres Salz ablehnen.
  - Iterationszahl mit führender Null (`i=0004096`, auch `i=01`): `PGR-E4005`
    (`posit-number` in RFC 5802). **Der Code weicht ab** (`leseIterationen` nimmt sie an).
    **Auftrag:** Code ändern, Test in `TestScramServerErste` für `i=0004096` und `i=01`;
    Mutation: Prüfung auf die führende Null entfernen. Der Zeilenumbruch in `s=` und `v=`
    bleibt ungültig (Test vorhanden).
  - Fehler der Schlüsselableitung: `PGR-E4005`, der Code stimmt; akzeptiertes Negativ ohne
    Test, Grund: nur im FIPS-Modus der Go-Laufzeit möglich, kein Lauf der Gates setzt ihn;
    die Spezifikation nennt es als Grenze und sagt nicht zu, es sei geprüft (`AGENTS.md`
    §3.11).
- **Bibliothek für SCRAM** — entschieden: Standardbibliothek (`crypto/pbkdf2`, `crypto/hmac`,
  `crypto/sha256`, `crypto/md5`, `crypto/rand`), keine weitere Bibliothek. Die SCRAM-Teile
  von `pgconn` sind nicht exportiert, `pgproto3` stellt nur die Nachrichten
  (`SASLInitialResponse`, `SASLResponse`, `PasswordMessage`). [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md)
  bleibt unberührt, es regelt `pgproto3`; `.a-check.yml` und `gomodguard` verlangen für die
  Standardbibliothek keinen Eintrag; keine neue ADR. Damit entfällt die Bedingung der
  Rückführung `in-progress` → `open` in §4.
- **Falsches Passwort** — der Server antwortet im Aufbau mit SQLSTATE-Klasse 28, das ist
  `PGR-E4005` nach der Regel *Aufbau* des Kerns (`LH-FA-20.a` *Aufbau*, Marke [K]); hier
  geprüft mit einer Anmeldung, dort mit einem unbekannten Benutzer.
- **Ablösung des Zwischenstands** — bis zu diesem Slice ist bei `play` jede
  Passwort-Anforderung `PGR-E4005` ohne Senden, und `PGWIRE_RECORDER_PASSWORD` bleibt
  unbeachtet (§6 von `slice-v1-abschluss-einspielen`, *Zwischenstand*). Die Tests des Kerns,
  die das für Klartext, MD5 oder SCRAM-SHA-256 prüfen, ändert dieser Slice; für die nicht
  unterstützten Verfahren bleiben sie. Geändert: `TestParsePlayFremdeUmgebung` (die Variable
  ist nicht mehr fremd) und `TestPlayVerbindung` (das Passwort steht in den Optionen); die
  Aufbau-Tests des Kerns (`TestEinspielAufbauFehler`) bleiben, weil sie keine Passwort-Quelle
  setzen: Ohne Passwort ist jede Anforderung von Klartext, MD5 und SASL `PGR-E4005` ohne Senden.

**Risiken:**

- SCRAM-SHA-256 rechnet `play` selbst; ein Fehler in der Berechnung fällt nur gegen einen
  echten Server auf, ein Testserver, der die Nachrichten vorgibt, rechnet mit demselben
  Fehler — **Ausgang:** offen bis Closure.
- Die Anmeldung ist nur gegen die Referenzversion 17 geprüft
  (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, §8) — **Ausgang:** offen bis
  Closure.

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

**Belege des Implementers** (Code-Commit `c8b7f44` auf `f9bc4d0`):

*Größe und Schichten.* `git diff --shortstat f9bc4d0 c8b7f44`: 18 Dateien, +1971 −77. Produkt-Code
+449 −57 (`anmeldung.go` neu +389, `einspielen.go` +16 −49, CLI `cli.go` +22 −2, `upstream.go`
+8 −1, `verbindung.go` +10 −3, Bootstrap `bootstrap.go` +4 −2), Tests +1434 −11, Runner der
Integrationstests +38 −1, erzeugt (`make abdeckung`) +39 −2, Plan §3 +11 −6. Zwei Schichten
(Upstream-Adapter, CLI-Adapter); der Bootstrap enthält nur das Reichen des Passworts und den
Kommentar zur Kopplung (F-569), keine Logik, `git diff f9bc4d0 c8b7f44 -- internal/bootstrap/bootstrap.go`
ist +4 −2. Keine Rückführung aus §4 ist eingetreten: SCRAM-SHA-256 baut mit der
Standardbibliothek (`crypto/pbkdf2`, `crypto/hmac`, `crypto/sha256`, `crypto/md5`, `crypto/rand`),
`go.mod` und `.a-check.yml` sind unverändert.

*Randformen.* Jede Operation des Diffs ist gegen §6 gemessen: Senden (Antwort auf eine
Anforderung; scheitert es, `PGR-E4002`, *Verbindungsende und Senden*), Empfangen (Code der
Nachricht `R`, *Anmelde-Nachrichten*), Schließen (nach einem Fehler ohne `Terminate`, *Abbruch im
Aufbau*), Fehler (`PGR-E4005` oder `PGR-E4002` nach der Tabelle in *Anmeldung* und *SCRAM-Austausch*),
Aufräumen (`Verbinde` schließt bei jedem Fehler), Zufall (18 Bytes, Base64, *SCRAM-Austausch*),
Zeit (keine eigene Frist, *Frist*). Keine Operation entscheidet eine Randform außerhalb von §6.
Die Kombinationen und Lesarten, die der erste Lauf offen ließ, sind in §6 entschieden (Randformen aus
der Rückgabe des Implementers) und durch Tests festgelegt (Tabelle, Zeilen *Nacharbeit*).

*Weg der Mutanten.* Je Mutant eine frische Kopie von `internal/` (Python `shutil.copyfile`, ohne
Übernahme der mtime, gleichwertig `cp -r` ohne `-p`) unter dem Scratch-Verzeichnis, genau eine
Ersetzung im Wortlaut, danach im Image der Stufe `test` (`pgwire-recorder:test`) mit der Kopie
als Bind-Mount über `/src/internal` (nicht über den Build-Kontext, die mtime ist ohne Wirkung):
`gofmt -l` leer, `go vet` ohne Befund, dann `go test -count=1` des Pakets. Rot heißt: ein
genannter Test schlug mit der genannten Meldung fehl, nicht Build, Vet oder gofmt, und nicht über
eine Gesamtfrist. Die E2E-Mutanten liefen in einer Kopie aller verfolgten Dateien außer
`.harness/` mit `make test-integration` (Container, Netz und Volume räumt der Runner ab; `docker ps -a`,
`docker network ls` und `docker volume ls` mit `pgr-it` danach leer). 84 Unit-Mutanten (76 auf dem
Stand `c8b7f44`, 8 in der Nacharbeit): 83 rot, 1 grün und äquivalent (Zeile am Ende der Tabelle); 5 E2E-Mutanten, alle rot.

| Zusage (§6) | Mutation | roter Test |
|---|---|---|
| Klartext: das Passwort unverändert (*Passwort*, *Verfahren*) | `strings.TrimSpace` auf das gesendete Passwort | `TestAnmeldungKlartext` (Passwort `"  ä%41$ ss "`) |
| Klartext: der Rumpf ist nur der Code, jedes Byte dahinter macht die Anforderung nicht lesbar | Prüfung `len(rest) != 0` entfernt | `TestAnmeldungAnforderungNichtLesbar/Klartext_mit_Rest` |
| MD5: Benutzer ist der `user` der gesendeten Startup-Daten, ohne `user` leer | Benutzer fest `postgres` | `TestAnmeldungMD5/anderer_Benutzer`, `/ohne_Benutzer`, `/Benutzer_u`; `TestAnmeldungWeitereAnforderung` |
| MD5: MD5 aus Passwort und Benutzer, in dieser Reihenfolge | `benutzer + passwort` | `TestAnmeldungMD5` (fünf Vektoren aus `hashlib`, nicht aus dem Code) |
| MD5: gefolgt vom Salz | Salz nicht angehängt | `TestAnmeldungMD5` (fünf Fälle, darunter `anderes_Salz`) |
| MD5: die Antwort beginnt mit `md5` | Präfix leer | `TestAnmeldungMD5` (fünf Fälle) |
| MD5: genau vier Byte Salz (mehr als vier ist nicht lesbar) | `len(rest) < 4` | `TestAnmeldungAnforderungNichtLesbar/MD5_mit_fünf_Byte` |
| MD5: genau vier Byte Salz (weniger als vier ist nicht lesbar) | `len(rest) > 4` | `TestAnmeldungAnforderungNichtLesbar/MD5_ohne_Salz`, `/MD5_mit_drei_Byte` |
| kein Passwort, Klartext verlangt: `PGR-E4005`, nichts gesendet | Prüfung `fehlendesPasswort` in `klartext` entfernt | `TestAnmeldungOhnePasswort/Klartext` (gesendet `p\x00…`), `TestAnmeldungVerbinde`, `TestEinspielAufbauFehler/Klartext` |
| kein Passwort, MD5 verlangt: `PGR-E4005`, nichts gesendet | Prüfung in `md5` entfernt | `TestAnmeldungOhnePasswort/MD5`, `TestEinspielAufbauFehler/MD5` |
| kein Passwort, SCRAM verlangt: `PGR-E4005`, nichts gesendet | Prüfung in `sasl` entfernt | `TestAnmeldungOhnePasswort/SCRAM`, `TestEinspielAufbauFehler/SASL` |
| ein fehlendes Passwort ist `PGR-E4005` (Code) | `CodeUpstream` statt `CodeLogin` | `TestAnmeldungOhnePasswort`, `TestAnmeldungVerbinde`, `TestEinspielAufbauFehler` |
| ein leeres Passwort gilt als keines | Bedingung `a.passwort == ""` entfernt | `TestAnmeldungOhnePasswort`, `TestAnmeldungVerbinde`, `TestEinspielAufbauFehler` |
| kein Passwort wird gesendet, wenn der Server keines verlangt | bei `AuthenticationOk` ein `PasswordMessage` nachgeschoben | `TestAnmeldungKlartext`, `TestAnmeldungMD5`, `TestAnmeldungScram` (Bytes nach der Antwort), `TestAnmeldungKeinPasswortVerlangt`; E2E `TestE2EPlayAnmeldung/*`, `TestE2EPlayAnmeldungOhneVerlangen` |
| SCRAM: `SCRAM-SHA-256` genau in dieser Schreibweise | `strings.EqualFold` | `TestAnmeldungNichtUnterstuetzt/Kleinbuchstaben` |
| SCRAM: der Name ganz, nicht als Präfix (`-PLUS` allein ist `PGR-E4005`) | `strings.HasPrefix` | `TestAnmeldungNichtUnterstuetzt/nur_PLUS`, `/mit_Anhängsel` |
| SCRAM: gewählt auch wenn nicht das erste der Liste | nur das erste Verfahren geprüft | `TestAnmeldungScram/SCRAM_nach_PLUS` |
| SASL: kein Byte hinter dem leeren Namen | `lesbar` immer wahr | `TestAnmeldungAnforderungNichtLesbar/SASL_mit_Bytes_danach`, `/SASL_mit_zwei_Enden` |
| SASL: jeder Name mit NUL abgeschlossen | fehlendes NUL ohne Fehler | `TestAnmeldungAnforderungNichtLesbar/SASL_ohne_NUL` |
| ein nicht unterstütztes Verfahren ist `PGR-E4005`, nichts gesendet | `CodeUpstream` statt `CodeLogin` | `TestAnmeldungNichtUnterstuetzt` (Kerberos, SCM, GSS, SSPI, unbekannt, SASL ohne SCRAM), `TestEinspielAufbauFehler` |
| kein Wechsel des Verfahrens: nach einer Antwort auf Klartext ist eine weitere Anforderung `PGR-E4002` | `beantwortet` in `klartext` nicht gesetzt | `TestAnmeldungWeitereAnforderung/Klartext_dann_*` (drei Fälle) |
| kein Wechsel des Verfahrens: nach einer Antwort auf MD5 | `beantwortet` in `md5` nicht gesetzt | `TestAnmeldungWeitereAnforderung/MD5_dann_*` (drei Fälle) |
| eine weitere Anforderung ist `PGR-E4002` (Code) | `CodeLogin` statt `CodeUpstream` | `TestAnmeldungWeitereAnforderung`, `TestAnmeldungNachScramAbschluss/{Klartext,MD5,SASL}` |
| eine Fortsetzung (8, 11, 12) ohne laufenden Austausch ist `PGR-E4002` | Fall entfernt | `TestAnmeldungFortsetzungOhneAustausch` (drei Fälle), `TestAnmeldungNachScramAbschluss` |
| gescheitertes Senden der Antwort ist `PGR-E4002` (Senden) | Fehler von `conn.Write` verworfen | `TestAnmeldungSendenScheitert` (die Meldung nennt „nicht zu senden“; der Code allein wäre auch der des Lesens danach) |
| gescheitertes Senden der Antwort ist `PGR-E4002` (Code) | `CodeLogin` statt `CodeUpstream` | `TestAnmeldungSendenScheitert` |
| SCRAM: die Nonce besteht aus 18 Zufallsbytes | 16 Bytes | `TestAnmeldungScram` (Nonce `AQIDBAUGBwgJCgsMDQ4PEBES` genau), `TestAnmeldungScramFehler`, `TestAnmeldungNonce` |
| SCRAM: die Nonce ist je Verbindung eine andere | feste Bytes statt `crypto/rand` | `TestAnmeldungNonce` („zwei Verbindungen mit derselben Nonce“) |
| SCRAM: erste Nachricht mit leerem Benutzer `n=` | `n=u` | `TestAnmeldungScram`, `TestAnmeldungScramFehler` |
| SCRAM: erste Nachricht mit dem Kopf `n,,` | `y,,` | `TestAnmeldungScram`, `TestAnmeldungScramFehler` |
| SCRAM: die Antwort trägt `c=biws` | `c=eSws` | `TestAnmeldungScram`, `TestAnmeldungScramFehler` |
| SCRAM: der Beweis geht aus dem Passwort hervor | Passwort um ein Zeichen verlängert | `TestAnmeldungScram`, `TestAnmeldungScramFehler` |
| SCRAM: die Auth-Message beginnt mit dem client-first-bare | ohne diesen Teil | `TestAnmeldungScram`, `TestAnmeldungScramFehler` |
| SCRAM: Beweis = ClientKey XOR ClientSignature | nur ClientSignature | `TestScramBeweisRFC7677` (Beweis aus RFC 7677), `TestAnmeldungScram`; E2E `TestE2EPlayAnmeldung/scram-sha-256_*` (echter Server, 28P01) |
| SCRAM: die Serversignatur kommt aus dem Server Key | Label `Client Key` | `TestScramBeweisRFC7677` (Signatur aus RFC 7677), `TestAnmeldungScram` |
| SCRAM: PBKDF2 mit den Iterationen des Servers | Iterationen + 1 | `TestScramBeweisRFC7677`, `TestAnmeldungScram` |
| erste Nachricht des Servers: genau drei Attribute (zu viele) | `len(teile) < 3` | `TestScramServerErste/vier_Attribute`, `/Komma_am_Ende`, `TestAnmeldungScramFehler` |
| erste Nachricht des Servers: r=, s=, i= in dieser Reihenfolge | Attribute in jeder Reihenfolge gelesen | `TestScramServerErste/andere_Reihenfolge`, `TestAnmeldungScramFehler/erste_Nachricht:_andere_Reihenfolge` |
| Servernonce beginnt mit der eigenen | Präfixprüfung entfernt | `TestScramServerErste/r_ohne_eigene_Nonce`, `TestAnmeldungScramFehler/…fremde_Nonce` |
| Servernonce ist länger als die eigene (gleich lang ist ein Fehler) | `len(r) < len(nonce)`; ebenso die Längenprüfung ganz entfernt | `TestScramServerErste/r_gleich_der_eigenen_Nonce`, `TestAnmeldungScramFehler/erste_Nachricht:_Nonce_nicht_länger` |
| Salz ist gültiges Base64 | Ergebnis der Dekodierung nicht geprüft | `TestScramServerErste` (drei Fälle), `TestAnmeldungScramFehler/erste_Nachricht:_Salz_kein_Base64` |
| Salz: Base64 ohne Zeilenumbruch | Prüfung auf `\r\n` entfernt | `TestScramServerErste/s_mit_Zeilenumbruch` |
| Salz: Base64 mit Auffüllung | `RawStdEncoding` mit abgeschnittenem `=` | `TestScramServerErste/s_ohne_Auffüllung` |
| Iterationen höchstens 10 000 000 (über der Grenze) | `n > Grenze+1` | `TestScramServerErste/Iterationen_über_der_Obergrenze`, `TestAnmeldungScramFehler/…Iterationen_über_der_Obergrenze` |
| Iterationen: genau 10 000 000 gilt | `n >= Grenze` | `TestScramServerErste/Iterationen_genau_an_der_Obergrenze` |
| Iterationen mindestens 1 | `n >= 0` | `TestScramServerErste/Iterationen_0`, `TestAnmeldungScramFehler/…Iterationen_0` |
| Iterationen nur Ziffern (ohne Vorzeichen) | Ziffernprüfung entfernt | `TestScramServerErste/Iterationen_keine_Zahl`, `/mit_Vorzeichen`, `/mit_Exponent` |
| Iterationen: ein Wert über dem Bereich von `int` läuft nicht über | Prüfung erst nach der Schleife | `TestScramServerErste/Iterationen_2_hoch_64_plus_5` (liefe nach dem Überlauf als 5) |
| Abschluss: die Signatur wird geprüft | Vergleich `hmac.Equal` entfällt | `TestAnmeldungScramFehler/Abschluss:_falsche_Signatur`, `…_ein_Byte_kürzer`, `…_länger`, `TestAnmeldungOhneGeheimnis/SCRAM_falsche_Signatur` |
| Abschluss: genau `v=` (ohne `v=` kein Abschluss) | Präfix nicht verlangt | `TestAnmeldungScramFehler/Abschluss:_ohne_v=`, `TestAnmeldungScram` |
| Abschluss: kein weiteres Attribut hinter der Signatur | Text ab dem Komma abgeschnitten | `TestAnmeldungScramFehler/Abschluss:_v=_mit_weiterem_Attribut` |
| Abschluss beendet den Austausch (danach `PGR-E4002`) | Stufe bleibt `Abschluss erwartet` | `TestAnmeldungNachScramAbschluss` (Code 11, 12, 8, 3, 5, 10), `TestAnmeldungScram` |
| im Austausch ist Code 12 vor Code 11 ein Fehler (auch mit gültigem Inhalt einer ersten Nachricht) | Code 12 an der ersten Stelle angenommen | `TestAnmeldungScramUnvorgesehen/Code_12_mit_dem_Inhalt_der_ersten_Nachricht` |
| im Austausch ist Code 11 nach der Antwort ein Fehler (auch mit gültigem Inhalt eines Abschlusses) | Code 11 an der zweiten Stelle angenommen | `TestAnmeldungScramUnvorgesehen/Code_11_mit_dem_Inhalt_des_Abschlusses` |
| im Austausch ist jede unvorgesehene Nachricht `PGR-E4005` (`AuthenticationOk`, Anforderungen 3, 5, 10, 7, unbekannt, 8) | der Austausch gilt nie als laufend | `TestAnmeldungScram`, `TestAnmeldungScramFehler`, `TestAnmeldungScramUnvorgesehen` (14 Fälle: E4002 statt E4005) |
| ein Fehler im Austausch ist `PGR-E4005` (Code) | `CodeUpstream` statt `CodeLogin` | `TestAnmeldungScramFehler` (alle Fälle), `TestAnmeldungScramUnvorgesehen` |
| falsches Passwort ist `PGR-E4005` mit SQLSTATE in der Meldung, andere Klasse und `ReadyForQuery` vor `AuthenticationOk` `PGR-E4002`, Verbindungsende `PGR-E4002`, nach Klartext, MD5 und SCRAM (erste Nachricht, Antwort) | (Einstufung liegt beim Kern, hier geprüft mit einem Passwort; der Kern-Mutant steht in `slice-v1-abschluss-einspielen`) | `TestAnmeldungAbgelehnt` (4 Abläufe × 4 Ausgänge), E2E `TestE2EPlayAnmeldungFehler` (drei Verfahren × falsch, leer, Anfang des richtigen) |
| Nacharbeit, Frage 1: eine weitere Anforderung, auch nicht lesbar, ist `PGR-E4002` (Code bestimmt die Art), Klartext mit Rest | Prüfung `beantwortet` in `klartext` hinter die Lesbarkeit | `TestAnmeldungWeitereAnforderungArt/{Klartext,MD5}_dann_Klartext_mit_Rest` |
| Nacharbeit, Frage 1: dasselbe für MD5 (drei Byte Salz, ohne Salz) | Prüfung in `md5` hinter die Lesbarkeit | `TestAnmeldungWeitereAnforderungArt/{Klartext,MD5}_dann_MD5_mit_drei_Byte_Salz`, `/…_MD5_ohne_Salz` |
| Nacharbeit, Frage 1: dasselbe für SASL (Liste ohne Ende) | Prüfung in `sasl` hinter die Lesbarkeit | `TestAnmeldungWeitereAnforderungArt/{Klartext,MD5}_dann_SASL_ohne_Ende` |
| Nacharbeit, Frage 1: alle drei Verfahren zugleich | die allgemeine Prüfung entfernt, in jedem Verfahren hinter die Lesbarkeit (gofmt-sauber) | die acht Fälle der drei Zeilen darüber |
| Nacharbeit, Frage 2: ein nicht unterstütztes Verfahren (Code 7, 9, unbekannt) nach einer Antwort bleibt `PGR-E4005` | `beantwortet` vor dem Verfahren prüfen | `TestAnmeldungWeitereAnforderungArt/{Klartext,MD5}_dann_Code_7`, `/…_Code_9`, `/…_unbekannter_Code` (sechs Fälle) |
| Nacharbeit, Frage 3: leeres Salz (`s=`) ist gültig und der Austausch gelingt mit dem Beweis aus `hashlib` | `len(salz) == 0` lehnt ab | `TestAnmeldungScramLeeresSalz` |
| Nacharbeit, Frage 4: Iterationen ohne führende Null (`i=0004096`, `i=01`, `i=010`) | Prüfung `HasPrefix(text, "0")` entfernt | `TestScramServerErste/Iterationen_mit_führenden_Nullen`, `/…mit_einer_führenden_Null`, `/Iterationen_10_mit_führender_Null` |
| Nacharbeit, Frage 4: nur eine Ausprägung der führenden Null abgefangen | Prüfung auf `"00"` statt `"0"` | `TestScramServerErste/Iterationen_mit_einer_führenden_Null`, `/Iterationen_10_mit_führender_Null` (der Fall `i=0004096` allein fängt es nicht) |
| Passwort nicht in Meldung, Log, Ursachenkette: Einspielziel (Formatierung) | `Format` gibt den Wert aus | `TestAnmeldungOhneGeheimnis` (letzte Schleife: `%v`, `%+v`, `%#v`, `%s`, `%q`, `%x`) |
| Passwort nicht in Meldung, Log, Ursachenkette: abgeleiteter Wert (erwartete Signatur) | Signatur in die Meldung | `TestAnmeldungOhneGeheimnis/SCRAM_falsche_Signatur` |
| Passwort nicht in Meldung, Log, Ursachenkette: Inhalt einer Nachricht (Abschluss) | Nachrichtentext in die Meldung | `TestAnmeldungOhneGeheimnis/SCRAM_Abschluss_ohne_v=` |
| Passwort nicht in Meldung, Log, Ursachenkette: Inhalt einer Nachricht (erste Nachricht) | Nachrichtentext in die Meldung | `TestAnmeldungOhneGeheimnis/SCRAM_erste_Nachricht_mit_Inhalt` |
| `Verbinde` reicht `Password` an die Anmeldung | `zugang{}` statt `zugang{passwort: …}` | `TestAnmeldungVerbinde`, `TestAnmeldungNonce`, `TestAnmeldungSignal` |
| Bootstrap reicht das Passwort weiter | `Password:` entfernt | `TestRunPlayPasswort` (vier Fälle), `TestRunPlayErstesSignalInAnmeldung`, `TestRunPlayZweitesSignalInAnmeldung`; E2E `TestE2EPlayAnmeldung/*`, `TestE2EPlayAnmeldungFehler/*` |
| Passwort aus `PGWIRE_RECORDER_PASSWORD` (host:port und Verbindung ohne Passwortteil) | Variable nicht gelesen | `TestPlayPasswort` (drei Fälle), `TestRunPlayPasswort/Variable_bei_host:port` |
| Variable auch bei host:port | bei host:port auf `""` gesetzt | `TestPlayPasswort/host:port_mit_Variable` |
| Variable bei einer Verbindung ohne Passwortteil | jede Verbindung gilt als mit Passwortteil | `TestPlayPasswort/Verbindung_ohne_Passwortteil_mit_Variable`, `/…ohne_Benutzer_mit_Variable` |
| der Platzhalter geht der Variable vor | Variable geht vor, wenn gesetzt | `TestPlayPasswort/Platzhalter_und_Variable_gesetzt`, `/Platzhalter_und_--user`; E2E `TestE2EPlayAnmeldung/*_Platzhalter_vor_der_Variable` |
| Passwort aus dem Platzhalter | Passwortteil ignoriert | `TestPlayVerbindung`, `TestPlayPasswort` (drei Fälle) |
| Passwort unverändert, ohne Kürzen (Platzhalter) | `TrimSpace` | `TestPlayPasswort/Platzhalter_und_Variable_gesetzt` |
| Passwort unverändert, ohne Kürzen (Variable) | `TrimSpace` | `TestPlayPasswort/host:port_mit_Variable` |
| Optionen von `play` geben das Passwort bei keiner Formatierung aus | `Format` gibt den Wert aus | `TestPlayOptionenOhnePasswort` (sieben Verben, Wert, Zeiger, Liste, Abbildung) |
| Hilfe nennt `PGWIRE_RECORDER_PASSWORD` | Name aus dem Text entfernt | `TestParsePlayHilfe` |
| `Verbinde` bricht ab, wenn ctx endet (zweites Signal), auch nach dem Passwort und im SCRAM-Austausch | `context.AfterFunc` durch Attrappe ersetzt | `TestAnmeldungSignal/nach_dem_Passwort`, `/im_SCRAM-Austausch` („endet binnen 5 s nicht“), `TestEinspielAufbauAbgebrochen` |
| nach einem Fehler der Anmeldung kein `Terminate` | `Terminate` vor dem Schließen gesendet | `TestAnmeldungVerbinde`, `TestEinspielAufbauFehler` (sechs Fälle: Bytes nach dem Startup) |
| das erste Signal lässt die Anmeldung zu Ende laufen | `Verbinde` bekommt den ctx des ersten Signals | `TestRunPlayErstesSignalInAnmeldung` („play endet nach dem ersten Signal vor der Antwort des Servers“, Exit-Code 4) |
| das zweite Signal schließt ohne `Terminate`, Exit-Code 0 | (`TestRunPlayZweitesSignalInAnmeldung` hält es; Mutant `AfterFunc` oben) | `TestRunPlayZweitesSignalInAnmeldung`, `TestAnmeldungSignal` |
| E2E: SCRAM gegen den echten Server | Beweis nur ClientSignature (Zeile oben) | `TestE2EPlayAnmeldung/scram-sha-256_Variable`, `/…Platzhalter_vor_der_Variable` |
| E2E: MD5 gegen den echten Server | `benutzer + passwort` (Zeile oben) | `TestE2EPlayAnmeldung/md5_Variable`, `/md5_Platzhalter_vor_der_Variable` |
| E2E: das Passwort erreicht den Server über den Bootstrap | `Password:` entfernt (Zeile oben) | `TestE2EPlayAnmeldung` (sechs Fälle), `TestE2EPlayAnmeldungFehler` |
| E2E: ein unverlangtes Passwort wird nicht gesendet (der Server wertet es als ungültige Nachricht) | Passwort bei `AuthenticationOk` gesendet | `TestE2EPlayAnmeldungOhneVerlangen`, `TestE2EPlayAnmeldung/*_Variable` |
| E2E: Platzhalter vor Variable am echten Server | Variable geht vor | `TestE2EPlayAnmeldung/scram-sha-256_Platzhalter_vor_der_Variable`, `/md5_…`, `/password_…` |

Grüner Mutant, äquivalent: Nach einem Abbruch durch das zweite Signal (`!stop()` in `Verbinde`)
hat `context.AfterFunc` die Verbindung schon geschlossen; ein `Terminate`, das dort nachgeschoben
wird (`verbinde-abbruch-mit-terminate`), geht nirgendwohin. Die Grenze trägt `Verbinde` (§6,
*Abbruch im Aufbau*: „keine weitere Nachricht“, geschlossen); eine Test-Idee dazu gibt es nicht,
weil der Mutant über die Schnittstelle nichts verändert.

Im ersten Lauf der Mutanten grün und deshalb nachgearbeitet (vor dem Code-Commit): (1) Der
Test `TestAnmeldungOhneGeheimnis` war durch eine Bearbeitung der Datei verloren; vier
Geheimnis-Mutanten blieben grün, bis er wieder dastand. (2) `TestAnmeldungSendenScheitert` prüfte nur
den Code `PGR-E4002`, den auch das folgende Lesen liefert; er prüft jetzt die Meldung. (3) Die
Fälle „Code 12 vor Code 11“ und „Code 11 nach der Antwort“ trugen den Inhalt der falschen
Nachricht und waren über den Inhalt abgelehnt, nicht über den Code (äquivalent im Ergebnis); sie
tragen jetzt auch den Inhalt der Stelle, die den Code erwartet (gleich welcher Inhalt, ein anderer
Code ist ein Fehler). (4) Die Obergrenze der Iterationen kippte bei `2^64 + 5` nicht, weil der
Test nur einen Wert nahm, der nach dem Überlauf außerhalb des Bereichs lag. (5) Die Prüfung auf
einen leeren Text in `leseIterationen` war überflüssig (`n >= 1` trägt sie, der Mutant blieb grün);
sie ist entfernt.

*Läufe.* Am Arbeitsbaum vor dem Code-Commit `c8b7f44`: `make test` grün; `make test-integration`
grün (darunter `TestE2EPlayAnmeldung` mit 6 Fällen, `TestE2EPlayAnmeldungOhneVerlangen`,
`TestE2EPlayAnmeldungFehler` mit 9 Fällen gegen PostgreSQL 17 mit `scram-sha-256`, `md5` und
`password` in `pg_hba.conf`); `make lint` grün (0 Befunde); `make abdeckung` geschrieben;
`make gates` grün (Exit-Code 0, 3 min 50 s). Die Gates auf sauberem Baum nach dem Commit dieses
Abschnitts stehen im Bericht an den Reviewer.

*Fragen an den Architect (Randform · Frage), entschieden am 2026-10-10 (§6, Randformen aus der Rückgabe des Implementers).* Frage 1 bis 4 sind umgesetzt, jede mit einem Test, der sie festlegt, und einer roten Mutation (Tabelle, Zeilen *Nacharbeit*); Frage 5 bleibt ein akzeptiertes Negativ ohne Test.

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

- `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft` (1×) — dieser Slice baut die
  Aufbau-Schleife in `Open` für den Anmeldeaustausch von `play` um; die Adresse für den Test
  bleibt `slice-v1-abschluss-anmeldung` (dort §1 und §6), hier keine Zusage.
- `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×) — die Anmeldung läuft nur
  gegen die Referenzversion 17 (§1 Abgrenzung, §6 Risiko); der Eintrag bleibt unter der
  Schwelle, Beleg erst bei Closure, falls ein Review eine Abweichung findet.
- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×), `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel`
  (1×) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (1×) — betreffen den Schnitt des
  Gebers und werden dort gezählt (§8 von `slice-v1-abschluss-einspielen`); dieser Slice
  übernimmt nur die Anmeldung mit drei Liefer-Punkten.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (18×, `AGENTS.md` §3.12),
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (6×),
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (20×, §3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (26×, §3.11),
  `BEO-REPO/plan-folgt-korrektur-nicht` (21×, §3.9) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (4×, §3.13) — verkörpert; die Randformen in §6
  sind vor dem Code entschieden, je Zusage eine Mutation mit Beleg in §7, die Geber
  `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-anmeldung` und
  `slice-v1-abschluss-antwortvergleich` zeigen im selben Commit hierher.

Nachgezählt beim Eintragen von F-568 und F-569 aus `slice-v1-abschluss-einspielen-laufsteuerung`
(2026-10-10, `AGENTS.md` §3.13): Die Bindung ist eine Randform in DoD-Punkt 2, der Kommentar im
Bootstrap eine Zeile in einer Datei, die §3 schon führt; drei Liefer-Punkte, zwei Schichten nach
der Teilung dieses Plans (CLI-Adapter, Upstream-Adapter; der Bootstrap reicht weiter, §1). Der
Bootstrap zählt nur mit Logik, reine Verdrahtung nicht (Entscheidung des Nutzers vom
2026-10-10, `AGENTS.md` §3.13); hier ist er Verdrahtung, also zwei Schichten, ohne die Sendungen
und mit ihnen (`BEO-REPO/schichtteilung-je-plan-verschieden`, verkörpert).

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-einspielen-anmeldung-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 2. Schichten: zwei Schichten nach der Teilung dieses Plans (Upstream-Adapter, CLI-Adapter; der Bootstrap reicht weiter, §1; der Bootstrap zählt als reine Verdrahtung nicht, siehe die Entscheidung unten).

Entscheidung des Nutzers zur Schichtzählung (2026-10-10, nach Prüfung des Architect vor dem ersten Code-Commit, `AGENTS.md` §3.13): Der Bootstrap (`internal/bootstrap`, `composition_root` in `.a-check.yml`) zählt nur dann als Schicht, wenn er Logik enthält; reine Verdrahtung (ein Feld durchreichen, ein Kommentar) ist keine Schicht. Dieser Slice ändert dort das Reichen des Passworts und den Kommentar zur Kopplung aus F-569: Verdrahtung. Schichten: zwei (Upstream-Adapter, CLI-Adapter); Liefer-Punkte: 2; der Plan liegt nicht über der Grenze, der Implementer beginnt nach §4 *Start*. Der Schnitt in §4 (Klartext und MD5 hier, SCRAM-SHA-256 eigener Slice) bleibt Rückführung, nicht Plan. Kommt Logik in den Bootstrap, geht der Slice nach §4 zurück. Die Obergrenze der SCRAM-Iterationen von 10 000 000 (§6) ist vom Nutzer bestätigt und bleibt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
