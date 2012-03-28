open OUT, ">", "dna.gen";

for( $idx = 0; $idx < 4096; $idx++) {
	print OUT map { ("A","C","G","T")[rand 4] } 1..(8 + rand 3);
	print OUT "\n";
}

close OUT;