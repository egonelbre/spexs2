;; Anything you type in here will be executed
;; immediately with the results shown on the
;; right.

;; Anything you type in here will be executed
;; immediately with the results shown on the
;; right.

(require '[ clojure.core.reducers :as r ])

; a parallel grouping function
(defn group-map-by [g f coll]
  (r/fold 
   (r/monoid (partial merge-with into) (constantly {}))
   (fn [ret x]
     (let [k (g x)]
       (assoc ret k (conj (get ret k []) (f x)))))
   coll))

; these are the minimal requirements for a dataset
(defprotocol Dataset
  "This is interface for searching data in datasets"
  (all   [this] "return all possible positions on the dataset")
  (walk  [this pos] "return coll of steps from pos"))

(defrecord Query [pattern positions])
(defrecord Step [token position])

; create an empty query for a dataset
(defn- empty-query [dataset]
  (Query. [] (all dataset)))

; create a child query for parent given a token and positions
(defn- child-query [parent [token positions]]
  (Query. (conj (:pattern parent) token) positions))

; simple extension function
(defn walk-extend [dataset positions]
  (let [steps (mapcat #(walk dataset %) positions)]
    (group-map-by :token :position steps)))

; group extender
(defn- select-merged [m ks]
  (mapcat second (select-keys m ks)))

(defn extend-grouper [extend groups]
  (fn [dataset positions]
    (let [ extended (extend dataset)
           groupings (for [[token items] groups] 
                        [token (select-merged extended items)])]
      (apply merge extended groupings))))

; function to combine multiple extension functions
(defn combine-extenders [extenders]
  (fn [dataset positions] 
    (apply merge-with concat (map #(% dataset positions) extenders))))

; spexs algorithm
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


; here is an example how to implement a dataset
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

(spexs {:dataset example :in [] :out [] 
        :extend walk-extend
        :extend? #(> (count (:positions %)) 3) 
        :output? #(> (count (:pattern %)) 2)})
