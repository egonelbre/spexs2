package utils

/*
from HyperQuick algorithm for discrete hypergeometric distribution

function InvJm(n,x,N,m):double;
	InvJm:=(1-x/(m+1))/(1-(n-1-x)/(N-1-m))
end function

begin
# loop fixed number of steps
	s:=1.0;
	for k:=x to M-2 do
		s:=s*InvJm(n,x,N,k)+1.0;
	end for
	# loop variable number of steps according to accuracy e
	ak:=s;
	bk:=s;
	k:=M-2;
	epsk:=2*e;
	while (k<N-(n-x)-1) and (epsk>e) do
		ck:=ak/bk;
		k:=k+1;
		jjm:=InvJm(n,x,N,k);
		bk:=bk*jjm+1.0;
		ak:=ak*jjm;
		epsk:=(N-(n-x)-1-k)*(ck-ak/bk);
	end while
	result:=1-(ak/bk-epsk/2);
end
*/


func invJm(n, x, N, m int) float64 {
	a := 1.0 - float64(x)/float64(m+1)
	b := 1.0 - float64(n-1-x)/float64(N-1-m)
	return a/b
}

const e = 1e-40

/*
	n = number of uniformly chosen balls
	x = at most how many are black
	N = number of balls
	M = black balls
*/

func hyperQuick(n, x, N, M int) float64 {
	s := float64(1.0)
	for k := x; k <= M-2; k += 1 {
		s = s * invJm(n,x,N,k) + 1.0
	}

	ak := s
	bk := s
	epsk := 2*e

	for k := M-2; (k < N-(n-x)) && (epsk > e); k += 1 {
		ck := ak/bk
		jjm := invJm(n,x,N,k)
		bk := bk*jjm + 1.0
		ak := ak*jjm
		epsk = float64(N-(n-x)-1-k)*(ck-ak/bk)
	}
	return 1.0 - (ak/bk - epsk/2)
}

func HyperQuickSplit(o int, r int, O int, R int) float64 {
	// o-1, O, R, o+r
	// 
	//return hyperQuick(r+o, o, O+R, R)
	return hyperQuick(o+r, r, O+R, O)
}
