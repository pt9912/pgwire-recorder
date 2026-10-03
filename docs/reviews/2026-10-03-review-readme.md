# Review-Report: README.md (ungestagte Änderung) — 2026-10-03

**Review-Art:** Code — Dokument-Review gegen Vorlage, Baseline-Regel und Repo-Wahrheit (Maintainability); keine DoD-Prüfung.

**Gegenstand:** `git diff HEAD -- README.md` (ausgefüllt aus `project-readme.template.md`)

**Skill:** `.harness/skills/reviewer.md` @ a19ce26 · <!-- d-check:ignore (Skill-Pfad) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `v6.13.0` · `templates/project-readme.template.md`
- `v6.13.0` · `regelwerk/modul-02-harness-bootstrap.md` §Ziel-Form: Projekt-README
- ADR: [ADR-0007](../plan/adr/0007-strict-replay.md) (Strict Replay, für die Aussage im Kerngedanken)
- Lastenheft §1 Zweck und Geltungsbereich, LH-Abschnitt zur Abweichung im Replay
- `AGENTS.md` (Hard Rules, §2 Source Precedence)
- Vergleich: `/Development/pg-change-feed/README.md`, `/Development/d-migrate/README.md`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die erste Zeile unter der Überschrift ist die Rolle-Zeile („Rang 6 der Source Precedence …, Baseline-Regelwerk `modul-02-…`"). Sie ist Vorlagen-/Harness-Text für Agenten und das Erste, was ein Mensch liest; der Projektzweck folgt erst nach ihr. Die Vorbild-Repos beginnen mit Tagline bzw. Zweck-Satz. | Maintainability; `v6.13.0` · `regelwerk/modul-02-harness-bootstrap.md` §Ziel-Form: Projekt-README | README.md, Kopf (unter `# pgwire-recorder`) | nein — kein Gate prüft Leser-Eignung | Agenten-Text als Einstieg für Menschen |
| F-2 | MEDIUM | Die Aussage „Die Alternative ist, die Datenbank im Test auf Treiberebene nachzubauen; dabei prüft der Test nicht mehr, was die Anwendung tatsächlich auf dem Netz sendet und empfängt" steht weder im Lastenheft noch in Spezifikation oder ADRs. Lastenheft §1 nennt als Ausgangssituation nur den Aufwand einer realen Instanz; die Spezifikation erwähnt „Mock-Regeln" nur als Nicht-Ziel. Die Alternative samt Bewertung ist erfunden. | Maintainability (Wahrheitsgehalt, Rang-6-Regel „verweist, erfindet nicht"); Lastenheft §1 | README.md, Abschnitt „Warum pgwire-recorder?" | nein — Beleg nur per Lesen; `make docs-check` prüft keine Aussagen | Aussage ohne kanonischen Beleg |
| F-3 | LOW | „Aufzeichnen statt Nachbauen" im Kerngedanken übernimmt die in F-2 unbelegte Gegenfigur „Nachbauen" als Leitidee. Der belegte Teil (nur Aufgezeichnetes, strikte Reihenfolge, nie geratene Antwort) deckt sich mit Lastenheft (Abweichung im Replay) und [ADR-0007](../plan/adr/0007-strict-replay.md). | Maintainability; [ADR-0007](../plan/adr/0007-strict-replay.md) | README.md, Abschnitt „Kerngedanke" | nein | Aussage ohne kanonischen Beleg |
| F-4 | LOW | „Was ist …?" und „Warum …?" geben Lastenheft §1 Ausgangssituation und Produktziel fast wörtlich wieder („aus Sicht der Anwendung beobachtbare Verhalten der PostgreSQL-Kommunikation", „aufwendiger … langsamer, komplexer und weniger deterministisch"). Ein Verweis auf Lastenheft §1 fehlt an dieser Stelle; er steht nur weiter unten unter „Verträge". | Maintainability (Rang 6 dupliziert nicht) | README.md, Abschnitte „Was ist pgwire-recorder?" und „Warum pgwire-recorder?" | nein | Duplikation statt Verweis |
| F-5 | LOW | „Was kann ich heute tun?" nennt feste Zahlen („zwölf angenommene Entscheidungen", „vier Wellen, 14 Slices"). Sie stimmen heute (12 ADR-Dateien mit Status Accepted, 4 Welle-Dateien, 14 Slice-Dateien in `open/`), driften aber beim nächsten ADR oder Slice, ohne dass ein Gate es meldet. Das Target `make help` als einziger Einstieg wird nicht genannt. | Maintainability; `AGENTS.md` §3.7 (Zustandsfelder: Zustand mit auflösbarem Anker) | README.md, Abschnitt „Was kann ich heute tun?" | nein | Zahl im Prosa-Stand driftet |
| F-6 | LOW | Der Abschnitt „Was macht es vertrauenswürdig?" und der Rolle-Satz tragen interne Kennungen (Hard Rules, Source Precedence, `LH-*`-IDs, Gates) ohne Erklärung für Erstleser. Das ist Vorlagenbestand, aber für Projektleser undurchsichtig. | Maintainability | README.md, Abschnitt „Was macht es vertrauenswürdig?" | nein | Interne Kennungen ohne Erklärung |
| F-7 | INFO | Fehlt ohne Vorlagenverstoß: Lizenz (kein `LICENSE` im Repo), Quellort/Issue-Tracker (kein Git-Remote konfiguriert; `pg-change-feed` führt „Source & issues"), Zielgruppen-Abschnitt (`d-migrate`: „Who is it for?"; hier nur ein Halbsatz), englische Fassung/Sprachumschalter (beide Vorbilder zweisprachig). Kein Lizenz-/Quell-Eintrag erfinden, solange die Entscheidung des Auftraggebers fehlt. | Maintainability | README.md, gesamt | nein | README-Bestandteil fehlt |
| F-8 | INFO | Positiv geprüft: `make gates` und `make docs-check` laufen grün (58 Dateien, 0 Befunde), alle Links lösen auf, `make gates` führt tatsächlich nur `docs-check` und `baseline-verify` als Checks aus, das Benutzerhandbuch existiert und ist als „noch nicht veröffentlicht" kompatibel zur Aussage „Spezifikationsphase". Die Annahme (Lastenheft §1), Zielgruppe und Produktziel sind belegt. | Maintainability | README.md, gesamt | ja — `make docs-check` | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| README.md, Vorlagenkonformität (alle sechs Abschnitte, keine `<…>`-Platzhalter, Template-Hinweis-Block entfernt) | geprüft, ohne Befund |
| README.md, Links (`AGENTS.md`, `harness/README.md`, `spec/*`, `docs/plan/adr/`, `docs/plan/planning/`, `docs/reviews/`, Benutzerhandbuch) | geprüft, ohne Befund |
| README.md, genannte Make-Targets (`make gates`, `make docs-check`, `make baseline-verify`) | geprüft, ohne Befund |
| README.md, Source-Precedence-Rolle: kein ADR-/Slice-/Wellen-Inhalt wiederholt | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Agenten-Text als Einstieg für Menschen · Aussage ohne kanonischen Beleg · Duplikation statt Verweis · Zahl im Prosa-Stand driftet · Interne Kennungen ohne Erklärung · README-Bestandteil fehlt

## Empfehlung zum Einstieg (Hinweis, keine Anordnung)

Vorbild-Repos: `pg-change-feed` beginnt mit Titel, fetter Tagline im Blockquote, Quellort-Zeile und dann „Was ist …?"; `d-migrate` mit Titel, fetter Tagline, Badges und einem Überblicksabsatz. Beide stellen den Zweck vor jede Prozess-Aussage.

- Erster Satz unter der Überschrift als Tagline, abgeleitet aus Benutzerhandbuch §1 und Lastenheft §1: „> **Zeichnet die Kommunikation zwischen Ihrer Anwendung und PostgreSQL auf und spielt sie später ohne Datenbank wieder ab.**"
- Die Rolle-Zeile ist Vorlagen-Norm (kein Ausfüll-Hinweis) und bleibt stehen, gehört aber hinter den Projektzweck, etwa als Schlusszeile unter „Was macht es vertrauenswürdig?" oder als kleiner Hinweis direkt vor diesem Abschnitt.

**Verdikt-Hinweis.** Die Rolle-Zeile komplett zu löschen wäre ein Verstoß gegen die Vorlage; Umstellen ist keiner.

## Verdikt

**Merge-blockierend:** ja — wegen F-2 (unbelegte Aussage über eine Alternative im README). F-1 blockiert typischerweise ebenfalls; Entscheidung über die Platzierung der Rolle-Zeile liegt beim Implementer bzw. Architect. LOW und INFO blockieren nicht.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
