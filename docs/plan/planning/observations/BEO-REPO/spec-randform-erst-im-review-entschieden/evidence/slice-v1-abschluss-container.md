**Vorgang:** slice-v1-abschluss-container
**Fund:** Review F-7 (LOW): `binaries.sh` baut in ein Stage-Verzeichnis außerhalb von `DIST`, verschiebt erst nach dem letzten Build und räumt es bei `EXIT`, `SIGINT` und `SIGTERM`; `SPEC-055` Punkt 3 und 9 und §6 des Plans nannten es nicht, der Implementer hatte es im Code entschieden. Der Architect trug es in Punkt 3 und 9 nach.
