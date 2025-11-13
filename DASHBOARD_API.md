# Documentation API Dashboard

## Vue d'ensemble

Cette API fournit une suite complète d'endpoints pour alimenter le dashboard du frontend Angular. L'API est conçue pour reproduire exactement les fonctionnalités des services TypeScript du frontend.

## Endpoints Disponibles

### Base URL: `/api/dashboard/main`

Tous les endpoints suivants nécessitent les paramètres de requête :
- `entreprise_uuid` (requis) : UUID de l'entreprise
- `pos_uuid` (requis) : UUID du point de vente
- `start_date` (optionnel) : Date de début au format ISO 8601 (2006-01-02T15:04:05Z07:00)
- `end_date` (optionnel) : Date de fin au format ISO 8601 (2006-01-02T15:04:05Z07:00)

---

## 1. Statistiques Principales

### `GET /api/dashboard/main/stats`

Retourne les statistiques principales du dashboard.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`

**Réponse :**
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

---

## 2. Données Graphique de Ventes

### `GET /api/dashboard/main/sales-chart`

Retourne les données pour le graphique de courbe de ventes.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "dates": ["08h", "09h", "10h", ...], // ou ["01/11", "02/11", ...]
  "totalCommandes": [12, 15, 8, ...],
  "montantVendu": [1500.00, 2250.00, 1200.00, ...],
  "gainObtenu": [300.00, 450.00, 240.00, ...]
}
```

**Note :** Si `start_date` et `end_date` correspondent au même jour, les données sont groupées par heure (00h-23h). Sinon, elles sont groupées par jour.

---

## 3. Graphique Donut des Plats

### `GET /api/dashboard/main/plat-chart`

Retourne les données pour le graphique donut des plats les plus vendus.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "labels": ["Pizza Margherita", "Burger Classic", "Salade César", ...],
  "series": [15000.00, 12500.00, 8900.00, ...],
  "percentages": [35.5, 29.2, 20.8, ...]
}
```

---

## 4. Graphique Donut des Produits

### `GET /api/dashboard/main/product-chart`

Retourne les données pour le graphique donut des produits les plus vendus.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "labels": ["Coca-Cola", "Eau Minérale", "Jus d'Orange", ...],
  "series": [8500.00, 6200.00, 4100.00, ...],
  "percentages": [42.5, 31.0, 20.5, ...]
}
```

---

## 5. Alertes de Stock

### `GET /api/dashboard/main/stock-alerts`

Retourne les produits en alerte de stock (rupture et avertissement).

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Réponse :**
```json
[
  {
    "uuid": "abc123",
    "name": "Produit A",
    "reference": "REF001",
    "unite_vente": "pièce",
    "stock": 0,
    "alertType": "rupture",
    "image": "path/to/image.jpg",
    "prix_vente": 25.50
  },
  {
    "uuid": "def456",
    "name": "Produit B",
    "reference": "REF002",
    "unite_vente": "kg",
    "stock": 3,
    "alertType": "avertissement",
    "image": "",
    "prix_vente": 15.00
  }
]
```

---

## 6. Données de Rotation de Stock

### `GET /api/dashboard/main/stock-rotation`

Retourne les données de taux de rotation de stock.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "productNames": ["Produit A", "Produit B", "Produit C", ...],
  "rotationRates": [12.5, 8.7, 6.2, ...],
  "categories": ["Excellente (≥6)", "Très bonne (4-6)", "Bonne (2-4)", ...],
  "colors": ["#28a745", "#6f42c1", "#007bff", ...]
}
```

---

## 7. Statistiques des Plats

### `GET /api/dashboard/main/plat-statistics`

Retourne les statistiques détaillées des plats.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "totalPlats": 25,
  "totalClients": 1580,
  "quantitesVendues": 2450,
  "chiffresAffaires": 89500.00
}
```

---

## 8. Statistiques des Livraisons

### `GET /api/dashboard/main/livraison-statistics`

Retourne les statistiques des livraisons.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "totalLivraisons": 250,
  "enCours": 15,
  "effectuees": 220,
  "annulees": 15,
  "enCoursPercentage": 6.0,
  "effectueesPercentage": 88.0,
  "annuleesPercentage": 6.0,
  "totalRevenu": 1250.00,
  "revenuMoyen": 5.00,
  "tauxReussite": 88.0
}
```

---

## 9. Données des Zones de Livraison

### `GET /api/dashboard/main/livraison-zones`

Retourne le top 5 des zones de livraison.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`

**Réponse :**
```json
[
  {
    "zoneName": "Centre-ville",
    "nombreLivraisons": 85,
    "revenu": 425.00
  },
  {
    "zoneName": "Quartier Nord",
    "nombreLivraisons": 65,
    "revenu": 325.00
  }
]
```

---

## 10. Performance des Livreurs

### `GET /api/dashboard/main/livreur-performance`

Retourne le top 5 des livreurs les plus performants.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`

**Réponse :**
```json
[
  {
    "uuid": "liv123",
    "name": "Jean Dupont",
    "totalLivraisons": 45,
    "effectuees": 42,
    "enCours": 2,
    "annulees": 1,
    "tauxReussite": 93.3
  }
]
```

---

## 11. Statistiques de la Caisse

### `GET /api/dashboard/main/caisse-statistics`

Retourne les statistiques complètes de la caisse.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "soldeCaisse": 15750.00,
  "totalEntrees": 25000.00,
  "totalSorties": 9250.00,
  "nombreTransactions": 150,
  "moyenneEntree": 312.50,
  "moyenneSortie": 154.17,
  "ratioEntreeSortie": 2.70,
  "tauxLiquidite": 63.0,
  "evolutionJournaliere": 0.0,
  "evolutionPercentage": 0.0,
  "tendance": "stable",
  "jourLePlusActif": "Vendredi",
  "heureLaPlusActive": "14h00",
  "nombreTransactionsParJour": 12.5
}
```

---

## 12. Flux de Trésorerie

### `GET /api/dashboard/main/flux-tresorerie`

Retourne les données de flux de trésorerie.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "dates": ["0h", "1h", "2h", ...], // ou ["01/11", "02/11", ...]
  "entrees": [500.00, 0.00, 0.00, ...],
  "sorties": [0.00, 150.00, 0.00, ...],
  "soldes": [500.00, 350.00, 350.00, ...]
}
```

---

## 13. Répartition des Transactions

### `GET /api/dashboard/main/repartition-transactions`

Retourne la répartition des transactions par catégorie.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`

**Réponse :**
```json
{
  "labels": ["Ventes", "Achats", "Charges", "Autres"],
  "values": [15000.00, 8500.00, 2000.00, 500.00],
  "percentages": [57.7, 32.7, 7.7, 1.9],
  "colors": ["#28a745", "#007bff", "#ffc107", "#dc3545"]
}
```

---

## 14. Top Transactions

### `GET /api/dashboard/main/top-transactions`

Retourne les meilleures entrées et sorties.

**Paramètres requis :**
- `entreprise_uuid`
- `pos_uuid`

**Paramètres optionnels :**
- `start_date`
- `end_date`
- `limit` (défaut: 5)

**Réponse :**
```json
{
  "topEntrees": [
    {
      "libelle": "Vente produits",
      "montant": 2500.00,
      "type": "Entree",
      "date": "2023-11-13T14:30:00Z",
      "reference": "REF123456"
    }
  ],
  "topSorties": [
    {
      "libelle": "Achat marchandises",
      "montant": 1200.00,
      "type": "Sortie",
      "date": "2023-11-13T10:15:00Z",
      "reference": "REF789012"
    }
  ]
}
```

---

## Exemples d'Usage

### Frontend Angular Service

```typescript
// dashboard-statistics.service.ts
async calculateDashboardStats(
  entreprise_uuid: string,
  pos_uuid: string,
  startDate?: Date,
  endDate?: Date
): Promise<DashboardStats> {
  const params = new URLSearchParams({
    entreprise_uuid,
    pos_uuid
  });
  
  if (startDate) params.append('start_date', startDate.toISOString());
  if (endDate) params.append('end_date', endDate.toISOString());
  
  const response = await fetch(`/api/dashboard/main/stats?${params}`);
  return response.json();
}
```

### Gestion des Erreurs

Tous les endpoints retournent des erreurs standardisées :

```json
{
  "error": "Les paramètres entreprise_uuid et pos_uuid sont requis"
}
```

Status codes :
- `200` : Succès
- `400` : Paramètres manquants ou invalides
- `500` : Erreur serveur

---

## Notes Techniques

1. **Dates** : Toutes les dates sont au format ISO 8601 UTC
2. **Montants** : Tous les montants sont arrondis à 2 décimales
3. **Pourcentages** : Arrondis à 2 décimales
4. **Performance** : Les requêtes utilisent des index PostgreSQL pour optimiser les performances
5. **Cache** : Aucun cache n'est implémenté - les données sont calculées à la demande

---

## Compatibilité

Cette API est entièrement compatible avec le frontend Angular existant et reproduit fidèlement le comportement des services TypeScript originaux.