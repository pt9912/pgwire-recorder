# Slice slice-v1-abschluss-verbindungen-platzhalter: Benannte Verbindungen und Platzhalter laden

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md), [ADR-0014](../../adr/0014-konfigurationsdatei.md)

**Berührte Spec-Stellen:** `LH-FA-17.a` · `SPEC-014` · `SPEC-034` · `ARC-005`

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

**Ziel:** Das Laden der Konfigurationsdatei prüft die benannten Verbindungen (`connections`) unabhängig vom Kommando: jede URL nach ihrer Grammatik samt Platzhaltern und `sslmode`, den Namen jeder Verbindung, `$$` und `${` in jedem Wert der Datei, den Schlüssel `upstream` gegen die Namen der Verbindungen und jedes Klartext-Passwort (`PGR-E2006`), auch bei `config show`.

**Übernimmt:** `slice-v1-abschluss-konfigurationsdatei` — dessen DoD-Punkt 3 (*Verbindungen
und Platzhalter*) nach dem vorab benannten Schnitt aus seinem §4 (*zu groß*, eingetreten am
2026-10-09; dort §1, *Abgegeben*). Dorthin kam der Punkt aus `slice-v1-abschluss-konfiguration`
und davor aus `slice-v1-abschluss-betrieb`. Nach dem Schnitt dieses Slice vom 2026-10-09
(§4, Entscheidung des Nutzers) bleibt davon hier das *Laden der Verbindungen*:

- die Zerlegung und die Grammatik der URL einer Verbindung samt der Form der Platzhalter in
  jedem Teil und `sslmode` (`disable`, `require`), in jeder Verbindung der Datei, auch einer
  nicht benutzten (`PGR-E2004`);
- der Name einer Verbindung ohne `:` und `@` (V-121, §6);
- `$$` in jedem Wert der Datei und ein `${` außerhalb einer URL (`PGR-E2004`);
- der Schlüssel `upstream` beim Laden gegen die Namen aller Verbindungen der Datei, sonst in
  der Form `host:port` (`PGR-E2004`), auch bei `config show`;
- `PGR-E2006` — ein Klartext-Passwort in jeder Verbindung, auch bei `config show`; die
  Konstante liegt seit jenem Slice ohne Erzeuger in der Code-Tabelle, dieser Slice erzeugt
  sie;
- die Randformen dazu aus §6 jenes Slice, die Entscheidungen des Architect vom 2026-10-08
  (Rückgaben 5, 6, 7 und 10, §6 unten) und die aus seiner Prüfung dieses Plans vom
  2026-10-09 vor dem Code, soweit sie das Laden betreffen (A1 bis A6, A7 und A8 für den
  Schlüssel, A10, A11; §6 unten), dazu die Rückgaben L1 bis L6 des Implementers vom
  2026-10-09 und das Gegenlesen dazu, die Entscheidungen zu F-520, F-521, F-522 und F-528
  aus dem Review (§6 unten).

**Abgegeben** an `slice-v1-abschluss-upstream-verbinden` (dort §1, *Übernimmt*, mit der
Kennung dieses Slice), nach dem vorab benannten Schnitt aus §4 (*zu groß*, eingetreten am
2026-10-09 vor dem Code): das *Benutzen der Verbindung* nach der Zusammenführung — Name und
Form `host:port` bei `--upstream` und `PGWIRE_RECORDER_UPSTREAM`, `sslmode=require` bei
`record`, das Einsetzen der Variablen, `PGR-E2005`, der Port nach dem Einsetzen, `record`
verbindet zu Host und Port der Verbindung, und der Teil des Benutzerhandbuchs zu
Verbindungen, `sslmode`, Platzhaltern, `PGR-E2005` und `PGR-E2006`; dazu die Randformen
Rückgabe 8, A7 und A8 für Option und Umgebung, A9 und A12.

**Aufsetzen.** Der Slice setzt auf dem Laden von `slice-v1-abschluss-konfigurationsdatei` auf:
Dort ist `connections:` eine Abbildung, ein Name einer Verbindung hat die Form eines Werts
ohne Steuerzeichen, und ein Wert ist ein Skalar (dort Rückgaben 4 und 11); bis zu diesem
Slice gilt ein Wert mit `$` wörtlich, und eine URL wird nicht zerlegt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Wahl, Laden und Form der Datei, die Priorität über drei Quellen, `config show`, die Hilfe
  und die Konstanten `PGR-E2004` bis `PGR-E2006` — `slice-v1-abschluss-konfigurationsdatei`;
  er liegt vor diesem Slice, und dieser setzt auf ihm auf.
- Das Benutzen der Verbindung nach der Zusammenführung (oben, *Abgegeben*) —
  `slice-v1-abschluss-upstream-verbinden`; er setzt die Zerlegung voraus, die dieser Slice
  liefert, weil eingesetzt wird, was das Laden zerlegt hat. Bis zu dessen Closure verbindet
  `record` mit einem Verbindungsnamen in `--upstream` nicht zu dessen Host und Port
  (Risiko in §6).
- Das Benutzerhandbuch zu Verbindungen und Platzhaltern —
  `slice-v1-abschluss-upstream-verbinden` (oben, *Abgegeben*); es beschreibt den Zielstand
  beider Slices und wird einmal nachgezogen, wenn er geliefert ist.
- Die Wirkung einer Verbindung bei `play` (Benutzer und Datenbank der URL mit dem Vorrang von
  `--user` und `--database`, das Passwort aus dem eingesetzten Platzhalter, TLS nach
  `sslmode`) — `slice-v1-abschluss-einspielen` (dort §1, *Übernommen aus*
  `slice-v1-abschluss-upstream-verbinden`); `play` gibt es erst mit ihm. Hier wird `sslmode`
  geprüft.
- TLS zum Upstream bei `record` — nicht Teil des Produkts in dieser Welle
  ([welle-v1-abschluss](../welle-v1-abschluss.md) §6); `sslmode=require` bei `record` prüft
  `slice-v1-abschluss-upstream-verbinden`.
- Code im Kern, im PGWire-Adapter, in den Driven-Adaptern und im Bootstrap —
  Schicht-Abgrenzung: Der Slice ändert nur das Laden im CLI-Adapter; die Code-Tabelle im
  Model bleibt unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Jede Verbindung der Datei, auch eine nicht benutzte, hat die Form der
      URL aus `LH-FA-17.a`, zerlegt mit Platzhaltern nach A1 (Schema, Benutzerteil, Host
      auch als IPv6 in eckigen Klammern, Port mit Default `5432` und Wert 1 bis 65535 wie
      geschrieben, Datenbank, Parameter, Prozent-Dekodierung wörtlicher Teile und der
      Parameter), der einzige Parameter `sslmode` nimmt `disable` und `require`, jede andere
      Form ist `PGR-E2004` mit dem Namen der Verbindung; die Form der Platzhalter prüft das
      Laden in jedem Teil. Ein Name einer Verbindung mit `:` oder `@` ist `PGR-E2004`, die
      Meldung nennt als Stelle nur `connections` (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `$$` steht in jedem Wert der Datei für `$`, von links gelesen; ein `${`
      außerhalb einer URL, das nicht aus `$$` hervorgeht, ist `PGR-E2004`; Kommandozeile und
      Umgebung kennen weder Platzhalter noch `$$`. Der Schlüssel `upstream` ist beim Laden an
      seiner Stelle der Name einer Verbindung der Datei, auch einer, die nach ihm steht, oder
      hat die Form `host:port`, sonst `PGR-E2004`, auch bei `config show` und einem anderen
      Kommando. `config show` zeigt Platzhalter unaufgelöst und `$$` wie geschrieben (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Ein Klartext-Passwort ist `PGR-E2006` in jeder Verbindung der Datei, auch
      bei `config show` — ein Passwortteil, der nicht genau ein `${VAR}` ist, auch ein leerer
      und ein fehlerhafter Platzhalter dort, und ein Parameter `password`, auch dekodiert;
      die Prüfung folgt der Reihenfolge aus `LH-FA-17.a` innerhalb einer URL (Test). Beleg in
      §7 für alle drei Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
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
| `spec/spezifikation.md` (`LH-FA-17.a`) | keine Änderung durch den Implementer | Die Randformen in §6 sind entschieden (Architect, 2026-10-08 und 2026-10-09, vor dem Code; die vom 2026-10-09 stehen seit dem Commit dieser Prüfung in `LH-FA-17.a`); eine Randform, die §6 nicht nennt, geht an den Architect zurück (`AGENTS.md` §3.12) |
| `internal/adapters/driving/cli` (Laden der Datei) | update | Zerleger der URL mit Platzhaltern (die Zerlegung der Standardbibliothek lehnt `${PORT}` ab), Grammatik, Prozent-Dekodierung, `sslmode`; Name einer Verbindung ohne `:` und `@`; `$$` in jedem Wert der Datei, `${` außerhalb einer URL; Schlüssel `upstream` gegen die Namen oder in der Form `host:port`; Klartext-Passwort `PGR-E2006`. Das Ergebnis der Zerlegung trägt, was `slice-v1-abschluss-upstream-verbinden` einsetzt (Teile, Platzhalter in der Reihenfolge der URL) |
| `internal/adapters/driving/cli` (Unit-Tests) | update | Happy/Boundary/Negative nach `LH-FA-17.a`, je Randform aus §6, die dieser Slice trägt, ein Fall |
| `test/integration` | update | `config show` mit Klartext-Passwort, ungültiger URL und ungültigem Schlüssel `upstream` endet mit Exit `2` und zeigt nichts |
| `docs/user/abdeckung-*.md` | update | über `make abdeckung` aus den Deklarationen der neuen Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-konfigurationsdatei` liegt in `done/`
(Laden der Datei, Form von `connections:`, Konstanten `PGR-E2005` und `PGR-E2006`). Schritt 5
der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) (Schnitt vom
2026-10-09). Die Randformen aus §6 entschied der Architect am 2026-10-08 vor dem Code; vor
dem ersten Code-Commit prüfte er die Liste am 2026-10-09 noch einmal gegen den Zuschnitt
(`AGENTS.md` §3.12).

**Erneuter Start** (`next` → `in-progress`) nach der Rückführung unten: dieser Plan im
geschnittenen Zuschnitt liegt auf dem Hauptzweig. Code gibt es noch keinen; die Arbeit ist
DoD-Punkt 1 bis 3.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar, oder eine Änderung im Kern wird nötig. Schnitt dann nach der Phase,
  nicht nach dem DoD-Punkt (Vorschlag des Architect vom 2026-10-09, §6 *Risiken*): *Laden der
  Verbindungen* und *Benutzen der Verbindung*. **Eingetreten am 2026-10-09** (Grund unten).
  Für den geschnittenen Zuschnitt: Der Zerleger allein hebt den Diff über eine
  Review-Sitzung; dann zurück an den Planner, geschnitten wird dann nicht weiter nach Phase,
  weil Zerlegung, Form der Platzhalter und Klartext-Passwort derselbe Zerleger sind.
- `in-progress` → `open` (blockiert — Carveout?): Eine Zusage aus `LH-FA-17.a` zu Verbindungen
  oder Platzhaltern widerspricht dem Laden von `slice-v1-abschluss-konfigurationsdatei` (etwa
  der Form eines Werts) und verlangt eine neue Entscheidung; dann zuerst die Entscheidung.

**Grund der Rückführung `in-progress` → `next`, eingetreten am 2026-10-09:** Der Architect
schätzte den Diff vor dem ersten Code-Commit auf 1050 bis 1300 Zeilen (§6 *Risiken*), so groß
wie `slice-v1-abschluss-konfigurationsdatei`, als er zur Zerlegung zurückging; das ist die
erste Bedingung oben, gemessen an der Schätzung statt am Diff. Der Nutzer entschied am
2026-10-09 den vorab benannten Schnitt nach der Phase: *Laden der Verbindungen* bleibt hier
(rund 650 bis 750 Zeilen), *Benutzen der Verbindung* geht an
`slice-v1-abschluss-upstream-verbinden` (§1, *Abgegeben*). Der Übergang läuft formal über
`next/` wie beim Geber (Entscheidung des Nutzers vom 2026-10-09 zu V-117), als reiner
`git mv` nach diesem Commit, dann `next` → `in-progress`.

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

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist. Alle zogen mit dem
Schnitt vom 2026-10-09 aus §6 von `slice-v1-abschluss-konfigurationsdatei` hierher; dorthin
kamen die ersten aus §6 von `slice-v1-abschluss-konfiguration`, geprüft vom Architect am
2026-10-08 vor dem Code. Die Rückgaben unten entschied er am 2026-10-08 vor dem Code, die
Nummern sind die aus jenem Plan. Offen ist keine; *Name einer Verbindung in der Form einer URL
mit Passwort* (unten) entschied der Nutzer am 2026-10-09 vor dem Code. Am 2026-10-09 prüfte der
Architect die Liste gegen den damaligen Zuschnitt und den Stand von `datei.go` und `cli.go` vor
dem ersten Code-Commit und entschied die fehlenden Randformen A1 bis A12 in `LH-FA-17.a`. Mit
dem Schnitt dieses Slice vom 2026-10-09 (§4) zogen die Randformen des Benutzens nach §6 von
`slice-v1-abschluss-upstream-verbinden`: Name bei `--upstream`, `sslmode=require` bei `record`,
die Variablen der benutzten Verbindung (nicht gesetzt, leer, in Teilen, die `record`
ignoriert), Einsetzen und Rekursion, `config show` ohne `PGR-E2005`, die Reihenfolge am Ende,
Rückgabe 8, A7 und A8 für Option und Umgebung, A9, A12 und das zweite akzeptierte Negativ.
Hier stehen die des Ladens. Die Rückgaben L1 bis L6 des Implementers vom 2026-10-09 und das
Gegenlesen entschied der Architect am selben Tag vor dem Code (unten).

*Verbindungen und Platzhalter beim Laden*

- **Klartext-Passwort** — `PGR-E2006`, in jeder Verbindung der Datei; `LH-FA-17.a`.
- **`$${VAR}`, Form des Namens** — `$$` steht für `$` in jedem Wert der Datei; Name
  `[A-Za-z_][A-Za-z0-9_]*`, sonst `PGR-E2004`; Entscheidung des Nutzers vom 2026-10-08,
  `LH-FA-17.a`. Hinweis an den Implementer: Die Zerlegung der Standardbibliothek lehnt einen
  Port `${PORT}` ab; zerlegt wird mit Platzhaltern.
- **Platzhalter außerhalb einer URL** — `PGR-E2004` (Tabelle *Fehler*); `LH-FA-17.a`.
- **Name einer Verbindung in der Form einer URL mit Passwort** (V-121 der Verifikation von
  `slice-v1-abschluss-konfigurationsdatei`) — ein Name enthält weder `:` noch `@`; sonst
  `PGR-E2004`, und die Meldung nennt als Stelle nur `connections`, nicht den Namen;
  Entscheidung des Nutzers vom 2026-10-09 nach Empfehlung des Architect, `LH-FA-17.a`
  (*Benannte Verbindungen*). Sie engt die Entscheidung zu Rückgabe 11 ein („sonst jeder Name,
  auch mit Leerraum“). Grund: Ohne `@` trägt ein Name keinen Benutzerteil und damit kein
  Passwort, ohne `:` ist er nie mit `host:port` bei `--upstream` verwechselbar.
  *Ort im Code:* dieser Slice, an der Prüfung des Namens neben der auf Steuerzeichen, die
  `slice-v1-abschluss-konfigurationsdatei` geliefert hat. Grund: Jener Slice ist geschlossen,
  und die Regel gehört zu den Verbindungen; bis zur Closure dieses Slice nimmt der Stand
  einen solchen Namen an (Risiko unten).
  *Tests* (`TestDateiUngueltig` bzw. der Test der Verbindungen): `connections:` mit den Namen
  `postgresql://u:GEHEIM@h/db`, `a:b` und `a@b`, je mit gültigem Wert, sind `PGR-E2004`, die
  Meldung beginnt mit `Konfigurationsdatei: connections:` und enthält weder den Namen noch
  `GEHEIM`; `config show` mit dieser Datei endet mit Exit `2` und zeigt nichts; die Namen
  `staging`, `a.b-c` und `mit Leerraum` bleiben gültig. *Mutation:* Prüfung auf `@` entfernt,
  rot über `a@b`; Prüfung auf `:` entfernt, rot über `a:b`; Stelle mit dem Namen, rot über
  `GEHEIM`.

*Anzeige und Fehler, soweit sie das Laden der Verbindungen betreffen*

- **`config show` und Platzhalter** — `config show` setzt keine Variable ein und zeigt
  Platzhalter unaufgelöst; `LH-FA-17.a`. Dass `config show` kein `PGR-E2005` meldet, sagt
  `slice-v1-abschluss-upstream-verbinden` zu, der das Einsetzen liefert.
- **`config show` und Geheimnisse** — keine Maskierung, und keine nötig: Ein Passwort steht in
  keiner gültigen Datei (`PGR-E2006`, mit diesem Slice), Platzhalter erscheinen unaufgelöst;
  `LH-FA-17` und `LH-FA-17.a` (*Anzeige*). Die Randform selbst bleibt in
  `slice-v1-abschluss-konfigurationsdatei`; hier kommt ihre Voraussetzung.
- **Reihenfolge, Teil des Ladens** — innerhalb einer URL in der Reihenfolge ihrer Teile, in
  der Reihenfolge der Datei; `LH-FA-17.a`. Den Teil am Ende nach der Zusammenführung liefert
  `slice-v1-abschluss-upstream-verbinden`, den übrigen Teil
  `slice-v1-abschluss-konfigurationsdatei`.

*Rückgaben des Implementers vom 2026-10-08*, entschieden vom Architect am 2026-10-08 vor dem
Code (`f19b440`), alle in `LH-FA-17.a`, abgeleitet aus den Grundsätzen dort: Text des Skalars
zählt, Laden ohne Kommando, geschlossene Fehlertabelle, URL vor dem Einsetzen zerlegt, Abbruch
beim ersten Fehler. Die Rückgaben 1 bis 4, 9, 11 und 12 liefert
`slice-v1-abschluss-konfigurationsdatei`, Rückgabe 8 `slice-v1-abschluss-upstream-verbinden`.

5. **Nicht benutzte Verbindung** (Teil der Verbindungen aus Rückgabe 5) — mitgeprüft, auch
   Grammatik, `sslmode` und Klartext-Passwort; vom Kommando hängen nur `sslmode=require` bei
   `record` und die Variablen der Platzhalter ab, beide bei
   `slice-v1-abschluss-upstream-verbinden`. Dass der Abschnitt eines anderen Kommandos
   mitgeprüft wird, liefert `slice-v1-abschluss-konfigurationsdatei`.
6. **URL ohne Port** — Port `5432`; die Grammatik lässt den Port weg. Ein Port sind Ziffern
   mit Wert 1 bis 65535, sonst `PGR-E2004`.
7. **Platzhalter-Syntax in ignorierten Teilen** — geprüft beim Laden, `PGR-E2004`;
   unbeachtet ist dort nur die Variable. Im Passwortteil ist ein fehlerhafter Platzhalter
   ein Klartext-Passwort (`PGR-E2006`).
10. **URL-Sonderformen** — was die Grammatik nicht zulässt, ist `PGR-E2004`: Schema
    `postgres://`, leerer Host, fehlende oder leere Datenbank, Fragment, Parameter ohne `=`,
    ein Parameter zweimal (auch `sslmode`). Wörtliche Teile werden prozent-dekodiert, ein
    ungültiges Escape ist `PGR-E2004`; Platzhalter und `$$` gelten vor der Dekodierung, ein
    eingesetzter Wert wird nicht dekodiert. Innerhalb einer URL gilt die Reihenfolge ihrer
    Teile.

*Prüfung des Architect vom 2026-10-09 vor dem Code* — je Randform, wo sie entschieden ist; neu
entschieden heißt: seit dem Commit dieser Prüfung in `LH-FA-17.a`, abgeleitet aus den
Grundsätzen dort (Grammatik geschlossen, URL vor dem Einsetzen zerlegt, jeder gesetzte Wert
geprüft, Laden ohne Kommando, nie ein Wert in der Meldung).

Bestätigt, schon vorher entschieden, soweit das Laden sie trägt: URL-Grammatik (*Benannte
Verbindungen*), Port-Default `5432` (Rückgabe 6), Prozent-Dekodierung wörtlicher Teile
(Rückgabe 10), `sslmode` mit `disable` und `require`, `${VAR}` und `$$`, `PGR-E2006` auch bei
`config show` (*Anzeige*), Name ohne `:` und `@` (V-121). Die übrigen Bestätigungen
(`sslmode=require` bei `record`, Auflösen des Namens, `PGR-E2005`, Port nach dem Einsetzen,
Passwort aus `${VAR}` mit `@` oder `:`) stehen bei `slice-v1-abschluss-upstream-verbinden`.
Ein Benutzer ohne Passwort ist gültig; woher `play` dann das Passwort nimmt, entscheidet
`LH-FA-17.a` (`PGWIRE_RECORDER_PASSWORD`) und liefert `slice-v1-abschluss-einspielen`.

Neu entschieden, soweit das Laden sie trägt:

- **A1 Zerlegung der URL** — von links; der Teil mit Benutzer, Host und Port reicht bis zum
  ersten `/`, `?` oder `#`; das letzte `@` darin trennt den Benutzerteil ab, das erste `:`
  im Benutzerteil das Passwort; die Datenbank reicht bis zum `?`, auch mit weiterem `/`;
  Parameter trennt `&`, ein `?` ohne Parameter und ein leerer Parameter sind Parameter ohne
  `=` (`PGR-E2004`).
- **A2 Leerer Benutzer, leeres Passwort** — leerer Benutzer (`postgresql://@h/db`,
  `postgresql://:${PW}@h/db`) ist `PGR-E2004`; ein leeres Passwort hinter `:`
  (`postgresql://u:@h/db`) ist ein Klartext-Passwort (`PGR-E2006`), weil es nicht genau ein
  `${VAR}` ist. Bestätigt durch den Nutzer am 2026-10-09 nach Empfehlung des Architect.
- **A3 IPv6-Host** — in eckigen Klammern, die Klammern gehören nicht zum Host; hinter `]`
  nur `:` mit Port oder das Ende; ohne Klammern endet der Host am ersten `:` (`::1` ohne
  Klammern ist damit ein ungültiger Wert, mehrere Hosts mit `,` ebenso über den Port). Das
  Zusammensetzen zu `host:port` für `record` liefert `slice-v1-abschluss-upstream-verbinden`.
- **A4 Port** — wie geschrieben, nicht dekodiert, führende Nullen erlaubt (`05432` ist
  5432), ein `:` ohne Port ist `PGR-E2004`.
- **A5 Dekodierung** — auch Name und Wert eines Parameters vor dem Vergleich
  (`sslmode=%72equire` ist `require`, `pass%77ord` ein Parameter `password`); ein Escape, das
  ein Steuerzeichen ergibt (`%00`), ist `PGR-E2004`.
- **A6 Platzhalter in den Parametern** — die Parameter prüft das Laden vor dem Einsetzen;
  `sslmode=${M}` ist ein ungültiger `sslmode` (`PGR-E2004`), `password=${PW}` ein
  Klartext-Passwort.
- **A7 Form `host:port`**, Teil des Schlüssels `upstream` — Host nicht leer (`:5432`
  ungültig), IPv6 in eckigen Klammern, Port in der Form aus A4; wie geschrieben, ohne
  Dekodierung. Dieselbe Form bei `--upstream` und `PGWIRE_RECORDER_UPSTREAM` prüft
  `slice-v1-abschluss-upstream-verbinden`; die Prüfung ist eine, beide benutzen sie.
- **A8 Schlüssel `upstream`** — geprüft beim Laden an seiner Stelle in der Datei gegen die
  Namen aller Verbindungen der Datei, auch einer, die nach ihm steht (`PGR-E2004`, auch bei
  `config show` und einem anderen Kommando). Option und Umgebungsvariable nach der
  Zusammenführung und die benutzte Verbindung trägt `slice-v1-abschluss-upstream-verbinden`.
- **A10 `$$` und `${` außerhalb einer URL** — von links gelesen (`$$${VAR}` ist `$` und ein
  Platzhalter); außerhalb einer URL ist jedes `${`, das nicht aus `$$` hervorgeht,
  `PGR-E2004` *Platzhalter außerhalb einer URL*, gleich ob seine Form gültig ist; die
  Wertemenge prüft den Text nach `$$`; `config show` zeigt `$$` wie geschrieben.
- **A11 Kommandozeile und Umgebung** — kennen weder Platzhalter noch `$$`
  (`--output='a$${X}'` ist der Pfad wie geschrieben). Als `VAR` gilt jeder Name der Form,
  auch `${PGWIRE_RECORDER_PASSWORD}`; ohne Ausnahme, weil die Variable nur gelesen wird.

*Rückgaben des Implementers vom 2026-10-09 (Laden)*, vor dem Code, entschieden vom Architect am
2026-10-09 in `LH-FA-17.a` (*Benannte Verbindungen*); Grundsätze: Grammatik geschlossen, Laden
prüft jede Verbindung, nie ein Wert in der Meldung, Steuerzeichen wie beim Namen. Tests in
`datei_test.go` (Test der Verbindungen bzw. `TestDateiUngueltig`), je Fall Code und Stelle
genau, die Meldung ohne den Wert.

- **L1 `$` im Namen einer Verbindung** — der Name ist kein Wert: Er gilt wörtlich, `$$` und
  Platzhalter gelten für ihn nicht, und er enthält kein `$` (`PGR-E2004`, Stelle nur
  `connections`, wie V-121). Der Schlüssel `upstream` ist ein Wert: `$$` gilt, dann der
  Vergleich; weil kein Name ein `$` trägt, trifft ein Wert mit `$` nie einen Namen, und
  `upstream: ${X}` ist ein Platzhalter außerhalb einer URL (A10). Grund: dieselbe Linie wie
  `:` und `@` (V-121), ohne sie hätte ein Name zwei Schreibweisen. Engt Rückgabe 11 weiter
  ein; Entscheidung des Architect, der Nutzer kann sie zurücknehmen (Bericht).
  *Test:* Namen `a$b`, `a$$b` und `${X}` je `PGR-E2004` mit Stelle `connections`, ohne den
  Namen; `upstream: a$$b` mit einer Verbindung `ab` ist `PGR-E2004` an `record.upstream`.
  *Mutation:* Prüfung auf `$` entfernt, rot über `a$b`; `$$` auf den Namen angewandt, rot
  über `a$$b`.
- **L2 Steuerzeichen im geschriebenen Text einer URL** — ungültiger Wert (`PGR-E2004`), auch
  der Tabulator, dieselbe Menge wie beim Namen (C0, `U+007F`, C1), geprüft vor der Zerlegung,
  also auch im Passwortteil `PGR-E2004` (wie das akzeptierte Negativ unten).
  *Test:* `"postgresql://h\x01/db"`, `"postgresql://h/d\tb"`, `"postgresql://u:p\x9fw@h/db"`
  je `PGR-E2004` an `connections.<Name>`. *Mutation:* Prüfung entfernt, rot über alle drei;
  Tabulator ausgenommen, rot über den zweiten.
- **L3 Prozent-Dekodierung** — (a) ein dekodierter Teil, der kein gültiges UTF-8 ist
  (`%FF`), ist ungültig; (b) `%C2%80` ergibt U+0080 aus C1 und ist ungültig; (c) dieselbe
  Menge wie beim Namen, auch `%09`. *Test:* Datenbank `%FF`, Benutzer `%C2%80`, Host `%09`
  je `PGR-E2004`; Datenbank `%C3%A4` (`ä`) gültig. *Mutation:* UTF-8-Prüfung entfernt, rot
  über `%FF`; nur C0 geprüft, rot über `%C2%80`.
- **L4 `upstream: a@b` vor `connections:` mit einer Verbindung `a@b`** — ein ungültiger Name
  zählt nicht als Name; der erste Fehler ist der, der in der Datei zuerst steht: hier
  `record.upstream` (weder Name noch `host:port`). Steht `connections:` zuerst, ist es
  `connections`. *Test:* beide Reihenfolgen, je Code und Stelle genau. *Mutation:* gegen
  alle Schlüssel statt gegen die gültigen Namen verglichen, rot über die erste Reihenfolge.
- **L5 Port mit Platzhalter und wörtlichem Text** — die Lesart des Implementers trägt nur
  halb: Beim Laden sind die wörtlichen Zeichen eines Ports mit Platzhalter Ziffern
  (`h:5${P}` gültig, `h:x${P}` und `h:$$${P}` `PGR-E2004`); Form und Wert als Ganzes prüft
  der Start nach dem Einsetzen (`slice-v1-abschluss-upstream-verbinden`). Grund: Das Laden
  prüft jede Verbindung, auch eine nicht benutzte, und `x${P}` wird nie ein Port.
  *Test:* die drei Fälle. *Mutation:* Ziffernprüfung bei Platzhalter entfernt, rot über
  `h:x${P}`.
- **L6 Passwort mit ungültigem Escape** — bestätigt: Der Passwortteil und der Wert des
  Parameters `password` werden nicht dekodiert, ob sie Klartext sind, entscheidet der
  geschriebene Text; `u:%ZZ@h` und `?password=%ZZ` sind `PGR-E2006`. Ein ungültiges Escape im
  Namen eines Parameters (`?pass%ZZ=x`) bleibt `PGR-E2004`. *Test:* die drei Fälle.
  *Mutation:* Passwort vor der Klartext-Prüfung dekodiert, rot über `u:%ZZ@h` (dann
  `PGR-E2004`).

*Gegengelesen, was der Implementer als entschieden las:*

- **`[abc` ohne `]`** — trug nicht; jetzt ausdrücklich: `[` ohne `]`, `[]` und `[` oder `]` an
  anderer Stelle des Hosts sind `PGR-E2004`. *Test:* `[abc`, `[]`, `a]b`. *Mutation:*
  Klammerprüfung entfernt, rot über `[abc`.
- **`SSLMODE` und `PASSWORD`** — bestätigt, jetzt ausdrücklich: Parameternamen gelten genau
  in der Schreibweise, beide sind unbekannte Parameter (`PGR-E2004`), `PASSWORD=x` also nicht
  `PGR-E2006`; beide beenden den Start ohne den Wert. *Test:* `?SSLMODE=disable`,
  `?PASSWORD=x` je `PGR-E2004`. *Mutation:* Vergleich ohne Groß- und Kleinschreibung, rot über
  beide.
- **`$${VAR}` im Passwort** — bestätigt: Der geschriebene Text ist nicht genau ein `${VAR}`,
  also Klartext (`PGR-E2006`). *Test:* `u:$${PW}@h`. *Mutation:* `$$` vor der Klartext-Prüfung
  ausgewertet, rot.
- **Prüfreihenfolge in der URL** — jetzt ausdrücklich: Steuerzeichen, Schema, Benutzer,
  Passwort, Host, Port, Datenbank, die Parameter in ihrer Reihenfolge, Fragment; je Teil die
  Form seiner Platzhalter vor dem Rest. *Test:* `postgres://:p@/` ist das Schema;
  `postgresql://u:p@[x/db?a=b` ist das Passwort (`PGR-E2006`) vor dem Host.
  *Mutation:* Host vor dem Passwort geprüft, rot über den zweiten Fall.

*Befunde des Reviews vom 2026-10-09 an den Architect* (`docs/reviews/2026-10-09-review-slice-v1-abschluss-verbindungen-platzhalter.md`),
entschieden vom Architect am 2026-10-09 in `LH-FA-17.a` (*Benannte Verbindungen*,
*Geheimnisse*); Grundsätze: Grammatik geschlossen, Laden prüft, was ohne Einsetzen prüfbar
ist, ein eingesetzter Wert bleibt ungeprüft. Je Fall Code, Stelle und Grund genau, die Meldung
ohne den Wert.

- **F-521 Inhalt des Hosts.**
  - *URL, in Klammern:* wörtlich, ohne Platzhalter, eine IPv6-Adresse in Textform (RFC 4291),
    auch mit Zone hinter `%25` (RFC 6874); IPv4-Adresse, Name oder Platzhalter in Klammern
    sind `PGR-E2004`. Grund: `[${H}]` ist vor dem Einsetzen nicht prüfbar, ein IPv6-Host aus
    einer Variable steht ohne Klammern (`${H}`) und wird beim Zusammensetzen geklammert.
  - *URL, ohne Klammern:* nach der Dekodierung kein Leerraum (Unicode `White_Space`) und keines
    von `@`, `:`, `/`, `?`, `#`, `[`, `]`, `%`; Steuerzeichen schon nach L2 und L3. Mit
    Platzhalter prüft das Laden die wörtlichen Zeichen, den eingesetzten Wert nicht. Ob ein
    Name auflöst oder ein gültiger DNS-Name ist, prüft der Start nicht (akzeptiertes Negativ:
    `a..b` lädt und scheitert beim Verbinden, `PGR-E4002`; eine DNS-Grammatik wäre eine zweite,
    strengere Regel als die des Resolvers).
  - *`host:port` (Schlüssel `upstream`, im Folge-Slice auch Option und Umgebung):* wie
    geschrieben, ohne Dekodierung; dieselben Regeln, Zone hinter `%`; dazu kein Steuerzeichen.
  - *Beim Laden geprüft:* alles oben. *Erst nach dem Einsetzen, im Folge-Slice:* nichts am
    Host, nur das Zusammensetzen (F-528).
  - *Test:* URL `[::1]`, `[fe80::1%25eth0]`, `[::ffff:1.2.3.4]` gültig; `[abc]`, `[a b]`,
    `[GEHEIM]:5`, `[1.2.3.4]`, `[${H}]`, `a%3Ab`, `h%2Fx`, `a%20b`, `a%25b` je `PGR-E2004`
    ohne `GEHEIM`; `db-${N}.example` gültig, `a b${N}` ungültig. `upstream` `"[::1]:5"`,
    `"h:5"` gültig; `"GEHEIM@h:5"`, `"GEHEIM h:5"`, `"h\x01:5"`, `"[GEHEIM]:5"`, `"a/b:5"`
    je `PGR-E2004` an `record.upstream` ohne `GEHEIM`.
  - *Mutation:* IPv6-Prüfung in Klammern entfernt, rot über `[abc]`; Platzhalter in Klammern
    zugelassen, rot über `[${H}]`; Prüfung der Zeichen ohne Klammern nach der Dekodierung
    statt davor, rot über `a%3Ab`; Leerraum zugelassen, rot über `a%20b` und `"GEHEIM h:5"`;
    `@` in `host:port` zugelassen, rot über `"GEHEIM@h:5"`; Steuerzeichen in `host:port`
    zugelassen, rot über `"h\x01:5"`.
- **F-522 Zwei Fälle von `password`.** Je Parameter zuerst der Name (Platzhalter, Escape),
  dann ist `password` ein Klartext-Passwort, auch ohne `=` und mit leerem Wert; danach `=`,
  bekannt, doppelt, Wert. `?password` ist damit `PGR-E2006`. `?pass${X}word=x` ist
  `PGR-E2004` (Platzhalter im Namen, F-520), kein Klartext: Der Name ist nicht `password`.
  *Test:* `?password`, `?password=` je `PGR-E2006`; `?pass${X}word=x` `PGR-E2004`.
  *Mutation:* `=` vor dem Namen geprüft, rot über `?password`; Platzhalter im Namen wie Text
  behandelt, rot über `?pass${X}word=x` (dann unbekannt mit anderem Grund).
- **F-520 Platzhalter neben wörtlichem Text in Parametern.** Ein Platzhalter an irgendeiner
  Stelle des Werts von `sslmode`, auch neben Text, ist ein ungültiger `sslmode`
  (`?sslmode=re${X}quire` ist `PGR-E2004`); ein Platzhalter an irgendeiner Stelle im Namen
  eines Parameters ist ein ungültiger Wert (`?ssl${X}mode=disable`, `?pass${X}word=x`). Den
  Rest von F-520 (Anker in `platzhalterName`) setzt der Implementer um. *Test:* die drei
  Fälle, je Grund genau. *Mutation:* M5 und M6 aus dem Review, beide rot.
- **F-528 Klammer in der Ablage** — braucht das Laden nicht: Nach F-521 enthält ein Host ohne
  Klammern kein `:`, auch nicht nach der Dekodierung (`a%3Ab` ist ungültig); ein abgelegter
  Host mit `:` stand also in Klammern oder kommt aus einer Variable. Das Zusammensetzen
  klammert jeden Host mit `:` und ist damit ohne Merkmal eindeutig. Als Randform eingetragen
  in §6 von `slice-v1-abschluss-upstream-verbinden`.

*Akzeptiertes Negativ der Prüfung* (keine Folgepflicht, einmalig und harmlos):

- Ein `/`, `?` oder `#` wörtlich im Passwort beendet den Teil mit Benutzer, Host und Port;
  die URL ist dann `PGR-E2004` statt `PGR-E2006`. Beide beenden den Start, keine Meldung
  nennt den Wert, `config show` zeigt nichts.

*Grenze der Tests* (Implementer, 2026-10-09; zwei äquivalente Mutanten, §7):

- Ob die Wertemenge einer Option den Text vor oder nach `$$` prüft, ist über die
  Schnittstelle nicht unterscheidbar, solange keine Wertemenge ein `$` annimmt (`text` nimmt
  jeden nicht leeren Wert, `$$` ergibt nie einen leeren); gesichert ist, dass die Option den
  Text nach `$$` erhält.
- `$$` auf einen Namen einer Verbindung angewandt ändert nichts, solange ein Name kein `$`
  enthält (L1): Das Ergebnis enthält wieder ein `$` und wird abgelehnt.

**Risiken:**

- Die Zerlegung der URL mit Platzhaltern braucht einen eigenen Zerleger, weil die
  Standardbibliothek `${PORT}` als Port ablehnt (Hinweis oben); Grammatik, Dekodierung und
  Einsetzen in einem Zerleger hoben den Diff über eine Review-Sitzung. *Schätzung des
  Architect vom 2026-10-09* gegen den Stand von `datei.go` (551 Zeilen) und `cli.go` und den
  Liefer-Commit des Gebers (`7a80393`, 1392 Zeilen, davon rund 590 Code und 665 Tests): Code rund
  350 bis 450 Zeilen (Zerleger mit Platzhaltern 200, `$$` und `${` in jedem Wert 40,
  Klartext 30, `host:port` und Auflösen 60, Einsetzen und Prüfungen am Ende 60), Unit-Tests
  550 bis 700 (je Randform aus §6 ein Fall, zwölf neue), Integration 80, Handbuch und
  Abdeckung 70 — zusammen **1050 bis 1300 Zeilen**, so groß wie der Geber, als er zur
  Zerlegung zurückging; der Fall von `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`.
  Entscheidung des Nutzers vom 2026-10-09: vor dem Code nach §4 geschnitten —
  **Ausgang:** eingetreten, Folge-Slice `slice-v1-abschluss-upstream-verbinden`.
- Auch der geschnittene Zuschnitt (*Laden*, nach der Schätzung rund 650 bis 750 Zeilen,
  davon der Zerleger mit Platzhaltern rund 200 Zeilen Code) kann über einer Review-Sitzung
  liegen, weil die Randformen des Ladens die meisten Testfälle tragen (A1 bis A6, A10,
  Rückgaben 6, 7 und 10) — **Ausgang:** offen bis Closure.
- Die Form, in der das Laden die zerlegte URL ablegt, ist die Schnittstelle zu
  `slice-v1-abschluss-upstream-verbinden`: Trägt sie die Platzhalter je Teil und in der
  Reihenfolge der URL nicht (A12 dort), muss der Folge-Slice den Zerleger ändern —
  **Ausgang:** offen bis Closure.
- Bis zur Closure dieses Slice lehnt der gelieferte Stand ein Klartext-Passwort in der Datei
  und einen Namen einer Verbindung mit `:` oder `@` nicht ab, und `config show` zeigt beide (Risiko in §6 von
  `slice-v1-abschluss-konfigurationsdatei`); DoD-Punkt 3 schließt die Lücke beim Passwort,
  DoD-Punkt 1 die beim Namen —
  **Ausgang:** offen bis Closure.

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

- **Belege zur DoD (Implementer):** siehe *Belege des Implementers* unten.
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

### Belege des Implementers

**Stand.** Geliefert sind DoD-Punkt 1 bis 3 in vier Commits: `f0180b6` (DoD-Punkt 1),
`d7e1833` (DoD-Punkt 2), `f9b9584` (DoD-Punkt 3 mit dem Integrationstest), `a3bcad0` (ein
Testfall zu L1 und L4, der einen Mutanten fängt). Neu ist
`internal/adapters/driving/cli/verbindung.go` (Zerleger, `host:port`, Klartext), geändert
`datei.go` (Verbindungen, `$$` und `${`, Schlüssel `upstream`), Tests in
`verbindung_test.go`, `leser_test.go` (Werte der Wertemenge `text` in der Form `host:port`,
weil `record.upstream` in der Datei sie verlangt), `export_test.go` und
`test/integration/konfiguration_e2e_test.go`. Keine Randform außerhalb von §6 entschieden;
die sechs Rückgaben dieses Laufs (L1 bis L6) entschied der Architect vor dem ersten
Code-Commit (`7a4dc1c`). Das Handbuch ist nicht berührt (§1, *Abgegeben*).

**Größe.** Der Diff `539ebee..a3bcad0` ohne Plan umfasst 991 hinzugefügte und 19 entfernte
Zeilen: rund 530 Code (davon `verbindung.go` 440 mit Kommentaren), rund 460 Tests, 12 Zeilen
Abdeckung. Das liegt über der Schätzung von 650 bis 750 Zeilen (§6 *Risiken*, zweiter
Punkt); ob es in eine Review-Sitzung passt, urteilt das Review.

**Schnittstelle zum Folge-Slice** (§6 *Risiken*, dritter Punkt). `datei.verbindungen` trägt je
Verbindung in der Reihenfolge der Datei den Namen und je Teil (Benutzer, Passwort, Host, Port,
Datenbank) eine Folge von Stücken in der Reihenfolge der URL: wörtlicher Text nach `$$` und,
außer im Port, dekodiert, oder ein Platzhalter mit dem Namen seiner Variable; dazu `sslmode`,
ohne Port das Stück `5432`. `TestVerbindungZerlegung` prüft diese Form über `Zerlege` in
`export_test.go`.

**Läufe.**

| Lauf | Stand | Ergebnis |
|---|---|---|
| `make test` | `f0180b6`, `d7e1833`, `f9b9584`, `a3bcad0` je vor dem Commit | Exit 0 |
| `make lint` | `f9b9584`, `a3bcad0` | Exit 0, kein Befund |
| `make test-integration` | `f9b9584` | Exit 0, `TestE2EConfigShowVerbindungUngueltig` PASS |
| `make abdeckung` | je Commit | Tabellen nachgezogen, `make abdeckung-check` im Gate-Lauf grün |
| `make gates` | `a3bcad0` (Code-Stand) | Exit 0 |

**Weg der Mutanten.** Je Mutant ein neues Verzeichnis unter dem Scratch-Pfad der Sitzung, in
das `go.mod`, `go.sum`, `cmd/` und `internal/` mit `shutil.copyfile` kopiert werden (keine
mtime übernommen, kein `cp -p`); die Änderung setzt ein Skript, das die Trefferzahl 1 prüft;
dann `go test -count=1 ./internal/adapters/driving/cli/` im Image der Stufe `deps` des
`Dockerfile`, die Kopie per Bind-Mount statt Build-Kontext, ohne Netz; danach wird genau
dieses Verzeichnis gelöscht. Die drei Mutanten des Integrationstests liefen in je einer Kopie
des Repos (`cp -r`, ohne `-p`) mit `make test-integration`, danach diese drei Verzeichnisse
gelöscht. Der Arbeitsbaum blieb unberührt. Gefahren: 85 Mutanten, 83 rot, 2 grün (äquivalent,
unten). Ein Mutant, der nicht übersetzte, ist neu formuliert und dann rot gesehen.

**DoD-Punkt 1 — Zerlegung und Grammatik der URL, Name einer Verbindung.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| Steuerzeichen im geschriebenen Text ungültig (L2) | Prüfung entfernt | `TestVerbindungUngueltig` |
| auch der Tabulator (L2) | Tabulator ausgenommen | `TestVerbindungUngueltig` |
| auch DEL und C1 (L2) | nur C0 geprüft | `TestVerbindungUngueltig` |
| Schema genau `postgresql://` | Prüfung entfernt | `TestVerbindungUngueltig`, `TestVerbindungName` |
| Benutzerteil am letzten `@` (A1) | erstes `@` | `TestVerbindungZerlegung` |
| leerer Benutzer ungültig (A2) | Prüfung entfernt | `TestVerbindungUngueltig` |
| IPv6 in eckigen Klammern ohne die Klammern (A3) | Klammern nicht erkannt | `TestVerbindungUngueltig`, `TestVerbindungZerlegung` |
| `[` ohne `]` ungültig | offene Klammer als Host angenommen | `TestVerbindungUngueltig` |
| `[` oder `]` an anderer Stelle des Hosts ungültig | Prüfung entfernt | `TestVerbindungUngueltig` |
| hinter `]` nur `:` mit Port oder das Ende (A3) | Prüfung entfernt | `TestVerbindungUngueltig` |
| leerer Host ungültig, auch `[]` | Prüfung entfernt | `TestVerbindungUngueltig` |
| ohne Port `5432` (Rückgabe 6) | Default `5433` | `TestVerbindungZerlegung` |
| `:` ohne Port ungültig (A4) | Prüfung entfernt (dann anderer Grund) | `TestVerbindungUngueltig` |
| Port höchstens 65535 | Grenze 65536 | `TestVerbindungUngueltig` |
| Port mindestens 1 | Grenze 0 | `TestVerbindungUngueltig` |
| führende Nullen erlaubt (A4) | führende Null abgelehnt | `TestVerbindungZerlegung` |
| Port nur Ziffern | Ziffernprüfung entfernt | `TestVerbindungUngueltig` (Fall `h:+5`) |
| Port nicht dekodiert (A4) | Port dekodiert | `TestVerbindungUngueltig` (Fall `h:%35`) |
| Port mit Platzhalter: wörtliche Zeichen Ziffern (L5) | Prüfung entfernt | `TestVerbindungUngueltig` (Fälle `h:x${P}`, `h:$$${P}`) |
| Datenbank fehlt: ungültig | Prüfung entfernt | `TestVerbindungUngueltig`, `TestVerbindungName` |
| Datenbank leer: ungültig | Prüfung entfernt | `TestVerbindungUngueltig` |
| Datenbank bis `?`, auch mit `/` (A1) | Datenbank endet am `/` | `TestVerbindungZerlegung` |
| wörtliche Teile prozent-dekodiert (Rückgabe 10) | keine Dekodierung | `TestVerbindungUngueltig`, `TestVerbindungZerlegung` |
| ungültiges Escape ungültig | Hex-Prüfung entfernt | `TestVerbindungUngueltig` |
| dekodiert gültiges UTF-8 (L3) | Prüfung entfernt | `TestVerbindungUngueltig` (Fall `%FF`) |
| dekodiert ohne C0, DEL, C1 (L3, A5) | nur C0 geprüft | `TestVerbindungUngueltig` (Fall `%C2%80`) |
| Parameter ohne `=` ungültig, auch `?` und leerer | Prüfung entfernt | `TestVerbindungUngueltig` |
| unbekannter Parameter ungültig | übergangen | `TestVerbindungUngueltig` |
| Parametername genau in der Schreibweise | ohne Groß- und Kleinschreibung | `TestVerbindungUngueltig` (Fall `SSLMODE`) |
| `sslmode` höchstens einmal | Prüfung entfernt | `TestVerbindungUngueltig` |
| Parametername dekodiert (A5) | Name nicht dekodiert | `TestVerbindungUngueltig`, `TestVerbindungZerlegung` |
| Platzhalter oder leerer Wert bei `sslmode` ungültig (A6) | beide angenommen | `TestVerbindungUngueltig` |
| `sslmode` nur `disable` und `require` | `prefer` angenommen | `TestVerbindungUngueltig` |
| Wert genau in der Schreibweise | `Require` angenommen | `TestVerbindungUngueltig` |
| `sslmode` übernommen | Wert nicht gemerkt | `TestVerbindungZerlegung` |
| Fragment ungültig | Prüfung entfernt | `TestVerbindungUngueltig` |
| Form der Platzhalter in jedem Teil (Rückgabe 7) | ungültiges `${` als Text | `TestVerbindungUngueltig` |
| `$$` ist `$` in der URL | `$$` bleibt `$$` | `TestVerbindungZerlegung` |
| `$$` von links | `$$` nicht ausgewertet | `TestVerbindungZerlegung` |
| Platzhalter je Teil in der Reihenfolge der URL (Schnittstelle) | Text vor einem Platzhalter verloren | `TestVerbindungZerlegung` |
| Passwort als Platzhalter abgelegt | Passwort als Text | `TestVerbindungZerlegung` |
| jede Verbindung geprüft, auch eine nicht benutzte | Ergebnis von `zerlegeURL` übergangen | `TestVerbindungUngueltig`, `TestVerbindungName` |
| Name ohne `@` (V-121) | `@` erlaubt | `TestVerbindungName` |
| Name ohne `:` (V-121) | `:` erlaubt | `TestVerbindungName` |
| Name ohne `$` (L1) | `$` erlaubt | `TestVerbindungName` |
| Meldung nennt nur `connections`, nicht den Namen | Name in der Stelle | `TestVerbindungName` |

**DoD-Punkt 2 — `$$` und `${` in jedem Wert, Schlüssel `upstream`.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| `$$` ist `$` in jedem Wert | `$$` bleibt `$$` | `TestDateiDollar`, `TestDateiUpstream` |
| `${` außerhalb einer URL ungültig, gleich ob gültig geformt (A10) | angenommen | `TestDateiPlatzhalterAusserhalb` |
| von links gelesen: `$$${X}` ist `$` und ein Platzhalter (A10) | `${` nach `$` übergangen | `TestDateiPlatzhalterAusserhalb` |
| die Option erhält den Text nach `$$` | Text wie geschrieben gemerkt | `TestDateiDollar`, `TestDateiUpstream` |
| Kommandozeile ohne `$$` (A11) | `$$` auf die Kommandozeile angewandt | `TestDateiDollar` |
| Umgebung ohne `$$` (A11) | `$$` auf die Umgebung angewandt | `TestDateiDollar` |
| Schlüssel `upstream` geprüft (A8) | Prüfung entfernt | `TestDateiUpstreamUngueltig` |
| `upstream` darf eine Verbindung nennen | nur `host:port` | `TestDateiUpstream`, `TestDateiUpstreamUngueltig` |
| `upstream` darf `host:port` sein | nur Namen | `TestDateiUpstream`, `TestDateiDollar`, `TestLeserAlleOptionen` |
| nur gültige Namen zählen (L4) | jeder Schlüssel unter `connections` zählt | `TestDateiUpstreamUngueltig` |
| auch eine Verbindung nach dem Schlüssel (A8) | keine Namen gesammelt | `TestDateiUpstream`, `TestDateiUpstreamUngueltig` |
| Name genau in der Schreibweise | Name auch groß geschrieben | `TestDateiUpstreamUngueltig` |
| ein Name mit `$` zählt nicht, auch nach `$$` (L1, L4) | `$$` beim Sammeln entfernt | `TestDateiUpstreamUngueltig` |
| `host:port`: Host nicht leer (A7) | leerer Host angenommen | `TestDateiUpstreamUngueltig` |
| `host:port`: `[` und `]` nur als Klammern | Prüfung entfernt | `TestDateiUpstreamUngueltig` |
| `host:port`: Port der Form aus A4 | jeder nicht leere Port | `TestDateiUpstreamUngueltig` |
| `host:port`: IPv6 in eckigen Klammern | Klammern nicht erkannt | `TestDateiUpstream` |
| `host:port`: Host endet am ersten `:` | am letzten `:` | `TestDateiUpstreamUngueltig` |
| `host:port`: wie geschrieben, ohne Dekodierung | Port dekodiert | `TestDateiUpstreamUngueltig` |
| auch bei `config show` und `replay` | durch die Prüfung oben mitgedeckt: dieselben Mutanten färben die Fälle `config show` und `replay` in `TestDateiUpstreamUngueltig` | `TestDateiUpstreamUngueltig` |

**DoD-Punkt 3 — Klartext-Passwort.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| Passwort nicht genau `${VAR}` ist Klartext | Prüfung entfernt | `TestVerbindungKlartext` |
| Code `PGR-E2006` | Code `PGR-E2004` | `TestVerbindungKlartext` |
| auch leeres Passwort (A2) | leeres angenommen | `TestVerbindungKlartext` |
| genau ein `${VAR}`, nicht nur enthalten | Muster nicht verankert | `TestVerbindungKlartext` |
| `$${VAR}` ist Klartext | `$$` vor der Prüfung ausgewertet | `TestVerbindungKlartext` |
| Passwort nicht dekodiert, ungültiges Escape ist Klartext (L6) | Passwort vor der Prüfung dekodiert | `TestVerbindungKlartext` |
| Passwort am ersten `:` (A1) | am letzten `:` | `TestVerbindungKlartext` |
| Parameter `password` ist Klartext | Prüfung entfernt | `TestVerbindungKlartext` |
| auch dekodiert geschrieben (A5) | Name wie geschrieben verglichen | `TestVerbindungKlartext` |
| nur genau `password` (Gegenlesen) | ohne Groß- und Kleinschreibung | `TestVerbindungKlartext`, `TestVerbindungUngueltig` |
| Wert von `password` nicht dekodiert (L6) | Wert dekodiert, Escape-Fehler als `PGR-E2004` | `TestVerbindungKlartext` |
| Passwort vor dem Host (Reihenfolge) | Host vor dem Passwort | `TestVerbindungKlartext` |
| Passwort vor Datenbank, Parametern und Fragment | Pfad vor dem Benutzerteil | `TestVerbindungKlartext` |
| Benutzer vor dem Passwort | Klartext vor dem Benutzer | `TestVerbindungKlartext` |
| Steuerzeichen vor dem Klartext (L2) | Klartext vor den Steuerzeichen | `TestVerbindungKlartext` |
| `config show` endet mit `PGR-E2006`, Exit 2, ohne Ausgabe (Integration) | Code `PGR-E2004` | `TestE2EConfigShowVerbindungUngueltig` |
| `config show` mit ungültiger URL endet mit Exit 2 (Integration) | Ergebnis von `zerlegeURL` übergangen | `TestE2EConfigShowVerbindungUngueltig` |
| `config show` mit ungültigem `upstream` endet mit Exit 2 (Integration) | Prüfung von `upstream` entfernt | `TestE2EConfigShowVerbindungUngueltig` |

**Grüne Mutanten, eingeordnet.**

- *Wertemenge prüft den Text vor `$$`* — äquivalent: Keine Wertemenge nimmt ein `$` an, und
  `$$` ergibt nie einen leeren Text; die Grenze steht in §6 (*Grenze der Tests*).
- *`$$` auf den Namen einer Verbindung angewandt* (Mutation aus L1) — äquivalent: Das
  Ergebnis enthält wieder ein `$` und wird abgelehnt; die Grenze steht in §6. Der verwandte
  Mutant beim Sammeln der Namen für `upstream` ändert das Verhalten und war zuerst grün;
  `a3bcad0` fügt den Fall hinzu, der ihn rot färbt (Tabelle DoD-Punkt 2).
- Eine Längenprüfung vor `strconv.Atoi` im Port war grün, weil `Atoi` eine zu große Zahl
  ohnehin als Grenze von `int` liefert; sie ist vor `f0180b6` entfernt, und der Mutant an
  ihrer Stelle (Ziffernprüfung entfernt) ist rot (Fall `h:+5`).

**Kommentare** (`AGENTS.md` §3.7, §3.11). Die neuen Kommentare beschreiben den Ist-Zustand; jede
Zusage darin hat eine Zeile oben. Abdeckungs-Deklarationen nennen nur Fälle, die der Test
enthält.

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `cfc4d6b` gesichtet, beim Schnitt vom
2026-10-09 am Stand `c297f9b` nachgesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×) — dieser Slice entsteht aus dem Schnitt von
  `slice-v1-abschluss-konfigurationsdatei`, der durch Übernahmen über eine Review-Sitzung
  wuchs. Er selbst übernimmt nur dessen DoD-Punkt 3 und bekommt drei Liefer-Punkte daraus; er
  ist kein weiterer Beleg. Ob der Geber bei seiner Closure der dritte ist, entscheidet dessen
  Closure.
- `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` (1×) — der Schnitt des Gebers und der
  Schnitt dieses Slice vom 2026-10-09 (§4) laufen formal über `next/` (Entscheidung des
  Nutzers vom 2026-10-09 zu V-117); kein weiterer Beleg, solange der `git mv` dieses Slice
  nach `next/` und zurück in der Historie steht.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (16×, verkörpert in `AGENTS.md` §3.12)
  und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (5×, verkörpert) — alle
  Randformen in §6 sind vor dem Code entschieden; eine weitere geht an den Architect zurück.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (17×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (23×, §3.11) — je Zusage eine Mutation,
  Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (19×, §3.9) — jede Korrektur zieht §1, §3 und §6 im
  selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — der Geber nennt diesen Slice in
  §1 unter *Abgegeben*; mit dem Schnitt vom 2026-10-09 nennt dieser Slice
  `slice-v1-abschluss-upstream-verbinden` in §1 unter *Abgegeben*, der ihn in §1 unter
  *Übernimmt* führt, und `slice-v1-abschluss-einspielen` führt die Wirkung bei `play` unter
  *Übernommen aus* `slice-v1-abschluss-upstream-verbinden`, alles im selben Commit.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (1×) — der Schnitt des Gebers gab diesem
  Slice eine Hälfte an der Grenze einer Review-Sitzung (§6 *Risiken*, 1050 bis 1300 Zeilen),
  und dieser Plan nannte den zweiten Schnitt als vorab benannte Rückführung. Anders als im
  Beleg trat sie vor dem ersten Code-Commit ein, nicht nach Liefer-Commits; ob das ein
  weiterer Beleg ist, entscheidet die Closure dieses Slice.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor dem
Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
