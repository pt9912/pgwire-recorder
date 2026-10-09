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
  Grenzen *SASLprep* und *Klartext ohne TLS*;
- die Randformen [A] aus §6 jenes Slice (§6 unten).

**Aufsetzen.** Der Slice setzt auf dem Kern `slice-v1-abschluss-einspielen` auf: Dort sind
das Kommando `play`, die Optionen, der Verbindungsaufbau ohne Passwort und ohne TLS, die
Einstufung einer Fehlerantwort im Aufbau (SQLSTATE-Klasse 28 `PGR-E4005`, sonst
`PGR-E4002`) und `PGR-E2005` für eine Variable im Passwortteil geliefert. Bis zu diesem Slice
ist jede Passwort-Anforderung des Servers bei `play` ein nicht unterstütztes Verfahren
(`PGR-E4005`, nichts gesendet), und `PGWIRE_RECORDER_PASSWORD` bleibt unbeachtet
(Zwischenstand in §6 jenes Slice); diesen Zwischenstand ersetzt dieser Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

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
      unterstütztem Verfahren sendet `play` nichts; das Passwort steht in keiner Meldung und
      keiner Log-Zeile ([`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)) (Test). Beleg in §7 für Punkt 1 und 2: je Zusage
      Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
- [ ] Das Benutzerhandbuch beschreibt in §5 *Konfigurationsdatei* das Passwort einer
      Verbindung bei `play` wie geliefert (aus dem Platzhalter, sonst
      `PGWIRE_RECORDER_PASSWORD`, wenn die Verbindung kein Passwort schreibt; Befund F-533)
      und nennt die Grenzen *SASLprep* und *Klartext ohne TLS*; die Abdeckungstabellen sind
      über `make abdeckung` nachgezogen.
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
| `internal/adapters/driven/postgres` | update | Anmeldung als Client im Aufbau von `play`: Klartext, MD5, SCRAM-SHA-256 ohne Channel Binding; nicht unterstützte Verfahren und Fehler im SCRAM-Austausch als `PGR-E4005`; ersetzt den Zwischenstand des Kerns |
| `internal/adapters/driving/cli` | update | Passwort aus dem eingesetzten Platzhalter der benutzten Verbindung, sonst aus `PGWIRE_RECORDER_PASSWORD` (leer gilt als nicht gesetzt); nie in Meldung oder Log |
| `internal/bootstrap` | update | das Passwort an den Upstream-Adapter von `play` reichen |
| `test/integration` | update | Server mit `password`, `md5` und `scram-sha-256` in `pg_hba.conf`; Happy/Negative nach LH-FA-20 |
| `internal/adapters/driven/postgres` (Tests) | update | Fehler im SCRAM-Austausch und nicht unterstützte Verfahren gegen einen Testserver, der die Nachrichten vorgibt |
| `docs/user/benutzerhandbuch.md` | update | §5 *Konfigurationsdatei*: Passwort einer Verbindung bei `play` (F-533), Grenzen *SASLprep* und *Klartext ohne TLS* |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen` liegt in `done/`
(Kommando `play`, Verbindungsaufbau ohne Passwort, Einstufung der Fehler im Aufbau). Schritt 9
der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) (Schnitt vom
2026-10-09). Die Randformen aus §6 entschied der Architect am 2026-10-09 vor dem Code in
`LH-FA-20.a` *Anmeldung*; vor dem ersten Code-Commit prüft er die Liste gegen den gelieferten
Kern, besonders gegen die Form, in der der Kern den Aufbau in `Open` einstuft
(`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar. Schnitt dann: Klartext und MD5 hier, SCRAM-SHA-256 als eigener
  Slice; jedes Verfahren wirkt für sich, beide Teile sind einzeln lieferbar.
- `in-progress` → `open` (blockiert — Carveout?): SCRAM-SHA-256 lässt sich mit
  `pgproto3` und der Standardbibliothek nicht bauen und verlangt eine weitere Bibliothek im
  Upstream-Adapter; dann zuerst die Entscheidung (ADR) und die Regel in `.a-check.yml`.

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
- **Falsches Passwort** — der Server antwortet im Aufbau mit SQLSTATE-Klasse 28, das ist
  `PGR-E4005` nach der Regel *Aufbau* des Kerns (`LH-FA-20.a` *Aufbau*, Marke [K]); hier
  geprüft mit einer Anmeldung, dort mit einem unbekannten Benutzer.
- **Ablösung des Zwischenstands** — bis zu diesem Slice ist bei `play` jede
  Passwort-Anforderung `PGR-E4005` ohne Senden, und `PGWIRE_RECORDER_PASSWORD` bleibt
  unbeachtet (§6 von `slice-v1-abschluss-einspielen`, *Zwischenstand*). Die Tests des Kerns,
  die das für Klartext, MD5 oder SCRAM-SHA-256 prüfen, ändert dieser Slice; für die nicht
  unterstützten Verfahren bleiben sie.

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

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
