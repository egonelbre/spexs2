
FLAGS=-stats -procs=2 -verbose=true
SPX=../spexs2 $(FLAGS)

GOFLAGS=-gcflags -B

rebuild:
	cd .. && go build $(GOFLAGS) .

dna: rebuild
	$(SPX) -conf=dna/small.json inp=dna/dna.small

events: rebuild
	$(SPX) -conf=events/conf.json inp=events/errors.txt ref=events/events.txt

proteins.10k: rebuild
	$(SPX) -conf=proteins/conf.json inp=proteins/g21_10k.inp ref=proteins/g27_10k.ref	

proteins.30k: rebuild
	$(SPX) -conf=proteins/conf.json inp=proteins/g21_30k.inp ref=proteins/g27_30k.ref	

proteins.large: rebuild
	$(SPX) -conf=proteins/conf.json inp=proteins/g21_large.inp ref=proteins/g27_large.ref	

text: rebuild
	$(SPX) -conf=text/conf.json inp=text/text.inp ref=text/text.ref

transcripts: rebuild
	$(SPX) -conf=transcript/conf.json inp=transcript/transcripts.inp 

yeast: rebuild
	$(SPX) -conf=yeast/conf.json inp=yeast/Yeast_-600_+2_W_cluster_1599945.fa ref=yeast/Yeast_-600_+2_W_random_1000_all.fa 
