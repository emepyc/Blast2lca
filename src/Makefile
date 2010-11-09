# Copyright 2010 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# prerequisite: GOROOT and GOARCH must be defined

# defines $(GC) (compiler), $(LD) (linker) and $(O) (architecture)
include $(GOROOT)/src/Make.inc

# name of the app
TARG=blast2lca

# build "main" executable
blast2lca: tax blast
	$(GC) blast2lca.go
	$(LD) -o $@ blast2lca.$O
	@-echo "Done. Executable is: $@"

clean:
	@-rm -rf *.[$(OS)o] *.a [$(OS)].out 

tax: taxonomy.go
	$(GC) $<

blast: blastm8.go
	$(GC) $<



