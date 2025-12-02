// Encounters Calculator JavaScript

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Configure HTMX
    htmx.config.requestClass = 'loading';
    htmx.config.historyEnabled = true;

    // Add loading indicators
    document.addEventListener('htmx:beforeRequest', function(evt) {
        const target = evt.target;
        if (target.classList.contains('btn')) {
            target.style.opacity = '0.7';
            target.style.pointerEvents = 'none';
            target.textContent = 'Caricamento...';
        }
    });

    document.addEventListener('htmx:afterRequest', function(evt) {
        const target = evt.target;
        if (target.classList.contains('btn')) {
            target.style.opacity = '1';
            target.style.pointerEvents = 'auto';
        }
    });

    // Handle errors gracefully
    document.addEventListener('htmx:responseError', function(evt) {
        console.error('HTMX Request failed:', evt.detail);
        showNotification('Errore nel caricamento. Riprova.', 'error');
    });
});

// Utility functions
window.encountersUtils = {
    // Show notifications
    showNotification: function(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;

        // Add styles
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 1rem 1.5rem;
            border-radius: 0.5rem;
            color: white;
            z-index: 1000;
            font-weight: 500;
            background: ${type === 'error' ? '#EF4444' : type === 'success' ? '#22C55E' : '#3B82F6'};
            box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
        `;

        document.body.appendChild(notification);

        // Auto remove after 3 seconds
        setTimeout(() => {
            if (notification.parentElement) {
                notification.parentElement.removeChild(notification);
            }
        }, 3000);
    },

    // Validate party levels
    validatePartyLevels: function(levels) {
        if (!levels || levels.length === 0) {
            return 'Inserisci almeno un livello';
        }

        for (let level of levels) {
            if (level < 1 || level > 20) {
                return 'I livelli devono essere tra 1 e 20';
            }
        }

        return null;
    },

    // Format XP numbers
    formatXP: function(xp) {
        return new Intl.NumberFormat('it-IT').format(xp);
    },

    // Get difficulty color class
    getDifficultyClass: function(difficulty) {
        switch (difficulty.toLowerCase()) {
            case 'easy': case 'facile': return 'difficulty-easy';
            case 'medium': case 'medio': return 'difficulty-medium';
            case 'hard': case 'difficile': return 'difficulty-hard';
            case 'deadly': case 'letale': return 'difficulty-deadly';
            default: return '';
        }
    }
};

// Make showNotification available globally
window.showNotification = window.encountersUtils.showNotification;