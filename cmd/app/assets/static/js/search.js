// Cinemeta search integration via backend proxy
const CINEMETA_SEARCH_API = '/api/search?q=';

function normalizeSearchResults(data) {
    if (!data) return [];
    if (Array.isArray(data.results)) return data.results;
    if (Array.isArray(data.metas)) return data.metas;
    return [];
}

function setSearchLoading(isLoading) {
    const searchResults = document.getElementById('search-results');
    if (!searchResults) return;

    if (isLoading) {
        searchResults.innerHTML = '<div class="search-results-loading">Buscando...</div>';
        searchResults.style.display = 'block';
        searchResults.classList.add('active');
    } else {
        const loading = searchResults.querySelector('.search-results-loading');
        if (loading) {
            loading.remove();
        }
    }
}

function clearSearchResults() {
    const searchResults = document.getElementById('search-results');
    if (!searchResults) return;
    searchResults.innerHTML = '';
    searchResults.style.display = 'none';
    searchResults.classList.remove('active');
}

function escapeHtml(str) {
    return str.replace(/[&<>]/g, function(m) {
        if (m === '&') return '&amp;';
        if (m === '<') return '&lt;';
        if (m === '>') return '&gt;';
        return m;
    });
}

function renderSearchResults(results, query) {
    const searchResults = document.getElementById('search-results');
    if (!searchResults) return;

    const items = normalizeSearchResults(results);
    if (!items.length) {
        searchResults.innerHTML = `<div class="search-results-empty">Nenhum resultado encontrado para "${escapeHtml(query)}".</div>`;
        searchResults.style.display = 'block';
        searchResults.classList.add('active');
        return;
    }

    const listItems = items.map(item => {
        const title = item.name || item.title || 'Sem título';
        const year = item.year ? ` (${item.year})` : '';
        const subtitle = item.type ? item.type.toUpperCase() : '';
        
        // Prioriza logo, depois poster. Se nenhum existir, usa placeholder
        const imageUrl = item.logo || item.poster || '';
        const logoHtml = imageUrl 
            ? `<div class="search-result-logo-wrap"><img class="search-result-logo" src="${imageUrl}" alt="${escapeHtml(title)}" loading="lazy"/></div>`
            : `<div class="search-result-logo-wrap"><div class="search-result-placeholder">🎬</div></div>`;

        return `
            <a class="search-result-item" href="/detail/${item.type}/${item.id}">
                ${logoHtml}
                <div class="search-result-meta">
                    <span class="title">${escapeHtml(title)}${year}</span>
                    <span class="meta">${escapeHtml(subtitle)}</span>
                </div>
            </a>
        `;
    }).join('');

    searchResults.innerHTML = `
        <div class="search-results-header">Resultados para "${escapeHtml(query)}"</div>
        <div class="search-results-list">${listItems}</div>
    `;
    searchResults.style.display = 'block';
    searchResults.classList.add('active');
}

function toggleSearchBar() {
    const searchBar = document.getElementById('search-bar');
    const searchInput = document.getElementById('cinemeta-search-input');
    if (!searchBar || !searchInput) return;

    const isActive = searchBar.classList.toggle('active');
    if (isActive) {
        searchInput.focus();
        searchInput.select();
    } else {
        clearSearchResults();
    }
}

async function performSearch(query) {
    const trimmed = query.trim();
    if (!trimmed) {
        clearSearchResults();
        return;
    }

    setSearchLoading(true);
    try {
        const response = await fetch(CINEMETA_SEARCH_API + encodeURIComponent(trimmed));
        if (!response.ok) {
            throw new Error('Erro ao buscar');
        }
        const data = await response.json();
        renderSearchResults(data, trimmed);
    } catch (err) {
        const searchResults = document.getElementById('search-results');
        if (!searchResults) return;
        searchResults.innerHTML = `<div class="search-results-error">Erro ao buscar. Tente novamente.</div>`;
        searchResults.style.display = 'block';
        searchResults.classList.add('active');
    } finally {
        setSearchLoading(false);
    }
}

function initSearchUI() {
    const toggleButton = document.getElementById('search-toggle');
    const searchInput = document.getElementById('cinemeta-search-input');

    if (toggleButton) {
        toggleButton.addEventListener('click', function (event) {
            event.stopPropagation();
            toggleSearchBar();
        });
    }

    if (searchInput) {
        searchInput.addEventListener('keydown', function (event) {
            if (event.key === 'Enter') {
                event.preventDefault();
                performSearch(searchInput.value);
            }
        });
    }

    document.addEventListener('click', function (event) {
        const searchBar = document.getElementById('search-bar');
        const toggleButton = document.getElementById('search-toggle');
        if (!searchBar || !toggleButton) return;

        if (!searchBar.contains(event.target) && event.target !== toggleButton) {
            searchBar.classList.remove('active');
        }
    });
}

window.addEventListener('DOMContentLoaded', initSearchUI);
