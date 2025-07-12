# Système de Gestion des Abonnements - ipos-stock-api

Ce document décrit le système complet de gestion des abonnements intégré dans l'API ipos-stock.

## Architecture

Le système est divisé en 3 contrôleurs principaux :

### 1. Subscriptions Controller (`controllers/subscriptions/subscription.controller.go`)
Gère le processus de souscription en 3 étapes :

**Étape 1 - Informations de l'entreprise :**
- `POST /api/subscriptions/create`
- Collecte : nom, email, téléphone, type d'entreprise, RCCM, etc.

**Étape 2 - Sélection du plan :**
- `PUT /api/subscriptions/:id/plan`
- Choix du plan, durée, calcul du prix

**Étape 3 - Paiement :**
- `PUT /api/subscriptions/:id/payment` (initiation)
- `PUT /api/subscriptions/:id/confirm` (confirmation)

**Autres endpoints :**
- `GET /api/subscriptions/plans` - Plans disponibles
- `GET /api/subscriptions/:id` - Détails d'un abonnement

### 2. Subscription Admin Controller (`controllers/subscriptions/subscription-admin.controller.go`)
Interface d'administration pour gérer tous les abonnements :

- `GET /api/subscriptions/admin/` - Liste paginée des abonnements
- `GET /api/subscriptions/admin/stats` - Statistiques globales
- `GET /api/subscriptions/admin/metrics` - Métriques détaillées
- `PUT /api/subscriptions/admin/:id/approve` - Approuver un abonnement
- `PUT /api/subscriptions/admin/:id/reject` - Rejeter un abonnement
- `PUT /api/subscriptions/admin/:id/suspend` - Suspendre un abonnement
- `PUT /api/subscriptions/admin/:id` - Mise à jour administrative

### 3. Abonnements Controller (`controllers/abonnements/abonnement.controller.go`)
Interface utilisateur pour gérer ses propres abonnements :

- `GET /api/abonnements/my` - Mes abonnements
- `GET /api/abonnements/current` - Abonnement actuel
- `GET /api/abonnements/stats` - Mes statistiques
- `PUT /api/abonnements/:id/renew` - Renouveler
- `PUT /api/abonnements/:id/cancel` - Annuler
- `GET /api/abonnements/:id/history` - Historique

## Modèles de Données

### Subscription (`models/subscription.go`)
Structure principale contenant :
- Informations entreprise (nom, RCCM, email, etc.)
- Plan sélectionné (ID, nom, durée, prix)
- Statut (pending, active, expired, etc.)
- Informations de paiement
- Dates importantes (début, fin, expiration)

### SubscriptionHistory
Historique des actions sur un abonnement :
- Actions : created, approved, rejected, suspended, renewed, etc.
- Dates et détails de chaque action

### SubscriptionPlan
Plans d'abonnement disponibles :
- Basic (29.99 USD/mois)
- Professional (59.99 USD/mois)
- Enterprise (99.99 USD/mois)

## Statuts des Abonnements

```go
const (
    StatusPending           = "pending"            // En attente
    StatusPaymentPending    = "payment_pending"    // Paiement en attente
    StatusPaymentFailed     = "payment_failed"     // Paiement échoué
    StatusValidationPending = "validation_pending" // Validation admin en attente
    StatusActive            = "active"             // Actif
    StatusSuspended         = "suspended"          // Suspendu
    StatusCancelled         = "cancelled"          // Annulé
    StatusExpired           = "expired"            // Expiré
    StatusRejected          = "rejected"           // Rejeté
)
```

## Flux de Souscription

1. **Création** : L'utilisateur soumet les informations de son entreprise
2. **Sélection plan** : Choix du plan et de la durée
3. **Paiement** : Traitement du paiement (mobile money, carte, etc.)
4. **Validation admin** : Un administrateur approuve ou rejette l'abonnement
5. **Activation** : L'abonnement devient actif

## Fonctionnalités Avancées

- **Renouvellement automatique** : Possibilité d'activer le renouvellement auto
- **Historique complet** : Traçabilité de toutes les actions
- **Métriques** : Statistiques détaillées pour les admins
- **Remises** : Remises automatiques pour les abonnements longs (6+ mois)
- **Multi-devises** : Support de différentes devises
- **Notifications** : Système de notification pour les expirations

## Configuration

Les modèles sont automatiquement migrés via GORM dans `database/connection.go`.

## Tests

Pour tester le système :

1. Compiler : `go build .`
2. Lancer : `./ipos-stock-api.exe`
3. Utiliser les endpoints via Postman ou l'interface Angular

## Sécurité

- Authentification requise pour la plupart des endpoints
- Validation des données avec le package `utils/validateStruct.go`
- Middleware d'autorisation pour les fonctions admin
