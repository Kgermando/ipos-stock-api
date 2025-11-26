# Documentation API Dashboard - Statistiques Caisse

## Vue d'ensemble

Cette API fournit des endpoints pour les statistiques de la caisse, conformément au service TypeScript/Angular du frontend.

---

## 📊 Endpoints Caisse

### 1. Statistiques Principales de la Caisse

**Endpoint:** `GET /api/dashboard/main/caisse-statistics`

**Paramètres Query:**
- `entreprise_uuid` (requis): UUID de l'entreprise
- `pos_uuid` (requis): UUID du point de vente
- `start_date` (optionnel): Date de début (format ISO 8601)
- `end_date` (optionnel): Date de fin (format ISO 8601)

**Réponse:**
```json
{
  "soldeCaisse": 15000.50,
  "totalEntrees": 20000.00,
  "totalSorties": 5000.00,
  "montantDebut": 500.50,
  "nombreTransactions": 150,
  "moyenneEntree": 250.00,
  "moyenneSortie": 100.00,
  "ratioEntreeSortie": 4.0,
  "tauxLiquidite": 75.00,
  "evolutionJournaliere": 500.00,
  "evolutionPercentage": 3.45,
  "tendance": "hausse",
  "jourLePlusActif": "Mercredi",
  "heureLaPlusActive": "14h00",
  "nombreTransactionsParJour": 25.5
}
```

**Description des champs:**
- `soldeCaisse`: Solde actuel (totalEntrees + montantDebut - totalSorties)
- `totalEntrees`: Somme de toutes les entrées
- `totalSorties`: Somme de toutes les sorties
- `montantDebut`: Montant d'ouverture/dépôt initial
- `nombreTransactions`: Nombre total de transactions
- `moyenneEntree`: Montant moyen des entrées
- `moyenneSortie`: Montant moyen des sorties
- `ratioEntreeSortie`: Ratio entrées/sorties
- `tauxLiquidite`: (soldeCaisse / totalEntrees) * 100
- `evolutionJournaliere`: Différence avec la période précédente
- `evolutionPercentage`: Pourcentage d'évolution
- `tendance`: 'hausse' | 'baisse' | 'stable'
- `jourLePlusActif`: Jour avec le plus de transactions
- `heureLaPlusActive`: Heure avec le plus de transactions
- `nombreTransactionsParJour`: Moyenne de transactions par jour

---

### 2. Flux de Trésorerie

**Endpoint:** `GET /api/dashboard/main/flux-tresorerie`

**Paramètres Query:**
- `entreprise_uuid` (requis): UUID de l'entreprise
- `pos_uuid` (requis): UUID du point de vente
- `start_date` (requis): Date de début (format ISO 8601)
- `end_date` (requis): Date de fin (format ISO 8601)

**Réponse:**
```json
{
  "dates": ["01/12", "02/12", "03/12"],
  "entrees": [1500.00, 2000.00, 1800.00],
  "sorties": [800.00, 900.00, 750.00],
  "soldes": [700.00, 1800.00, 2850.00]
}
```

**Comportement:**
- **Un seul jour** (startDate == endDate): Affichage par **heure** (0h-23h)
- **Plusieurs jours**: Affichage par **jour** (format DD/MM)
- Le solde est cumulé et inclut le MontantDebut

**Exemple pour une journée:**
```json
{
  "dates": ["0h", "1h", "2h", ..., "23h"],
  "entrees": [0, 0, 0, ..., 500.00],
  "sorties": [0, 0, 0, ..., 100.00],
  "soldes": [500.50, 500.50, 500.50, ..., 900.50]
}
```

---

### 3. Répartition des Transactions

**Endpoint:** `GET /api/dashboard/main/repartition-transactions`

**Paramètres Query:**
- `entreprise_uuid` (requis): UUID de l'entreprise
- `pos_uuid` (requis): UUID du point de vente
- `start_date` (optionnel): Date de début
- `end_date` (optionnel): Date de fin

**Réponse:**
```json
{
  "labels": ["Entrées", "Sorties", "Montant Initial"],
  "values": [20000.00, 5000.00, 500.50],
  "percentages": [78.43, 19.60, 1.97],
  "colors": ["#28a745", "#dc3545", "#007bff"]
}
```

**Codes couleur:**
- `#28a745`: Vert (Entrées)
- `#dc3545`: Rouge (Sorties)
- `#007bff`: Bleu (Montant Initial)

---

### 4. Top Transactions

**Endpoint:** `GET /api/dashboard/main/top-transactions`

**Paramètres Query:**
- `entreprise_uuid` (requis): UUID de l'entreprise
- `pos_uuid` (requis): UUID du point de vente
- `start_date` (optionnel): Date de début
- `end_date` (optionnel): Date de fin
- `limit` (optionnel, défaut=5): Nombre de transactions à retourner

**Réponse:**
```json
{
  "topEntrees": [
    {
      "libelle": "Vente produits",
      "montant": 5000.00,
      "type": "Entree",
      "date": "2024-11-25T14:30:00Z",
      "reference": "REF-001"
    }
  ],
  "topSorties": [
    {
      "libelle": "Achat fournitures",
      "montant": 1500.00,
      "type": "Sortie",
      "date": "2024-11-25T10:15:00Z",
      "reference": "REF-002"
    }
  ]
}
```

---

### 5. Analyse des Catégories

**Endpoint:** `GET /api/dashboard/main/analyse-categories`

**Paramètres Query:**
- `entreprise_uuid` (requis): UUID de l'entreprise
- `pos_uuid` (requis): UUID du point de vente
- `start_date` (optionnel): Date de début
- `end_date` (optionnel): Date de fin

**Réponse:**
```json
[
  {
    "categorie": "Vente produits",
    "totalMontant": 15000.00,
    "nombreTransactions": 45,
    "pourcentage": 65.20,
    "moyenne": 333.33,
    "tendance": "stable"
  },
  {
    "categorie": "Achat fournitures",
    "totalMontant": 8000.00,
    "nombreTransactions": 20,
    "pourcentage": 34.80,
    "moyenne": 400.00,
    "tendance": "stable"
  }
]
```

**Description:**
- Groupement par `libelle` de transaction
- Tri décroissant par `totalMontant`
- `tendance`: Pour le moment retourne 'stable' (implémentation future: comparaison avec période précédente)

---

### 6. Prévisions de Trésorerie

**Endpoint:** `GET /api/dashboard/main/previsions-tresorerie`

**Paramètres Query:**
- `entreprise_uuid` (requis): UUID de l'entreprise
- `pos_uuid` (requis): UUID du point de vente
- `nombre_jours` (optionnel, défaut=7): Nombre de jours à prévoir

**Réponse:**
```json
[
  {
    "date": "26/11",
    "previsionEntree": 1800.00,
    "previsionSortie": 850.00,
    "previsionSolde": 16950.00,
    "confiance": 95
  },
  {
    "date": "27/11",
    "previsionEntree": 1750.00,
    "previsionSortie": 900.00,
    "previsionSolde": 17800.00,
    "confiance": 90
  }
]
```

**Algorithme:**
- Basé sur les moyennes des 30 derniers jours
- Variation aléatoire de ±20% pour simuler la réalité
- Niveau de confiance décroît de 5% par jour (max 95%, min 50%)
- Le solde est cumulatif

---

## 🔄 Format des dates

Toutes les dates doivent être au format **ISO 8601**:
```
2024-11-25T00:00:00Z
2024-11-25T23:59:59Z
```

---

## ⚠️ Gestion des erreurs

### Erreur 400 - Paramètres manquants
```json
{
  "error": "Les paramètres entreprise_uuid et pos_uuid sont requis"
}
```

### Erreur 400 - Format de date invalide
```json
{
  "error": "Format de date de début invalide"
}
```

---

## 📝 Notes d'implémentation

### Calcul du Solde de Caisse
```
soldeCaisse = totalEntrees + montantDebut - totalSorties
```

### Calcul de l'Évolution
- Compare la période actuelle avec la période précédente (même durée)
- Tendance:
  - `"stable"`: variation < 5%
  - `"hausse"`: variation ≥ 5% positive
  - `"baisse"`: variation < -5%

### Types de transactions
- `"Entree"`: Entrée d'argent
- `"Sortie"`: Sortie d'argent
- `"MontantDebut"`: Montant d'ouverture/dépôt initial

---

## 🎯 Exemples d'utilisation

### Obtenir les stats du jour actuel
```bash
GET /api/dashboard/main/caisse-statistics?entreprise_uuid=abc123&pos_uuid=def456&start_date=2024-11-25T00:00:00Z&end_date=2024-11-25T23:59:59Z
```

### Obtenir le flux sur 7 jours
```bash
GET /api/dashboard/main/flux-tresorerie?entreprise_uuid=abc123&pos_uuid=def456&start_date=2024-11-18T00:00:00Z&end_date=2024-11-25T23:59:59Z
```

### Obtenir le flux horaire d'une journée
```bash
GET /api/dashboard/main/flux-tresorerie?entreprise_uuid=abc123&pos_uuid=def456&start_date=2024-11-25T00:00:00Z&end_date=2024-11-25T23:59:59Z
```

### Obtenir les prévisions sur 14 jours
```bash
GET /api/dashboard/main/previsions-tresorerie?entreprise_uuid=abc123&pos_uuid=def456&nombre_jours=14
```

---

## 🏗️ Architecture

### Fichiers modifiés/créés:

1. **controllers/dashboard/dashboard-main.go**
   - `getCaisseStatistics()`: Calcul des statistiques principales
   - `getFluxTresorerieData()`: Génération des données de flux
   - `getRepartitionTransactionsData()`: Répartition par type
   - `getTopTransactions()`: Top entrées/sorties
   - `analyseCategories()`: Analyse par libellé
   - `genererPrevisions()`: Prévisions basées sur l'historique
   - `calculateCaisseEvolution()`: Calcul de l'évolution
   - `getJourLePlusActif()`: Analyse temporelle
   - `getHeureLaPlusActive()`: Analyse temporelle
   - `getFluxParHeure()`: Flux horaire
   - `getFluxParJour()`: Flux journalier

2. **models/dashboard.go**
   - Structures de données pour toutes les réponses

3. **routes/routes.go**
   - Enregistrement des endpoints

---

## ✅ Conformité avec le Frontend Angular

Cette implémentation Go est **100% compatible** avec le service TypeScript suivant:
- `CaisseStatisticsService` (caisse-statistics.service.ts)
- Tous les calculs et formats de données correspondent exactement
- Les noms des champs JSON respectent la convention camelCase du frontend

---

## 🚀 Améliorations futures possibles

1. **Prévisions avancées**
   - Utiliser des algorithmes de machine learning
   - Prendre en compte les tendances saisonnières
   - Analyser les jours de la semaine

2. **Analyse des catégories**
   - Comparer avec la période précédente pour déterminer la tendance réelle
   - Détecter les anomalies

3. **Cache**
   - Mettre en cache les statistiques fréquemment demandées
   - Invalider le cache lors de nouvelles transactions

4. **Performance**
   - Ajouter des index sur les colonnes `created_at`, `type_transaction`, `caisse_uuid`
   - Pagination pour les grandes quantités de données

---

## 📞 Support

Pour toute question ou problème, veuillez contacter l'équipe de développement.
