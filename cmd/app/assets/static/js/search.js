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

    const cards = items.map(item => {
        const title = item.name || item.title || 'Sem título';
        const year = item.year || '';
        const type = item.type ? item.type.toUpperCase() : '';
        // Usa poster (formato retrato) em vez de logo horizontal
        const posterUrl = item.poster || '';
        const placeholder = posterUrl
            ? `<img src="${posterUrl}" alt="${escapeHtml(title)}" loading="lazy" />`
            : `<div style="aspect-ratio:2/3; background:#1a1a1a; display:flex; align-items:center; justify-content:center; color:#555; font-size:2rem;">🎬</div>`;

        return `
            <div class="movie-card search-card">
                <a href="/detail/${item.type}/${item.id}">
                    ${placeholder}
                    <div class="card-info">
                        <h3>${escapeHtml(title)}</h3>
                        <p>${year} ${type}</p>
                    </div>
                </a>
            </div>
        `;
    }).join('');

    searchResults.innerHTML = `
        <div class="search-results-header">Resultados para "${escapeHtml(query)}"</div>
        <div class="catalog-grid search-results-grid">${cards}</div>
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
