#!/bin/bash -eu
# jextract --dump-includes scripts/jextract.sh shared/libenry.h
#### Extracted from: shared/libenry.h

package=tech.sourced.enry.internal
output="$2/generated-sources/jextract"
lib="$1"

exec jextract \
    --include-struct GetLanguageByContent_return \
    --include-struct GetLanguageByEmacsModeline_return \
    --include-struct GetLanguageByExtension_return \
    --include-struct GetLanguageByFilename_return \
    --include-struct GetLanguageByModeline_return \
    --include-struct GetLanguageByShebang_return \
    --include-struct GetLanguageByVimModeline_return \
    --include-typedef GoInt \
    --include-typedef GoInt16 \
    --include-typedef GoInt32 \
    --include-typedef GoInt64 \
    --include-typedef GoInt8 \
    --include-typedef GoSlice \
    --include-typedef GoString \
    --include-typedef GoUint \
    --include-typedef GoUint16 \
    --include-typedef GoUint32 \
    --include-typedef GoUint64 \
    --include-typedef GoUint8 \
    --include-typedef _GoString_ \
    --include-function GetColor \
    --include-function GetLanguage \
    --include-function GetLanguageByContent \
    --include-function GetLanguageByEmacsModeline \
    --include-function GetLanguageByExtension \
    --include-function GetLanguageByFilename \
    --include-function GetLanguageByModeline \
    --include-function GetLanguageByShebang \
    --include-function GetLanguageByVimModeline \
    --include-function GetLanguageExtensions \
    --include-function GetLanguageType \
    --include-function GetLanguages \
    --include-function GetLanguagesByContent \
    --include-function GetLanguagesByEmacsModeline \
    --include-function GetLanguagesByExtension \
    --include-function GetLanguagesByFilename \
    --include-function GetLanguagesByModeline \
    --include-function GetLanguagesByShebang \
    --include-function GetLanguagesByVimModeline \
    --include-function GetMimeType \
    --include-function IsBinary \
    --include-function IsConfiguration \
    --include-function IsDocumentation \
    --include-function IsDotFile \
    --include-function IsGenerated \
    --include-function IsImage \
    --include-function IsTest \
    --include-function IsVendor \
    --header-class-name GoEnry \
    --output "$output" \
    --target-package "$package" \
    --library enry \
    "$lib/libenry.h"
