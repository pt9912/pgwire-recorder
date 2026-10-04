# Fremdimporte im Kern meldet das Architektur-Gate nur über die Tech-Liste

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

a-check meldet einen Import des Domain Models nur dann, wenn die importierte
Bibliothek in der `tech`-Liste von `.a-check.yml` steht. Ein Import der
Standardbibliothek wie `os` oder einer nicht gelisteten Bibliothek geht mit
0 Befunden durch; `spec/architecture.md` §2 verlangt ein Domain Model ohne
Dateisystem- und Drittbibliothekszugriff.
