# Slice slice-erster-release-homebrew-nachweis: Homebrew-Nachweis (Abnahmeszenario 11)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag.

**Bezug:** [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)

**Berührte Spec-Stellen:** `LH-FA-19.a` · `SPEC-042` · `SPEC-067` (Punkt 19 (7), Nachtrag)

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

**Ziel:** Nach dem ersten stabilen Release liegt die Formel im Tap `pt9912/homebrew-pgwire-recorder`, und die dokumentierten Befehle installieren das Werkzeug auf macOS und Linux.

**Übernommen aus `slice-erster-release-veroeffentlichung`** (dort §1, Abgrenzung): die Homebrew-Formel im Tap und ihr Nachweis; das sind Ziel und DoD-Punkt 1 oben.

**Übernommen aus `slice-v1-abschluss-homebrew`** (dort §1 und §6, `AGENTS.md` §3.13): (1) der Abschnitt zu Homebrew in §2 *Installation* des Benutzerhandbuchs und die Zeile in `README.md`, jeweils nur mit Befehlen, die gegen den echten Tap gelaufen sind; (2) die **Formel des ersten stabilen Tags** im echten Tap: Sie entsteht durch den Job `homebrew` des Laufs jenes Tags (`SPEC-067` Punkt 19); schlägt der Job fehl, ist der Wiederanlauf des Jobs oder der Nachtrag von Hand mit den Token in der Umgebung der Weg (Punkt 19 (7), Handlung des Betreibers). Für `v0.1.0` gibt es keine Formel: Seine Binaries tragen das alte Schema (Entscheidung des Nutzers vom 2026-10-20, [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md)); (3) die Belege, die die Gegenprobe jenes Slice nicht liefert: das Schreibrecht des Tap-Tokens, die Annahme der Formel durch Homebrew, das Binary auf macOS und Linux, der Wortlaut von `gh release download`.

**Bindung aus `slice-v1-abschluss-binary-namen`** (`AGENTS.md` §3.13): Der Slice setzt einen veröffentlichten stabilen Tag voraus, dessen Binaries `pgwire-recorder-<os>-<arch>` heißen. Einen solchen Tag bereitet `slice-v1-abschluss-release-v0-1-1` vor (Welle `welle-v1-nachtrag`, Reihenfolge 4: `VERSION`, Eintrag der Freigabe, Bezug im Handbuch); gesetzt wird er vom Nutzer nach neuem Ja. **Bindung aus `slice-v1-abschluss-release-v0-1-1`** (`AGENTS.md` §3.13): Der Eintrag der Freigabe nennt das Ergebnis des Jobs `homebrew` (Lauf-Kennung, grün oder rot); Wiederanlauf des Jobs, Nachtrag von Hand und alle Belege des Taps gehören diesem Slice. Die Abschnitte „Binary vom GitHub-Release“ und „Container-Image aus der Registry“ in §2 *Installation* des Handbuchs und die Zeile zur Version in `README.md` schreibt `slice-v1-abschluss-release-v0-1-1`; dieser Slice ergänzt nur den Abschnitt zu Homebrew und die Plattform-Zeile (dieselbe Datei, andere Abschnitte, der Nachweis folgt dem Release-Slice).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Erzeugen der Formel — `slice-v1-abschluss-homebrew`; hier wird das fertige Verfahren mit dem echten Release angewandt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): Der Tap trägt die Formel des ersten stabilen Tags mit dem Schema ohne Version (`v0.1.1`, `slice-v1-abschluss-release-v0-1-1`; Commit im Tap als Beleg; erzeugt vom Job `homebrew` des Tag-Laufs, bei Fehler per Wiederanlauf oder Nachtrag von Hand nach `SPEC-067` Punkt 19 (7)), und Abnahmeszenario 11 ist auf macOS und Linux nachgewiesen: Installation mit den dokumentierten Befehlen, `pgwire-recorder version` meldet die installierte Version, ein Replay-Lauf startet. Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung). Aus `slice-doku-ist-stand` (dort §6, *Installation*): Der Slice liefert die Befehle des Homebrew-Abschnitts in §2 *Installation* des Handbuchs, belegt mit der Installation aus dem echten Tap auf macOS und Linux, und die Plattform-Zeile in §1 *Voraussetzungen* (macOS).
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
| Tap `pt9912/homebrew-pgwire-recorder` | neu (extern) | Formel des ersten stabilen Tags mit dem neuen Schema, durch den Job `homebrew` (Wiederanlauf oder Nachtrag von Hand bei Fehler) |
| `docs/user/benutzerhandbuch.md` | update | Befehle gegen die echte Installation prüfen |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-homebrew` ist `done` (das Verfahren der Formel; Reihenfolge in §5 von [welle-v1-nachtrag](../welle-v1-nachtrag.md), Entscheidung des Nutzers vom 2026-10-11), `slice-v1-abschluss-binary-namen` ist `done`, und **ein stabiler Tag nach beiden ist veröffentlicht**: Seine Binaries heißen `pgwire-recorder-<os>-<arch>`, und der Lauf hatte den Job `homebrew` (Handlung des Nutzers nach einem neuen Ja zum Tag; `v0.1.0` zählt nicht, Entscheidung des Nutzers vom 2026-10-20). Der Slice, der diesen Tag vorbereitet und belegt, ist `slice-v1-abschluss-release-v0-1-1`; er liegt in `done/` (Start-Trigger), und der Tag ist veröffentlicht.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: der Nachweis braucht eine zweite Plattform-Umgebung, die nicht bereitsteht — zurück zur Zerlegung.
- `in-progress` → `open`: Der Tap lässt sich nicht beschreiben (Zugang) — extern klären.

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

- Das Verhalten von Homebrew gegenüber Drittanbieter-Taps (`brew trust`) ist erst mit dem echten Tap belegbar — **Ausgang:** offen bis Closure.
- Aus `slice-v1-abschluss-release-v0-1-1` (Bindung, §1): Ergebnis des Jobs `homebrew` im Lauf des Tags — grün: der Commit im Tap ist der Beleg des ersten Liefer-Punkts; rot: Wiederanlauf nur dieses Jobs oder Nachtrag von Hand, kein neuer Tag (`SPEC-067` Punkt 19 (7)); die Handbuch-Abschnitte zu Binary und Image schreibt der Release-Slice — **Ausgang:** offen bis Closure.
- Aus `slice-v1-abschluss-homebrew` (§6 dort), je mit Ausgang bei Closure dieses Slice: (a) das Schreibrecht des Tap-Tokens — das Geheimnis `HOMEBREW_TAP_GITHUB_TOKEN` (Fine-grained PAT nur für den Tap, Contents: Read and write, Entscheidung des Nutzers vom 2026-10-20) zeigt der Job `homebrew` des Laufs, der die erste Formel erzeugt; läuft die Formel stattdessen per Nachtrag mit dem eigenen Token des Nutzers, ist das Recht des Geheimnisses nicht belegt, und der Ausgang ist *weiter offen* (Register); (b) die Annahme der Formel durch Homebrew (Syntax, die Zeile `version`, Installation einer Datei ohne Archiv, `test do`) — eine Korrektur geht über die Vorlage und einen erneuten Nachtrag, ohne neuen Tag; (c) der Wortlaut von `gh release download` mit dem echten `gh` — **Ausgang:** offen bis Closure.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Die Dokumentation ist eine Schicht; der Handbuch- und README-Teil liegt im ersten Liefer-Punkt und in denselben zwei Dateien, kein neuer Liefer-Punkt. Schichten: zwei Schichten (Tap mit Formel, Dokumentation). Nachgezählt am 2026-10-20 mit der Sendung aus `slice-v1-abschluss-homebrew` (Architect, `AGENTS.md` §3.13): **Liefer-Punkte: 2** (1: Formel im echten Tap; 2: Abnahmeszenario 11 mit Handbuch- und README-Teil); Schichten bleiben zwei (der Tap als Ergebnis des Jobs oder einer Handlung des Betreibers, Dokumentation). Nachgezählt am 2026-10-20 mit der Bindung aus `slice-v1-abschluss-binary-namen`: Der Nachtrag für `v0.1.0` entfällt, die Formel entsteht durch den Job; **Liefer-Punkte bleiben 2, Schichten zwei**; der Slice liegt unter der Grenze.

Nachgezählt am 2026-10-20 mit der Bindung aus `slice-v1-abschluss-release-v0-1-1` (Planner, `AGENTS.md` §3.13): Der Start-Trigger nennt den Slice, der Tag kommt von ihm; die Abgrenzung der Handbuch-Abschnitte (Binary und Image beim Release-Slice, Homebrew hier) verkleinert den Teil dieses Slice; **Liefer-Punkte bleiben 2, Schichten zwei**; der Slice liegt unter der Grenze.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
