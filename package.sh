VER=0.6-140623
OS="linux"
ARCH="386 amd64"

dir="`pwd`"
mkdir build

goxc -pv=$VER -wd="$dir/skynet-matrix" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="xc archive"

goxc -pv=$VER -wd="$dir/skynet-agent" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="xc archive"

