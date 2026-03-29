const API_URL = 'http://localhost:8080/api/verify';

// --- DOM Elements ---
const verifyBtn    = document.getElementById('verify-btn');
const claimInput   = document.getElementById('claim-input');
const charCount    = document.getElementById('char-count');
const resultSection = document.getElementById('result-section');

// Result Elements
const verdictHeader = document.getElementById('verdict-header');
const verdictIcon   = document.getElementById('verdict-icon');
const verdictBadge  = document.getElementById('verdict-badge');
const verdictTitle  = document.getElementById('verdict-title');
const reasoningText = document.getElementById('reasoning-text');
const explanationText = document.getElementById('explanation-text');
const sourcesList   = document.getElementById('sources-list');
const toxicityBar   = document.getElementById('toxicity-bar');
const spamBar       = document.getElementById('spam-bar');
const toxicityVal   = document.getElementById('toxicity-val');
const spamVal       = document.getElementById('spam-val');
const toxPct        = document.getElementById('tox-pct');
const spamPct       = document.getElementById('spam-pct');

// --- Character Count ---
claimInput.addEventListener('input', () => {
    charCount.textContent = claimInput.value.length;
});

// Enter key to submit (Ctrl+Enter or Shift+Enter for newline)
claimInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && !e.shiftKey && !e.ctrlKey) {
        e.preventDefault();
        verifyBtn.click();
    }
});

// --- Verify Button ---
verifyBtn.addEventListener('click', async () => {
    const text = claimInput.value.trim();
    if (!text) {
        claimInput.focus();
        claimInput.style.borderColor = '#ff4d6d';
        setTimeout(() => { claimInput.style.borderColor = ''; }, 1500);
        return;
    }

    setLoading(true);
    resultSection.classList.add('hidden');

    try {
        const response = await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ text })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || `Server error: ${response.status}`);
        }

        renderResult(data);

    } catch (error) {
        showError(error.message);
    } finally {
        setLoading(false);
    }
});

// --- Loading State ---
function setLoading(loading) {
    const btnContent = verifyBtn.querySelector('.btn-content');
    const btnLoader  = verifyBtn.querySelector('.btn-loader');
    verifyBtn.disabled = loading;
    btnContent.classList.toggle('hidden', loading);
    btnLoader.classList.toggle('hidden', !loading);
}

// --- Render Results ---
function renderResult(data) {
    const status = (data.status || '').toUpperCase();

    // Determine verdict config
    let cls, icon, label, title;

    if (status === 'TRUE') {
        cls   = 'true';
        icon  = '✅';
        label = '✓ Verified True';
        title = 'Statement Confirmed as Accurate';
    } else if (status === 'FALSE') {
        cls   = 'false';
        icon  = '❌';
        label = '✗ Verified False';
        title = 'Factual Inaccuracy Detected';
    } else {
        cls   = 'partial';
        icon  = '⚠️';
        label = '~ Partially True';
        title = 'Statement Contains Mixed Accuracy';
    }

    // Set verdict icon
    verdictIcon.textContent = icon;
    verdictIcon.className = `verdict-icon ${cls}`;

    // Set badge
    verdictBadge.className = `verdict-badge ${cls}`;
    verdictBadge.textContent = label;

    // Set title
    verdictTitle.textContent = title;

    // Set reasoning
    reasoningText.textContent = data.reasoning || 'No summary available.';

    // Set detailed explanation
    explanationText.textContent = data.detailed_explanation || 'No detailed analysis available.';

    // Scores
    const toxScore  = Math.round((data.toxicity_score || 0) * 100);
    const spamScore = Math.round((data.spam_score || 0) * 100);

    toxicityVal.textContent = `${toxScore}%`;
    spamVal.textContent     = `${spamScore}%`;
    toxPct.textContent      = `${toxScore}%`;
    spamPct.textContent     = `${spamScore}%`;

    // Score colors
    toxicityVal.style.color = scoreColor(toxScore);
    spamVal.style.color     = scoreColor(spamScore);

    // Animate bars (delay so transition fires)
    setTimeout(() => {
        toxicityBar.style.width = `${toxScore}%`;
        spamBar.style.width     = `${spamScore}%`;
    }, 100);

    // Sources
    sourcesList.innerHTML = '';
    const sources = data.sources && data.sources.length > 0
        ? data.sources
        : ['General AI Knowledge Base'];

    sources.forEach((src, i) => {
        const item = document.createElement('div');
        item.className = 'source-item';
        item.innerHTML = `
            <div class="source-number">${i + 1}</div>
            <span>${src}</span>
        `;
        sourcesList.appendChild(item);
    });

    // Show result section
    resultSection.classList.remove('hidden');
    resultSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

// Returns color based on score value
function scoreColor(pct) {
    if (pct < 20) return '#10d98a';
    if (pct < 50) return '#f59e0b';
    return '#ff4d6d';
}

// --- Error Display ---
function showError(message) {
    // Show error in result section
    verdictIcon.textContent = '⚡';
    verdictIcon.className = 'verdict-icon false';
    verdictBadge.className = 'verdict-badge false';
    verdictBadge.textContent = '✗ Audit Failed';
    verdictTitle.textContent = 'Error Processing Request';
    reasoningText.textContent = message;
    explanationText.textContent = 'Please check that the backend server is running on port 8080 and your API key is valid.';
    sourcesList.innerHTML = '';
    toxicityBar.style.width = '0%';
    spamBar.style.width = '0%';
    toxicityVal.textContent = '—';
    spamVal.textContent = '—';
    toxPct.textContent = '—';
    spamPct.textContent = '—';

    resultSection.classList.remove('hidden');
    resultSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
}
