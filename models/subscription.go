package models

import (
	"time"

	"gorm.io/gorm"
)

type SubscriptionStatus string

const (
	StatusPending           SubscriptionStatus = "pending"
	StatusPaymentPending    SubscriptionStatus = "payment_pending"
	StatusPaymentFailed     SubscriptionStatus = "payment_failed"
	StatusValidationPending SubscriptionStatus = "validation_pending"
	StatusActive            SubscriptionStatus = "active"
	StatusSuspended         SubscriptionStatus = "suspended"
	StatusCancelled         SubscriptionStatus = "cancelled"
	StatusExpired           SubscriptionStatus = "expired"
	StatusRejected          SubscriptionStatus = "rejected"
)

type Subscription struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Status et progression
	Step   int                `gorm:"not null;default:1" json:"step"` // Étape du processus (1-3)
	Status SubscriptionStatus `gorm:"not null;default:pending" json:"status"`

	// Information de l'entreprise
	TypeEntreprise string `gorm:"not null" json:"type_entreprise"` // PME, GE, Particulier
	Name           string `gorm:"not null" json:"name"`
	Rccm           string `json:"rccm"`
	IdNat          string `json:"idnat"`
	NImpot         string `json:"nimpot"`
	Email          string `json:"email"`                     // Email officiel
	Telephone      string `gorm:"not null" json:"telephone"` // Telephone officiel
	Manager        string `gorm:"not null" json:"manager"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Country        string `json:"country"`

	// Information de payment
	Planid    string  `gorm:"not null" json:"planid"`               // ID du plan d'abonnement
	PlanName  string  `gorm:"not null" json:"plan_name"`            // Nom du plan d'abonnement
	Amount    float64 `gorm:"not null" json:"amount"`               // Montant de l'abonnement
	Currency  string  `gorm:"not null;default:USD" json:"currency"` // Devise de l'abonnement, default USD
	Duration  int     `gorm:"not null" json:"duration"`             // Durée de l'abonnement en mois
	Features  string  `gorm:"type:text" json:"features"`            // JSON des fonctionnalités du plan
	MaxUsers  int     `json:"max_users"`                            // Nombre maximum d'utilisateurs
	MaxPOS    int     `json:"max_pos"`                              // Nombre maximum de points de vente
	StorageGB int     `json:"storage_gb"`                           // Stockage en GB

	// Information de paiement
	PaymentMethod    string    `json:"payment_method"`    // Méthode de paiement
	PaymentReference string    `json:"payment_reference"` // Référence de paiement
	PaymentStatus    string    `json:"payment_status"`    // Statut du paiement
	TransactionID    string    `json:"transaction_id"`    // ID de transaction
	PaymentDate      time.Time `json:"payment_date"`      // Date de paiement
	MobileOperator   string    `json:"mobile_operator"`   // Opérateur mobile si applicable

	// Information de validation
	ValidateBy     string    `json:"validate_by"`     // UUID de l'utilisateur qui a validé l'abonnement
	ValidationDate time.Time `json:"validation_date"` // Date de validation de l'abonnement
	Notes          string    `json:"notes"`           // Notes sur l'abonnement
	RejectedReason string    `json:"rejected_reason"` // Raison du rejet de l'abonnement

	// Dates importantes
	StartDate time.Time `json:"start_date"` // Date de début d'abonnement
	EndDate   time.Time `json:"end_date"`   // Date de fin d'abonnement
	ExpiresAt time.Time `json:"expires_at"` // Date d'expiration

	// Renouvellement
	AutoRenewal     bool      `json:"auto_renewal"`      // Renouvellement automatique
	NextBillingDate time.Time `json:"next_billing_date"` // Prochaine date de facturation
}

// SubscriptionPlan représente un plan d'abonnement
type SubscriptionPlan struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Currency    string   `json:"currency"`
	Duration    int      `json:"duration"` // en mois
	MaxUsers    int      `json:"max_users"`
	MaxPOS      int      `json:"max_pos"`
	StorageGB   int      `json:"storage_gb"`
	Features    []string `json:"features"`
	Popular     bool     `json:"popular"`
}

// PaymentMethod représente une méthode de paiement
type PaymentMethod struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Icon         string `json:"icon"`
	Instructions string `json:"instructions"`
}

// MobileOperator représente un opérateur mobile
type MobileOperator struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Instructions string `json:"instructions"`
}

// SubscriptionHistory représente l'historique des abonnements
type SubscriptionHistory struct {
	UUID             string    `gorm:"type:varchar(255);primary_key" json:"uuid"`
	SubscriptionUUID string    `gorm:"not null" json:"subscription_uuid"`
	Action           string    `gorm:"not null" json:"action"` // created, upgraded, downgraded, renewed, suspended, cancelled, reactivated
	OldPlan          string    `json:"old_plan"`
	NewPlan          string    `json:"new_plan"`
	Amount           float64   `json:"amount"`
	ActionDate       time.Time `gorm:"not null" json:"action_date"`
	UserUUID         string    `gorm:"not null" json:"user_uuid"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// SubscriptionNotification représente les notifications d'abonnement
type SubscriptionNotification struct {
	UUID             string    `gorm:"type:varchar(255);primary_key" json:"uuid"`
	SubscriptionUUID string    `gorm:"not null" json:"subscription_uuid"`
	Type             string    `gorm:"not null" json:"type"` // expiration_warning, payment_failed, upgrade_available, trial_ending, suspended, renewed
	Title            string    `gorm:"not null" json:"title"`
	Message          string    `gorm:"not null" json:"message"`
	SendDate         time.Time `json:"send_date"`
	Read             bool      `gorm:"default:false" json:"read"`
	Priority         string    `gorm:"not null" json:"priority"` // low, medium, high, urgent
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// SubscriptionMetrics représente les métriques d'abonnement
type SubscriptionMetrics struct {
	TotalSubscriptions   int     `json:"total_subscriptions"`
	ActiveSubscriptions  int     `json:"active_subscriptions"`
	TrialSubscriptions   int     `json:"trial_subscriptions"`
	ExpiredSubscriptions int     `json:"expired_subscriptions"`
	MonthlyRevenue       float64 `json:"monthly_revenue"`
	AnnualRevenue        float64 `json:"annual_revenue"`
	TrialConversionRate  float64 `json:"trial_conversion_rate"`
	RetentionRate        float64 `json:"retention_rate"`
	ChurnRate            float64 `json:"churn_rate"`
	ARPU                 float64 `json:"arpu"` // Average Revenue Per User
	LTV                  float64 `json:"ltv"`  // Lifetime Value
}

// Promotion représente les codes promotionnels
type Promotion struct {
	UUID          string    `gorm:"type:varchar(255);primary_key" json:"uuid"`
	Code          string    `gorm:"unique;not null" json:"code"`
	Description   string    `gorm:"not null" json:"description"`
	Type          string    `gorm:"not null" json:"type"` // percentage, fixed_amount, free_trial_extension
	Value         float64   `gorm:"not null" json:"value"`
	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null" json:"end_date"`
	MaxUsages     int       `json:"max_usages"`
	CurrentUsages int       `gorm:"default:0" json:"current_usages"`
	EligiblePlans string    `gorm:"type:text" json:"eligible_plans"` // JSON array
	Conditions    string    `json:"conditions"`
	Active        bool      `gorm:"default:true" json:"active"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
