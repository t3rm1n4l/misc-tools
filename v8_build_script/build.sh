#!/bin/bash -xe
# Description: Build libv8 on nix platforms

git clone https://github.com/v8/v8.git
cd v8
git checkout 3.21.17

git clone https://github.com/svn2github/gyp build/gyp
(cd build/gyp; git checkout a3e2a5caf24a1e0a45401e09ad131210bf16b852)

git clone https://github.com/svn2github/icu46 third_party/icu
(cd third_party/icu; git checkout 73170776491a3e38f68ab0367f3121256a3cc289)

cat <<"EOF" | git apply
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
diff --git a/src/ia32/lithium-codegen-ia32.cc b/src/ia32/lithium-codegen-ia32.cc
index 06628ff..a6935a4 100644
--- a/src/ia32/lithium-codegen-ia32.cc
+++ b/src/ia32/lithium-codegen-ia32.cc
@@ -444,6 +444,7 @@ bool LCodeGen::GenerateJumpTable() {
   return !is_aborted();
 }
 
+#pragma GCC diagnostic ignored "-Wuninitialized"
 
 bool LCodeGen::GenerateDeferredCode() {
   ASSERT(is_generating());
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
diff --git a/tools/js2c.py b/tools/js2c.py
index 9492b00..c15f635 100644
--- a/tools/js2c.py
+++ b/tools/js2c.py
@@ -1,4 +1,4 @@
-#!/usr/bin/env python
+#!/usr/bin/env python2.6
 #
 # Copyright 2012 the V8 project authors. All rights reserved.
 # Redistribution and use in source and binary forms, with or without
EOF

cat <<"EOF" | (cd build/gyp; git apply)
diff --git a/gyp b/gyp
index b53a6dd..b298550 100755
--- a/gyp
+++ b/gyp
@@ -5,4 +5,4 @@
 
 set -e
 base=$(dirname "$0")
-exec python "${base}/gyp_main.py" "$@"
+exec python2.6 "${base}/gyp_main.py" "$@"
diff --git a/pylib/gyp/generator/make.py b/pylib/gyp/generator/make.py
index b88a433..5c7310b 100644
--- a/pylib/gyp/generator/make.py
+++ b/pylib/gyp/generator/make.py
@@ -135,7 +135,7 @@ quiet_cmd_alink = AR($(TOOLSET)) $@
 cmd_alink = rm -f $@ && $(AR.$(TOOLSET)) crs $@ $(filter %.o,$^)
 
 quiet_cmd_alink_thin = AR($(TOOLSET)) $@
-cmd_alink_thin = rm -f $@ && $(AR.$(TOOLSET)) crsT $@ $(filter %.o,$^)
+cmd_alink_thin = rm -f $@ && $(AR.$(TOOLSET)) crs $@ $(filter %.o,$^)
 
 # Due to circular dependencies between libraries :(, we wrap the
 # special "figure out circular dependencies" flags around the entire
EOF

make -j8 library=shared i18nsupport=off native
