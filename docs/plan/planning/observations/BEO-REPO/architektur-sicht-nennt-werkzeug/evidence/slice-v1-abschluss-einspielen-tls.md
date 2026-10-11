**Vorgang:** slice-v1-abschluss-einspielen-tls
**Fund:** Die Sicht (`spec/architecture.md` §6) nannte `go list -deps ./cmd/...` und das Linken des Binaries (Review F-601, HIGH). Berichtigt vom Architect in `3a2f1f3`: Die Sicht sagt nur noch „darf nur Testcode importieren; das ausgelieferte Programm enthält sie nicht“, den Aufruf trägt der Test.
