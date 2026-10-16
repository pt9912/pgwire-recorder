# Review-Report: slice-v1-abschluss-verbindungsende-doku — 2026-10-17

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten); der Gegenstand ist Doku (Handbuch, README), geprüft gegen das gebaute Binary.

**Gegenstand:** `slice-v1-abschluss-verbindungsende-doku` (Schritt 42, Handbuch und README zum Verbindungsende), Diff `e9f61167..f61e3458` (Pfade `docs/user`, `README.md`, dazu Plan und Roadmap).

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec08 · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-v1-abschluss-verbindungsende-doku` (§1, §3, §4, §6, §7 Belege des Implementers, §8)
- [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md); Geber `slice-v1-abschluss-verbindungsende-format` bis `slice-v1-abschluss-verbindungsende-upstream` (geschlossen)
- `LH-FA-25`, `LH-FA-24`, `LH-FA-20`, `LH-FA-22` (Boundary zur Version der Aufzeichnung)
- `AGENTS.md` (Hard Rules, §3.7, §3.9, §3.11, §3.13, §3.14)
- Eigene Probe am gebauten Binary (`make build`, Image `pgwire-recorder:dev`, `postgres:17-alpine` aus `harness/mk/integration.mk`, eigenes Docker-Netz, ein Lauf zur Zeit, Container und Netz danach entfernt): `record` mit `pg_terminate_backend` (YAML und SQLite), `replay` derselben Aufzeichnung, `play` mit und ohne `--allow-recorded-errors`, mit `--compare-responses`, mit `--continue-on-error` und in Kombination gegen die Fixtures `play-ende-*.yaml`, Leerlauf mit `idle_session_timeout`, Abbruch der Datenbank (`docker kill`), Version 1 mit `connection_closed` und Version 3. Zeilen 1, 4 bis 8, 9, 12 bis 20 der Zusagen-Tabelle in §7 (13 von 22; Zeilen 2, 10, 11 nur über Tests der Services) wurden nachgefahren; jede Ausgabe stimmt mit dem Handbuch überein (Wortlaut der Log-Zeilen, Exit-Codes 0/2/3/4/5, `eingespielt`/`verglichen`/`abweichend`, `version|2`, `1|1|connection_closed`, `server|{}`).

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill
— Tag einsetzen, denn diese Zeile wandert in den eingefrorenen Report.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | DoD-Punkt 1 verlangt im Handbuch den Satz „ältere Programmversionen lehnen sie ab“ (Version 2); das Handbuch enthält ihn nicht (Implementer: nicht belegbar, §7 *Nicht aufgenommen* (a)), die DoD ist unverändert und widerspricht damit dem Diff und dem Satz in §1 „was sie nicht zeigt, bleibt draußen“. Die Aussage stützt sich auf eine Zusage des Lastenhefts zum Produkt (`LH-FA-22` Boundary), nicht auf eine Probe am Binary; im Handbuch deckt sie nur die Zeile `PGR-E3002` („stammt aus einer anderen Programmversion“) ab. Der Verifier misst gegen einen Wortlaut, den der Slice bewusst nicht liefert. | `AGENTS.md` §3.9; `AGENTS.md` §3.11; `LH-FA-22` | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-doku.md` · „ältere Programmversionen lehnen sie ab“ (§2, Punkt 1) gegen §7 „Nicht aufgenommen, weil nicht belegbar“ | nein — kein Gate liest DoD gegen Handbuch; `grep` im Handbuch nach „ältere“ zeigt keinen Treffer | DoD folgt dem Diff nicht (Plan-Mismatch, `BEO-REPO/plan-folgt-korrektur-nicht`) |
| F-2 | INFO | Die Beispiele des neuen Abschnitts (Kommandozeilen, Ausgaben, YAML-Auszug, `sqlite3`-Abfragen) liegen nicht als Datei im Repo und laufen in keinem Test oder Gate gegen das Binary; belegt sind sie durch den Handlauf des Implementers (§7), den dieser Review unabhängig mit gleichem Ergebnis wiederholt hat. Der Plan entscheidet das so (§3: unter `test/integration` entsteht nichts); eine spätere Änderung am Binary färbt kein Gate rot, wenn der Abschnitt dann nicht mehr stimmt. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „Das Ende der Verbindung durch die Datenbank“, Unterabschnitt „Vorgehen“ | nein — ein Gate dafür gibt es nicht (§1 schließt einen Sensor aus) | Beispiel nur von Hand belegt |
| F-3 | INFO | Die Zeile `PGR-E5004` nennt für ein Ende zwei Abweichungen (anderer SQLSTATE; Aufzeichnung zeigt eines, die Datenbank nicht); die dritte Form (die Datenbank endet, die Aufzeichnung an der Stelle nicht, mit Art `Fehler` oder `Nachrichtenart`) steht nur im Abschnitt „Antworten beim Einspielen vergleichen“. Beobachtet: `play-ende-grenze-ohne-ende.yaml` mit Vergleich gibt Art `Nachrichtenart`, eine Aufzeichnung mit `ERROR` an der Stelle gibt Art `Fehler` mit beiden SQLSTATE. Die Zeile sagt „kann“, widerspricht dem Abschnitt nicht. | Maintainability | `docs/user/benutzerhandbuch.md` · „Auch ein Ende der Verbindung durch die Datenbank kann abweichen“ | nein | Tabellenzeile unvollständig gegen den Abschnitt |
| F-4 | INFO | Der Abschnitt nennt „Ersatz-Server“ für `PANIC` und für `ERROR` mit Abbruch (`PGR-E6001`), ohne ihn zu erklären; ein Leser kann beides nicht nachfahren. Belegt ist beides durch die Tests der Services (`internal/hexagon/services/record_ende_test.go`), nicht durch einen Lauf gegen das Binary; das Handbuch sagt es an beiden Stellen. Zusage und Prüfung decken sich, die Lesbarkeit leidet. | Maintainability | `docs/user/benutzerhandbuch.md` · „Ersatz-Server geprüft, der `PANIC` sendet, nicht mit PostgreSQL“ | nein | Prüfhilfe im Handbuch unerklärt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/user/benutzerhandbuch.md`, neuer Abschnitt „Das Ende der Verbindung durch die Datenbank“ (Vorgehen, Ergebnis, Hinweise zum Aufzeichnen, Wiedergeben, Einspielen, Fehlerbilder) | geprüft gegen das gebaute Binary, ohne weiteren Befund; Quantoren („jede“, „nie“, „alle“, „nur“, „auch“, „gleich welcher SQLSTATE“, „für `FATAL` wie `PANIC` und für beide Formate“, „mehrere Sitzungen mit einem Ende“, „andere Verbindungen laufen weiter“) tragen je eine Zeile in §7 oder einen Test der Services |
| `docs/user/benutzerhandbuch.md`, angepasste Sätze (Einleitung, Wiedergeben, SQLite-Tabelle `meta`, Einspiel-Hinweise, Vergleichs-Hinweise, Glossar) | geprüft, ohne Befund; keine Widersprüche zu Fehlerzeilen und Rest des Handbuchs; `version: 1` im Beispiel zu `--record-empty-sessions` bleibt wahr (Aufzeichnung ohne Ende) |
| Fehlercode-Tabelle `PGR-E3002`, `PGR-E3003`, `PGR-E4003`, `PGR-E5004`, `PGR-E6001` gegen `internal/hexagon/model/fehler.go` und Läufe | geprüft, ohne Befund; Codes vorhanden, Exit-Codes 3, 3, 4, 5 am Binary beobachtet, 6 über Tests der Services |
| §1 Ausschluss und Bindung (`-upstream-extended`): Extended-Aufzeichnung mit Ende gegen PostgreSQL, Herunterfahren mit Ende, „nur noch Schließen“/`Send` nach dem Ende | geprüft per Suche im Diff (erweitert, Extended, Sync, Flush, Herunterfahren, nur noch); keine Zusage im Handbuch oder README |
| AGENTS.md §3.11 / §3.7: Chronik, „noch nicht/kommt mit“, Verweise auf Spezifikation, Lastenheft, ADR, Slice, Welle, Review im Handbuch | geprüft per Suche im Diff, ohne Befund; README verweist nicht neu |
| `README.md` | geprüft, ohne Befund; Absatz stimmt mit dem Handbuch und mit den Läufen überein |
| Messung der Marke (§3.14) | nachgerechnet: `git diff --numstat e9f61167 HEAD -- docs/user README.md test/integration` ergibt 9+4 (README) und 244+30 (Handbuch), Summe 287 bei Marke 600 und Meldeschwelle 480; nicht erreicht |
| Gates | `make docs-check` (932 Dateien, 0 Befunde), `make kopf-check` (Exit 0), `make gates-schnell` (Exit 0) am Stand `f61e3458`; Arbeitsbaum danach sauber |
| Plan §3.13 (Adressen, Nachzählung) | geprüft, ohne Befund; der Diff trägt keine neue Sendung ein, §8 zählt 2 Liefer-Punkte, eine Schicht |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** DoD folgt dem Diff nicht (Plan-Mismatch, `BEO-REPO/plan-folgt-korrektur-nicht`) · Beispiel nur von Hand belegt · Tabellenzeile unvollständig gegen den Abschnitt · Prüfhilfe im Handbuch unerklärt

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM). Das Handbuch selbst ist gegen das Binary und die Zusagen-Tabelle stimmig; der Mismatch liegt im Plan: Der Planner entscheidet, ob die DoD den Satz „ältere Programmversionen lehnen Version 2 ab“ streicht (das Handbuch sagt nur zu, was ein Test oder eine Probe trägt) oder ob er mit eigener Prüfung ins Handbuch kommt. Zuständigkeit F-1: Planner (DoD, §1, Entscheidung des Nutzers); F-2 bis F-4: zur Kenntnis, Implementer nach Ermessen, keine Pflicht.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt: F-1 an den Planner); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
