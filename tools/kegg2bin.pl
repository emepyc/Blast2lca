#!/usr/bin/perl

# TO-DO: Implement this in Go with 2 goroutines,
# one entering KEGG pathways and the other entering Taxonomy IDs
# each one their proper slot

use strict;
use warnings;
use File::Tail;

my ($gi_taxid_file, $gi_kegg_file) = @ARGV;
my $ttail = File::Tail->new(name=>$gi_taxid_file,tail=>1) or croak "$!";
my $last_tax_line = $ttail->read();
die "I can't read the $gi_taxid_file, is it empty?\n" unless (defined $last_tax_line);

my ($last_tax_gi) = split /\t/, $last_tax_line;

my $bin = "\0" x (4*18 * ($last_val_gi+1));

open my $dict_fh, "<", $gi_kegg_file or die $!;
while (<$dict_fh>) {
    chomp;
    my ($key,$val) = split /\t/;
    substr(
}
