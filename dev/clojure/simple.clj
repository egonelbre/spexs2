(ns spexs
	(:require clojure.set))

(def input [\A \C \G \G \G \T \A \C nil \C \G \G \G \A \T nil])
(def full-set (set (range 0 (count text))))

; generates first layer
(defn identify [i] 
	{ (str (nth input i)) #{(inc i)}})

; merges { \A #{ 1 2 3 } } 
(defn merge-idents [x y]
	(merge-with union x y))

; increments each position and groups by the position
(defn simplex [positions] 
	(let [
		looks		(map identify positions)
		extended	(reduce merge-idents looks)
		result		(dissoc extended (str nil))
	] result))

; generates prefixer function for map
(defn prefixer [prefix]
	(fn [[suffix positions]] 
		{ (str prefix suffix) positions } ))

(defn spexs [[letter positions]]
	(let [
		extensions 		(simplex positions)
		sub-extensions	(reduce merge (map spexs extensions))
		fully-named     (reduce merge (map (prefixer letter) sub-extensions))
		result			(merge fully-named extensions)
	] result))

(defn run-spexs [positions]
	(let [result (spexs [nil positions])
		counts (reduce merge (map (fn [[x y]] {x (count y)} ) result))
	] counts))
