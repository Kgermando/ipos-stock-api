# Exemples de Requêtes API Dashboard

Ce fichier contient des exemples de requêtes pour tester l'API Dashboard.

## Configuration de Base

```bash
# Variables d'environnement pour les tests
ENTREPRISE_UUID="550e8400-e29b-41d4-a716-446655440000"
POS_UUID="6ba7b810-9dad-11d1-80b4-00c04fd430c8"
BASE_URL="http://localhost:8000/api/dashboard/main"
START_DATE="2023-11-01T00:00:00Z"
END_DATE="2023-11-13T23:59:59Z"
```

## 1. Statistiques Principales

```bash
curl -X GET "${BASE_URL}/stats?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 2. Graphique de Ventes (par jour)

```bash
curl -X GET "${BASE_URL}/sales-chart?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=2023-11-01T00:00:00Z&end_date=2023-11-13T23:59:59Z" \
  -H "Content-Type: application/json"
```

## 3. Graphique de Ventes (par heure - même jour)

```bash
curl -X GET "${BASE_URL}/sales-chart?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=2023-11-13T00:00:00Z&end_date=2023-11-13T23:59:59Z" \
  -H "Content-Type: application/json"
```

## 4. Graphique Donut des Plats

```bash
curl -X GET "${BASE_URL}/plat-chart?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 5. Graphique Donut des Produits

```bash
curl -X GET "${BASE_URL}/product-chart?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 6. Alertes de Stock

```bash
curl -X GET "${BASE_URL}/stock-alerts?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}" \
  -H "Content-Type: application/json"
```

## 7. Rotation de Stock

```bash
curl -X GET "${BASE_URL}/stock-rotation?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 8. Statistiques des Plats

```bash
curl -X GET "${BASE_URL}/plat-statistics?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 9. Statistiques des Livraisons

```bash
curl -X GET "${BASE_URL}/livraison-statistics?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 10. Zones de Livraison

```bash
curl -X GET "${BASE_URL}/livraison-zones?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 11. Performance des Livreurs

```bash
curl -X GET "${BASE_URL}/livreur-performance?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 12. Statistiques de la Caisse

```bash
curl -X GET "${BASE_URL}/caisse-statistics?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 13. Flux de Trésorerie (par jour)

```bash
curl -X GET "${BASE_URL}/flux-tresorerie?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=2023-11-01T00:00:00Z&end_date=2023-11-13T23:59:59Z" \
  -H "Content-Type: application/json"
```

## 14. Flux de Trésorerie (par heure)

```bash
curl -X GET "${BASE_URL}/flux-tresorerie?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=2023-11-13T00:00:00Z&end_date=2023-11-13T23:59:59Z" \
  -H "Content-Type: application/json"
```

## 15. Répartition des Transactions

```bash
curl -X GET "${BASE_URL}/repartition-transactions?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}" \
  -H "Content-Type: application/json"
```

## 16. Top Transactions

```bash
curl -X GET "${BASE_URL}/top-transactions?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}&limit=5" \
  -H "Content-Type: application/json"
```

## Tests avec HTTPie

Si vous préférez HTTPie, voici les mêmes requêtes :

```bash
# Statistiques principales
http GET localhost:8000/api/dashboard/main/stats \
  entreprise_uuid=="${ENTREPRISE_UUID}" \
  pos_uuid=="${POS_UUID}" \
  start_date=="${START_DATE}" \
  end_date=="${END_DATE}"

# Graphique de ventes
http GET localhost:8000/api/dashboard/main/sales-chart \
  entreprise_uuid=="${ENTREPRISE_UUID}" \
  pos_uuid=="${POS_UUID}" \
  start_date=="${START_DATE}" \
  end_date=="${END_DATE}"

# Alertes de stock
http GET localhost:8000/api/dashboard/main/stock-alerts \
  entreprise_uuid=="${ENTREPRISE_UUID}" \
  pos_uuid=="${POS_UUID}"
```

## Tests avec Postman

Collection Postman JSON :

```json
{
  "info": {
    "name": "Dashboard API Tests",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8000/api/dashboard/main"
    },
    {
      "key": "entreprise_uuid",
      "value": "550e8400-e29b-41d4-a716-446655440000"
    },
    {
      "key": "pos_uuid",
      "value": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
    },
    {
      "key": "start_date",
      "value": "2023-11-01T00:00:00Z"
    },
    {
      "key": "end_date",
      "value": "2023-11-13T23:59:59Z"
    }
  ],
  "item": [
    {
      "name": "Get Dashboard Stats",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{baseUrl}}/stats?entreprise_uuid={{entreprise_uuid}}&pos_uuid={{pos_uuid}}&start_date={{start_date}}&end_date={{end_date}}",
          "host": ["{{baseUrl}}"],
          "path": ["stats"],
          "query": [
            {"key": "entreprise_uuid", "value": "{{entreprise_uuid}}"},
            {"key": "pos_uuid", "value": "{{pos_uuid}}"},
            {"key": "start_date", "value": "{{start_date}}"},
            {"key": "end_date", "value": "{{end_date}}"}
          ]
        }
      }
    }
  ]
}
```

## Validation des Réponses

### Script de test automatique

```bash
#!/bin/bash

ENTREPRISE_UUID="550e8400-e29b-41d4-a716-446655440000"
POS_UUID="6ba7b810-9dad-11d1-80b4-00c04fd430c8"
BASE_URL="http://localhost:8000/api/dashboard/main"
START_DATE="2023-11-01T00:00:00Z"
END_DATE="2023-11-13T23:59:59Z"

echo "Testing Dashboard API endpoints..."

# Test 1: Statistiques principales
echo "1. Testing dashboard stats..."
response=$(curl -s "${BASE_URL}/stats?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}")
echo "Response: $response"

# Test 2: Alertes de stock
echo "2. Testing stock alerts..."
response=$(curl -s "${BASE_URL}/stock-alerts?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}")
echo "Response: $response"

# Test 3: Graphique de ventes
echo "3. Testing sales chart..."
response=$(curl -s "${BASE_URL}/sales-chart?entreprise_uuid=${ENTREPRISE_UUID}&pos_uuid=${POS_UUID}&start_date=${START_DATE}&end_date=${END_DATE}")
echo "Response: $response"

echo "Tests completed!"
```

## Réponses Attendues

### Exemple de réponse pour `/stats`

```json
{
  "totalArticles": 150,
  "articlesRuptureStock": 5,
  "articlesRuptureStockPercentage": 3,
  "totalVentes": 1250,
  "totalVentesPercentage": 85,
  "totalMontantVendu": 45000.00,
  "totalMontantVenduPercentage": 75
}
```

### Exemple de réponse pour `/stock-alerts`

```json
[
  {
    "uuid": "abc-123-def-456",
    "name": "Produit en rupture",
    "reference": "REF001",
    "unite_vente": "pièce",
    "stock": 0,
    "alertType": "rupture",
    "image": "",
    "prix_vente": 25.50
  }
]
```

### Exemple de réponse vide

```json
{
  "labels": [],
  "series": [],
  "percentages": []
}
```

## Codes d'Erreur

### Paramètres manquants (400)

```json
{
  "error": "Les paramètres entreprise_uuid et pos_uuid sont requis"
}
```

### Format de date invalide (400)

```json
{
  "error": "Format de date de début invalide"
}
```

### Dates requises (400)

```json
{
  "error": "Les dates de début et de fin sont requises"
}
```