# Review-Report: slice-harness-d-check-v0-85 — 2026-10-09

**Review-Art:** Code-Review (Pin, Workflow-Commands, Reviewer-Skill, Plan-Berichtigung, neuer Nehmer-Slice). Geprüft wird gegen Plan §1, §3, §6 und §8, die Hard Rules (`AGENTS.md` §3.3, §3.6, §3.9 bis §3.13) und den Pin selbst (Modul 10). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff afc2528..d892be1`:

- `4c7ae6a`: `d-check.mk`, `DCHECK_IMAGE` und `DCHECK_DIGEST` auf `v0.85.0`.
- `79d7bb8`: `.harness/skills/reviewer.md`, MEDIUM *Nehmer nicht nachgezählt*.
- `825b71e`: Plan §7, Belege des Implementers (Pin, Gegenprobe, Freshness-Audit, Läufe).
- `dfb7e88`: Planner, Berichtigung von DoD-Punkt 1, Plan §1, §2, §3, §6, §8; neuer Slice `slice-harness-lh-links-pflicht` in `open/`; Drift-Log der Roadmap.
- `d892be1`: Block *Strenges Doc-Gate* in `.claude/commands/implement-slice.md`, `plan-welle.md`, `close-welle.md`; Plan §3 und §7.

**Skill:** `.harness/skills/reviewer.md` am Stand `d892be1`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- Plan von `slice-harness-d-check-v0-85` ganz am Stand `d892be1`, dazu der Plan-Diff je Commit und die Fassung am Stand `afc2528`. §7 *Belege des Implementers* gelesen, den Bericht des Implementers nicht.
- Plan von `slice-harness-lh-links-pflicht` ganz (Nehmer, `open/`).
- `.d-check.yml` ganz, `d-check.mk`.
- `AGENTS.md` §3.3, §3.6, §3.9 bis §3.13 (§3.13 mit *Nachzählen beim Eintragen*).
- Baseline `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice und §Ein Slice, dessen Gegenstand ein anderer übernimmt.
- Register: `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, `BEO-REPO/slice-waechst-durch-uebernahmen` (`observation.md`).
- Keine Lastenheft-Kennung im Scope; keine ADR im Diff oder in den Commit-Messages. Vorheriger Report zur Form: `slice-harness-lint` (Review). Die höchste vergebene Nummer vor diesem Lauf war F-554.

**Ausgeführte Läufe:**

- Pin selbst nachgeprüft: `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.85.0` meldet den OCI-Index `sha256:c07f1fe6053b1f790c4a1e01a76bcf4d3a8ff85e6eb609fe1aaaaf6ab6f09abe` (amd64 `eb73e50a…`, arm64 `69493255…`); `docker image inspect` auf den Digest nennt denselben RepoDigest. Gleich `DCHECK_DIGEST`.
- Freshness-Audit selbst nachgeprüft: `gh release list -R pt9912/ai-harness-course` nennt `v6.17.0` (Latest, 2026-10-07) vor `v6.16.0`. Asset `lab-regelwerk.zip` von `v6.17.0` nur in den Scratchpad geladen (`sha256 afe50df8…`), `diff -r` gegen `.harness/baseline/v6.16.0/`: über Quell-Zeile und Tag hinaus nur `grundlagen-harness-dateien.md`, `modul-13-quality-gates.md` (Disjunktheit der Teile eines Gate-Index), `regelwerk/README.md` (Stand-Zeile) und `templates/.d-check.yml` (`authority-disjoint`). Keine Datei neu oder entfallen. Deckt sich mit §7.
- Neun frische Kopien des Arbeitsbaums an `d892be1` unter dem Scratch-Präfix `rev-dcheck-`, je jeder Eintrag der obersten Ebene außer `.git` mit `cp -r`, Mutation angehängt, Lauf mit dem gepinnten Digest, Kopie danach gelöscht. Tabelle unten. Nichts lief im Repo.
- `make docs-check` vor dem Commit.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-555 | MEDIUM | Der Block *Strenges Doc-Gate* sagt in allen drei Commands zu, `docs-check` sei rot bei einer `ADR-`-Kennung ohne Link „in einer gescannten `.md`, auch in Inline-Code“. Eine Kennung `ADR-` mit vier Ziffern in einem umzäunten Code-Block ist kein Befund (M5: 0 Befunde, Exit 0); die Gegenprobe in §7 mutiert nur Fließtext und Inline-Code, und DoD-Punkt 1 schreibt „`ADR-` ohne Link überall“. *Failure-Szenario:* Ein Planner oder Implementer schreibt eine ADR-Kennung in ein Beispiel im Code-Block (Commit-Message, Hilfetext), liest den Block als Zusage, dass das Gate jede unverlinkte Kennung fängt, und das Gate bleibt grün. Die Klasse ist verkörpert und wiederkehrend; der Diff fasst den Block gerade, damit er nur zusagt, was geprüft ist. | `AGENTS.md` §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `.claude/commands/implement-slice.md` · „gescannten `.md`, auch in Inline-Code (`id-unlinked`)“ (Z. 39); `.claude/commands/plan-welle.md` · „gescannten `.md`, auch in Inline-Code, und bei einem toten“ (Z. 29); `.claude/commands/close-welle.md` · „`.md`, auch in Inline-Code, und bei einem toten Linkziel“ (Z. 27); Plan · „`ADR-` ohne Link überall, auch in Inline-Code“ | ja (M5 grün) | Kommentar-Zusage weiter als die Prüfung |
| F-556 | MEDIUM | Geber und Nehmer werden im selben Commit `dfb7e88` mit zwei verschiedenen Schicht-Teilungen gezählt. Der Geber führt die Spezifikation als eigene Schicht (§1 *Schicht-Abgrenzung* „Produkt-Code, Tests unter `test/` und die Spezifikation“) und Skill und Commands als „Agenten-Anweisungen“; der Nehmer fasst Spezifikation, Pläne und Commands zu „Dokumentation“ zusammen und kommt so auf zwei Schichten. Nach der Teilung des Gebers berührt der Nehmer Gate-Konfiguration (`.d-check.yml`), Spezifikation (§11, neue Stelle), Agenten-Anweisungen (drei Commands) und den Doku-Bestand, also mehr als zwei. *Failure-Szenario:* Der Nehmer geht mit einer Zählung nach `next/`, die nur unter einer Teilung hält, die der Geber selbst nicht verwendet; die Größe findet erst die Messung beim Start oder der Architect, das Muster aus `BEO-REPO/slice-waechst-durch-uebernahmen`. | `AGENTS.md` §3.13 *Nachzählen beim Eintragen*; Reviewer-Skill MEDIUM *Nehmer nicht nachgezählt*; Baseline `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice | `docs/plan/planning/open/slice-harness-lh-links-pflicht.md` · „(Gate-Konfiguration `.d-check.yml`; Dokumentation: Spezifikation, Pläne, Commands)“; `docs/plan/planning/in-progress/slice-harness-d-check-v0-85.md` · „(Gate-Konfiguration; Agenten-Anweisungen: Skill und Commands)“ | nein (Urteil über die Schicht-Teilung) | Schicht-Zählung des Nehmers nicht reproduzierbar |
| F-557 | LOW | `implement-slice.md` ordnet den Code `matrix-forbidden` auch einer blanken oder in Inline-Code geschriebenen ADR-Kennung aus `spec/` zu. Rot ist das, aber mit `id-unlinked` (M2); die Klasse `adr` der Matrix führt kein `token:`, `matrix-forbidden` meldet nur der Link. Für Slice (M1), Welle (M8) und `MR-` trifft die Zuordnung zu. Nicht MEDIUM: Die Zusage „rot“ hält, falsch ist nur der genannte Code. | `AGENTS.md` §3.11 | `.claude/commands/implement-slice.md` · „`MR-`-Kennung, auch blank oder in Inline-Code (`matrix-forbidden`)“ (Z. 41) | ja (M2) | Befund-Code im Hilfetext falsch zugeordnet |
| F-558 | LOW | Die Berichtigung begründet die Unerfüllbarkeit damit, dass „die erste Abgrenzung unten `.d-check.yml` unverändert lässt“. Am Stand `afc2528` hielt die erste Abgrenzung nur die Zeile `modules:` fest, nicht die `ids`-Muster; ausgeschlossen war `.d-check.yml` nur durch die Schicht-Abgrenzung („ändert einen Pin, einen Skill und diesen Plan“) und durch §3. Die Folgerung trägt, die genannte Stelle nicht. | `AGENTS.md` §3.9; Plan §1 | `docs/plan/planning/in-progress/slice-harness-d-check-v0-85.md` · „die erste Abgrenzung unten `.d-check.yml` unverändert lässt“ (Z. 53) | nein (Lesen) | Begründung zeigt auf die falsche Abgrenzung |
| F-559 | LOW | Zwei Plan-Stellen folgen dem Gelieferten nicht. §3 sagt, die MEDIUM-Klasse *Adresse nimmt nicht an* sei „um das Nachzählen … ergänzt“; geliefert ist eine eigene Klasse *Nehmer nicht nachgezählt*, die bestehende ist unverändert. §1 überschreibt „Herkunft der drei Teile“ und führt seit `dfb7e88` vier. | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | `docs/plan/planning/in-progress/slice-harness-d-check-v0-85.md` · „MEDIUM *Adresse nimmt nicht an* um das Nachzählen nach `AGENTS.md` §3.13 ergänzt“ (Z. 138); · „**Herkunft der drei Teile**“ (Z. 39) | nein (Lesen) | Plan folgt der Korrektur nicht |
| F-560 | LOW | Die neue Klasse stimmt mit `AGENTS.md` §3.13 *Nachzählen beim Eintragen* überein: Auslöser ist das Eintragen einer Sendung in §1 oder DoD des Nehmers, Prüfstelle §8 des Nehmers, Grenze drei Punkte und zwei Schichten nach Modul 5, „läge darüber“ ist erfasst. Eine Bedingung der Regel fehlt im Auslöser: Die Regel verlangt die Zählung „im selben Commit“, die Klasse prüft nur, dass §8 sie nennt. Ein Diff, der die Sendung in einem Commit einträgt und die Zählung in einem späteren nachträgt, ist nach der Klasse ohne Befund. | `AGENTS.md` §3.13 *Nachzählen beim Eintragen* („zählt im selben Commit“) | `.harness/skills/reviewer.md` · „und §8 des Nehmers nennt die Zählung seiner Liefer-Punkte“ | nein (Lesen) | Skill-Klasse enger als die verkörperte Regel |
| F-561 | INFO | Der Block in `implement-slice.md` nennt den Verweis auf eine superseded ADR nicht mehr. Das Gate meldet einen Link darauf in jeder Datei außerhalb von `docs/plan/adr/README.md` und `docs/reviews/**` als `matrix-inactive` (M6). Nach §3.11 ist das Weglassen zulässig; ein Agent, der auf diesen Befund läuft, findet im Block aber weder den Code noch einen Ausweg. Lebende Verweise auf die beiden superseded ADRs gibt es heute nicht. Für den Implementer. | Maintainability | `.claude/commands/implement-slice.md` · „Nicht geprüft: `LH-`-Kennungen“ | ja (M6) | Aktiver Befund des Gates im Hilfetext nicht genannt |
| F-562 | INFO | Zur Zuordnung bei der Closure (§6, erstes Risiko). Die Fehllesung „bares `LH-`-Token ist rot“ kam aus einem auskommentierten Muster in `.d-check.yml` und aus dem alten Block der Commands, nicht aus einer aktiven Regel, die das Werkzeug anders auswertet. `observation.md` von `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` beschreibt den zweiten Fall (Werkzeug wertet eine aktive Regel anders aus); der erste Fall liegt näher an `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. Mit einem dritten Beleg stünde der erste Eintrag an der Schwelle. Für die Closure. | Plan §6; Baseline `v6.16.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register | `docs/plan/planning/in-progress/slice-harness-d-check-v0-85.md` · „ob das ein Beleg zu diesem Eintrag ist, entscheidet die Closure“ | nein (Urteil) | — (Hinweis zur Zuordnung) |
| F-563 | INFO | Für den Verifier, zu den Schwerpunkten. **Pin:** Digest gehört zum Tag (eigener Abruf, oben); unveränderte Kopie 436 Dateien, 0 Befunde. **Gegenprobe:** Slice aus `spec/` auch in Inline-Code rot (M1), Welle aus `spec/` rot (M8, im Block nicht zugesagt), `ADR-` ohne Link in einem Command rot (M7), `LH-` blank grün (M3), `MR-` in Inline-Code außerhalb von `spec/` grün (M4) — die Nicht-Zusagen halten. **Audit:** Release-Liste und Delta wie in §7; Ausgang für `MR-000` und `MR-001` nachvollziehbar, Regelblock-Tabelle ohne Umbenennung. **Berichtigung:** vom Planner nach Rückgabe durch den Implementer, mit Entscheidung des Nutzers, Drift-Log-Zeile und Sendung; der Implementer hat DoD-Punkt 1 nicht selbst geändert. **Nehmer:** §1 und DoD führen die Sendung mit der Kennung des Gebers, kein Ausschluss trifft sie, er liegt in `open/` und startet erst nach dem Geber. **§3.6:** keine Lockerung, `.d-check.yml` unverändert; die DoD verlor eine Zusage über ein Verhalten, das das Gate nie hatte. | Plan §2, §4; `AGENTS.md` §3.6, §3.13 | `d-check.mk` · `DCHECK_DIGEST ?= sha256:c07f1fe6` | ja (Läufe oben) | — (Negativbefund) |

**Eigene Mutationen** (jede in einer frischen Kopie an `d892be1`, angehängt an das Ende der Datei, Image per Digest `c07f1fe6…`; die unveränderte Kopie: 436 Dateien, 0 Befunde, Exit 0):

| ID | Mutation | Ergebnis |
|---|---|---|
| M1 | `spec/architecture.md`: Slice-Kennung in Inline-Code | rot, `slice- matrix-forbidden`, Exit 1 |
| M2 | `spec/architecture.md`: `ADR-` mit vier Ziffern in Inline-Code | rot, nur `id-unlinked`, kein `matrix-forbidden`, Exit 1 (F-557) |
| M3 | `harness/README.md`: Lastenheft-Kennung blank | grün, 0 Befunde (Nicht-Zusage hält) |
| M4 | `harness/README.md`: `MR-001` in Inline-Code | grün, 0 Befunde (Nicht-Zusage hält) |
| M5 | `harness/README.md`: `ADR-` mit vier Ziffern in einem umzäunten Code-Block | grün, 0 Befunde (F-555) |
| M6 | `harness/README.md`: Link auf die superseded ADR 0021 | rot, `matrix-inactive`, Exit 1 (F-561) |
| M7 | `.claude/commands/plan-welle.md`: `ADR-` mit vier Ziffern blank | rot, `id-unlinked`, Exit 1 |
| M8 | `spec/architecture.md`: Welle-Kennung blank | rot, `welle- matrix-forbidden`, Exit 1 |

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `d-check.mk` | geprüft, ohne Befund. Nur `DCHECK_IMAGE` und `DCHECK_DIGEST` geändert, Digest selbst gegen die Registry nachgeprüft; der Kopfkommentar bleibt wahr. |
| `.d-check.yml` | geprüft, ohne Befund. Nicht im Diff, wie §1 sagt; Hard Rule 3.6 nicht berührt. |
| `.claude/commands/` | geprüft. Befunde F-555, F-557, F-561. Sonst sagt der Block nur zu, was eine Mutation rot zeigt; die Nicht-Zusagen halten (M3, M4). |
| `.harness/skills/reviewer.md` | geprüft. Befund F-560; sonst deckungsgleich mit `AGENTS.md` §3.13 *Nachzählen beim Eintragen*, Herkunfts-Anker richtig. |
| Plan `slice-harness-d-check-v0-85` | geprüft. Befunde F-558, F-559, F-562. Berichtigung nach dem Code regelgerecht: Rückgabe durch den Implementer in §7, Entscheidung des Nutzers datiert, §1, §2, §3, §6, §8 im selben Commit nachgezogen (`AGENTS.md` §3.9), Drift-Log-Zeile. |
| Plan `slice-harness-lh-links-pflicht` | geprüft. Befund F-556. Annahme der Sendung (§3.13) ohne Befund; Randformen vor dem Code in §6 genannt und dem Architect zugewiesen (§3.12). |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft, ohne Befund. Drift-Log-Zeile nennt Umplanung, Datum und Grund, keine Schließung. |
| Freshness-Audit (§7) | geprüft, ohne Befund. Release-Liste und Delta selbst nachgefahren, deckungsgleich; am vendored Baum nichts geändert (`make docs-check` und `baseline-verify` laufen mit `v6.16.0`). |
| Hard Rule 3.3 | geprüft, ohne Befund. Kein Move im Diff. |
| Hard Rule 3.10 | geprüft. Je Zusage des Blocks eine Zeile der Gegenprobe; Lücke F-555 (Code-Block). |
| Hard Rule 3.12 | geprüft, ohne Befund. Kein neuer Vertrag im Geber; keine Randform im Code-Commit. |
| Commit-Messages `4c7ae6a` bis `d892be1` | geprüft, ohne Befund. Jede nennt `slice-harness-d-check-v0-85` und `MR-000`, keine eine Struktur-Kennung. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 3 |

**Summary:** 0 HIGH · 2 MEDIUM · 4 LOW · 3 INFO (F-555 Block sagt „ohne Link in einer gescannten `.md`“ rot zu, Code-Block grün, M5; F-556 Geber und Nehmer zählen Schichten mit verschiedenen Teilungen; F-557 `matrix-forbidden` für ADR aus `spec/` falsch zugeordnet; F-558 Berichtigung zeigt auf die falsche Abgrenzung; F-559 §3 und §1 folgen dem Gelieferten nicht; F-560 Skill-Klasse ohne „im selben Commit“; F-561 `matrix-inactive` im Block nicht genannt; F-562 Zuordnung der LH-Fehllesung; F-563 Pin, Audit, Annahme und Berichtigung bestätigt). Wiederkehrende Klassen: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-555, F-557), `BEO-REPO/slice-waechst-durch-uebernahmen` (F-556), `BEO-REPO/plan-folgt-korrektur-nicht` (F-559).

**Finding-Klassen dieses Laufs:** Kommentar-Zusage weiter als die Prüfung · Schicht-Zählung des Nehmers nicht reproduzierbar · Befund-Code im Hilfetext falsch zugeordnet · Begründung zeigt auf die falsche Abgrenzung · Plan folgt der Korrektur nicht · Skill-Klasse enger als die verkörperte Regel · Aktiver Befund des Gates im Hilfetext nicht genannt

## Verdikt

**Merge-blockierend:** nein, aber F-555 und F-556 brauchen vor der Closure einen Ausgang. Der Pin stimmt, der Audit stimmt, die Berichtigung ist regelgerecht, und der Nehmer nimmt an. Der Block der Commands sagt nach der Korrektur noch eine Form zu, die die Gegenprobe nicht abdeckt und das Gate nicht fängt.

**Übergabe:** F-555 und F-557 an den Implementer (Block enger oder Fall in der Gegenprobe; Code-Zuordnung). F-556 an den Planner (Schicht-Teilung des Nehmers mit der des Gebers abgleichen, bei mehr als zwei Schichten nach §3.13 entscheiden). F-558 und F-559 an den Planner bzw. Implementer (Plan-Text). F-560 an den Planner (die Klasse ist Gegenstand von DoD-Punkt 3; ob „im selben Commit“ in den Auslöser gehört). F-561 an den Implementer. F-562 an die Closure. F-563 an den Verifier. Ein Widerspruch zu F-555 oder F-556 läuft über den Architect, nicht über Herabstufung. Dieser Report ersetzt keine Verifikation.
