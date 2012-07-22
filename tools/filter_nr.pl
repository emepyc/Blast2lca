#!/usr/bin/perl

use strict;
use warnings;
use Bio::LITE::Taxonomy::NCBI;

my ($downdir,$nr,$term) = @ARGV;

my $nodes = "$downdir/nodes.dmp";
my $names = "$downdir/names.dmp";
my $dict  = "$downdir/gi_taxid_prot.dmp.bin";

my $taxidMapper = Bio::LITE::Taxonomy::NCBI::Gi2taxid->new(dict => "dict.in");

my $taxDB = Bio::LITE::Taxonomy::NCBI->new(
    nodes => $nodes,
    names => $names,
    dict  => $dict,
    save_mem => 1
    );

my $level = $taxDB->get_level_from_name($term);

my $nrFH;
# Accept nr in compressed or uncompressed file
if ( ($nr =~ /\.gz$/) or ($nr =~ /\.Z$/) ){
  open $nrFH, "zcat $nr |" or die $!;
} else {
  open $nrFH, "<", $nr or die $!;
}

my $seq = 1;
{
  local $/="\n>";

  while (<$nrFH>){
    chomp;
    s/^>// if ($seq == 1);
    print STDERR "Processing sequence $seq          \r";
    $seq++;
    /^\S+\|(\d+)\|/;
    my $id = $1;
    my $t = $taxDB->get_term_at_level_from_gi($id,$level);
    print ">$_\n" if (($t eq $term) and (defined $id));
  }
}
close $nrFH;
