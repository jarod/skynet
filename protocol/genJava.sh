export SKYNET_JAVA=~/workspace/java/skynet-client
GEN_DIR=$SKYNET_JAVA/src/main/java
rm -rf $GEN_DIR/skynet/proto
protoc --java_out=$GEN_DIR skynet.proto