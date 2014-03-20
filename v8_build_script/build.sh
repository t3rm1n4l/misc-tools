#!/bin/bash -xe
# Description: Build libv8 on nix platforms

git clone https://github.com/v8/v8.git
cd v8
git checkout 3.25.18

git clone https://github.com/svn2github/gyp build/gyp
(cd build/gyp; git checkout a3e2a5caf24a1e0a45401e09ad131210bf16b852)

git clone https://github.com/svn2github/icu46 third_party/icu
(cd third_party/icu; git checkout 73170776491a3e38f68ab0367f3121256a3cc289)

cat <<EOF | git apply
diff --git a/src/typing.cc b/src/typing.cc
index b62e909..263bf55 100644
--- a/src/typing.cc
+++ b/src/typing.cc
@@ -35,6 +35,7 @@
 namespace v8 {
 namespace internal {

+#pragma GCC diagnostic ignored "-Wuninitialized"

 AstTyper::AstTyper(CompilationInfo* info)
     : info_(info),
EOF

make -j8 library=shared 18nsupport=off native
