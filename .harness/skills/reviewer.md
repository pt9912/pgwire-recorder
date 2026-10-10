# Reviewer-Skill — pgwire-recorder

* Status: Accepted
* Bezug: `AGENTS.md` §3 (Harte Regeln) und §6 (Minimal Agent Workflow)
* Gilt für: den Review-Lauf nach Schritt 8 des Minimal Agent Workflow; ein eigenes Make-Target gibt es nicht

## Kontext-Eingang (Pflicht)

Was der Reviewer *immer* mitbringt, bevor er den Diff liest:

- Diff des PR
- `spec/lastenheft.md` (für referenzierte `LH-*`-IDs)
- ADRs, deren ID im PR oder in der Commit-Message vorkommt
- `AGENTS.md` §"Hard Rules"
- vorherige Findings am gleichen Modul (letzte ~5 PRs)
- §7 des Slice-Plans, *Belege des Implementers* (Mutationstabelle, Befundzeilen,
  Läufe); der Bericht des Implementers liegt dem Review nicht vor
  (seit slice-harness-blackbox-kern)

Ohne diesen Block sieht der Reviewer den Code, aber nicht *die Verträge, gegen
die er prüft*.

## Klassifikation

Jeder Anker HIGH/MEDIUM/LOW hat eine *konkrete* Liste — nicht generisch. INFO ist
bewusst kurz (Ergänzungs-Kanal, nicht Hauptkanal).

**HIGH** — eines der folgenden:
- ADR-Verstoß (Layer, Tool, Hard Rule)
- Sicherheits-Anti-Pattern (Injection, fehlende Auth-Prüfung)
- Korrektheitsfehler im *kritischen* Pfad (Record-/Replay-Pfad: unpassende oder geratene Antwort, verletzte Reihenfolge, Byte-Verlust zwischen Aufzeichnung und Wiedergabe)
- Suppression eines Gates (`#noqa`, `//nolint`, `[SuppressMessage]`) ohne ADR
- **Norm nur im Template-Kommentar** — eine Regel steht im `<!-- -->`-Block
  eines `.template.md` und nirgends sonst. Sie ist beim Adopter weg, sobald er
  die Kommentare entfernt. Kein Gate fängt das (siehe Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Template-Schichtung)
- **Kommentar trägt keine der Kommentar-Klassen** — ein Kommentar in Code, Config
  oder Skript beschreibt die verworfene Alternative („Ohne X wäre …"), einen
  abwesenden Text („früher stand hier …") oder bricht mitten im Satz ab, weil
  eine Teilersetzung den Rest stehen ließ. Kein Gate fängt das (siehe
  Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein Kommentar trägt)
- **Zustandsfeld trägt Chronik** — eine `Stand`-/`Status`-Zelle (Roadmap,
  Beobachtungs-Register, Meilenstein) erzählt, wie der Zustand entstand, statt
  Zustand und Beleg als Anker zu nennen; oder ein Drift-Log protokolliert
  Schließungen und erreichte Meilensteine. Kein Gate fängt das (siehe
  Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein Kommentar trägt,
  *Dieselbe Regel für Zustandsfelder*)
- **Spec-Stratum nennt Artefakt unterhalb oder außerhalb** — Lastenheft,
  Spezifikation oder Architektur nennen eine ADR, einen Slice, eine Welle, einen
  Commit oder (Lastenheft) die darunterliegenden Spec-Dateien; die Sicht trägt
  zudem keine Sprach-Artefakte (Code, Dateiendungen, Bibliotheksnamen,
  sprachspezifische Konstrukte, `AGENTS.md` §3.4)
- **ADR oder Move verletzt Immutabilität/Commit-Regel** — eine `Accepted`-ADR
  wird inhaltlich überschrieben (`AGENTS.md` §3.5), oder ein Move und eine
  Inhaltsänderung stehen im selben Commit (`AGENTS.md` §3.3)
- **Core-Reinheit verletzt** — `internal/hexagon/**` importiert Adapter,
  `pgproto3`, eine YAML-Bibliothek oder das Dateisystem, oder `pgproto3` liegt
  außerhalb der beiden PGWire-Adapter

**MEDIUM** — eines der folgenden:
- unklare Fehlerbehandlung am Rand des Spec-Bereichs
- fehlende Negativtests bei neuem öffentlichem Vertrag
- Wiederholung eines Musters, das schon zweimal LOW war
- **Randform im Code-Commit** — eine Randform erscheint im selben Commit wie eine
  Änderung an Code, Tests oder Gates neu in §6 des Slice-Plans, oder §6 weist sie dem
  Review statt dem Architect zu; mindestens MEDIUM, gegen `AGENTS.md` §3.12
  (seit slice-harness-blackbox-kern)
- **Adresse nimmt nicht an** — der Diff nennt einen anderen Slice neu oder geändert als
  Adresse (Folge-Slice, Risiko-Ausgang *eingetreten*, Register *geplant*, Abgrenzung der
  Klasse 1), und §1 oder DoD des Nehmers führt die Sendung nicht mit der Kennung des
  Gebers, ein Punkt unter *Ausdrücklich NICHT* trifft sie, oder der Nehmer liegt in
  `done/`. Jede solche Adresse wird gegen §1 und DoD des Nehmers gelesen; mindestens
  MEDIUM, gegen `AGENTS.md` §3.13 (seit slice-lint-bestand-kern-driven)
- **Nehmer nicht nachgezählt** — der Diff trägt eine Sendung in einen Nehmer ein (§1 oder
  DoD des Nehmers), und §8 des Nehmers nennt die Zählung seiner Liefer-Punkte und
  Schichten mit der Sendung nicht im selben Commit, oder der Nehmer liegt mit ihr über drei Liefer-Punkten
  oder zwei Schichten. Gezählt wird nach Baseline-Regelwerk `modul-05-planning-harness.md`
  §Ziel-Form: Slice; mindestens MEDIUM, gegen `AGENTS.md` §3.13 *Nachzählen beim
  Eintragen* (seit slice-v1-abschluss-einspielen)

**LOW** — *mit Konventions-Anker* (ADR, Hard Rule, Linter-Regel, Eintrag im
Reviewer-Skill): stilistisch unschön ohne semantische Auswirkung, einmalige
Tippfehler, unbenutzte Imports.

**INFO** — Hinweis ohne erwartete Aktion (z. B. „diese Stelle hat ein passendes
ArchUnit-Pendant, das du nicht kennst").

## Was dieser Skill NICHT macht

- Keine Lösungsvorschläge („schreib das so") — Reviewer kategorisiert,
  Implementer entscheidet.
- Kein Refactoring-Vorschlag, der über den Diff hinausgeht.
- Keine Verifikation gegen DoD — das ist Verifier-Aufgabe (Modul 11).
- Keine Validation gegen reale Bedürfnisse — das ist Validator-Aufgabe.
- **Kein Stil-Polizist:** Formatierung oder Benennung ohne Konventions-Anker ist
  kein Finding.
- **Kein HIGH- oder MEDIUM-Finding ohne Failure-Szenario:** was sich nicht als
  konkretes Versagen erzählen lässt, wird nicht als HIGH oder MEDIUM gemeldet.

Wenn etwas auffällt, das in diese Kategorien gehört: ein INFO-Finding mit Verweis
auf die zuständige Rolle.

## Output-Schema

Jedes Finding:

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: ADR-ID, `LH-*`-ID, Hard-Rule-Name oder „Maintainability"
- `pfad`: Datei · wörtliches, in der Datei eindeutig auffindbares Kurzzitat
  der Stelle als Anker; die Zeile darf als Lesehilfe dazu, ist aber nicht der Anker
- `befund`: 1–2 Sätze, beobachtbar, ohne Lösungsvorschlag
- `verifizierbar`: ja/nein — gibt es einen Gate-Lauf, der es bestätigen würde?
- `klasse`: stabile Kurz-Bezeichnung des Fehlermusters, z. B. „Tie-Break in
  sortierender Operation nicht dokumentiert" — speist den
  Steering-Loop-Zähler (siehe §Pflege)

Zusätzlich am Ende: eine Zeile „geprüft, ohne Befund" pro betrachtetem
Verzeichnis (Negativbefund-Zeile — sonst ist „keine Findings" nicht von „nicht
geprüft" unterscheidbar). Report-Gerüst für den ganzen Lauf:
`.harness/baseline/<tag>/templates/docs/reviews/review-report.template.md` (vendored), ein Report pro Lauf, Folgeläufe als
neue Datei statt Überschreibung. Der Dateiname trägt Datum, Rolle und die **volle**
Slice-Kennung (`<YYYY-MM-DD>-review-<slice-Kennung>.md`); ein Folgelauf
hängt `-r2`, `-r3` an die Kennung (`<YYYY-MM-DD>-review-<slice-Kennung>-r2.md`).

## Pflege (Steering-Loop)

Bei dreimaligem Auftreten desselben Findings:

- ist die Kategorie noch richtig? → Klassifikation schärfen
- gibt es einen ADR/`AGENTS.md`-Eintrag, der das verhindert hätte?
  → Folge-ADR oder `AGENTS.md`-Update
- gibt es eine Fitness Function, die das prüfen würde? → Modul 13, Gate hinzufügen

Diese Skill-Datei wird **nicht** überschrieben, sondern versioniert
(ADR-Hard-Rule, Modul 4).
