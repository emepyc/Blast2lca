#!/usr/bin/perl

## TODO:
## Use env perl instead of hardcode the perl binary
## Allow input of both nucl and prot gi_taxid files
## Allow gi_taxid input files in compressed format

use strict;
use warnings;

my $gi_taxid_file = shift @ARGV;
my $last_line = qx/tail -n 1 $gi_taxid_file/;
my ($last_val) = split /\t/,$last_line;

my $bin= "\0" x (4 * ($last_val+1));

open my $dict_fh, "<", $gi_taxid_file or die $!;
while (<$dict_fh>){
  my ($key,$val) = split /\t/;
  substr($bin,$key*4,4,pack("N",$val));
}
close $dict_fh;

print $bin;
