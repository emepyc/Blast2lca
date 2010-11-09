#!/usr/bin/perl

use strict;
use warnings;

my $gi_taxid_file = shift @ARGV;
my $last_line = qx/tail -n 1 $gi_taxid_file/;
my ($last_val) = split /\t/,$last_line;

my $bin=0;
substr($bin,$_*4,4,pack ("N",0)) for (0..$last_val);

open my $dict_fh, "<", $gi_taxid_file or die $!;
while (<$dict_fh>){
  my ($key,$val) = split /\t/;
  substr($bin,$key*4,4,pack("N",$val));
}
close $dict_fh;

print $bin;
