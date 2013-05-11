
(defprotocol Dataset
  (all   [this] "return all possible positions on the dataset")
  (token [this at] "return the token at position")
  (walk  [this from] "return neighbor positions at pos"))


(defn- posify [row-index row-item]
  (map (fn [pos] [row-index pos]) (range (count row-item))))

(defrecord SequenceDatset [items]
  Dataset
  (all   [this]
         (mapcat posify (range) (:items this)))
  (token [this [row pos]] 
         (nth (nth (:items this) row) pos))
  (walk  [this [row pos]] 
         (let [row-item (nth (:items this) row)]
           (if (> (count row-item) pos)
             [row (inc pos)]
             []))))

(def example (SequenceDatset. ["ACGT" "CGATA" "CFASG" "FAEG"]))

(:items example)

(def example-db ["ACGT" "CGTA" "ACCGTGC" "CGTACG"])

(defn all-positions 
  "returns a list of all positions in a database"
  [db]
  (mapcat (fn [i row] 
            (map (fn [v] [i v]) 
                 (range (count row))))
          (range) db))

(defn simple-step
  "finds the following positions in a db "
  [db [row i]]
  (let [row-item (db row)]
    (if (> (count row-item) i)
      [(Walk. (nth row-item i) [row (inc i)])]
      [])))

(simple-step example-db [0 1])

(group-by :tok (mapcat #(step example-db %) (all-positions example-db)))

(defrecord Query [pattern positions])
(defrecord Walk [token to])

(defn empty-query [db] 
  (Query. [] (all-positions db))

(defn simple-extend [db query]
  (let [walks (mapcat #(simple-step db %) (:pos query))
        grouped (group-by :tok edges)]
    (map (fn [[tok walks]]
           (Query. (conj (:pattern query) tok)
                   (map :to walks))) grouped)))

(def e (empty-query example-db))

(def level-1 (simple-extend example-db e))
(def level-2 (simple-extend example-db (first level-1)))

(defn spexs [*in* *out* db])
