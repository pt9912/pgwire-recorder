# Slice slice-harness-commit-struktur-id: Commit-Träger lehnt Struktur-Kennungen ab

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-07: nach `slice-harness-lint` und vor `slice-harness-abdeckung-gate`
(WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** kein Lastenheft-Bezug — der Commit-Träger ist ein Werkzeug der
Prüfumgebung, keine Zusage des Produkts. Bindung an Entscheidungen:
[ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) (benannte
Slice-Kennungen im Träger, Lesebereich bis zur Scissors-Zeile, ohne Kommentarzeilen),
[ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) (benannte
Welle-Kennungen), [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) (Proposed, Architect vor dem Code:
Ablehnung von Struktur-Kennungen, ergänzt beide; Annahme durch den Nutzer offen).

**Berührte Spec-Stellen:** [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--struktur-kennungen-im-commit-träger-commit-msg) (Commit-Träger, vom Architect vor dem Code
geschrieben) · `spezifikation.md` §12 (*Historie*)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Commit-Träger `.githooks/commit-msg` lehnt eine Commit-Message ab, die
eine Struktur-Kennung der Spezifikation oder der Sicht (`SPEC-<NNN>`, `ARC-<NNN>`)
nennt: Exit 1 und je Nennung eine Meldung mit Zeile und Kennung. Grundlage ist
`AGENTS.md` §5 Regel 1 (Struktur-IDs adressieren innerhalb der Spec und gehören nicht
in die Commit-Message). Regel und Randformen stehen in `SPEC-050`, vom Architect vor
dem ersten Code-Commit geschrieben (`AGENTS.md` §3.12), die Entscheidung in
[ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md); `make hook-gegenprobe` führt je Randform einen Fall. Die Prüfung liegt
im Träger selbst (`SPEC-050` Punkt 1); `tools/harness/commit-msg-traceability.sh` bleibt
unberührt.

**Herkunft:** `BEO-REPO/commit-nennt-struktur-kennung` steht mit dem Nachtrag von F-336
bei 3× (`slice-replay-semantik-meldungscodes` V-61, `slice-extended-query-replay`
F-336 zu `c8ebf08`, `slice-harness-blackbox-einstieg` F-454). Entscheidung des Nutzers
vom 2026-10-07: Ausgang *geplant* auf diesen Slice. Dreimal hat das Review das Muster
erst nach dem Commit gefunden; ein Sensor im Commit-Pfad fängt es davor.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Bestehende Commits umschreiben — Bestand bleibt bewusst stehen: Die betroffenen
  Commits sind gepusht, und ein Umschreiben der Historie bräche jeden Klon; die
  Prüfung wirkt ab dem aktivierten Träger.
- Weitere ausgeschlossene Klassen (etwa `BEO-…`, `CO-…`, Review-Nummern `F-…`,
  `V-…`) — anderer Vorgang: `AGENTS.md` §5 Regel 1 schließt nur Struktur-IDs aus;
  eine weitere Klasse wäre eine neue Regel mit eigener Herkunft.
- Ob eine zugelassene Kennung wahr ist, über
  [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) und
  [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) hinaus —
  Bestand bleibt stehen: Die Grenze im Kopf der Prüfung gilt weiter.
- Ein Träger außerhalb des Klons (CI, Push-Schutz) — anderer Vorgang:
  `git commit --no-verify` und ein Klon ohne `make hooks-install` bleiben ungeprüft,
  wie `harness/mk/hooks-install.mk` sagt.
- Die Annahme nach [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) und [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) in die
  Spezifikation übertragen — Bestand bleibt bewusst stehen: `SPEC-050` legt nur die
  Ablehnung und ihren Vorrang fest. Eine Übertragung müsste die Wortgrenze der Annahme
  entscheiden, die der Träger heute nur vorn prüft (`slice-harness-lint_x` nimmt er als
  `slice-harness-lint` an); das wäre eine Änderung der Annahme, nicht dieser Slice.
- Die mitgelieferte Prüfung `tools/harness/commit-msg-traceability.sh` — Bestand bleibt
  stehen: Ein erneuter Bootstrap schreibt sie neu; die Ablehnung liegt darum im Träger.
- Produkt-Code, Lastenheft und Sicht — Schicht-Abgrenzung: Der Slice ändert Träger,
  Gegenprobe, Doku und in der Spezifikation nur §11 und §12 (Architect).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Vertrag und Träger: `SPEC-050` führt die Ablehnung und jede Randform aus §6 als
      Entscheidung, dazu die Grenzen (Architect, vor dem ersten Code-Commit, mit
      Historien-Zeile in §12 und [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)). Der Träger lehnt danach jede Message
      ab, die nach `SPEC-050` eine Struktur-Kennung nennt, mit Ausgabe und Ausgang nach
      Punkt 6, auch wenn eine zugelassene Kennung daneben steht (Punkt 5); jede Message
      ohne Struktur-Kennung und jede nach Punkt 4 ausgenommene behandelt er wie heute.
- [ ] Gegenprobe: `make hook-gegenprobe` führt je Randform aus §6 mindestens einen Fall
      mit der Erwartung aus §11 (abgelehnt oder angenommen); die 15 vorhandenen Fälle
      bleiben mit ihrer Erwartung. Je Zusage ist die Mutation gesehen (`AGENTS.md`
      §3.10). Der Beleg des Implementers — je Zusage Zusage · Mutation · roter Fall,
      dazu eine Probe-Message mit Struktur-Kennung gegen den mit `make hooks-install`
      aktivierten Träger, Ausgabe und Exit — steht in §7 dieses Plans.
- [ ] Doku: Der Kopfkommentar des Trägers (und der Prüfung, wenn sie sich ändert)
      nennt die Ablehnung und ihre Grenzen nur so weit, wie die Gegenprobe prüft
      (`AGENTS.md` §3.11); `harness/README.md` nennt sie in §Sensors (Zeile
      `make hook-gegenprobe`, Bindung an [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)) und in §Traceability
      rules; `AGENTS.md` §5 Regel 1 nennt den Träger mit Herkunfts-Anker
      `seit slice-harness-commit-struktur-id`; `harness/mk/hook-gegenprobe.mk` nennt
      die Ablehnung im Hilfetext.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; `BEO-REPO/commit-nennt-struktur-kennung`
      wechselt in `state.md` von *geplant* auf *verkörpert*, mit Zielort und
      Herkunfts-Anker.
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
| `spec/spezifikation.md` §11, §12 | update (Architect, vor dem Code, geliefert) | `SPEC-050` mit der Ablehnung, den Entscheidungen zu den Randformen aus §6 und den Grenzen; Historien-Zeile |
| `docs/plan/adr/0035-…` und ADR-Index | neu (Architect, vor dem Code, geliefert) | [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) ergänzt [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) und [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) um die Ablehnung; Entscheidung und Gründe, `Schärft:` `SPEC-050` (`AGENTS.md` §3.8); Status Proposed bis zur Annahme durch den Nutzer |
| `.githooks/commit-msg` | update | Ablehnung nach `SPEC-050` **vor** dem frühen Exit 0 für benannte Slices und Wellen; Kopfkommentar |
| `tools/hook/commit-msg-gegenprobe.sh` | update | je Randform aus §6 ein Fall mit der Erwartung aus §11; die Tabelle im Kopf folgt; die 15 vorhandenen Fälle bleiben |
| `harness/mk/hook-gegenprobe.mk` | update | Hilfetext des Ziels nennt die Ablehnung |
| `harness/README.md`, `AGENTS.md` §5 | update | §Sensors und §Traceability rules; Regel 1 nennt ihren Träger |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint` liegt in `done/` (WIP-Limit 1).
Reihenfolge nach Entscheidung des Nutzers vom 2026-10-05, 2026-10-06 und 2026-10-07:
`slice-harness-lint-werkzeug`, die vier Umstellungs-Slices,
`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`, `slice-harness-lint`,
dieser Slice, `slice-harness-integration-wait`, `slice-harness-meldungskatalog-gate`, `slice-harness-abdeckung-gate`,
`slice-harness-coverage`, `slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`, `slice-harness-mutation`. Zwischen diesem Slice und
`slice-harness-integration-wait` besteht keine technische Abhängigkeit; dieser steht
vorn, weil jeder weitere Commit das Muster wiederholen kann. Vor dem ersten
Code-Commit entscheidet der Architect die Randformen aus §6 in §11 und ob eine
Folge-ADR nötig ist (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Entscheidungen in §11
  verlangen ein Lesen der Message über einen Zeilen-Ausdruck hinaus (etwa Backticks
  über Zeilengrenzen oder Code-Blöcke), und Träger, Gegenprobe und Doku sind nicht in
  einer Review-Sitzung prüfbar; dann trägt dieser Slice die Ablehnung im Fließtext,
  die Sonderformen gehen in einen eigenen Slice vor `slice-harness-abdeckung-gate`.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer nimmt [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) nicht
  an; dann zuerst diese Entscheidung. Die frühere Bedingung (Prüfung in der
  mitgelieferten Datei) ist mit `SPEC-050` Punkt 1 entfallen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit `hook-gegenprobe`, die Probe-Message mit
Struktur-Kennung vom aktivierten Träger abgelehnt (Beleg in §7), Closure-Notiz mit
Lerneintrag (neuer Sensor, Auslöser `BEO-REPO/commit-nennt-struktur-kennung`, 3×).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des Vertrags** (`AGENTS.md` §3.12) — entschieden vom Architect am
2026-10-08 vor dem ersten Code-Commit, in [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--struktur-kennungen-im-commit-träger-commit-msg); die Gründe in
[ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md). Was dort nicht steht, entscheidet der Implementer nicht, er gibt es
zurück.

- **Kommentarzeilen** — nicht gelesen, auch eingerückt: erstes Zeichen nach Leerzeichen
  und Tabs ist `#` (Punkt 2, derselbe Lesebereich wie die Annahme); ein anderes
  `core.commentChar` wird gelesen (Grenze).
- **Diff nach `git commit -v`** — unter der Scissors-Zeile nicht gelesen; nur die genaue
  Zeile mit je 24 Bindestrichen beendet den Bereich, eine andere Form nicht (Punkt 2).
- **Backticks, Pfad, Anker, Zitat** — Backticks, Zitat, Klammern, Linktext und Pfad
  stellen nicht frei, die Kennung zählt; Wortzeichen sind Buchstabe, Ziffer und
  Unterstrich, `/`, `#`, `-`, `.` grenzen ab. Ein Anker ist klein geschrieben und keine
  Nennung (Punkt 3). Abweichend von der Annahme nach [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md), die `-` als
  Wortzeichen führt: Eine Struktur-Kennung hat kein Suffix mit `-`, das sie zu einem
  anderen Namen machte.
- **Merge- und Revert-Messages** — ausgenommen, gleich ob von git erzeugt: Betreff
  beginnt mit `Merge ` oder `Revert ` (Punkt 4). Eine selbst geschriebene Message dieser
  Form bleibt ungeprüft (Grenze, akzeptiertes Negativ: im Bestand keine Merge- oder
  Revert-Message, fünf Betreffs mit Struktur-Kennung, deren Revert sonst scheiterte).
- **Groß- und Kleinschreibung, Ziffernzahl, Platzhalter** — nur `SPEC-NNN`/`ARC-NNN` in
  Großschreibung mit genau drei Ziffern; `spec-049`, `Spec-049`, `SPEC-49`,
  `SPEC-0491`, `SPEC-<NNN>`, `ARC-*` sind keine (Punkt 3). Ob sie existiert, ist ohne
  Belang; in einem Bereich zählen beide Enden.
- **Struktur-Kennung neben einer zugelassenen** — die Ablehnung geht vor, vor dem
  frühen Exit 0 des Trägers und vor der mitgelieferten Prüfung (Punkt 5). Ohne lesbare
  Datei bleibt der Weg wie heute.
- **Ausgabe und Ausgang** — je Nennung eine Zeile `commit-msg: Zeile <n>: Struktur-Kennung
  <Kennung>` auf stderr, Zeilennummer der Datei, je Zeile und Kennung einmal, danach
  eine feste Schlusszeile; stdout leer; Ausgang 1 wie bei fehlender Kennung (Punkt 6).

**Bestand** (`git log --format=%B`, Stand `25325be`, 489 Commits): Nach Punkt 3 nennen
27 Messages eine Struktur-Kennung, alle großgeschrieben mit drei Ziffern im Fließtext
oder Betreff, fünf davon im Betreff; keine kleine Form, kein Anker, kein Platzhalter,
keine Merge- oder Revert-Message. Die genaue Schreibweise fängt damit jede Form des
Bestands; die Commits bleiben, wie sie sind (§1).

**Offene Punkte für den Architect — entschieden:**

- **Ort der Prüfung** — der Träger `.githooks/commit-msg` (`SPEC-050` Punkt 1). Er liegt
  an einem Pfad, den der Bootstrap nur belegt, wenn er frei ist, und bleibt bei einem
  erneuten Lauf stehen (`harness/conventions.md`, Regelblock `grundlagen-traceability.md`;
  `harness/mk/hooks-install.mk`). Kein eigenes Skript: Der Ort gehört dem Repo schon.
- **Umfang der neuen Kennung** — nur die Ablehnung, ihr Lesebereich und ihr Vorrang;
  die Annahme bleibt in [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) und im Kopf des Trägers (§1,
  Bestand bleibt stehen). `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum`
  trifft diesen Slice nicht erneut: Seine neue Festlegung steht im Technik-Stratum.
- **Folge-ADR** — ja, [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md): Die Ablehnung ändert für Messages mit
  Struktur-Kennung, was die beiden Accepted ADRs zusagen (durchlassen, unverändert
  weiterreichen). Sie ergänzt beide und ersetzt keine; beide bleiben unverändert
  (ADR-Index §Beziehungen).

**Risiken:**

- **Falscher Treffer sperrt einen legitimen Commit** — etwa ein Pfad oder Anker, der
  eine Struktur-Kennung als Teil trägt; die Entscheidungen zu Backticks, Pfad, Anker und
  Schreibweise und ihre Fälle in der Gegenprobe fangen das. — **Ausgang:** — (bei Closure)
- **Annahme bricht beim Umbau** — die Ablehnung tritt neben die Annahme, nicht an ihre
  Stelle; die 15 vorhandenen Fälle bleiben (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`,
  1×). — **Ausgang:** — (bei Closure)
- **Regel greift nur im aktivierten Klon** — ohne `make hooks-install` oder mit
  `--no-verify` bleibt die Message ungeprüft, und das Review bleibt der zweite Leser;
  der Kopf des Trägers nennt das als Grenze. — **Ausgang:** — (bei Closure)

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

- **Belege zur DoD (Implementer):** <je Zusage Zusage · Mutation · roter Fall; die
  Probe-Message gegen den aktivierten Träger mit Ausgabe und Exit>
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
`docs/plan/planning/observations/BEO-REPO/` am Stand `cceb17f` gesichtet, dazu der
Nachtrag von F-336 im selben Commit wie dieser Plan (Zähler = Dateien unter
`evidence/`). Treffer:

- `BEO-REPO/commit-nennt-struktur-kennung` (3× mit dem Nachtrag, Ausgang *geplant* auf
  diesen Slice) — der Gegenstand.
- `BEO-REPO/werkzeug-commit-ohne-zugelassene-kennung` (2×, Stand *behoben*, Träger
  `.githooks/commit-msg`) — derselbe Träger; die Ablehnung darf die Annahme nicht
  brechen (Risiko in §6).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (13×, verkörpert in `AGENTS.md`
  §3.10) — die Ablehnung ist eine neue Zusage; daher die Mutation je Zusage in DoD-Punkt 2.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (16×, verkörpert in §3.11) —
  Kopfkommentar, Hilfetext und `harness/README.md` sagen nur zu, was die Gegenprobe
  prüft (DoD-Punkt 3).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (11×, verkörpert in §3.12) und
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (3×, verkörpert) — darum
  stehen die Randformen in §6 offen und werden vor dem Code entschieden.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — dieser Slice ist
  wellenlos; §4 *Start* nennt den Architect-Schritt vor dem Code.
- `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum` (1×) — die Regeln des
  Trägers stehen heute nur in den ADRs und im Skriptkopf; offener Punkt *Umfang der
  neuen Kennung* in §6. Entschieden: Die Annahme bleibt stehen, die neue Festlegung
  steht in `SPEC-050`; der Eintrag trifft diesen Slice nicht (§6).
- `BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen` (1×) — offener
  Punkt *Ort der Prüfung* in §6.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — ein Risiko in §6.
- `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht` (3×, verkörpert) — darum
  stehen die Belege der DoD in §7, nicht im Bericht.

Außer dem Gegenstand erreicht keiner der Einträge mit diesem Slice die Schwelle 3×;
keine neue Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
