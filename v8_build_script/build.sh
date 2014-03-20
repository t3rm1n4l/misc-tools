#!/bin/bash -xe
# Description: Build libv8 on nix platforms

git clone https://github.com/v8/v8.git
cd v8
git checkout 3.21.17

git clone https://github.com/svn2github/gyp build/gyp
(cd build/gyp; git checkout a3e2a5caf24a1e0a45401e09ad131210bf16b852)

git clone https://github.com/svn2github/icu46 third_party/icu
(cd third_party/icu; git checkout 73170776491a3e38f68ab0367f3121256a3cc289)

cat <<EOF | git apply
diff --git a/build/standalone.gypi b/build/standalone.gypi
index f183331..cd078f9 100644
--- a/build/standalone.gypi
+++ b/build/standalone.gypi
@@ -258,7 +258,7 @@
           'GCC_INLINES_ARE_PRIVATE_EXTERN': 'YES',
           'GCC_SYMBOLS_PRIVATE_EXTERN': 'YES',      # -fvisibility=hidden
           'GCC_THREADSAFE_STATICS': 'NO',           # -fno-threadsafe-statics
-          'GCC_TREAT_WARNINGS_AS_ERRORS': 'YES',    # -Werror
+          'GCC_TREAT_WARNINGS_AS_ERRORS': 'NO',    # -Werror
           'GCC_VERSION': 'com.apple.compilers.llvmgcc42',
           'GCC_WARN_ABOUT_MISSING_NEWLINE': 'YES',  # -Wnewline-eof
           'GCC_WARN_NON_VIRTUAL_DESTRUCTOR': 'YES', # -Wnon-virtual-dtor
diff --git a/src/typing.cc b/src/typing.cc
index 34bb64b..661feb6 100644
--- a/src/typing.cc
+++ b/src/typing.cc
@@ -33,6 +33,7 @@
 namespace v8 {
 namespace internal {

+#pragma GCC diagnostic ignored "-Wuninitialized"

 AstTyper::AstTyper(CompilationInfo* info)
     : info_(info),
EOF

make -j8 library=shared i18nsupport=off native
