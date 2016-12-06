VER=0.11-161206
OS="linux"
ARCH="amd64"

dir="`pwd`"
mkdir build

goxc -pv=$VER -wd="$dir/cmd/skynet-matrix" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="xc archive"

goxc -pv=$VER -wd="$dir/cmd/skynet-agent" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="xc archive"

