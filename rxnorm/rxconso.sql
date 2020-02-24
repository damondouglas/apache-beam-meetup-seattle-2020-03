SELECT
    DISTINCT(LOWER(str)) as name,
    rxcui
FROM `bigquery-public-data.nlm_rxnorm.rxnconso_*`
WHERE sab = "RXNORM"