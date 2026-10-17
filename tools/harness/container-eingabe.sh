# container-eingabe — Prüfung der Eingabe von `binaries` und `test-container`
# (SPEC-055 Punkt 1 und 2). Wird von tools/harness/binaries.sh und
# tools/harness/run-container.sh gelesen; der Aufrufer setzt `name` und `fehler=0`.
# Die Marken `# P1…` und `# P2…` am Zeilenende benennen die Form für die Gegenprobe.
BIN_OS="linux darwin windows"
BIN_ARCH="amd64 arm64"
ERF_BINARIES="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64"
ERF_IMAGES="linux/amd64 linux/arm64"

eingabe() { echo "$name: Eingabe: $*" >&2; fehler=1; }

# liste_pruefen <Bezeichnung> <erforderliche Einträge> <Eintrag>…
liste_pruefen() {
  local was="$1" erf="$2" e os arch gesehen=" " p
  shift 2
  [ "$#" -gt 0 ] || eingabe "die Liste $was ist leer" # P1a
  for e; do
    [[ "$e" == */* ]] || { eingabe "Eintrag ohne '/' in der Liste $was: $e"; continue; } # P1b1
    os="${e%%/*}"
    arch="${e#*/}"
    [ -n "$os" ] || eingabe "leeres Betriebssystem im Eintrag $e der Liste $was" # P1b2
    [ -n "$arch" ] || eingabe "leere Architektur im Eintrag $e der Liste $was" # P1b3
    [[ -z "$os" || " $BIN_OS " == *" $os "* ]] || eingabe "Betriebssystem '$os' unbekannt im Eintrag $e der Liste $was" # P1c1
    [[ -z "$arch" || " $BIN_ARCH " == *" $arch "* ]] || eingabe "Architektur '$arch' unbekannt im Eintrag $e der Liste $was" # P1c2
    [[ "$gesehen" != *" $e "* ]] || eingabe "Eintrag doppelt in der Liste $was: $e" # P1d
    gesehen+="$e "
  done
  for p in $erf; do
    [[ "$gesehen" == *" $p "* ]] || eingabe "$p fehlt in der Liste $was" # P1e
  done
}

# images_in_binaries <Liste der Binaries als ein String> <Image-Plattform>…
images_in_binaries() {
  local bin=" $1 " e
  shift
  for e; do
    [[ "$bin" == *" $e "* ]] || eingabe "Image-Plattform $e steht nicht in der Liste der Binaries" # P1f
  done
}

version_pruefen() {
  local erst="0-9A-Za-z" rest="0-9A-Za-z.+_-" # P2zeichen
  [[ "${VERSION-}" =~ ^[$erst][$rest]*$ ]] || eingabe "VERSION '${VERSION-}' ist ungültig" # P2
}
