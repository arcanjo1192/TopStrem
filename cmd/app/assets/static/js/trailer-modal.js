(function() {
    const modal = document.getElementById('trailer-modal');
    const iframe = document.getElementById('trailer-iframe');
    const closeBtn = document.querySelector('.modal-close');
    const trailerButton = document.querySelector('.trailer-button');

    if (trailerButton) {
        trailerButton.addEventListener('click', function() {
            const trailerUrl = this.getAttribute('data-trailer-url');
            if (trailerUrl) {
                iframe.src = trailerUrl;
                modal.style.display = 'flex';
            }
        });
    }

    function closeModal() {
        modal.style.display = 'none';
        iframe.src = '';
    }

    if (closeBtn) closeBtn.addEventListener('click', closeModal);
    if (modal) modal.addEventListener('click', function(e) {
        if (e.target === modal) closeModal();
    });
})();