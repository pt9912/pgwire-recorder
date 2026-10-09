# Verifikation: slice-harness-d-check-v0-85 — 2026-10-09

**Rolle:** Verifier (Modul 11). Geprüft wird, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-09-review-slice-harness-d-check-v0-85.md`, F-555 bis F-563); F-563 ging an mich, F-562 ordne ich nur ein.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-d-check-v0-85.md`, Kopf und §1 bis §8, gegen den Gesamt-Diff `afc2528..cc2b61c`. Implementer: `4c7ae6a`, `79d7bb8`, `825b71e`, `d892be1`, Nacharbeit `cc2b61c`. Planner: `dfb7e88`, Nacharbeit `6db272b`. Review: `13b17dc`. Dazu die beiden Nehmer `slice-harness-lh-links-pflicht` und `slice-harness-lh-links-bestand` (beide `open/`).

**Skill:** `.harness/skills/reviewer.md` am Stand `cc2b61c` (Gegenstand von DoD-Punkt 3, nicht Urteilsgrundlage dieses Laufs)
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

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

- der Plan ganz, besonders §1 (*Herkunft*, *Berichtigung*, *Ausdrücklich NICHT*), §2 (DoD), §3, §6 und §7 (*Belege des Implementers*, beide Gegenproben-Tabellen, *Nacharbeit*, Freshness-Audit, *Läufe*)
- die Nehmer `slice-harness-lh-links-pflicht` und `slice-harness-lh-links-bestand` ganz (§1, DoD, §4, §6, §8 mit *Schichtteilung*)
- `.d-check.yml` ganz, `d-check.mk`, der Diff der drei Commands `.claude/commands/implement-slice.md`, `plan-welle.md`, `close-welle.md` und von `.harness/skills/reviewer.md`
- `AGENTS.md` §3.9 bis §3.13 (§3.13 mit *Nachzählen beim Eintragen*)
- Baseline `v6.16.0` · `regelwerk/modul-02-harness-bootstrap.md` §Freshness-Audit der vendored Baseline; `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen

Keine Lastenheft-Kennung und keine ADR im Scope (Kopf: *Berührte Spec-Stellen* —). Beim Start war der Arbeitsbaum sauber, HEAD `cc2b61c`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.** Je Mutation eine frische Kopie des Arbeitsbaums im Scratchpad unter dem Präfix `ver-dcheck-` (jeder Eintrag der obersten Ebene außer `.git` mit `cp -r`), Mutation angehängt bzw. per `sed` an einer Zeile eingefügt, Lauf `docker run --rm --network none` mit dem Digest aus `d-check.mk` (`c07f1fe6…`), danach nur diese Kopie gelöscht. Unveränderte Kopie: `438 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Nichts lief im Repo außer `make gates` und `make docs-check`.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: d-check `v0.85.0`, Gegenprobe, Block *Strenges Doc-Gate*. **Teilweise bestätigt; abgelehnt in der Teil-Zusage über den Block (V-134).**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `DCHECK_IMAGE` = `ghcr.io/pt9912/d-check:v0.85.0`, `DCHECK_DIGEST` = `sha256:c07f1fe6053b1f790c4a1e01a76bcf4d3a8ff85e6eb609fe1aaaaf6ab6f09abe` | `d-check.mk` Z. 7–8 am Stand `cc2b61c`; Diff ändert nur diese zwei Zeilen | bestätigt |
| der Digest gehört zum Tag (Multi-Arch-Index) | eigener Abruf `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.85.0`: MediaType `application/vnd.oci.image.index.v1+json`, Digest `sha256:c07f1fe6…f09abe`, Manifeste `linux/amd64` `sha256:eb73e50a…` und `linux/arm64` `sha256:69493255…` — gleich §7 | bestätigt |
| `make docs-check` 0 Befunde | eigener Lauf in `make gates` am Stand `cc2b61c`: `d-check: 438 Datei(en) geprüft, 0 Befund(e)` | bestätigt |
| Gegenprobe: toter Anker, `ADR-` ohne Link, totes Linkziel, Link aus Spec-Stratum auf ADR je rot | eigene Kopien: `anchor-missing`, `id-unlinked`, `target-missing`, `matrix-forbidden … spec-straten → adr`, je Exit 1 | bestätigt |
| Belege in §7 je Zeile (Mutation · Befundzeile) | alle 14 Zeilen der drei Tabellen nachgefahren (Anker, Linkziel, ADR-Link aus `spec/`, ADR in Inline-Code in `spec/` → `id-unlinked`, `MR-` in `spec/` in Inline-Code, Slice in `spec/`, `LH-` blank grün, `MR-` außerhalb `spec/` grün, Pfad in Inline-Code grün, Link auf superseded ADR `matrix-inactive`, Link auf die ersetzende ADR grün, Ausnahme `docs/reviews/` grün, Zaun `~~~` grün, Zaun in `spec/` grün): jedes Ergebnis wie in §7 | bestätigt |
| Block sagt nur zu, was `.d-check.yml` aktiv prüft und eine Zeile der Gegenprobe rot zeigt („heute: `ADR-` ohne Link im Fließtext und in Inline-Code jeder gescannten `.md`“) | eigene Mutationen außerhalb der Gegenprobe: ADR-Kennung in einer **Überschrift** (`harness/README.md`) grün; ADR-Kennung blank im **Fließtext einer ADR** (Abschnitt *Kontext* von ADR 0002) und im **ADR-Index** grün; Slice-Kennung in `spec/architecture.md` und Link auf eine superseded ADR in `harness/README.md` je in einem Abschnitt **`## Geschichte`** grün. Weder DoD-Punkt 1 noch der Block nennen diese Fälle unter *nicht* | **abgelehnt**, V-134 |
| Fundstellen: kein gleichlautender Block unter `.claude/agents/` und `.harness/skills/` | `grep -rn -i "Strenges Doc-Gate"` über `.claude/`, `.harness/skills/`, `AGENTS.md`: nur die drei Commands | bestätigt |

Gegenkontrollen, die die Zusagen des Blocks stützen: ADR-Kennung in einem eingerückten Code-Block (vier Leerzeichen) rot — der Block nimmt nur Zäune aus, das stimmt; Zaun mit vier Backticks grün; ADR-Kennung in einer Tabellenzelle rot; in `docs/reviews/` rot; im Abschnitt `## Geschichte` von `harness/README.md` rot (`ids` kennt die Abschnitts-Ausnahme nicht, nur `matrix`).

### Punkt 2: Freshness-Audit. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Release-Liste: neuester Tag gegen `v6.16.0` | eigener Lauf `gh release list -R pt9912/ai-harness-course`: `v6.17.0` Latest 2026-10-07T06:44Z, davor `v6.16.0` 2026-10-06T12:49Z; Delta genau ein Release | bestätigt |
| Delta inhaltlich | `lab-regelwerk.zip` von `v6.17.0` in den Scratchpad (`sha256 afe50df8…`, gleich `SHA256SUMS` des Releases), `diff -r` gegen `.harness/baseline/v6.16.0/`: dieselben Dateien (nur `SHA256SUMS` liegt zusätzlich im vendored Baum); Zeilen über Quell-Zeile und Tag hinaus nur in `modul-13-quality-gates.md`, `grundlagen-harness-dateien.md` (Disjunktheit der Teile eines Gate-Index), `regelwerk/README.md` (Stand-Zeile *Kurs-Welle 160*) und `templates/.d-check.yml` (`authority-disjoint`) — deckt sich mit §7 | bestätigt |
| je Eintrag des Adaptions-Blocks ein Ausgang | `MR-000` *bleibt gültig*, `MR-001` (aufgelöst) *bleibt gültig* ohne Nachfolger; aktive keine. Beide Ausgänge aus der Menge von Modul 2 und vom Delta getragen | bestätigt |
| Regelblock-Tabelle: Umbenennung, Zugang, Wegfall | keiner; Werte der berührten Zeilen bleiben, weil `.d-check.yml` weder `targets` noch `authority` führt (eigene Lesung) | bestätigt |
| ohne Änderung am vendored Baum | `make baseline-verify` in `make gates`: `v6.16.0 OK — 54 Dateien`; der Diff berührt `.harness/baseline/` nicht | bestätigt |

### Punkt 3: Reviewer-Skill *Nehmer nicht nachgezählt*. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| §Klassifikation (MEDIUM) führt die Klasse | `.harness/skills/reviewer.md` am Stand `cc2b61c`, unter MEDIUM nach *Adresse nimmt nicht an* | bestätigt |
| Auslöser deckungsgleich mit `AGENTS.md` §3.13 *Nachzählen beim Eintragen* | Regel: wer einträgt, „zählt im selben Commit dessen Liefer-Punkte und Schichten nach, mit der Sendung“, Grenze drei und zwei nach Modul 5, „läge der Nehmer darüber, geht der Punkt an den Planner“. Klasse: Sendung in §1 oder DoD des Nehmers, §8 nennt die Zählung „nicht im selben Commit“, oder der Nehmer liegt mit ihr darüber; Zählung nach Modul 5 §Ziel-Form: Slice; mindestens MEDIUM; Herkunfts-Anker `seit slice-v1-abschluss-einspielen` gleich dem der Regel. F-560 ist damit erledigt | bestätigt |
| Wortlaut DoD-Punkt 3 = Wortlaut Skill | beide mit „im selben Commit“, beide „über drei Punkten oder zwei Schichten“ | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt durch eigenen Lauf; Beleg in §7 nur an `79d7bb8` (V-136).**

Eigener Lauf am Stand `cc2b61c`, Arbeitsbaum sauber: Exit 0. Darin `baseline-verify: v6.16.0 OK — 54 Dateien`, `d-check: 438 Datei(en) geprüft, 0 Befund(e)`, `run-integration-tests: gruen`, `a-check-negativ`, `commit-msg-gegenprobe`, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`, `lint-gegenprobe` je `gruen`.

### Übrige DoD-Punkte (konstant je Slice)

Review-Report liegt vor (`13b17dc`). Closure-Notiz, Register, Risiko-Ausgänge und Paarungen sind Sache der Closure und hier nicht zu bestätigen; beide Risiken in §6 tragen heute *offen bis Closure*.

## 2. Plan gegen Diff

| Stelle | Ergebnis |
|---|---|
| §1 Ziel und Herkunft | folgt dem Diff: Pin, Block in drei Commands, Audit in §7, neue Skill-Klasse; *Herkunft der drei Teile und der Berichtigung* (F-559 erledigt); Berichtigung nennt jetzt Schicht-Abgrenzung und §3 als Grund (F-558 erledigt) |
| §1 Schicht-Abgrenzung | „Der Slice ändert einen Pin, einen Skill, einen Block in drei Workflow-Commands und diesen Plan“ — der Diff legt dazu zwei Slice-Pläne in `open/` an und schreibt zwei Drift-Log-Zeilen der Roadmap. Nach der Schichtteilung ist das Planung, keine Schicht; als Aufzählung des Geänderten ist der Satz aber unvollständig (V-135) |
| §3 | Zeilen für `d-check.mk`, Skill, Commands, §7 folgen dem Diff; die Nehmer-Pläne und die Roadmap fehlen (V-135). `.d-check.yml` ist nicht im Diff, wie §1 sagt |
| §6 | Randformen „keine“ trägt: kein neuer Vertrag im Code-Sinn, der Block ist Text über ein bestehendes Gate. Beide Risiken mit Stand; Risiko 2 verweist auf einen Upgrade-Slice, den es noch nicht gibt — als Ausgang bei der Closure zu benennen |
| §8 | Schichtteilung identisch in allen drei Plänen (wörtlich verglichen); Geber zählt eine Schicht (Harness) und drei Liefer-Punkte — nachvollzogen |

**§3.13 für beide Nehmer, nachvollzogen:**

| Nehmer | Annahme (§1 und DoD mit Kennung des Gebers, kein Ausschluss trifft, nicht in `done/`) | Nachzählen im selben Commit | Zählung nach der Schichtteilung |
|---|---|---|---|
| `slice-harness-lh-links-pflicht` (Geber `slice-harness-d-check-v0-85`) | §1 *Sendung aus slice-harness-d-check-v0-85*, DoD *Muster aktiv* und *Aussage über das Gate* mit Kennung; kein Punkt unter *Ausdrücklich NICHT* trifft; `open/` | Sendung in `dfb7e88` mit Zählung in §8 (damals mit abweichender Teilung, F-556); neu gezählt in `6db272b` | drei Liefer-Punkte; Schichten: Spezifikation (neue Stelle §11, Links) und Harness (`.d-check.yml`, Commands, `harness/`) = zwei; Pläne und Roadmap Planung. Nachvollzogen |
| `slice-harness-lh-links-bestand` (Geber `slice-harness-lh-links-pflicht`) | §1 *Sendung aus slice-harness-lh-links-pflicht*, beide DoD-Punkte mit Kennung; der Ausschluss *Wortlaut einer ADR ändern* trifft die Sendung nicht, solange ein Link keine Aussage ändert — die Entscheidung dazu ist vor `next` dem Nutzer zugewiesen (§6, §4); `open/` | Anlage, Sendung und §8-Zählung in demselben Commit `6db272b` | zwei Liefer-Punkte; Schichten Entscheidungen (`docs/plan/adr/`) und Nutzer- und Wartungs-Doku = zwei. Nachvollzogen |

Die Messung in §1 von `slice-harness-lh-links-bestand` (elf Kennungen in vier ADRs, sechs in `docs/maintainer/releasing.md`) habe ich nicht nachgemessen; sie wird beim Start neu gemessen (§4 dort).

## 3. Übergaben des Reviews

| Finding | Übergabe | Stand an `cc2b61c` |
|---|---|---|
| F-555 (Code-Block) | Implementer | erledigt: alle drei Commands nehmen Zäune aus drei Backticks und `~~~` aus, DoD-Punkt 1 ebenso; Gegenprobe in §7 je Zaun und in `spec/`, von mir nachgefahren. Dieselbe Klasse tritt an drei weiteren Formen wieder auf (V-134) |
| F-556 (Schichtteilung) | Planner | erledigt: gemeinsame Teilung in drei Plänen, Nehmer geschnitten, Drift-Log-Zeile; §3.13 nachvollzogen (Abschnitt 2) |
| F-557 (Code bei ADR aus `spec/`) | Implementer | erledigt: `implement-slice.md` nennt `id-unlinked`; Mutation nachgefahren |
| F-558 (Begründung der Berichtigung) | Planner | erledigt (§1) |
| F-559 (§3, §1 Überschrift) | Planner/Implementer | erledigt |
| F-560 („im selben Commit“) | Planner | erledigt in Skill und DoD-Punkt 3 (Punkt 3 oben) |
| F-561 (superseded ADR) | Implementer | erledigt: Code `matrix-inactive`, Ausweg „ersetzende ADR“, Ausnahmen ADR-Index und `docs/reviews/`; nachgefahren. Die Ausnahme der Abschnitte `## Geschichte` fehlt (V-134) |
| F-562 (Zuordnung) | Closure | eingeordnet: Die LH-Fehllesung kam aus einem auskommentierten Muster und dem alten Block, beide Versionen identisch; sie ist kein Beleg für `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (Werkzeug wertet eine *aktive* Regel anders aus), sondern für `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. V-134 dagegen trifft eine aktive Regel (`link-policy: always` wirkt nicht in Überschriften und nicht im Zielverzeichnis); ob die Closure V-134 dem ersten Eintrag zuordnet (dann 3×), ist ihr Urteil |
| F-563 (an den Verifier) | Verifier | Pin, Audit, Berichtigung und Annahme selbst nachgeprüft (Abschnitte 1 und 2) |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-134 | MEDIUM | DoD-Punkt 1 und der Block *Strenges Doc-Gate* sagen rot zu für „`ADR-` ohne Link im Fließtext und in Inline-Code jeder gescannten `.md`“ bzw. „in einer gescannten `.md`“; grün bleiben eine ADR-Kennung in einer Überschrift, eine blanke ADR-Kennung im Fließtext einer Datei unter `docs/plan/adr/` (ADR und Index). `implement-slice.md` sagt zudem `matrix-forbidden` für eine Slice-Kennung aus `spec/` und `matrix-inactive` für einen Link auf eine superseded ADR außerhalb von Index und `docs/reviews/` zu; in einem Abschnitt `## Geschichte` bleiben beide grün (`exclude-sections` in `.d-check.yml`). Keiner der Fälle steht unter *nicht*, keiner in der Gegenprobe. *Failure-Szenario:* Ein Architect schreibt eine ADR mit blanker Kennung einer anderen ADR, oder ein Agent schreibt eine Kennung in eine Überschrift, liest den Block als Zusage und das Gate bleibt grün — dieselbe Klasse wie F-555, nach dessen Nacharbeit. | `AGENTS.md` §3.11; Plan §2 DoD-Punkt 1 („sagt … nur zu, was … eine Zeile der Gegenprobe in §7 rot zeigt“); `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `docs/plan/planning/in-progress/slice-harness-d-check-v0-85.md` §2 · „in Inline-Code jeder gescannten `.md`“; `.claude/commands/implement-slice.md` · „in einer gescannten `.md`, im Fließtext und in Inline-Code (`id-unlinked`)“; `.claude/commands/plan-welle.md`, `.claude/commands/close-welle.md` · „ohne Link in einer gescannten `.md`, im Fließtext und in Inline-Code“ | ja (fünf Kopien: Überschrift in `harness/README.md`, *Kontext* von ADR 0002, ADR-Index, `## Geschichte` mit Slice-Kennung in `spec/architecture.md`, `## Geschichte` mit Link auf ADR 0021 in `harness/README.md` — je 0 Befunde, Exit 0) | Kommentar-Zusage weiter als die Prüfung |
| V-135 | LOW | §1 *Schicht-Abgrenzung* zählt auf, was der Slice ändert („einen Pin, einen Skill, einen Block in drei Workflow-Commands und diesen Plan“), und §3 führt dieselben Dateien; der Diff legt zusätzlich zwei Slice-Pläne in `open/` an und schreibt zwei Drift-Log-Zeilen der Roadmap. Die Schichtzahl ändert das nicht (Planung ist keine Schicht), die Aufzählung ist aber nicht mehr der Diff. | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | `docs/plan/planning/in-progress/slice-harness-d-check-v0-85.md` §1 · „und diesen Plan“; §3 · „dieser Plan, §7“ | nein (Lesen: `git diff --stat afc2528..cc2b61c`) | Plan folgt der Korrektur nicht |
| V-136 | LOW | §7 *Läufe* belegt `make gates` nur am Stand `79d7bb8`; nach Berichtigung, Block und beiden Nacharbeiten steht dort „Der letzte Lauf nach dem Commit dieser Belege steht im Bericht“. Der Bericht ist kein Artefakt. Mein eigener Lauf an `cc2b61c` ist grün (Punkt 4). | Auftrag des Verifiers (Belege in §7, seit slice-harness-blackbox-kern); Plan §2 „`make gates` grün“ | `docs/plan/planning/in-progress/slice-harness-d-check-v0-85.md` §7 · „Der letzte Lauf nach dem Commit dieser Belege steht im Bericht“ | ja (§7 lesen) | DoD-Beleg nur im Bericht, nicht in §7 |
| V-137 | INFO | Für den Nehmer `slice-harness-lh-links-pflicht`: Das `ids`-Modul prüft unter `v0.85.0` keine Überschriften und kein Zielverzeichnis eines Musters (V-134). Für das LH-Muster wären das Überschriften jeder Datei und `spec/lastenheft.md`; die Randform *Datei der Definition* in §6 dort nennt das Zielverzeichnis, Überschriften nennt sie nicht. Für den Planner, vor dem Architect. | `AGENTS.md` §3.12 | `docs/plan/planning/open/slice-harness-lh-links-pflicht.md` §6 · „**Datei der Definition:**“ | ja (Kopie mit Kennung in Überschrift) | — (Hinweis an einen Nehmer) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Pin in `d-check.mk` | geprüft, ohne Befund. Tag und Index-Digest selbst gegen die Registry geprüft; Kopfkommentar wahr; nur zwei Zeilen geändert. |
| `.d-check.yml` | geprüft, ohne Befund. Nicht im Diff, wie §1 sagt; `AGENTS.md` §3.6 nicht berührt. |
| Gegenprobe in §7 | geprüft, ohne Befund. Alle Zeilen der drei Tabellen nachgefahren, jedes Ergebnis wie notiert. |
| Freshness-Audit | geprüft, ohne Befund. Release-Liste und Delta selbst nachgefahren; vendored Baum unverändert. |
| Reviewer-Skill | geprüft, ohne Befund. Klasse deckungsgleich mit `AGENTS.md` §3.13 *Nachzählen beim Eintragen*. |
| §3.13 für beide Nehmer | geprüft, ohne Befund. Annahme und Nachzählen je im Commit der Sendung, Zählung nach der gemeinsamen Teilung nachvollzogen. |
| Commit-Messages `4c7ae6a` bis `cc2b61c` | geprüft, ohne Befund. Jede nennt `slice-harness-d-check-v0-85` und `MR-000`, keine eine Struktur-Kennung. |
| Hard Rule 3.3 | geprüft, ohne Befund. Kein Move im Diff. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Kommentar-Zusage weiter als die Prüfung · Plan folgt der Korrektur nicht · DoD-Beleg nur im Bericht, nicht in §7

Zuordnung für die Closure: V-134 ist eine Ausprägung von `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` und berührt `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (siehe F-562 oben); V-135 von `BEO-REPO/plan-folgt-korrektur-nicht`.

## Verdikt

**Urteil je DoD-Liefer-Punkt:** Punkt 1 teilweise bestätigt — Pin, Digest, 0 Befunde und Gegenprobe halten; die Teil-Zusage „der Block sagt nur zu, was geprüft ist“ ist **abgelehnt** (V-134). Punkt 2 bestätigt. Punkt 3 bestätigt. Punkt 4 durch eigenen Lauf bestätigt (Beleg in §7 veraltet, V-136).

**Closure-blockierend:** V-134 ja, bis Block und DoD-Punkt 1 die drei Formen ausnehmen oder die Gegenprobe sie rot zeigt; V-135 und V-136 nein, sie gehen vor der Closure in den Plan.

**Übergabe:**

- V-134 an den Implementer: den Block in den drei Commands und, über den Planner, DoD-Punkt 1 auf das Geprüfte fassen (Überschriften, `docs/plan/adr/`, Abschnitte `## Geschichte` unter *nicht*), je Form eine Zeile der Gegenprobe in §7. Ein Widerspruch läuft über den Architect.
- V-135 an den Planner: §1 *Schicht-Abgrenzung* und §3 um die beiden Nehmer-Pläne und die Roadmap ergänzen.
- V-136 an den Implementer: den letzten `make gates`-Lauf mit Stand und Exit-Code in §7 eintragen.
- V-137 an den Planner (Randform im Nehmer `slice-harness-lh-links-pflicht`).
- F-562 und die Zuordnung oben an die Closure.

Dieser Bericht ersetzt kein Review und keine Validierung.

---

## Nachprüfung — 2026-10-09

**Gegenstand:** die Nacharbeit `24df901` (V-134, V-135), `6565310` (V-137) und `23a5314` (V-136) am Stand `23a5314`; Diff `9293191..23a5314` berührt nur die drei Commands, den Plan dieses Slice und `slice-harness-lh-links-pflicht`. Pin, Skill und Audit sind nicht im Diff, ihr Urteil von oben bleibt. Beim Start war der Arbeitsbaum sauber, HEAD `23a5314`.

**Proben wie oben**, Präfix `ver2-dcheck-`: je Mutation eine frische Kopie (jeder Eintrag der obersten Ebene außer `.git` mit `cp -r` ohne `-p`), Mutation per `printf` angehängt bzw. als neue Datei, Lauf `docker run --rm --network none` mit dem Digest aus `d-check.mk` (`c07f1fe6…`), danach nur diese Kopie gelöscht. Unveränderte Kopie: `439 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Rot heißt unten mindestens 1 Befund, Exit 1; grün 0 Befunde, Exit 0. Die Kennung `ADR-` ist die mit vier Ziffern 0001, die Slice-Kennung die dieses Slice.

### Stichprobe der Gegenprobe in §7 (Nacharbeit nach V-134)

Acht Zeilen nachgefahren, jedes Ergebnis wie in §7: Absatz rot (`id-unlinked`); `## Probe` mit der Kennung grün; Kennung als Link-Text grün; Kennung blank in ADR 0002 grün; `## Geschichte` mit Slice-Kennung in `spec/architecture.md` grün; `## Geschichte` mit Link auf ADR 0021 in `harness/README.md` grün; Link auf ADR 0021 aus ADR 0010 grün; `## 9. Geschichte` rot (`matrix-forbidden`). Die fünf Formen von V-134 nennen Block und DoD-Punkt 1 jetzt unter *nicht*, je mit grüner Zeile in §7. **V-134 ist erledigt.**

### Weitere Formen, aktiv gesucht

| Form | Mutation | Ergebnis | Block sagt |
|---|---|---|---|
| verschachtelte Liste, nummerierte Liste, Liste im Zitat, Zitat im Zitat | Kennung blank darin | je rot | rot (Listenpunkt, Zitat) — trägt |
| Fett, kursiv, durchgestrichen, Folgezeile eines Absatzes, Tabelle ohne Rand-Pipe, CRLF | Kennung blank darin | je rot | rot — trägt |
| Kennung neben einem Link, neben einem Link auf eine andere ADR, nach früherem Link auf dieselbe ADR | im selben Absatz bzw. später | je rot | rot — trägt (jedes Vorkommen zählt) |
| Text eines Referenz-Links `[…][r]`, Kurzform `[…]` mit Definition, Fußnote, mehrzeiliger HTML-Kommentar, `<details>`, `<span>`, Doppel-Backtick-Inline-Code, Front-Matter `---` und `+++` | Kennung blank darin | je rot | Referenz-Link-Text: *nicht* (V-140); übrige ungenannt oder rot — trägt |
| Überschrift in einem Listenpunkt (`- ## …`) und in einem Zitat (`> ## …`) | Kennung darin | je rot | *nicht* (`#`-Überschrift) — Ausnahme weiter als das Gate (V-140) |
| Abschnitt `Geschichte` als Setext-Überschrift; `### Geschichte`; `# Geschichte` mit Unterabschnitt | Slice-Kennung in `spec/architecture.md` | rot; grün; grün (auch im Unterabschnitt) | Setext: *nicht* — weiter als das Gate (V-140); Ebenen: trägt |
| Zaun mit vier Tilden, Zaun im Zitat, Zaun im Listenpunkt | Kennung darin | je grün | *nicht* (Zaun) — trägt |
| offener Zaun ohne Ende | Kennung darin | rot, `fence-unclosed` | ungenannt — der Block sagt weniger zu, trägt |
| Titel eines Links `[x](conventions.md "…")` | Kennung im Titel, `harness/README.md` | **grün** | rot (Absatz) — V-138 |
| Slice-Kennung als Link-Text, als Bild-Text, im Titel eines Links, `MR-001` als Link-Text, je in `spec/` mit Ziel in `spec/` | angehängt an `spec/architecture.md` bzw. `spec/spezifikation.md` | je **grün** | rot (`matrix-forbidden`) — V-138 |
| `<a href>` aus `spec/` auf ADR 0001; externe URL auf ADR 0001 | an `spec/architecture.md` | je **grün** | rot (Link aus `spec/` auf eine ADR) — V-138 |
| `<a href>` auf ADR 0021; externe URL auf ADR 0021 | an `harness/README.md` | je **grün** | rot (`matrix-inactive`) — V-138 |
| Anker in ein Ziel, das keine `.md` ist (`../Makefile#nope`, `../go.mod#L1`, `a.txt#nope`) | an `harness/README.md` | je **grün** | rot (toter Anker), alle drei Commands — V-138 |
| Link-Referenz, Ziel in spitzen Klammern, mit Titel, mit Query, absoluter Pfad, Verzeichnis, Anker in derselben Datei, Anker in Großschreibung | totes Ziel bzw. toter Anker | je rot | rot — trägt |
| Referenz-Link, Bild, Überschrift, Anker aus `spec/` auf ADR 0001 bzw. auf ADR 0021 | Link | je rot | rot — trägt |
| Text eines HTML-Links `<a href="conventions.md">…</a>` | `ADR-`-Kennung als Text | rot | ungenannt (*HTML-Links* meint das Ziel) — trägt |
| `docs/vendor/`, `docs/node_modules/`, `build/`, `docs/dist/`, `docs/target/`, `docs/.venv/`, `docs/__pycache__/`, `vendor/`, `node_modules/` | neue `.md` mit der Kennung blank | je **grün**, 439 Dateien (nicht gezählt) | gescannt, rot — V-139 |
| `third_party/`, `.foo/`, `.githooks/`, `test/testdata/`, `docs/out/`, `bin/`, `tmp/`, `venv/`, `.cache/`, `coverage/` | dasselbe | je rot, 440 Dateien | gescannt — trägt |
| `harness/probe.MD`, `harness/probe.markdown` | dasselbe | je grün, 439 Dateien | „jede `.md`“ — trägt wörtlich |
| `ADR\-0001`, `slice\-harness` (maskierter Bindestrich) | im Absatz bzw. in `spec/` | je grün | — (das Token steht nicht im Quelltext; nur notiert, V-140) |
| URL eines Links, `file:///…`, `mailto:` | Kennung in der URL bzw. totes Ziel | je grün | *nicht* (externe URL) — trägt |
| `../../../etc/passwd` als Ziel | Link | rot, `repo-escape` | ungenannt — trägt |
| Zeilen-Marker `<!-- d-check:ignore … -->` in der Zeile davor bzw. danach | Kennung blank | je rot | nur „in einer Zeile mit“ — trägt |

**Zur Richtung:** Der Block erklärt alles nicht als rot Genannte für ungeprüft. Eine Ausnahme, die das Gate in einer Form doch rot zeigt (V-140), bricht daher keine Zusage; eine Form, die der Block rot nennt und die grün bleibt (V-138, V-139), bricht DoD-Punkt 1. Die Lücken in V-138 liegen fast alle in `implement-slice.md`: Dort stehen *externe URLs, HTML-Links* nur beim Spiegelstrich der Linkziele und *Text eines Links oder Bildes* nur beim Spiegelstrich der `ADR-`-Kennung; die Spiegelstriche zu `matrix-forbidden` und `matrix-inactive` erben sie nicht, und *Für alle vier* nennt nur Zäune. In `plan-welle.md` und `close-welle.md` stehen die Ausnahmen ungebunden; dort trägt V-138 nur der Link-Titel und der Anker in ein Ziel ohne `.md`.

### V-135 bis V-137

| Befund | Stand an `23a5314` |
|---|---|
| V-135 | erledigt: §1 *Schicht-Abgrenzung* nennt die beiden Nehmer-Pläne und die zwei Drift-Log-Zeilen, §3 führt beide als Zeilen (Planung) |
| V-136 | erledigt: §7 *Läufe* nennt `make gates` am Stand `6565310`, Exit 0, mit den Zeilen der Teil-Gates; `23a5314` fügt nur diesen Absatz an (6 Zeilen, nur der Plan). Mein eigener Lauf an `23a5314` steht unten |
| V-137 | erledigt: `slice-harness-lh-links-pflicht` §6 nennt die Ausnahmen des `ids`-Moduls mit der Kennung des Gebers und die Datei der Definition samt Unterordnern; §8 begründet, warum die Zählung bleibt (eine Frage an den Architect, kein Liefer-Punkt) — nachvollzogen. Die Liste der nicht gescannten Pfade dort ist unvollständig wie im Block (V-139) |

### Urteil je DoD-Liefer-Punkt (aktualisiert)

- **Punkt 1** — Pin, Digest, 0 Befunde, Gegenprobe der vier Klassen: bestätigt (unverändert). Block und DoD-Wortlaut: V-134 erledigt; die Teil-Zusage „sagt nur zu, was … rot zeigt“ ist **weiter abgelehnt**, enger als zuvor (V-138, V-139).
- **Punkt 2** — bestätigt (nicht im Diff der Nacharbeit).
- **Punkt 3** — bestätigt (nicht im Diff der Nacharbeit).
- **Punkt 4** — eigener Lauf `make gates` an `23a5314`, Arbeitsbaum sauber: Exit 0; darin `baseline-verify: v6.16.0 OK — 54 Dateien`, `d-check: 439 Datei(en) geprüft, 0 Befund(e)`, `run-integration-tests: gruen`, je `gruen` `a-check-negativ`, `commit-msg-gegenprobe`, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`, `lint-gegenprobe` — bestätigt. Der Beleg in §7 (Stand `6565310`) ist für den Code-Stand gültig: Danach änderte nur `23a5314` einen Absatz des Plans.

## Findings der Nachprüfung

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-138 | MEDIUM | Der Block sagt rot zu für Formen, die grün bleiben: (a) in `implement-slice.md` beim Spiegelstrich `matrix-forbidden` eine Slice- oder `MR-`-Kennung in `spec/` als Text oder Titel eines Links oder als Text eines Bildes (Ziel in `spec/`), und einen Link aus `spec/` auf eine ADR als `<a href>` oder als externe URL; beim Spiegelstrich `matrix-inactive` den Link auf eine superseded ADR als `<a href>` oder externe URL. Die Ausnahmen *HTML-Links, externe URLs* und *Text eines Links oder Bildes* stehen nur bei den Spiegelstrichen der Linkziele bzw. der `ADR-`-Kennung, *Für alle vier* nennt nur Zäune. (b) In allen drei Commands den Titel eines Links (`[x](ziel "…")`) mit einer `ADR-`-Kennung — Fließtext im Absatz, *nicht* nennt nur den Text — und einen toten Anker in ein Ziel, das keine `.md` ist (`../go.mod#L1`). *Failure-Szenario:* Ein Implementer verlinkt aus `spec/` per `<a href>` auf eine ADR oder setzt eine Slice-Kennung als Link-Text in die Spezifikation, liest den Spiegelstrich als Zusage und das Gate bleibt grün — dieselbe Klasse wie V-134, nach dessen Nacharbeit. Weg, der die Klasse schließt statt der Form: die Ausnahmen nach Art des Befunds bündeln statt je Spiegelstrich — *externe URLs und HTML-Links als Ziel* für jeden Link-Befund (`target-missing`, `anchor-missing`, Link aus `spec/`, `matrix-inactive`; je grün gemessen), *Text und Titel eines Inline-Links oder Bildes* für jeden Kennungs-Befund (`id-unlinked`, Slice- und `MR-`-Kennung in `spec/`; je grün gemessen), dazu „Anker nur in `.md`-Zielen“; je eine grüne Zeile in §7. Nicht ausgenommen ist der Text eines HTML-Links: `<a href="…">` mit der `ADR-`-Kennung als Text ist rot. | `AGENTS.md` §3.11; Plan §2 DoD-Punkt 1 („sagt … nur zu, was … eine Zeile der Gegenprobe in §7 rot zeigt“); `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `.claude/commands/implement-slice.md` · „bei einem Link aus `spec/` auf eine ADR und einer Slice- oder `MR-`-Kennung in `spec/`, auch in“; · „bei einem Link auf eine superseded ADR (`matrix-inactive`“; · „Für alle vier nicht geprüft: umzäunte Code-Blöcke“; `.claude/commands/plan-welle.md`, `.claude/commands/close-welle.md` · „als Text eines Links oder Bildes“ und „toten relativen Linkziel oder Anker“ | ja (zwölf Kopien, Tabelle oben, je 0 Befunde, Exit 0) | Kommentar-Zusage weiter als die Prüfung |
| V-139 | LOW | „Gescannt wird jede `.md` außer `.harness/**` …, `.tmp/**` und `**/*.template.md`“ (alle drei Commands) und DoD-Punkt 1 „nicht gescannte Pfade (`scan.ignore`)“: d-check lässt zusätzlich, ohne Eintrag in `scan.ignore`, Verzeichnisse namens `vendor`, `node_modules`, `build`, `dist`, `target`, `.venv` und `__pycache__` in jeder Tiefe aus; eine `.md` darin wird nicht gezählt. Heute liegt keines im Repo (`find`), ein `vendor/` kann aber `go mod vendor` anlegen. Dieselbe unvollständige Liste trägt `slice-harness-lh-links-pflicht` §6 *Ausnahmen des `ids`-Moduls* (an den Planner). | `AGENTS.md` §3.11; Plan §2 DoD-Punkt 1 | `.claude/commands/implement-slice.md` · „Gescannt wird jede `.md` außer“; `plan-welle.md`, `close-welle.md` · „(jede außer `.harness/**`,“; Plan §2 · „nicht gescannte Pfade (`scan.ignore`)“ | ja (neun Kopien, je 439 Dateien, 0 Befunde) | Kommentar-Zusage weiter als die Prüfung |
| V-140 | INFO | Drei genannte Ausnahmen sind weiter als das Gate (Richtung harmlos, der Block sagt weniger zu): eine `#`-Überschrift in einem Listenpunkt oder Zitat ist rot; der Text eines Referenz-Links `[…][r]` ist rot; eine Setext-Überschrift `Geschichte` nimmt nichts aus. Dazu grün, aber nicht zu fassen: ein maskierter Bindestrich (`ADR\-0001`) — das Token steht nicht im Quelltext. Kein Handlungsbedarf für die Closure; bei einer Neufassung genügt „`#`-Überschrift am Zeilenanfang“, „Text eines Inline-Links“. | Plan §2 DoD-Punkt 1 („nennt … die Ausnahmen klassenweise“) | `.claude/commands/implement-slice.md` · „Nicht geprüft: in einer `#`-Überschrift, als Text eines Links oder Bildes“; · „dessen Überschrift genau `Geschichte` lautet“ | ja (Tabelle oben) | — |

## Summary der Nachprüfung

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen:** Kommentar-Zusage weiter als die Prüfung (V-138, V-139; zweiter Lauf in Folge, Zuordnung zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` wie V-134).

## Verdikt der Nachprüfung

V-134 bis V-137 erledigt. **DoD-Punkt 1 bleibt in der Teil-Zusage über den Block abgelehnt** (V-138), Punkte 2 bis 4 bestätigt. **Closure-blockierend:** V-138 ja, bis die Ausnahmen nach Art des Befunds für alle Spiegelstriche gelten, Link-Titel und Anker in Nicht-`.md`-Ziele genannt sind und §7 je eine grüne Zeile führt; V-139 nein, aber vor der Closure in Block, DoD-Wortlaut und Nehmer-§6 nachzuziehen; V-140 nein.

**Übergabe:** V-138 und V-139 (Block) an den Implementer, DoD-Wortlaut über den Planner; V-139 (Nehmer `slice-harness-lh-links-pflicht` §6) an den Planner; V-140 zur Kenntnis. Da die Form-Suche zweimal neue Lücken fand, empfehle ich die Ausnahmen je Art des Befunds (Link oder Kennung) statt je Spiegelstrich — sie schließt die Klasse, nicht die nächste Form.
