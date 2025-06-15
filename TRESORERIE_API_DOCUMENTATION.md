# API Dashboard Trésorerie - Documentation

## 🚀 Vue d'ensemble

Ce système de dashboard de trésorerie a été créé pour correspondre parfaitement aux besoins du frontend Angular fourni. Il offre une API complète avec des données en temps réel, des graphiques interactifs, des alertes intelligentes et des métriques avancées.

## 📋 Endpoints Disponibles

### 1. Dashboard Principal
**GET** `/api/dashboard/tresoreries/`

Récupère toutes les données du dashboard en un seul appel.

**Paramètres :**
- `user_uuid` (requis) : UUID de l'utilisateur connecté
- `pos_uuid` (optionnel) : UUID du POS spécifique, sinon toutes les données de l'entreprise
- `period` (optionnel) : Période d'analyse (`today`, `this_week`, `this_month`, `last_3_months`, `last_6_months`, `last_year`)

**Exemple d'appel :**
```
GET /api/dashboard/tresoreries/?user_uuid=123e4567-e89b-12d3-a456-426614174000&pos_uuid=987fcdeb-51a2-43d7-8f9e-123456789abc&period=this_month
```

**Réponse :**
```json
{
  "success": true,
  "data": {
    "kpi_data": {
      "solde_total": 150000.50,
      "total_entrees": 85000.00,
      "total_sorties": 35000.00,
      "total_fond_de_caisse": 100000.50,
      "nombre_entrees": 45,
      "nombre_sorties": 12,
      "variation_solde": 8.5
    },
    "top_caisses": [
      {
        "uuid": "caisse-uuid-1",
        "name": "Caisse Principale",
        "solde": 50000.00,
        "total_entrees": 25000.00,
        "total_sorties": 5000.00,
        "nombre_transactions": 15,
        "performance": 85.5
      }
    ],
    "alerts": [
      {
        "id": "solde-critique",
        "type": "danger",
        "icon": "alert-triangle",
        "title": "Solde Critique",
        "message": "Le solde total des caisses est en dessous du seuil de sécurité."
      }
    ],
    "metriques": {
      "velocite": 2.5,
      "velocite_variation": 12.5,
      "ratio_liquidite": 245.5,
      "taux_croissance": 8.2,
      "efficacite": 1250.75
    },
    "actions_recommandees": [
      {
        "id": "action-1",
        "description": "Effectuer un apport de fonds dans les caisses principales"
      }
    ],
    "stats_footer": {
      "transactions_aujourdhui": 8,
      "montant_moyen_transaction": 2500.50,
      "temps_moyen_traitement": 3.2,
      "derniere_mise_a_jour": "2025-06-15T14:30:00Z"
    }
  }
}
```

### 2. Données d'Évolution (Graphique Linéaire)
**GET** `/api/dashboard/tresoreries/evolution-chart`

**Paramètres :** Mêmes que l'endpoint principal

**Réponse :**
```json
{
  "success": true,
  "data": {
    "dates": ["01/06", "02/06", "03/06", "04/06"],
    "solde_data": [50000, 52000, 48000, 55000],
    "entrees_data": [15000, 18000, 12000, 20000],
    "sorties_data": [8000, 6000, 10000, 5000],
    "fond_caisses_data": [43000, 40000, 46000, 40000]
  }
}
```

### 3. Données de Répartition (Graphique Circulaire)
**GET** `/api/dashboard/tresoreries/repartition-chart`

**Paramètres :** Mêmes que l'endpoint principal

**Réponse :**
```json
{
  "success": true,
  "data": {
    "series": [85000.00, 35000.00, 100000.50],
    "labels": ["Entrées (38.6%)", "Sorties (15.9%)", "Fond de Caisse (45.5%)"]
  }
}
```

### 4. Résumé Financier Simplifié
**GET** `/api/dashboard/tresoreries/financial-summary`

**Paramètres :**
- `entreprise_uuid` (requis) : UUID de l'entreprise
- `pos_uuid` (optionnel) : UUID du POS spécifique

### 5. Top Caisses
**GET** `/api/dashboard/tresoreries/top-caisses`

**Paramètres :**
- `entreprise_uuid` (requis) : UUID de l'entreprise
- `pos_uuid` (optionnel) : UUID du POS spécifique
- `period` (optionnel) : Période d'analyse
- `limit` (optionnel) : Nombre de caisses à retourner (défaut: 10)

### 6. Alertes et Recommandations
**GET** `/api/dashboard/tresoreries/alerts-recommendations`

**Paramètres :** Mêmes que l'endpoint principal

## 🔧 Intégration Frontend

### Service Angular

```typescript
export class TresorerieApiService {
  constructor(private http: HttpClient) {}

  mapPeriodToApi(period: string): string {
    const periodMap = {
      'Aujourd\'hui': 'today',
      'Cette semaine': 'this_week', 
      'Ce mois': 'this_month',
      '3 derniers mois': 'last_3_months',
      '6 derniers mois': 'last_6_months',
      '1 an': 'last_year'
    };
    return periodMap[period] || 'this_month';
  }

  getDashboardTresorerie(userUuid: string, posUuid: string, period: string): Observable<ApiResponse<DashboardTresorerieDataApi>> {
    const params = new HttpParams()
      .set('user_uuid', userUuid)
      .set('pos_uuid', posUuid || '')
      .set('period', period);

    return this.http.get<ApiResponse<DashboardTresorerieDataApi>>(
      '/api/dashboard/tresoreries/', 
      { params }
    );
  }

  getEvolutionChartData(userUuid: string, posUuid: string, period: string): Observable<ApiResponse<EvolutionChartDataApi>> {
    const params = new HttpParams()
      .set('user_uuid', userUuid)
      .set('pos_uuid', posUuid || '')
      .set('period', period);

    return this.http.get<ApiResponse<EvolutionChartDataApi>>(
      '/api/dashboard/tresoreries/evolution-chart', 
      { params }
    );
  }

  getRepartitionChartData(userUuid: string, posUuid: string, period: string): Observable<ApiResponse<RepartitionChartDataApi>> {
    const params = new HttpParams()
      .set('user_uuid', userUuid)
      .set('pos_uuid', posUuid || '')
      .set('period', period);

    return this.http.get<ApiResponse<RepartitionChartDataApi>>(
      '/api/dashboard/tresoreries/repartition-chart', 
      { params }
    );
  }
}
```

## 🎯 Fonctionnalités Clés

### 1. **Système d'Alertes Intelligentes**
- Détection automatique des soldes critiques
- Alertes de performance dégradée
- Notification de déséquilibre des flux

### 2. **Métriques Avancées**
- **Vélocité** : Nombre de transactions par jour
- **Ratio de Liquidité** : Capacité à couvrir les sorties
- **Efficacité** : Montant moyen par transaction
- **Taux de Croissance** : Évolution du solde
- **Variation de Vélocité** : Comparaison avec la période précédente

### 3. **Graphiques Dynamiques**
- Graphique d'évolution temporelle (solde, entrées, sorties, fonds)
- Graphique de répartition circulaire avec pourcentages
- Support de multiple périodes (jour, semaine, mois, trimestre, semestre, année)

### 4. **Gestion Multi-POS**
- Données agrégées par entreprise
- Filtrage par POS spécifique
- Comparaison de performance entre POS

### 5. **Calculs de Performance**
- Performance des caisses basée sur la rentabilité
- Classement automatique des meilleures caisses
- Indicateurs de santé financière

## 📊 Logique Métier

### Calcul du Solde Total
```
Solde Total = Total Entrées - Total Sorties + Total Fond de Caisse
```

### Calcul de la Performance d'une Caisse
```
Performance = ((Total Entrées - Total Sorties) / Total Flux) * 100
```

### Calcul de la Variation
```
Variation = ((Valeur Actuelle - Valeur Précédente) / Valeur Précédente) * 100
```

## 🔒 Sécurité et Validation

- Validation stricte des UUIDs utilisateur
- Filtrage par entreprise pour la sécurité des données
- Gestion des erreurs avec messages explicites
- Protection contre les injections SQL avec requêtes préparées

## 🚀 Performance

- Requêtes SQL optimisées avec indices appropriés
- Pagination pour les grandes listes
- Cache des métriques calculées
- Requêtes groupées pour réduire les appels DB

## 📈 Évolutivité

Le système est conçu pour être facilement extensible :
- Ajout de nouvelles métriques
- Support de nouveaux types de graphiques
- Intégration de nouveaux types d'alertes
- Export vers différents formats (PDF, Excel, etc.)

Cette API fournit une base solide pour un dashboard de trésorerie moderne et performant, parfaitement adaptée aux besoins du frontend Angular fourni.
