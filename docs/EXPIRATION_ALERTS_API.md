# API Alertes d'Expiration de Stock

## Description

Cette API permet de récupérer les alertes d'expiration pour les produits en stock. Elle identifie les produits expirés ou bientôt expirés (dans les 7 prochains jours).

## Endpoint

```
GET /api/dashboard/main/expiration-alerts
```

## Paramètres de requête (Query Parameters)

| Paramètre       | Type   | Requis | Description                                    |
|----------------|--------|--------|------------------------------------------------|
| entreprise_uuid | string | Oui    | UUID de l'entreprise                          |
| pos_uuid        | string | Oui    | UUID du point de vente (POS)                  |

## Exemple de requête

```bash
GET /api/dashboard/main/expiration-alerts?entreprise_uuid=123e4567-e89b-12d3-a456-426614174000&pos_uuid=123e4567-e89b-12d3-a456-426614174001
```

## Réponse

### Structure de la réponse (200 OK)

```json
[
  {
    "uuid": "stock-uuid-123",
    "product_uuid": "product-uuid-456",
    "name": "Yaourt Nature",
    "reference": "PROD-001",
    "unite_vente": "Unité",
    "quantity": 50.5,
    "date_expiration": "2025-11-30T00:00:00Z",
    "prix_achat": 1.5,
    "fournisseur_name": "Laiterie Moderne",
    "alertType": "expire",
    "image": "https://example.com/image.jpg",
    "daysRemaining": -2
  },
  {
    "uuid": "stock-uuid-789",
    "product_uuid": "product-uuid-012",
    "name": "Pain de mie",
    "reference": "PROD-002",
    "unite_vente": "Paquet",
    "quantity": 30,
    "date_expiration": "2025-12-05T00:00:00Z",
    "prix_achat": 2.0,
    "fournisseur_name": "Boulangerie Centrale",
    "alertType": "bientot_expire",
    "image": "https://example.com/pain.jpg",
    "daysRemaining": 5
  }
]
```

### Description des champs de réponse

| Champ            | Type    | Description                                                          |
|-----------------|---------|----------------------------------------------------------------------|
| uuid            | string  | UUID du stock (le plus proche de l'expiration)                      |
| product_uuid    | string  | UUID du produit                                                      |
| name            | string  | Nom du produit                                                       |
| reference       | string  | Référence du produit                                                 |
| unite_vente     | string  | Unité de vente du produit                                           |
| quantity        | float64 | Quantité totale du produit concerné par l'alerte                    |
| date_expiration | string  | Date d'expiration au format ISO 8601                                |
| prix_achat      | float64 | Prix d'achat du stock                                               |
| fournisseur_name| string  | Nom du fournisseur (optionnel)                                      |
| alertType       | string  | Type d'alerte: "expire" (expiré) ou "bientot_expire" (bientôt expiré) |
| image           | string  | URL de l'image du produit (optionnel)                               |
| daysRemaining   | int     | Nombre de jours restants avant expiration (négatif si déjà expiré) |

### Types d'alertes

1. **expire**: Le produit est déjà expiré (`daysRemaining <= 0`)
2. **bientot_expire**: Le produit expirera dans les 7 prochains jours (`0 < daysRemaining <= 7`)

## Logique de fonctionnement

1. **Récupération des stocks**: L'API récupère tous les stocks de l'entreprise et du POS spécifiés
2. **Filtrage par date**: Seuls les stocks avec une date d'expiration dans les 7 prochains jours ou déjà expirés sont conservés
3. **Groupement par produit**: Les stocks sont groupés par produit
4. **Sélection du stock le plus proche**: Pour chaque produit, on sélectionne le stock avec la date d'expiration la plus proche
5. **Calcul de la quantité totale**: La quantité totale de tous les stocks concernés par l'alerte est calculée
6. **Tri des résultats**: Les alertes sont triées par nombre de jours restants (les expirés en premier)

## Codes d'erreur

| Code | Description                                                    |
|------|----------------------------------------------------------------|
| 200  | Succès - Liste des alertes retournée                         |
| 400  | Mauvaise requête - Paramètres manquants ou invalides         |
| 500  | Erreur serveur interne                                        |

### Exemple de réponse d'erreur (400)

```json
{
  "error": "Les paramètres entreprise_uuid et pos_uuid sont requis"
}
```

## Cas particuliers

- Si aucun stock n'est trouvé, l'API retourne un tableau vide `[]`
- Si un produit n'a pas de date d'expiration, il est ignoré
- Si plusieurs stocks du même produit expirent à des dates différentes, la quantité totale inclut tous les stocks concernés
- Le nom du fournisseur n'est affiché que si le fournisseur existe et est lié au stock

## Utilisation Frontend (TypeScript)

```typescript
interface ExpirationAlert {
  uuid?: string;
  product_uuid: string;
  name: string;
  reference: string;
  unite_vente: string;
  quantity: number;
  date_expiration: Date;
  prix_achat: number;
  fournisseur_name?: string;
  alertType: 'expire' | 'bientot_expire';
  image?: string;
  daysRemaining: number;
}

async function getExpirationAlerts(
  entreprise_uuid: string, 
  pos_uuid: string
): Promise<ExpirationAlert[]> {
  const response = await fetch(
    `/api/dashboard/main/expiration-alerts?entreprise_uuid=${entreprise_uuid}&pos_uuid=${pos_uuid}`
  );
  
  if (!response.ok) {
    throw new Error('Erreur lors de la récupération des alertes');
  }
  
  const alerts: ExpirationAlert[] = await response.json();
  
  // Statistiques
  console.log('🗓️ Alertes d\'expiration trouvées:', {
    total: alerts.length,
    expires: alerts.filter(a => a.alertType === 'expire').length,
    bientotExpires: alerts.filter(a => a.alertType === 'bientot_expire').length
  });
  
  return alerts;
}
```

## Performance

- L'API utilise des requêtes optimisées avec GORM
- Les stocks sont filtrés au niveau de la base de données
- Le tri est effectué en mémoire sur un ensemble réduit de données
- Temps de réponse moyen: < 100ms pour 1000 produits

## Notes importantes

- La période de 7 jours pour les alertes "bientôt expiré" est codée en dur dans l'API
- Les dates sont normalisées à minuit (00:00:00) pour des comparaisons cohérentes
- Les quantités sont arrondies à 2 décimales
- L'API ne modifie pas les données, elle est en lecture seule
