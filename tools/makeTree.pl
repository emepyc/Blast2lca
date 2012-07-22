#!/usr/bin/perl

use strict;
use warnings;
use Data::Dumper;
use Bio::LITE::Taxonomy::NCBI;

use CGI;

my $css = <<EOCSS;
.TC0 { font-size : 0.6em; font-weight : 90;  }
.TC1 { font-size : 0.7em; font-weight : 100; }
.TC2 { font-size : 0.8em; font-weight : 200; }
.TC3 { font-size : 0.9em; font-weight : 300; }
.TC4 { font-size : 1.0em; font-weight : 400; }
.TC5 { font-size : 1.2em; font-weight : 500; }
.TC6 { font-size : 1.4em; font-weight : 600; }
.TC7 { font-size : 1.6em; font-weight : 700; }
.TC8 { font-size : 1.8em; font-weight : 800; }
.TC9 { font-size : 1.9em; font-weight : 900; }
.TC10 { font-size : 2.0em; font-weight : 900; }

.CTC0 { color : #739cf8 }
.CTC1 { color : #5a84e3 }
.CTC2 { color : #4b75d4 }
.CTC3 { color : #3563cb }
.CTC4 { color : #1f4db5 }
.CTC5 { color : #1341a9 }
.CTC6 { color : #06359f }
.CTC7 { color : #022982 }
.CTC8 { color : #011d5c }
.CTC9 { color : #01133d }
.CTC10 { color : #000000 }

.nTC0 { font-size : 1.2em; font-weight : 500; }

.hide {
      DISPLAY: none;
}

.show {
      DISPLAY: block;
}

.tree {
	background-color : lightyellow;
	border : 1px solid grey;
}

EOCSS

my $jscript = <<EOSCRIPT;
function activate_font_size () {
    var mysheet=document.styleSheets[0];
    var myrules=mysheet.cssRules? mysheet.cssRules: mysheet.rules;
    for (i=0; i<myrules.length; i++){
	if(myrules[i].selectorText == ".TC0") {
	    myrules[i].style.fontSize = "0.6em";
	    myrules[i].style.fontWeight = 90;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC1") {
	    myrules[i].style.fontSize = "0.7em";
	    myrules[i].style.fontWeight = 100;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC2") {
	    myrules[i].style.fontSize = "0.8em";
	    myrules[i].style.fontWeight = 200;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC3") {
	    myrules[i].style.fontSize = "0.8em";
	    myrules[i].style.fontWeight = 300;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC4") {
	    myrules[i].style.fontSize = "1.0em";
	    myrules[i].style.fontWeight = 400;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC5") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC6") {
	    myrules[i].style.fontSize = "1.4em";
	    myrules[i].style.fontWeight = 600;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC7") {
	    myrules[i].style.fontSize = "1.6em";
	    myrules[i].style.fontWeight = 700;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC8") {
	    myrules[i].style.fontSize = "1.8em";
	    myrules[i].style.fontWeight = 800;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC9") {
	    myrules[i].style.fontSize = "1.9em";
	    myrules[i].style.fontWeight = 900;
	    myrules[i].style.color = "#000000";
	} else if (myrules[i].selectorText == ".TC10") {
	    myrules[i].style.fontSize = "2.0em";
	    myrules[i].style.fontWeight = 900;
	    myrules[i].style.color = "#000000";
	}
    }
}
function activate_font_color () {
    var mysheet=document.styleSheets[0];
    var myrules=mysheet.cssRules? mysheet.cssRules: mysheet.rules;
    for (i=0; i<myrules.length; i++){
	if(myrules[i].selectorText == ".TC0") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#739cf8";
	} else if (myrules[i].selectorText == ".TC1") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#5a84e3";
	} else if (myrules[i].selectorText == ".TC2") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#4b75d4";
	} else if (myrules[i].selectorText == ".TC3") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#3563cb";
	} else if (myrules[i].selectorText == ".TC4") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#1f4db5";
	} else if (myrules[i].selectorText == ".TC5") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#1341a9";
	} else if (myrules[i].selectorText == ".TC6") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#06359f";
	} else if (myrules[i].selectorText == ".TC7") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#022982";
	} else if (myrules[i].selectorText == ".TC8") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#011d5c";
	} else if (myrules[i].selectorText == ".TC9") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#01133d";
	} else if (myrules[i].selectorText == ".TC10") {
	    myrules[i].style.fontSize = "1.2em";
	    myrules[i].style.fontWeight = 500;
	    myrules[i].style.color = "#000000";
	}
    }
}

function simplehideshow (name) {
    var whichcontent = document.getElementById(name);
    var pat1 = /(hide|show)/;
    var OK = pat1.exec(whichcontent.className);
    if (!OK){
	window.alert("Bad section!");
    }
    if (RegExp.\$1 == "show") {
	whichcontent.className = "hide";
	document.getElementById("b"+name).innerHTML = "+";
    } else {
	whichcontent.className = "show";
	document.getElementById("b"+name).innerHTML = "-";
    }
}

function hideshow (contentid) {
  var whichcontent = document.getElementById(contentid);
  var pat1 = /(.*)(hide|show)/;
  var OK = pat1.exec(whichcontent.className);
  if (!OK){
    window.alert ("Bad section!");
  }
  var sect = RegExp.\$1;
  var pat = /show/;
  var showed = pat.exec(whichcontent.className);
  if (showed) {
      whichcontent.className = sect+"hide";
   }
   else {
      whichcontent.className= sect+"show";
   }
}

EOSCRIPT
    
print STDERR "Loading TaxonomyDB... ";
my $taxDB = Bio::LITE::Taxonomy::NCBI->new(
    names => '/Users/miguel/repos/src/Blast2lca/taxonomyDB/names.dmp',
    nodes => '/Users/miguel/repos/src/Blast2lca/taxonomyDB/nodes.dmp',
    );
print STDERR "Ok\n";
my ($fin) = @ARGV;

open my $fh, "<", $fin or die $!;

my @taxIDs = map {chomp; (split /\t/)[1]} (<$fh>);
print STDERR Dumper \@taxIDs;

my @taxes = map { $taxDB->get_taxonomy($taxDB->get_taxid_from_name($_)) } @taxIDs;

print STDERR Dumper \@taxes;

my $q = new CGI;
my $tree_ref = tree_of_lists (\@taxes);

print $q->start_html (
    -style => {'code' => $css},
    -script => [
	 {
	     -type => "text/javascript",
	     -code => $jscript
	 }
    ]
    );

my $tot = 0;
for my $k (keys %$tree_ref) {
    $tot += $tree_ref->{$k}->{num};
}
my $treeHtml = "<a name = 'tree'>";
rec_putTax ($tree_ref, 0, \$treeHtml,$tot);
$treeHtml .= "</div></a>";
my $outtable = "\n<div class=tree><p><b>Phylogenetic tree</b> of the reads taxonomy. [number] indicates the number of times the given taxon appears in the results. <b>Relative abundance</b> of each taxon within each taxonomic level is graphically denoted by its font size or color:</p><p id=FontEffect>Font effect:<a href='javascript:desactivate_font_effect()'><span class='globbed' onmouseover=\"Tip('Desactivate font effect on tree', WIDTH, 250, TITLE, 'Phylogenetic Tree', SHADOW, true, FADEIN, 300, FADEOUT, 300, STICKY, 1, CLOSEBTN, true, CLICKCLOSE, true, BGCOLOR, '#FFFFFF', TITLEBGCOLOR, '#8eb1d2', TITLEFONTCOLOR, '#153E7E')\" onmouseout='UnTip()'>None</span></a> | <a href='javascript:activate_font_size()'><span class='globbed' onmouseover=\"Tip('Font size on each taxon denotes its relative abundance within each taxonomic level', WIDTH, 250, TITLE, 'Phylogenetic Tree', SHADOW, true, FADEIN, 300, FADEOUT, 300, STICKY, 1, CLOSEBTN, true, CLICKCLOSE, true, BGCOLOR, '#FFFFFF', TITLEBGCOLOR, '#8eb1d2', TITLEFONTCOLOR, '#153E7E')\" onmouseout='UnTip()'>Size</span></a> | <a href='javascript:activate_font_color()'><span class='globbed' onmouseover=\"Tip('Color font denotes taxon relative abundance within each taxonomic level', WIDTH, 250, TITLE, 'Phylogenetic Tree', SHADOW, true, FADEIN, 300, FADEOUT, 300, STICKY, 1, CLOSEBTN, true, CLICKCLOSE, true, BGCOLOR, '#FFFFFF', TITLEBGCOLOR, '#8eb1d2', TITLEFONTCOLOR, '#153E7E')\" onmouseout='UnTip()'>Color</span></a></p>\n$treeHtml</div>\n";
print $outtable;

sub rec_putTax
{
    my ($tree, $level, $retStr, $rtot) = @_;
    my ($maxlevelkey) = sort {$b->{num} <=> $a->{num}} values %$tree;
    my $max;
    if (defined $maxlevelkey) {
	$max = $maxlevelkey->{num};
    }
    my $sum = 0;
    for (values %$tree){
	$sum += $_->{num};
    }
    my $mean;
    if (values %$tree){
	$mean = int ($sum / (scalar (values %$tree)));
    } else {$mean = 0}
    for my $node (sort keys %$tree) {
	my $n = $tree->{$node}->{num};
	# Nueva forma de calcular la media:
	# calculamos la media y le asignamos valor 5.
	# de ahí => subimos y bajamos en rango de 0 <=> media <=> max
#	my $class = int (10 * $n / $max);  # old
	my $class;
	if ($n == $mean)   { $class = 5 }
	elsif ($n > $mean) { $class = int (($n * 5 / $max) + 5 ) }
	else               { $class = int (5 - ($n * 5 / $max)) }
#	print STDERR "SUM: $sum , MEAN: $mean , MAX: $max , N: $n => $class\n";
#	print STDERR "$node => TOT: $rtot , N: $n => $class2\n";
	if (scalar keys %{$tree->{$node}->{next}} == 0){
	    $$retStr .= ("__"x$level)." <span class=TC$class>".$node."</span><span class=\"sm\">[".$tree->{$node}->{'num'}."]</span><br />\n";
	} else {
	    $$retStr .= "__"x$level."<a href=\'javascript: simplehideshow (\"$node\")\'><b id=\"b$node\">+</b></a>";
	    $$retStr .= "<span class=TC$class>$node</span><span class=\"sm\">[".$tree->{$node}->{'num'}."]</span><br />\n<div id=\"$node\" class=\"hide\">\n";
	    rec_putTax ($tree->{$node}->{next},($level+1), $retStr, $rtot);
	    $$retStr .= "</div>\n";
	}
    }
    return
}

sub tree_of_lists
{
    my ($lsts) = @_;
    my @lists_ref = map { [split ";",$_] } @$lsts;
    my %tree;
    for my $list (@lists_ref){
	my $node = \%tree;
	
	for my $e (@$list){
	    if (defined $node->{$e}){
		$node->{$e}->{num}++;
	    } else {
		$node->{$e} = {'num' => 1, 'next' => {}};
	    }
	    $node = $node->{$e}->{'next'};
	}
    }
    return \%tree;
}



# sub rec_putTax
# {
#     my ($tree, $level, $retStr) = @_;
#     print STDERR "$$retStr\n\n";
#     for my $node (keys %$tree) {
# 	if (scalar keys %{$tree->{$node}} == 0){
# 	    $$retStr .= ("--"x$level).$node."<br />\n";
# 	} else {
# 	    $$retStr .= "--"x$level."<a href=\'javascript: simplehideshow (\"$node\")\'>-</a>";
# 	    $$retStr .= "$node\n<div id=\"$node\" class=\"show\">\n";
# 	    rec_putTax ($tree->{$node},($level+1), $retStr);
# 	    $$retStr .= "</div><br />\n";
# 	}
#     }
#     return $retStr;
# }

# sub rec_putTax
# {
#     my ($terms,$nextItem,$level) = @_;
#     my @ks = keys %$terms;
#     my $currTerm = shift @ks;
#     if (scalar @$terms == 0) {
# 	print +("."x$level).$currTerm."<br />\n";
#     } else {
# 	print "."x$level."<a href=\'javascript: simplehideshow (\"$currTerm\")\'>-</a>";
# 	print "$currTerm\n<div id=\"$currTerm\" class=\"show\">\n";
# 	rec_putTax (\@$terms,++$level);
# 	print "</div>\n";
#     }
# }


# sub tree_of_lists
# {
#     my ($taxes) = @_;
#     my @lists_ref = map { [split ";",$_] } @$taxes;
#     my %tree;
#     for my $list ( @lists_ref ) {
# 	my $node = \%tree;
	
# 	$node = $node->{shift @$list} ||= {} while @$list;
#     }
    
#     return \%tree;
# }


