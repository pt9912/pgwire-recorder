# Slice slice-v1-abschluss-tls-client: TLS zum Client bei record und replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-23`](../../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client), [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0018](../../adr/0018-tls-zum-client.md)

**Berührte Spec-Stellen:** `LH-FA-23.a` · `LH-FA-05.c` · `LH-FA-05.e` · `SPEC-018` · `SPEC-019` · `SPEC-020` · `SPEC-034` · `SPEC-033` · `ARC-006` · `ARC-005`

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

**Ziel:** `record` und `replay` nehmen auf Wunsch TLS von Clients an (`--tls-cert`, `--tls-key`), weisen unverschlüsselte Clients dann ab, sofern sie nicht zugelassen sind (`--allow-plaintext`), und verhalten sich ohne Konfiguration wie bisher.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- TLS zum Upstream im Record-Modus — Out-of-Scope von LH-FA-23; die Verbindung zum Server bleibt unverschlüsselt.
- Prüfung von Client-Zertifikaten und Zertifikatsverwaltung — Out-of-Scope von LH-FA-23.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-23`](../../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client): Ein Client verbindet sich verschlüsselt, `record` zeichnet auf und `replay` liefert die Aufzeichnung aus, jeweils wie bei einer unverschlüsselten Verbindung; dieselbe Aufzeichnung ist über beide Verbindungsarten gleich (Abnahmeszenario 15; Integrationstest).
- [ ] Ohne Konfiguration wird `SSLRequest` mit `N` beantwortet; mit Konfiguration wird ein unverschlüsselter Client (auch nach `GSSENCRequest`) mit `PGR-E6003` abgewiesen, mit `--allow-plaintext` zugelassen, ein Klartext-`CancelRequest` bleibt `PGR-W3001`; mehrere Verbindungen nacheinander funktionieren; `--tls-key` allein und `--allow-plaintext` ohne `--tls-cert` sind `PGR-E2001` (Test).
- [ ] Nicht lesbare, ungültige, abgelaufene, verschlüsselte oder nicht passende Dateien beenden den Start mit `PGR-E2007` und Exit-Code 2, ein gescheiterter Handshake endet mit `PGR-W3002` ohne Wirkung auf den Exit-Code; Schlüssel erscheinen weder in Logs noch in `config show` (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/pgwire` | update | TLS-Terminierung, Aushandlung nach `SSLRequest`, Abweisung unverschlüsselter Clients |
| `internal/adapters/driving/cli` | update | Optionen `--tls-cert`, `--tls-key`, `--allow-plaintext`, Fehlerabbildung auf `PGR-E2001` und `PGR-E2007` |
| `test/integration` | update | Happy/Boundary/Negative nach LH-FA-23 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`, `slice-v1-abschluss-sessions` ist `done`, und [ADR-0018](../../adr/0018-tls-zum-client.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Handshake und Abweisung der Klartext-Verbindungen sprengen den Slice — zurück zur Zerlegung.
- `in-progress` → `open`: Die TLS-Bibliothek verlangt eine Abhängigkeit außerhalb des PGWire-Adapters — Carveout.

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

- Test-Zertifikate müssen im Testlauf erzeugt werden; ein fest eingecheckter Schlüssel wäre ein Geheimnis im Repository — **Ausgang:** offen bis Closure.
- Clients mit eigener Zertifikatsprüfung scheitern an einem selbst signierten Zertifikat; das Handbuch muss den Weg beschreiben — **Ausgang:** offen bis Closure.

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

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
