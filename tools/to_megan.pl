#!/usr/bin/env perl

use strict;
use warnings;

my %bytaxon;
while (<>) {
    chomp;
    my @flds = split;
    $bytaxon{$flds[1]}++;
}

for my $taxon (keys %bytaxon) {
    print "$taxon\t$bytaxon{$taxon}\n";
}

