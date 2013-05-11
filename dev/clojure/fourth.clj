;; Anything you type in here will be executed
;; immediately with the results shown on the
;; right.

(defn- group-map-by [g f coll]
  "Returns a map of the elements of coll keyed by the result of 
  g on each element. The value at each key will be a vector of the
  corresponding mapped elements by f, in the order they appeared in coll."
  (persistent!
   (reduce (fn [ret x]
             (let [k (g x)]
               (assoc! ret k (conj (get ret k []) (f x)))))
           (transient {}) coll)))

; ===================================
; Define structures for the algorithm
(defprotocol Dataset
  "This is interface for searching data in datasets"
  (all   [this] "return all possible positions on the dataset")
  (walk  [this pos] "return coll of steps from pos"))

(defrecord Query [pattern positions])
(defrecord Step [token position])

(defn empty-query [dataset]
  (Query. [] (all dataset)))

(defn- child-query [parent [token positions]]
  (Query. (conj (:pattern parent) token) positions))

; extender is a function that returns a dictionary token -> positions
; eg. {t1 [p1 p2 p3] t2 [p4 p5 p6]}
; where t_ are tokens, p_ are positions

(defn combine-extenders [extenders]
  (fn [dataset positions] 
    (apply merge-with concat (map #(% dataset positions) extenders))))

(defn walk-extend [dataset positions]
  (let [steps (mapcat #(walk dataset %) positions)]
   (group-map-by :token :position steps)))

; recursive spexs algorithm
(defn spexs-step [ds q extend]
  (map #(child-query q %) (extend ds (:positions q))))

(defn spexs [{
    ds  :dataset ; dataset
    in  :in      ; input coll
    out :out     ; output coll
    extend  :extend  ; position extender function
    extend? :extend? ; query filter for further extension
    output? :output? ; query filter for output
  }]
  (let [e (empty-query ds)]
    (loop [in (conj in e)
           out out]
      (if-not (empty? in)
        (let [[q & qs] in
              querys (spexs-step ds q extend)
              new-in  (concat qs  (filter extend? querys))
              new-out (concat out (filter output? querys))]
          (recur new-in new-out))
        out))))

; ===================================
; Example implementation of a Dataset

(defn- posify [row-index row-item]
  (map (fn [pos] [row-index pos]) (range (count row-item))))

(defrecord SequenceDataset [items]
  (token [this [row pos]] 
       (nth (nth (:items this) row) pos))
  Dataset
  (all   [this]
         (mapcat posify (range) (:items this)))
  (walk  [this [row i]] 
         (let [row-item (nth (:items this) row)]
           (if (> (count row-item) i)
             [(Step. (token this [row i]) [row (inc i)])]
             []))))

; example
(def example (SequenceDataset. ["ACGT" "CGATA" "AGCTTCGA" "GCGTAA"]))

(spexs {:dataset example :input [] :output [] 
        :extend walk-extend
        :extend? #(> (count (:positions %)) 3) 
        :output? #(> (count (:pattern %)) 2)})
