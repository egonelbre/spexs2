;; Anything you type in here will be executed
;; immediately with the results shown on the
;; right.

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
      [(make-edge (nth row-item i) [row (inc i)])]
      [])))

(simple-step example-db [0 0])

(group-by :tok (mapcat #(simple-step example-db %) (all-positions example-db)))

(defn make-query [pattern positions]
  {:pat pattern :pos positions})

(defn make-edge [token position]
  {:tok token :pos position})

(defn empty-query [db] (make-query "" (all-positions db))

(defn simple-extend [db query]
  (let [stepped (mapcat #(simple-step db %) (:pos query))
        grouped (group-by :tok stepped)]
    (map (fn [[tok els]]
           (make-query (str (:pat query) tok)
                       (map :pos els))) grouped)))

(def e (empty-query example-db))

(def level-1 (simple-extend example-db e))
(def level-2 (simple-extend example-db (first level-1)))

(defn spexs [*in* *out* db])
