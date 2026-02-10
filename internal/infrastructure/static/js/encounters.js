// Encounters Calculator JavaScript

// Patreon banner dismiss
function initPatreonBanner() {
    var banner = document.getElementById('patreon-banner');
    if (!banner) return;
    if (localStorage.getItem('patreon-banner-dismissed')) {
        banner.classList.add('patreon-banner--hidden');
        return;
    }
    var dismissBtn = banner.querySelector('.patreon-banner-dismiss');
    if (dismissBtn) {
        dismissBtn.addEventListener('click', function() {
            banner.classList.add('patreon-banner--hidden');
            localStorage.setItem('patreon-banner-dismissed', '1');
        });
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    initPatreonBanner();
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

// Monster selection tracking
window.selectedMonsters = [];

function addMonster(btn) {
    const row = btn.closest('.monster-row');
    const name = row.dataset.name;
    const xp = parseInt(row.dataset.xp, 10);
    const id = row.dataset.id;

    window.selectedMonsters.push({ id, name, xp });
    updateSelectedMonstersUI();
}

function removeMonster(index) {
    window.selectedMonsters.splice(index, 1);
    updateSelectedMonstersUI();
}

function updateSelectedMonstersUI() {
    const list = document.getElementById('selected-monsters-list');
    const countEl = document.getElementById('selected-count');
    const usedEl = document.getElementById('xp-used');
    const remainingEl = document.getElementById('xp-remaining');

    if (!list) return;

    const totalUsed = window.selectedMonsters.reduce((sum, m) => sum + m.xp, 0);
    const maxXPInput = document.querySelector('input[name="max_xp"]');
    const budget = maxXPInput ? parseInt(maxXPInput.value, 10) : 0;

    countEl.textContent = window.selectedMonsters.length;
    usedEl.textContent = window.encountersUtils.formatXP(totalUsed);

    const remaining = budget - totalUsed;
    remainingEl.innerHTML = 'Rimanenti: <strong>' + window.encountersUtils.formatXP(remaining) + '</strong>';
    remainingEl.classList.toggle('over-budget', remaining < 0);

    list.innerHTML = window.selectedMonsters.map((m, i) =>
        '<div class="selected-monster-item">' +
            '<span>' + m.name + ' (PE ' + window.encountersUtils.formatXP(m.xp) + ')</span>' +
            '<button type="button" onclick="removeMonster(' + i + ')">✕</button>' +
        '</div>'
    ).join('');
}

// Monster row accordion expand/collapse
document.addEventListener('click', function(evt) {
    // Don't toggle when clicking the + button
    if (evt.target.closest('.monster-add-btn')) return;

    var row = evt.target.closest('.monster-row');
    if (!row) return;

    var detailRow = row.nextElementSibling;
    if (!detailRow || !detailRow.classList.contains('monster-detail-row')) return;

    var isExpanded = detailRow.classList.contains('expanded');

    // Collapse all expanded rows (accordion behavior)
    document.querySelectorAll('.monster-detail-row.expanded').forEach(function(el) {
        el.classList.remove('expanded');
        el.previousElementSibling.classList.remove('expanded');
    });

    // Toggle clicked row (if it wasn't already expanded)
    if (!isExpanded) {
        detailRow.classList.add('expanded');
        row.classList.add('expanded');
    }
});

// Reset selected monsters when a new calculation is made
document.addEventListener('htmx:afterSwap', function(evt) {
    if (evt.detail.target && evt.detail.target.id === 'result-container') {
        window.selectedMonsters = [];
    }
});