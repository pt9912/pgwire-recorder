# Slice slice-v1-abschluss-einspielen: Einspielen einer Aufzeichnung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten), [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md), [ADR-0014](../../adr/0014-konfigurationsdatei.md), [ADR-0016](../../adr/0016-einspielen-anmeldung-und-tls.md), [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md)

**Berührte Spec-Stellen:** `LH-FA-20.a` · `LH-FA-17.a` · `LH-FA-03.b` · `LH-FA-02.b` · `LH-FA-01.a` · `LH-FA-12.a` · `LH-FA-14.a` · `SPEC-017` · `SPEC-022` · `SPEC-026` · `SPEC-028` · `SPEC-033` · `SPEC-034` · `SPEC-041` · `ARC-002` · `ARC-003` · `ARC-005` · `ARC-007` · `ARC-009`

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

**Ziel:** `pgwire-recorder play` führt die einfachen Anfragen einer Aufzeichnung gegen einen PostgreSQL-Server aus, der weder Passwort noch TLS verlangt, liest seine Optionen aus Kommandozeile, Umgebung und dem Abschnitt `play:`, bricht bei einer Fehlerantwort des Servers ab und endet bei einem Abbruchsignal nach der laufenden Interaktion.

**Übernommen aus `slice-replay-semantik-mismatch`:** `play` kennt die Option `--fail-on-unconsumed` nicht (`PGR-E2001`, Exit-Code 2) und lässt ihre Umgebungsvariable `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` unbeachtet, auch mit ungültigem Wert (`LH-FA-03.b` §Andere Kommandos, `LH-FA-17.a`); der Test dafür gehört zu den Optionen von `play`.

**Übernommen aus `slice-replay-semantik-meldungscodes`:** `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` bei `play`, mit derselben Wertemenge und Strenge wie bei `record` und `replay` (Optionstabelle, `LH-FA-17.a`).

**Bereinigung aus `slice-replay-semantik-fehlerreplay`:** Die Teile von `LH-FA-20.a` Schritt 3, die aufgezeichnete Interaktionen ohne `ReadyForQuery` behandeln, sind seit der Entscheidung in `LH-FA-02.b` §Fehlerantwort vor dem Abbruch nicht mehr erreichbar: Der Recorder schreibt keine solche Interaktion, und der Leser lehnt sie als beschädigt ab (`PGR-E3003`). Bereinigt hat sie der Architect vor dem Code (2026-10-09, Schritt 3 von `LH-FA-20.a`); die Zeile *Aufgezeichnetes Ende* der Tabelle gilt nur mit Vergleich und bleibt bei `slice-v1-abschluss-antwortvergleich` (dort §6). Die DoD sagt für sie nichts zu.

**Übernommen aus `slice-v1-abschluss-konfiguration`** (dort §1, Abgrenzung): Die Optionen von `play` werden über den allgemeinen Leser jenes Slice gelesen, in derselben Priorität wie bei `record` und `replay`; das gehört zu den Optionen von `play` oben.

**Übernommen aus `slice-v1-abschluss-konfigurationsdatei`** (dort §1, Abgrenzung): der Abschnitt `play:` der Konfigurationsdatei, gelesen über das Laden jenes Slice; er gehört zu den Optionen von `play`.

**Übernommen aus `slice-v1-abschluss-upstream-verbinden`** (dort §1, Abgrenzung; hervorgegangen aus den Schnitten von `slice-v1-abschluss-konfigurationsdatei` und `slice-v1-abschluss-verbindungen-platzhalter` vom 2026-10-09; bis zum zweiten Schnitt gab `slice-v1-abschluss-verbindungen-platzhalter` den Punkt): die Wirkung einer benannten Verbindung bei `play`, die `slice-v1-abschluss-verbindungen-platzhalter` beim Laden prüft und jener Slice auflöst und einsetzt — Host und Port, Benutzer und Datenbank der URL mit dem Vorrang von `--user` und `--database` (`LH-FA-17.a`, *Wirkung einer URL*); das gehört zum Ziel oben. Das Passwort aus dem eingesetzten Platzhalter und TLS zum Upstream nach `sslmode` sind seit dem Schnitt vom 2026-10-09 abgegeben (unten, *Abgegeben*). Dazu, mit der Kennung `slice-v1-abschluss-upstream-verbinden` (dort §6, U8): Bei `play` werden die Platzhalter aller Teile der benutzten Verbindung eingesetzt, auch in Benutzer, Passwort und Datenbank, über das Einsetzen jenes Slice; eine nicht gesetzte oder leere Variable dort ist `PGR-E2005` mit der Verbindung und dem Namen der ersten Variable in der Reihenfolge der URL (Test bei den Optionen von `play`). Befund F-533 aus dessen Review (im Benutzerhandbuch §5 *Konfigurationsdatei* die Wirkung einer Verbindung bei `play` wie geliefert, das Passwort und `sslmode=require`; jener Slice beschreibt dort nur `record`) ist seit dem Schnitt vom 2026-10-09 ganz abgegeben, geteilt nach diesen beiden Teilen (unten, *Abgegeben*). Dazu, mit der Kennung `slice-v1-abschluss-upstream-verbinden` (Befund V-125 aus dessen Verifikation): in Beispiel und Abschnittsliste von §5 *Konfigurationsdatei* wieder der Abschnitt `play:` (mit `upstream` und `input`), den jener Slice durch `record:` ersetzt hat, weil das Laden `play:` vor diesem Slice als unbekannten Schlüssel ablehnt; das Beispiel als Datei startet mit `config show` ohne Meldung.

**Abgegeben** beim Schnitt vom 2026-10-09 (Prüfung des Architect in §6, Risiko *Größe*; Entscheidung des Nutzers), je Teil mit seinen Randformen aus der Prüfung des Architect (Marken [A] und [T]), die seit dem Schnitt in §6 des Nehmers stehen:

- an `slice-v1-abschluss-einspielen-anmeldung` (dort §1, *Übernimmt*): die Passwortquellen (Passwort aus dem eingesetzten Platzhalter, sonst `PGWIRE_RECORDER_PASSWORD`), Klartext, MD5 und SCRAM-SHA-256, `PGR-E4005` der Anmeldung, im Handbuch der Passwort-Teil aus F-533; dazu die Anmeldung beim Einspielen, die `slice-v1-abschluss-anmeldung` (dort §1, Abgrenzung) hierher gegeben hatte;
- an `slice-v1-abschluss-einspielen-tls` (dort §1, *Übernimmt*): `--upstream-tls`, `sslmode=require` der benutzten Verbindung, `--upstream-ca` mit `PGR-E2007`, `PGR-E4005` von TLS und Zertifikat, im Handbuch der `sslmode`-Teil aus F-533.

Beide hängen nur an diesem Slice, nicht aneinander. Bis zu ihnen gilt bei `play` der Zwischenstand in §6 (*Zwischenstand*).

**Abgegeben** beim zweiten Schnitt vom 2026-10-09 (Option O3 des Architect, §6 Risiko *Größe
des Kerns*; Entscheidung des Nutzers), je Teil mit seinen Randformen aus der Prüfung des
Architect (Marken [L], [E] und [L·E]), die seit dem Schnitt in §6 des Nehmers stehen:

- an `slice-v1-abschluss-einspielen-laufsteuerung` (dort §1, *Übernimmt*): die Optionen
  `--continue-on-error`, `--allow-recorded-errors` und `--finish-session-on-interrupt` mit
  allem, was erst durch sie erreichbar ist — Fortsetzung nach einer Fehlerantwort, erwarteter
  Fehler, Exit-Code nach einem früheren Fehler, Session zu Ende nach dem Signal; aus der DoD
  bis zum Schnitt die Zusagen zu den drei Optionen in Punkt 2 und 3;
- an `slice-v1-abschluss-einspielen-extended` (dort §1, *Übernimmt*): das Einspielen einer
  Extended-Interaktion — Gruppen, Warten nach `Sync` und `Flush`, Abbruch ohne weitere
  Gruppe, Signal innerhalb der Interaktion —, aus Ziel und DoD-Punkt 1 bis zum Schnitt
  „einfach und Extended“, die Ablösung des Zwischenstands *Aufzeichnung mit
  Extended-Interaktion* und die Randformen [L·E], weil er in der Reihenfolge nach der
  Laufsteuerung liegt (§5 von [welle-v1-abschluss](../welle-v1-abschluss.md)).

Die Laufsteuerung hängt nur an diesem Slice, der Slice für Extended an diesem und an der
Laufsteuerung ([L·E]); Anmeldung und TLS hängen an keinem der beiden. Bis zur Laufsteuerung
bricht jeder Fehler ab, bis zum Slice für Extended gilt der Zwischenstand in §6.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Anmeldung mit Passwort — `slice-v1-abschluss-einspielen-anmeldung` (oben, *Abgegeben*); sie setzt den Verbindungsaufbau dieses Slice voraus und ist für sich lieferbar. Hier läuft jeder Test gegen einen Server, der kein Passwort verlangt.
- TLS zum Server und die eigene Zertifizierungsstelle — `slice-v1-abschluss-einspielen-tls` (oben, *Abgegeben*); ebenso. Hier baut `play` keine TLS-Verbindung auf.
- Die Optionen der Laufsteuerung und was erst durch sie erreichbar ist —
  `slice-v1-abschluss-einspielen-laufsteuerung` (oben, *Abgegeben*); ohne sie ist jeder
  Fehler ein Abbruch und das Signal endet nach der laufenden Interaktion, beides Zielverhalten
  ohne Option. Hier sind die drei Optionen unbekannt (§6, *Optionen der Laufsteuerung*).
- Das Einspielen einer Extended-Interaktion — `slice-v1-abschluss-einspielen-extended` (oben,
  *Abgegeben*); hier ist sie ein Startfehler `PGR-E6001` (§6, *Zwischenstand*).
- Zeitangaben und zeitgetreues Einspielen — `slice-v1-abschluss-zeitangaben` (dort §1, *Übernommen aus* diesem Slice); es setzt das Einspielen voraus.
- Vergleich der Serverantworten mit der Aufzeichnung — `slice-v1-abschluss-antwortvergleich` (dort §1 und §6); er setzt das Einspielen voraus.
- Paralleles Einspielen — Out-of-Scope von [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung).



## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): Eine Aufzeichnung mit DDL- und DML-Anweisungen als einfache Anfragen wird gegen eine leere Instanz, die weder Passwort noch TLS verlangt, eingespielt, und die Datenbank enthält danach deren Wirkung (Abnahmeszenario 12); mehrere Sessions laufen über eigene Verbindungen nacheinander; die Optionen von `play` wirken aus Kommandozeile, Umgebung und dem Abschnitt `play:`, `--user` und `--database` gehen Benutzer und Datenbank der benutzten Verbindung und der Aufzeichnung vor, eine nicht gesetzte oder leere Variable in Benutzer, Passwort oder Datenbank der benutzten Verbindung ist `PGR-E2005` (U8), ein gewöhnliches Argument `PGR-E2001`; `--fail-on-unconsumed` ist bei `play` unbekannt (`PGR-E2001`) und ihre Umgebungsvariable bleibt dort unbeachtet, auch mit ungültigem Wert; `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` wirken bei `play` mit derselben Wertemenge, Strenge und Schwelle wie bei `record` und `replay`, ein ungültiger Wert ist `PGR-E2001`, auch in der Umgebungsvariable neben gültiger Option; bis zu den Folge-Slices gilt der Zwischenstand aus §6 (Integrationstest). Das Benutzerhandbuch zeigt in Beispiel und Abschnittsliste von §5 *Konfigurationsdatei* den Abschnitt `play:` (mit `upstream` und `input`), und das Beispiel als Datei startet mit `config show` ohne Meldung (V-125).
- [ ] Eine Fehlerantwort des Servers bricht ab (`PGR-E4004`, Exit-Code 4, keine weitere Nachricht, `Terminate`); ein Verbindungsfehler bricht ab; im Aufbau sind ein nicht erreichbarer Server und eine Fehlerantwort außerhalb der SQLSTATE-Klasse 28 (etwa eine fehlende Datenbank) `PGR-E4002`, eine der Klasse 28 (etwa ein unbekannter Benutzer) `PGR-E4005` (Test). Beleg in §7 für Punkt 1 bis 3: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
- [ ] `SIGINT` und `SIGTERM` beenden nach der laufenden Interaktion, ein erstes Signal im Aufbau nach dem Aufbau, ohne Interaktion; ein zweites Signal beendet sofort; der Exit-Code ist 0, weil ohne die Optionen der Laufsteuerung jeder Fehler vorher abbricht; ohne Vergleich ist jedes Verbindungsende nach dem ersten `ReadyForQuery` `PGR-E4003`; die Rangfolge mit Vergleich (Exit-Code 5) prüft `slice-v1-abschluss-antwortvergleich` (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` (Play-Service), `internal/hexagon/ports/driving` | neu | Einspiel-Use-Case für einfache Anfragen mit Abbruch bei jedem Fehler, Abbruchsignal und dem Startfehler für eine Extended-Interaktion (§6, *Zwischenstand*); nutzt nur Driven Ports |
| `internal/adapters/driving/cli` | update | Kommando `play`, Optionen, Abschnitt `play:` der Konfigurationsdatei, Einsetzen aller Teile der benutzten Verbindung (U8); `--fail-on-unconsumed` bei `play` unbekannt, ihre Umgebungsvariable unbeachtet (Test, `LH-FA-03.b`); `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` bei `play` wie bei `record` und `replay` (übernommen aus `slice-replay-semantik-meldungscodes`, Test) |
| `internal/adapters/driven/postgres` | update | Verbindungsaufbau als Client ohne Passwort und ohne TLS mit der Einstufung nach `LH-FA-20.a` *Aufbau*, die `Query` einer einfachen Anfrage senden und die Antworten bis `ReadyForQuery` lesen |
| `internal/bootstrap` | update | `play` verdrahten: Upstream-Adapter mit Adresse, Benutzer und Datenbank, Play-Service, Signale (erstes, zweites), Exit-Code; die Optionen von `play` und die Signale gehen an den Play-Service, der über Fortsetzung und Session-Ende entscheidet (`slice-v1-abschluss-einspielen-laufsteuerung` §1, Schicht-Abgrenzung) |
| `test/integration` | update | Happy/Boundary/Negative nach LH-FA-20, gegen einen Server ohne Passwort und ohne TLS |
| `docs/user/benutzerhandbuch.md` | update | §5 *Konfigurationsdatei*: der Abschnitt `play:` in Beispiel und Abschnittsliste (V-125), übernommen aus `slice-v1-abschluss-upstream-verbinden`; Passwort und `sslmode` bei `play` (F-533) beschreiben die Folge-Slices |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`, und `slice-v1-abschluss-herunterfahren`, `slice-v1-abschluss-konfiguration`, `slice-v1-abschluss-konfigurationsdatei`, `slice-v1-abschluss-verbindungen-platzhalter` und `slice-v1-abschluss-upstream-verbinden` sind `done` (Signalbehandlung, allgemeiner Leser, Konfigurationsdatei, Laden und Benutzen benannter Verbindungen; alle aus `slice-v1-abschluss-betrieb` hervorgegangen), dazu nach der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) `slice-v1-abschluss-schreiben`, und [ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nach dem zweiten Schnitt nicht in einer Review-Sitzung prüfbar. Schnitt dann: DoD-Punkt 3 (Abbruchsignal, zweites Signal) als eigener Slice; bis zu ihm braucht `play` eine Regel für das Signal, die der Architect vor dem Schnitt entscheidet (§6, *Größe des Kerns*). Eingetreten sind beide früher benannten Bedingungen, jeweils vor dem ersten Code-Commit, und der Slice blieb beide Male nach Entscheidung des Nutzers in `in-progress/`: *die Authentifizierungsverfahren des Servers sprengen den Slice* beim ersten Schnitt vom 2026-10-09 (Anmeldung und TLS an `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls`), *der Diff von DoD-Punkt 1 und 2 ist nicht in einer Review-Sitzung prüfbar* beim zweiten, nach der Prüfung des Architect am Code (Laufsteuerung und Extended an `slice-v1-abschluss-einspielen-laufsteuerung` und `slice-v1-abschluss-einspielen-extended`; §1, *Abgegeben*).

- `in-progress` → `open`: Das Recording enthält Interaktionen, die sich nicht einspielen lassen — Carveout.


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

**Randformen** (`AGENTS.md` §3.12) — je Operation des Laufs, wo sie entschieden ist. Geprüft
vom Architect am 2026-10-09 vor dem ersten Code-Commit, gegen `LH-FA-20.a`, `LH-FA-17.a`,
`LH-FA-14.a`, `SPEC-034`, [ADR-0016](../../adr/0016-einspielen-anmeldung-und-tls.md),
[ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md),
[ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) und den Bestand
(`cli.go`: `optionen`, `leserKommandos`, `lies`; `verbindung.go`: `einsetzen`, `adresseRecord`;
`postgres/upstream.go`: `Open`, `lies`, `verbindungsende`; `bootstrap.go`: `Run`, `fail`). Neu
entschieden heißt: im Commit dieser Prüfung in `LH-FA-20.a` *Randformen je Schritt*, sonst an
der genannten Stelle. Die Marke in eckigen Klammern sagt, zu welchem Teil des Schnitts vom
2026-10-09 (Risiko *Größe* unten) die Randform gehört: **[K]** Kern, also dieser Slice; die
Randformen **[A]** (Anmeldung) und **[T]** (TLS) stehen seit dem Schnitt in §6 von
`slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls`. Offen ist
keine.

Für den zweiten Schnitt (Entscheidung des Nutzers vom 2026-10-09, Option O3 des Architect,
Risiko *Größe des Kerns* unten) markiert der Architect am 2026-10-09 zusätzlich: **[L]**
Laufsteuerung, das sind `--continue-on-error`, `--allow-recorded-errors` und
`--finish-session-on-interrupt` mit allem, was erst durch sie erreichbar ist (Fortsetzung nach
einem Fehler, erwarteter Fehler, Exit-Code nach einem früheren Fehler, Session zu Ende nach dem
Signal); **[E]** Extended, das Einspielen einer Extended-Interaktion (Gruppen, Warten, Signal
und Abbruch innerhalb ihrer). Eine Randform mit **[L·E]** betrifft beide; sie liefert und
testet der zweite der beiden Slices in der Reihenfolge, die der Planner setzt. Beim Kern bleiben
das Abbruchsignal ohne Option und das Einspielen einfacher Anfragen. Die Stelle je Randform
ändert der Schnitt nicht; den Schnitt selbst (Kennungen, §1, DoD, Abgabe) macht der Planner.
Geschnitten vom Planner am 2026-10-09: Reihenfolge Kern, `slice-v1-abschluss-einspielen-laufsteuerung`,
`slice-v1-abschluss-einspielen-extended`; die Randformen [L·E] liefert damit der Slice für
Extended. Die Randformen [L] stehen seit dem Schnitt in §6 der Laufsteuerung, die [E] und
[L·E] in §6 des Slice für Extended; hier bleiben die mit [K], und wo eine von ihnen einen
Teil [L] oder [E] hat, nennt sie ihn mit dem Nehmer.

*Lesen der Optionen (CLI-Adapter)*

- **Optionen von `play`** [K] — `--upstream`, `--input`, `--user`, `--database`,
  `--log-level`, `--config` am allgemeinen Leser, in der Reihenfolge der Optionstabelle,
  Abschnitt `play:`; bestätigt, `LH-FA-17.a`. Werte von `--user` und `--database`: jeder nicht
  leere Text (`artText`). `--upstream-tls` und `--upstream-ca` liefert
  `slice-v1-abschluss-einspielen-tls` (*Zwischenstand* unten).
- **Optionen der Laufsteuerung** (`--continue-on-error`, `--allow-recorded-errors`,
  `--finish-session-on-interrupt`) — liefert `slice-v1-abschluss-einspielen-laufsteuerung`
  (dort §6, [L]). Bis dahin gilt für sie dasselbe wie für die *Optionen der Folge-Slices* darunter: unbekannt (`PGR-E2001`,
  unbekannte Option, `SPEC-034`), ihre Schlüssel in `play:` unbekannt (`PGR-E2004`,
  `LH-FA-17.a` Fehlertabelle), ihre Umgebungsvariablen unbeachtet (`LH-FA-17.a`, „Die
  Umgebungsvariable einer Option, die das Kommando nicht kennt, bleibt unbeachtet“). Ein
  Zwischenstand mit eigener Regel ist dafür nicht nötig: Ohne die drei Optionen bricht jeder
  Fehler ab (Tabelle der Fehlerregeln, Spalte *Abbruch*), und das Signal endet nach der
  laufenden Interaktion (Schritt 7); beides ist Zielverhalten ohne Option.
- **Optionen der Folge-Slices** (`--keep-timing`, `--timing-mode`, `--timing-reference`,
  `--compare-responses`) — bis zu `slice-v1-abschluss-zeitangaben` und
  `slice-v1-abschluss-antwortvergleich` unbekannt (`PGR-E2001`), ihre Schlüssel im Abschnitt
  `play:` unbekannt (`PGR-E2004`), ihre Umgebungsvariablen unbeachtet; derselbe Stand wie bei
  `record` (`RecordOptions`: „soweit dieser Stand sie kennt“). Beide Slices führen die Optionen
  in ihrem §3 (`internal/adapters/driving/cli`); keine Zuweisung.
- **`--fail-on-unconsumed` und `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED`** [K] — unbekannt
  beziehungsweise unbeachtet, auch mit ungültigem Wert; bestätigt, `LH-FA-03.b` §Andere
  Kommandos.
- **Gewöhnliches Argument** [K] — `PGR-E2001`; neu entschieden in `LH-FA-01.a`.
- **Abschnitt `play:` bei anderen Kommandos** [K] — mit `play` am allgemeinen Leser prüft das
  Laden ihn auch bei `record` und `replay`; bestätigt, `LH-FA-17.a` („unabhängig vom
  Kommando“).
- **Variablen aller Teile** (U8 aus `slice-v1-abschluss-upstream-verbinden`) [K] —
  `PGR-E2005` mit Verbindung und erster Variable in der Reihenfolge der URL, auch in einem
  Teil, den `--user` oder `--database` überschreibt, und im Passwort auch, wenn der Server
  keines verlangt; neu entschieden (Überschreiben) in `LH-FA-17.a` *Wirkung einer URL*, der
  Rest bestätigt (`einsetzen` über alle Teile).

*Start (nach den Optionen)* — alle [K]

- **Aufzeichnung nicht ladbar** [K] — wie bei `replay`, Exit-Code 3; bestätigt, Schritt 1.
- **Aufzeichnung ohne Session mit Interaktion** [K] — keine Verbindung, Exit-Code 0; neu
  entschieden in `LH-FA-20.a` *Start*. Eine Session nur aus Lebendprüfungen wird eingespielt
  (bestätigt, `LH-FA-12.a`).
- **Ausgabe eines Startfehlers** [K] — Zeile beim Prozessende (`fail`), keine Verbindung; neu
  entschieden in `LH-FA-20.a` *Start*.
- **Signal während des Starts** [K] — bricht ihn nicht ab, danach keine Session; neu
  entschieden in `LH-FA-20.a` *Start*.
- **Interaktion einer Art, die `play` nicht einspielt** [K] — Startfehler `PGR-E6001`,
  Exit-Code 6, keine Verbindung; Einzelheiten im *Zwischenstand* unten (Extended).

*Verbindungsaufbau je Session (Upstream-Adapter)*

- **Startup-Daten** [K] — aufgezeichnete Parameter unverändert, `user` und `database` nach dem
  Vorrang Option, URL, Aufzeichnung; ohne Wert fehlt der Parameter; Serverparameter nicht
  gesendet; neu entschieden in `LH-FA-20.a` *Startup-Daten*.
- **Server nicht erreichbar, Verbindungsende vor dem ersten `ReadyForQuery` ohne
  Fehlerantwort, unerwartete oder nicht lesbare Nachricht im Aufbau** [K] — `PGR-E4002`; neu
  entschieden in `LH-FA-20.a` *Aufbau*. *Hinweis:* `Open` meldet ein Verbindungsende im Aufbau
  heute als `PGR-E4003` und eine Anmeldeanforderung als `PGR-E6001` (für `record`); `play`
  braucht die Einstufung nach `LH-FA-20.a`.
- **Fehlerantwort im Aufbau** [K] — Klasse 28 `PGR-E4005`, sonst `PGR-E4002`, gleich welcher
  Schweregrad; neu entschieden in `LH-FA-20.a` *Aufbau* (die Tabelle nannte `FATAL`).
- **`BackendKeyData`, `CancelRequest`** [K] — verworfen, nie gesendet; neu entschieden in
  `LH-FA-20.a` *Aufbau*.
*Zwischenstand bis zu den Folge-Slices* — kein Stand des Produkts, den die Spezifikation
beschreibt, sondern der Stand dieses Slice, bis `slice-v1-abschluss-einspielen-anmeldung`,
`slice-v1-abschluss-einspielen-tls` und `slice-v1-abschluss-einspielen-extended` geliefert sind; vom Planner am 2026-10-09 nach dem Vorbild
der *Optionen der Folge-Slices* oben gesetzt (derselbe Stand wie bei `record`, kein stilles
Herabstufen auf eine Verbindung ohne Passwort oder ohne TLS). Vom Architect am 2026-10-09 vor
dem ersten Code-Commit geprüft (`AGENTS.md` §3.12): Jeder Punkt folgt aus einer allgemeinen
Regel der Spezifikation, angewandt auf den Umfang, den dieser Slice liefert; keine Regel nennt
einen Slice oder einen Stand. Die Folge-Slices ersetzen die Punkte und ändern die Tests dazu
(dort §1 *Aufsetzen* und §6 *Ablösung des Zwischenstands*); die Spezifikation ändern sie dafür
nicht.

- **`--upstream-tls`, `--upstream-ca`** [K] — bei `play` unbekannt (`PGR-E2001`), ihre
  Schlüssel im Abschnitt `play:` unbekannt (`PGR-E2004`), ihre Umgebungsvariablen unbeachtet;
  bestätigt, `LH-FA-17.a` (Umgebungsvariable einer Option, die das Kommando nicht kennt;
  *Fehler*, unbekannter Schlüssel) und `PGR-E2001` (unbekannte Option), wie bei den
  *Optionen der Folge-Slices* oben.
- **`sslmode=require` der benutzten Verbindung** [K] — bei `play` `PGR-E2004` wie bei `record`
  (Bestand aus `slice-v1-abschluss-upstream-verbinden`); `play` baut keine Verbindung ohne TLS
  auf, die TLS verlangt. Neu entschieden in `LH-FA-17.a` *Wirkung einer URL* (jedes Kommando
  ohne TLS zum Upstream, an der Stelle von `record` in der Reihenfolge, Meldung mit der
  Verbindung).
- **Passwort-Anforderung des Servers** (Klartext, MD5, SASL) [K] — wie ein nicht unterstütztes
  Verfahren: `PGR-E4005`, nichts gesendet, auch mit einem Passwort aus dem Platzhalter; neu
  entschieden (allgemein gefasst) in `LH-FA-20.a` *Anmeldung* („ein Verfahren, das `play` nicht
  unterstützt“), Abbruch sofort nach der Tabellenzeile *Verfahren nicht unterstützt*.
- **`PGWIRE_RECORDER_PASSWORD`** [K] — unbeachtet; eine Variable im Passwortteil der benutzten
  Verbindung wird dennoch eingesetzt und ist, wenn sie fehlt, `PGR-E2005` (U8 oben).
  Bestätigt: Ein Passwort geht nur in ein Verfahren (`LH-FA-20.a` *Anmeldung*), und ohne
  unterstütztes Verfahren geht keines; die Variable hat damit keine beobachtbare Wirkung, eine
  Prüfung ihres Werts kennt `LH-FA-17.a` nicht. `config show` listet sie wie jede aktive
  Variable (`LH-FA-17.a` *Anzeige*).
- **Aufzeichnung mit Extended-Interaktion** [K], abgelöst von
  `slice-v1-abschluss-einspielen-extended` [E] — `play` spielt bis zu ihm nur einfache
  Anfragen ein; eine Extended-Interaktion ist damit eine
  *Interaktion einer Art, die `play` nicht einspielt*. Vom Architect am 2026-10-09 vor dem
  ersten Code-Commit neu entschieden, allgemein gefasst in `LH-FA-20.a` *Art der Interaktion*
  und `SPEC-034` (`PGR-E6001`), je Operation:
  - *Erkennen beim Laden*: nach dem Laden (Schritt 1), vor der ersten Verbindung, über alle
    Sessions; ein Ladefehler (Exit-Code 3) geht vor; nicht erst beim Senden.
  - *Gemischte Sessions*: Eine Extended-Interaktion irgendwo, auch in einer späteren Session
    oder nach einfachen Anfragen derselben Session, verhindert jede Verbindung; kein
    Teil-Einspielen, weil `play` eine Datenbank verändert und ein halber Lauf nicht
    wiederholbar ist.
  - *Meldungstext*: `PGR-E6001` mit `id` der Session, `sequence` und Art der ersten solchen
    Interaktion (Reihenfolge der Sessions, dann der Interaktionen); keine Anfrage, keine
    Parameterwerte (`SPEC-033`).
  - *Exit-Code*: 6 (Klasse des Codes, `SPEC-034`), Startfehler, auch nach einem Signal
    während des Starts.
  - *Log*: Zeile beim Prozessende (`LH-FA-14.a` Startfehler), keine Log-Zeile `error` und
    keine Zeile `info` zum Start.
  - *Aufzeichnung ohne Session mit Interaktion*: unverändert Exit-Code 0, sie enthält keine.
  `slice-v1-abschluss-einspielen-extended` ersetzt den Punkt durch das Zielverhalten (Schritt 3 für
  Extended-Interaktionen) und ändert die Tests dazu; die Spezifikation ändert er dafür nicht.
  *Akzeptiertes Negativ:* Im Zielstand hat die Regel *Art der Interaktion* keinen Fall, weil
  `play` jede Art des Formats einspielt; sie bleibt als allgemeine Regel stehen, damit der
  Folge-Slice nur das Zielverhalten liefert (Entscheidung des Nutzers), und gilt für eine
  künftige Art des Formats, die `play` noch nicht einspielt.

*Interaktion (Play-Service und Upstream-Adapter)*

- **Andere Server-Nachrichten** (`NoticeResponse`, `ParameterStatus`, `NotificationResponse`,
  Zeilen) [K] — gelesen und verworfen; neu entschieden in `LH-FA-20.a` *Interaktion*.
- **`Copy…Response`, nicht lesbare Nachricht** [K] — `PGR-E6001`, Exit-Code 6, sofort; neu
  entschieden in `LH-FA-20.a` (Tabellenzeile, *Interaktion*) und `SPEC-034`.
- **Fehlerantwort `FATAL` oder `PANIC`** [K] — sofort `PGR-E4003`, kein `PGR-E4004`, kein
  Warten auf das Ende; Schweregrad aus `V`, ohne es `S`; neu entschieden in `LH-FA-20.a`
  *Interaktion*. Damit hat das zweite Risiko unten seinen Ort.
- **Senden scheitert** [K] — `PGR-E4003`, kein Weiterlesen; **Lesen endet** nach dem ersten
  `ReadyForQuery` — `PGR-E4003`; neu entschieden beziehungsweise bestätigt (Tabelle) in
  `LH-FA-20.a`.
- **Abbruch nach `PGR-E4004`** [K] — keine weitere Antwort gelesen, keine weitere Nachricht,
  `Terminate` („keine weitere Gruppe“ ist [E], `slice-v1-abschluss-einspielen-extended`); neu entschieden (Lesen) in `LH-FA-20.a` *Interaktion*, sonst bestätigt,
  Schritt 6.
- **Frist** [K] — keine eigene, auf Antworten wie auf den Aufbau; neu entschieden in
  `LH-FA-20.a` *Interaktion*.

*Meldungen und Log (Bootstrap und Play-Service)*

- **Fehler nach dem Start** [K] — je Fehler (je `PGR-E4004` bei `--continue-on-error` ist [L],
  `slice-v1-abschluss-einspielen-laufsteuerung`) eine Log-Zeile `error` mit `code` und `error`; nennt `id`, `sequence`, SQLSTATE und `M` der
  Fehlerantwort, keine weiteren Felder, nie Passwort oder Parameterwerte; neu entschieden in
  `LH-FA-20.a` *Meldungen* und `LH-FA-14.a`.
- **Zeilen der Stufe `info`** [K] — Start (`upstream=host:port`, Aufzeichnung; nie Benutzer,
  Passwort, Datenbank), Ende, erstes Abbruchsignal; neu entschieden in `LH-FA-14.a` und
  `LH-FA-20.a` *Meldungen*.

*Ende einer Session, Signal, Exit-Code*

- **`Terminate` und Schließen scheitern** [K] — keine Meldung, Exit-Code unverändert; offene
  Transaktion setzt der Server zurück; neu entschieden in `LH-FA-20.a` *Ende einer Session*.
- **Erstes Signal im Aufbau** [K] — der Aufbau läuft zu Ende, ein Fehler darin zählt; danach
  keine Interaktion (mit `--finish-session-on-interrupt` die ganze Session ist [L],
  `slice-v1-abschluss-einspielen-laufsteuerung`); neu entschieden in `LH-FA-20.a`
  *Abbruchsignal*. **Zwischen Sessions** [K] — keine neue (Session läuft ab ihrem Aufbau). In
  einer Extended-Interaktion bis zu ihrem `ReadyForQuery` ist [E],
  `slice-v1-abschluss-einspielen-extended`.
- **Exit-Code beim Abbruchsignal** [K] — 0; ohne die Optionen der Laufsteuerung bricht jeder
  Fehler vorher ab, „4 nach einem früheren Fehler“ ist [L],
  `slice-v1-abschluss-einspielen-laufsteuerung`; bestätigt, Tabellenzeile *Abbruchsignal*.
- **Zweites Signal** [K] — `Terminate`, soweit sofort möglich, Schließen, Ende; die
  unterbrochene Interaktion ist kein Fehler, Exit-Code nach der Zeile *Abbruchsignal* (0 ohne
  vorherigen Fehler); weitere Signale ohne Wirkung; neu entschieden in `LH-FA-20.a`
  *Abbruchsignal* (die Tabelle deckte das zweite Signal schon, der Satz macht es ausdrücklich).

*Akzeptierte Negative der Prüfung vom 2026-10-09* (keine Folgepflicht, einmalig und harmlos):

- **Große Ergebnismenge** — keine Zusage, dass `play` eine Antwort streamt: `SPEC-032` ist ein
  Ziel ohne Zusage, und dieselbe Antwort hielt `record` beim Aufzeichnen schon im Speicher.
- **NUL in `--user` oder `--database`** — nur über ein Escape in der Konfigurationsdatei
  erreichbar; der Server lehnt den Aufbau ab (`PGR-E4002` oder `PGR-E4005`), keine eigene
  Prüfung.
- **Mehrere Fehlerantworten in einer Interaktion** — PostgreSQL sendet höchstens eine je
  `Sync` beziehungsweise `Query`; „je Fehlerantwort eine Meldung“ ist über mehrere
  Interaktionen prüfbar, innerhalb einer nicht.
- **Kein Abbruch der Anfrage beim zweiten Signal** — `play` sendet kein `CancelRequest`, der
  Server führt die Anfrage zu Ende, bis er das Verbindungsende bemerkt; kein Fehler des
  Einspielens.
- **TLS-Version und Verfahren** — ging mit TLS an `slice-v1-abschluss-einspielen-tls` (dort §6).
- **Weitere aufgezeichnete Startup-Parameter** (etwa `replication`) — unverändert gesendet,
  wie aufgezeichnet; was der Server daraus macht, ist seine Antwort.

**Risiken:**

- Das Einspielen verändert eine Datenbank; ein Fehlgebrauch gegen eine falsche Instanz ist durch das Handbuch nur gewarnt — **Ausgang:** offen bis Closure.
- Eine Serverantwort mit `FATAL` beendet die Verbindung; die Behandlung gemäß Spezifikation ist erst im Test belegbar — **Ausgang:** offen bis Closure.
- **Größe** (Prüfung des Architect vom 2026-10-09): Der Slice berührte CLI-Adapter, Play-Service
  mit Port, Upstream-Adapter und Bootstrap — mehr als zwei Schichten — und trug mit Anmeldung
  (drei Verfahren, SCRAM selbst), TLS mit eigener Zertifizierungsstelle, Fehlerregeln,
  Signalen, Konfiguration und Handbuch rund 80 Zusagen; geschätzt **2500 bis 3500 Zeilen**.
  `slice-v1-abschluss-upstream-verbinden` lag mit rund 30 Zusagen bei 712 Zeilen,
  `slice-v1-abschluss-konfigurationsdatei` mit rund 70 bei 1066 und war die Grenze einer
  Review-Sitzung. Geschnitten nach Entscheidung des Nutzers vom 2026-10-09 in drei Slices nach
  den Marken oben (§1, *Abgegeben*) — **Ausgang:** eingetreten:
  `slice-v1-abschluss-einspielen-anmeldung`, `slice-v1-abschluss-einspielen-tls`.
- **Größe des Kerns** (Planner, 2026-10-09): Auch nach dem Schnitt berührt der Slice vier
  Schichten (CLI-Adapter, Play-Service mit Port, Upstream-Adapter, Bootstrap) und trägt 29 der
  44 markierten Randformen. Nach diesem Anteil, grob und nicht gemessen, liegt die Schätzung
  des Architect für den Kern bei 1700 bis 2300 Zeilen und damit über der Grenze von
  `slice-v1-abschluss-konfigurationsdatei`; der Schnitt nach DoD-Punkten steht in §4 als
  vorab benannte Rückführung. Prüfung des Architect am Code (2026-10-09): Bestand trägt den
  allgemeinen Leser, das Einsetzen, Query, Send und Receive des Upstream-Adapters und die
  Signale in `main`; neu sind Play-Service, Port, die Einstufung des Aufbaus für `play`, das
  Warten je Gruppe, Fehler- und Signalsteuerung, Bootstrap und Tests, geschätzt 2200 bis 2800
  Zeilen mit Tests. Der Schnitt aus §4 (DoD-Punkt 3) allein reicht nicht und verlangte eine
  Regel für das Signal bis zum Folge-Slice; Optionen und Empfehlung lagen dem Nutzer vor.
  Entscheidung des Nutzers vom 2026-10-09: Option O3, Laufsteuerung [L] und Extended [E] als
  eigene Slices, das Abbruchsignal ohne Option bleibt im Kern; die Marken stehen oben, der
  Zwischenstand für Extended ist vor dem Code entschieden (*Zwischenstand*), für die
  Laufsteuerung ist keiner nötig (*Optionen der Laufsteuerung*). Stand nach dem Schnitt des
  Planners (2026-10-09): drei Liefer-Punkte, nur einfache Anfragen, Abbruch bei jedem Fehler
  und das Signal ohne Option; die vier Randformen, die nur [L] oder [E] trugen, und die
  Teile [L], [E] und [L·E] von sechs weiteren sind abgegeben (§6 oben). Die Schichten bleiben
  vier — CLI-Adapter, Play-Service mit Port, Upstream-Adapter, Bootstrap —, über der Grenze
  von zwei; das trägt die Entscheidung des Nutzers für O3, nicht die Größenregel. Die Zeilen
  sind nach dem Schnitt nicht neu geschätzt; reicht auch dieser Kern nicht in eine
  Review-Sitzung, gilt die Rückführung in §4 — **Ausgang:** eingetreten:
  `slice-v1-abschluss-einspielen-laufsteuerung`, `slice-v1-abschluss-einspielen-extended`.


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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel
`REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `f7c9c79` beim Schnitt vom 2026-10-09
gesichtet (Zähler = Dateien unter `evidence/`; ein Beleg entsteht bei der Closure dieses
Slice, nicht vorher). Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×: `slice-v1-abschluss-betrieb`,
  `slice-v1-abschluss-konfiguration`; Stand offen) — passt auf diesen Slice: Er war die
  Adresse von sechs Abgrenzungen (§1, *Übernommen aus*), DoD-Punkt 1 bündelte zuletzt
  Abnahmeszenario 12, drei Anmeldeverfahren, TLS mit eigener Zertifizierungsstelle und die
  Optionen, bei drei Liefer-Punkten und vier Schichten; die Größe fand erst die Prüfung des
  Architect. Mit dem Beleg dieses Slice stünde der Eintrag bei **3×** und wäre eine Lücke, die
  einen eigenen Folge-Slice braucht. Dem Nutzer am 2026-10-09 vorgelegt; Ausgang offen bis zu
  seiner Entscheidung.
- `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` (1×: `slice-v1-abschluss-konfiguration`)
  — passt: Der Schnitt geschah im selben Zug, der Slice bleibt in `in-progress/`, der Rest ging
  an zwei neue Slices in `next/` (§4). Mit dem Beleg dieses Slice stünde er bei 2×.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (1×) — kann passen: Der Kern liegt nach der
  groben Schätzung in §6 (*Größe des Kerns*) über der Grenze einer Review-Sitzung, §4 nennt
  den zweiten Schnitt vorab. Beleg erst, wenn diese Rückführung eintritt.
- `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft` (1×) — der Kern baut die
  Aufbau-Schleife in `Open` für die Einstufung nach `LH-FA-20.a` *Aufbau* um; die Adresse für
  den Test bleibt `slice-v1-abschluss-anmeldung`, hier keine Zusage.
- `BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt` (1×) — betrifft `record`; für `play`
  regelt `LH-FA-20.a` *Aufbau* den Fall (§6), kein Beleg.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (18×, `AGENTS.md` §3.12),
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (6×),
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (20×, §3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (26×, §3.11),
  `BEO-REPO/plan-folgt-korrektur-nicht` (21×, §3.9) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (4×, §3.13) — verkörpert; die Randformen in §6
  sind vor dem Code entschieden (auch der *Zwischenstand*, vom Architect am
  2026-10-09 bestätigt), je Zusage eine Mutation mit Beleg in §7, und die Nehmer
  `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls` nennen
  diesen Slice unter *Übernimmt*, im selben Commit wie §1 hier.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF. Produktionscode für `record` und
`replay` liegt vor, `play` entsteht mit diesem Slice; Spezifikation und Entscheidungen gehen
ihm voraus (`LH-FA-20.a`, [ADR-0016](../../adr/0016-einspielen-anmeldung-und-tls.md),
[ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md)).

