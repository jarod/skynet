VER=0.8-140702
OS="linux"
ARCH="amd64"

dir="`pwd`"
mkdir build

goxc -pv=$VER -wd="$dir/skynet-matrix" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="xc archive"

goxc -pv=$VER -wd="$dir/skynet-agent" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="xc archive"

