# Slice slice-v1-abschluss-konfiguration: Konfiguration über Kommandozeile, Umgebung und Datei

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0014](../../adr/0014-konfigurationsdatei.md), [ADR-0027](../../adr/0027-yaml-bibliothek.md), [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) (Proposed)

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

**Ziel:** Jede Option von `record` und `replay` ist per Kommandozeile, Umgebungsvariable und Konfigurationsdatei setzbar, in der Priorität Kommandozeile vor Umgebungsvariable vor Datei vor Default; die Hilfe geht jeder Prüfung vor, und ein Fehler in jeder Quelle ist `PGR-E2001` oder `PGR-E2004` bis `PGR-E2006`.

**Übernimmt:** `slice-v1-abschluss-betrieb` — dessen Teil *Konfiguration* (DoD-Punkt 2 ohne
`--output` und `--force`, DoD-Punkt 3; Entscheidung des Nutzers vom 2026-10-08, Schnitt nach
F-345; dort §7 `Gegenstand:`). Im Einzelnen, je mit dem ursprünglichen Geber:

- **Aus `slice-replay-semantik-mismatch`:** der Schlüssel `fail_on_unconsumed` im Abschnitt
  `replay:`; die Prüfung jeder gesetzten Umgebungsvariable einer Option des Kommandos, auch
  wenn die Kommandozeile vorgeht (`PGR-E2001`, `LH-FA-17.a`), im allgemeinen Leser für alle
  Optionen; die Hilfe vor jeder Prüfung von Optionen, Umgebungsvariablen und
  Konfigurationsdatei (`LH-FA-01.a`), auch für `config show` und `--config`, samt `--` als
  Ende der Optionen.
- **Aus `slice-replay-semantik-meldungscodes`:** der Schlüssel `log_level` auf der obersten
  Ebene der Konfigurationsdatei (`LH-FA-17.a`), mit derselben Wertemenge wie `--log-level`;
  ein ungültiger Wert ist `PGR-E2004`.
- **Aus `slice-v1-abschluss-herunterfahren`** (dort §1, Abgrenzung): der Schlüssel der Frist
  `--shutdown-timeout` in der Konfigurationsdatei, und dass der allgemeine Leser auch
  `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` liest.

**Optionen, die der Stand kennt.** Der allgemeine Leser ist die einzige Stelle, an der eine
Option von `record` und `replay` angemeldet wird; Kommandozeile, Umgebungsvariable und
Schlüssel der Datei folgen aus dieser einen Anmeldung. Die Schlüssel, die die Datei annimmt,
sind die der Optionen, die der Stand kennt; ein Schlüssel einer Option aus der Tabelle in
`LH-FA-17.a`, die der Stand noch nicht kennt (etwa `format`, `tls_cert` oder der Abschnitt
`play:`), ist bis zu ihrem Slice ein unbekannter Schlüssel (`PGR-E2004`). Ein Test läuft über
alle angemeldeten Optionen (drei Quellen, Priorität), sodass eine später angemeldete Option
ihn ohne eigenen Plan-Punkt mitnimmt.

**Ort des Lesens.** Die Datei liest der CLI-Adapter (`ARC-005`). Die YAML-Bibliothek ist dort
heute nicht zulässig (`ARC-013`, [ADR-0027](../../adr/0027-yaml-bibliothek.md), `tech`-Regel in
`.a-check.yml`); [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md)
schlägt vor, sie dort zuzulassen. Bis sie `Accepted` ist und die Sicht nachgezogen ist, gibt
es keinen Code-Commit (§4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `--output`, `--force` und die strengen Werte für `--force` — `slice-v1-abschluss-schreiben`;
  er folgt diesem Slice und liest `--force` über den allgemeinen Leser von hier.
- Die Optionen von `play` — `slice-v1-abschluss-einspielen` (dort §1: `--fail-on-unconsumed`
  und `--log-level` bei `play`); er setzt auf dem allgemeinen Leser dieses Slice auf.
- Signale und die Frist selbst — `slice-v1-abschluss-herunterfahren`; hier nur ihre
  Quellen.
- TLS zum Upstream nach `sslmode` — `slice-v1-abschluss-einspielen`; nur `play` verbindet
  mit TLS zum Server. Hier wird `sslmode` geprüft, und `sslmode=require` ist bei `record`
  `PGR-E2004` (`LH-FA-17.a`).
- Code im Kern, im PGWire-Adapter und in den Driven-Adaptern — Schicht-Abgrenzung: Der
  Slice ändert den CLI-Adapter und den Bootstrap; dazu die `tech`-Regel des Architektur-Gates
  (Gate-Konfiguration, keine Schicht).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Für jede Option von `record` und `replay`, auch `--shutdown-timeout`, gilt
      die Priorität Kommandozeile vor Umgebungsvariable vor Konfigurationsdatei vor Default
      (`SPEC-007`); ein allgemeiner Leser prüft jede gesetzte Umgebungsvariable einer Option
      des Kommandos, auch wenn die Kommandozeile vorgeht (`PGR-E2001`), und ersetzt die
      Einzel-Leser von `--fail-on-unconsumed`, `--log-level` und `--shutdown-timeout` (Test).
- [ ] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Die Hilfe geht jeder Prüfung von Optionen,
      Umgebungsvariablen und Konfigurationsdatei vor, auch für `config show` und
      `--config`, samt `--` als Ende der Optionen; `config show` zeigt die gewählte Datei,
      ohne einen aufgelösten Wert (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Benannte Verbindungen (`connections`, `--upstream <Name>`, `sslmode`),
      `--config`, `PGWIRE_RECORDER_CONFIG`, die Standarddatei, der Schlüssel
      `fail_on_unconsumed` im Abschnitt `replay:`, der Schlüssel `log_level` auf der obersten
      Ebene (Wertemenge und Strenge von `--log-level`, ungültig `PGR-E2004`) und `$${VAR}`
      verhalten sich wie spezifiziert; eine ungültige Datei, eine nicht gesetzte Variable und
      ein Klartext-Passwort sind `PGR-E2004` bis `PGR-E2006` (Test). Beleg in §7 für alle
      drei Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
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
| `spec/architecture.md` | update (Architect, vor dem Code) | `ARC-013`, §2 *Zusätzliche Einschränkungen* und §6: YAML-Bibliothek auch in `ARC-005`, nach `Accepted` von [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) |
| `spec/spezifikation.md` | update (Architect, vor dem Code) | `LH-FA-17.a`: die Antworten des Nutzers auf die offenen Randformen in §6 |
| `.a-check.yml` | update | `tech`-Regel für `go.yaml.in/yaml` und `gopkg.in/yaml`: Recording-Adapter und CLI-Adapter ([ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md)) |
| `internal/adapters/driving/cli` | update | allgemeiner Leser, an dem jede Option einmal angemeldet wird, für Kommandozeile, Umgebungsvariable und Schlüssel der Datei; er prüft jeden gesetzten Wert und ersetzt die Einzel-Leser von `--fail-on-unconsumed`, `--log-level` und `--shutdown-timeout`; Priorität; Hilfe vor jeder Prüfung; Kommando `config show`; Wahl und Laden der Datei mit `connections`, `fail_on_unconsumed`, `log_level`, `${VAR}` und `$${VAR}`; Auflösen des Namens bei `--upstream` zu `host:port` |
| `internal/bootstrap` | update | `config show` ausführen (Ausgabe auf `stdout`); die zusammengeführten Optionen an die Use Cases geben |
| `internal/adapters/driving/cli` (Unit-Tests), `internal/bootstrap` (Tests), `test/integration` | update | Happy/Boundary/Negative nach `LH-FA-17.a` und `LH-FA-01.a`; ein Test über alle angemeldeten Optionen; die vorhandenen Tests der drei Einzel-Leser bleiben unverändert und grün (Risiko in §6) |
| `docs/user/benutzerhandbuch.md` | update | Konfiguration über drei Quellen und ihre Priorität |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-herunterfahren` liegt in `done/` (der allgemeine Leser
übernimmt dessen Umgebungsvariable). Schritt 3 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)
(Entscheidung des Nutzers vom 2026-10-08). Vor dem ersten Code-Commit prüft der Architect
die Randformen aus §6 und entscheidet die offenen (`AGENTS.md` §3.12).

**Erster Code-Commit** erst, wenn alle drei gelten: [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) ist `Accepted`, die Sicht
ist nachgezogen (§3), und jede Randform in §6, die dort *offen* steht, ist in `LH-FA-17.a`
entschieden.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar, oder eine Änderung im Upstream-Adapter oder im Kern wird nötig.
  Schnitt dann: *Leser* (DoD-Punkt 1 und die Hilfe aus DoD-Punkt 2, ohne Datei) und
  *Datei* (DoD-Punkt 3 und `config show`); der zweite setzt den ersten voraus.
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

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist; vom Architect am
2026-10-08 vor dem Code geprüft. *Offen* heißt: Frage an den Nutzer, Empfehlung des Architects
dahinter; der Implementer beginnt den betroffenen Teil erst, wenn die Antwort in `LH-FA-17.a`
steht (§4).

*Umgebung*

- **Gesetzte, aber leere Umgebungsvariable einer Option**, auch `PGWIRE_RECORDER_CONFIG` —
  gilt als nicht gesetzt; `LH-FA-17.a`.
- **Umgebungsvariable mit Präfix ohne passende Option** und die einer fremden Option —
  unbeachtet; `LH-FA-17.a`. `config show` nennt ihren Namen.
- **Ungültige Umgebungsvariable neben gesetzter Kommandozeile** — `PGR-E2001`; `LH-FA-17.a`.
- **Groß- und Kleinschreibung im Namen einer Umgebungsvariable** — folgt dem Betriebssystem;
  unter Windows unterscheidet es nicht. Akzeptiertes Negativ: keine Regel, kein Test, weil
  das Produkt die Umgebung nicht selbst liest, sondern über das Betriebssystem.

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
- **Schlüssel `config` in einem Abschnitt** — *offen.* Empfehlung: unbekannter Schlüssel
  (`PGR-E2004`); eine Datei nennt keine Datei.
- **Wertform in der Datei** (YAML `true` gegen `"true"`, `True`, `yes`, `1`; `shutdown_timeout: 0`
  als Zahl; leerer Wert, `null`, Liste oder Abbildung als Wert; Anker und Aliase) — *offen.*
  Empfehlung: Maßgeblich ist der Text des Skalars, geprüft mit derselben Wertemenge wie
  Option und Umgebungsvariable (`true` und `"true"` gültig, `True`, `yes` und `1` nicht;
  `0` und `"5s"` gültig); ein leerer Wert, `null`, eine Liste oder Abbildung an Stelle eines
  Werts und Anker, Aliase und Merge-Schlüssel sind `PGR-E2004`, wie bei der Aufzeichnung
  (`SPEC-001`).
- **Form der Dauer** — dieselbe in Option, Umgebungsvariable und Datei; ungültig `PGR-E2001`
  bzw. `PGR-E2004`; `LH-FA-17.a` (*Dauer*). Ob `0` als YAML-Zahl gilt, hängt an der Wertform
  oben.
- **`log_level` ungültig** — `PGR-E2004`; `LH-FA-17.a`.
- **Relative Pfade in der Datei** (`input`, `output`, später Zertifikate) — *offen.*
  Empfehlung: relativ zum aktuellen Verzeichnis wie auf der Kommandozeile, nicht zum
  Verzeichnis der Datei.

*Verbindungen und Platzhalter*

- **Name bei `--upstream`** — nur bei genauer Übereinstimmung; ein Wert, der weder Name noch
  `host:port` ist, ist ungültig (`PGR-E2001`, aus der Datei `PGR-E2004`); `LH-FA-17.a`.
- **`sslmode=require` bei `record`** — `PGR-E2004` der benutzten Verbindung; `LH-FA-17.a`.
- **Klartext-Passwort** — `PGR-E2006`, in jeder Verbindung der Datei; `LH-FA-17.a`.
- **Variable eines Platzhalters der benutzten Verbindung nicht gesetzt** — `PGR-E2005`;
  `LH-FA-17.a`.
- **Variable eines Platzhalters gesetzt, aber leer** — *offen.* Empfehlung: wie bei Optionen
  nicht gesetzt, also `PGR-E2005`.
- **Platzhalter in Teilen der URL, die `record` ignoriert** (Benutzer, Passwort, Datenbank) —
  *offen.* Empfehlung: unbeachtet wie die Teile selbst; `record` löst nur Platzhalter in
  Host und Port auf und meldet nur für sie `PGR-E2005`.
- **`$${VAR}`, Rekursion, Form des Namens, Einsetzen in die URL** — *offen*; `LH-FA-17.a`
  sagt nur, dass `$${VAR}` wörtlich bleibt. Empfehlung: `$$` steht für ein `$`, also ergibt
  `$${VAR}` den Text `${VAR}`; ein `$` vor anderem Zeichen bleibt stehen; dieselbe Lesart in
  jedem Wert der Datei. Der Name ist `[A-Za-z_][A-Za-z0-9_]*`; `${` ohne `}` oder mit anderem
  Namen ist `PGR-E2004`. Eingesetzt wird einmal, ein eingesetzter Wert wird nicht erneut
  ausgewertet. Die URL wird vor dem Einsetzen zerlegt, der Wert steht unverändert in seinem
  Teil, sodass ein Passwort mit `@`, `:`, `/` oder `%` trägt.

*Anzeige und Fehler*

- **`config show` und Geheimnisse** — keine Maskierung, und keine nötig: Ein Passwort steht
  in keiner gültigen Datei (`PGR-E2006`), Platzhalter erscheinen unaufgelöst, Umgebungs-
  variablen nur mit Namen; `LH-FA-17` (Lastenheft: *die Anzeige zeigt sie nicht*) und
  `LH-FA-17.a` (*Anzeige*). `LH-RB-01` gilt der Aufzeichnung (keine Maskierungsfunktion für
  deren Inhalt, `SPEC-033`), nicht der Anzeige der Konfiguration; kein Widerspruch.
- **`config show`: `PGR-E2005`, aktive Umgebungsvariablen** — kommt nicht vor; gesetzt und
  nicht leer, auch ohne passende Option, geprüft nur `PGWIRE_RECORDER_CONFIG`; `LH-FA-17.a`.
- **Form der Ausgabe von `config show`** — *offen*; `LH-FA-17.a` sagt *eingerückter Baum* und
  *nennt die Datei*. Empfehlung: erste Zeile nennt den Pfad der gewählten Datei, danach ihr
  Inhalt als YAML mit zwei Leerzeichen Einzug in der Reihenfolge der Datei, ohne Kommentare,
  danach die Namen der aktiven Umgebungsvariablen, je einer pro Zeile, nach Namen sortiert;
  ohne Datei eine Zeile, dass keine gefunden wurde, und die Namen. Alles auf `stdout`.
- **Hilfe nach einer ungültigen Option, auch bei `config show` und `--config`** — die Hilfe
  geht vor; `LH-FA-01.a`.
- **Reihenfolge bei mehreren Fehlern** — *offen.* Empfehlung: Der Start bricht beim ersten
  Fehler ab, mit einer Meldung. Geprüft wird in dieser Reihenfolge: Kommandozeile;
  Umgebungsvariablen der Optionen des Kommandos in der Reihenfolge der Tabelle in
  `LH-FA-17.a`; Datei (Wahl, Lesen, YAML, dann Schlüssel, Werte und Klartext-Passwörter in
  der Reihenfolge der Datei); zuletzt nach der Zusammenführung Pflichtoptionen, Kombinationen,
  `--upstream` und `PGR-E2005` der benutzten Verbindung.

**Risiken:**

- Der allgemeine Leser ersetzt drei Einzel-Leser; eine Prüfung eines Einzel-Lesers kann
  dabei wegfallen, statt in den allgemeinen überzugehen
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Gegenmittel: Die vorhandenen Tests der
  drei Einzel-Leser bleiben unverändert und grün (§3) — **Ausgang:** offen bis Closure.
- [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) wird nicht angenommen; dann liest der CLI-Adapter keine
  YAML-Datei, und der Slice geht nach `open` zurück (§4, blockiert), bis eine andere
  Entscheidung den Ort des Lesens trägt — **Ausgang:** offen bis Closure.

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `3864e44` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — der allgemeine Leser ersetzt
  Einzel-Leser; ein Risiko in §6.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (12×, verkörpert in `AGENTS.md`
  §3.12) — darum nennt §6 die Randformen, und der Architect prüft sie vor dem Code.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (14×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (19×, §3.11) — je Zusage eine
  Mutation, Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (17×, §3.9) — dieser Plan entsteht aus einem
  Schnitt; jede Korrektur zieht §1, §3 und §6 im selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — jede Übernahme oben steht
  mit der Kennung des Gebers; die Geber zeigen im selben Commit hierher.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor
dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
