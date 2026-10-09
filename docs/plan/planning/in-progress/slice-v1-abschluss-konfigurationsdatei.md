# Slice slice-v1-abschluss-konfigurationsdatei: Konfigurationsdatei, benannte Verbindungen und config show

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md), [ADR-0014](../../adr/0014-konfigurationsdatei.md), [ADR-0027](../../adr/0027-yaml-bibliothek.md), [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md)

**Berührte Spec-Stellen:** `LH-FA-01.a` · `LH-FA-17.a` · `LH-FA-03.b` · `LH-FA-14.a` · `SPEC-007` · `SPEC-008` · `SPEC-012` · `SPEC-014` · `SPEC-020` · `SPEC-033` · `SPEC-034` · `SPEC-046` · `ARC-005` · `ARC-013`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Jede Option von `record` und `replay`, die am allgemeinen Leser angemeldet ist, ist auch über die Konfigurationsdatei setzbar, in der Priorität Kommandozeile vor Umgebungsvariable vor Datei vor Default, samt benannten Verbindungen und Platzhaltern; `config show` zeigt die gewählte Datei, die Hilfe geht jeder Prüfung vor, und ein Fehler der Datei ist `PGR-E2004` bis `PGR-E2006`.

**Übernimmt:** `slice-v1-abschluss-konfiguration` — dessen Datei-Hälfte nach dem Schnitt aus
seinem §4 (*zu groß*, eingetreten am 2026-10-08): DoD-Punkt 2 für `config show` und `--config`
(Hilfe vor jeder Prüfung, `config show`) und DoD-Punkt 3 (Datei, Verbindungen, Platzhalter,
`PGR-E2004` bis `PGR-E2006`); dort §1, Abgrenzung. Im Einzelnen, je mit dem ursprünglichen
Geber, über `slice-v1-abschluss-konfiguration` und davor `slice-v1-abschluss-betrieb`:

- **Aus `slice-replay-semantik-mismatch`:** der Schlüssel `fail_on_unconsumed` im Abschnitt
  `replay:`; die Hilfe vor jeder Prüfung von Optionen, Umgebungsvariablen und
  Konfigurationsdatei (`LH-FA-01.a`) für `config show` und `--config`, samt `--` als Ende der
  Optionen.
- **Aus `slice-replay-semantik-meldungscodes`:** der Schlüssel `log_level` auf der obersten
  Ebene der Konfigurationsdatei (`LH-FA-17.a`), mit derselben Wertemenge wie `--log-level`;
  ein ungültiger Wert ist `PGR-E2004`.
- **Aus `slice-v1-abschluss-herunterfahren`** (dort §1, Abgrenzung): der Schlüssel der Frist
  `--shutdown-timeout` in der Konfigurationsdatei.

**Schlüssel der Datei.** Die Schlüssel, die die Datei annimmt, folgen aus der Anmeldung am
allgemeinen Leser von `slice-v1-abschluss-konfiguration`; keine Option wird ein zweites Mal
angemeldet. Ein Schlüssel einer Option aus der Tabelle in `LH-FA-17.a`, die der Stand noch
nicht kennt (etwa `format`, `tls_cert` oder der Abschnitt `play:`), ist bis zu ihrem Slice ein
unbekannter Schlüssel (`PGR-E2004`). `output` und `force` sind am Leser angemeldet und damit
hier Schlüssel des Abschnitts `record:` (Rückgabe 12 in §6). Der
Test über alle angemeldeten Optionen läuft dann über drei Quellen.

**Ort des Lesens.** Die Datei liest der CLI-Adapter (`ARC-005`). Die YAML-Bibliothek ist dort
nur zum Lesen der Konfigurationsdatei zulässig (`ARC-013`,
[ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), ergänzt
[ADR-0027](../../adr/0027-yaml-bibliothek.md)); kein Typ der Bibliothek verlässt den
CLI-Adapter. Die Gegenprobe des Architektur-Gates dafür liefert
`slice-v1-abschluss-konfiguration`; dass ein Import im CLI-Adapter nur der Datei dient, prüft
kein Gate, das bleibt Review.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der allgemeine Leser für Kommandozeile und Umgebungsvariable, die Hilfe von `record` und
  `replay` und die Gegenprobe des Architektur-Gates — `slice-v1-abschluss-konfiguration`;
  er liegt vor diesem Slice, und dieser setzt auf ihm auf.
- Das atomare Schreiben und das Verhalten bei vorhandenem `--output` —
  `slice-v1-abschluss-schreiben`; er folgt diesem Slice. `--output` und `--force` selbst
  meldet `slice-v1-abschluss-konfiguration` am Leser an (Entscheidung zu F-496).
- Die Optionen von `play` und der Abschnitt `play:` — `slice-v1-abschluss-einspielen` (dort
  §1); er liest sie über den allgemeinen Leser und die Datei dieses Slice.
- TLS zum Upstream nach `sslmode` — `slice-v1-abschluss-einspielen`; nur `play` verbindet
  mit TLS zum Server. Hier wird `sslmode` geprüft, und `sslmode=require` ist bei `record`
  `PGR-E2004` (`LH-FA-17.a`).
- Code im Kern, im PGWire-Adapter und in den Driven-Adaptern — Schicht-Abgrenzung: Der
  Slice ändert den CLI-Adapter und den Bootstrap. **Einzige Ausnahme:** drei Konstanten
  `PGR-E2004` bis `PGR-E2006` in der Code-Tabelle `internal/hexagon/model/fehler.go`, ohne
  Logik — die Tabelle ist die eine Code-Tabelle im Quelltext
  ([ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md), Konsequenzen), und dort steht
  schon `PGR-E2001`, das nur der CLI-Adapter und der Bootstrap erzeugen (Entscheidung des
  Architect vom 2026-10-08, §6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Die Hilfe geht jeder Prüfung von Optionen,
      Umgebungsvariablen und Konfigurationsdatei vor, auch für `config show` und `--config`,
      samt `--` als Ende der Optionen; `config show` zeigt die gewählte Datei in der Form aus
      `LH-FA-17.a`, ohne einen aufgelösten Wert (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `--config`, `PGWIRE_RECORDER_CONFIG` und die Standarddatei wählen genau
      eine Datei; für jede am allgemeinen Leser angemeldete Option gilt die Priorität
      Kommandozeile vor Umgebungsvariable vor Konfigurationsdatei vor Default (`SPEC-007`),
      auch für den Schlüssel `fail_on_unconsumed` im Abschnitt `replay:`, den Schlüssel der
      Frist und den Schlüssel `log_level` auf der obersten Ebene (Wertemenge und Strenge von
      `--log-level`); eine ungültige Datei ist `PGR-E2004` (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Benannte Verbindungen (`connections`, `--upstream <Name>`, `sslmode`),
      `${VAR}` und `$${VAR}` verhalten sich wie spezifiziert; eine nicht gesetzte Variable ist
      `PGR-E2005`, ein Klartext-Passwort `PGR-E2006` (Test). Beleg in §7 für alle drei Punkte:
      je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
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
| `spec/spezifikation.md` (`LH-FA-17.a`) | update (Architect, vor dem Code, erledigt am 2026-10-08) | Entscheidung der zwölf Rückgaben aus §6 (`AGENTS.md` §3.12) |
| `internal/adapters/driving/cli` | update | Wahl und Laden der Datei (`--config`, `PGWIRE_RECORDER_CONFIG`, Standarddatei); Schlüssel aus der Anmeldung am allgemeinen Leser als dritte Quelle der Priorität; `fail_on_unconsumed`, `log_level`, Frist; `connections`, `${VAR}`, `$${VAR}`, `sslmode`; Auflösen des Namens bei `--upstream` zu `host:port`; Kommando `config show` und seine Hilfe |
| `internal/hexagon/model/fehler.go` | update | nur die drei Konstanten `PGR-E2004` bis `PGR-E2006` in der Code-Tabelle (§1, Ausnahme; §6) |
| `internal/bootstrap` | update | `config show` ausführen (Ausgabe auf `stdout`); die zusammengeführten Optionen an die Use Cases geben |
| `internal/adapters/driving/cli` (Unit-Tests), `internal/bootstrap` (Tests), `test/integration` | update | Happy/Boundary/Negative nach `LH-FA-17.a` und `LH-FA-01.a`; der Test über alle angemeldeten Optionen läuft über drei Quellen |
| `docs/user/benutzerhandbuch.md` | update, falls abweichend | §4 *Die gewählte Konfigurationsdatei anzeigen* und §5 *Konfigurationsdatei* beschreiben den Zielstand; nachgezogen wird, was die Entscheidungen aus §6 ändern |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-konfiguration` liegt in `done/` (der
allgemeine Leser, aus dessen Anmeldung die Schlüssel folgen). Schritt 4 der Reihenfolge in §5
von [welle-v1-abschluss](../welle-v1-abschluss.md) (Entscheidung des Nutzers vom 2026-10-08).
Die zwölf Randformen und den Ort der Konstanten aus §6 entschied der Architect am 2026-10-08
vor dem Code (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar, oder eine Änderung im Kern über die Konstanten hinaus wird nötig.
  Schnitt dann: *Datei* (DoD-Punkt 1 und 2) und *Verbindungen und Platzhalter* (DoD-Punkt 3);
  der zweite setzt den ersten voraus.
- `in-progress` → `open` (blockiert — Carveout?): Die Form der Datei nach [ADR-0014](../../adr/0014-konfigurationsdatei.md) trägt
  einen Schlüssel nicht, den eine Option braucht, ohne eine neue Entscheidung; dann zuerst
  die Entscheidung.

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

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist. Die entschiedenen
zogen mit dem Schnitt aus §6 von `slice-v1-abschluss-konfiguration` hierher; der Architect
prüfte sie dort am 2026-10-08 vor dem Code. Die zwölf Rückgaben des Implementers und den Ort
der Konstanten entschied er am 2026-10-08 hier, vor dem Code (unten); offen ist keine.

*Umgebung, soweit sie die Datei betrifft*

- **Gesetzte, aber leere `PGWIRE_RECORDER_CONFIG`** — gilt als nicht gesetzt; `LH-FA-17.a`.
- **Umgebungsvariable mit Präfix ohne passende Option** — unbeachtet; `config show` nennt ihren
  Namen; `LH-FA-17.a`.

*Datei: Wahl und Form*

- **Ort der Standarddatei je Betriebssystem** — nur `.pgwire-recorder.yaml` im aktuellen
  Verzeichnis, auf jedem Betriebssystem, kein Benutzer- oder Systemverzeichnis; `LH-FA-17.a`.
- **Standarddatei fehlt / Datei aus `--config` oder `PGWIRE_RECORDER_CONFIG` fehlt** — keine
  Datei gelesen / `PGR-E2004`; `LH-FA-17.a`. Standarddatei vorhanden, aber nicht lesbar:
  `PGR-E2004` (Tabelle *Fehler*).
- **`PGWIRE_RECORDER_CONFIG` nennt eine fehlende Datei, `--config` geht vor** — gelesen und
  geprüft wird nur die gewählte Datei; `LH-FA-17.a` (*genau eine Datei*).
- **Groß- und Kleinschreibung von Schlüsseln** — Schlüssel heißen wie die Optionen;
  `Log_Level` ist ein unbekannter Schlüssel (`PGR-E2004`); `LH-FA-17.a`.
- **Unbekannter Schlüssel, Schlüssel im falschen Abschnitt** — `PGR-E2004`; `LH-FA-17.a`.
  Schlüssel einer Option, die der Stand noch nicht kennt: unbekannt (§1).
- **Doppelter Schlüssel** — ungültiges YAML, `PGR-E2004`; `LH-FA-17.a`. Hinweis an den
  Implementer: Das Lesen über die Node-Schnittstelle der Bibliothek lehnt Doppelte nicht
  selbst ab.
- **Schlüssel `config` in einem Abschnitt** — gibt es nicht, unbekannter Schlüssel
  (`PGR-E2004`); Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.
- **Wertform in der Datei** (YAML `true` gegen `"true"`, `True`, `yes`, `1`; `shutdown_timeout: 0`
  als Zahl; leerer Wert, `null`, Liste oder Abbildung als Wert; Anker und Aliase) — Text des
  Skalars mit der Wertemenge der Option; leerer Wert, `null`, Liste, Abbildung, Anker, Aliase
  und Merge-Schlüssel `PGR-E2004`; Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.
- **Form der Dauer** — dieselbe in Option, Umgebungsvariable und Datei; ungültig `PGR-E2001`
  bzw. `PGR-E2004`; `LH-FA-17.a` (*Dauer*). `0` ohne Anführungszeichen ist gültig (Wertform
  oben).
- **`log_level` ungültig** — `PGR-E2004`; `LH-FA-17.a`.
- **Relative Pfade in der Datei** (`input`, `output`, später Zertifikate) — relativ zum
  aktuellen Verzeichnis; Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.

*Verbindungen und Platzhalter*

- **Name bei `--upstream`** — nur bei genauer Übereinstimmung; ein Wert, der weder Name noch
  `host:port` ist, ist ungültig (`PGR-E2001`, aus der Datei `PGR-E2004`); `LH-FA-17.a`.
- **`sslmode=require` bei `record`** — `PGR-E2004` der benutzten Verbindung; `LH-FA-17.a`.
- **Klartext-Passwort** — `PGR-E2006`, in jeder Verbindung der Datei; `LH-FA-17.a`.
- **Variable eines Platzhalters der benutzten Verbindung nicht gesetzt** — `PGR-E2005`;
  `LH-FA-17.a`.
- **Variable eines Platzhalters gesetzt, aber leer** — nicht gesetzt, `PGR-E2005`;
  Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.
- **Platzhalter in Teilen der URL, die `record` ignoriert** (Benutzer, Passwort, Datenbank) —
  unbeachtet; `record` löst nur Platzhalter in Host und Port auf; Entscheidung des Nutzers
  vom 2026-10-08, `LH-FA-17.a`.
- **`$${VAR}`, Rekursion, Form des Namens, Einsetzen in die URL** — `$$` steht für `$` in
  jedem Wert der Datei; Name `[A-Za-z_][A-Za-z0-9_]*`, sonst `PGR-E2004`; einmal eingesetzt;
  URL vor dem Einsetzen zerlegt, Wert unverändert in seinem Teil; Entscheidung des Nutzers
  vom 2026-10-08, `LH-FA-17.a`. Hinweis an den Implementer: Die Zerlegung der
  Standardbibliothek lehnt einen Port `${PORT}` ab; zerlegt wird mit Platzhaltern.

*Anzeige und Fehler*

- **`config show` und Geheimnisse** — keine Maskierung, und keine nötig: Ein Passwort steht
  in keiner gültigen Datei (`PGR-E2006`), Platzhalter erscheinen unaufgelöst, Umgebungs-
  variablen nur mit Namen; `LH-FA-17` (Lastenheft: *die Anzeige zeigt sie nicht*) und
  `LH-FA-17.a` (*Anzeige*). `LH-RB-01` gilt der Aufzeichnung (keine Maskierungsfunktion für
  deren Inhalt, `SPEC-033`), nicht der Anzeige der Konfiguration; kein Widerspruch.
- **`config show`: `PGR-E2005`, aktive Umgebungsvariablen** — kommt nicht vor; gesetzt und
  nicht leer, auch ohne passende Option, geprüft nur `PGWIRE_RECORDER_CONFIG`; `LH-FA-17.a`.
- **Form der Ausgabe von `config show`** — erste Zeile Pfad der Datei oder dass keine
  gefunden wurde, danach der Inhalt als YAML (zwei Leerzeichen, Reihenfolge der Datei, ohne
  Kommentare), danach die Namen der aktiven Umgebungsvariablen sortiert, alles auf `stdout`;
  Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.
- **Hilfe nach einer ungültigen Option, auch bei `config show` und `--config`** — die Hilfe
  geht vor; `LH-FA-01.a`.
- **Reihenfolge bei mehreren Fehlern** — Abbruch beim ersten Fehler; Reihenfolge
  Kommandozeile, Umgebungsvariablen nach der Tabelle, Datei in ihrer Reihenfolge, zuletzt
  Pflichtoptionen, Kombinationen, `--upstream` und `PGR-E2005`; Entscheidung des Nutzers vom
  2026-10-08, `LH-FA-17.a`.


*Rückgaben des Implementers vom 2026-10-08* (aus dem abgebrochenen Lauf an
`slice-v1-abschluss-konfiguration`), entschieden vom Architect am 2026-10-08 vor dem Code,
alle in `LH-FA-17.a`, abgeleitet aus den Grundsätzen dort: Text des Skalars zählt, Laden ohne
Kommando (`config show` prüft die Datei ohne Kommando, also prüft jedes Kommando sie ganz),
geschlossene Fehlertabelle, URL vor dem Einsetzen zerlegt, Abbruch beim ersten Fehler.

1. **Leere Datei oder nur Kommentare** — setzt nichts, kein Fehler; keine Ursache der
   Fehlertabelle trifft zu.
2. **Mehrere YAML-Dokumente** — ein zweites Dokument ist ungültiges YAML (`PGR-E2004`); es
   still zu übergehen, ließe Eingabe unbeachtet ([ADR-0014](../../adr/0014-konfigurationsdatei.md): jede unbekannte Eingabe ist ein
   Startfehler). Ein einzelnes `---` vor dem einen Dokument ist zulässig.
3. **Oberste Ebene keine Abbildung** — `PGR-E2004`.
4. **Leerer Abschnitt, leeres `connections:`** — `{}` setzt nichts; ohne Inhalt oder `null`
   ist `PGR-E2004` wie ein leerer Wert. Akzeptiertes Negativ: Ein Abschnitt, dessen Inhalt
   ganz auskommentiert ist, wird damit ungültig; das folgt der Strenge bei `null`
   (Entscheidung des Nutzers), und die Meldung nennt den Abschnitt.
5. **Abschnitt eines anderen Kommandos** — mitgeprüft, ebenso jede nicht benutzte Verbindung;
   vom Kommando hängen nur `sslmode=require` bei `record` und die Variablen der Platzhalter
   ab.
6. **URL ohne Port** — Port `5432`; die Grammatik lässt den Port weg. Ein Port sind Ziffern
   mit Wert 1 bis 65535, sonst `PGR-E2004`.
7. **Platzhalter-Syntax in ignorierten Teilen** — geprüft beim Laden, `PGR-E2004`;
   unbeachtet ist dort nur die Variable. Im Passwortteil ist ein fehlerhafter Platzhalter
   ein Klartext-Passwort (`PGR-E2006`).
8. **Eingesetzter Wert macht seinen Teil ungültig** — nur der Port wird nach dem Einsetzen
   geprüft: `PGR-E2004`, im letzten Schritt direkt nach den Variablen der benutzten
   Verbindung (`PGR-E2005`). Den Host prüft der Start nicht; ein Host mit `/` scheitert beim
   Verbindungsaufbau.
9. **YAML-Tags** — jeder ausdrücklich geschriebene Tag ist `PGR-E2004`, wie Anker und Aliase;
   ein Tag behauptet eine Bedeutung jenseits des Texts.
10. **URL-Sonderformen** — was die Grammatik nicht zulässt, ist `PGR-E2004`: Schema
    `postgres://`, leerer Host, fehlende oder leere Datenbank, Fragment, Parameter ohne `=`,
    ein Parameter zweimal (auch `sslmode`). Wörtliche Teile werden prozent-dekodiert, ein
    ungültiges Escape ist `PGR-E2004`; Platzhalter und `$$` gelten vor der Dekodierung, ein
    eingesetzter Wert wird nicht dekodiert. Prozent-Dekodierung fehlte in der Liste; sie
    folgt aus „Wert unverändert, auch mit `%`“ und ist mit entschieden. Innerhalb einer URL
    gilt die Reihenfolge ihrer Teile.
11. **Name einer Verbindung** — Form eines Werts (Text des Skalars; leer, `null`, Tag
    ungültig), dazu kein Steuerzeichen, weil die Meldung den Namen in einer Zeile nennt
    (`SPEC-034`); sonst jeder Name, auch mit Leerraum.
12. **`output` und `force` in der Datei** — Schlüssel des Abschnitts `record:` mit diesem
    Slice, weil `slice-v1-abschluss-konfiguration` beide am Leser anmeldet (Entscheidung des
    Architect vom 2026-10-08 zu F-496; sie ersetzt die frühere Antwort *unbekannt bis
    `slice-v1-abschluss-schreiben`*). `LH-FA-17.a` führt beide in der Tabelle, kein
    Sonderfall.

*Ort der Konstanten* — im Model, `internal/hexagon/model/fehler.go`, als drei Konstanten ohne
Logik (§1, Ausnahme). Grund: [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) führt eine Code-Tabelle im Quelltext, die mit dem
Katalog gleich bleibt, und der Bestand legt dort schon Codes ab, die nur ein Adapter erzeugt
(`PGR-E2001` im CLI-Adapter, `PGR-E4001` im PGWire-Adapter, `PGR-E3001` im Recording-Adapter); eine zweite Tabelle im
CLI-Adapter wäre eine zweite Quelle. Keine ADR nötig: Die Entscheidung folgt aus
[ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) und dem Bestand.

Aus `slice-v1-abschluss-konfiguration` zieht kein Risiko mit: Dessen Risiko der Einzel-Leser
bleibt dort, das zu [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) ist entfallen (die ADR ist `Accepted`).

**Risiken:**

- Die Datei-Hälfte ist auf 900 bis 1300 Zeilen geschätzt (Bericht des Implementers vom
  2026-10-08); das liegt an der Grenze einer Review-Sitzung. Gegenmittel: der vorab benannte
  Schnitt in §4 — **Ausgang:** offen bis Closure.
- Die Entscheidungen der zwölf Rückgaben ändern `LH-FA-17.a`; das Benutzerhandbuch beschreibt
  die Datei schon im Zielstand und kann abweichen (§3) — **Ausgang:** offen bis Closure.

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

- **Belege zur DoD (Implementer):** <je Zusage: Zusage · Mutation · roter Test (`AGENTS.md` §3.10)>
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `a968790` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen` (1×) — `slice-v1-abschluss-konfiguration` wuchs
  durch Übernahmen bei drei Liefer-Punkten, bis die Datei-Hälfte nicht mehr in eine
  Review-Sitzung passte; dieser Slice entsteht aus dem Schnitt. Der Beleg gehört in die
  Closure jenes Slice (dann 2×); hier ein Risiko in §6.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (14×, verkörpert in `AGENTS.md` §3.12)
  und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (4×, verkörpert) — der
  Implementer hielt an und gab zwölf Randformen zurück, statt sie im Code zu entscheiden;
  der Architect entschied sie vor dem Code (§6).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (15×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (21×, §3.11) — je Zusage eine
  Mutation, Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (17×, §3.9) — dieser Plan entsteht aus einem
  Schnitt; jede Korrektur zieht §1, §3 und §6 im selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — die Übernahme oben steht mit
  der Kennung des Gebers; Geber und die Nehmer `slice-v1-abschluss-schreiben` und
  `slice-v1-abschluss-einspielen` zeigen im selben Commit hierher.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor
dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
