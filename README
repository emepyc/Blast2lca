Blast2lca 0.400 (Dec 2013) README file

TUTORIAL
========

Blast2lca is a simple and fast tool to help in the taxonomical classification of biological sequences based on BLAST searches.
It has been primarily used as part of computational pipelines in the analysis of the taxonomic content of metagenomic samples.
If you want to try blast2lca follow these steps:

1.- Installation:

1.1.- Install Go (programming language) from http://golang.org

1.2.- Set the $GOPATH environmental variable. For example:
   $ mkdir $HOME/gocode
   $ export GOPATH=$HOME/gocode
   [ Long story here: http://golang.org/cmd/go/#hdr-GOPATH_environment_variable ]

1.3.- Include the $GOPATH/bin directory in your $PATH (even if it doesn't exist yet)
   $ export PATH=$PATH:$GOPATH

1.4.- You can add the following lines to your ~/.profile to make these changes permanent.
     export GOPATH=$HOME/gocode
     export PATH=$PATH:$GOPATH

1.5.- Install Blast2lca with the following commands:
  $ go get github.com/emepyc/Blast2lca/blast2lca
  $ go get github.com/emepyc/Blast2lca/gitaxid2bin

1.6.- Make sure that your $GOPATH and $PATH variables are set correctly:
  $ which blast2lca
  $ which gitaxid2bin
  Both commands should give you the path to the just installed tools.

2.- Downloads:

2.1.- The NCBI's taxonomy database from ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz
  $ wget ftp://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump.tar.gz

2.2.- The NCBI's GI to taxonomy mapping from ftp://ftp.ncbi.nih.gov/pub/taxonomy/gi_taxid_[prot/nucl].dmp.gz
  $ wget wget ftp://ftp.ncbi.nlm.nih.gov/pub/taxonomy/gi_taxid_prot.dmp.gz

2.3.- Uncompress & untar taxdump.tar.gz to get the names.dmp and nodes.dmp files.
  $ gunzip gi_taxid_prot.dmp.gz # or gi_taxid_nucl.dmp.gz
  $ tar -xzvf taxdump.tar.gz

3.- Format the database files: (resulting file is gi_taxid_prot.bin)
  $ gitaxid2bin ncbi_taxonomy/gi_taxid_prot.dmp # or ncbi_taxonomy/gi_taxid_nucl.dmp

4.- Run the program:
  $ blast2lca -savemem -dict gi_taxid_prot.bin -nodes nodes.dmp -names names.dmp sample/metagenome.blout.127 > metagenome.lca # or gi_taxid_nucl.bin

Note: You should always use the correct gi_taxid_[prot|nucl].dmp file that matches with your blast results. ie. if your blast has been run against the nr database you should use the gi_taxid_nucl.dmp file, otherwise blast2lca won't be able to map any of your GIs in the dict file.


INTRODUCTION
============

Blast2lca is a tool designed to help in the taxonomical assignment of biological sequences through BLAST (http://www.ncbi.nlm.nih.gov/pmc/articles/PMC146917/?tool=pubmed) based homology searches. For this assignment it calculates the lowest common ancestor (LCA: http://en.wikipedia.org/wiki/Lowest_common_ancestor) of the best scoring homologs. It works in a similar way as the MEGAN tool (see below for details).

Its input is the output of the BLAST program (tab formatted with "-m8" option in blastall or "-outfmt 6" in blast+ programs). The "subject" field must contain its GI identifier in the form "gi|xxxxxxxxxx|" where "x"s are the actual identifier (this is likely to happen if you use NCBI databases like nr or nt).

You will also need the Taxonomy database from the NCBI (the taxdump.[gz|tar.gz|tar.Z] and the gi_taxid_[nucl|prot].[zip|dmp.gz] files. These files can be downloaded from ftp://ftp.ncbi.nih.gov/pub/taxonomy/. See HOW TO RUN THE PROGRAM below.


A note on MEGAN:
-----------------
MEGAN (http://www-ab.informatik.uni-tuebingen.de/software/megan) is a nice tool for global analysis and taxonomic representation of metagenomic sequences but while working with it we noticed a few issues that prevent its use for our specific needs:

- Although written in Java, MEGAN is not a very fast tool, or at least, not fast enough when you are analyzing many >10e5 sequences. Working with 454-Titanium runs we noticed long hangs and crashes.

- MEGAN focuses in the global analysis of the input sequences. This means that it needs to track information of all the sequences being analyzed to show the final trees, etc. But many times we prefer to get the LCA of each individual sequence to further analyze them using different statistical methods. Unfortunately, when exporting the LCA of each individual sequence MEGAN doesn't inform you about the taxonomic rank or level the LCA is, making it hard i) to know how specific the LCA is and ii) to postprocess this LCA information (for example, group by "family" level).

- It works with the full (pairwise) output of BLAST instead of the more compacted tabular form. This means that if you are working with many sequences i) BLAST will be slower and ii) Your input to MEGAN will be orders of magnitud bigger.

- It is not obvious how to update the taxonomy information that MEGAN uses. The NCBI's Taxonomy database is dynamic and interpreting a recent BLAST result with an outdated taxonomy leads to a greater number of unassigned hits and sequences.

Some of these concerns were commented with the MEGAN development team (Dr. Huson's lab) and we are greatful for the inclusion of some of them in the MEGAN package (for example, now MEGAN is able to deal with a modified BLAST tabular format that must include an extra column with the Taxid of the subject GIs, for us this is far better than having to run all the BLAST with full/pairwise output). But finally we decided to build a tool that tackles these specific issues.

Blast2lca is:

+ Fast: It is multicore. Has been used in a 24-nodes server obtaining > 20x accelerations. We are however trying to improve the speed of each individual thread.

+ Simple: Its input is a tabular blast result (-m8) and the NCBI's taxonomy database. Its output is the sequence IDs and their LCA. There is only one aditional step: you have to reformat the gi_taxid[nucl|prot].[zip|dmp.gz] file that serves the GI to taxid mappings (available at the NCBI site (ftp://ftp.ncbi.nih.gov/pub/taxonomy/) to speed up things.

+ Transparent: Since you use the NCBI's Taxonomy database directly you can be always up-to-date with the current releases.

Of course, Blast2lca doesn't give you all the other goodies that MEGAN offers you (fancy tree visualization, etc...).



COMMON USAGE:
=============

1.- Taxonomy DB preparation:
----------------------------
The first time you use this software (and in general everytime you want to update the Taxonomy database) you should prepare the Taxonomy database for its usage with Blast2lca [See above for details]

2.- Run the program:
--------------------
$ ./blast2lca [options] <blast>
Where options are:

      --version:
              Prints the current version and exits

      --help:
              Prints the usage of the software and exits

      --nprocs:
              Number of CPUs to use (defaults to 4)

      --savemem:
              The data structures used by the program are stored in the hard disk
              instead of in memory. For now, this option is recommended.

      --nodes:
              Path to the nodes.dmp file downloaded from the NCBI's Taxonomy DB
              Defaults to "nodes.dmp"

      --names:
              Path to the names.dmp file downloaded from the NCBI's Taxonomy DB
              Defaults to "names.dmp"

      --dict:
              Path to the gi2taxid binary file you have obtained from the previous step

      --levels:
             The taxonomic levels you want from the LCA.
             If the LCA of a sequence is lower than the specified level, you will get this instead.
             If the LCA of a sequence is higher than the specified level, you will get the true LCA but with the letters "uc_" preppended ("uc" stands for unclassified).
	     You can pass a series of taxonomy levels separated by ":", see the example below

Example:
$ ./blast2lca -names names.dmp -nodes nodes.dmp -dict gi_taxid_prot.bin -levels=superkingdom:phylum:class:family blastm8.txt > lca.txt

3.- Output:
----------
For each input query sequence you will obtain its header and its taxons according to the LCA (separated by a tab).
Example:
GCQ6XTU01A3NBL	Bacteria;Firmicutes;Bacilli;Streptococcaceae
GCQ6XTU01A3N9G	Bacteria;Cyanobacteria;uc_Cyanobacteria;uc_Chroococcales
GCQ6XTU01A3N88	Bacteria;Proteobacteria;Gammaproteobacteria;Pasteurellaceae
GCQ6XTU01A3N8G	Bacteria;Proteobacteria;Deltaproteobacteria;Myxococcaceae
GCQ6XTU01A3NC4	Bacteria;Firmicutes;Bacilli;Streptococcaceae

As a fist global visualization of the result you can run the script makeTree.pl (under the tools/ directory) to generate a tree in HTML (+CSS +JavaScript) that can be visualized with any web browser. See the tools/makeTree.README file for (a bit) more information.

You can also convert the output of blast2lca in a format compatible with MEGAN using the script located in tools/to_megan.pl.


BUGS & CONTACT:
===============
For bug reports, feature requests or comments, please send an email to emepyc@gmail.com



