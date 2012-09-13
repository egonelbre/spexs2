(ns spexs
	(:require clojure.set))

(def input [\A \C \G \G \G \T \A \C nil \C \G \G \G \A \T nil])
(def full-set (set (range 0 (count input))))

(defn simple-stepper [input]
	(fn [i]
		(let [ sym (nth input i) ]
			(case sym
				nil {}
				{ sym #{(inc i)} })
	)))

(defn star-stepper [input]
	; create map from i -> closest letters
	; reduce by { let -> i + 1 }, with merge max
	(fn [i]
		; return
		{}))

; merges map of sets by using union on the sets
(defn merge-steps [x y]
	(merge-with clojure.set/union x y))

(defn merge-steppers [a b]
	(fn [i] (merge-steps (a i) (b i))))

; extends positions to a map of
;   sym => next positions
(defn extend [step positions]
	(reduce merge-steps (map step positions)))

;==========================================================================

(defn make-querys [parent querys]
	)

(defn expand-query [query]
		(make-querys query (extend simple-step (:pos query))))

(defn expand-group [query] 0)

(defn expand [query]
	(case (query :type)
		:query (expand-query query)
		:group (expand-group query)
	))

(defn spexs-step [query expand combine]
	(let [
		stepped (expand query)
		groups (combine stepped)
		all (merge stepped groups)
		result {
			:input	(filter expandable? all)
			:output	(filter outputtable? all)
		}
	] result))

(defn spexs [query]
	(let [
		extensions (extend simple-step (:pos query)) ; { sym #{ poss } )
		sub-querys (make-querys extensions)
		group-querys (make-groups sub-querys)
	] result))


(defn extend-query [query]
	(query-extension (extend simple-step (:pos query))))

(defn extend-group [group]
	(query-group (let [
			extensions (map force-extend (:links group))
		] result)))

